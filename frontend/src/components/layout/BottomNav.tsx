import { NavLink } from 'react-router-dom';
import { Home, Dumbbell, History, User } from 'lucide-react';

const navItems = [
  { to: '/', icon: Home, label: 'Home' },
  { to: '/workout', icon: Dumbbell, label: 'Workout' },
  { to: '/history', icon: History, label: 'History' },
  { to: '/profile', icon: User, label: 'Profile' },
];

export function BottomNav() {
  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 md:hidden">
      {/* Industrial edge line */}
      <div className="h-[2px] bg-gradient-to-r from-transparent via-accent to-transparent" />

      {/* Main nav bar with subtle noise texture overlay */}
      <div className="bg-surface-elevated/95 backdrop-blur-lg border-t border-white/5">
        <div className="flex items-center justify-around h-16 max-w-lg mx-auto px-2">
          {navItems.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) =>
                `group relative flex flex-col items-center justify-center w-16 h-14 rounded-lg transition-all duration-200 ${
                  isActive
                    ? 'text-accent'
                    : 'text-muted hover:text-foreground'
                }`
              }
            >
              {({ isActive }) => (
                <>
                  {/* Active indicator bar */}
                  {isActive && (
                    <span className="absolute -top-[2px] left-1/2 -translate-x-1/2 w-8 h-[2px] bg-accent rounded-full" />
                  )}

                  {/* Icon with scale effect */}
                  <span className={`transition-transform duration-200 ${isActive ? 'scale-110' : 'group-hover:scale-105'}`}>
                    <Icon
                      size={22}
                      strokeWidth={isActive ? 2.5 : 2}
                      className="transition-all duration-200"
                    />
                  </span>

                  {/* Label */}
                  <span className={`mt-1 text-[10px] font-medium tracking-wide uppercase transition-all duration-200 ${
                    isActive ? 'opacity-100' : 'opacity-70 group-hover:opacity-100'
                  }`}>
                    {label}
                  </span>
                </>
              )}
            </NavLink>
          ))}
        </div>
      </div>

      {/* Safe area padding for notched devices */}
      <div className="h-safe-area-inset-bottom bg-surface-elevated" />
    </nav>
  );
}
