import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../../contexts/useAuth';
import { useEnrollment, useLiftMaxes, useLifts } from '../../hooks';
import { useProgram } from '../../hooks/useProgram';
import {
  LogOut,
  User,
  Copy,
  Check,
  AlertTriangle,
  Dumbbell,
  ChevronRight,
  TrendingUp,
  AlertCircle,
  Loader2,
} from 'lucide-react';
import type { LiftMax, Lift, MaxType } from '../../api/types';

function MaxTypeBadge({ type }: { type: MaxType }) {
  const isOneRM = type === 'ONE_RM';
  return (
    <span
      className={`
        inline-flex items-center px-1.5 py-0.5
        text-[10px] font-bold uppercase tracking-wider
        rounded
        ${isOneRM
          ? 'bg-accent/15 text-accent'
          : 'bg-success/15 text-success'
        }
      `}
    >
      {isOneRM ? '1RM' : 'TM'}
    </span>
  );
}

interface LiftMaxSummary {
  liftId: string;
  liftName: string;
  isCompetitionLift: boolean;
  oneRM?: LiftMax;
  trainingMax?: LiftMax;
}

function getLiftMaxSummaries(maxes: LiftMax[], lifts: Lift[]): LiftMaxSummary[] {
  const summaryMap = new Map<string, LiftMaxSummary>();

  // Get competition lifts first, then others
  for (const lift of lifts) {
    summaryMap.set(lift.id, {
      liftId: lift.id,
      liftName: lift.name,
      isCompetitionLift: lift.isCompetitionLift,
    });
  }

  // Populate with most recent maxes
  for (const max of maxes) {
    const summary = summaryMap.get(max.liftId);
    if (!summary) continue;

    if (max.type === 'ONE_RM') {
      if (!summary.oneRM || new Date(max.effectiveDate) > new Date(summary.oneRM.effectiveDate)) {
        summary.oneRM = max;
      }
    } else {
      if (!summary.trainingMax || new Date(max.effectiveDate) > new Date(summary.trainingMax.effectiveDate)) {
        summary.trainingMax = max;
      }
    }
  }

  // Return only lifts that have at least one max, sorted by competition lifts first
  return Array.from(summaryMap.values())
    .filter((s) => s.oneRM || s.trainingMax)
    .sort((a, b) => {
      if (a.isCompetitionLift !== b.isCompetitionLift) {
        return a.isCompetitionLift ? -1 : 1;
      }
      return a.liftName.localeCompare(b.liftName);
    });
}

function EnrollmentCard({ userId }: { userId: string }) {
  const { data: enrollment, isLoading, error } = useEnrollment(userId);

  if (isLoading) {
    return (
      <div className="bg-surface border border-border rounded-lg p-6 mb-6">
        <div className="flex items-center justify-center py-4">
          <Loader2 className="w-5 h-5 text-accent animate-spin" />
        </div>
      </div>
    );
  }

  if (error || !enrollment) {
    return (
      <Link
        to="/programs"
        className="
          block bg-surface border border-border rounded-lg p-6 mb-6
          hover:border-accent/30 hover:bg-surface-elevated
          transition-all duration-200
        "
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 rounded-lg bg-muted/10 border border-border flex items-center justify-center">
              <Dumbbell className="w-6 h-6 text-muted" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-foreground">No Program</h2>
              <p className="text-sm text-muted">Browse available programs to get started</p>
            </div>
          </div>
          <ChevronRight className="w-5 h-5 text-muted" />
        </div>
      </Link>
    );
  }

  const { program, state, enrollmentStatus } = enrollment;
  const progressPercentage = (state.currentWeek / program.cycleLengthWeeks) * 100;

  return (
    <Link
      to={`/programs/${program.slug}`}
      className="
        block bg-surface border border-border rounded-lg p-6 mb-6
        hover:border-accent/30 hover:bg-surface-elevated
        transition-all duration-200
      "
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3">
          <div className="w-12 h-12 rounded-lg bg-accent/10 border border-accent/20 flex items-center justify-center flex-shrink-0">
            <Dumbbell className="w-6 h-6 text-accent" />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-foreground">{program.name}</h2>
            <p className="text-sm text-muted">
              {enrollmentStatus === 'BETWEEN_CYCLES'
                ? 'Cycle complete - ready for next'
                : `Week ${state.currentWeek} of ${program.cycleLengthWeeks} · Cycle ${state.currentCycleIteration}`}
            </p>
          </div>
        </div>
        <ChevronRight className="w-5 h-5 text-muted flex-shrink-0 mt-1" />
      </div>

      {/* Progress bar */}
      {enrollmentStatus !== 'BETWEEN_CYCLES' && (
        <div className="mt-4">
          <div className="h-2 bg-surface-elevated rounded-full overflow-hidden">
            <div
              className="h-full bg-accent rounded-full transition-all duration-500"
              style={{ width: `${progressPercentage}%` }}
            />
          </div>
        </div>
      )}
    </Link>
  );
}

