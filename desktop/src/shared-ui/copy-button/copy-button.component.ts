import { Component, DestroyRef, inject, input, OnDestroy, output, signal } from "@angular/core";
import { ButtonComponent } from "../../button/button/button.component";

/**
 * Generic copy-to-clipboard icon button.
 *
 * Wraps `<app-button>` with an internal "copied" state that flips the icon
 * and tooltip for a brief window after a successful copy. Intended to be
 * reusable anywhere we want a compact clipboard action — dialogs, table
 * cells, read-only detail views, etc.
 *
 * Usage:
 *   <app-copy-button [text]="someString"></app-copy-button>
 *   <app-copy-button [text]="jsonBlob" tooltip="Copy JSON"></app-copy-button>
 */
@Component({
  selector: "app-copy-button",
  templateUrl: "./copy-button.component.html",
  styleUrl: "./copy-button.component.scss",
  imports: [ButtonComponent],
})
export class CopyButtonComponent implements OnDestroy {
  public readonly text = input.required<string>();

  public readonly tooltip = input<string>("Copy");

  public readonly successTooltip = input<string>("Copied!");

  public readonly icon = input<string>("content_copy");

  public readonly successIcon = input<string>("check");

  public readonly successDurationMs = input<number>(1500);

  public readonly disabled = input<boolean>(false);

  public readonly copied = output<string>();

  protected readonly isCopied = signal(false);

  private resetTimeoutId: ReturnType<typeof setTimeout> | null = null;

  private readonly destroyRef = inject(DestroyRef);

  constructor() {
    this.destroyRef.onDestroy(() => this.clearResetTimer());
  }

  public ngOnDestroy(): void {
    this.clearResetTimer();
  }

  public async copy(): Promise<void> {
    const value = this.text();
    if (!navigator?.clipboard?.writeText) {
      // Non-secure-context or unsupported browser — no-op rather than throw
      // so the button doesn't break the page.
      return;
    }
    try {
      await navigator.clipboard.writeText(value);
    } catch {
      return;
    }
    this.isCopied.set(true);
    this.copied.emit(value);
    this.clearResetTimer();
    this.resetTimeoutId = setTimeout(() => {
      this.isCopied.set(false);
      this.resetTimeoutId = null;
    }, this.successDurationMs());
  }

  private clearResetTimer(): void {
    if (this.resetTimeoutId !== null) {
      clearTimeout(this.resetTimeoutId);
      this.resetTimeoutId = null;
    }
  }
}
