import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Container } from '../../components/layout';
import { useProgram } from '../../hooks/useProgram';
import { useEnrollment, useEnrollInProgram, useUnenroll } from '../../hooks/useCurrentUser';
import { useAuth } from '../../contexts/useAuth';
import { ConfirmDialog } from '../../components/ui/ConfirmDialog';
import {
  ArrowLeft,
  AlertCircle,
  Calendar,
  Repeat,
  CheckCircle2,
  Zap,
  ArrowRight,
  LogOut,
} from 'lucide-react';
import type { ProgramDetail, EnrollmentProgram } from '../../api/types';

function LoadingSkeleton() {
  return (
    <div className="animate-fade-in">
      {/* Back button skeleton */}
      <div className="h-5 w-32 bg-surface-elevated rounded mb-8 animate-pulse" />

      {/* Header skeleton */}
      <div className="mb-8">
        <div className="h-10 w-3/4 bg-surface-elevated rounded mb-4 animate-pulse" />
        <div className="h-5 w-full bg-surface-elevated rounded mb-2 animate-pulse" />
        <div className="h-5 w-2/3 bg-surface-elevated rounded animate-pulse" />
      </div>

      {/* Stats skeleton */}
      <div className="grid grid-cols-2 gap-4 mb-8">
        {[1, 2].map((i) => (
          <div
            key={i}
            className="bg-surface border border-border rounded-lg p-5 animate-pulse"
          >
            <div className="h-4 w-20 bg-surface-elevated rounded mb-3" />
            <div className="h-7 w-16 bg-surface-elevated rounded" />
          </div>
        ))}
      </div>

      {/* Weeks skeleton */}
      <div className="bg-surface border border-border rounded-lg p-5 animate-pulse">
        <div className="h-5 w-32 bg-surface-elevated rounded mb-6" />
        <div className="flex gap-3">
          {[1, 2, 3, 4].map((i) => (
            <div
              key={i}
              className="flex-1 h-16 bg-surface-elevated rounded"
            />
          ))}
        </div>
      </div>

      {/* CTA skeleton */}
      <div className="mt-8 h-14 bg-surface-elevated rounded-lg animate-pulse" />
    </div>
  );
}

function NotFoundState() {
  return (
    <div className="text-center py-16 animate-fade-in">
      <div className="mx-auto w-20 h-20 rounded-full bg-surface border-2 border-border flex items-center justify-center mb-6">
        <AlertCircle className="w-10 h-10 text-muted" />
      </div>
      <h2 className="text-2xl font-bold text-foreground mb-3">
        Program Not Found
      </h2>
      <p className="text-muted max-w-sm mx-auto mb-8">
        This program doesn't exist or may have been removed.
      </p>
      <Link
        to="/programs"
        className="
          inline-flex items-center gap-2
          px-6 py-3
          bg-surface border border-border rounded-lg
          text-foreground font-medium
          hover:border-accent/30 hover:bg-surface-elevated
          transition-all duration-200
        "
      >
        <ArrowLeft size={18} />
        Back to Programs
      </Link>
    </div>
  );
}

function ErrorState({ message }: { message: string }) {
  return (
    <div className="text-center py-16 animate-fade-in">
      <div className="mx-auto w-20 h-20 rounded-full bg-error/10 border-2 border-error/20 flex items-center justify-center mb-6">
        <AlertCircle className="w-10 h-10 text-error" />
      </div>
      <h2 className="text-2xl font-bold text-foreground mb-3">
        Failed to Load Program
      </h2>
      <p className="text-muted max-w-sm mx-auto mb-8">{message}</p>
      <Link
        to="/programs"
        className="
          inline-flex items-center gap-2
          px-6 py-3
          bg-surface border border-border rounded-lg
          text-foreground font-medium
          hover:border-accent/30 hover:bg-surface-elevated
          transition-all duration-200
        "
      >
        <ArrowLeft size={18} />
        Back to Programs
      </Link>
    </div>
  );
}

interface WeekBlockProps {
  weekNumber: number;
  totalWeeks: number;
  index: number;
}

