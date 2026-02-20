import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const buttonVariants = cva(
  'inline-flex items-center justify-center whitespace-nowrap text-[0.8125rem] font-semibold uppercase tracking-[0.0625em] ring-offset-background transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        default:
          'bg-primary text-primary-foreground shadow-offset-sm hover:-translate-x-px hover:-translate-y-px hover:shadow-offset-sm-hover active:translate-x-0 active:translate-y-0 active:shadow-none',
        destructive:
          'bg-destructive text-destructive-foreground shadow-offset-sm hover:-translate-x-px hover:-translate-y-px hover:shadow-offset-sm-hover active:translate-x-0 active:translate-y-0 active:shadow-none',
        accent:
          'bg-accent text-accent-foreground shadow-offset-sm hover:-translate-x-px hover:-translate-y-px hover:shadow-offset-sm-hover active:translate-x-0 active:translate-y-0 active:shadow-none',
        outline:
          'bg-transparent text-primary border-2 border-primary hover:bg-primary hover:text-primary-foreground',
        secondary:
          'bg-secondary text-secondary-foreground hover:bg-secondary/80',
        ghost:
          'bg-transparent text-[hsl(var(--text-secondary))] hover:text-primary',
        link: 'text-primary underline-offset-4 hover:underline',
      },
      size: {
        default: 'px-5 py-2.5',
        sm: 'px-3.5 py-1.5 text-[0.6875rem]',
        lg: 'px-8 py-3',
        icon: 'h-10 w-10',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, type, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button';
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        type={asChild ? type : type ?? 'button'}
        {...props}
      />
    );
  }
);
Button.displayName = 'Button';

export { Button, buttonVariants };
