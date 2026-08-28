import { AbstractControl, ValidationErrors, ValidatorFn } from "@angular/forms";

/**
 * Bounds a duration setting edited as a number + unit pair: a whole number, at
 * least 1, at most `max` (expressed in the currently selected unit).
 *
 * It replaces `Validators.min` / `Validators.max` here for two reasons:
 *
 * - **Whole numbers.** The API stores these as a Go `int`. A fractional entry
 *   like `1.5` fails `json.Unmarshal` outright, so the request comes back as an
 *   unparseable-body error rather than a field-level message. `type="number"`
 *   accepts decimals, so the form has to reject them.
 * - **Visible messages.** `BaseInputComponent` only maps `required` / `email` /
 *   `duplicate` / `min` to text, so a `Validators.max` failure renders as an
 *   empty `mat-error` — a red field with no explanation. Emitting the message as
 *   the error *value* takes the `typeof value === "string"` path in that
 *   component, which renders it verbatim.
 */
export function durationValueValidator(max: number): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const raw = control.value;

    // Emptiness is Validators.required's job, not ours.
    if (raw === null || raw === undefined || raw === "") {
      return null;
    }

    const value = Number(raw);
    if (!Number.isFinite(value)) {
      return { duration: "Must be a whole number." };
    }

    if (!Number.isInteger(value)) {
      return { duration: "Must be a whole number." };
    }

    if (value < 1) {
      return { duration: "Must be at least 1." };
    }

    if (value > max) {
      return { duration: `Must be at most ${max}.` };
    }

    return null;
  };
}
