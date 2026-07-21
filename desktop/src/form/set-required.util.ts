import { AbstractControl, Validators } from "@angular/forms";

// Adds or removes the required validator on a control and revalidates it without emitting a
// valueChanges event. Shared by the quick-scan dialog and group receipt settings, which toggle
// required per the group's quick-scan configuration.
export function setRequired(control: AbstractControl | null, required: boolean): void {
  if (!control) {
    return;
  }

  if (required) {
    control.addValidators(Validators.required);
  } else {
    control.removeValidators(Validators.required);
  }

  control.updateValueAndValidity({ emitEvent: false });
}
