import { DestroyRef } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { FilterOperation } from "../open-api/index";
import { buildFieldFormGroup, listenForBetweenOperation } from "./filter-form";

export function buildReceiptFilterForm(filter: any, destroyRef: DestroyRef): FormGroup {
  const formGroup = new FormGroup({
    date: buildFieldFormGroup(
      filter?.date?.value,
      filter?.date?.operation,
      destroyRef,
      filter?.date?.operation === FilterOperation.Between
    ),
    amount: buildFieldFormGroup(
      filter?.amount?.value,
      filter?.amount?.operation,
      destroyRef,
      filter?.amount?.operation === FilterOperation.Between
    ),
    name: buildFieldFormGroup(
      filter?.name?.value,
      filter?.name?.operation,
      destroyRef
    ),
    paidBy: buildFieldFormGroup(
      filter?.paidBy?.value ?? [],
      filter?.paidBy?.operation,
      destroyRef,
      true
    ),
    categories: buildFieldFormGroup(
      filter?.categories?.value ?? [],
      filter?.categories?.operation,
      destroyRef,
      true
    ),
    tags: buildFieldFormGroup(
      filter?.tags?.value ?? [],
      filter?.tags?.operation,
      destroyRef,
      true
    ),
    status: buildFieldFormGroup(
      filter?.status?.value ?? [],
      filter?.status?.operation,
      destroyRef,
      true
    ),
    group: buildFieldFormGroup(
      filter?.group?.value ?? [],
      filter?.group?.operation,
      destroyRef,
      true
    ),
    resolvedDate: buildFieldFormGroup(
      filter?.resolvedDate?.value,
      filter?.resolvedDate?.operation,
      destroyRef,
      filter?.resolvedDate?.operation === FilterOperation.Between
    ),
    createdAt: buildFieldFormGroup(
      filter?.createdAt?.value,
      filter?.createdAt?.operation,
      destroyRef,
      filter?.createdAt?.operation === FilterOperation.Between
    ),
  });

  listenForBetweenOperation(formGroup, "amount", destroyRef);
  listenForBetweenOperation(formGroup, "date", destroyRef);
  listenForBetweenOperation(formGroup, "resolvedDate", destroyRef);
  listenForBetweenOperation(formGroup, "createdAt", destroyRef);


  return formGroup;
}
