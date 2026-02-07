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
  Check,
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
import { Separator } from '@/components/ui/separator';
import { useEditorStore, useUIStore, usePreviewScalePercent } from '@/stores';
import type { TemplateId } from '@/types/cv';

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
  const currentPrimaryColor = cvContent?.styling?.primaryColor || '#2563eb';

  return (
    <div className="h-12 border-b bg-background px-2 flex items-center gap-1">
      {/* Back to Dashboard */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button variant="ghost" size="icon" className="h-8 w-8" asChild>
            <Link href="/">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
        </TooltipTrigger>
        <TooltipContent>Back to Dashboard</TooltipContent>
      </Tooltip>

      <Separator orientation="vertical" className="h-6 mx-1" />

      {/* Left Panel Toggle */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant={leftSidebarOpen ? 'secondary' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={toggleLeftSidebar}
          >
            <PanelLeft className="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>
          {leftSidebarOpen ? 'Hide section palette' : 'Show section palette'}
        </TooltipContent>
      </Tooltip>

      <Separator orientation="vertical" className="h-6 mx-1" />

      {/* Undo/Redo */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            disabled={!canUndo()}
            onClick={undo}
          >
            <Undo2 className="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>Undo (Ctrl+Z)</TooltipContent>
      </Tooltip>

      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            disabled={!canRedo()}
            onClick={redo}
          >
            <Redo2 className="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>Redo (Ctrl+Y)</TooltipContent>
      </Tooltip>

      <Separator orientation="vertical" className="h-6 mx-1" />

      {/* Zoom Controls */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={zoomOut}
          >
            <ZoomOut className="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>Zoom out</TooltipContent>
      </Tooltip>

      <Button
        variant="ghost"
        size="sm"
        className="h-8 px-2 min-w-[60px] text-xs"
        onClick={resetZoom}
      >
        {scalePercent}
      </Button>

      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={zoomIn}
          >
            <ZoomIn className="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>Zoom in</TooltipContent>
      </Tooltip>

      <Separator orientation="vertical" className="h-6 mx-1" />

      {/* Template Selector */}
      <Tooltip>
        <TooltipTrigger asChild>
          <select
            value={currentTemplateId}
            onChange={(e) => updateTemplateId(e.target.value as TemplateId)}
            className="h-8 px-2 text-xs border rounded bg-background hover:bg-accent cursor-pointer focus:outline-none focus:ring-2 focus:ring-ring"
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
          <label className="h-8 w-8 flex items-center justify-center cursor-pointer rounded hover:bg-accent transition-colors relative">
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
      <div className="flex items-center gap-2 text-sm text-muted-foreground mr-2">
        {saveStatus === 'saving' && (
          <>
            <Loader2 className="h-4 w-4 animate-spin" />
            <span className="hidden sm:inline">Saving...</span>
          </>
        )}
        {saveStatus === 'saved' && !isDirty && (
          <>
            <Check className="h-4 w-4 text-green-600" />
            <span className="hidden sm:inline">Saved</span>
          </>
        )}
        {saveStatus === 'error' && (
          <Tooltip>
            <TooltipTrigger asChild>
              <button
                onClick={() => saveNow?.()}
                className="flex items-center gap-2 text-destructive hover:text-destructive/80 transition-colors"
              >
                <AlertCircle className="h-4 w-4" />
                <span className="hidden sm:inline">Save failed — click to retry</span>
              </button>
            </TooltipTrigger>
            <TooltipContent>
              {saveError || 'Unknown error'}
            </TooltipContent>
          </Tooltip>
        )}
        {(saveStatus === 'idle' && isDirty) && (
          <span className="hidden sm:inline">Unsaved changes</span>
        )}
      </div>

      {/* Manual Save */}
      {isDirty && saveStatus !== 'saving' && (
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={() => saveNow?.()}
            >
              <Save className="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>Save now (Ctrl+S)</TooltipContent>
        </Tooltip>
      )}

      <Separator orientation="vertical" className="h-6 mx-1" />

      {/* Export */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant="outline"
            size="sm"
            className="h-8"
            onClick={openExportModal}
          >
            <Download className="h-4 w-4 mr-2" />
            <span className="hidden sm:inline">Export</span>
          </Button>
        </TooltipTrigger>
        <TooltipContent>Export as PDF, DOCX, or JSON</TooltipContent>
      </Tooltip>

      <Separator orientation="vertical" className="h-6 mx-1" />

      {/* Right Panel Toggle */}
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant={rightSidebarOpen ? 'secondary' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={toggleRightSidebar}
          >
            <PanelRight className="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>
          {rightSidebarOpen ? 'Hide properties' : 'Show properties'}
        </TooltipContent>
      </Tooltip>
    </div>
  );
}
