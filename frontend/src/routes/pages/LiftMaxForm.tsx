import { useState } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { Container } from '../../components/layout';
import { useAuth } from '../../contexts/useAuth';
import { useLifts, useLiftMaxes, useCreateLiftMax, useUpdateLiftMax } from '../../hooks';
import {
  Loader2,
  ChevronLeft,
  Save,
  AlertCircle,
  Info,
} from 'lucide-react';
import type { LiftMax, Lift } from '../../api/types';

function formatDateForInput(dateString?: string): string {
  const date = dateString ? new Date(dateString) : new Date();
  return date.toISOString().split('T')[0];
}

function formatDateForApi(dateString: string): string {
  // Convert YYYY-MM-DD to RFC3339 format for Go backend
  return new Date(dateString + 'T00:00:00Z').toISOString();
}

interface FormState {
  liftId: string;
  value: string;
  effectiveDate: string;
}

interface LiftMaxFormInnerProps {
  isEditing: boolean;
  existingMax: LiftMax | null;
  lifts: Lift[];
  onSubmit: (form: FormState) => Promise<void>;
  isSaving: boolean;
}

function LiftMaxFormInner({
  isEditing,
  existingMax,
  lifts,
  onSubmit,
  isSaving,
}: LiftMaxFormInnerProps) {
  const [form, setForm] = useState<FormState>(() => {
    if (existingMax) {
      return {
        liftId: existingMax.liftId,
        value: existingMax.value.toString(),
        effectiveDate: formatDateForInput(existingMax.effectiveDate),
      };
    }
    return {
      liftId: '',
      value: '',
      effectiveDate: formatDateForInput(),
    };
  });
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const valueNum = parseFloat(form.value);
    if (!form.liftId) {
      setError('Please select a lift');
      return;
    }
    if (!valueNum || valueNum <= 0) {
      setError('Please enter a valid weight');
      return;
    }

    try {
      await onSubmit(form);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save lift max');
    }
  };

  // Calculate what the TM will be
  const valueNum = parseFloat(form.value);
  const trainingMax = valueNum > 0 ? Math.round(valueNum * 0.9 * 4) / 4 : null;

  return (
    <form onSubmit={handleSubmit} className="max-w-lg">
      {/* Error message */}
      {error && (
        <div className="mb-6 p-4 bg-error/10 border border-error/20 rounded-lg flex items-start gap-3">
          <AlertCircle className="w-5 h-5 text-error flex-shrink-0 mt-0.5" />
          <p className="text-sm text-error">{error}</p>
        </div>
      )}

      {/* Lift Selection */}
      <div className="mb-6">
        <label htmlFor="liftId" className="block text-sm font-medium text-foreground mb-2">
          Lift
        </label>
        <select
          id="liftId"
          value={form.liftId}
          onChange={(e) => setForm((prev) => ({ ...prev, liftId: e.target.value }))}
          disabled={isEditing}
          className="
            w-full px-4 py-3
            bg-surface border border-border rounded-lg
            text-foreground
            focus:outline-none focus:border-accent focus:ring-1 focus:ring-accent
            disabled:opacity-60 disabled:cursor-not-allowed
          "
        >
          <option value="">Select a lift...</option>
          {lifts.map((lift) => (
            <option key={lift.id} value={lift.id}>
              {lift.name}
            </option>
          ))}
        </select>
      </div>

      {/* Weight Value */}
      <div className="mb-6">
        <label htmlFor="value" className="block text-sm font-medium text-foreground mb-2">
          1RM (One Rep Max) in lbs
        </label>
        <input
          type="number"
          id="value"
          value={form.value}
          onChange={(e) => setForm((prev) => ({ ...prev, value: e.target.value }))}
          min="0"
          step="2.5"
          className="
            w-full px-4 py-3
            bg-surface border border-border rounded-lg
            text-foreground text-lg font-bold tabular-nums
            focus:outline-none focus:border-accent focus:ring-1 focus:ring-accent
            [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none
          "
          placeholder="0"
        />
        <p className="mt-2 text-xs text-muted">
          Enter your actual or estimated one-rep max.
        </p>

        {/* Training Max Info */}
        {trainingMax && (
          <div className="mt-4 p-4 bg-success/5 border border-success/20 rounded-lg">
            <div className="flex items-start gap-3">
              <Info className="w-5 h-5 text-success mt-0.5" />
              <div className="flex-1">
                <p className="text-sm font-medium text-foreground">
                  Training Max: <span className="text-success">{trainingMax} lbs</span>
                </p>
                <p className="text-xs text-muted mt-1">
                  Automatically calculated at 90% of your 1RM
                </p>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Effective Date */}
      <div className="mb-8">
        <label htmlFor="effectiveDate" className="block text-sm font-medium text-foreground mb-2">
          Date
        </label>
        <input
          type="date"
          id="effectiveDate"
          value={form.effectiveDate}
          onChange={(e) => setForm((prev) => ({ ...prev, effectiveDate: e.target.value }))}
          className="
            w-full px-4 py-3
            bg-surface border border-border rounded-lg
            text-foreground
            focus:outline-none focus:border-accent focus:ring-1 focus:ring-accent
          "
        />
        <p className="mt-2 text-xs text-muted">
          When was this max achieved or recorded?
        </p>
      </div>

      {/* Submit Button */}
      <div className="flex gap-3">
        <button
          type="submit"
          disabled={isSaving}
          className="
            flex items-center justify-center gap-2
            flex-1 py-3.5 px-6
            bg-accent rounded-lg
            text-background font-bold
            hover:bg-accent-light
            disabled:opacity-50 disabled:cursor-not-allowed
            transition-colors
          "
        >
          {isSaving ? (
            <Loader2 size={18} className="animate-spin" />
          ) : (
            <Save size={18} />
          )}
          {isSaving ? 'Saving...' : 'Save'}
        </button>
        <Link
          to="/lift-maxes"
          className="
            py-3.5 px-6
            bg-surface border border-border rounded-lg
            text-foreground font-medium
            hover:bg-surface-elevated
            transition-colors
          "
        >
          Cancel
        </Link>
      </div>
    </form>
  );
}

export default function LiftMaxForm() {
  const navigate = useNavigate();
  const { id } = useParams<{ id?: string }>();
  const { userId } = useAuth();

  const isEditing = Boolean(id) && id !== 'new';

  const { data: liftsData, isLoading: liftsLoading } = useLifts();
  const { data: maxesData, isLoading: maxesLoading } = useLiftMaxes(userId ?? undefined);

  const createMax = useCreateLiftMax(userId ?? undefined);
  const updateMax = useUpdateLiftMax(userId ?? undefined);

  const lifts = liftsData?.data ?? [];
  // Only allow editing 1RM entries
  const existingMax = isEditing && maxesData?.data
    ? maxesData.data.find((m: LiftMax) => m.id === id && m.type === 'ONE_RM') ?? null
    : null;

  const handleSubmit = async (form: FormState) => {
    const valueNum = parseFloat(form.value);
    const effectiveDateRfc3339 = formatDateForApi(form.effectiveDate);

    if (isEditing && id) {
      await updateMax.mutateAsync({
        liftMaxId: id,
        request: {
          value: valueNum,
          effectiveDate: effectiveDateRfc3339,
        },
      });
    } else {
      await createMax.mutateAsync({
        liftId: form.liftId,
        value: valueNum,
        effectiveDate: effectiveDateRfc3339,
      });
    }
    navigate('/lift-maxes');
  };

  const isLoading = liftsLoading || (isEditing && maxesLoading);
  const isSaving = createMax.isPending || updateMax.isPending;

  if (isLoading) {
    return (
      <div className="py-6 md:py-8">
        <Container>
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-8 h-8 text-accent animate-spin" />
          </div>
        </Container>
      </div>
    );
  }

  // If editing a TM directly, redirect to list (shouldn't happen but safety check)
  if (isEditing && maxesData?.data) {
    const targetMax = maxesData.data.find((m: LiftMax) => m.id === id);
    if (targetMax && targetMax.type === 'TRAINING_MAX') {
      navigate('/lift-maxes');
      return null;
    }
  }

  return (
    <div className="py-6 md:py-8">
      <Container>
        {/* Header */}
        <div className="mb-8">
          <Link
            to="/lift-maxes"
            className="inline-flex items-center gap-1 text-sm text-muted hover:text-accent transition-colors mb-4"
          >
            <ChevronLeft size={16} />
            Back to Lift Maxes
          </Link>

          <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground">
            {isEditing ? 'Edit 1RM' : 'Record 1RM'}
          </h1>
          <p className="mt-2 text-muted">
            {isEditing
              ? 'Update your one-rep max. Training max will be recalculated automatically.'
              : 'Record your one-rep max. A training max at 90% will be set automatically.'}
          </p>
        </div>

        {/* Form - key prop resets state when switching between new/edit */}
        <LiftMaxFormInner
          key={existingMax?.id ?? 'new'}
          isEditing={isEditing}
          existingMax={existingMax}
          lifts={lifts}
          onSubmit={handleSubmit}
          isSaving={isSaving}
        />
      </Container>
    </div>
  );
}
