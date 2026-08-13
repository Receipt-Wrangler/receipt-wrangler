import { AbstractControl, ValidatorFn, Validators } from "@angular/forms";

// Adds or removes the required validator on a control and revalidates it without emitting a
// valueChanges event. Shared by the quick-scan dialog and group receipt settings, which toggle
// required per the group's quick-scan configuration.
//
// `validator` overrides which validator is toggled, for fields whose "required" differs from the
// stock one (e.g. a trim-aware required for a value the backend trims before its own emptiness
// check). It MUST be a stable reference across calls -- removeValidators matches by identity, so a
// freshly built validator would never be removed. Validators compose additively, so any others on
// the control survive every toggle.
export function setRequired(
  control: AbstractControl | null,
  required: boolean,
  validator: ValidatorFn = Validators.required
): void {
  if (!control) {
    return;
  }

  if (required) {
    control.addValidators(validator);
  } else {
    control.removeValidators(validator);
  }

  control.updateValueAndValidity({ emitEvent: false });
}
