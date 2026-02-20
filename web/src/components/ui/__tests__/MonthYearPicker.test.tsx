import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, setupUser } from '@/test/utils';
import { MonthYearPicker } from '../month-year-picker';

// Ensure ResizeObserver mock is fresh before each test
// (vi.clearAllMocks in global afterEach can affect the global mock)
beforeEach(() => {
  global.ResizeObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
  }));
});

describe('MonthYearPicker', () => {
  it('renders placeholder when no value', () => {
    render(<MonthYearPicker onChange={vi.fn()} />);
    expect(screen.getByText('Select date')).toBeInTheDocument();
  });

  it('renders custom placeholder', () => {
    render(<MonthYearPicker onChange={vi.fn()} placeholder="Pick a month" />);
    expect(screen.getByText('Pick a month')).toBeInTheDocument();
  });

  it('displays formatted date for existing value', () => {
    render(<MonthYearPicker value="2024-01" onChange={vi.fn()} />);
    expect(screen.getByText('Jan 2024')).toBeInTheDocument();
  });

  it('displays formatted date for different months', () => {
    render(<MonthYearPicker value="2025-12" onChange={vi.fn()} />);
    expect(screen.getByText('Dec 2025')).toBeInTheDocument();
  });

  it('opens popover and shows month/year selects on click', async () => {
    const user = setupUser();
    render(<MonthYearPicker onChange={vi.fn()} />);

    await user.click(screen.getByText('Select date'));

    expect(screen.getByText('Month')).toBeInTheDocument();
    expect(screen.getByText('Year')).toBeInTheDocument();
  });

  it('emits YYYY-MM when both month and year are selected', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<MonthYearPicker onChange={onChange} />);

    await user.click(screen.getByText('Select date'));

    // Select month first (March = 03)
    const monthSelect = screen.getByRole('combobox', { name: /month/i });
    const yearSelect = screen.getByRole('combobox', { name: /year/i });

    await user.selectOptions(monthSelect, '03');
    // Month selected but no year yet — should not emit
    expect(onChange).not.toHaveBeenCalled();

    await user.selectOptions(yearSelect, '2025');
    // Now both are set — should emit
    expect(onChange).toHaveBeenCalledWith('2025-03');
  });

  it('clear resets to empty string', async () => {
    const onChange = vi.fn();
    const user = setupUser();
    render(<MonthYearPicker value="2024-06" onChange={onChange} />);

    // Click the X button to clear
    const clearButton = screen.getByRole('button', { name: /clear date/i });
    await user.click(clearButton);

    expect(onChange).toHaveBeenCalledWith('');
  });

  it('disabled state prevents interaction', () => {
    render(<MonthYearPicker value="2024-01" onChange={vi.fn()} disabled />);

    const trigger = screen.getByRole('button');
    expect(trigger).toBeDisabled();
  });

  it('does not show clear button when disabled', () => {
    render(<MonthYearPicker value="2024-01" onChange={vi.fn()} disabled />);

    expect(screen.queryByLabelText(/clear date/i)).not.toBeInTheDocument();
  });
});
