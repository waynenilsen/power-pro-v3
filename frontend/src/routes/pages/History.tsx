import { useState } from 'react';
import { Container } from '../../components/layout';
import { useAuth } from '../../contexts/useAuth';
import { useWorkoutSessions, useEnrollment } from '../../hooks';
import {
  Loader2,
  CheckCircle2,
  XCircle,
  Calendar,
  Clock,
  ChevronLeft,
  ChevronRight,
  Dumbbell,
  Filter,
  Trophy,
  Flame,
} from 'lucide-react';
import type { WorkoutSession, WorkoutSessionStatus } from '../../api/types';

type FilterStatus = 'ALL' | WorkoutSessionStatus;

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  });
}

function formatTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
  });
}

function calculateDuration(startedAt: string, finishedAt?: string): string | null {
  if (!finishedAt) return null;
  const start = new Date(startedAt);
  const end = new Date(finishedAt);
  const diffMs = end.getTime() - start.getTime();
  const diffMins = Math.round(diffMs / 60000);

  if (diffMins < 60) {
    return `${diffMins} min`;
  }
  const hours = Math.floor(diffMins / 60);
  const mins = diffMins % 60;
  return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
}

function StatusBadge({ status }: { status: WorkoutSessionStatus }) {
  const config = {
    COMPLETED: {
      icon: CheckCircle2,
      label: 'Completed',
      className: 'bg-success/15 text-success border-success/30',
    },
    ABANDONED: {
      icon: XCircle,
      label: 'Abandoned',
      className: 'bg-error/15 text-error border-error/30',
    },
    IN_PROGRESS: {
      icon: Flame,
      label: 'In Progress',
      className: 'bg-warning/15 text-warning border-warning/30 animate-pulse',
    },
  }[status];

  const Icon = config.icon;

  return (
    <span
      className={`
        inline-flex items-center gap-1.5 px-2.5 py-1
        text-xs font-bold uppercase tracking-wider
        rounded-md border
        ${config.className}
      `}
    >
      <Icon size={12} />
      {config.label}
    </span>
  );
}

function WorkoutSessionCard({ session }: { session: WorkoutSession }) {
  const duration = calculateDuration(session.startedAt, session.finishedAt);
  const isCompleted = session.status === 'COMPLETED';
  const isAbandoned = session.status === 'ABANDONED';

  return (
    <div
      className={`
        group relative
        bg-surface border rounded-xl p-5
        transition-all duration-300
        hover:bg-surface-elevated
        ${isCompleted ? 'border-success/20 hover:border-success/40' : ''}
        ${isAbandoned ? 'border-error/20 hover:border-error/40 opacity-75' : ''}
        ${!isCompleted && !isAbandoned ? 'border-warning/30' : ''}
      `}
    >
      {/* Completion glow effect for completed workouts */}
      {isCompleted && (
        <div className="absolute inset-0 rounded-xl bg-success/5 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
      )}

      <div className="relative flex items-start justify-between gap-4">
        {/* Left side - Main info */}
        <div className="flex-1 min-w-0">
          {/* Week/Day header */}
          <div className="flex items-center gap-3 mb-3">
            <div className="w-10 h-10 rounded-lg bg-accent/10 border border-accent/20 flex items-center justify-center flex-shrink-0">
              <Dumbbell className="w-5 h-5 text-accent" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-foreground">
                Week {session.weekNumber}, Day {session.dayIndex + 1}
              </h3>
              <div className="flex items-center gap-2 text-sm text-muted">
                <Calendar size={14} />
                <span>{formatDate(session.startedAt)}</span>
                <span className="text-border">•</span>
                <span>{formatTime(session.startedAt)}</span>
              </div>
            </div>
          </div>

          {/* Duration (if available) */}
          {duration && (
            <div className="flex items-center gap-2 text-sm text-muted mt-2">
              <Clock size={14} className="text-accent" />
              <span className="font-medium">{duration}</span>
            </div>
          )}
        </div>

        {/* Right side - Status */}
        <div className="flex flex-col items-end gap-2">
          <StatusBadge status={session.status} />
          {isCompleted && (
            <Trophy size={20} className="text-success/50" />
          )}
        </div>
      </div>
    </div>
  );
}

