'use client';

import {
  useSensor,
  useSensors,
  PointerSensor,
  KeyboardSensor,
  closestCenter,
} from '@dnd-kit/core';
import {
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';

/**
 * Standard sensors configuration for drag and drop
 * - Pointer sensor with 8px activation distance (prevents accidental drags)
 * - Keyboard sensor with sortable coordinates
 */
export function useDndSensors() {
  return useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );
}

/**
 * Collision detection strategy
 */
export const collisionDetection = closestCenter;

/**
 * Sort strategy for vertical lists
 */
export const sortStrategy = verticalListSortingStrategy;

/**
 * Helper to reorder array items
 */
export function arrayMove<T>(array: T[], from: number, to: number): T[] {
  const newArray = [...array];
  const [removed] = newArray.splice(from, 1);
  newArray.splice(to, 0, removed);
  return newArray;
}
