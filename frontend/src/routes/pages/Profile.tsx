import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/useAuth';
import { LogOut, User, Copy, Check, AlertTriangle } from 'lucide-react';

export default function Profile() {
  const navigate = useNavigate();
  const { userId, logout } = useAuth();
  const [showConfirm, setShowConfirm] = useState(false);
  const [copied, setCopied] = useState(false);

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
