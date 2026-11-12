import { MatDialog } from "@angular/material/dialog";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { ClaimReceiptsResult } from "../../open-api";
import { SnackbarService } from "../../services/snackbar.service";
import { receiptNoun } from "../claim.utils";
import {
  ClaimReceiptsResultDialogComponent,
  ClaimReceiptsResultDialogData,
} from "./claim-receipts-result-dialog.component";

export type ClaimReceiptsAction = "add" | "remove";

export interface ReportClaimReceiptsResultOptions {
  action: ClaimReceiptsAction;
  claimName: string;
  /** What was asked for, so each refusal can be named; unknown ids fall back to "Receipt #id". */
  requested: { id: number; name?: string }[];
  result: ClaimReceiptsResult;
}

/**
 * The one way an add/remove-receipts result is reported, so the bulk action on
 * the receipts table and the claim page say it the same way: a snackbar when
 * everything went through, otherwise a dialog naming each refused receipt.
 */
export function reportClaimReceiptsResult(
  matDialog: MatDialog,
  snackbarService: SnackbarService,
  { action, claimName, requested, result }: ReportClaimReceiptsResultOptions
): void {
  const done = result.receiptIds.length;
  const verb = action === "add" ? "Added" : "Removed";
  const preposition = action === "add" ? "to" : "from";

  if (result.failures.length === 0) {
    snackbarService.success(`${verb} ${done} ${receiptNoun(done)} ${preposition} ${claimName}`);
    return;
  }

  matDialog.open<ClaimReceiptsResultDialogComponent, ClaimReceiptsResultDialogData>(
    ClaimReceiptsResultDialogComponent,
    {
      ...DEFAULT_DIALOG_CONFIG,
      data: {
        headerText:
          action === "add"
            ? "Some receipts could not be added"
            : "Some receipts could not be removed",
        summary: `${verb} ${done} of ${requested.length} ${receiptNoun(requested.length)} ${preposition} ${claimName}.`,
        failures: result.failures.map((failure) => ({
          name:
            requested.find((receipt) => receipt.id === failure.receiptId)?.name ??
            `Receipt #${failure.receiptId}`,
          message: failure.message,
        })),
      },
    }
  );
}
