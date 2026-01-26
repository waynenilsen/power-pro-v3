import { Check, Clock, MessageSquare } from 'lucide-react';
import type { WorkoutExercise, GeneratedSet } from '../../api/types';

interface ExerciseCardProps {
  exercise: WorkoutExercise;
  exerciseIndex: number;
  completedSets: Set<string>;
  onSetToggle: (exerciseIndex: number, setNumber: number) => void;
  isReadOnly?: boolean;
}

interface SetRowProps {
  set: GeneratedSet;
  exerciseIndex: number;
  isCompleted: boolean;
  onToggle: () => void;
  isReadOnly?: boolean;
}

function SetRow({ set, isCompleted, onToggle, isReadOnly }: SetRowProps) {
  const isWorkSet = set.isWorkSet;

  return (
    <div
      className={`
        group flex items-center gap-3 py-3 px-4
        border-b border-border/50 last:border-b-0
        transition-all duration-150
        ${isCompleted ? 'bg-success/5' : ''}
        ${!isReadOnly && !isCompleted ? 'hover:bg-surface-elevated' : ''}
      `}
    >
      {/* Checkbox / Completion indicator */}
      {!isReadOnly && (
        <button
          onClick={onToggle}
          className={`
            relative flex items-center justify-center
            w-11 h-11 rounded-lg
            border-2 transition-all duration-200
            focus:outline-none focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 focus-visible:ring-offset-surface
            ${isCompleted
              ? 'bg-success border-success text-background'
              : isWorkSet
                ? 'border-accent/40 hover:border-accent hover:bg-accent/10'
                : 'border-border hover:border-muted hover:bg-surface-elevated'
            }
            active:scale-95
          `}
          aria-label={isCompleted ? 'Mark set incomplete' : 'Mark set complete'}
        >
          {isCompleted && <Check className="w-5 h-5" strokeWidth={3} />}
        </button>
      )}

      {/* Set number indicator */}
      <div
        className={`
          w-8 h-8 rounded-md flex items-center justify-center text-sm font-bold
          ${isWorkSet
            ? 'bg-accent/15 text-accent border border-accent/30'
            : 'bg-surface-elevated text-muted border border-border'
          }
          ${isCompleted ? 'opacity-50' : ''}
        `}
      >
        {set.setNumber}
      </div>

      {/* Weight */}
      <div className={`flex-1 ${isCompleted ? 'opacity-50' : ''}`}>
        <span
          className={`
            text-lg font-bold tabular-nums tracking-tight
            ${isWorkSet ? 'text-foreground' : 'text-muted'}
            ${isCompleted ? 'line-through decoration-success/50' : ''}
          `}
        >
          {set.weight}
        </span>
        <span className={`text-sm ml-1 ${isWorkSet ? 'text-muted' : 'text-muted/60'}`}>
          lbs
        </span>
      </div>

      {/* Reps */}
      <div className={`text-right ${isCompleted ? 'opacity-50' : ''}`}>
        <span
          className={`
            text-lg font-bold tabular-nums
            ${isWorkSet ? 'text-foreground' : 'text-muted'}
            ${isCompleted ? 'line-through decoration-success/50' : ''}
          `}
        >
          {set.targetReps}
        </span>
        <span className={`text-sm ml-1 ${isWorkSet ? 'text-muted' : 'text-muted/60'}`}>
          reps
        </span>
      </div>

      {/* Set type badge */}
      <div
        className={`
          hidden sm:block w-16 text-right text-xs font-medium uppercase tracking-wider
          ${isWorkSet ? 'text-accent' : 'text-muted/50'}
          ${isCompleted ? 'opacity-50' : ''}
        `}
      >
        {isWorkSet ? 'Work' : 'Warm'}
      </div>
    </div>
  );
}

export function ExerciseCard({
  exercise,
  exerciseIndex,
  completedSets,
  onSetToggle,
  isReadOnly = false,
}: ExerciseCardProps) {
  const totalSets = exercise.sets.length;
  const completedCount = exercise.sets.filter(
    set => completedSets.has(`${exerciseIndex}-${set.setNumber}`)
  ).length;
  const allCompleted = completedCount === totalSets;

  return (
    <div
      className={`
        bg-surface border rounded-xl overflow-hidden
        transition-all duration-300
        ${allCompleted
          ? 'border-success/30 shadow-[0_0_20px_-5px_rgba(34,197,94,0.15)]'
          : 'border-border hover:border-border-accent'
        }
      `}
    >
      {/* Header */}
      <div className="px-4 py-4 border-b border-border bg-surface-elevated/50">
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1 min-w-0">
            {/* Exercise name */}
            <h3
              className={`
                text-lg font-bold tracking-tight
                ${allCompleted ? 'text-success' : 'text-foreground'}
                transition-colors duration-300
              `}
            >
              {exercise.lift.name}
            </h3>

            {/* Progress indicator */}
            <div className="flex items-center gap-2 mt-1">
              <span className="text-sm text-muted tabular-nums">
                {completedCount}/{totalSets} sets
              </span>
              {allCompleted && (
                <span className="text-xs font-medium text-success uppercase tracking-wider animate-fade-in">
                  Complete
                </span>
              )}
            </div>
          </div>

          {/* Rest timer badge */}
          {exercise.restSeconds && (
            <div className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg bg-surface border border-border text-muted">
              <Clock className="w-3.5 h-3.5" />
              <span className="text-xs font-medium tabular-nums">
                {Math.floor(exercise.restSeconds / 60)}:{String(exercise.restSeconds % 60).padStart(2, '0')}
              </span>
            </div>
          )}
        </div>

        {/* Notes */}
        {exercise.notes && (
          <div className="flex items-start gap-2 mt-3 p-2.5 rounded-lg bg-surface border border-border/50">
            <MessageSquare className="w-4 h-4 text-muted flex-shrink-0 mt-0.5" />
            <p className="text-sm text-muted leading-relaxed">{exercise.notes}</p>
          </div>
        )}
      </div>

      {/* Sets list */}
      <div className="divide-y divide-border/30">
        {exercise.sets.map((set) => {
          const setKey = `${exerciseIndex}-${set.setNumber}`;
          const isCompleted = completedSets.has(setKey);

          return (
            <SetRow
              key={setKey}
              set={set}
              exerciseIndex={exerciseIndex}
              isCompleted={isCompleted}
              onToggle={() => onSetToggle(exerciseIndex, set.setNumber)}
              isReadOnly={isReadOnly}
            />
          );
        })}
      </div>

      {/* Progress bar at bottom */}
      <div className="h-1 bg-surface-elevated">
        <div
          className={`
            h-full transition-all duration-500 ease-out
            ${allCompleted ? 'bg-success' : 'bg-accent'}
          `}
          style={{ width: `${(completedCount / totalSets) * 100}%` }}
        />
      </div>
    </div>
  );
}

export default ExerciseCard;
