'use client';

import { useCallback } from 'react';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { useEditorStore } from '@/stores';
import type { SummaryContent } from '@/types/cv';

interface SummaryEditorProps {
  sectionId: string;
  content: SummaryContent;
}

export function SummaryEditor({ sectionId, content }: SummaryEditorProps) {
  const { updateSectionContent } = useEditorStore();

  const handleChange = useCallback(
    (text: string) => {
      updateSectionContent(sectionId, { text });
    },
    [sectionId, updateSectionContent]
  );

  const charCount = (content.text || '').length;
  const recommendedMin = 150;
  const recommendedMax = 300;

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="summary">Professional Summary</Label>
        <Textarea
          id="summary"
          value={content.text || ''}
          onChange={(e) => handleChange(e.target.value)}
          placeholder="Write a brief summary of your professional background, key skills, and career objectives..."
          className="min-h-[150px] resize-y"
        />
      </div>

      <div className="flex items-center justify-between text-sm">
        <p className="text-muted-foreground">
          Recommended: {recommendedMin}-{recommendedMax} characters
        </p>
        <p
          className={
            charCount < recommendedMin || charCount > recommendedMax
              ? 'text-sienna'
              : 'text-green-500'
          }
        >
          {charCount} characters
        </p>
      </div>

      <div className="rounded-lg bg-muted/50 p-4 text-sm">
        <p className="font-medium mb-2">Tips for a great summary:</p>
        <ul className="space-y-1 text-muted-foreground">
          <li>• Start with your job title and years of experience</li>
          <li>• Highlight 2-3 key achievements or skills</li>
          <li>• Mention the value you bring to employers</li>
          <li>• Keep it concise and impactful</li>
        </ul>
      </div>
    </div>
  );
}
