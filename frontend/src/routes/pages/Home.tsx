import { Link } from 'react-router-dom';
import { Container } from '../../components/layout';
import { useAuth } from '../../contexts/useAuth';
import { useEnrollment, useAdvanceState } from '../../hooks';
import { BookOpen, ChevronRight, Dumbbell, Play, RotateCcw, Trophy, Loader2 } from 'lucide-react';
import type { EnrollmentStatus, WorkoutSessionSummary } from '../../api/types';

function EnrolledProgramCard({
  programName,
  programSlug,
  currentWeek,
  currentDayIndex,
  cycleLengthWeeks,
  currentCycleIteration,
  enrollmentStatus,
  currentWorkoutSession,
  onStartNextCycle,
  isAdvancing,
}: {
  programName: string;
  programSlug: string;
  currentWeek: number;
  currentDayIndex?: number;
  cycleLengthWeeks: number;
  currentCycleIteration: number;
  enrollmentStatus: EnrollmentStatus;
  currentWorkoutSession?: WorkoutSessionSummary;
  onStartNextCycle: () => void;
  isAdvancing: boolean;
}) {
  const progressPercentage = (currentWeek / cycleLengthWeeks) * 100;
  const hasActiveSession = currentWorkoutSession?.status === 'IN_PROGRESS';
  const isBetweenCycles = enrollmentStatus === 'BETWEEN_CYCLES';

  // Format day display (1-indexed for display)
  const dayDisplay = currentDayIndex !== undefined ? currentDayIndex + 1 : null;

  return (
    <div className="space-y-4">
      {/* Current program card */}
      <div className="bg-surface border border-border rounded-lg p-5">
        <div className="flex items-start justify-between gap-4 mb-4">
          <div>
            <p className="text-xs uppercase tracking-wider text-muted mb-1">
              Current Program
            </p>
            <h2 className="text-xl font-bold text-foreground">{programName}</h2>
          </div>
          <div className="w-10 h-10 rounded-lg bg-accent/10 border border-accent/20 flex items-center justify-center flex-shrink-0">
            <Dumbbell className="w-5 h-5 text-accent" />
          </div>
        </div>

        {/* Current position indicator */}
        {!isBetweenCycles && (
          <div className="mb-3 px-3 py-2 bg-surface-elevated rounded-md border border-border">
            <p className="text-sm font-medium text-foreground">
              Week {currentWeek}{dayDisplay !== null && `, Day ${dayDisplay}`}
            </p>
            <p className="text-xs text-muted">Cycle {currentCycleIteration}</p>
          </div>
        )}

        {/* Progress */}
        <div className="mb-4">
          <div className="flex justify-between text-sm mb-2">
            <span className="text-muted">
              Cycle {currentCycleIteration} · Week {currentWeek} of {cycleLengthWeeks}
            </span>
            <span className="text-accent font-medium tabular-nums">
              {Math.round(progressPercentage)}%
            </span>
          </div>
          <div className="h-2 bg-surface-elevated rounded-full overflow-hidden">
            <div
              className="h-full bg-accent rounded-full transition-all duration-500"
              style={{ width: `${progressPercentage}%` }}
            />
          </div>
        </div>

        {/* Active session indicator */}
        {hasActiveSession && (
          <div className="mb-4 px-3 py-2 bg-accent/10 border border-accent/20 rounded-md">
            <p className="text-sm font-medium text-accent">Workout in progress</p>
            <p className="text-xs text-muted">
              Started {new Date(currentWorkoutSession!.startedAt).toLocaleTimeString([], { hour: 'numeric', minute: '2-digit' })}
            </p>
          </div>
        )}

        {/* View details link */}
        <Link
          to={`/programs/${programSlug}`}
          className="text-sm text-muted hover:text-accent transition-colors"
        >
          View program details →
        </Link>
      </div>

      {/* Cycle complete message */}
      {isBetweenCycles && (
        <div className="bg-surface border border-accent/30 rounded-lg p-5 text-center">
          <div className="w-12 h-12 mx-auto mb-3 rounded-full bg-accent/10 border border-accent/20 flex items-center justify-center">
            <Trophy className="w-6 h-6 text-accent" />
          </div>
          <h3 className="text-lg font-bold text-foreground mb-1">Cycle Complete!</h3>
          <p className="text-sm text-muted mb-4">
            You've finished cycle {currentCycleIteration}. Ready for the next one?
          </p>
          <button
            onClick={onStartNextCycle}
            disabled={isAdvancing}
            className="
              group flex items-center justify-center gap-3
              w-full py-4 px-6
              bg-accent rounded-lg
              text-background font-bold text-lg
              hover:bg-accent-light
              active:scale-[0.98]
              transition-all duration-200
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            {isAdvancing ? (
              <Loader2 size={22} className="animate-spin" />
            ) : (
              <RotateCcw size={22} className="group-hover:scale-110 transition-transform" />
            )}
            <span>{isAdvancing ? 'Starting...' : 'Start Next Cycle'}</span>
          </button>
        </div>
      )}

      {/* Start/Continue workout CTA */}
      {!isBetweenCycles && (
        <Link
          to="/workout"
          className="
            group flex items-center justify-center gap-3
            w-full py-4 px-6
            bg-accent rounded-lg
            text-background font-bold text-lg
            hover:bg-accent-light
            active:scale-[0.98]
            transition-all duration-200
          "
        >
          <Play size={22} className="group-hover:scale-110 transition-transform" />
          <span>{hasActiveSession ? 'Continue Workout' : 'Start Workout'}</span>
        </Link>
      )}
    </div>
  );
}

