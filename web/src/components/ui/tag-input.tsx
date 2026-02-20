'use client';

import { useState, useCallback, type KeyboardEvent, type ClipboardEvent } from 'react';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';

interface TagInputProps {
  value: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
  id?: string;
}

export function TagInput({
  value,
  onChange,
  placeholder = 'Type and press Enter',
  id,
}: TagInputProps) {
  const [inputValue, setInputValue] = useState('');

  const addTags = useCallback(
    (raw: string) => {
      const newTags = raw
        .split(/[,\n]/)
        .map((tag) => tag.trim())
        .filter((tag) => tag !== '' && !value.includes(tag));

      if (newTags.length > 0) {
        onChange([...value, ...newTags]);
      }
      setInputValue('');
    },
    [value, onChange]
  );

  const removeTag = useCallback(
    (index: number) => {
      onChange(value.filter((_, i) => i !== index));
    },
    [value, onChange]
  );

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (inputValue.trim()) {
        addTags(inputValue);
      }
    } else if (e.key === 'Backspace' && inputValue === '' && value.length > 0) {
      removeTag(value.length - 1);
    }
  };

  const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
    const pasted = e.clipboardData.getData('text');
    if (pasted.includes(',') || pasted.includes('\n')) {
      e.preventDefault();
      addTags(pasted);
    }
  };

  return (
    <div
      className={cn(
        'flex flex-wrap gap-1.5 border-2 border-input bg-secondary p-2 transition-all duration-150 focus-within:border-primary focus-within:bg-card focus-within:shadow-offset-sm'
      )}
    >
      {value.map((tag, index) => (
        <span
          key={`${tag}-${index}`}
          className="inline-flex items-center gap-1 bg-primary/10 text-primary px-2 py-0.5 text-sm"
        >
          {tag}
          <button
            type="button"
            onClick={() => removeTag(index)}
            className="text-primary/60 hover:text-primary"
            aria-label={`Remove ${tag}`}
          >
            <X className="h-3 w-3" />
          </button>
        </span>
      ))}
      <input
        id={id}
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        onPaste={handlePaste}
        placeholder={value.length === 0 ? placeholder : ''}
        className="flex-1 min-w-[120px] bg-transparent text-sm outline-none placeholder:text-muted-foreground"
      />
    </div>
  );
}
