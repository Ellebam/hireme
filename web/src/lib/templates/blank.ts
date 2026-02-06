/**
 * Blank Template
 *
 * Empty CV structure with all common sections
 * but no pre-filled content.
 */

import type { CVContent } from '@/types/cv';
import { generateId } from '@/lib/utils';

export function blankTemplate(): CVContent {
  return {
    schemaVersion: '1.0.0',
    templateId: 'blank',
    locale: 'en',
    title: 'My CV',
    sections: [
      {
        id: generateId(),
        type: 'personal',
        order: 0,
        visible: true,
        content: {
          firstName: '',
          lastName: '',
          jobTitle: '',
          email: '',
          phone: '',
          location: '',
          links: [],
        },
      },
      {
        id: generateId(),
        type: 'summary',
        order: 1,
        visible: true,
        content: {
          text: '',
        },
      },
      {
        id: generateId(),
        type: 'experience',
        order: 2,
        visible: true,
        title: 'Experience',
        content: {
          entries: [],
        },
      },
      {
        id: generateId(),
        type: 'education',
        order: 3,
        visible: true,
        title: 'Education',
        content: {
          entries: [],
        },
      },
      {
        id: generateId(),
        type: 'skills',
        order: 4,
        visible: true,
        title: 'Skills',
        content: {
          categories: [],
        },
      },
      {
        id: generateId(),
        type: 'languages',
        order: 5,
        visible: true,
        title: 'Languages',
        content: {
          entries: [],
        },
      },
    ],
    styling: {
      primaryColor: '#2563eb',
      secondaryColor: '#64748b',
      fontFamily: 'inter',
      fontSize: 'medium',
      lineHeight: 'normal',
      showIcons: true,
    },
  };
}
