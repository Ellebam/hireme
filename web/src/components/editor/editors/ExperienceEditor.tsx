'use client';

import { useState, useCallback } from 'react';
import { Plus, Trash2, Edit, GripVertical } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
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
import type { ExperienceContent, ExperienceEntry } from '@/types/cv';
import { generateId } from '@/lib/utils';

interface ExperienceEditorProps {
  sectionId: string;
  content: ExperienceContent;
}

export function ExperienceEditor({ sectionId, content }: ExperienceEditorProps) {
  const { updateSectionContent } = useEditorStore();
  const [editingEntry, setEditingEntry] = useState<ExperienceEntry | null>(null);
  const [isNew, setIsNew] = useState(false);

  const entries = content.entries || [];

  const addEntry = useCallback(() => {
    const newEntry: ExperienceEntry = {
      id: generateId(),
      company: '',
      position: '',
      location: '',
      startDate: '',
      endDate: null,
      current: false,
      description: '',
      highlights: [],
    };
    setEditingEntry(newEntry);
    setIsNew(true);
  }, []);

  const editEntry = useCallback((entry: ExperienceEntry) => {
    setEditingEntry({ ...entry });
    setIsNew(false);
  }, []);

  const saveEntry = useCallback(
    (entry: ExperienceEntry) => {
      let newEntries: ExperienceEntry[];
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
          Add Experience
        </Button>
      </div>

      {entries.length === 0 ? (
        <div className="text-center py-8 border border-dashed rounded-lg">
          <p className="text-muted-foreground mb-2">No work experience added</p>
          <Button variant="outline" size="sm" onClick={addEntry}>
            <Plus className="h-4 w-4 mr-1" />
            Add Your First Job
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          {entries.map((entry) => (
            <ExperienceEntryCard
              key={entry.id}
              entry={entry}
              onEdit={() => editEntry(entry)}
              onDelete={() => deleteEntry(entry.id)}
            />
          ))}
        </div>
      )}

      {/* Edit Modal */}
      <ExperienceEntryModal
        entry={editingEntry}
        open={!!editingEntry}
        onClose={() => setEditingEntry(null)}
        onSave={saveEntry}
      />
    </div>
  );
}

interface ExperienceEntryCardProps {
  entry: ExperienceEntry;
  onEdit: () => void;
  onDelete: () => void;
}

function ExperienceEntryCard({ entry, onEdit, onDelete }: ExperienceEntryCardProps) {
  return (
    <div className="flex items-center gap-3 p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors group">
      <div className="cursor-grab text-muted-foreground hover:text-foreground">
        <GripVertical className="h-4 w-4" />
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-medium truncate">{entry.position || 'Position'}</p>
        <p className="text-sm text-muted-foreground truncate">
          {entry.company}
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

interface ExperienceEntryModalProps {
  entry: ExperienceEntry | null;
  open: boolean;
  onClose: () => void;
  onSave: (entry: ExperienceEntry) => void;
}

function ExperienceEntryModal({ entry, open, onClose, onSave }: ExperienceEntryModalProps) {
  const [formData, setFormData] = useState<ExperienceEntry | null>(null);
  const [highlightsText, setHighlightsText] = useState('');

  // Sync form data when entry changes
  useState(() => {
    if (entry) {
      setFormData({ ...entry });
      setHighlightsText((entry.highlights || []).join('\n'));
    }
  });

  if (!entry) return null;

  const data = formData || entry;

  const updateField = <K extends keyof ExperienceEntry>(
    field: K,
    value: ExperienceEntry[K]
  ) => {
    setFormData((prev) => (prev ? { ...prev, [field]: value } : null));
  };

  const handleSave = () => {
    if (!formData) return;
    const highlights = highlightsText
      .split('\n')
      .map((h) => h.trim())
      .filter(Boolean);
    onSave({ ...formData, highlights });
  };

  return (
    <Dialog open={open} onOpenChange={(o) => !o && onClose()}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {formData?.company ? `Edit ${formData.company}` : 'Add Work Experience'}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* Company & Position */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="company">Company</Label>
              <Input
                id="company"
                value={data.company}
                onChange={(e) => updateField('company', e.target.value)}
                placeholder="Acme Inc."
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="position">Position</Label>
              <Input
                id="position"
                value={data.position}
                onChange={(e) => updateField('position', e.target.value)}
                placeholder="Senior Developer"
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
              placeholder="San Francisco, CA"
            />
          </div>

          {/* Dates */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="startDate">Start Date</Label>
              <Input
                id="startDate"
                type="month"
                value={data.startDate}
                onChange={(e) => updateField('startDate', e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="endDate">End Date</Label>
              <Input
                id="endDate"
                type="month"
                value={data.endDate || ''}
                onChange={(e) => updateField('endDate', e.target.value || null)}
                disabled={data.current}
              />
            </div>
          </div>

          {/* Current Position Toggle */}
          <div className="flex items-center space-x-2">
            <Switch
              id="current"
              checked={data.current || false}
              onCheckedChange={(checked) => {
                updateField('current', checked);
                if (checked) updateField('endDate', null);
              }}
            />
            <Label htmlFor="current">I currently work here</Label>
          </div>

          {/* Description */}
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={data.description || ''}
              onChange={(e) => updateField('description', e.target.value)}
              placeholder="Brief description of your role and responsibilities..."
              className="min-h-[80px]"
            />
          </div>

          {/* Highlights */}
          <div className="space-y-2">
            <Label htmlFor="highlights">
              Key Achievements
              <span className="text-muted-foreground font-normal ml-2">
                (one per line)
              </span>
            </Label>
            <Textarea
              id="highlights"
              value={highlightsText}
              onChange={(e) => setHighlightsText(e.target.value)}
              placeholder="Increased team velocity by 30%&#10;Led migration to microservices&#10;Mentored 5 junior developers"
              className="min-h-[100px]"
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
