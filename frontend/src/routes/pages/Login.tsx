import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../contexts/useAuth';
import { Dumbbell, ChevronDown, Loader2, Zap, User } from 'lucide-react';

export default function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, createUser, isLoading } = useAuth();

  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || '/';
  const [showExistingId, setShowExistingId] = useState(false);
  const [existingId, setExistingId] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreateAccount = async () => {
    setError(null);
    setIsCreating(true);
    try {
      const userId = createUser();
      login(userId);
      navigate(from, { replace: true });
    } catch {
      setError('Failed to create account. Please try again.');
    } finally {
      setIsCreating(false);
    }
  };

  const handleExistingLogin = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (!existingId.trim()) {
      setError('Please enter a valid user ID');
      return;
    }
    login(existingId.trim());
    navigate(from, { replace: true });
  };

  const loading = isLoading || isCreating;

  return (
    <div className="min-h-screen bg-background relative overflow-hidden flex items-center justify-center p-4">
      {/* Dramatic background elements */}
      <div className="absolute inset-0 bg-grid-pattern opacity-[0.03]" />

      {/* Diagonal accent lines */}
      <div
        className="absolute top-0 right-0 w-[200%] h-1 bg-gradient-to-r from-transparent via-accent/30 to-transparent origin-top-right"
        style={{ transform: 'rotate(-35deg) translateY(-50vh)' }}
      />
      <div
        className="absolute bottom-0 left-0 w-[200%] h-px bg-gradient-to-r from-transparent via-accent/20 to-transparent origin-bottom-left"
        style={{ transform: 'rotate(-35deg) translateY(30vh)' }}
      />

      {/* Radial glow behind card */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-accent/5 rounded-full blur-[120px] pointer-events-none" />

      {/* Main card */}
      <div className="relative w-full max-w-md animate-fade-in">
        {/* Top accent bar */}
        <div className="h-1 bg-gradient-to-r from-accent via-accent-light to-accent rounded-t-lg" />

        <div className="bg-surface border border-border rounded-b-lg p-8 md:p-10">
          {/* Logo and branding */}
          <div className="text-center mb-10">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-accent/10 border border-accent/20 mb-6 relative group">
              <Dumbbell className="w-8 h-8 text-accent" strokeWidth={2.5} />
              <div className="absolute inset-0 rounded-2xl bg-accent/20 blur-xl opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
            </div>

            <h1 className="text-3xl font-bold tracking-tight text-foreground mb-2">
              POWER<span className="text-accent">PRO</span>
            </h1>
            <p className="text-muted text-sm tracking-wide uppercase">
              Strength Programming System
            </p>
          </div>

          {/* Welcome message */}
          <div className="mb-8 text-center">
            <h2 className="text-xl font-semibold text-foreground mb-2">
              Welcome, Athlete
            </h2>
            <p className="text-muted text-sm leading-relaxed">
              Track your lifts. Follow proven programs.<br />
              Build unstoppable strength.
            </p>
          </div>

          {/* Error display */}
          {error && (
            <div className="mb-6 p-3 bg-error/10 border border-error/20 rounded-lg text-error text-sm text-center animate-fade-in">
              {error}
            </div>
          )}

          {/* Primary CTA - Create Account */}
          <button
            onClick={handleCreateAccount}
            disabled={loading}
            className="
              w-full py-4 px-6 rounded-lg
              bg-accent hover:bg-accent-light
              text-background font-semibold text-lg
              transition-all duration-200
              flex items-center justify-center gap-3
              disabled:opacity-50 disabled:cursor-not-allowed
              hover:shadow-[0_0_30px_rgba(249,115,22,0.3)]
              active:scale-[0.98]
              group
            "
          >
            {loading ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                <span>Creating Account...</span>
              </>
            ) : (
              <>
                <Zap className="w-5 h-5 transition-transform group-hover:scale-110" />
                <span>Create New Account</span>
              </>
            )}
          </button>

          {/* Divider */}
          <div className="relative my-8">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-border" />
            </div>
            <div className="relative flex justify-center">
              <span className="px-4 bg-surface text-muted text-xs uppercase tracking-wider">
                or
              </span>
            </div>
          </div>

          {/* Collapsible Existing ID Section */}
          <div className="border border-border rounded-lg overflow-hidden">
            <button
              onClick={() => setShowExistingId(!showExistingId)}
              className="
                w-full py-3 px-4
                flex items-center justify-between
                text-muted hover:text-foreground
                transition-colors duration-200
                bg-surface-elevated/50 hover:bg-surface-elevated
              "
            >
              <span className="flex items-center gap-2 text-sm">
                <User className="w-4 h-4" />
                Use Existing ID
                <span className="text-xs text-muted/60">(dev)</span>
              </span>
              <ChevronDown
                className={`w-4 h-4 transition-transform duration-200 ${showExistingId ? 'rotate-180' : ''}`}
              />
            </button>

            {showExistingId && (
              <form onSubmit={handleExistingLogin} className="p-4 bg-surface-elevated/30 animate-fade-in">
                <div className="mb-4">
                  <label htmlFor="userId" className="block text-xs text-muted uppercase tracking-wider mb-2">
                    User ID
                  </label>
                  <input
                    type="text"
                    id="userId"
                    value={existingId}
                    onChange={(e) => setExistingId(e.target.value)}
                    placeholder="Enter your existing user ID"
                    className="
                      w-full py-3 px-4 rounded-lg
                      bg-background border border-border
                      text-foreground placeholder:text-muted/50
                      focus:border-accent focus:outline-none
                      transition-colors duration-200
                      text-sm font-mono
                    "
                    disabled={loading}
                  />
                </div>
                <button
                  type="submit"
                  disabled={loading || !existingId.trim()}
                  className="
                    w-full py-3 px-4 rounded-lg
                    bg-surface-elevated border border-border
                    text-foreground font-medium
                    hover:border-accent/50 hover:bg-surface-elevated/80
                    transition-all duration-200
                    disabled:opacity-50 disabled:cursor-not-allowed
                    text-sm
                  "
                >
                  Continue with ID
                </button>
              </form>
            )}
          </div>

          {/* Footer note */}
          <p className="mt-8 text-center text-xs text-muted/60">
            Your data is stored locally on your device.
          </p>
        </div>

        {/* Bottom decorative element */}
        <div className="mt-6 flex justify-center gap-1">
          {[...Array(5)].map((_, i) => (
            <div
              key={i}
              className="w-1.5 h-1.5 rounded-full bg-accent/30"
              style={{
                opacity: i === 2 ? 1 : 0.3 + (i === 1 || i === 3 ? 0.2 : 0),
              }}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