function BrowseProgramsCard() {
  return (
    <Link
      to="/programs"
      className="
        group flex items-center gap-4
        bg-surface border border-border rounded-lg
        p-4 sm:p-5
        hover:border-accent/30 hover:bg-surface-elevated
        transition-all duration-200
      "
    >
      <div className="w-12 h-12 rounded-lg bg-accent/10 border border-accent/20 flex items-center justify-center flex-shrink-0">
        <BookOpen className="w-6 h-6 text-accent" />
      </div>
      <div className="flex-1 min-w-0">
        <h2 className="text-lg font-semibold text-foreground group-hover:text-accent transition-colors">
          Browse Programs
        </h2>
        <p className="text-sm text-muted">
          Explore available training programs
        </p>
      </div>
      <ChevronRight
        size={20}
        className="flex-shrink-0 text-muted group-hover:text-accent group-hover:translate-x-0.5 transition-all"
      />
    </Link>
  );
}

function LoadingState() {
  return (
    <div className="flex items-center justify-center py-8">
      <Loader2 className="w-6 h-6 text-accent animate-spin" />
    </div>
  );
}

export default function Home() {
  const { userId } = useAuth();
  const { data: enrollment, isLoading } = useEnrollment(userId ?? undefined);
  const advanceState = useAdvanceState(userId ?? undefined);

  const isEnrolled = enrollment &&
    (enrollment.enrollmentStatus === 'ACTIVE' || enrollment.enrollmentStatus === 'BETWEEN_CYCLES');

  const handleStartNextCycle = () => {
    advanceState.mutate({ advanceType: 'week' });
  };

  return (
    <div className="py-6 md:py-8">
      <Container>
        <h1 className="text-2xl md:text-3xl font-bold tracking-tight">
          Welcome to <span className="text-accent">PowerPro</span>
        </h1>
        <p className="mt-2 text-muted">
          Track your powerlifting progress and get stronger.
        </p>

        <div className="mt-8">
          {isLoading ? (
            <LoadingState />
          ) : isEnrolled && enrollment ? (
            <EnrolledProgramCard
              programName={enrollment.program.name}
              programSlug={enrollment.program.slug}
              currentWeek={enrollment.state.currentWeek}
              currentDayIndex={enrollment.state.currentDayIndex}
              cycleLengthWeeks={enrollment.program.cycleLengthWeeks}
              currentCycleIteration={enrollment.state.currentCycleIteration}
              enrollmentStatus={enrollment.enrollmentStatus}
              currentWorkoutSession={enrollment.currentWorkoutSession}
              onStartNextCycle={handleStartNextCycle}
              isAdvancing={advanceState.isPending}
            />
          ) : (
            <BrowseProgramsCard />
          )}
        </div>
      </Container>
    </div>
  );
}
