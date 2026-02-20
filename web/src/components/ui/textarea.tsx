import * as React from 'react';
import { cn } from '@/lib/utils';

export interface TextareaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {}

const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, ...props }, ref) => {
    return (
      <textarea
        className={cn(
          'flex min-h-[80px] w-full resize-y border-2 border-input bg-secondary px-3 py-[9px] text-sm text-foreground transition-all duration-150 placeholder:text-muted-foreground focus-visible:outline-none focus-visible:border-primary focus-visible:bg-card focus-visible:shadow-offset-sm disabled:cursor-not-allowed disabled:opacity-50',
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Textarea.displayName = 'Textarea';

export { Textarea };
