import { ComponentRef, DestroyRef, InputSignalWithTransform } from "@angular/core";
import { MatDialogRef } from "@angular/material/dialog";

/** A component's signal inputs, keyed by name, typed by the value each accepts. */
export type DialogInputs<T> = {
  [K in keyof T as T[K] extends InputSignalWithTransform<any, any> ? K : never]?:
    T[K] extends InputSignalWithTransform<any, infer WriteT> ? WriteT : never;
};

/**
 * Sets the signal inputs of a component opened with MatDialog. The dialog
 * creates the component without template bindings, and a signal input cannot
 * be assigned, so openers pass its inputs through here. Call it straight after
 * `open()`: the dialog's first change detection runs later, so the values are
 * in place before the component first reads them, required inputs included.
 *
 * Keys and values are checked against the component's input() declarations at
 * each call site, by the app build. Every key is optional, though: Angular's
 * types do not mark a required input apart from a defaulted one, so an opener
 * must pass a dialog's required inputs itself.
 */
export function setDialogInputs<T>(
  dialogRef: MatDialogRef<T>,
  inputs: DialogInputs<T>
): void {
  const componentRef = getDialogComponentRef(dialogRef, "setDialogInputs");

  for (const [name, value] of Object.entries(inputs)) {
    componentRef.setInput(name, value);
  }
}

/**
 * The `DestroyRef` of the component a dialog was opened with. It fires when
 * MatDialog destroys that component on close, so anything an opener builds
 * *for* the dialog — a filter form and its subscriptions, say — should end on
 * this rather than on the opener's own `DestroyRef`, which can outlive many
 * dialogs.
 */
export function getDialogDestroyRef<T>(dialogRef: MatDialogRef<T>): DestroyRef {
  return getDialogComponentRef(dialogRef, "getDialogDestroyRef").injector.get(DestroyRef);
}

function getDialogComponentRef<T>(dialogRef: MatDialogRef<T>, caller: string): ComponentRef<T> {
  const componentRef = dialogRef.componentRef;
  if (!componentRef) {
    throw new Error(
      `${caller} needs a dialog opened with a component, not a template.`
    );
  }

  return componentRef;
}
