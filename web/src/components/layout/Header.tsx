'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Menu, X, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { useState } from 'react';

const navItems = [
  { href: '/', label: 'Dashboard' },
  { href: '/editor', label: 'Editor' },
  { href: '/templates', label: 'Templates' },
];

export function Header() {
  const pathname = usePathname();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="sticky top-0 z-50 w-full h-[60px] bg-card border-b-2 border-ink">
      <div className="flex h-full items-center px-9">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-3">
          <div className="w-[34px] h-[34px] border-2 border-accent flex items-center justify-center font-serif font-bold text-[17px] text-accent -rotate-[2deg]">
            H
          </div>
          <span className="font-serif text-2xl font-bold text-primary tracking-[-0.03em]">
            HireMe
          </span>
        </Link>

        {/* Desktop Navigation */}
        <nav className="ml-8 hidden md:flex items-center gap-0">
          {navItems.map((item) => {
            const isActive = item.href === '/'
              ? pathname === '/'
              : pathname === item.href || pathname.startsWith(`${item.href}/`);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'relative px-5 py-[18px] text-xs font-semibold uppercase tracking-[0.125em] transition-colors duration-200',
                  isActive
                    ? 'text-primary after:absolute after:bottom-[-2px] after:left-5 after:right-5 after:h-0.5 after:bg-accent'
                    : 'text-[hsl(var(--text-secondary))] hover:text-primary'
                )}
              >
                {item.label}
              </Link>
            );
          })}
        </nav>

        {/* Spacer */}
        <div className="flex-1" />

        {/* New CV Button + User Avatar */}
        <div className="hidden md:flex items-center gap-3">
          <Button asChild size="sm">
            <Link href="/editor">
              <Plus className="h-3.5 w-3.5 mr-1.5" />
              New CV
            </Link>
          </Button>
          <div className="w-8 h-8 border-2 border-border flex items-center justify-center font-serif font-semibold text-[13px] text-sienna cursor-pointer">
            U
          </div>
        </div>

        {/* Mobile Menu Button */}
        <Button
          variant="ghost"
          size="icon"
          className="md:hidden"
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          aria-label={mobileMenuOpen ? 'Close menu' : 'Open menu'}
          aria-expanded={mobileMenuOpen}
        >
          {mobileMenuOpen ? (
            <X className="h-5 w-5" />
          ) : (
            <Menu className="h-5 w-5" />
          )}
        </Button>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden border-t-2 border-ink bg-card">
          <nav className="flex flex-col p-4 gap-1">
            {navItems.map((item) => {
              const isActive = item.href === '/'
                ? pathname === '/'
                : pathname === item.href || pathname.startsWith(`${item.href}/`);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setMobileMenuOpen(false)}
                  className={cn(
                    'px-5 py-3 text-xs font-semibold uppercase tracking-[0.125em] transition-colors duration-200',
                    isActive
                      ? 'text-primary border-l-[3px] border-accent bg-[hsl(var(--vermillion-pale))]'
                      : 'text-[hsl(var(--text-secondary))] hover:text-primary border-l-[3px] border-transparent'
                  )}
                >
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>
      )}
    </header>
  );
}
