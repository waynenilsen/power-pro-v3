import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Container } from '../../components/layout';
import { useAuth } from '../../contexts/useAuth';
import { useLiftMaxes, useLifts, useDeleteLiftMax } from '../../hooks';
import {
  Loader2,
  Plus,
  ChevronLeft,
  Dumbbell,
  AlertCircle,
  Trash2,
  Edit2,
  TrendingUp,
  Calendar,
} from 'lucide-react';
import type { LiftMax, Lift, MaxType } from '../../api/types';
import { ConfirmDialog } from '../../components/ui/ConfirmDialog';

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
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

interface LiftMaxCardProps {
  liftMax: LiftMax;
  liftName: string;
  onEdit: () => void;
  onDelete: () => void;
}

function LiftMaxCard({ liftMax, liftName, onEdit, onDelete }: LiftMaxCardProps) {
  return (
    <div
      className="
        group relative
        bg-surface border border-border rounded-xl p-5
        transition-all duration-200
        hover:bg-surface-elevated hover:border-border
      "
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-10 h-10 rounded-lg bg-accent/10 border border-accent/20 flex items-center justify-center flex-shrink-0">
              <Dumbbell className="w-5 h-5 text-accent" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-foreground">{liftName}</h3>
              <div className="flex items-center gap-2 mt-1">
                <MaxTypeBadge type={liftMax.type} />
                <span className="text-muted text-xs flex items-center gap-1">
                  <Calendar size={12} />
                  {formatDate(liftMax.effectiveDate)}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div className="text-right">
          <p className="text-2xl font-black tabular-nums text-foreground">
            {liftMax.value}
            <span className="text-sm font-medium text-muted ml-1">lbs</span>
          </p>
        </div>
      </div>

      {/* Action buttons */}
      <div className="absolute top-3 right-3 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          onClick={onEdit}
          className="
            p-2 rounded-lg
            bg-surface-elevated border border-border
            text-muted hover:text-accent hover:border-accent/30
            transition-all duration-200
          "
          title="Edit"
        >
          <Edit2 size={14} />
        </button>
        <button
          onClick={onDelete}
          className="
            p-2 rounded-lg
            bg-surface-elevated border border-border
            text-muted hover:text-error hover:border-error/30
            transition-all duration-200
          "
          title="Delete"
        >
          <Trash2 size={14} />
        </button>
      </div>
    </div>
  );
}

interface GroupedMaxes {
  liftId: string;
  liftName: string;
  maxes: LiftMax[];
}

function groupMaxesByLift(maxes: LiftMax[], lifts: Lift[]): GroupedMaxes[] {
  const liftMap = new Map(lifts.map((l) => [l.id, l.name]));
  const grouped = new Map<string, LiftMax[]>();

  for (const max of maxes) {
    const existing = grouped.get(max.liftId) ?? [];
    grouped.set(max.liftId, [...existing, max]);
  }

  return Array.from(grouped.entries())
    .map(([liftId, liftMaxes]) => ({
      liftId,
      liftName: liftMap.get(liftId) ?? 'Unknown Lift',
      maxes: liftMaxes.sort((a, b) =>
        new Date(b.effectiveDate).getTime() - new Date(a.effectiveDate).getTime()
      ),
    }))
    .sort((a, b) => a.liftName.localeCompare(b.liftName));
}

function EmptyState() {
  return (
    <div className="text-center py-16">
      <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-surface-elevated border border-border flex items-center justify-center">
        <TrendingUp className="w-10 h-10 text-muted" />
      </div>
      <h3 className="text-xl font-bold text-foreground mb-2">No lift maxes recorded</h3>
      <p className="text-muted max-w-sm mx-auto mb-6">
        Add your first lift max to start tracking your strength progress.
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
        <Plus size={18} />
        Add First Max
      </Link>
    </div>
  );
}

