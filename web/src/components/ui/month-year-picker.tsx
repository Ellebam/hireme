'use client';

import { useState, useEffect, useId } from 'react';
import { CalendarDays, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Popover, PopoverTrigger, PopoverContent } from './popover';

const MONTHS = [
  'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
];

const currentYear = new Date().getFullYear();
const MIN_YEAR = 1960;
const MAX_YEAR = currentYear + 10;

const YEARS = Array.from(
  { length: MAX_YEAR - MIN_YEAR + 1 },
  (_, i) => MAX_YEAR - i
);

function formatMonthYear(value: string): string {
  if (!value) return '';
  const [year, month] = value.split('-');
  if (!year || !month) return value;
  const monthIndex = parseInt(month, 10) - 1;
  if (monthIndex < 0 || monthIndex > 11) return value;
  return `${MONTHS[monthIndex]} ${year}`;
}

export interface MonthYearPickerProps {
  value?: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  placeholder?: string;
  id?: string;
}

export function MonthYearPicker({
  value,
  onChange,
  disabled = false,
  placeholder = 'Select date',
  id,
}: MonthYearPickerProps) {
  const [open, setOpen] = useState(false);
  const [selectedMonth, setSelectedMonth] = useState<string>('');
  const [selectedYear, setSelectedYear] = useState<string>('');
  const generatedId = useId();
  const monthSelectId = `${generatedId}-month`;
  const yearSelectId = `${generatedId}-year`;

  useEffect(() => {
    if (value) {
      const [year, month] = value.split('-');
      setSelectedYear(year || '');
      setSelectedMonth(month || '');
    } else {
      setSelectedMonth('');
      setSelectedYear('');
    }
  }, [value]);

  const handleMonthChange = (month: string) => {
    setSelectedMonth(month);
    if (month && selectedYear) {
      onChange(`${selectedYear}-${month}`);
      setOpen(false);
    }
  };

  const handleYearChange = (year: string) => {
    setSelectedYear(year);
    if (selectedMonth && year) {
      onChange(`${year}-${selectedMonth}`);
      setOpen(false);
    }
  };

  const handleClear = () => {
    setSelectedMonth('');
    setSelectedYear('');
    onChange('');
    setOpen(false);
  };

  const displayValue = value ? formatMonthYear(value) : '';

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <div className="flex w-full items-center gap-1">
        <PopoverTrigger asChild>
          <button
            id={id}
            type="button"
            disabled={disabled}
            className={cn(
              'flex h-10 flex-1 items-center gap-2 border-2 border-input bg-secondary px-3 py-[9px] text-sm transition-all duration-150 focus-visible:outline-none focus-visible:border-primary focus-visible:bg-card focus-visible:shadow-offset-sm disabled:cursor-not-allowed disabled:opacity-50',
              !displayValue && 'text-muted-foreground'
            )}
          >
            <CalendarDays className="h-4 w-4 text-muted-foreground" />
            {displayValue || placeholder}
          </button>
        </PopoverTrigger>
        {displayValue && !disabled && (
          <button
            type="button"
            aria-label="Clear date"
            onClick={handleClear}
            className="text-muted-foreground hover:text-foreground"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>
      <PopoverContent align="start" className="w-[240px] p-3">
        <div className="space-y-3">
          <div className="space-y-1.5">
            <label htmlFor={monthSelectId} className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Month
            </label>
            <select
              id={monthSelectId}
              value={selectedMonth}
              onChange={(e) => handleMonthChange(e.target.value)}
              className="flex h-9 w-full border-2 border-input bg-secondary px-2 py-1 text-sm focus:outline-none focus:border-primary"
            >
              <option value="">—</option>
              {MONTHS.map((month, index) => (
                <option key={month} value={String(index + 1).padStart(2, '0')}>
                  {month}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-1.5">
            <label htmlFor={yearSelectId} className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              Year
            </label>
            <select
              id={yearSelectId}
              value={selectedYear}
              onChange={(e) => handleYearChange(e.target.value)}
              className="flex h-9 w-full border-2 border-input bg-secondary px-2 py-1 text-sm focus:outline-none focus:border-primary"
            >
              <option value="">—</option>
              {YEARS.map((year) => (
                <option key={year} value={String(year)}>
                  {year}
                </option>
              ))}
            </select>
          </div>
          {value && (
            <button
              type="button"
              onClick={handleClear}
              className="text-xs text-muted-foreground hover:text-foreground underline"
            >
              Clear date
            </button>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
