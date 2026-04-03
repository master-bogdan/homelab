import { renderWithProviders, screen } from '@/test/test-utils';

import { SectionCard } from '../SectionCard';

describe('SectionCard', () => {
  it('renders the title, description, and content', () => {
    renderWithProviders(
      <SectionCard description="Module-ready card wrapper" title="Section title">
        <div>Section content</div>
      </SectionCard>
    );

    expect(screen.getByRole('heading', { name: 'Section title' })).toBeInTheDocument();
    expect(screen.getByText('Module-ready card wrapper')).toBeInTheDocument();
    expect(screen.getByText('Section content')).toBeInTheDocument();
  });
});