function WeekBlock({ weekNumber, totalWeeks, index }: WeekBlockProps) {
  // Calculate intensity visual - later weeks appear "heavier"
  const intensity = (weekNumber / totalWeeks) * 100;
  const delay = index * 50;

  return (
    <div
      className="
        group relative flex-1 min-w-0
        animate-slide-up
      "
      style={{ animationDelay: `${delay}ms` }}
    >
      {/* Week bar */}
      <div
        className="
          relative overflow-hidden
          bg-surface border border-border rounded-lg
          p-4 h-full min-h-[80px]
          hover:border-accent/40
          transition-all duration-300
        "
      >
        {/* Intensity fill */}
        <div
          className="
            absolute bottom-0 left-0 right-0
            bg-gradient-to-t from-accent/20 to-transparent
            transition-all duration-500
            group-hover:from-accent/30
          "
          style={{ height: `${Math.max(30, intensity)}%` }}
        />

        {/* Content */}
        <div className="relative z-10 flex flex-col items-center justify-center h-full">
          <span className="text-xs uppercase tracking-widest text-muted mb-1">
            Week
          </span>
          <span className="text-2xl font-bold text-foreground tabular-nums">
            {weekNumber}
          </span>
        </div>
      </div>
    </div>
  );
}

type DialogType = 'enroll' | 'switch' | 'unenroll' | null;

interface ProgramContentProps {
  program: ProgramDetail;
  enrollmentStatus: 'not-enrolled' | 'enrolled-this' | 'enrolled-other';
  currentProgram?: EnrollmentProgram;
  onEnroll: () => void;
  onUnenroll: () => void;
  onSwitch: () => void;
  isEnrolling: boolean;
  isUnenrolling: boolean;
  enrollError: Error | null;
  unenrollError: Error | null;
}

