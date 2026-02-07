'use client';

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
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Separator } from '@/components/ui/separator';
import { useEditorStore, useUIStore, usePreviewScalePercent } from '@/stores';

export function EditorToolbar() {
  const { canUndo, canRedo, undo, redo, saveStatus, isDirty } = useEditorStore();
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

  return (
    <div className="h-12 border-b bg-background px-2 flex items-center gap-1">
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

      {/* Spacer */}
      <div className="flex-1" />

      {/* Save Status */}
      <div className="flex items-center gap-2 text-sm text-muted-foreground mr-2">
        {saveStatus === 'saving' && (
          <>
            <Loader2 className="h-4 w-4 animate-spin" />
            <span>Saving...</span>
          </>
        )}
        {saveStatus === 'saved' && !isDirty && (
          <>
            <Check className="h-4 w-4 text-green-600" />
            <span>Saved</span>
          </>
        )}
        {saveStatus === 'error' && (
          <>
            <AlertCircle className="h-4 w-4 text-destructive" />
            <span className="text-destructive">Error</span>
          </>
        )}
        {(saveStatus === 'idle' && isDirty) && (
          <span>Unsaved changes</span>
        )}
      </div>

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
            Export
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
