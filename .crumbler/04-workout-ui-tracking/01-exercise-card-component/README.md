# Exercise Card Component

Create a reusable ExerciseCard component for displaying workout exercises with their sets.

## Component: `ExerciseCard`

**Location:** `frontend/src/components/workout/ExerciseCard.tsx`

## Props Interface

```typescript
interface ExerciseCardProps {
  exercise: WorkoutExercise;  // From src/api/types.ts
  exerciseIndex: number;
  completedSets: Set<string>; // Set of "exerciseIndex-setNumber" keys
  onSetToggle: (exerciseIndex: number, setNumber: number) => void;
  isReadOnly?: boolean;       // For history view - no checkboxes
}
```

## Features

1. **Exercise Header**
   - Exercise name (prominent)
   - Notes display if present
   - Rest period indicator if restSeconds exists

2. **Sets Display**
   - Each set shows: set number, weight (lbs), target reps
   - Work sets vs warmup sets visually differentiated (isWorkSet field)
   - Checkbox for marking set completion (unless readOnly)
   - Completed sets have visual feedback (strikethrough or dimmed)

3. **Visual Design**
   - Dark surface card with border
   - Large touch targets for checkboxes (min 44px)
   - Clear weight/reps typography (tabular-nums for alignment)
   - Accent color for work sets
   - Muted styling for warmup sets

## API Types Reference

```typescript
interface WorkoutExercise {
  prescriptionId: string;
  lift: { id: string; name: string; slug: string; };
  sets: GeneratedSet[];
  notes?: string;
  restSeconds?: number;
}

interface GeneratedSet {
  setNumber: number;
  weight: number;
  targetReps: number;
  isWorkSet: boolean;
}
```

## Implementation Notes

- Use /frontend-design skill for styling
- Checkboxes should be easily tappable during workout
- Weight should be formatted with "lbs" suffix
- Set completion state is managed by parent component