function ProgramContent({
  program,
  enrollmentStatus,
  currentProgram,
  onEnroll,
  onUnenroll,
  onSwitch,
  isEnrolling,
  isUnenrolling,
  enrollError,
  unenrollError,
}: ProgramContentProps) {
  const weeks = program.cycle.weeks;
  const [openDialog, setOpenDialog] = useState<DialogType>(null);
  const [showSuccess, setShowSuccess] = useState(false);

  const handleConfirmEnroll = () => {
    onEnroll();
    setOpenDialog(null);
    setShowSuccess(true);
    setTimeout(() => setShowSuccess(false), 3000);
  };

  const handleConfirmSwitch = () => {
    onSwitch();
    setOpenDialog(null);
    setShowSuccess(true);
    setTimeout(() => setShowSuccess(false), 3000);
  };

  const handleConfirmUnenroll = () => {
    onUnenroll();
    setOpenDialog(null);
  };

  const error = enrollError || unenrollError;

  return (
    <div className="animate-fade-in">
      {/* Back link */}
      <Link
        to="/programs"
        className="
          inline-flex items-center gap-2
          text-sm text-muted
          hover:text-accent
          transition-colors duration-200
          mb-6
        "
      >
        <ArrowLeft size={16} />
        <span>All Programs</span>
      </Link>

      {/* Program header */}
      <header className="mb-8">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-3">
          {program.name}
        </h1>
        {program.description && (
          <p className="text-lg text-muted leading-relaxed max-w-2xl">
            {program.description}
          </p>
        )}
      </header>

      {/* Stats grid */}
      <div className="grid grid-cols-2 gap-4 mb-8">
        {/* Cycle info */}
        <div className="bg-surface border border-border rounded-lg p-5">
          <div className="flex items-center gap-2 text-muted mb-2">
            <Repeat size={16} />
            <span className="text-xs uppercase tracking-wider">Cycle</span>
          </div>
          <div className="text-xl font-semibold text-foreground">
            {program.cycle.name}
          </div>
        </div>

        {/* Duration */}
        <div className="bg-surface border border-border rounded-lg p-5">
          <div className="flex items-center gap-2 text-muted mb-2">
            <Calendar size={16} />
            <span className="text-xs uppercase tracking-wider">Duration</span>
          </div>
          <div className="text-xl font-semibold text-foreground">
            <span className="tabular-nums">{program.cycle.lengthWeeks}</span>
            <span className="text-muted font-normal ml-1">weeks</span>
          </div>
        </div>
      </div>

      {/* Week breakdown */}
      <section className="bg-surface border border-border rounded-lg p-5 mb-8">
        <h2 className="text-sm uppercase tracking-wider text-muted mb-5 flex items-center gap-2">
          <span className="w-1 h-4 bg-accent rounded-full" />
          Program Structure
        </h2>

        {weeks.length > 0 ? (
          <div className="flex gap-2 sm:gap-3 overflow-x-auto pb-2">
            {weeks.map((week, index) => (
              <WeekBlock
                key={week.id}
                weekNumber={week.weekNumber}
                totalWeeks={program.cycle.lengthWeeks}
                index={index}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-8 text-muted">
            <Calendar className="w-8 h-8 mx-auto mb-2 opacity-50" />
            <p>Week structure not yet defined</p>
          </div>
        )}
      </section>

      {/* Additional info */}
      {(program.weeklyLookup || program.dailyLookup || program.defaultRounding) && (
        <section className="bg-surface border border-border rounded-lg p-5 mb-8">
          <h2 className="text-sm uppercase tracking-wider text-muted mb-4 flex items-center gap-2">
            <span className="w-1 h-4 bg-accent rounded-full" />
            Program Details
          </h2>
          <dl className="grid gap-3 text-sm">
            {program.weeklyLookup && (
              <div className="flex justify-between items-center py-2 border-b border-border">
                <dt className="text-muted">Weekly Variation</dt>
                <dd className="font-medium text-foreground">{program.weeklyLookup.name}</dd>
              </div>
            )}
            {program.dailyLookup && (
              <div className="flex justify-between items-center py-2 border-b border-border">
                <dt className="text-muted">Daily Variation</dt>
                <dd className="font-medium text-foreground">{program.dailyLookup.name}</dd>
              </div>
            )}
            {program.defaultRounding && (
              <div className="flex justify-between items-center py-2">
                <dt className="text-muted">Weight Rounding</dt>
                <dd className="font-medium text-foreground">{program.defaultRounding} lbs</dd>
              </div>
            )}
          </dl>
        </section>
      )}

      {/* Success message */}
      {showSuccess && (
        <div className="mb-4 p-4 bg-success/10 border border-success/30 rounded-lg text-success text-center animate-fade-in">
          <CheckCircle2 className="inline-block w-5 h-5 mr-2" />
          Successfully enrolled in {program.name}!
        </div>
      )}

      {/* Error message */}
      {error && (
        <div className="mb-4 p-4 bg-error/10 border border-error/30 rounded-lg text-error text-center animate-fade-in">
          <AlertCircle className="inline-block w-5 h-5 mr-2" />
          {error.message || 'An error occurred. Please try again.'}
        </div>
      )}

      {/* Enrollment CTA */}
      <div className="mt-8">
        {enrollmentStatus === 'enrolled-this' ? (
          <div className="space-y-3">
            <div
              className="
                flex items-center justify-center gap-3
                w-full py-4 px-6
                bg-success/10 border-2 border-success/30 rounded-lg
                text-success font-semibold
              "
            >
              <CheckCircle2 size={22} />
              <span>Currently Enrolled</span>
            </div>
            <button
              type="button"
              onClick={() => setOpenDialog('unenroll')}
              disabled={isUnenrolling}
              className="
                flex items-center justify-center gap-2
                w-full py-3 px-6
                bg-surface border border-border rounded-lg
                text-muted font-medium text-sm
                hover:border-error/50 hover:text-error
                transition-all duration-200
                cursor-pointer
                disabled:opacity-50 disabled:cursor-not-allowed
              "
            >
              <LogOut size={16} />
              <span>Unenroll from Program</span>
            </button>
          </div>
        ) : enrollmentStatus === 'enrolled-other' ? (
          <button
            type="button"
            onClick={() => setOpenDialog('switch')}
            disabled={isEnrolling || isUnenrolling}
            className="
              group flex items-center justify-center gap-3
              w-full py-4 px-6
              bg-surface border-2 border-accent/50 rounded-lg
              text-accent font-semibold
              hover:bg-accent hover:text-background hover:border-accent
              transition-all duration-300
              cursor-pointer
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            {isEnrolling || isUnenrolling ? (
              <>
                <span className="w-5 h-5 border-2 border-current border-t-transparent rounded-full animate-spin" />
                <span>Switching...</span>
              </>
            ) : (
              <>
                <Zap size={20} className="group-hover:animate-pulse" />
                <span>Switch to this Program</span>
                <ArrowRight size={18} className="group-hover:translate-x-1 transition-transform" />
              </>
            )}
          </button>
        ) : (
          <button
            type="button"
            onClick={() => setOpenDialog('enroll')}
            disabled={isEnrolling}
            className="
              group flex items-center justify-center gap-3
              w-full py-4 px-6
              bg-accent rounded-lg
              text-background font-bold text-lg
              hover:bg-accent-light
              active:scale-[0.98]
              transition-all duration-200
              cursor-pointer
              disabled:opacity-70 disabled:cursor-not-allowed
            "
          >
            {isEnrolling ? (
              <>
                <span className="w-5 h-5 border-2 border-current border-t-transparent rounded-full animate-spin" />
                <span>Enrolling...</span>
              </>
            ) : (
              <>
                <Zap size={22} />
                <span>Enroll in Program</span>
                <ArrowRight size={20} className="group-hover:translate-x-1 transition-transform" />
              </>
            )}
          </button>
        )}
      </div>

      {/* Confirmation Dialogs */}
      <ConfirmDialog
        isOpen={openDialog === 'enroll'}
        title="Start Training"
        message={`Begin your training with ${program.name}? You can track your workouts and progression once enrolled.`}
        confirmLabel="Start Training"
        cancelLabel="Not Yet"
        variant="default"
        isLoading={isEnrolling}
        onConfirm={handleConfirmEnroll}
        onCancel={() => setOpenDialog(null)}
      />

      <ConfirmDialog
        isOpen={openDialog === 'switch'}
        title="Switch Programs"
        message={`Switch from ${currentProgram?.name || 'your current program'} to ${program.name}? Your current enrollment progress will be reset.`}
        confirmLabel="Switch Program"
        cancelLabel="Keep Current"
        variant="destructive"
        isLoading={isEnrolling || isUnenrolling}
        onConfirm={handleConfirmSwitch}
        onCancel={() => setOpenDialog(null)}
      />

      <ConfirmDialog
        isOpen={openDialog === 'unenroll'}
        title="Unenroll from Program"
        message={`Are you sure you want to unenroll from ${program.name}? Your enrollment progress will be lost.`}
        confirmLabel="Unenroll"
        cancelLabel="Stay Enrolled"
        variant="destructive"
        isLoading={isUnenrolling}
        onConfirm={handleConfirmUnenroll}
        onCancel={() => setOpenDialog(null)}
      />
    </div>
  );
}

export default function ProgramDetails() {
  const { id } = useParams<{ id: string }>();
  const { userId } = useAuth();

  const { data: program, isLoading, error } = useProgram(id);
  const { data: enrollmentData } = useEnrollment(userId ?? undefined);

  const enrollMutation = useEnrollInProgram(userId ?? undefined);
  const unenrollMutation = useUnenroll(userId ?? undefined);

  // Determine enrollment status
  let enrollmentStatus: 'not-enrolled' | 'enrolled-this' | 'enrolled-other' = 'not-enrolled';
  let currentProgram: EnrollmentProgram | undefined;

  if (enrollmentData) {
    currentProgram = enrollmentData.program;
    const currentProgramId = enrollmentData.program.id;
    if (currentProgramId === id) {
      enrollmentStatus = 'enrolled-this';
    } else {
      enrollmentStatus = 'enrolled-other';
    }
  }

  const handleEnroll = () => {
    if (id) {
      enrollMutation.mutate({ programId: id });
    }
  };

  const handleUnenroll = () => {
    unenrollMutation.mutate();
  };

  const handleSwitch = async () => {
    // For switching, we unenroll first then enroll in the new program
    // The backend might handle this atomically in the future, but for now we chain them
    if (id) {
      await unenrollMutation.mutateAsync();
      enrollMutation.mutate({ programId: id });
    }
  };

  return (
    <div className="py-6 md:py-8 min-h-screen">
      <Container>
        {isLoading && <LoadingSkeleton />}

        {error && (
          <ErrorState
            message={error instanceof Error ? error.message : 'An unexpected error occurred'}
          />
        )}

        {!isLoading && !error && !program && <NotFoundState />}

        {!isLoading && !error && program && (
          <ProgramContent
            program={program}
            enrollmentStatus={enrollmentStatus}
            currentProgram={currentProgram}
            onEnroll={handleEnroll}
            onUnenroll={handleUnenroll}
            onSwitch={handleSwitch}
            isEnrolling={enrollMutation.isPending}
            isUnenrolling={unenrollMutation.isPending}
            enrollError={enrollMutation.error}
            unenrollError={unenrollMutation.error}
          />
        )}
      </Container>
    </div>
  );
}
