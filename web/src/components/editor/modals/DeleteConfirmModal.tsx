'use client';

import { AlertTriangle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { useUIStore, useEditorStore } from '@/stores';
import { SECTION_LABELS } from '@/types/cv';

export function DeleteConfirmModal() {
  const { deleteConfirmModalOpen, deleteSectionId, closeDeleteConfirmModal } =
    useUIStore();
  const { cvContent, deleteSection } = useEditorStore();

  const section = deleteSectionId
    ? cvContent?.sections.find((s) => s.id === deleteSectionId)
    : null;

  const sectionName = section
    ? section.title || SECTION_LABELS[section.type]
    : 'this section';

  const handleDelete = () => {
    if (deleteSectionId) {
      deleteSection(deleteSectionId);
      closeDeleteConfirmModal();
    }
  };

  return (
    <Dialog
      open={deleteConfirmModalOpen}
      onOpenChange={(open) => !open && closeDeleteConfirmModal()}
    >
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-full bg-destructive/10">
              <AlertTriangle className="h-5 w-5 text-destructive" />
            </div>
            <div>
              <DialogTitle>Delete Section</DialogTitle>
              <DialogDescription className="mt-1">
                Are you sure you want to delete &ldquo;{sectionName}&rdquo;?
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="py-2">
          <p className="text-sm text-muted-foreground">
            This action cannot be undone. All content in this section will be
            permanently removed.
          </p>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={closeDeleteConfirmModal}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            Delete
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
