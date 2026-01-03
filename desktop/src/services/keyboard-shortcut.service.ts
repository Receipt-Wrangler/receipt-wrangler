import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { ITEM_MANAGEMENT_SHORTCUTS } from '../constants/keyboard-shortcuts.constant';

export interface KeyboardShortcut {
  key: string;
  ctrlKey?: boolean;
  shiftKey?: boolean;
  altKey?: boolean;
  metaKey?: boolean;
  action: string;
  description?: string;
}

export interface KeyboardShortcutEvent {
  action: string;
  event: KeyboardEvent;
}

@Injectable({
  providedIn: 'root'
})
export class KeyboardShortcutService {
  private shortcuts: Map<string, KeyboardShortcut> = new Map();
  private shortcutTriggered$ = new Subject<KeyboardShortcutEvent>();
  private keyboardHintTimeout: any;
  private showHint$ = new Subject<boolean>();

  public shortcutTriggered = this.shortcutTriggered$.asObservable();
  public showHint = this.showHint$.asObservable();

  constructor() {
    this.initializeDefaultShortcuts();
  }

  private initializeDefaultShortcuts(): void {
    // Register item management shortcuts from constants
    ITEM_MANAGEMENT_SHORTCUTS.forEach(shortcut => {
      this.registerShortcut(shortcut);
    });
  }

  public registerShortcut(shortcut: KeyboardShortcut): void {
    const key = this.getShortcutKey(shortcut);
    this.shortcuts.set(key, shortcut);
  }

  public unregisterShortcut(shortcut: KeyboardShortcut): void {
    const key = this.getShortcutKey(shortcut);
    this.shortcuts.delete(key);
  }

  public handleKeyboardEvent(event: KeyboardEvent): boolean {
    const key = this.getEventKey(event);
    const shortcut = this.shortcuts.get(key);

    if (shortcut) {
      event.preventDefault();
      this.shortcutTriggered$.next({
        action: shortcut.action,
        event: event
      });

      // Show hint for certain shortcuts
      if (this.shouldShowHint(shortcut)) {
        this.triggerHint();
      }

      return true;
    }

    return false;
  }

  public getShortcuts(): KeyboardShortcut[] {
    return Array.from(this.shortcuts.values());
  }

  public getShortcutsByAction(action: string): KeyboardShortcut[] {
    return Array.from(this.shortcuts.values()).filter(s => s.action === action);
  }

  private getShortcutKey(shortcut: KeyboardShortcut): string {
    const parts: string[] = [];
    if (shortcut.ctrlKey) parts.push('ctrl');
    if (shortcut.shiftKey) parts.push('shift');
    if (shortcut.altKey) parts.push('alt');
    if (shortcut.metaKey) parts.push('meta');
    parts.push(shortcut.key.toLowerCase());
    return parts.join('+');
  }

  private getEventKey(event: KeyboardEvent): string {
    const parts: string[] = [];
    if (event.ctrlKey) parts.push('ctrl');
    if (event.shiftKey) parts.push('shift');
    if (event.altKey) parts.push('alt');
    if (event.metaKey) parts.push('meta');
    parts.push(event.key.toLowerCase());
    return parts.join('+');
  }

  private shouldShowHint(shortcut: KeyboardShortcut): boolean {
    // Show hint for item-related shortcuts
    return ['ADD_ITEM', 'SUBMIT_AND_CONTINUE', 'SUBMIT_AND_FINISH'].includes(shortcut.action);
  }

  private triggerHint(): void {
    this.showHint$.next(true);
    clearTimeout(this.keyboardHintTimeout);
    this.keyboardHintTimeout = setTimeout(() => {
      this.showHint$.next(false);
    }, 2000);
  }

  public clearHintTimeout(): void {
    clearTimeout(this.keyboardHintTimeout);
    this.showHint$.next(false);
  }
}