function LiftMaxesSection({ userId, enrolledProgramId }: { userId: string; enrolledProgramId?: string }) {
  const { data: maxesData, isLoading: maxesLoading, error: maxesError } = useLiftMaxes(userId);
  const { data: liftsData, isLoading: liftsLoading } = useLifts();
  const { data: programDetail } = useProgram(enrolledProgramId);

  const isLoading = maxesLoading || liftsLoading;
  const maxes = maxesData?.data ?? [];
  const lifts = liftsData?.data ?? [];

  const summaries = getLiftMaxSummaries(maxes, lifts);

  // Get required lifts from the enrolled program, or fall back to competition lifts
  const requiredLiftNames = programDetail?.liftRequirements ??
    lifts.filter((l) => l.isCompetitionLift).map((l) => l.name);

  // Find lifts that have training maxes
  const liftsWithTrainingMax = summaries.filter((s) => s.trainingMax).map((s) => s.liftName);

  // Find required lifts that are missing training maxes
  const missingRequiredLifts = requiredLiftNames.filter(
    (name) => !liftsWithTrainingMax.includes(name)
  );

  return (
    <div className="bg-surface border border-border rounded-lg p-6 mb-6">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-lg bg-accent/10 border border-accent/20 flex items-center justify-center">
            <TrendingUp className="w-6 h-6 text-accent" />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-foreground">Lift Maxes</h2>
            <p className="text-sm text-muted">Your recorded maxes</p>
          </div>
        </div>
        <Link
          to="/lift-maxes"
          className="
            text-sm font-medium text-accent
            hover:text-accent-light
            transition-colors
          "
        >
          View All →
        </Link>
      </div>

      {/* Warning for missing required lift maxes */}
      {missingRequiredLifts.length > 0 && (
        <div className="mb-4 p-3 bg-warning/10 border border-warning/20 rounded-lg flex items-start gap-2">
          <AlertCircle className="w-4 h-4 text-warning flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-sm font-medium text-foreground">Missing Training Maxes</p>
            <p className="text-xs text-muted mt-0.5">
              Add maxes for: {missingRequiredLifts.join(', ')}
            </p>
          </div>
        </div>
      )}

      {isLoading ? (
        <div className="flex items-center justify-center py-6">
          <Loader2 className="w-5 h-5 text-accent animate-spin" />
        </div>
      ) : maxesError ? (
        <div className="text-center py-4">
          <p className="text-sm text-error">Failed to load lift maxes</p>
        </div>
      ) : summaries.length === 0 ? (
        <div className="text-center py-6">
          <p className="text-sm text-muted mb-3">No lift maxes recorded yet</p>
          <Link
            to="/lift-maxes/new"
            className="
              inline-flex items-center gap-1 text-sm font-medium
              text-accent hover:text-accent-light
              transition-colors
            "
          >
            Add your first max →
          </Link>
        </div>
      ) : (
        <div className="space-y-2">
          {summaries.slice(0, 5).map((summary) => (
            <Link
              key={summary.liftId}
              to={`/lift-maxes/${summary.liftId}/history`}
              className="
                flex items-center justify-between p-3
                bg-surface-elevated rounded-lg border border-transparent
                hover:border-accent/20
                transition-all duration-200
              "
            >
              <div className="flex items-center gap-2">
                <span className="font-medium text-foreground text-sm">
                  {summary.liftName}
                </span>
                {summary.isCompetitionLift && (
                  <span className="text-[10px] uppercase tracking-wider text-accent font-bold">
                    Main
                  </span>
                )}
              </div>
              <div className="flex items-center gap-3">
                {summary.trainingMax && (
                  <div className="flex items-center gap-1.5">
                    <MaxTypeBadge type="TRAINING_MAX" />
                    <span className="text-sm font-bold tabular-nums text-foreground">
                      {summary.trainingMax.value}
                    </span>
                  </div>
                )}
                {summary.oneRM && (
                  <div className="flex items-center gap-1.5">
                    <MaxTypeBadge type="ONE_RM" />
                    <span className="text-sm font-bold tabular-nums text-foreground">
                      {summary.oneRM.value}
                    </span>
                  </div>
                )}
              </div>
            </Link>
          ))}

          {summaries.length > 5 && (
            <Link
              to="/lift-maxes"
              className="
                block text-center py-2
                text-sm text-muted hover:text-accent
                transition-colors
              "
            >
              +{summaries.length - 5} more
            </Link>
          )}
        </div>
      )}
    </div>
  );
}