function LoadingState() {
  return (
    <div className="flex items-center justify-center py-16">
      <div className="flex flex-col items-center gap-4">
        <div className="relative">
          <div className="w-16 h-16 rounded-2xl bg-accent/10 border border-accent/20 flex items-center justify-center">
            <Dumbbell className="w-8 h-8 text-accent animate-pulse" />
          </div>
          <Loader2 className="absolute -bottom-1 -right-1 w-6 h-6 text-accent animate-spin" />
        </div>
        <p className="text-muted text-sm font-medium tracking-wide uppercase">
          Loading lift maxes...
        </p>
      </div>
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
        Failed to load lift maxes
      </h3>
      <p className="text-sm text-muted max-w-sm mx-auto">{message}</p>
    </div>
  );
}

export default function LiftMaxes() {
  const navigate = useNavigate();
  const { userId } = useAuth();
  const { data: maxesData, isLoading: maxesLoading, error: maxesError } = useLiftMaxes(userId ?? undefined);
  const { data: liftsData, isLoading: liftsLoading } = useLifts();
  const deleteMax = useDeleteLiftMax(userId ?? undefined);

  const [deleteTarget, setDeleteTarget] = useState<{ id: string; liftName: string } | null>(null);

  const isLoading = maxesLoading || liftsLoading;
  const maxes = maxesData?.data ?? [];
  const lifts = liftsData?.data ?? [];

  const groupedMaxes = groupMaxesByLift(maxes, lifts);

  const handleDeleteConfirm = async () => {
    if (deleteTarget) {
      await deleteMax.mutateAsync(deleteTarget.id);
      setDeleteTarget(null);
    }
  };

  return (
    <div className="py-6 md:py-8">
      <Container>
        {/* Header */}
        <div className="mb-8">
          <Link
            to="/profile"
            className="inline-flex items-center gap-1 text-sm text-muted hover:text-accent transition-colors mb-4"
          >
            <ChevronLeft size={16} />
            Back to Profile
          </Link>

          <div className="flex items-start justify-between gap-4">
            <div>
              <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground">
                Lift Maxes
              </h1>
              <p className="mt-2 text-muted">
                Track your 1RM and training maxes for each lift.
              </p>
            </div>

            <Link
              to="/lift-maxes/new"
              className="
                flex items-center gap-2 py-2.5 px-4
                bg-accent rounded-lg
                text-background font-bold text-sm
                hover:bg-accent-light
                transition-colors
                flex-shrink-0
              "
            >
              <Plus size={18} />
              <span className="hidden sm:inline">Add Max</span>
            </Link>
          </div>
        </div>

        {/* Content */}
        {isLoading && <LoadingState />}

        {maxesError && (
          <ErrorState
            message={maxesError instanceof Error ? maxesError.message : 'An unexpected error occurred'}
          />
        )}

        {!isLoading && !maxesError && maxes.length === 0 && <EmptyState />}

        {!isLoading && !maxesError && groupedMaxes.length > 0 && (
          <div className="space-y-6">
            {groupedMaxes.map((group) => (
              <div key={group.liftId}>
                <h2 className="text-lg font-semibold text-foreground mb-3 flex items-center gap-2">
                  {group.liftName}
                  <span className="text-sm font-normal text-muted">
                    ({group.maxes.length} {group.maxes.length === 1 ? 'record' : 'records'})
                  </span>
                </h2>
                <div className="space-y-3">
                  {group.maxes.map((max) => (
                    <LiftMaxCard
                      key={max.id}
                      liftMax={max}
                      liftName={group.liftName}
                      onEdit={() => navigate(`/lift-maxes/${max.id}/edit`)}
                      onDelete={() => setDeleteTarget({ id: max.id, liftName: group.liftName })}
                    />
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </Container>

      {/* Delete confirmation dialog */}
      <ConfirmDialog
        isOpen={deleteTarget !== null}
        onCancel={() => setDeleteTarget(null)}
        onConfirm={handleDeleteConfirm}
        title="Delete Lift Max"
        message={`Are you sure you want to delete this ${deleteTarget?.liftName} max? This action cannot be undone.`}
        confirmLabel="Delete"
        variant="destructive"
        isLoading={deleteMax.isPending}
      />
    </div>
  );
}
