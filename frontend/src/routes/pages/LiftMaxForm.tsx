import { useState } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { Container } from '../../components/layout';
import { useAuth } from '../../contexts/useAuth';
import { useLifts, useLiftMaxes, useCreateLiftMax, useUpdateLiftMax } from '../../hooks';
import {
  Loader2,
  ChevronLeft,
  Save,
  Calculator,
  AlertCircle,
} from 'lucide-react';
import type { MaxType, LiftMax, Lift } from '../../api/types';

function formatDateForInput(dateString?: string): string {
  const date = dateString ? new Date(dateString) : new Date();
  return date.toISOString().split('T')[0];
}

interface FormState {
  liftId: string;
  type: MaxType;
  value: string;
  effectiveDate: string;
}

function TrainingMaxCalculator({
  value,
  onCalculate,
}: {
  value: string;
  onCalculate: (tmValue: number) => void;
}) {
  const oneRM = parseFloat(value);
  if (!oneRM || oneRM <= 0) return null;

  const trainingMax = Math.round(oneRM * 0.9);

  return (
    <div className="mt-4 p-4 bg-accent/5 border border-accent/20 rounded-lg">
      <div className="flex items-start gap-3">
        <Calculator className="w-5 h-5 text-accent mt-0.5" />
        <div className="flex-1">
          <p className="text-sm font-medium text-foreground">
            Training Max Calculator
          </p>
          <p className="text-xs text-muted mt-1">
            90% of {oneRM} lbs = <span className="font-bold text-accent">{trainingMax} lbs</span>
          </p>
          <button
            type="button"
            onClick={() => onCalculate(trainingMax)}
            className="
              mt-2 text-xs font-semibold text-accent
              hover:text-accent-light
              transition-colors
            "
          >
            Use {trainingMax} lbs as Training Max →
          </button>
        </div>
      </div>
    </div>
  );
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
        type: existingMax.type,
        value: existingMax.value.toString(),
        effectiveDate: formatDateForInput(existingMax.effectiveDate),
      };
    }
    return {
      liftId: '',
      type: 'ONE_RM',
      value: '',
      effectiveDate: formatDateForInput(),
    };
  });
  const [error, setError] = useState<string | null>(null);
  const [showCalculator, setShowCalculator] = useState(false);

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

  const handleCalculatorUse = (tmValue: number) => {
    setForm((prev) => ({
      ...prev,
      type: 'TRAINING_MAX',
      value: tmValue.toString(),
    }));
    setShowCalculator(false);
  };

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

      {/* Max Type */}
      <div className="mb-6">
        <label className="block text-sm font-medium text-foreground mb-2">
          Max Type
        </label>
        <div className="flex gap-3">
          <button
            type="button"
            onClick={() => {
              setForm((prev) => ({ ...prev, type: 'ONE_RM' }));
              setShowCalculator(true);
            }}
            disabled={isEditing}
            className={`
              flex-1 py-3 px-4 rounded-lg
              font-medium text-sm
              border transition-all duration-200
              disabled:cursor-not-allowed
              ${form.type === 'ONE_RM'
                ? 'bg-accent text-background border-accent'
                : 'bg-surface border-border text-muted hover:text-foreground hover:border-accent/30'
              }
            `}
          >
            1RM (One Rep Max)
          </button>
          <button
            type="button"
            onClick={() => {
              setForm((prev) => ({ ...prev, type: 'TRAINING_MAX' }));
              setShowCalculator(false);
            }}
            disabled={isEditing}
            className={`
              flex-1 py-3 px-4 rounded-lg
              font-medium text-sm
              border transition-all duration-200
              disabled:cursor-not-allowed
              ${form.type === 'TRAINING_MAX'
                ? 'bg-success text-background border-success'
                : 'bg-surface border-border text-muted hover:text-foreground hover:border-success/30'
              }
            `}
          >
            Training Max
          </button>
        </div>
        <p className="mt-2 text-xs text-muted">
          {form.type === 'ONE_RM'
            ? 'Your actual one-rep max (tested or estimated).'
            : 'Your working max, typically 85-90% of your 1RM.'}
        </p>
      </div>

      {/* Weight Value */}
      <div className="mb-6">
        <label htmlFor="value" className="block text-sm font-medium text-foreground mb-2">
          Weight (lbs)
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

        {/* Training Max Calculator (only for 1RM) */}
        {form.type === 'ONE_RM' && showCalculator && (
          <TrainingMaxCalculator value={form.value} onCalculate={handleCalculatorUse} />
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
  const existingMax = isEditing && maxesData?.data
    ? maxesData.data.find((m: LiftMax) => m.id === id) ?? null
    : null;

  const handleSubmit = async (form: FormState) => {
    const valueNum = parseFloat(form.value);

    if (isEditing && id) {
      await updateMax.mutateAsync({
        liftMaxId: id,
        request: {
          value: valueNum,
          effectiveDate: form.effectiveDate,
        },
      });
    } else {
      await createMax.mutateAsync({
        liftId: form.liftId,
        type: form.type,
        value: valueNum,
        effectiveDate: form.effectiveDate,
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
            {isEditing ? 'Edit Lift Max' : 'Add Lift Max'}
          </h1>
          <p className="mt-2 text-muted">
            {isEditing
              ? 'Update your recorded max for this lift.'
              : 'Record a new personal record or training max.'}
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
