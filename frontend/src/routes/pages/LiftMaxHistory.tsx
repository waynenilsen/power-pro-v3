import { Link, useParams } from 'react-router-dom';
import { Container } from '../../components/layout';
import { useAuth } from '../../contexts/useAuth';
import { useLiftMaxes, useLifts } from '../../hooks';
import {
  Loader2,
  ChevronLeft,
  Dumbbell,
  AlertCircle,
  TrendingUp,
  TrendingDown,
  Minus,
  Calendar,
} from 'lucide-react';
import type { LiftMax, MaxType } from '../../api/types';

function formatDate(dateString: string): string {
  // Parse the date as UTC to avoid timezone issues
  // The backend stores dates at midnight UTC, so we extract just the date part
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  });
}

function MaxTypeBadge({ type }: { type: MaxType }) {
  const isOneRM = type === 'ONE_RM';
  return (
    <span
      className={`
        inline-flex items-center px-2 py-0.5
        text-xs font-bold uppercase tracking-wider
        rounded-md border
        ${isOneRM
          ? 'bg-accent/15 text-accent border-accent/30'
          : 'bg-success/15 text-success border-success/30'
        }
      `}
    >
      {isOneRM ? '1RM' : 'TM'}
    </span>
  );
}

interface HistoryEntryProps {
  max: LiftMax;
  previousMax?: LiftMax;
  isFirst: boolean;
}

