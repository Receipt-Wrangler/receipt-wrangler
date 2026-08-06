import { AbstractControl, ValidationErrors, ValidatorFn } from "@angular/forms";

/**
 * Rejects anything that is not an absolute http(s) URL, mirroring the backend's
 * `isValidAbsoluteUrl` (`api/internal/commands/upsert_system_settings_command.go`)
 * so the user gets inline feedback instead of a 400 snackbar on save.
 *
 * A blank control is left to `Validators.required` — the backend likewise only
 * validates the format once a value is present. Whitespace-only is treated as
 * invalid rather than blank: the backend trims before its own emptiness check,
 * so spaces would otherwise satisfy `Validators.required` here and still be
 * rejected server-side as missing.
 *
 * Embedded credentials (`https://user:token@host`) are rejected for the same
 * reason the backend rejects them: both URLs this guards are published to
 * unauthenticated clients.
 */
export function absoluteUrlValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const raw = control.value;
    if (raw === null || raw === undefined || raw === "") {
      return null;
    }

    const trimmed = String(raw).trim();
    if (trimmed.length === 0) {
      return { url: true };
    }

    // `new URL()` is more lenient than Go's `url.Parse`: it normalizes
    // authority-less forms (`https:host/api`, `https:/host/api`, and
    // backslash variants) into a valid URL, whereas `url.Parse` leaves Host
    // empty and the backend rejects them. Requiring the literal prefix keeps
    // the two in agreement. Case-insensitive on purpose -- `url.Parse`
    // lowercases the scheme, so `HTTPS://host` is valid server-side.
    if (!/^https?:\/\//i.test(trimmed)) {
      return { url: true };
    }

    let parsed: URL;
    try {
      parsed = new URL(trimmed);
    } catch {
      return { url: true };
    }

    const isHttp = parsed.protocol === "http:" || parsed.protocol === "https:";
    if (!isHttp || !parsed.hostname || parsed.username || parsed.password) {
      return { url: true };
    }

    return null;
  };
}