export default function Profile() {
  const navigate = useNavigate();
  const { userId, logout } = useAuth();
  const [showConfirm, setShowConfirm] = useState(false);
  const [copied, setCopied] = useState(false);

  // Fetch enrollment to get the enrolled program ID
  const { data: enrollment } = useEnrollment(userId ?? undefined);
  const enrolledProgramId = enrollment?.program?.id;

  const handleLogout = () => {
    logout();
    navigate('/login', { replace: true });
  };

  const handleCopyId = async () => {
    if (userId) {
      await navigator.clipboard.writeText(userId);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="p-4 md:p-6 max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-foreground mb-6">Profile</h1>

      {/* Program Enrollment Section */}
      {userId && <EnrollmentCard userId={userId} />}

      {/* Lift Maxes Section */}
      {userId && <LiftMaxesSection userId={userId} enrolledProgramId={enrolledProgramId} />}

      {/* User ID Section */}
      <div className="bg-surface border border-border rounded-lg p-6 mb-6">
        <div className="flex items-center gap-3 mb-4">
          <div className="w-12 h-12 rounded-full bg-accent/10 border border-accent/20 flex items-center justify-center">
            <User className="w-6 h-6 text-accent" />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-foreground">Your Account</h2>
            <p className="text-sm text-muted">User ID for this device</p>
          </div>
        </div>

        <div className="bg-surface-elevated rounded-lg p-4 border border-border">
          <label className="block text-xs text-muted uppercase tracking-wider mb-2">
            User ID
          </label>
          <div className="flex items-center gap-2">
            <code className="flex-1 text-sm font-mono text-foreground break-all">
              {userId}
            </code>
            <button
              onClick={handleCopyId}
              className="
                p-2 rounded-lg
                text-muted hover:text-foreground
                hover:bg-white/5
                transition-colors duration-200
              "
              title="Copy User ID"
            >
              {copied ? (
                <Check className="w-4 h-4 text-accent" />
              ) : (
                <Copy className="w-4 h-4" />
              )}
            </button>
          </div>
          <p className="mt-3 text-xs text-muted">
            Save this ID to access your data on another device.
          </p>
        </div>
      </div>

      {/* Sign Out Section */}
      <div className="bg-surface border border-border rounded-lg p-6">
        <h2 className="text-lg font-semibold text-foreground mb-2">Sign Out</h2>
        <p className="text-sm text-muted mb-4">
          Sign out of your account on this device. Make sure to save your User ID first if you want to access your data again.
        </p>

        {!showConfirm ? (
          <button
            onClick={() => setShowConfirm(true)}
            className="
              flex items-center gap-2 py-3 px-4 rounded-lg
              bg-surface-elevated border border-border
              text-foreground font-medium
              hover:border-error/50 hover:text-error
              transition-all duration-200
            "
          >
            <LogOut className="w-4 h-4" />
            Sign Out
          </button>
        ) : (
          <div className="bg-error/10 border border-error/20 rounded-lg p-4 animate-fade-in">
            <div className="flex items-start gap-3 mb-4">
              <AlertTriangle className="w-5 h-5 text-error flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-sm font-medium text-foreground">
                  Are you sure you want to sign out?
                </p>
                <p className="text-xs text-muted mt-1">
                  You'll need your User ID to access this account again.
                </p>
              </div>
            </div>
            <div className="flex gap-3">
              <button
                onClick={handleLogout}
                className="
                  flex-1 py-2.5 px-4 rounded-lg
                  bg-error hover:bg-error/90
                  text-white font-medium text-sm
                  transition-colors duration-200
                "
              >
                Yes, Sign Out
              </button>
              <button
                onClick={() => setShowConfirm(false)}
                className="
                  flex-1 py-2.5 px-4 rounded-lg
                  bg-surface-elevated border border-border
                  text-foreground font-medium text-sm
                  hover:bg-surface-elevated/80
                  transition-colors duration-200
                "
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
