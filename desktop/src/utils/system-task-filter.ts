import { FormGroup } from "@angular/forms";
import { FilterOperation } from "../open-api/index";
import { buildFieldFormGroup, listenForBetweenOperation } from "./filter";

export function buildSystemTaskFilterForm(filter: any, thisContext: any): FormGroup {
  const formGroup = new FormGroup({
    type: buildFieldFormGroup(
      filter?.type?.value ?? [],
      filter?.type?.operation,
      thisContext,
      true
    ),
    status: buildFieldFormGroup(
      filter?.status?.value ?? [],
      filter?.status?.operation,
      thisContext,
      true
    ),
    ranByUserId: buildFieldFormGroup(
      filter?.ranByUserId?.value ?? [],
      filter?.ranByUserId?.operation,
      thisContext,
      true
    ),
    startedAt: buildFieldFormGroup(
      filter?.startedAt?.value,
      filter?.startedAt?.operation,
      thisContext,
      filter?.startedAt?.operation === FilterOperation.Between
    ),
    endedAt: buildFieldFormGroup(
      filter?.endedAt?.value,
      filter?.endedAt?.operation,
      thisContext,
      filter?.endedAt?.operation === FilterOperation.Between
    ),
  });

  listenForBetweenOperation(formGroup, "startedAt", thisContext);
  listenForBetweenOperation(formGroup, "endedAt", thisContext);

  return formGroup;
}