function FilterButton({
  status,
  currentFilter,
  onClick,
  count,
}: {
  status: FilterStatus;
  currentFilter: FilterStatus;
  onClick: () => void;
  count?: number;
}) {
  const isActive = status === currentFilter;
  const labels: Record<FilterStatus, string> = {
    ALL: 'All',
    COMPLETED: 'Completed',
    ABANDONED: 'Abandoned',
    IN_PROGRESS: 'Active',
  };

  return (
    <button
      onClick={onClick}
      className={`
        px-4 py-2 rounded-lg text-sm font-semibold
        transition-all duration-200
        ${isActive
          ? 'bg-accent text-background'
          : 'bg-surface-elevated text-muted hover:text-foreground hover:bg-surface border border-border'
        }
      `}
    >
      {labels[status]}
      {count !== undefined && count > 0 && (
        <span className={`ml-1.5 ${isActive ? 'text-background/70' : 'text-muted'}`}>
          ({count})
        </span>
      )}
    </button>
  );
}

function EmptyState({ filter, isEnrolled }: { filter: FilterStatus; isEnrolled: boolean }) {
  const messages: Record<FilterStatus, { title: string; subtitle: string }> = {
    ALL: {
      title: 'No workouts yet',
      subtitle: isEnrolled
        ? 'Start your first workout to begin tracking your progress.'
        : 'Enroll in a program to start tracking your workouts.',
    },
    COMPLETED: {
      title: 'No completed workouts',
      subtitle: 'Finish a workout session to see it here.',
    },
    ABANDONED: {
      title: 'No abandoned workouts',
      subtitle: "That's great! Keep up the consistency.",
    },
    IN_PROGRESS: {
      title: 'No active workout',
      subtitle: 'Start a new workout session from the home page.',
    },
  };

  const { title, subtitle } = messages[filter];

  return (
    <div className="text-center py-16">
      <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-surface-elevated border border-border flex items-center justify-center">
        <Dumbbell className="w-10 h-10 text-muted" />
      </div>
      <h3 className="text-xl font-bold text-foreground mb-2">{title}</h3>
      <p className="text-muted max-w-sm mx-auto">{subtitle}</p>
    </div>
  );
}

function Pagination({
  page,
  totalPages,
  onPageChange,
}: {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}) {
  if (totalPages <= 1) return null;

  return (
    <div className="flex items-center justify-center gap-2 mt-8">
      <button
        onClick={() => onPageChange(page - 1)}
        disabled={page <= 1}
        className="
          p-2 rounded-lg
          bg-surface border border-border
          text-muted hover:text-foreground
          disabled:opacity-50 disabled:cursor-not-allowed
          transition-colors
        "
      >
        <ChevronLeft size={20} />
      </button>

      <div className="px-4 py-2 text-sm font-medium text-muted">
        Page <span className="text-foreground tabular-nums">{page}</span> of{' '}
        <span className="text-foreground tabular-nums">{totalPages}</span>
      </div>

      <button
        onClick={() => onPageChange(page + 1)}
        disabled={page >= totalPages}
        className="
          p-2 rounded-lg
          bg-surface border border-border
          text-muted hover:text-foreground
          disabled:opacity-50 disabled:cursor-not-allowed
          transition-colors
        "
      >
        <ChevronRight size={20} />
      </button>
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
          Loading history...
        </p>
      </div>
    </div>
  );
}

