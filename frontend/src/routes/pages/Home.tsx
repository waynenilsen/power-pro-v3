import { Link } from 'react-router-dom';
import { Container } from '../../components/layout';
import { useAuth } from '../../contexts/useAuth';
import { useEnrollment } from '../../hooks/useCurrentUser';
import { BookOpen, ChevronRight, Dumbbell, Play, Loader2 } from 'lucide-react';

function EnrolledProgramCard({
  programName,
  programSlug,
  currentWeek,
  cycleLengthWeeks,
  currentCycleIteration,
}: {
  programName: string;
  programSlug: string;
  currentWeek: number;
  cycleLengthWeeks: number;
  currentCycleIteration: number;
}) {
  const progressPercentage = (currentWeek / cycleLengthWeeks) * 100;

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

        {/* View details link */}
        <Link
          to={`/programs/${programSlug}`}
          className="text-sm text-muted hover:text-accent transition-colors"
        >
          View program details →
        </Link>
      </div>

      {/* Start workout CTA */}
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
        <span>Start Workout</span>
      </Link>
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

  const isEnrolled = enrollment?.data && enrollment.data.enrollmentStatus === 'ACTIVE';

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
          ) : isEnrolled && enrollment.data ? (
            <EnrolledProgramCard
              programName={enrollment.data.program.name}
              programSlug={enrollment.data.program.slug}
              currentWeek={enrollment.data.state.currentWeek}
              cycleLengthWeeks={enrollment.data.program.cycleLengthWeeks}
              currentCycleIteration={enrollment.data.state.currentCycleIteration}
            />
          ) : (
            <BrowseProgramsCard />
          )}
        </div>
      </Container>
    </div>
  );
}
