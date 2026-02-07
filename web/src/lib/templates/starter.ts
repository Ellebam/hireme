/**
 * Starter Template
 *
 * A complete CV template with example content.
 * Shows users how a filled CV looks and demonstrates
 * all editor features.
 */

import type { CVContent } from '@/types/cv';
import { generateId } from '@/lib/utils';

export function starterTemplate(): CVContent {
  return {
    schemaVersion: '1.0.0',
    templateId: 'modern',
    locale: 'en',
    title: 'My Professional CV',
    sections: [
      {
        id: generateId(),
        type: 'personal',
        order: 0,
        visible: true,
        content: {
          firstName: 'Alex',
          lastName: 'Johnson',
          jobTitle: 'Software Engineer',
          email: 'alex.johnson@email.com',
          phone: '+1 (555) 123-4567',
          location: 'San Francisco, CA',
          links: [
            {
              type: 'linkedin',
              url: 'https://linkedin.com/in/alexjohnson',
              label: 'LinkedIn',
            },
            {
              type: 'github',
              url: 'https://github.com/alexjohnson',
              label: 'GitHub',
            },
          ],
        },
      },
      {
        id: generateId(),
        type: 'summary',
        order: 1,
        visible: true,
        content: {
          text: 'Software engineer with 5+ years of experience building web applications. Passionate about clean code, user experience, and continuous learning. Experienced in full-stack development with a focus on React and Node.js.',
        },
      },
      {
        id: generateId(),
        type: 'experience',
        order: 2,
        visible: true,
        title: 'Work Experience',
        content: {
          entries: [
            {
              id: generateId(),
              company: 'Tech Solutions Inc.',
              position: 'Senior Software Engineer',
              location: 'San Francisco, CA',
              startDate: '2022-03',
              endDate: null,
              current: true,
              description:
                'Lead development of customer-facing web applications serving 100k+ users.',
              highlights: [
                'Architected new microservices platform reducing deployment time by 60%',
                'Mentored team of 4 junior developers',
                'Implemented CI/CD pipeline improving release frequency by 3x',
              ],
            },
            {
              id: generateId(),
              company: 'Digital Agency Co.',
              position: 'Full Stack Developer',
              location: 'Los Angeles, CA',
              startDate: '2019-06',
              endDate: '2022-02',
              current: false,
              description:
                'Built custom web solutions for clients across various industries.',
              highlights: [
                'Delivered 15+ client projects on time and within budget',
                'Introduced automated testing reducing bugs by 40%',
              ],
            },
          ],
        },
      },
      {
        id: generateId(),
        type: 'education',
        order: 3,
        visible: true,
        title: 'Education',
        content: {
          entries: [
            {
              id: generateId(),
              institution: 'University of California',
              degree: 'B.S.',
              field: 'Computer Science',
              location: 'Berkeley, CA',
              startDate: '2015-09',
              endDate: '2019-05',
              current: false,
              grade: '3.8 GPA',
              description: 'Graduated with honors. Focus on software engineering and algorithms.',
            },
          ],
        },
      },
      {
        id: generateId(),
        type: 'skills',
        order: 4,
        visible: true,
        title: 'Skills',
        content: {
          categories: [
            {
              id: generateId(),
              name: 'Programming Languages',
              skills: [
                { name: 'TypeScript', level: 'expert' },
                { name: 'JavaScript', level: 'expert' },
                { name: 'Python', level: 'advanced' },
                { name: 'Go', level: 'intermediate' },
              ],
            },
            {
              id: generateId(),
              name: 'Frameworks & Libraries',
              skills: [
                { name: 'React', level: 'expert' },
                { name: 'Next.js', level: 'advanced' },
                { name: 'Node.js', level: 'advanced' },
                { name: 'Express', level: 'advanced' },
              ],
            },
            {
              id: generateId(),
              name: 'Tools & Platforms',
              skills: [
                { name: 'Git', level: 'expert' },
                { name: 'Docker', level: 'advanced' },
                { name: 'AWS', level: 'intermediate' },
                { name: 'PostgreSQL', level: 'advanced' },
              ],
            },
          ],
        },
      },
      {
        id: generateId(),
        type: 'languages',
        order: 5,
        visible: true,
        title: 'Languages',
        content: {
          entries: [
            { id: generateId(), language: 'English', proficiency: 'native' },
            { id: generateId(), language: 'Spanish', proficiency: 'intermediate' },
          ],
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
