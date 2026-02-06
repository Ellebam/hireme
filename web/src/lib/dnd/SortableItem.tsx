'use client';

import React from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { cn } from '@/lib/utils';

interface SortableItemProps {
  id: string;
  children: React.ReactNode;
  className?: string;
  disabled?: boolean;
}

export function SortableItem({
  id,
  children,
  className,
  disabled = false,
}: SortableItemProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id, disabled });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        'relative',
        isDragging && 'opacity-50 z-50',
        className
      )}
      {...attributes}
      {...listeners}
    >
      {children}
    </div>
  );
}

/**
 * Hook to get sortable props for custom implementations
 */
export function useSortableItem(id: string, disabled = false) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
    isOver,
  } = useSortable({ id, disabled });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return {
    ref: setNodeRef,
    style,
    isDragging,
    isOver,
    dragHandleProps: {
      ...attributes,
      ...listeners,
    },
  };
}
