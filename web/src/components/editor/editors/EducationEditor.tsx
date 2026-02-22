'use client';

import { useState, useCallback, useEffect } from 'react';
import { Plus, Trash2, Edit, GripVertical, GraduationCap } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { MonthYearPicker } from '@/components/ui/month-year-picker';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { useEditorStore } from '@/stores';
import type { EducationContent, EducationEntry } from '@/types/cv';
import { generateId } from '@/lib/utils';

interface EducationEditorProps {
  sectionId: string;
  content: EducationContent;
}

export function EducationEditor({ sectionId, content }: EducationEditorProps) {
  const { updateSectionContent } = useEditorStore();
  const [editingEntry, setEditingEntry] = useState<EducationEntry | null>(null);
  const [isNew, setIsNew] = useState(false);

  const entries = content.entries || [];

  const addEntry = useCallback(() => {
    const newEntry: EducationEntry = {
      id: generateId(),
      institution: '',
      degree: '',
      field: '',
      location: '',
      startDate: '',
      endDate: null,
      current: false,
      grade: '',
      description: '',
    };
    setEditingEntry(newEntry);
    setIsNew(true);
  }, []);

  const editEntry = useCallback((entry: EducationEntry) => {
    setEditingEntry({ ...entry });
    setIsNew(false);
  }, []);

  const saveEntry = useCallback(
    (entry: EducationEntry) => {
      let newEntries: EducationEntry[];
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
          Add Education
        </Button>
      </div>

      {entries.length === 0 ? (
        <div className="text-center py-8 border-2 border-dashed border-border rounded-lg">
          <GraduationCap className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
          <p className="text-muted-foreground mb-2">No education added</p>
          <Button variant="outline" size="sm" onClick={addEntry}>
            <Plus className="h-4 w-4 mr-1" />
            Add Education
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          {entries.map((entry) => (
            <EducationEntryCard
              key={entry.id}
              entry={entry}
              onEdit={() => editEntry(entry)}
              onDelete={() => deleteEntry(entry.id)}
            />
          ))}
        </div>
      )}

      {/* Edit Modal */}
      <EducationEntryModal
        entry={editingEntry}
        open={!!editingEntry}
        onClose={() => setEditingEntry(null)}
        onSave={saveEntry}
      />
    </div>
  );
}

interface EducationEntryCardProps {
  entry: EducationEntry;
  onEdit: () => void;
  onDelete: () => void;
}

function EducationEntryCard({ entry, onEdit, onDelete }: EducationEntryCardProps) {
  return (
    <div className="flex items-center gap-3 p-3 rounded-lg border-2 bg-card hover:bg-[hsl(var(--vermillion-pale))] transition-colors group">
      <div className="cursor-grab text-muted-foreground hover:text-foreground">
        <GripVertical className="h-4 w-4" />
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-medium truncate">
          {entry.degree}
          {entry.field && ` in ${entry.field}`}
        </p>
        <p className="text-sm text-muted-foreground truncate">
          {entry.institution}
          {entry.startDate && ` · ${entry.startDate}`}
          {entry.current ? ' - Present' : entry.endDate ? ` - ${entry.endDate}` : ''}
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

interface EducationEntryModalProps {
  entry: EducationEntry | null;
  open: boolean;
  onClose: () => void;
  onSave: (entry: EducationEntry) => void;
}

function EducationEntryModal({ entry, open, onClose, onSave }: EducationEntryModalProps) {
  const [formData, setFormData] = useState<EducationEntry | null>(null);

  // Sync form data when entry changes
  useEffect(() => {
    if (entry) {
      setFormData({ ...entry });
    }
  }, [entry]);

  if (!entry) return null;

  const data = formData || entry;

  const updateField = <K extends keyof EducationEntry>(
    field: K,
    value: EducationEntry[K]
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
            {formData?.institution ? `Edit ${formData.institution}` : 'Add Education'}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* Institution */}
          <div className="space-y-2">
            <Label htmlFor="institution">Institution</Label>
            <Input
              id="institution"
              value={data.institution}
              onChange={(e) => updateField('institution', e.target.value)}
              placeholder="Stanford University"
            />
          </div>

          {/* Degree & Field */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="degree">Degree</Label>
              <Input
                id="degree"
                value={data.degree}
                onChange={(e) => updateField('degree', e.target.value)}
                placeholder="Bachelor of Science"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="field">Field of Study</Label>
              <Input
                id="field"
                value={data.field || ''}
                onChange={(e) => updateField('field', e.target.value)}
                placeholder="Computer Science"
              />
            </div>
          </div>

          {/* Location */}
          <div className="space-y-2">
            <Label htmlFor="location">Location</Label>
            <Input
              id="location"
              value={data.location || ''}
              onChange={(e) => updateField('location', e.target.value)}
              placeholder="Stanford, CA"
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
                disabled={data.current}
              />
            </div>
          </div>

          {/* Currently Studying Toggle */}
          <div className="flex items-center space-x-2">
            <Switch
              id="current"
              checked={data.current || false}
              onCheckedChange={(checked) => {
                updateField('current', checked);
                if (checked) updateField('endDate', null);
              }}
            />
            <Label htmlFor="current">I am currently studying here</Label>
          </div>

          {/* Grade */}
          <div className="space-y-2">
            <Label htmlFor="grade">Grade / GPA</Label>
            <Input
              id="grade"
              value={data.grade || ''}
              onChange={(e) => updateField('grade', e.target.value)}
              placeholder="3.8 GPA, Magna Cum Laude"
            />
          </div>

          {/* Description */}
          <div className="space-y-2">
            <Label htmlFor="description">Additional Details</Label>
            <Textarea
              id="description"
              value={data.description || ''}
              onChange={(e) => updateField('description', e.target.value)}
              placeholder="Relevant coursework, activities, honors..."
              className="min-h-[80px]"
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
