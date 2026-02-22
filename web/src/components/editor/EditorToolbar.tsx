'use client';

import Link from 'next/link';
import {
  Undo2,
  Redo2,
  Download,
  PanelLeft,
  PanelRight,
  ZoomIn,
  ZoomOut,
  Loader2,
  AlertCircle,
  Save,
  ArrowLeft,
  Palette,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useEditorStore, useUIStore, usePreviewScalePercent } from '@/stores';
import type { TemplateId } from '@/types/cv';
import { cn } from '@/lib/utils';

const TEMPLATE_OPTIONS: { id: TemplateId; label: string }[] = [
  { id: 'classic', label: 'Classic' },
  { id: 'modern', label: 'Modern' },
  { id: 'visionary', label: 'Visionary' },
];

export function EditorToolbar() {
  const { canUndo, canRedo, undo, redo, saveStatus, saveError, saveNow, isDirty, cvContent, updateTemplateId, updateStyling } = useEditorStore();
  const {
    leftSidebarOpen,
    rightSidebarOpen,
    toggleLeftSidebar,
    toggleRightSidebar,
    zoomIn,
    zoomOut,
    resetZoom,
    openExportModal,
  } = useUIStore();
  const scalePercent = usePreviewScalePercent();

  const currentTemplateId = cvContent?.templateId || 'classic';
  const currentPrimaryColor = cvContent?.styling?.primaryColor || '#c0392b';

  return (
    <div className="h-[50px] border-b-2 border-ink bg-card px-5 flex items-center gap-1">
      {/* Back to Dashboard */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Link
            href="/"
            className="w-8 h-8 flex items-center justify-center border border-transparent bg-transparent text-[hsl(var(--text-secondary))] cursor-pointer transition-all duration-150 hover:text-primary hover:bg-[hsl(var(--vermillion-pale))]"
          >
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </TooltipTrigger>
        <TooltipContent>Back to Dashboard</TooltipContent>
      </Tooltip>

      <div className="w-px h-5 bg-border mx-1.5" />

      {/* Left Panel Toggle */}
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type="button"
            className={cn(
              'w-8 h-8 flex items-center justify-center border border-transparent bg-transparent cursor-pointer transition-all duration-150',
              leftSidebarOpen
                ? 'text-accent bg-[hsl(var(--vermillion-pale))]'
                : 'text-[hsl(var(--text-secondary))] hover:text-primary hover:bg-[hsl(var(--vermillion-pale))]'
            )}
            onClick={toggleLeftSidebar}
          >
            <PanelLeft className="h-4 w-4" />
          </button>
        </TooltipTrigger>
        <TooltipContent>
          {leftSidebarOpen ? 'Hide section palette' : 'Show section palette'}
        </TooltipContent>
      </Tooltip>

      <div className="w-px h-5 bg-border mx-1.5" />

      {/* Undo/Redo */}
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type="button"
            className="w-8 h-8 flex items-center justify-center border border-transparent bg-transparent text-[hsl(var(--text-secondary))] cursor-pointer transition-all duration-150 hover:text-primary hover:bg-[hsl(var(--vermillion-pale))] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent disabled:hover:text-[hsl(var(--text-secondary))]"
            disabled={!canUndo()}
            onClick={undo}
          >
            <Undo2 className="h-4 w-4" />
          </button>
        </TooltipTrigger>
        <TooltipContent>Undo (Ctrl+Z)</TooltipContent>
      </Tooltip>

      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type="button"
            className="w-8 h-8 flex items-center justify-center border border-transparent bg-transparent text-[hsl(var(--text-secondary))] cursor-pointer transition-all duration-150 hover:text-primary hover:bg-[hsl(var(--vermillion-pale))] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent disabled:hover:text-[hsl(var(--text-secondary))]"
            disabled={!canRedo()}
            onClick={redo}
          >
            <Redo2 className="h-4 w-4" />
          </button>
        </TooltipTrigger>
        <TooltipContent>Redo (Ctrl+Y)</TooltipContent>
      </Tooltip>

      <div className="w-px h-5 bg-border mx-1.5" />

      {/* Zoom Controls */}
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type="button"
            className="w-8 h-8 flex items-center justify-center border border-transparent bg-transparent text-[hsl(var(--text-secondary))] cursor-pointer transition-all duration-150 hover:text-primary hover:bg-[hsl(var(--vermillion-pale))]"
            onClick={zoomOut}
          >
            <ZoomOut className="h-4 w-4" />
          </button>
        </TooltipTrigger>
        <TooltipContent>Zoom out</TooltipContent>
      </Tooltip>

      <button
        type="button"
        className="font-mono text-xs text-[hsl(var(--text-secondary))] px-1.5 cursor-pointer hover:text-primary transition-colors duration-150"
        onClick={resetZoom}
      >
        {scalePercent}
      </button>

      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type="button"
            className="w-8 h-8 flex items-center justify-center border border-transparent bg-transparent text-[hsl(var(--text-secondary))] cursor-pointer transition-all duration-150 hover:text-primary hover:bg-[hsl(var(--vermillion-pale))]"
            onClick={zoomIn}
          >
            <ZoomIn className="h-4 w-4" />
          </button>
        </TooltipTrigger>
        <TooltipContent>Zoom in</TooltipContent>
      </Tooltip>

      <div className="w-px h-5 bg-border mx-1.5" />

      {/* Template Selector */}
      <Tooltip>
        <TooltipTrigger asChild>
          <select
            value={currentTemplateId}
            onChange={(e) => updateTemplateId(e.target.value as TemplateId)}
            className="h-8 px-2 text-xs border-2 border-input bg-secondary cursor-pointer focus:outline-none focus:border-primary focus:bg-card focus:shadow-offset-sm transition-all duration-150"
          >
            {TEMPLATE_OPTIONS.map((t) => (
              <option key={t.id} value={t.id}>
                {t.label}
              </option>
            ))}
          </select>
        </TooltipTrigger>
        <TooltipContent>CV template style</TooltipContent>
      </Tooltip>

      {/* Color Picker */}
      <Tooltip>
        <TooltipTrigger asChild>
          <label className="w-8 h-8 flex items-center justify-center border border-transparent bg-transparent text-[hsl(var(--text-secondary))] cursor-pointer transition-all duration-150 hover:text-primary hover:bg-[hsl(var(--vermillion-pale))] relative">
            <Palette className="h-4 w-4" />
            <input
              type="color"
              value={currentPrimaryColor}
              onChange={(e) => updateStyling({ primaryColor: e.target.value })}
              className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            />
          </label>
        </TooltipTrigger>
        <TooltipContent>Accent color</TooltipContent>
      </Tooltip>

      {/* Spacer */}
      <div className="flex-1" />

      {/* Save Status */}
      <div className="ml-auto flex items-center gap-1.5 font-mono text-[0.6875rem] text-muted-foreground uppercase tracking-[0.03em] mr-2">
        {saveStatus === 'saving' && (
          <>
            <Loader2 className="h-3.5 w-3.5 animate-spin" />
            <span className="hidden sm:inline">Saving...</span>
          </>
        )}
        {saveStatus === 'saved' && !isDirty && (
          <>
            <div className="w-1.5 h-1.5 rounded-full bg-[#3d7a3d]" />
            <span className="hidden sm:inline">Saved</span>
          </>
        )}
        {saveStatus === 'error' && (
          <Tooltip>
            <TooltipTrigger asChild>
              <button
                type="button"
                onClick={() => saveNow?.()}
                className="flex items-center gap-1.5 text-destructive hover:text-destructive/80 transition-colors"
              >
                <AlertCircle className="h-3.5 w-3.5" />
                <span className="hidden sm:inline">Save failed</span>
              </button>
            </TooltipTrigger>
            <TooltipContent>
              {saveError || 'Unknown error'} — click to retry
            </TooltipContent>
          </Tooltip>
        )}
        {(saveStatus === 'idle' && isDirty) && (
          <span className="hidden sm:inline">Unsaved</span>
        )}
      </div>

      {/* Manual Save */}
      {isDirty && saveStatus !== 'saving' && (
        <Tooltip>
          <TooltipTrigger asChild>
            <button
              type="button"
              className="w-8 h-8 flex items-center justify-center border border-transparent bg-transparent text-[hsl(var(--text-secondary))] cursor-pointer transition-all duration-150 hover:text-primary hover:bg-[hsl(var(--vermillion-pale))]"
              onClick={() => saveNow?.()}
            >
              <Save className="h-4 w-4" />
            </button>
          </TooltipTrigger>
          <TooltipContent>Save now (Ctrl+S)</TooltipContent>
        </Tooltip>
      )}

      <div className="w-px h-5 bg-border mx-1.5" />

      {/* Export */}
      <Button
        variant="default"
        size="sm"
        onClick={openExportModal}
      >
        <Download className="h-4 w-4 mr-2" />
        <span className="hidden sm:inline">Export</span>
      </Button>

      <div className="w-px h-5 bg-border mx-1.5" />

      {/* Right Panel Toggle */}
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type="button"
            className={cn(
              'w-8 h-8 flex items-center justify-center border border-transparent bg-transparent cursor-pointer transition-all duration-150',
              rightSidebarOpen
                ? 'text-accent bg-[hsl(var(--vermillion-pale))]'
                : 'text-[hsl(var(--text-secondary))] hover:text-primary hover:bg-[hsl(var(--vermillion-pale))]'
            )}
            onClick={toggleRightSidebar}
          >
            <PanelRight className="h-4 w-4" />
          </button>
        </TooltipTrigger>
        <TooltipContent>
          {rightSidebarOpen ? 'Hide properties' : 'Show properties'}
        </TooltipContent>
      </Tooltip>
    </div>
  );
}
