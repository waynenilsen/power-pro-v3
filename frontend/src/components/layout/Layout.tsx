import { Outlet } from 'react-router-dom';
import { BottomNav } from './BottomNav';
import { SideNav } from './SideNav';

export function Layout() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      {/* Side navigation for desktop */}
      <SideNav />

      {/* Main content area */}
      <main className="
        min-h-screen
        pb-20
        md:pb-0
        md:pl-20
        lg:pl-64
        transition-all
        duration-200
      ">
        {/* Subtle background pattern for depth */}
        <div className="fixed inset-0 bg-grid-pattern opacity-[0.02] pointer-events-none -z-10" />

        {/* Page content */}
        <div className="relative">
          <Outlet />
        </div>
      </main>

      {/* Bottom navigation for mobile */}
      <BottomNav />
    </div>
  );
}
