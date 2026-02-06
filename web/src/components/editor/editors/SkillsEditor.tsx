'use client';

import { useState, useCallback } from 'react';
import { Plus, Trash2, GripVertical, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useEditorStore } from '@/stores';
import type { SkillsContent, SkillCategory, Skill, SkillLevel } from '@/types/cv';
import { generateId } from '@/lib/utils';

const SKILL_LEVELS: { value: SkillLevel; label: string }[] = [
  { value: 'beginner', label: 'Beginner' },
  { value: 'intermediate', label: 'Intermediate' },
  { value: 'advanced', label: 'Advanced' },
  { value: 'expert', label: 'Expert' },
];

interface SkillsEditorProps {
  sectionId: string;
  content: SkillsContent;
}

export function SkillsEditor({ sectionId, content }: SkillsEditorProps) {
  const { updateSectionContent } = useEditorStore();
  const categories = content.categories || [];

  const addCategory = useCallback(() => {
    const newCategory: SkillCategory = {
      id: generateId(),
      name: 'New Category',
      skills: [],
    };
    updateSectionContent(sectionId, {
      categories: [...categories, newCategory],
    });
  }, [sectionId, categories, updateSectionContent]);

  const updateCategory = useCallback(
    (categoryId: string, updates: Partial<SkillCategory>) => {
      const newCategories = categories.map((cat) =>
        cat.id === categoryId ? { ...cat, ...updates } : cat
      );
      updateSectionContent(sectionId, { categories: newCategories });
    },
    [sectionId, categories, updateSectionContent]
  );

  const deleteCategory = useCallback(
    (categoryId: string) => {
      const newCategories = categories.filter((cat) => cat.id !== categoryId);
      updateSectionContent(sectionId, { categories: newCategories });
    },
    [sectionId, categories, updateSectionContent]
  );

  const addSkill = useCallback(
    (categoryId: string) => {
      const newCategories = categories.map((cat) => {
        if (cat.id === categoryId) {
          return {
            ...cat,
            skills: [...cat.skills, { name: '', level: undefined }],
          };
        }
        return cat;
      });
      updateSectionContent(sectionId, { categories: newCategories });
    },
    [sectionId, categories, updateSectionContent]
  );

  const updateSkill = useCallback(
    (categoryId: string, skillIndex: number, updates: Partial<Skill>) => {
      const newCategories = categories.map((cat) => {
        if (cat.id === categoryId) {
          const newSkills = [...cat.skills];
          newSkills[skillIndex] = { ...newSkills[skillIndex], ...updates };
          return { ...cat, skills: newSkills };
        }
        return cat;
      });
      updateSectionContent(sectionId, { categories: newCategories });
    },
    [sectionId, categories, updateSectionContent]
  );

  const deleteSkill = useCallback(
    (categoryId: string, skillIndex: number) => {
      const newCategories = categories.map((cat) => {
        if (cat.id === categoryId) {
          const newSkills = cat.skills.filter((_, i) => i !== skillIndex);
          return { ...cat, skills: newSkills };
        }
        return cat;
      });
      updateSectionContent(sectionId, { categories: newCategories });
    },
    [sectionId, categories, updateSectionContent]
  );

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          {categories.length} {categories.length === 1 ? 'category' : 'categories'}
        </p>
        <Button onClick={addCategory} size="sm">
          <Plus className="h-4 w-4 mr-1" />
          Add Category
        </Button>
      </div>

      {categories.length === 0 ? (
        <div className="text-center py-8 border border-dashed rounded-lg">
          <p className="text-muted-foreground mb-2">No skill categories added</p>
          <Button variant="outline" size="sm" onClick={addCategory}>
            <Plus className="h-4 w-4 mr-1" />
            Add Your First Category
          </Button>
        </div>
      ) : (
        <div className="space-y-4">
          {categories.map((category) => (
            <SkillCategoryCard
              key={category.id}
              category={category}
              onUpdateCategory={(updates) => updateCategory(category.id, updates)}
              onDeleteCategory={() => deleteCategory(category.id)}
              onAddSkill={() => addSkill(category.id)}
              onUpdateSkill={(index, updates) =>
                updateSkill(category.id, index, updates)
              }
              onDeleteSkill={(index) => deleteSkill(category.id, index)}
            />
          ))}
        </div>
      )}
    </div>
  );
}

