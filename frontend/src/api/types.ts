// Common types

export interface ApiErrorBody {
  code: string;
  message: string;
  details?: {
    validationErrors?: string[];
  };
}

export interface ApiError {
  error: ApiErrorBody;
}

export interface PaginationMeta {
  total: number;
  limit: number;
  offset: number;
  hasMore: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: PaginationMeta;
}

export interface PaginationParams {
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// Health

export interface HealthResponse {
  status: string;
}

// Lifts

export interface Lift {
  id: string;
  name: string;
  slug: string;
  isCompetitionLift: boolean;
  parentLiftId: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateLiftRequest {
  name: string;
  slug?: string;
  isCompetitionLift?: boolean;
  parentLiftId?: string;
}

export interface UpdateLiftRequest {
  name?: string;
  slug?: string;
  isCompetitionLift?: boolean;
  parentLiftId?: string;
  clearParentLift?: boolean;
}

// Lift Maxes

export type MaxType = 'ONE_RM' | 'TRAINING_MAX';

export interface LiftMax {
  id: string;
  userId: string;
  liftId: string;
  type: MaxType;
  value: number;
  effectiveDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface LiftMaxWithWarnings {
  data: LiftMax;
  warnings?: string[];
}

export interface CreateLiftMaxRequest {
  liftId: string;
  type: MaxType;
  value: number;
  effectiveDate?: string;
}

export interface UpdateLiftMaxRequest {
  value?: number;
  effectiveDate?: string;
}

export interface ConversionResult {
  originalValue: number;
  originalType: MaxType;
  convertedValue: number;
  convertedType: MaxType;
  percentage: number;
}

// Load Strategy

export interface LoadStrategy {
  type: 'PERCENT_OF';
  maxType: MaxType;
  percentage: number;
  lookupKey?: string;
  roundTo?: number;
}

// Set Schemes

export interface FixedSetScheme {
  type: 'FIXED';
  sets: number;
  reps: number;
  isAmrap?: boolean;
}

export interface RampSetScheme {
  type: 'RAMP';
  sets: Array<{
    percentage: number;
    reps: number;
  }>;
}

export type SetScheme = FixedSetScheme | RampSetScheme;

// Prescriptions

export interface Prescription {
  id: string;
  liftId: string;
  loadStrategy: LoadStrategy;
  setScheme: SetScheme;
  order: number;
  notes?: string;
  restSeconds?: number;
  createdAt: string;
  updatedAt: string;
}

export interface GeneratedSet {
  setNumber: number;
  weight: number;
  targetReps: number;
  isWorkSet: boolean;
}

export interface ResolvedPrescription {
  prescriptionId: string;
  lift: {
    id: string;
    name: string;
    slug: string;
  };
  sets: GeneratedSet[];
  notes?: string;
  restSeconds?: number;
}

// Days

export interface Day {
  id: string;
  name: string;
  slug: string;
  metadata?: Record<string, unknown>;
  programId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface DayPrescription {
  id: string;
  prescriptionId: string;
  order: number;
  createdAt: string;
}

export interface DayWithPrescriptions extends Day {
  prescriptions: DayPrescription[];
}

// Weeks

export interface Week {
  id: string;
  cycleId: string;
  weekNumber: number;
  name?: string;
  createdAt: string;
  updatedAt: string;
}

export interface WeekDay {
  id: string;
  dayId: string;
  position: number;
}

export interface WeekWithDays extends Week {
  days: WeekDay[];
}

// Cycles

export interface Cycle {
  id: string;
  name: string;
  lengthWeeks: number;
  createdAt: string;
  updatedAt: string;
}

// Lookups

export interface WeeklyLookup {
  id: string;
  name: string;
  entries: Record<string, number>;
  createdAt: string;
  updatedAt: string;
}

export interface DailyLookup {
  id: string;
  name: string;
  entries: Record<string, number>;
  createdAt: string;
  updatedAt: string;
}

// Programs

export interface ProgramListItem {
  id: string;
  name: string;
  slug: string;
  description?: string;
  cycleId: string;
  weeklyLookupId?: string;
  dailyLookupId?: string;
  defaultRounding?: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProgramDetail {
  id: string;
  name: string;
  slug: string;
  description?: string;
  cycle: {
    id: string;
    name: string;
    lengthWeeks: number;
    weeks: Array<{
      id: string;
      weekNumber: number;
    }>;
  };
  weeklyLookup?: {
    id: string;
    name: string;
  };
  dailyLookup?: {
    id: string;
    name: string;
  };
  defaultRounding?: number;
  createdAt: string;
  updatedAt: string;
}

// Progressions

export type ProgressionType = 'LINEAR' | 'CYCLE';

export interface Progression {
  id: string;
  name: string;
  type: ProgressionType;
  parameters: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface ProgramProgression {
  id: string;
  programId: string;
  progressionId: string;
  liftId?: string;
  priority: number;
  createdAt: string;
  updatedAt: string;
}

// Enrollment

export type EnrollmentStatus = 'ACTIVE' | 'BETWEEN_CYCLES' | 'QUIT';
export type CycleStatus = 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
export type WeekStatus = 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
export type WorkoutSessionStatus = 'IN_PROGRESS' | 'COMPLETED' | 'ABANDONED';

export interface EnrollmentState {
  currentWeek: number;
  currentCycleIteration: number;
  currentDayIndex?: number;
}

export interface EnrollmentProgram {
  id: string;
  name: string;
  slug: string;
  description?: string;
  cycleLengthWeeks: number;
  daysPerWeek: number;
}

export interface Enrollment {
  id: string;
  userId: string;
  program: EnrollmentProgram;
  state: EnrollmentState;
  enrolledAt: string;
  updatedAt: string;
}

export interface WorkoutSessionSummary {
  id: string;
  weekNumber: number;
  dayIndex: number;
  status: WorkoutSessionStatus;
  startedAt: string;
  finishedAt?: string;
}

export type EnrollmentResponse = Enrollment & {
  enrollmentStatus: EnrollmentStatus;
  cycleStatus: CycleStatus;
  weekStatus: WeekStatus;
  currentWorkoutSession?: WorkoutSessionSummary;
};

export interface EnrollRequest {
  programId: string;
}

// Workout Sessions

export interface WorkoutSession {
  id: string;
  userProgramStateId: string;
  weekNumber: number;
  dayIndex: number;
  status: WorkoutSessionStatus;
  startedAt: string;
  finishedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export type WorkoutSessionDataResponse = WorkoutSession;

export interface WorkoutSessionListResponse {
  data: WorkoutSession[];
  meta: PaginationMeta;
}

// Workouts

export interface WorkoutExercise {
  prescriptionId: string;
  lift: {
    id: string;
    name: string;
    slug: string;
  };
  sets: GeneratedSet[];
  notes?: string;
  restSeconds?: number;
}

export interface Workout {
  userId: string;
  programId: string;
  cycleIteration: number;
  weekNumber: number;
  daySlug: string;
  date: string;
  exercises: WorkoutExercise[];
}

// State Advancement

export type AdvanceType = 'day' | 'week';

export interface AdvanceStateRequest {
  advanceType: AdvanceType;
}

// Progression History

export interface ProgressionHistoryEntry {
  id: string;
  userId: string;
  liftId: string;
  progressionId: string;
  previousValue: number;
  newValue: number;
  delta: number;
  appliedAt: string;
  triggeredBy?: string;
  cycleIteration?: number;
  weekNumber?: number;
}

// Manual Progression

export interface ManualTriggerRequest {
  progressionId: string;
  liftId?: string;
  force?: boolean;
}

export interface ManualTriggerResult {
  progressionId: string;
  liftId: string;
  applied: boolean;
  skipped: boolean;
  result?: {
    previousValue: number;
    newValue: number;
    delta: number;
    maxType: string;
    appliedAt: string;
  };
}

export interface ManualTriggerResponse {
  results: ManualTriggerResult[];
  totalApplied: number;
  totalSkipped: number;
  totalErrors: number;
}
