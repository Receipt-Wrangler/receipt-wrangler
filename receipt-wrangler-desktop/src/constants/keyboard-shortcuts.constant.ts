import { KeyboardShortcut } from '../services/keyboard-shortcut.service';

/**
 * Keyboard shortcuts for item management
 */
export const ITEM_MANAGEMENT_SHORTCUTS: KeyboardShortcut[] = [
  {
    key: 'i',
    ctrlKey: true,
    action: 'ADD_ITEM',
    description: 'Add new item'
  },
  {
    key: 'Enter',
    ctrlKey: true,
    action: 'SUBMIT_AND_CONTINUE',
    description: 'Add item and continue'
  },
  {
    key: 'Enter',
    ctrlKey: true,
    shiftKey: true,
    action: 'SUBMIT_AND_FINISH',
    description: 'Add item and finish'
  },
  {
    key: 'Escape',
    action: 'CANCEL',
    description: 'Cancel current action'
  }
];

/**
 * Actions for keyboard shortcuts
 */
export const KEYBOARD_SHORTCUT_ACTIONS = {
  ADD_ITEM: 'ADD_ITEM',
  SUBMIT_AND_CONTINUE: 'SUBMIT_AND_CONTINUE',
  SUBMIT_AND_FINISH: 'SUBMIT_AND_FINISH',
  CANCEL: 'CANCEL'
} as const;

/**
 * Display shortcuts for UI components
 */
export const DISPLAY_SHORTCUTS = [
  {
    keys: 'Ctrl+Enter',
    action: 'Add Item',
    description: 'Add the current item and continue adding more'
  },
  {
    keys: 'Ctrl+Shift+Enter',
    action: 'Add & Done',
    description: 'Add the current item and close the form'
  },
  {
    keys: 'Escape',
    action: 'Cancel',
    description: 'Cancel the current operation'
  },
  {
    keys: 'Tab',
    action: 'Next Field',
    description: 'Navigate to the next form field'
  }
] as const;

export type KeyboardShortcutAction = typeof KEYBOARD_SHORTCUT_ACTIONS[keyof typeof KEYBOARD_SHORTCUT_ACTIONS];