import { Link } from 'react-router-dom';
import { Container } from '../../components/layout';
import { BookOpen, ChevronRight } from 'lucide-react';

export default function Home() {
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
        </div>
      </Container>
    </div>
  );
}
