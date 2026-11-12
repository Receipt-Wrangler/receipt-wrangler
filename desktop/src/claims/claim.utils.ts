import { HttpErrorResponse } from "@angular/common/http";
import { FormGroup } from "@angular/forms";
import {
  Permission,
  Receipt,
  UpsertClaimCommand,
  UpsertClaimLedgerEntryCommand,
} from "../open-api";
import { toCalendarDay } from "../utils/app-time-zone";
import { applyApiErrors } from "../utils/form.utils";
import { hasAll } from "../utils/permission.utils";

/*
 * Reimbursement claims (fork feature). Everything here is about the generated
 * `Claim` / `ReimbursementClaimService`, NOT the JWT `Claims` model or the
 * hand-written `ClaimsService` in src/services, which share only the word.
 */

export type ClaimStatusTone = "red" | "green" | "gray" | "yellow";

export interface ClaimStatusOption {
  value: string;
  displayValue: string;
  /** Feeds `app-status-chip`'s `customStatusColor`; undefined is the plain chip. */
  tone?: ClaimStatusTone;
}

/**
 * The statuses the server knows today. The wire field is a plain string, so a
 * status missing from this list still renders (as its raw code) rather than
 * breaking the page; only the label and the colour are authored here.
 */
export const CLAIM_STATUS_OPTIONS: ClaimStatusOption[] = [
  { value: "DRAFT", displayValue: "Draft" },
  { value: "SUBMITTED", displayValue: "Submitted", tone: "yellow" },
  { value: "PARTIALLY_PAID", displayValue: "Partially Paid", tone: "yellow" },
  { value: "PAID", displayValue: "Paid", tone: "green" },
  { value: "DISPUTED", displayValue: "Disputed", tone: "red" },
  { value: "CLOSED", displayValue: "Closed", tone: "gray" },
];

export const DEFAULT_CLAIM_STATUS = "DRAFT";

export const CLAIM_LEDGER_PAYMENT = "PAYMENT";
export const CLAIM_LEDGER_RECOUPMENT = "RECOUPMENT";
export const CLAIM_LEDGER_NOTE = "NOTE";

export interface ClaimLedgerTypeOption {
  value: string;
  displayValue: string;
  icon: string;
}

export const CLAIM_LEDGER_TYPE_OPTIONS: ClaimLedgerTypeOption[] = [
  { value: CLAIM_LEDGER_PAYMENT, displayValue: "Payment received", icon: "payments" },
  { value: CLAIM_LEDGER_RECOUPMENT, displayValue: "Recoupment (taken back)", icon: "undo" },
  { value: CLAIM_LEDGER_NOTE, displayValue: "Note", icon: "sticky_note_2" },
];

/** The label for a status, or the raw code for one this client does not know. */
export function claimStatusLabel(status: string | null | undefined): string {
  return CLAIM_STATUS_OPTIONS.find((option) => option.value === status)?.displayValue ?? status ?? "";
}

export function claimStatusTone(status: string | null | undefined): ClaimStatusTone | undefined {
  return CLAIM_STATUS_OPTIONS.find((option) => option.value === status)?.tone;
}

/**
 * The options a status select offers: the known statuses, plus the claim's own
 * when it is one this client does not know, so the select can still show the
 * value it holds instead of rendering blank.
 */
export function claimStatusOptionsFor(current: string | null | undefined): ClaimStatusOption[] {
  if (!current || CLAIM_STATUS_OPTIONS.some((option) => option.value === current)) {
    return CLAIM_STATUS_OPTIONS;
  }

  return [...CLAIM_STATUS_OPTIONS, { value: current, displayValue: current }];
}

export function claimLedgerTypeOption(type: string | null | undefined): ClaimLedgerTypeOption {
  return (
    CLAIM_LEDGER_TYPE_OPTIONS.find((option) => option.value === type) ?? {
      value: type ?? "",
      displayValue: type ?? "",
      icon: "help_outline",
    }
  );
}

/** Same rule as `claimStatusOptionsFor`, for a ledger entry's type. */
export function claimLedgerTypeOptionsFor(current: string | null | undefined): ClaimLedgerTypeOption[] {
  if (!current || CLAIM_LEDGER_TYPE_OPTIONS.some((option) => option.value === current)) {
    return CLAIM_LEDGER_TYPE_OPTIONS;
  }

  return [...CLAIM_LEDGER_TYPE_OPTIONS, claimLedgerTypeOption(current)];
}

/** Reading claims is gated on the same permission as reading the group's receipts. */
export function canReadClaims(
  groupPermissions: { [groupId: number]: string[] },
  groupId: number | null | undefined
): boolean {
  return !!groupId && hasAll(groupPermissions[groupId] ?? [], Permission.GroupReceiptsRead);
}

/** Every claim write is gated on updating the group's receipts, in the claim's own group. */
export function canUpdateClaims(
  groupPermissions: { [groupId: number]: string[] },
  groupId: number | null | undefined
): boolean {
  return !!groupId && hasAll(groupPermissions[groupId] ?? [], Permission.GroupReceiptsUpdate);
}

/** A trimmed string, or undefined when blank, so an emptied optional field is omitted. */
export function optionalText(value: unknown): string | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }

  const trimmed = String(value).trim();
  return trimmed === "" ? undefined : trimmed;
}

/**
 * A datepicker value as the bare `yyyy-MM-dd` the contract wants, or undefined.
 * The datepicker holds the LOCAL midnight of the picked day, so the day is read
 * on the browser's calendar: no time zone conversion, because it is a day, not
 * an instant.
 */
export function optionalCalendarDay(value: Date | string | null | undefined): string | undefined {
  if (value === null || value === undefined || value === "") {
    return undefined;
  }

  return toCalendarDay(value);
}

