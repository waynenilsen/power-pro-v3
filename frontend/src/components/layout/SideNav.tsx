import { NavLink } from 'react-router-dom';
import { Home, Dumbbell, History, User, Zap } from 'lucide-react';

const navItems = [
  { to: '/', icon: Home, label: 'Home' },
  { to: '/workout', icon: Dumbbell, label: 'Workout' },
  { to: '/history', icon: History, label: 'History' },
  { to: '/profile', icon: User, label: 'Profile' },
];

export function SideNav() {
  return (
    <nav className="hidden md:flex fixed left-0 top-0 bottom-0 z-50 w-20 lg:w-64 flex-col bg-surface-elevated border-r border-white/5">
      {/* Logo Section */}
      <div className="flex items-center h-16 px-4 lg:px-6 border-b border-white/5">
        <div className="flex items-center gap-3">
          <div className="relative">
            <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-accent to-accent-dark flex items-center justify-center">
              <Zap size={22} className="text-background" strokeWidth={2.5} />
            </div>
            {/* Glow effect */}
            <div className="absolute inset-0 rounded-lg bg-accent/20 blur-md -z-10" />
          </div>
          <span className="hidden lg:block text-lg font-bold tracking-tight text-foreground">
            Power<span className="text-accent">Pro</span>
          </span>
        </div>
      </div>

      {/* Navigation Items */}
      <div className="flex-1 py-6 px-3 lg:px-4">
        <ul className="space-y-1">
          {navItems.map(({ to, icon: Icon, label }) => (
            <li key={to}>
              <NavLink
                to={to}
                className={({ isActive }) =>
                  `group relative flex items-center gap-4 h-12 px-3 lg:px-4 rounded-lg transition-all duration-200 ${
                    isActive
                      ? 'bg-accent/10 text-accent'
                      : 'text-muted hover:bg-white/5 hover:text-foreground'
                  }`
                }
              >
                {({ isActive }) => (
                  <>
                    {/* Active indicator */}
                    {isActive && (
                      <span className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-6 bg-accent rounded-r-full" />
                    )}

                    {/* Icon */}
                    <span className={`flex-shrink-0 transition-transform duration-200 ${isActive ? '' : 'group-hover:scale-105'}`}>
                      <Icon
                        size={22}
                        strokeWidth={isActive ? 2.5 : 2}
                      />
                    </span>

                    {/* Label - hidden on collapsed state */}
                    <span className={`hidden lg:block text-sm font-medium tracking-wide transition-all duration-200 ${
                      isActive ? '' : 'opacity-80 group-hover:opacity-100'
                    }`}>
                      {label}
                    </span>
                  </>
                )}
              </NavLink>
            </li>
          ))}
        </ul>
      </div>

      {/* Bottom section - version/branding */}
      <div className="hidden lg:block px-6 py-4 border-t border-white/5">
        <p className="text-[10px] text-muted uppercase tracking-widest">
          Built for Strength
        </p>
      </div>
    </nav>
  );
}
