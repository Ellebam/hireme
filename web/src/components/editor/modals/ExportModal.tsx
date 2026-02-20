'use client';

import type React from 'react';
import { useState, useRef, useEffect } from 'react';
import { FileText, FileJson, FileType, Download, Loader2, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { useUIStore } from '@/stores';
import { useEditorStore } from '@/stores';
import { api } from '@/lib/api';
import type { ExportFormat } from '@/types/api';
import { cn } from '@/lib/utils';

const exportFormats: {
  id: ExportFormat;
  name: string;
  description: string;
  icon: React.ElementType;
}[] = [
  {
    id: 'pdf',
    name: 'PDF',
    description: 'Best for sharing and printing',
    icon: FileText,
  },
  {
    id: 'docx',
    name: 'Word (DOCX)',
    description: 'Editable document format',
    icon: FileType,
  },
  {
    id: 'json',
    name: 'JSON',
    description: 'Raw data for backup or import',
    icon: FileJson,
  },
];

export function ExportModal() {
  const { exportModalOpen, closeExportModal } = useUIStore();
  const { cv, cvContent } = useEditorStore();
  const [selectedFormat, setSelectedFormat] = useState<ExportFormat>('pdf');
  const [isExporting, setIsExporting] = useState(false);
  const [exportSuccess, setExportSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const closeTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Clear timeout on unmount
  useEffect(() => {
    return () => {
      if (closeTimeoutRef.current) {
        clearTimeout(closeTimeoutRef.current);
      }
    };
  }, []);

  const handleExport = async () => {
    if (!cv) return;

    setIsExporting(true);
    setError(null);
    setExportSuccess(false);

    try {
      if (selectedFormat === 'json') {
        // For JSON, just download the content directly
        const blob = new Blob([JSON.stringify(cvContent, null, 2)], {
          type: 'application/json',
        });
        downloadBlob(blob, `${cv.title || 'cv'}.json`);
      } else {
        // For PDF/DOCX, call the export API
        const response = await api.export.export(cv.id, selectedFormat);

        const filename = `${cv.title || 'cv'}.${selectedFormat}`;
        downloadBlob(response, filename);
      }

      setExportSuccess(true);
      closeTimeoutRef.current = setTimeout(() => {
        closeExportModal();
        setExportSuccess(false);
      }, 1500);
    } catch (err) {
      console.error('Export failed:', err);
      setError(
        selectedFormat === 'json'
          ? 'Failed to export JSON'
          : 'Export is not yet available. The Gotenberg service needs to be configured.'
      );
    } finally {
      setIsExporting(false);
    }
  };

  const downloadBlob = (blob: Blob, filename: string) => {
    const url = URL.createObjectURL(blob);
    const downloadLink = document.createElement('a');
    downloadLink.href = url;
    downloadLink.download = filename;
    document.body.appendChild(downloadLink);
    downloadLink.click();
    // Delay cleanup so the browser can start the download before the blob URL is revoked
    setTimeout(() => {
      document.body.removeChild(downloadLink);
      URL.revokeObjectURL(url);
    }, 150);
  };

  return (
    <Dialog open={exportModalOpen} onOpenChange={(open) => !open && closeExportModal()}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="font-serif text-[1.625rem] font-bold tracking-[-0.02em]">Export CV</DialogTitle>
          <DialogDescription>
            Choose a format to export your CV
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <div className="border-t border-dashed border-border">
            {exportFormats.map((format) => {
              const Icon = format.icon;
              const isSelected = selectedFormat === format.id;

              return (
                <button
                  type="button"
                  key={format.id}
                  onClick={() => setSelectedFormat(format.id)}
                  className={cn(
                    'w-full flex items-center gap-4 px-4 py-3 border-b border-dashed border-border text-left transition-all duration-150',
                    isSelected
                      ? 'bg-[hsl(var(--vermillion-pale))] border-l-[3px] border-l-accent pl-3.5'
                      : 'hover:bg-[hsl(var(--vermillion-pale))] border-l-[3px] border-l-transparent pl-3.5'
                  )}
                >
                  <Icon className={cn('h-5 w-5', isSelected ? 'text-accent' : 'text-[hsl(var(--text-secondary))]')} />
                  <div className="flex-1">
                    <p className="font-serif font-semibold text-sm">{format.name}</p>
                    <p className="text-[0.6875rem] text-[hsl(var(--text-secondary))]">
                      {format.description}
                    </p>
                  </div>
                  {isSelected && (
                    <div className="h-5 w-5 bg-accent flex items-center justify-center">
                      <Check className="h-3 w-3 text-accent-foreground" />
                    </div>
                  )}
                </button>
              );
            })}
          </div>
        </div>

        {error && (
          <div className="p-3 border-2 border-destructive bg-destructive/10 text-destructive text-sm">
            {error}
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={closeExportModal}>
            Cancel
          </Button>
          <Button onClick={handleExport} disabled={isExporting || exportSuccess}>
            {isExporting ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Exporting...
              </>
            ) : exportSuccess ? (
              <>
                <Check className="h-4 w-4 mr-2" />
                Downloaded!
              </>
            ) : (
              <>
                <Download className="h-4 w-4 mr-2" />
                Export {selectedFormat.toUpperCase()}
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
