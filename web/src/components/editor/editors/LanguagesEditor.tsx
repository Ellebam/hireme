'use client';

import { useCallback } from 'react';
import { Plus, Trash2, GripVertical, Languages } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useEditorStore } from '@/stores';
import type { LanguagesContent, LanguageEntry, LanguageProficiency } from '@/types/cv';

const PROFICIENCY_LEVELS: { value: LanguageProficiency; label: string }[] = [
  { value: 'native', label: 'Native' },
  { value: 'fluent', label: 'Fluent' },
  { value: 'advanced', label: 'Advanced' },
  { value: 'intermediate', label: 'Intermediate' },
  { value: 'basic', label: 'Basic' },
];

interface LanguagesEditorProps {
  sectionId: string;
  content: LanguagesContent;
}

export function LanguagesEditor({ sectionId, content }: LanguagesEditorProps) {
  const { updateSectionContent } = useEditorStore();
  const entries = content.entries || [];

  const addEntry = useCallback(() => {
    const newEntry: LanguageEntry = {
      language: '',
      proficiency: 'intermediate',
    };
    updateSectionContent(sectionId, {
      entries: [...entries, newEntry],
    });
  }, [sectionId, entries, updateSectionContent]);

  const updateEntry = useCallback(
    (index: number, updates: Partial<LanguageEntry>) => {
      const newEntries = entries.map((entry, i) =>
        i === index ? { ...entry, ...updates } : entry
      );
      updateSectionContent(sectionId, { entries: newEntries });
    },
    [sectionId, entries, updateSectionContent]
  );

  const deleteEntry = useCallback(
    (index: number) => {
      const newEntries = entries.filter((_, i) => i !== index);
      updateSectionContent(sectionId, { entries: newEntries });
    },
    [sectionId, entries, updateSectionContent]
  );

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          {entries.length} {entries.length === 1 ? 'language' : 'languages'}
        </p>
        <Button onClick={addEntry} size="sm">
          <Plus className="h-4 w-4 mr-1" />
          Add Language
        </Button>
      </div>

      {entries.length === 0 ? (
        <div className="text-center py-8 border border-dashed rounded-lg">
          <Languages className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
          <p className="text-muted-foreground mb-2">No languages added</p>
          <Button variant="outline" size="sm" onClick={addEntry}>
            <Plus className="h-4 w-4 mr-1" />
            Add Your First Language
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          {entries.map((entry, index) => (
            <LanguageEntryRow
              key={`${entry.language}-${entry.proficiency}-${index}`}
              entry={entry}
              onUpdate={(updates) => updateEntry(index, updates)}
              onDelete={() => deleteEntry(index)}
            />
          ))}
        </div>
      )}
    </div>
  );
}

interface LanguageEntryRowProps {
  entry: LanguageEntry;
  onUpdate: (updates: Partial<LanguageEntry>) => void;
  onDelete: () => void;
}

function LanguageEntryRow({ entry, onUpdate, onDelete }: LanguageEntryRowProps) {
  return (
    <div className="flex items-center gap-3 p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors group">
      <div className="cursor-grab text-muted-foreground hover:text-foreground">
        <GripVertical className="h-4 w-4" />
      </div>

      <div className="flex-1 grid grid-cols-2 gap-3">
        <Input
          value={entry.language}
          onChange={(e) => onUpdate({ language: e.target.value })}
          placeholder="Language name"
          className="h-9"
        />
        <select
          value={entry.proficiency}
          onChange={(e) =>
            onUpdate({ proficiency: e.target.value as LanguageProficiency })
          }
          className="h-9 rounded-md border border-input bg-background px-3 text-sm"
        >
          {PROFICIENCY_LEVELS.map((level) => (
            <option key={level.value} value={level.value}>
              {level.label}
            </option>
          ))}
        </select>
      </div>

      <Input
        value={entry.certification || ''}
        onChange={(e) => onUpdate({ certification: e.target.value || undefined })}
        placeholder="Certification (optional)"
        className="h-9 w-40"
      />

      <Button
        variant="ghost"
        size="icon"
        onClick={onDelete}
        className="text-destructive hover:text-destructive opacity-0 group-hover:opacity-100 transition-opacity"
      >
        <Trash2 className="h-4 w-4" />
      </Button>
    </div>
  );
}