function HistoryEntry({ max, previousMax, isFirst }: HistoryEntryProps) {
  const delta = previousMax ? max.value - previousMax.value : 0;
  const showDelta = previousMax && previousMax.type === max.type;

  return (
    <div className="relative pl-8 pb-8 last:pb-0">
      {/* Timeline line */}
      <div className="absolute left-[11px] top-3 bottom-0 w-0.5 bg-border last:hidden" />

      {/* Timeline dot */}
      <div
        className={`
          absolute left-0 top-1.5 w-6 h-6 rounded-full
          flex items-center justify-center
          ${isFirst
            ? 'bg-accent/20 border-2 border-accent'
            : 'bg-surface-elevated border-2 border-border'
          }
        `}
      >
        {showDelta && delta > 0 && (
          <TrendingUp size={12} className="text-success" />
        )}
        {showDelta && delta < 0 && (
          <TrendingDown size={12} className="text-error" />
        )}
        {(!showDelta || delta === 0) && !isFirst && (
          <Minus size={12} className="text-muted" />
        )}
        {isFirst && !showDelta && (
          <div className="w-2 h-2 rounded-full bg-accent" />
        )}
      </div>

      {/* Entry content */}
      <div
        className={`
          bg-surface border rounded-lg p-4
          ${isFirst ? 'border-accent/30' : 'border-border'}
        `}
      >
        <div className="flex items-start justify-between gap-4">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <MaxTypeBadge type={max.type} />
              <span className="text-xs text-muted flex items-center gap-1">
                <Calendar size={12} />
                {formatDate(max.effectiveDate)}
              </span>
            </div>
            {isFirst && (
              <p className="text-xs text-accent font-medium mt-1">Current</p>
            )}
          </div>

          <div className="text-right">
            <p className="text-xl font-black tabular-nums text-foreground">
              {max.value}
              <span className="text-sm font-medium text-muted ml-1">lbs</span>
            </p>
            {showDelta && delta !== 0 && (
              <p
                className={`text-xs font-bold ${delta > 0 ? 'text-success' : 'text-error'}`}
              >
                {delta > 0 ? '+' : ''}{delta} lbs
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function EmptyState({ liftName }: { liftName: string }) {
  return (
    <div className="text-center py-16">
      <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-surface-elevated border border-border flex items-center justify-center">
        <TrendingUp className="w-10 h-10 text-muted" />
      </div>
      <h3 className="text-xl font-bold text-foreground mb-2">No history for {liftName}</h3>
      <p className="text-muted max-w-sm mx-auto mb-6">
        Add your first max to start tracking progress.
      </p>
      <Link
        to="/lift-maxes/new"
        className="
          inline-flex items-center gap-2 py-3 px-6
          bg-accent rounded-lg
          text-background font-bold
          hover:bg-accent-light
          transition-colors
        "
      >
        Add Max
      </Link>
    </div>
  );
}

function LoadingState() {
  return (
    <div className="flex items-center justify-center py-16">
      <Loader2 className="w-8 h-8 text-accent animate-spin" />
    </div>
  );
}

function ErrorState({ message }: { message: string }) {
  return (
    <div className="text-center py-12">
      <div className="mx-auto w-16 h-16 rounded-full bg-error/10 border border-error/20 flex items-center justify-center mb-4">
        <AlertCircle className="w-8 h-8 text-error" />
      </div>
      <h3 className="text-lg font-semibold text-foreground mb-2">
        Failed to load history
      </h3>
      <p className="text-sm text-muted max-w-sm mx-auto">{message}</p>
    </div>
  );
}

function StatsCard({ label, value, subtext }: { label: string; value: string; subtext?: string }) {
  return (
    <div className="bg-surface border border-border rounded-xl p-4 text-center">
      <p className="text-2xl font-black tabular-nums text-foreground">{value}</p>
      <p className="text-xs text-muted uppercase tracking-wider mt-1">{label}</p>
      {subtext && <p className="text-xs text-accent mt-1">{subtext}</p>}
    </div>
  );
}

export default function LiftMaxHistory() {
  const { liftId } = useParams<{ liftId: string }>();
  const { userId } = useAuth();
  const { data: maxesData, isLoading: maxesLoading, error: maxesError } = useLiftMaxes(
    userId ?? undefined,
    { liftId }
  );
  const { data: liftsData, isLoading: liftsLoading } = useLifts();

  const isLoading = maxesLoading || liftsLoading;
  const maxes = maxesData?.data ?? [];
  const lift = liftsData?.data?.find((l) => l.id === liftId);
  const liftName = lift?.name ?? 'Unknown Lift';

  // Sort by date descending (most recent first)
  const sortedMaxes = [...maxes].sort(
    (a, b) => new Date(b.effectiveDate).getTime() - new Date(a.effectiveDate).getTime()
  );

  // Calculate stats
  const oneRMs = sortedMaxes.filter((m) => m.type === 'ONE_RM');
  const trainingMaxes = sortedMaxes.filter((m) => m.type === 'TRAINING_MAX');

  const currentOneRM = oneRMs[0]?.value;
  const currentTM = trainingMaxes[0]?.value;
  const firstOneRM = oneRMs[oneRMs.length - 1]?.value;
  const totalGain = currentOneRM && firstOneRM ? currentOneRM - firstOneRM : null;

  return (
    <div className="py-6 md:py-8">
      <Container>
        {/* Header */}
        <div className="mb-8">
          <Link
            to="/lift-maxes"
            className="inline-flex items-center gap-1 text-sm text-muted hover:text-accent transition-colors mb-4"
          >
            <ChevronLeft size={16} />
            Back to Lift Maxes
          </Link>

          <div className="flex items-start gap-4">
            <div className="w-12 h-12 rounded-xl bg-accent/10 border border-accent/20 flex items-center justify-center flex-shrink-0">
              <Dumbbell className="w-6 h-6 text-accent" />
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground">
                {liftName}
              </h1>
              <p className="mt-1 text-muted">
                Progress history for this lift
              </p>
            </div>
          </div>
        </div>

        {/* Content */}
        {isLoading && <LoadingState />}

        {maxesError && (
          <ErrorState
            message={maxesError instanceof Error ? maxesError.message : 'An unexpected error occurred'}
          />
        )}

        {!isLoading && !maxesError && sortedMaxes.length === 0 && (
          <EmptyState liftName={liftName} />
        )}

        {!isLoading && !maxesError && sortedMaxes.length > 0 && (
          <>
            {/* Stats */}
            <div className="grid grid-cols-3 gap-4 mb-8">
              <StatsCard
                label="Current 1RM"
                value={currentOneRM ? `${currentOneRM}` : '-'}
                subtext={currentOneRM ? 'lbs' : undefined}
              />
              <StatsCard
                label="Training Max"
                value={currentTM ? `${currentTM}` : '-'}
                subtext={currentTM ? 'lbs' : undefined}
              />
              <StatsCard
                label="Total Gain"
                value={totalGain !== null ? `${totalGain > 0 ? '+' : ''}${totalGain}` : '-'}
                subtext={totalGain !== null ? 'lbs' : undefined}
              />
            </div>

            {/* Timeline */}
            <div className="bg-surface-elevated border border-border rounded-xl p-6">
              <h2 className="text-lg font-bold text-foreground mb-6">History</h2>
              <div>
                {sortedMaxes.map((max, index) => (
                  <HistoryEntry
                    key={max.id}
                    max={max}
                    previousMax={sortedMaxes[index + 1]}
                    isFirst={index === 0}
                  />
                ))}
              </div>
            </div>
          </>
        )}
      </Container>
    </div>
  );
}
