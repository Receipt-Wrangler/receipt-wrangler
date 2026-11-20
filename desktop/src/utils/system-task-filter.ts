import { DestroyRef } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { FilterOperation, SystemTaskPagedRequestFilter } from "../open-api/index";
import { toDateWireFilter } from "./date-wire-filter";
import { buildFieldFormGroup, listenForBetweenOperation } from "./filter-form";

/**
 * The system task filter's reactive form — the same `{ operation, value }`
 * shape as `buildReceiptFilterForm`, built from the same shared helpers, so the
 * two dialogs behave identically.
 *
 * `destroyRef` is the calling component's: the helpers tie their subscriptions
 * to its lifetime.
 */
export function buildSystemTaskFilterForm(filter: any, destroyRef: DestroyRef): FormGroup {
  const formGroup = new FormGroup({
    type: buildFieldFormGroup(
      filter?.type?.value ?? [],
      filter?.type?.operation,
      destroyRef,
      true
    ),
    ranBy: buildFieldFormGroup(
      filter?.ranBy?.value ?? [],
      filter?.ranBy?.operation,
      destroyRef,
      true
    ),
    startedAt: buildFieldFormGroup(
      filter?.startedAt?.value,
      filter?.startedAt?.operation,
      destroyRef,
      filter?.startedAt?.operation === FilterOperation.Between
    ),
    endedAt: buildFieldFormGroup(
      filter?.endedAt?.value,
      filter?.endedAt?.operation,
      destroyRef,
      filter?.endedAt?.operation === FilterOperation.Between
    ),
  });

  listenForBetweenOperation(formGroup, "startedAt", destroyRef);
  listenForBetweenOperation(formGroup, "endedAt", destroyRef);

  return formGroup;
}

/** The system task filter's date fields — instants, resolved in the app zone. */
export const SYSTEM_TASK_DATE_FILTER_KEYS = ["startedAt", "endedAt"] as const;

/**
 * The filter as it goes on the wire: `startedAt` / `endedAt` values become the
 * calendar day the user picked, `yyyy-MM-dd`. See `toDateWireFilter`.
 */
export function toSystemTaskWireFilter(
  filter: SystemTaskPagedRequestFilter | undefined | null,
): SystemTaskPagedRequestFilter | undefined {
  return toDateWireFilter(filter, SYSTEM_TASK_DATE_FILTER_KEYS);
}
