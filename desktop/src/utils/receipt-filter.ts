import { AbstractControl, FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { untilDestroyed } from "@ngneat/until-destroy";
import { distinctUntilChanged, pairwise, startWith, tap } from "rxjs";
import { FilterOperation } from "../open-api/index";

export function buildReceiptFilterForm(filter: any, thisContext: any): FormGroup {
  const formGroup = new FormGroup({
    date: buildFieldFormGroup(
      filter?.date?.value,
      filter?.date?.operation,
      thisContext,
      filter?.date?.operation === FilterOperation.Between
    ),
    amount: buildFieldFormGroup(
      filter?.amount?.value,
      filter?.amount?.operation,
      thisContext,
      filter?.amount?.operation === FilterOperation.Between
    ),
    name: buildFieldFormGroup(
      filter?.name?.value,
      filter?.name?.operation,
      thisContext
    ),
    paidBy: buildFieldFormGroup(
      filter?.paidBy?.value ?? [],
      filter?.paidBy?.operation,
      thisContext,
      true
    ),
    categories: buildFieldFormGroup(
      filter?.categories?.value ?? [],
      filter?.categories?.operation,
      thisContext,
      true
    ),
    tags: buildFieldFormGroup(
      filter?.tags?.value ?? [],
      filter?.tags?.operation,
      thisContext,
      true
    ),
    status: buildFieldFormGroup(
      filter?.status?.value ?? [],
      filter?.status?.operation,
      thisContext,
      true
    ),
    resolvedDate: buildFieldFormGroup(
      filter?.resolvedDate?.value,
      filter?.resolvedDate?.operation,
      thisContext,
      filter?.resolvedDate?.operation === FilterOperation.Between
    ),
    createdAt: buildFieldFormGroup(
      filter?.createdAt?.value,
      filter?.createdAt?.operation,
      thisContext,
      filter?.createdAt?.operation === FilterOperation.Between
    ),
  });

  listenForBetweenOperation(formGroup, "amount", thisContext);
  listenForBetweenOperation(formGroup, "date", thisContext);
  listenForBetweenOperation(formGroup, "resolvedDate", thisContext);
  listenForBetweenOperation(formGroup, "createdAt", thisContext);


  return formGroup;
}

function listenForBetweenOperation(form: FormGroup, key: string, thisContext: any): void {
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

function betweenValidator(control: AbstractControl): { [key: string]: any } | null {
  const formArray = control as FormArray;

  if (formArray.value[0] > formArray.value[1] && formArray.value[1] !== null) {
    formArray.at(0).setErrors({ invalidValue: true });
  }

  if (formArray.value[0] < formArray.value[1]) {
    formArray.at(0).setErrors(null);
  }

  return null;
}

function buildFieldFormGroup(
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
