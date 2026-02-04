'use client';

import { useCallback } from 'react';
import { Plus, Trash2, Globe, Linkedin, Github, Twitter } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useEditorStore } from '@/stores';
import type { PersonalContent, ProfileLink, LinkType } from '@/types/cv';
import { generateId } from '@/lib/utils';

const LINK_ICONS: Record<LinkType, React.ReactNode> = {
  linkedin: <Linkedin className="h-4 w-4" />,
  github: <Github className="h-4 w-4" />,
  twitter: <Twitter className="h-4 w-4" />,
  website: <Globe className="h-4 w-4" />,
  portfolio: <Globe className="h-4 w-4" />,
  other: <Globe className="h-4 w-4" />,
};

interface PersonalEditorProps {
  sectionId: string;
  content: PersonalContent;
}

export function PersonalEditor({ sectionId, content }: PersonalEditorProps) {
  const { updateSectionContent } = useEditorStore();

  const updateField = useCallback(
    (field: keyof PersonalContent, value: string) => {
      updateSectionContent(sectionId, {
        ...content,
        [field]: value,
      });
    },
    [sectionId, content, updateSectionContent]
  );

  const addLink = useCallback(() => {
    const newLink: ProfileLink = {
      type: 'website',
      url: '',
      label: '',
    };
    updateSectionContent(sectionId, {
      ...content,
      links: [...(content.links || []), newLink],
    });
  }, [sectionId, content, updateSectionContent]);

  const updateLink = useCallback(
    (index: number, updates: Partial<ProfileLink>) => {
      const links = [...(content.links || [])];
      links[index] = { ...links[index], ...updates };
      updateSectionContent(sectionId, { ...content, links });
    },
    [sectionId, content, updateSectionContent]
  );

  const removeLink = useCallback(
    (index: number) => {
      const links = (content.links || []).filter((_, i) => i !== index);
      updateSectionContent(sectionId, { ...content, links });
    },
    [sectionId, content, updateSectionContent]
  );

  return (
    <div className="space-y-6">
      {/* Name */}
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="firstName">First Name</Label>
          <Input
            id="firstName"
            value={content.firstName || ''}
            onChange={(e) => updateField('firstName', e.target.value)}
            placeholder="John"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="lastName">Last Name</Label>
          <Input
            id="lastName"
            value={content.lastName || ''}
            onChange={(e) => updateField('lastName', e.target.value)}
            placeholder="Doe"
          />
        </div>
      </div>

      {/* Job Title */}
      <div className="space-y-2">
        <Label htmlFor="jobTitle">Job Title</Label>
        <Input
          id="jobTitle"
          value={content.jobTitle || ''}
          onChange={(e) => updateField('jobTitle', e.target.value)}
          placeholder="Senior Software Engineer"
        />
      </div>

      {/* Contact */}
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            value={content.email || ''}
            onChange={(e) => updateField('email', e.target.value)}
            placeholder="john@example.com"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="phone">Phone</Label>
          <Input
            id="phone"
            type="tel"
            value={content.phone || ''}
            onChange={(e) => updateField('phone', e.target.value)}
            placeholder="+1 (555) 123-4567"
          />
        </div>
      </div>

      {/* Location */}
      <div className="space-y-2">
        <Label htmlFor="location">Location</Label>
        <Input
          id="location"
          value={content.location || ''}
          onChange={(e) => updateField('location', e.target.value)}
          placeholder="San Francisco, CA"
        />
      </div>

      {/* Links */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <Label>Links</Label>
          <Button variant="outline" size="sm" onClick={addLink}>
            <Plus className="h-3 w-3 mr-1" />
            Add Link
          </Button>
        </div>

        {(content.links || []).map((link, index) => (
          <div key={index} className="flex items-center gap-2">
            <div className="flex items-center justify-center w-8 h-8 rounded bg-muted">
              {LINK_ICONS[link.type]}
            </div>
            <select
              value={link.type}
              onChange={(e) =>
                updateLink(index, { type: e.target.value as LinkType })
              }
              className="h-10 rounded-md border border-input bg-background px-2 text-sm"
            >
              <option value="linkedin">LinkedIn</option>
              <option value="github">GitHub</option>
              <option value="twitter">Twitter</option>
              <option value="website">Website</option>
              <option value="portfolio">Portfolio</option>
              <option value="other">Other</option>
            </select>
            <Input
              value={link.url}
              onChange={(e) => updateLink(index, { url: e.target.value })}
              placeholder="https://..."
              className="flex-1"
            />
            <Button
              variant="ghost"
              size="icon"
              onClick={() => removeLink(index)}
              className="text-destructive hover:text-destructive"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        ))}

        {(!content.links || content.links.length === 0) && (
          <p className="text-sm text-muted-foreground text-center py-4 border border-dashed rounded-lg">
            No links added yet
          </p>
        )}
      </div>
    </div>
  );
}
