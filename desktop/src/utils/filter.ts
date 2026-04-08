import { AbstractControl, FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { untilDestroyed } from "@ngneat/until-destroy";
import { distinctUntilChanged, pairwise, startWith, tap } from "rxjs";
import { FilterOperation } from "../open-api/index";

export function listenForBetweenOperation(form: FormGroup, key: string, thisContext: any): void {
  updateValidatorsOnOperationChange(true, form, key, null, form.get(key)?.get("operation")?.value);

  form
    .get(key)
    ?.get("operation")
    ?.valueChanges
    .pipe(
      startWith(form.get(key)?.get("operation")?.value),
      distinctUntilChanged(),
      pairwise(),
      untilDestroyed(thisContext),
      tap(([prev, curr]: [FilterOperation | null, FilterOperation | null]) => {
        updateValidatorsOnOperationChange(false, form, key, prev, curr);
      }),
    ).subscribe();
}

function updateValidatorsOnOperationChange(
  firstRun: boolean,
  form: FormGroup,
  key: string,
  prev: FilterOperation | null,
  curr: FilterOperation | null
)
  : void {
  const formBuilder = new FormBuilder();
  if (curr === FilterOperation.Between) {
    if (!firstRun) {
      (form.get(key) as FormGroup).removeControl("value");
      (form.get(key) as FormGroup).addControl("value", formBuilder.array([null, null]));
    }

    (form.get(key) as FormGroup).get("value.0")?.addValidators(Validators.required);
    (form.get(key) as FormGroup).get("value.1")?.addValidators(Validators.required);

    (form.get(key) as FormGroup).get("value")?.addValidators(betweenValidator);
  } else if (prev === FilterOperation.Between && curr !== FilterOperation.Between) {
    if (!firstRun) {
      (form.get(key) as FormGroup).removeControl("value");
      (form.get(key) as FormGroup).addControl("value", formBuilder.control(null));
    }
  }
}

export function betweenValidator(control: AbstractControl): { [key: string]: any } | null {
  const formArray = control as FormArray;

  if (formArray.value[0] > formArray.value[1] && formArray.value[1] !== null) {
    formArray.at(0).setErrors({ invalidValue: true });
  }

  if (formArray.value[0] < formArray.value[1]) {
    formArray.at(0).setErrors(null);
  }

  return null;
}

export function buildFieldFormGroup(
  value: string | string[] | number | any,
  operation: string | undefined,
  thisContext: any,
  isArray?: boolean,
): FormGroup {
  const formBuilder = new FormBuilder();
  let valueControl: AbstractControl;
  const operationControl = formBuilder.control(operation);

  if (operation === FilterOperation.Between) {
    valueControl = formBuilder.array([value?.[0], value?.[1]]);
  } else if (isArray) {
    valueControl = formBuilder.array(value);
  } else {
    valueControl = formBuilder.control(value);
  }

  valueControl.valueChanges.pipe(
    untilDestroyed(thisContext),
    startWith(value),
    tap((value) => {
      if (!!value && value?.length > 0) {
        operationControl.addValidators(Validators.required);
      } else {
        operationControl.removeValidators(Validators.required);
      }

      operationControl.updateValueAndValidity();
    })).subscribe();

  return formBuilder.group({
    operation: operationControl,
    value: valueControl
  });
}
