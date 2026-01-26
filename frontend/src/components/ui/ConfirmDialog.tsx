import { useEffect, useRef } from 'react';
import { AlertTriangle, X } from 'lucide-react';

export interface ConfirmDialogProps {
  isOpen: boolean;
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: 'default' | 'destructive';
  isLoading?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export function ConfirmDialog({
  isOpen,
  title,
  message,
  confirmLabel = 'Confirm',
  cancelLabel = 'Cancel',
  variant = 'default',
  isLoading = false,
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  const dialogRef = useRef<HTMLDivElement>(null);

  // Handle escape key
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen && !isLoading) {
        onCancel();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, isLoading, onCancel]);

  // Lock body scroll when open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  // Focus trap
  useEffect(() => {
    if (isOpen && dialogRef.current) {
      dialogRef.current.focus();
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const isDestructive = variant === 'destructive';

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm animate-fade-in"
        onClick={isLoading ? undefined : onCancel}
      />

      {/* Dialog */}
      <div
        ref={dialogRef}
        tabIndex={-1}
        className="
          relative z-10 w-full max-w-md mx-4
          bg-surface border border-border rounded-xl
          shadow-2xl
          animate-slide-up
        "
        role="dialog"
        aria-modal="true"
        aria-labelledby="dialog-title"
      >
        {/* Header */}
        <div className="flex items-start gap-4 p-5 border-b border-border">
          <div
            className={`
              w-10 h-10 rounded-full flex-shrink-0
              flex items-center justify-center
              ${isDestructive ? 'bg-error/10' : 'bg-accent/10'}
            `}
          >
            <AlertTriangle
              className={`w-5 h-5 ${isDestructive ? 'text-error' : 'text-accent'}`}
            />
          </div>
          <div className="flex-1">
            <h2
              id="dialog-title"
              className="text-lg font-semibold text-foreground"
            >
              {title}
            </h2>
            <p className="text-sm text-muted mt-1">{message}</p>
          </div>
          <button
            onClick={onCancel}
            disabled={isLoading}
            className="
              p-1 rounded-lg
              text-muted hover:text-foreground
              hover:bg-surface-elevated
              transition-colors duration-200
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            <X size={20} />
          </button>
        </div>

        {/* Actions */}
        <div className="flex gap-3 p-5">
          <button
            onClick={onCancel}
            disabled={isLoading}
            className="
              flex-1 py-3 px-4 rounded-lg
              bg-surface-elevated border border-border
              text-foreground font-medium
              hover:bg-surface-elevated/80
              transition-colors duration-200
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            {cancelLabel}
          </button>
          <button
            onClick={onConfirm}
            disabled={isLoading}
            className={`
              flex-1 py-3 px-4 rounded-lg
              font-medium
              transition-colors duration-200
              disabled:opacity-70 disabled:cursor-not-allowed
              ${
                isDestructive
                  ? 'bg-error hover:bg-error/90 text-white'
                  : 'bg-accent hover:bg-accent-light text-background'
              }
            `}
          >
            {isLoading ? (
              <span className="flex items-center justify-center gap-2">
                <span className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
                Processing...
              </span>
            ) : (
              confirmLabel
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
