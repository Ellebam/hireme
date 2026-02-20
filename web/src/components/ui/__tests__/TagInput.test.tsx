import { describe, it, expect, vi } from 'vitest';
import { render, screen, setupUser } from '@/test/utils';
import { TagInput } from '../tag-input';

describe('TagInput', () => {
  it('renders existing tags as pills', () => {
    render(<TagInput value={['React', 'Node.js']} onChange={vi.fn()} />);

    expect(screen.getByText('React')).toBeInTheDocument();
    expect(screen.getByText('Node.js')).toBeInTheDocument();
  });

  it('shows placeholder when empty', () => {
    render(<TagInput value={[]} onChange={vi.fn()} />);

    expect(
      screen.getByPlaceholderText('Type and press Enter')
    ).toBeInTheDocument();
  });

  it('hides placeholder when tags exist', () => {
    render(<TagInput value={['React']} onChange={vi.fn()} />);

    expect(
      screen.queryByPlaceholderText('Type and press Enter')
    ).not.toBeInTheDocument();
  });

  it('adds tag on Enter', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={['React', 'Node.js']} onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.type(input, 'PostgreSQL{Enter}');

    expect(onChange).toHaveBeenCalledWith(['React', 'Node.js', 'PostgreSQL']);
  });

  it('removes tag on X click', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={['React', 'Node.js']} onChange={onChange} />);

    await user.click(screen.getByLabelText('Remove React'));

    expect(onChange).toHaveBeenCalledWith(['Node.js']);
  });

  it('prevents duplicate tags', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={['React']} onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.type(input, 'React{Enter}');

    expect(onChange).not.toHaveBeenCalled();
  });

  it('ignores empty input on Enter', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={['React']} onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.click(input);
    await user.keyboard('{Enter}');

    expect(onChange).not.toHaveBeenCalled();
  });

  it('splits pasted comma-separated text into tags', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={[]} onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.click(input);
    await user.paste('React, Vue, Angular');

    expect(onChange).toHaveBeenCalledWith(['React', 'Vue', 'Angular']);
  });

  it('splits pasted newline-separated text into tags', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={[]} onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.click(input);
    await user.paste('React\nVue\nAngular');

    expect(onChange).toHaveBeenCalledWith(['React', 'Vue', 'Angular']);
  });

  it('deduplicates pasted tags against existing', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={['React']} onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.click(input);
    await user.paste('React, Vue');

    expect(onChange).toHaveBeenCalledWith(['React', 'Vue']);
  });

  it('backspace removes last tag when input is empty', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<TagInput value={['React', 'Vue']} onChange={onChange} />);

    const input = screen.getByRole('textbox');
    await user.click(input);
    await user.keyboard('{Backspace}');

    expect(onChange).toHaveBeenCalledWith(['React']);
  });

  it('renders custom placeholder', () => {
    render(
      <TagInput value={[]} onChange={vi.fn()} placeholder="Add skills" />
    );

    expect(screen.getByPlaceholderText('Add skills')).toBeInTheDocument();
  });
});
