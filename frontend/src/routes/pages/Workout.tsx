import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container } from '../../components/layout';
import { ExerciseCard } from '../../components/workout/ExerciseCard';
import { ConfirmDialog } from '../../components/ui/ConfirmDialog';
import { useAuth } from '../../contexts/useAuth';
import {
  useEnrollment,
  useCurrentWorkout,
  useCurrentWorkoutSession,
  useStartWorkoutSession,
  useFinishWorkoutSession,
  useAbandonWorkoutSession,
  useAdvanceState,
} from '../../hooks';
import {
  Loader2,
  Play,
  Flag,
  XCircle,
  ArrowLeft,
  Dumbbell,
  ChevronRight,
  Trophy,
  RotateCcw,
  Home,
  Settings,
} from 'lucide-react';

type DialogType = 'abandon' | 'finish' | 'advanceWeek' | 'startNextCycle' | null;

export default function Workout() {
  const navigate = useNavigate();
  const { userId } = useAuth();

  const { data: enrollment, isLoading: enrollmentLoading } = useEnrollment(userId ?? undefined);
  const { data: workout, isLoading: workoutLoading, error: workoutError } = useCurrentWorkout(userId ?? undefined);
  const { data: session, isLoading: sessionLoading } = useCurrentWorkoutSession(userId ?? undefined);

  const startSession = useStartWorkoutSession(userId ?? undefined);
  const finishSession = useFinishWorkoutSession(userId ?? undefined);
  const abandonSession = useAbandonWorkoutSession(userId ?? undefined);
  const advanceState = useAdvanceState(userId ?? undefined);

  const [completedSets, setCompletedSets] = useState<Set<string>>(new Set());
  const [activeDialog, setActiveDialog] = useState<DialogType>(null);
  const [showCompletionSummary, setShowCompletionSummary] = useState(false);

  const isLoading = enrollmentLoading || workoutLoading || sessionLoading;
  const isEnrolled = enrollment?.enrollmentStatus === 'ACTIVE';
  const hasActiveSession = session?.status === 'IN_PROGRESS';

  // Check if the error is about missing lift maxes
  const isMissingLiftMaxes = workoutError?.message?.toLowerCase().includes('missing lift max') ||
    workoutError?.message?.toLowerCase().includes('training max');

  const handleSetToggle = (exerciseIndex: number, setNumber: number) => {
    const key = `${exerciseIndex}-${setNumber}`;
    setCompletedSets(prev => {
      const next = new Set(prev);
      if (next.has(key)) {
        next.delete(key);
      } else {
        next.add(key);
      }
      return next;
    });
  };

  const handleStartWorkout = () => {
    startSession.mutate(undefined, {
      onSuccess: () => {
        setCompletedSets(new Set());
      },
    });
  };

  const handleFinishWorkout = () => {
    if (!session?.id) return;
    finishSession.mutate(session.id, {
      onSuccess: () => {
        // Advance to next day after finishing workout
        advanceState.mutate(
          { advanceType: 'day' },
          {
            onSettled: () => {
              setActiveDialog(null);
              setShowCompletionSummary(true);
            },
          }
        );
      },
    });
  };

  const handleAbandonWorkout = () => {
    if (!session?.id) return;
    abandonSession.mutate(session.id, {
      onSuccess: () => {
        setActiveDialog(null);
        setCompletedSets(new Set());
        navigate('/');
      },
    });
  };

  const handleAdvanceWeek = () => {
    advanceState.mutate(
      { advanceType: 'week' },
      {
        onSuccess: () => {
          setActiveDialog(null);
          setShowCompletionSummary(false);
          navigate('/');
        },
      }
    );
  };

  const handleGoHome = () => {
    setShowCompletionSummary(false);
    setCompletedSets(new Set());
    navigate('/');
  };

  // Calculate progress
  const totalSets = workout?.exercises.reduce((acc, ex) => acc + ex.sets.length, 0) ?? 0;
  const completedCount = completedSets.size;
  const progressPercentage = totalSets > 0 ? (completedCount / totalSets) * 100 : 0;

  // Determine if we're at end of week or cycle
  const currentDayIndex = enrollment?.state.currentDayIndex ?? 0;
  const currentWeek = enrollment?.state.currentWeek ?? 1;
  const cycleLengthWeeks = enrollment?.program.cycleLengthWeeks ?? 1;
  const daysPerWeek = enrollment?.program.daysPerWeek ?? 4;
  const isLastDayOfWeek = currentDayIndex >= daysPerWeek - 1;
  const isLastWeekOfCycle = currentWeek >= cycleLengthWeeks;

  // Loading state
  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <div className="relative">
            <div className="w-16 h-16 rounded-2xl bg-accent/10 border border-accent/20 flex items-center justify-center">
              <Dumbbell className="w-8 h-8 text-accent animate-pulse" />
            </div>
            <Loader2 className="absolute -bottom-1 -right-1 w-6 h-6 text-accent animate-spin" />
          </div>
          <p className="text-muted text-sm font-medium tracking-wide uppercase">
            Loading workout...
          </p>
        </div>
      </div>
    );
  }

  // Not enrolled state
  if (!isEnrolled) {
    return (
      <div className="py-8">
        <Container>
          <div className="max-w-md mx-auto text-center">
            <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-surface-elevated border border-border flex items-center justify-center">
              <Dumbbell className="w-10 h-10 text-muted" />
            </div>
            <h1 className="text-2xl font-bold text-foreground mb-3">
              No Active Program
            </h1>
            <p className="text-muted mb-8">
              You need to enroll in a training program before you can start working out.
            </p>
            <button
              onClick={() => navigate('/programs')}
              className="
                inline-flex items-center gap-2 px-6 py-3
                bg-accent text-background font-bold rounded-lg
                hover:bg-accent-light active:scale-[0.98]
                transition-all duration-200
              "
            >
              Browse Programs
              <ChevronRight size={18} />
            </button>
          </div>
        </Container>
      </div>
    );
  }

  // Missing lift maxes state
  if (isMissingLiftMaxes) {
    return (
      <div className="py-8">
        <Container>
          <div className="max-w-md mx-auto text-center">
            <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-warning/10 border border-warning/30 flex items-center justify-center">
              <Settings className="w-10 h-10 text-warning" />
            </div>
            <h1 className="text-2xl font-bold text-foreground mb-3">
              Set Up Your Training Maxes
            </h1>
            <p className="text-muted mb-8">
              Before you can generate workouts, you need to set your training maxes for the main lifts.
              This helps calculate the weights for your program.
            </p>
            <button
              onClick={() => navigate('/profile')}
              className="
                inline-flex items-center gap-2 px-6 py-3
                bg-accent text-background font-bold rounded-lg
                hover:bg-accent-light active:scale-[0.98]
                transition-all duration-200
              "
            >
              Go to Profile Settings
              <ChevronRight size={18} />
            </button>
          </div>
        </Container>
      </div>
    );
  }

  // Completion summary state
  if (showCompletionSummary) {
    return (
      <div className="py-8">
        <Container>
          <div className="max-w-lg mx-auto">
            {/* Success hero */}
            <div className="relative text-center mb-8">
              <div className="absolute inset-0 flex items-center justify-center opacity-10">
                <div className="w-64 h-64 rounded-full bg-success blur-3xl" />
              </div>
              <div className="relative">
                <div className="w-24 h-24 mx-auto mb-6 rounded-full bg-success/10 border-2 border-success/30 flex items-center justify-center">
                  <Trophy className="w-12 h-12 text-success" />
                </div>
                <h1 className="text-3xl font-black tracking-tight text-foreground mb-2">
                  Workout Complete!
                </h1>
                <p className="text-muted">
                  Great work. You finished today's training session.
                </p>
              </div>
            </div>

            {/* Stats card */}
            <div className="bg-surface border border-border rounded-xl p-6 mb-6">
              <div className="grid grid-cols-2 gap-4">
                <div className="text-center p-4 rounded-lg bg-surface-elevated">
                  <p className="text-3xl font-black tabular-nums text-accent">
                    {completedCount}
                  </p>
                  <p className="text-sm text-muted mt-1">Sets Completed</p>
                </div>
                <div className="text-center p-4 rounded-lg bg-surface-elevated">
                  <p className="text-3xl font-black tabular-nums text-accent">
                    {workout?.exercises.length ?? 0}
                  </p>
                  <p className="text-sm text-muted mt-1">Exercises</p>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="space-y-3">
              {isLastDayOfWeek && !isLastWeekOfCycle && (
                <button
                  onClick={() => setActiveDialog('advanceWeek')}
                  className="
                    w-full flex items-center justify-center gap-3
                    py-4 px-6 rounded-xl
                    bg-accent text-background font-bold text-lg
                    hover:bg-accent-light active:scale-[0.98]
                    transition-all duration-200
                  "
                >
                  <ChevronRight size={22} />
                  Advance to Week {currentWeek + 1}
                </button>
              )}

              {isLastDayOfWeek && isLastWeekOfCycle && (
                <button
                  onClick={() => setActiveDialog('startNextCycle')}
                  className="
                    w-full flex items-center justify-center gap-3
                    py-4 px-6 rounded-xl
                    bg-success text-background font-bold text-lg
                    hover:bg-success/90 active:scale-[0.98]
                    transition-all duration-200
                  "
                >
                  <RotateCcw size={22} />
                  Start Next Cycle
                </button>
              )}

              <button
                onClick={handleGoHome}
                className="
                  w-full flex items-center justify-center gap-3
                  py-4 px-6 rounded-xl
                  bg-surface border border-border text-foreground font-bold text-lg
                  hover:bg-surface-elevated active:scale-[0.98]
                  transition-all duration-200
                "
              >
                <Home size={22} />
                Go Home
              </button>
            </div>
          </div>
        </Container>
      </div>
    );
  }

  // Workout preview (no active session)
  if (!hasActiveSession) {
    return (
      <div className="py-6">
        <Container>
          {/* Header */}
          <div className="flex items-center gap-3 mb-6">
            <button
              onClick={() => navigate('/')}
              className="p-2 -ml-2 rounded-lg text-muted hover:text-foreground hover:bg-surface-elevated transition-colors"
            >
              <ArrowLeft size={20} />
            </button>
            <div>
              <p className="text-xs uppercase tracking-wider text-muted">
                Week {currentWeek} · Day {currentDayIndex + 1}
              </p>
              <h1 className="text-xl font-bold text-foreground">
                {workout?.daySlug?.replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase()) ?? 'Today\'s Workout'}
              </h1>
            </div>
          </div>

          {/* Workout preview */}
          <div className="space-y-4 mb-8">
            {workout?.exercises.map((exercise, index) => (
              <ExerciseCard
                key={`${exercise.prescriptionId}-${index}`}
                exercise={exercise}
                exerciseIndex={index}
                completedSets={new Set()}
                onSetToggle={() => {}}
                isReadOnly
              />
            ))}
          </div>

          {/* Start button */}
          <div className="sticky bottom-4">
            <button
              onClick={handleStartWorkout}
              disabled={startSession.isPending}
              className="
                w-full flex items-center justify-center gap-3
                py-4 px-6 rounded-xl
                bg-accent text-background font-black text-lg uppercase tracking-wide
                hover:bg-accent-light active:scale-[0.98]
                disabled:opacity-70 disabled:cursor-not-allowed
                transition-all duration-200
                shadow-lg shadow-accent/25
              "
            >
              {startSession.isPending ? (
                <>
                  <Loader2 size={22} className="animate-spin" />
                  Starting...
                </>
              ) : (
                <>
                  <Play size={22} />
                  Start Workout
                </>
              )}
            </button>
          </div>
        </Container>
      </div>
    );
  }

  // Active workout session
  return (
    <div className="py-6 pb-32">
      <Container>
        {/* Header with progress */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <p className="text-xs uppercase tracking-wider text-muted">
                Week {currentWeek} · Day {currentDayIndex + 1}
              </p>
              <h1 className="text-xl font-bold text-foreground">
                {workout?.daySlug?.replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase()) ?? 'Workout'}
              </h1>
            </div>
            <div className="text-right">
              <p className="text-2xl font-black tabular-nums text-accent">
                {completedCount}/{totalSets}
              </p>
              <p className="text-xs text-muted uppercase tracking-wider">Sets</p>
            </div>
          </div>

          {/* Progress bar */}
          <div className="h-2 bg-surface-elevated rounded-full overflow-hidden">
            <div
              className="h-full bg-accent rounded-full transition-all duration-500 ease-out"
              style={{ width: `${progressPercentage}%` }}
            />
          </div>
        </div>

        {/* Exercise cards */}
        <div className="space-y-4">
          {workout?.exercises.map((exercise, index) => (
            <ExerciseCard
              key={`${exercise.prescriptionId}-${index}`}
              exercise={exercise}
              exerciseIndex={index}
              completedSets={completedSets}
              onSetToggle={handleSetToggle}
            />
          ))}
        </div>
      </Container>

      {/* Floating action bar */}
      <div className="fixed bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-background via-background to-transparent pt-8">
        <Container>
          <div className="flex gap-3">
            <button
              onClick={() => setActiveDialog('abandon')}
              className="
                flex-shrink-0 p-4 rounded-xl
                bg-surface border border-border text-muted
                hover:bg-error/10 hover:border-error/30 hover:text-error
                active:scale-[0.98]
                transition-all duration-200
              "
            >
              <XCircle size={24} />
            </button>
            <button
              onClick={() => setActiveDialog('finish')}
              disabled={finishSession.isPending}
              className="
                flex-1 flex items-center justify-center gap-3
                py-4 px-6 rounded-xl
                bg-success text-background font-black text-lg uppercase tracking-wide
                hover:bg-success/90 active:scale-[0.98]
                disabled:opacity-70 disabled:cursor-not-allowed
                transition-all duration-200
                shadow-lg shadow-success/25
              "
            >
              {finishSession.isPending ? (
                <>
                  <Loader2 size={22} className="animate-spin" />
                  Finishing...
                </>
              ) : (
                <>
                  <Flag size={22} />
                  Finish Workout
                </>
              )}
            </button>
          </div>
        </Container>
      </div>

      {/* Dialogs */}
      <ConfirmDialog
        isOpen={activeDialog === 'abandon'}
        title="Abandon Workout?"
        message="Your progress for this session will not be saved. You can start the workout again later."
        confirmLabel="Abandon"
        cancelLabel="Keep Going"
        variant="destructive"
        isLoading={abandonSession.isPending}
        onConfirm={handleAbandonWorkout}
        onCancel={() => setActiveDialog(null)}
      />

      <ConfirmDialog
        isOpen={activeDialog === 'finish'}
        title="Finish Workout?"
        message={`You've completed ${completedCount} of ${totalSets} sets. Ready to finish this session?`}
        confirmLabel="Finish"
        cancelLabel="Not Yet"
        variant="default"
        isLoading={finishSession.isPending}
        onConfirm={handleFinishWorkout}
        onCancel={() => setActiveDialog(null)}
      />

      <ConfirmDialog
        isOpen={activeDialog === 'advanceWeek'}
        title="Advance to Next Week?"
        message={`Great work completing Week ${currentWeek}! Ready to move on to Week ${currentWeek + 1}?`}
        confirmLabel="Advance"
        cancelLabel="Stay Here"
        variant="default"
        isLoading={advanceState.isPending}
        onConfirm={handleAdvanceWeek}
        onCancel={() => setActiveDialog(null)}
      />

      <ConfirmDialog
        isOpen={activeDialog === 'startNextCycle'}
        title="Start Next Cycle?"
        message="You've completed this training cycle! Starting a new cycle will apply any programmed progressions to your training maxes."
        confirmLabel="Start Cycle"
        cancelLabel="Not Yet"
        variant="default"
        isLoading={advanceState.isPending}
        onConfirm={handleAdvanceWeek}
        onCancel={() => setActiveDialog(null)}
      />
    </div>
  );
}
