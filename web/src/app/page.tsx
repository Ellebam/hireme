'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Plus, MoreVertical, Trash2, Edit } from 'lucide-react';
import { AppShell } from '@/components/layout';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { api } from '@/lib/api';
import type { CV } from '@/types/api';
import { useEffect } from 'react';

export default function DashboardPage() {
  const [cvs, setCVs] = useState<CV[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadCVs() {
      try {
        const data = await api.cv.list();
        setCVs(data);
      } catch {
        setError('Failed to load CVs');
      } finally {
        setLoading(false);
      }
    }

    loadCVs();
  }, []);

  const handleDelete = async (id: string) => {
    try {
      await api.cv.delete(id);
      setCVs((prev) => prev.filter((cv) => cv.id !== id));
    } catch {
      setError('Failed to delete CV');
    }
  };

  return (
    <AppShell>
      <div className="max-w-[840px] px-10">
        {/* Hero Section */}
        <div className="pt-[72px]">
          <div className="flex items-center gap-3 mb-5 animate-slide-in opacity-0 [animation-delay:0.1s]">
            <div className="w-9 h-0.5 bg-accent" />
            <span className="font-mono text-[0.6875rem] font-medium uppercase tracking-[0.1875em] text-accent">
              Your Workspace
            </span>
          </div>
          <h1 className="font-serif text-[3.5rem] font-bold leading-[1.08] tracking-[-0.03em] text-primary mb-5 animate-slide-in opacity-0 [animation-delay:0.2s]">
            Write the story<br />they <em className="font-normal italic text-sienna">remember.</em>
          </h1>
          <p className="text-[1.0625rem] text-[hsl(var(--text-secondary))] leading-[1.7] max-w-[500px] mb-9 animate-slide-in opacity-0 [animation-delay:0.35s]">
            A CV is not a list. It&apos;s a carefully composed letter &mdash; part design, part craft, all you.
          </p>
          <div className="flex gap-3 mb-14 animate-slide-in opacity-0 [animation-delay:0.45s]">
            <Button asChild variant="accent">
              <Link href="/templates">
                <Plus className="mr-2 h-4 w-4" />
                Create CV
              </Link>
            </Button>
            <Button asChild variant="outline">
              <Link href="/templates">
                Templates &rarr;
              </Link>
            </Button>
          </div>
        </div>

        {/* Divider */}
        <div className="border-t-2 border-ink -mx-10 px-10" />

        {/* Content */}
        {loading ? (
          <DashboardSkeleton />
        ) : error ? (
          <ErrorState message={error} />
        ) : cvs.length > 0 ? (
          <DocumentList cvs={cvs} onDelete={handleDelete} />
        ) : (
          <EmptyState />
        )}
      </div>
    </AppShell>
  );
}

function DocumentList({ cvs, onDelete }: { cvs: CV[]; onDelete: (id: string) => void }) {
  return (
    <>
      {/* Section Header */}
      <div className="flex items-baseline justify-between py-6 pb-4">
        <h2 className="font-serif text-[1.625rem] font-semibold tracking-[-0.02em]">
          Documents
        </h2>
        <span className="font-mono text-xs text-muted-foreground bg-secondary px-2 py-0.5">
          {cvs.length}
        </span>
      </div>

      {/* Document Rows */}
      {cvs.map((cv, index) => (
        <DocumentRow
          key={cv.id}
          cv={cv}
          index={index}
          onDelete={() => onDelete(cv.id)}
        />
      ))}

      {/* Create Row */}
      <Link
        href="/templates"
        className="flex items-center gap-4 py-6 cursor-pointer text-muted-foreground transition-all duration-200 hover:text-accent"
      >
        <span className="font-serif text-2xl font-light">+</span>
        <span className="text-[0.8125rem] font-semibold uppercase tracking-[0.0625em]">
          Create New Document
        </span>
      </Link>
    </>
  );
}

function DocumentRow({ cv, index, onDelete }: { cv: CV; index: number; onDelete: () => void }) {
  const updatedAt = new Date(cv.updatedAt);
  const relativeTime = getRelativeTime(updatedAt);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  return (
    <>
      <Link
        href={`/editor/${cv.id}`}
        className="grid grid-cols-[48px_1fr_120px_100px_40px] items-center gap-5 py-5 border-b border-dashed border-border cursor-pointer transition-all duration-200 hover:bg-[hsl(var(--vermillion-pale))] hover:pl-3 hover:-ml-3 hover:-mr-3 hover:pr-3 hover:shadow-offset-md"
      >
        <span className="font-mono text-[0.8125rem] text-muted-foreground">
          {String(index + 1).padStart(2, '0')}
        </span>
        <div>
          <span className="font-serif text-[1.1875rem] font-semibold text-primary tracking-[-0.01em]">
            {cv.title}
          </span>
          <span className="block text-[0.8125rem] text-[hsl(var(--text-secondary))]">
            {cv.content?.sections?.length || 0} sections
          </span>
        </div>
        <span className="text-[0.8125rem] font-medium text-[hsl(var(--text-secondary))]">
          {cv.content?.templateId}
        </span>
        <span className="font-mono text-xs text-muted-foreground">
          {relativeTime}
        </span>
        <div className="flex justify-end" onClick={(e) => e.preventDefault()}>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link href={`/editor/${cv.id}`} className="flex items-center gap-2">
                  <Edit className="h-4 w-4" />
                  Edit
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-destructive focus:text-destructive"
                onClick={() => setDeleteDialogOpen(true)}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </Link>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="font-serif text-[1.625rem] font-bold tracking-[-0.02em]">
              Delete CV
            </DialogTitle>
            <DialogDescription>
              Are you sure you want to delete &quot;{cv.title}&quot;? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => {
                onDelete();
                setDeleteDialogOpen(false);
              }}
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

function EmptyState() {
  return (
    <div className="py-16 border-t border-dashed border-border">
      <div className="flex flex-col items-center justify-center text-center">
        <h3 className="font-serif text-xl font-semibold mb-2">No CVs yet</h3>
        <p className="text-[hsl(var(--text-secondary))] max-w-sm mb-6">
          Create your first professional CV and start landing your dream job.
        </p>
        <Button asChild>
          <Link href="/templates">
            <Plus className="mr-2 h-4 w-4" />
            Create Your CV
          </Link>
        </Button>
      </div>
    </div>
  );
}

function ErrorState({ message }: { message: string }) {
  return (
    <div className="py-16 border-t border-dashed border-border">
      <div className="flex flex-col items-center justify-center text-center">
        <p className="text-destructive font-medium mb-4">{message}</p>
        <Button variant="outline" onClick={() => window.location.reload()}>
          Try Again
        </Button>
      </div>
    </div>
  );
}

function DashboardSkeleton() {
  return (
    <div className="py-6">
      <div className="flex items-baseline justify-between pb-4">
        <Skeleton className="h-8 w-40" />
        <Skeleton className="h-5 w-8" />
      </div>
      <div className="border-b border-dashed border-border py-5">
        <div className="grid grid-cols-[48px_1fr_120px_100px_40px] items-center gap-5">
          <Skeleton className="h-4 w-6" />
          <div className="space-y-2">
            <Skeleton className="h-5 w-48" />
            <Skeleton className="h-4 w-24" />
          </div>
          <Skeleton className="h-4 w-16" />
          <Skeleton className="h-4 w-12" />
          <Skeleton className="h-8 w-8" />
        </div>
      </div>
    </div>
  );
}

function getRelativeTime(date: Date): string {
  const now = new Date();
  const diffMs = Math.max(0, now.getTime() - date.getTime());
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) return 'just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return date.toLocaleDateString();
}