interface SkillCategoryCardProps {
  category: SkillCategory;
  onUpdateCategory: (updates: Partial<SkillCategory>) => void;
  onDeleteCategory: () => void;
  onAddSkill: () => void;
  onUpdateSkill: (index: number, updates: Partial<Skill>) => void;
  onDeleteSkill: (index: number) => void;
}

function SkillCategoryCard({
  category,
  onUpdateCategory,
  onDeleteCategory,
  onAddSkill,
  onUpdateSkill,
  onDeleteSkill,
}: SkillCategoryCardProps) {
  const [isEditing, setIsEditing] = useState(false);

  return (
    <div className="rounded-lg border bg-card p-4 space-y-3">
      {/* Category Header */}
      <div className="flex items-center gap-2">
        <div className="cursor-grab text-muted-foreground">
          <GripVertical className="h-4 w-4" />
        </div>
        {isEditing ? (
          <Input
            value={category.name}
            onChange={(e) => onUpdateCategory({ name: e.target.value })}
            onBlur={() => setIsEditing(false)}
            onKeyDown={(e) => e.key === 'Enter' && setIsEditing(false)}
            autoFocus
            className="h-8 font-medium"
          />
        ) : (
          <button
            onClick={() => setIsEditing(true)}
            className="font-medium text-left hover:text-primary"
          >
            {category.name || 'Untitled Category'}
          </button>
        )}
        <div className="flex-1" />
        <Button
          variant="ghost"
          size="icon"
          onClick={onDeleteCategory}
          className="h-8 w-8 text-destructive hover:text-destructive"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>

      {/* Skills */}
      <div className="flex flex-wrap gap-2">
        {category.skills.map((skill, index) => (
          <SkillTag
            key={index}
            skill={skill}
            onUpdate={(updates) => onUpdateSkill(index, updates)}
            onDelete={() => onDeleteSkill(index)}
          />
        ))}
        <Button
          variant="outline"
          size="sm"
          onClick={onAddSkill}
          className="h-8"
        >
          <Plus className="h-3 w-3 mr-1" />
          Add Skill
        </Button>
      </div>
    </div>
  );
}

interface SkillTagProps {
  skill: Skill;
  onUpdate: (updates: Partial<Skill>) => void;
  onDelete: () => void;
}

function SkillTag({ skill, onUpdate, onDelete }: SkillTagProps) {
  const [isEditing, setIsEditing] = useState(!skill.name);

  if (isEditing) {
    return (
      <div className="flex items-center gap-1 px-2 py-1 rounded-md bg-muted">
        <Input
          value={skill.name}
          onChange={(e) => onUpdate({ name: e.target.value })}
          onBlur={() => skill.name && setIsEditing(false)}
          onKeyDown={(e) => e.key === 'Enter' && skill.name && setIsEditing(false)}
          autoFocus
          placeholder="Skill name"
          className="h-6 w-24 text-xs px-1"
        />
        <select
          value={skill.level || ''}
          onChange={(e) =>
            onUpdate({ level: (e.target.value as SkillLevel) || undefined })
          }
          className="h-6 text-xs border rounded px-1 bg-background"
        >
          <option value="">No level</option>
          {SKILL_LEVELS.map((level) => (
            <option key={level.value} value={level.value}>
              {level.label}
            </option>
          ))}
        </select>
        <button
          onClick={onDelete}
          className="text-muted-foreground hover:text-destructive"
        >
          <X className="h-3 w-3" />
        </button>
      </div>
    );
  }

  return (
    <div
      className="flex items-center gap-1 px-3 py-1 rounded-md bg-muted hover:bg-accent cursor-pointer group"
      onClick={() => setIsEditing(true)}
    >
      <span className="text-sm">{skill.name}</span>
      {skill.level && (
        <span className="text-xs text-muted-foreground capitalize">
          ({skill.level})
        </span>
      )}
      <button
        onClick={(e) => {
          e.stopPropagation();
          onDelete();
        }}
        className="opacity-0 group-hover:opacity-100 text-muted-foreground hover:text-destructive ml-1"
      >
        <X className="h-3 w-3" />
      </button>
    </div>
  );
}
