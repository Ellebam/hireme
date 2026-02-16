import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { CertificationsEditor } from '../CertificationsEditor';
import { mockCertificationsContent } from '@/test/mocks/cv';

describe('CertificationsEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders entry cards with name and issuer', () => {
    render(
      <CertificationsEditor sectionId="sec-1" content={mockCertificationsContent} />
    );

    expect(
      screen.getByText('AWS Solutions Architect')
    ).toBeInTheDocument();
    expect(screen.getByText(/Amazon Web Services/)).toBeInTheDocument();
    expect(
      screen.getByText('Google Cloud Professional')
    ).toBeInTheDocument();
    expect(screen.getByText(/Google · 2022-03/)).toBeInTheDocument();
  });

  it('shows entry count', () => {
    render(
      <CertificationsEditor sectionId="sec-1" content={mockCertificationsContent} />
    );

    expect(screen.getByText('2 entries')).toBeInTheDocument();
  });

  it('shows empty state when no entries', () => {
    render(
      <CertificationsEditor sectionId="sec-1" content={{ entries: [] }} />
    );

    expect(
      screen.getByText('No certifications added')
    ).toBeInTheDocument();
  });
});
