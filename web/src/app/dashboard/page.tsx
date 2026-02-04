'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Plus, FileText, Clock, MoreVertical, Trash2, Edit } from 'lucide-react';
import { AppShell } from '@/components/layout';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { api, ApiError } from '@/lib/api';
import type { CV } from '@/types/api';

export default function DashboardPage() {
  const [cv, setCV] = useState<CV | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadCV() {
      try {
        const data = await api.cv.get();
        setCV(data);
      } catch (err) {
        if (err instanceof ApiError && err.isNotFound) {
          // No CV exists yet - that's OK
          setCV(null);
        } else {
          setError('Failed to load CV');
          console.error('Error loading CV:', err);
        }
      } finally {
        setLoading(false);
      }
    }

    loadCV();
  }, []);

  return (
    <AppShell>
      <div className="container mx-auto py-8 px-4">
        {/* Page Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
            <p className="text-muted-foreground mt-1">
              Manage your CV and track your progress
            </p>
          </div>
          <Button asChild>
            <Link href="/editor">
              <Plus className="mr-2 h-4 w-4" />
              {cv ? 'Edit CV' : 'Create CV'}
            </Link>
          </Button>
        </div>

        {/* Content */}
        {loading ? (
          <DashboardSkeleton />
        ) : error ? (
          <ErrorState message={error} />
        ) : cv ? (
          <CVCard cv={cv} />
        ) : (
          <EmptyState />
        )}
      </div>
    </AppShell>
  );
}

function CVCard({ cv }: { cv: CV }) {
  const updatedAt = new Date(cv.updatedAt);
  const relativeTime = getRelativeTime(updatedAt);

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <Card className="group relative hover:shadow-md transition-shadow">
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-primary/10 text-primary">
                <FileText className="h-5 w-5" />
              </div>
              <div>
                <CardTitle className="text-lg">{cv.title}</CardTitle>
                <CardDescription className="flex items-center gap-1 mt-1">
                  <Clock className="h-3 w-3" />
                  Updated {relativeTime}
                </CardDescription>
              </div>
            </div>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-4">
            <span className="px-2 py-1 rounded-md bg-muted text-xs font-medium">
              {cv.content.templateId}
            </span>
            <span>·</span>
            <span>{cv.content.sections?.length || 0} sections</span>
          </div>
          <div className="flex gap-2">
            <Button className="flex-1" asChild>
              <Link href={`/editor/${cv.id}`}>
                <Edit className="mr-2 h-4 w-4" />
                Edit
              </Link>
            </Button>
            <Button variant="outline" size="icon">
              <Trash2 className="h-4 w-4 text-destructive" />
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Create New Card */}
      <Card className="border-dashed hover:border-primary/50 hover:bg-accent/50 transition-colors cursor-pointer">
        <Link href="/editor" className="flex flex-col items-center justify-center h-full min-h-[200px] p-6">
          <div className="p-3 rounded-full bg-muted mb-3">
            <Plus className="h-6 w-6 text-muted-foreground" />
          </div>
          <p className="font-medium text-muted-foreground">Create New CV</p>
        </Link>
      </Card>
    </div>
  );
}

function EmptyState() {
  return (
    <Card className="border-dashed">
      <CardContent className="flex flex-col items-center justify-center py-16">
        <div className="p-4 rounded-full bg-muted mb-4">
          <FileText className="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 className="text-xl font-semibold mb-2">No CV yet</h3>
        <p className="text-muted-foreground text-center max-w-sm mb-6">
          Create your first professional CV and start landing your dream job.
        </p>
        <Button asChild>
          <Link href="/editor">
            <Plus className="mr-2 h-4 w-4" />
            Create Your CV
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}

function ErrorState({ message }: { message: string }) {
  return (
    <Card className="border-destructive/50">
      <CardContent className="flex flex-col items-center justify-center py-16">
        <p className="text-destructive font-medium mb-4">{message}</p>
        <Button variant="outline" onClick={() => window.location.reload()}>
          Try Again
        </Button>
      </CardContent>
    </Card>
  );
}

function DashboardSkeleton() {
  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center gap-3">
            <Skeleton className="h-10 w-10 rounded-lg" />
            <div className="space-y-2">
              <Skeleton className="h-5 w-32" />
              <Skeleton className="h-4 w-24" />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-4 w-full mb-4" />
          <div className="flex gap-2">
            <Skeleton className="h-10 flex-1" />
            <Skeleton className="h-10 w-10" />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function getRelativeTime(date: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) return 'just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return date.toLocaleDateString();
}
