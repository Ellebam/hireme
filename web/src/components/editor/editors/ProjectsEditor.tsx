'use client';

import { useState, useCallback, useEffect } from 'react';
import { Plus, Trash2, Edit, GripVertical, FolderKanban } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { MonthYearPicker } from '@/components/ui/month-year-picker';
import { TagInput } from '@/components/ui/tag-input';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { useEditorStore } from '@/stores';
import type { ProjectsContent, ProjectEntry } from '@/types/cv';
import { generateId } from '@/lib/utils';

interface ProjectsEditorProps {
  sectionId: string;
  content: ProjectsContent;
}

export function ProjectsEditor({ sectionId, content }: ProjectsEditorProps) {
  const { updateSectionContent } = useEditorStore();
  const [editingEntry, setEditingEntry] = useState<ProjectEntry | null>(null);
  const [isNew, setIsNew] = useState(false);

  const entries = content.entries || [];

  const addEntry = useCallback(() => {
    const newEntry: ProjectEntry = {
      id: generateId(),
      name: '',
      role: '',
      description: '',
      url: '',
      technologies: [],
      startDate: '',
      endDate: null,
    };
    setEditingEntry(newEntry);
    setIsNew(true);
  }, []);

  const editEntry = useCallback((entry: ProjectEntry) => {
    setEditingEntry({ ...entry });
    setIsNew(false);
  }, []);

  const saveEntry = useCallback(
    (entry: ProjectEntry) => {
      let newEntries: ProjectEntry[];
      if (isNew) {
        newEntries = [...entries, entry];
      } else {
        newEntries = entries.map((e) => (e.id === entry.id ? entry : e));
      }
      updateSectionContent(sectionId, { entries: newEntries });
      setEditingEntry(null);
    },
    [sectionId, entries, isNew, updateSectionContent]
  );

  const deleteEntry = useCallback(
    (id: string) => {
      const newEntries = entries.filter((e) => e.id !== id);
      updateSectionContent(sectionId, { entries: newEntries });
    },
    [sectionId, entries, updateSectionContent]
  );

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          {entries.length} {entries.length === 1 ? 'entry' : 'entries'}
        </p>
        <Button onClick={addEntry} size="sm">
          <Plus className="h-4 w-4 mr-1" />
          Add Project
        </Button>
      </div>

      {entries.length === 0 ? (
        <div className="text-center py-8 border border-dashed rounded-lg">
          <FolderKanban className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
          <p className="text-muted-foreground mb-2">No projects added</p>
          <Button variant="outline" size="sm" onClick={addEntry}>
            <Plus className="h-4 w-4 mr-1" />
            Add Project
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          {entries.map((entry) => (
            <ProjectEntryCard
              key={entry.id}
              entry={entry}
              onEdit={() => editEntry(entry)}
              onDelete={() => deleteEntry(entry.id)}
            />
          ))}
        </div>
      )}

      {/* Edit Modal */}
      <ProjectEntryModal
        entry={editingEntry}
        open={!!editingEntry}
        onClose={() => setEditingEntry(null)}
        onSave={saveEntry}
      />
    </div>
  );
}

interface ProjectEntryCardProps {
  entry: ProjectEntry;
  onEdit: () => void;
  onDelete: () => void;
}

function ProjectEntryCard({ entry, onEdit, onDelete }: ProjectEntryCardProps) {
  return (
    <div className="flex items-center gap-3 p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors group">
      <div className="cursor-grab text-muted-foreground hover:text-foreground">
        <GripVertical className="h-4 w-4" />
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-medium truncate">{entry.name || 'Project'}</p>
        <p className="text-sm text-muted-foreground truncate">
          {entry.role}
          {entry.technologies && entry.technologies.length > 0 && (
            <>{entry.role ? ' · ' : ''}{entry.technologies.join(', ')}</>
          )}
        </p>
      </div>
      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <Button variant="ghost" size="icon" onClick={onEdit}>
          <Edit className="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          onClick={onDelete}
          className="text-destructive hover:text-destructive"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}

interface ProjectEntryModalProps {
  entry: ProjectEntry | null;
  open: boolean;
  onClose: () => void;
  onSave: (entry: ProjectEntry) => void;
}

function ProjectEntryModal({ entry, open, onClose, onSave }: ProjectEntryModalProps) {
  const [formData, setFormData] = useState<ProjectEntry | null>(null);

  // Sync form data when entry changes
  useEffect(() => {
    if (entry) {
      setFormData({ ...entry });
    }
  }, [entry]);

  if (!entry) return null;

  const data = formData || entry;

  const updateField = <K extends keyof ProjectEntry>(
    field: K,
    value: ProjectEntry[K]
  ) => {
    setFormData((prev) => (prev ? { ...prev, [field]: value } : null));
  };

  const handleSave = () => {
    if (!formData) return;
    onSave(formData);
  };

  return (
    <Dialog open={open} onOpenChange={(o) => !o && onClose()}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {formData?.name ? `Edit ${formData.name}` : 'Add Project'}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* Name & Role */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="name">Project Name</Label>
              <Input
                id="name"
                value={data.name}
                onChange={(e) => updateField('name', e.target.value)}
                placeholder="Open Source CLI Tool"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="role">Role</Label>
              <Input
                id="role"
                value={data.role || ''}
                onChange={(e) => updateField('role', e.target.value)}
                placeholder="Lead Developer"
              />
            </div>
          </div>

          {/* Description */}
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={data.description}
              onChange={(e) => updateField('description', e.target.value)}
              placeholder="Brief description of the project and your contributions..."
              className="min-h-[80px]"
            />
          </div>

          {/* Dates */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="startDate">Start Date</Label>
              <MonthYearPicker
                id="startDate"
                value={data.startDate || ''}
                onChange={(value) => updateField('startDate', value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="endDate">End Date</Label>
              <MonthYearPicker
                id="endDate"
                value={data.endDate || ''}
                onChange={(value) => updateField('endDate', value || null)}
              />
            </div>
          </div>

          {/* Technologies */}
          <div className="space-y-2">
            <Label htmlFor="technologies">Technologies</Label>
            <TagInput
              id="technologies"
              value={data.technologies || []}
              onChange={(tags) => updateField('technologies', tags)}
              placeholder="Type a technology and press Enter"
            />
          </div>

          {/* URL */}
          <div className="space-y-2">
            <Label htmlFor="url">Project URL</Label>
            <Input
              id="url"
              value={data.url || ''}
              onChange={(e) => updateField('url', e.target.value)}
              placeholder="https://github.com/example/project"
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleSave}>Save</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
