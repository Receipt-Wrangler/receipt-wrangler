import { AbstractControl, ValidationErrors, ValidatorFn } from "@angular/forms";

/**
 * A maximum-length validator that counts **Unicode code points**, not UTF-16
 * code units like `Validators.maxLength`.
 *
 * The Go backend measures its string limits in runes (e.g. the quick-scan
 * comment's `models.MaxCommentLength`, checked with `utf8.RuneCountInString`),
 * and a `varchar(n)` column counts characters on MySQL/Postgres. JavaScript's
 * `String.length` counts UTF-16 code units, so anything outside the BMP counts
 * twice: `"😀".repeat(500)` has a `length` of 1000 and would be rejected here
 * while the server accepts it happily. Counting code points keeps the client
 * from refusing input the API would take.
 *
 * Emits Angular's own `maxlength` error shape so existing error-message
 * mappings and `hasError("maxlength")` checks keep working unchanged.
 */
export function codePointMaxLengthValidator(maxLength: number): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const raw = control.value;
    if (raw === null || raw === undefined || raw === "") {
      return null;
    }

    // Array.from iterates by code point, so a surrogate pair counts once --
    // the same unit the backend's rune count uses.
    const actualLength = Array.from(String(raw)).length;
    if (actualLength <= maxLength) {
      return null;
    }

    return { maxlength: { requiredLength: maxLength, actualLength } };
  };
}

/**
 * `Validators.required` but whitespace-only is treated as blank.
 *
 * Use it wherever the backend trims a string before its own emptiness check --
 * otherwise `"   "` satisfies the form, and the request comes back as a 400
 * saying the field is missing. The quick-scan comment is one such field
 * (`QuickScanCommand` trims each entry on parse, so a whitespace-only comment
 * arrives empty and fails the group's required check).
 *
 * Emits the standard `required` error so the shared inputs render their usual
 * "<label> is required." message.
 */
export function trimmedRequiredValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const raw = control.value;
    if (raw === null || raw === undefined) {
      return { required: true };
    }

    return String(raw).trim().length === 0 ? { required: true } : null;
  };
}
