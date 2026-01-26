import { Link } from 'react-router-dom';
import { Container } from '../../components/layout';
import { usePrograms } from '../../hooks/usePrograms';
import { BookOpen, ChevronRight, AlertCircle } from 'lucide-react';
import type { ProgramListItem } from '../../api/types';

function ProgramCard({ program }: { program: ProgramListItem }) {
  return (
    <Link
      to={`/programs/${program.id}`}
      className="
        group block
        bg-surface border border-border rounded-lg
        p-4 sm:p-5
        hover:border-accent/30 hover:bg-surface-elevated
        transition-all duration-200
      "
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <h3 className="text-base sm:text-lg font-semibold text-foreground group-hover:text-accent transition-colors">
            {program.name}
          </h3>
          {program.description && (
            <p className="mt-1 text-sm text-muted line-clamp-2">
              {program.description}
            </p>
          )}
          <div className="mt-3 flex items-center gap-2 text-xs text-muted">
            <span className="uppercase tracking-wider">
              {program.slug}
            </span>
          </div>
        </div>
        <ChevronRight
          size={20}
          className="flex-shrink-0 text-muted group-hover:text-accent group-hover:translate-x-0.5 transition-all"
        />
      </div>
    </Link>
  );
}

function LoadingSkeleton() {
  return (
    <div className="space-y-3">
      {[1, 2, 3].map((i) => (
        <div
          key={i}
          className="bg-surface border border-border rounded-lg p-4 sm:p-5 animate-pulse"
        >
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1">
              <div className="h-5 bg-surface-elevated rounded w-2/3 mb-2" />
              <div className="h-4 bg-surface-elevated rounded w-full mb-1" />
              <div className="h-4 bg-surface-elevated rounded w-4/5 mt-3" />
            </div>
            <div className="w-5 h-5 bg-surface-elevated rounded" />
          </div>
        </div>
      ))}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="text-center py-12">
      <div className="mx-auto w-16 h-16 rounded-full bg-surface-elevated border border-border flex items-center justify-center mb-4">
        <BookOpen className="w-8 h-8 text-muted" />
      </div>
      <h3 className="text-lg font-semibold text-foreground mb-2">
        No programs available
      </h3>
      <p className="text-sm text-muted max-w-sm mx-auto">
        There are no training programs available at the moment. Check back later.
      </p>
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
        Failed to load programs
      </h3>
      <p className="text-sm text-muted max-w-sm mx-auto">
        {message}
      </p>
    </div>
  );
}

export default function Programs() {
  const { data, isLoading, error } = usePrograms();

  // Handle both array response (simple list) and paginated response ({data: [], ...})
  const programs = Array.isArray(data) ? data : (data?.data ?? []);

  return (
    <div className="py-6 md:py-8">
      <Container>
        <div className="mb-6">
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground">
            Training Programs
          </h1>
          <p className="mt-2 text-muted">
            Browse available powerlifting programs
          </p>
        </div>

        {isLoading && <LoadingSkeleton />}

        {error && (
          <ErrorState
            message={error instanceof Error ? error.message : 'An unexpected error occurred'}
          />
        )}

        {!isLoading && !error && programs.length === 0 && <EmptyState />}

        {!isLoading && !error && programs.length > 0 && (
          <div className="space-y-3">
            {programs.map((program) => (
              <ProgramCard key={program.id} program={program} />
            ))}
          </div>
        )}
      </Container>
    </div>
  );
}