function StatsHeader({
  totalCompleted,
  totalAbandoned,
}: {
  totalCompleted: number;
  totalAbandoned: number;
}) {
  const total = totalCompleted + totalAbandoned;
  const completionRate = total > 0 ? Math.round((totalCompleted / total) * 100) : 0;

  return (
    <div className="grid grid-cols-3 gap-4 mb-8">
      <div className="bg-surface border border-border rounded-xl p-4 text-center">
        <p className="text-3xl font-black tabular-nums text-foreground">{total}</p>
        <p className="text-xs text-muted uppercase tracking-wider mt-1">Total Workouts</p>
      </div>
      <div className="bg-surface border border-success/20 rounded-xl p-4 text-center">
        <p className="text-3xl font-black tabular-nums text-success">{totalCompleted}</p>
        <p className="text-xs text-muted uppercase tracking-wider mt-1">Completed</p>
      </div>
      <div className="bg-surface border border-border rounded-xl p-4 text-center">
        <p className="text-3xl font-black tabular-nums text-accent">{completionRate}%</p>
        <p className="text-xs text-muted uppercase tracking-wider mt-1">Success Rate</p>
      </div>
    </div>
  );
}

export default function History() {
  const { userId } = useAuth();
  const [page, setPage] = useState(1);
  const [filter, setFilter] = useState<FilterStatus>('ALL');
  const pageSize = 10;

  const { data: enrollment } = useEnrollment(userId ?? undefined);
  const { data: sessions, isLoading } = useWorkoutSessions(userId ?? undefined, {
    page,
    pageSize,
    status: filter === 'ALL' ? undefined : filter,
  });

  // Fetch all sessions for stats (separate query without pagination)
  const { data: allCompletedSessions } = useWorkoutSessions(userId ?? undefined, {
    status: 'COMPLETED',
    pageSize: 1,
  });
  const { data: allAbandonedSessions } = useWorkoutSessions(userId ?? undefined, {
    status: 'ABANDONED',
    pageSize: 1,
  });

  const isEnrolled = enrollment?.enrollmentStatus === 'ACTIVE' ||
    enrollment?.enrollmentStatus === 'BETWEEN_CYCLES';

  const totalCompleted = allCompletedSessions?.meta.total ?? 0;
  const totalAbandoned = allAbandonedSessions?.meta.total ?? 0;

  const handleFilterChange = (newFilter: FilterStatus) => {
    setFilter(newFilter);
    setPage(1); // Reset to first page when filter changes
  };

  return (
    <div className="py-6 md:py-8">
      <Container>
        {/* Page header */}
        <div className="mb-8">
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground">
            Workout History
          </h1>
          <p className="mt-2 text-muted">
            Track your training progress over time.
          </p>
        </div>

        {isLoading ? (
          <LoadingState />
        ) : (
          <>
            {/* Stats header - only show if there are workouts */}
            {(totalCompleted > 0 || totalAbandoned > 0) && (
              <StatsHeader
                totalCompleted={totalCompleted}
                totalAbandoned={totalAbandoned}
              />
            )}

            {/* Filter buttons */}
            <div className="flex items-center gap-2 mb-6 overflow-x-auto pb-2">
              <Filter size={16} className="text-muted flex-shrink-0" />
              <div className="flex gap-2">
                <FilterButton
                  status="ALL"
                  currentFilter={filter}
                  onClick={() => handleFilterChange('ALL')}
                />
                <FilterButton
                  status="COMPLETED"
                  currentFilter={filter}
                  onClick={() => handleFilterChange('COMPLETED')}
                  count={totalCompleted}
                />
                <FilterButton
                  status="ABANDONED"
                  currentFilter={filter}
                  onClick={() => handleFilterChange('ABANDONED')}
                  count={totalAbandoned}
                />
              </div>
            </div>

            {/* Sessions list */}
            {sessions?.data && sessions.data.length > 0 ? (
              <>
                <div className="space-y-4">
                  {sessions.data.map((session) => (
                    <WorkoutSessionCard key={session.id} session={session} />
                  ))}
                </div>

                <Pagination
                  page={page}
                  totalPages={Math.ceil(sessions.meta.total / sessions.meta.limit)}
                  onPageChange={setPage}
                />
              </>
            ) : (
              <EmptyState filter={filter} isEnrolled={isEnrolled} />
            )}
          </>
        )}
      </Container>
    </div>
  );
}
