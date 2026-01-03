import { AbstractControl, FormArray, FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators } from "@angular/forms";
import { Item, ItemStatus } from "../../open-api";

export function buildItemForm(item?: Item, receiptId?: string, isShare: boolean = true, syncAmountWithItems: boolean = false): FormGroup {
  const formGroup = new FormGroup({
    name: new FormControl(item?.name ?? "", Validators.required),
    chargedToUserId: new FormControl(item?.chargedToUserId ?? undefined, []),
    receiptId: new FormControl(Number(item?.receiptId ?? receiptId)),
    amount: new FormControl(item?.amount ?? undefined, [
      Validators.required,
      Validators.min(0),
      itemTotalValidator(isShare, syncAmountWithItems),
    ]),
    isTaxed: new FormControl(item?.IsTaxed ?? false),
    categories: new FormArray(item?.categories?.map((c) => new FormControl(c)) ?? []),
    tags: new FormArray(item?.tags?.map((t) => new FormControl(t)) ?? []),
    status: new FormControl(
      item?.status ?? ItemStatus.Open,
      Validators.required
    ),
    // Always include linkedItems FormArray, even if empty
    linkedItems: new FormArray(
      item?.linkedItems?.map(linkedItem => 
        buildItemForm(linkedItem, receiptId, true, syncAmountWithItems)
      ) ?? []
    ),
  });

  if (isShare) {
    formGroup.get("chargedToUserId")?.addValidators(Validators.required);
  }

  return formGroup;
}

function itemTotalValidator(isShare: boolean, syncAmountWithItems: boolean = false): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    // Check current sync state dynamically from the parent form
    const parentForm = control?.parent?.parent?.parent;
    const currentSyncState = parentForm?.get('syncAmountWithItems')?.value ?? syncAmountWithItems;
    
    // Skip validation if sync with items is currently enabled - the receipt total will adjust automatically
    if (currentSyncState) {
      return null;
    }

    const epsilon = 0.01;
    const errKey = "itemLargerThanTotal";
    const objectName = isShare ? "Share" : "Item";

    const formArray = control.parent?.parent as FormArray<FormControl<Item>>;
    const amountControl = control?.parent?.parent?.parent;
    let receiptTotal: number = 1;
    if (amountControl) {
      const amountValue = amountControl.get("amount")?.value;
      if (amountValue !== undefined) {
        receiptTotal = Number.parseFloat(amountValue ?? "1");
      }
    }

    if (!formArray) {
      return null;
    }

    let itemControls = formArray.controls;
    if (isShare) {
      itemControls = itemControls.filter(c => c.value.chargedToUserId !== null && c.value.chargedToUserId !== undefined);
    } else {
      itemControls = itemControls.filter(c => c.value.chargedToUserId === null || c.value.chargedToUserId === undefined);
    }

    const itemsAmounts: number[] = itemControls
      .map((c) => c.get("amount")?.value ?? 0)
      .map((amount: any) => Number.parseFloat(amount) ?? 1);
    const itemsTotal = itemsAmounts.reduce((a, b) => a + b);

    if (itemsTotal > receiptTotal + epsilon) {
      itemControls.forEach((c) => {
        if (c !== control) {
          c.get("amount")?.setErrors({
            [errKey]: `${objectName} sum cannot be larger than receipt total`,
          });
        }
      });
      return { [errKey]: `${objectName} sum cannot be larger than receipt total` };
    } else {
      itemControls.forEach((c) => {
        if (c.get("amount")?.errors && c.get("amount")?.hasError(errKey)) {
          let newErrors = c.get("amount")?.errors ?? {};
          delete newErrors[errKey];
          c.get("amount")?.setErrors(newErrors);
        }
      });
      return null;
    }
  };
}
