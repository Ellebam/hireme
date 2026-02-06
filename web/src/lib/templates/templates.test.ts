/**
 * Template Tests
 */

import { describe, it, expect } from 'vitest';
import { starterTemplate } from './starter';
import { blankTemplate } from './blank';
import { getTemplate, AVAILABLE_TEMPLATES } from './registry';
import { SCHEMA_VERSION } from '@/types/cv';

describe('starterTemplate', () => {
  it('returns valid CV content structure', () => {
    const template = starterTemplate();

    expect(template.schemaVersion).toBe(SCHEMA_VERSION);
    expect(template.templateId).toBe('modern');
    expect(template.locale).toBe('en');
    expect(template.sections).toBeInstanceOf(Array);
    expect(template.sections.length).toBeGreaterThan(0);
  });

  it('includes all required section types', () => {
    const template = starterTemplate();
    const sectionTypes = template.sections.map((s) => s.type);

    expect(sectionTypes).toContain('personal');
    expect(sectionTypes).toContain('summary');
    expect(sectionTypes).toContain('experience');
    expect(sectionTypes).toContain('education');
    expect(sectionTypes).toContain('skills');
    expect(sectionTypes).toContain('languages');
  });

  it('has populated content in personal section', () => {
    const template = starterTemplate();
    const personal = template.sections.find((s) => s.type === 'personal');

    expect(personal).toBeDefined();
    expect(personal?.content).toHaveProperty('firstName');
    expect(personal?.content).toHaveProperty('lastName');
    expect((personal?.content as { firstName: string }).firstName).toBeTruthy();
  });

  it('has experience entries', () => {
    const template = starterTemplate();
    const experience = template.sections.find((s) => s.type === 'experience');

    expect(experience).toBeDefined();
    const content = experience?.content as { entries: unknown[] };
    expect(content.entries.length).toBeGreaterThan(0);
  });

  it('generates unique IDs each call', () => {
    const template1 = starterTemplate();
    const template2 = starterTemplate();

    const id1 = template1.sections[0].id;
    const id2 = template2.sections[0].id;

    expect(id1).not.toBe(id2);
  });
});

describe('blankTemplate', () => {
  it('returns valid CV content structure', () => {
    const template = blankTemplate();

    expect(template.schemaVersion).toBe(SCHEMA_VERSION);
    expect(template.templateId).toBe('blank');
    expect(template.sections).toBeInstanceOf(Array);
  });

  it('includes all required section types', () => {
    const template = blankTemplate();
    const sectionTypes = template.sections.map((s) => s.type);

    expect(sectionTypes).toContain('personal');
    expect(sectionTypes).toContain('summary');
    expect(sectionTypes).toContain('experience');
    expect(sectionTypes).toContain('education');
    expect(sectionTypes).toContain('skills');
    expect(sectionTypes).toContain('languages');
  });

  it('has empty content in sections', () => {
    const template = blankTemplate();
    const personal = template.sections.find((s) => s.type === 'personal');

    expect(personal).toBeDefined();
    const content = personal?.content as { firstName: string };
    expect(content.firstName).toBe('');
  });

  it('has empty entries arrays', () => {
    const template = blankTemplate();
    const experience = template.sections.find((s) => s.type === 'experience');

    expect(experience).toBeDefined();
    const content = experience?.content as { entries: unknown[] };
    expect(content.entries).toEqual([]);
  });
});

describe('template registry', () => {
  it('has available templates', () => {
    expect(AVAILABLE_TEMPLATES.length).toBeGreaterThan(0);
  });

  it('each template has required metadata', () => {
    for (const meta of AVAILABLE_TEMPLATES) {
      expect(meta.id).toBeTruthy();
      expect(meta.name).toBeTruthy();
      expect(meta.description).toBeTruthy();
      expect(typeof meta.factory).toBe('function');
    }
  });

  it('getTemplate returns correct template', () => {
    const starter = getTemplate('starter');
    expect(starter).toBeDefined();
    expect(starter?.schemaVersion).toBe(SCHEMA_VERSION);

    const blank = getTemplate('blank');
    expect(blank).toBeDefined();
    expect(blank?.schemaVersion).toBe(SCHEMA_VERSION);
  });

  it('getTemplate returns null for unknown template', () => {
    const result = getTemplate('nonexistent');
    expect(result).toBeNull();
  });
});
