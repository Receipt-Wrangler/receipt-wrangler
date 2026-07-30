import { Injectable } from "@angular/core";
import { generateSecurePassword } from "../utils/password.utils";
import { SnackbarService } from "./snackbar.service";

export const PASSWORD_COPIED_MESSAGE =
  "Password generated and copied to clipboard";

export const PASSWORD_COPY_FAILED_MESSAGE =
  "Password generated, but copying to the clipboard failed";

/**
 * Generates a password and copies it to the clipboard, toasting the outcome.
 *
 * Keeps the clipboard/snackbar policy in one place so every "Generate
 * Password" control behaves identically, and so the generic `app-input` that
 * hosts the control stays free of clipboard handling.
 */
@Injectable({
  providedIn: "root",
})
export class PasswordGeneratorService {
  constructor(private snackbarService: SnackbarService) {}

  /**
   * Returns a new password **synchronously** — callers run under zoneless
   * change detection, where state set after an `await` would not re-render.
   * The clipboard write and its toast run as a detached side effect.
   */
  public generateAndCopy(length?: number): string {
    const password = generateSecurePassword(length);
    void this.copyAndNotify(password);

    return password;
  }

  public async copyAndNotify(password: string): Promise<void> {
    // Non-secure context or unsupported browser: the password is still in the
    // field, so tell the user to copy it manually rather than failing silently.
    if (!navigator?.clipboard?.writeText) {
      this.snackbarService.error(PASSWORD_COPY_FAILED_MESSAGE);
      return;
    }

    try {
      await navigator.clipboard.writeText(password);
    } catch {
      this.snackbarService.error(PASSWORD_COPY_FAILED_MESSAGE);
      return;
    }

    this.snackbarService.success(PASSWORD_COPIED_MESSAGE);
  }
}