/**
 * The masked currency input yields a string or a number depending on the
 * currency settings; the contract wants a decimal string, and the server reads a
 * blank as "no amount".
 */
export function optionalMoney(value: unknown): string | undefined {
  return optionalText(value);
}

export interface ClaimFormValue {
  name: string;
  referenceNumber: string | null;
  type: string | null;
  status: string;
  submittedDate: Date | string | null;
  dueDate: Date | string | null;
  claimedAmount: string | number | null;
  notes: string | null;
}

/**
 * PUT is a full replace of the claim's own fields, so an emptied optional field
 * is sent as absent, which clears it.
 */
export function buildUpsertClaimCommand(value: ClaimFormValue): UpsertClaimCommand {
  return {
    name: String(value.name ?? "").trim(),
    referenceNumber: optionalText(value.referenceNumber),
    type: optionalText(value.type),
    status: value.status,
    submittedDate: optionalCalendarDay(value.submittedDate),
    dueDate: optionalCalendarDay(value.dueDate),
    claimedAmount: optionalMoney(value.claimedAmount),
    notes: optionalText(value.notes),
  };
}

export interface ClaimLedgerEntryFormValue {
  type: string;
  date: Date | string | null;
  amount: string | number | null;
  note: string | null;
}

export function buildUpsertClaimLedgerEntryCommand(
  value: ClaimLedgerEntryFormValue
): UpsertClaimLedgerEntryCommand {
  return {
    type: value.type,
    date: optionalCalendarDay(value.date) ?? "",
    amount: optionalMoney(value.amount),
    note: optionalText(value.note),
  };
}

/**
 * Puts a claim endpoint's field-keyed 400 onto the matching form controls, via
 * the app's existing `applyApiErrors`, and marks them touched so the message
 * shows. Returns whether any control took a message.
 *
 * An `{ errorMsg }` body is left alone: the HTTP interceptor already toasts it.
 */
export function applyClaimApiErrors(form: FormGroup, error: unknown): boolean {
  if (!(error instanceof HttpErrorResponse) || error.status !== 400) {
    return false;
  }

  const body = error.error;
  if (!body || typeof body !== "object" || "errorMsg" in body) {
    return false;
  }

  applyApiErrors(form, error);
  const matched = Object.keys(body).filter((key) => !!form.get(key));
  matched.forEach((key) => form.get(key)?.markAsTouched());

  return matched.length > 0;
}

/**
 * The first message of a field-keyed 400 whose key matches no control, so a
 * caller can still say something rather than fail silently (e.g. `receiptIds`).
 */
export function unmatchedApiErrorMessage(form: FormGroup, error: unknown): string | undefined {
  if (!(error instanceof HttpErrorResponse) || error.status !== 400) {
    return undefined;
  }

  const body = error.error;
  if (!body || typeof body !== "object" || "errorMsg" in body) {
    return undefined;
  }

  const key = Object.keys(body).find((candidate) => !form.get(candidate));
  return key ? String(body[key]) : undefined;
}

/** The distinct groups a receipt selection spans, in first-seen order. */
export function receiptGroupIds(receipts: Pick<Receipt, "groupId">[]): number[] {
  return Array.from(new Set(receipts.map((receipt) => receipt.groupId)));
}

export function receiptNoun(count: number): string {
  return count === 1 ? "receipt" : "receipts";
}

/*
 * A receipt's own claim. The server sets `claimId` / `claimName` on every
 * receipt it returns, but only when the receipt is in a claim the viewer may
 * see, so their absence never proves a receipt is unclaimed: it means "nothing
 * to show". The server's refusal stays the authority for the hidden case.
 */
export type ReceiptClaimRef = Pick<Receipt, "claimId" | "claimName">;

/** Whether the receipt is in a (visible) claim other than `claimId`. */
export function isInOtherClaim(receipt: ReceiptClaimRef, claimId: number | undefined): boolean {
  return !!receipt.claimId && receipt.claimId !== claimId;
}

export interface ClaimedReceiptGroup {
  claimId: number;
  claimName: string;
  receiptNames: string[];
}

/**
 * The receipts of a selection that are already in a visible claim other than
 * `targetClaimId`, grouped by that claim in first-seen order. A receipt already
 * in the target claim is left out: the server counts it as added.
 */
export function claimedReceiptsByClaim(
  receipts: (ReceiptClaimRef & Pick<Receipt, "name">)[],
  targetClaimId?: number
): ClaimedReceiptGroup[] {
  const groups = new Map<number, ClaimedReceiptGroup>();
  for (const receipt of receipts) {
    if (!isInOtherClaim(receipt, targetClaimId)) {
      continue;
    }

    const claimId = receipt.claimId as number;
    const group = groups.get(claimId) ?? {
      claimId,
      claimName: receipt.claimName || `Claim #${claimId}`,
      receiptNames: [],
    };
    group.receiptNames.push(receipt.name);
    groups.set(claimId, group);
  }

  return Array.from(groups.values());
}

/**
 * The receipt form's warning when the group field is changed on a receipt in a
 * claim: the server unlinks a receipt that moves group, because a claim holds
 * one group's receipts. Undefined while the group is unchanged or blank.
 */
export function claimGroupMoveWarning(
  receipt: (ReceiptClaimRef & Pick<Receipt, "groupId">) | undefined,
  groupId: unknown
): string | undefined {
  if (!receipt?.claimName || groupId === null || groupId === undefined || groupId === "") {
    return undefined;
  }

  if (Number(groupId) === receipt.groupId) {
    return undefined;
  }

  return `Moving this receipt to another group will remove it from claim '${receipt.claimName}'.`;
}
