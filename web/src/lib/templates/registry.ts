/**
 * Template Registry
 *
 * Central registry for all available CV templates.
 * Allows easy lookup and extension.
 */

import type { CVContent } from '@/types/cv';
import { starterTemplate } from './starter';
import { blankTemplate } from './blank';

export interface TemplateMeta {
  id: string;
  name: string;
  description: string;
  factory: () => CVContent;
}

export const AVAILABLE_TEMPLATES: TemplateMeta[] = [
  {
    id: 'starter',
    name: 'Starter Template',
    description: 'A complete example CV with sample content. Great for seeing how your CV could look.',
    factory: starterTemplate,
  },
  {
    id: 'blank',
    name: 'Blank Template',
    description: 'Start from scratch with empty sections. For those who prefer to build from zero.',
    factory: blankTemplate,
  },
];

/**
 * Get a template by ID
 * @param templateId - Template identifier
 * @returns CVContent or null if not found
 */
export function getTemplate(templateId: string): CVContent | null {
  const template = AVAILABLE_TEMPLATES.find((t) => t.id === templateId);
  return template ? template.factory() : null;
}

/**
 * Get the default template for new CVs
 */
export function getDefaultTemplate(): CVContent {
  return starterTemplate();
}
