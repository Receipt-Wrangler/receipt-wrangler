import { AbstractControl, ValidationErrors, ValidatorFn } from "@angular/forms";

/**
 * Bounds a duration setting edited as a number + unit pair: a whole number, at
 * least `min`, at most `max` (both expressed in the currently selected unit).
 *
 * `min` defaults to 1 because most of these settings only need a positive whole
 * number. Pass it when the API enforces a real floor — without it a value the
 * server rejects passes client validation and comes back as a bare 400 with
 * nothing attached to the field.
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
export function durationValueValidator(max: number, min: number = 1): ValidatorFn {
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

    if (value < min) {
      return { duration: `Must be at least ${min}.` };
    }

    if (value > max) {
      return { duration: `Must be at most ${max}.` };
    }

    return null;
  };
}
