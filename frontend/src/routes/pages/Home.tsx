import { Container } from '../../components/layout';

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
      </Container>
    </div>
  );
}
