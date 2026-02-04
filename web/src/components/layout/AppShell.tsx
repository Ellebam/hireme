'use client';

import { TooltipProvider } from '@/components/ui/tooltip';
import { Header } from './Header';

interface AppShellProps {
  children: React.ReactNode;
  /** Whether to show the header (default: true) */
  showHeader?: boolean;
  /** Full height layout for editor */
  fullHeight?: boolean;
}

export function AppShell({
  children,
  showHeader = true,
  fullHeight = false,
}: AppShellProps) {
  return (
    <TooltipProvider delayDuration={300}>
      <div className={fullHeight ? 'flex flex-col h-screen' : 'min-h-screen'}>
        {showHeader && <Header />}
        <main className={fullHeight ? 'flex-1 overflow-hidden' : 'flex-1'}>
          {children}
        </main>
      </div>
    </TooltipProvider>
  );
}
