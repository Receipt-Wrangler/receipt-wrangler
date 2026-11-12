import { Component, computed, inject, signal } from "@angular/core";
import { rxResource, takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from "@angular/material/dialog";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { Store } from "@ngxs/store";
import { finalize, map, Observable, of, startWith, switchMap, take } from "rxjs";
import { InputModule } from "../../input";
import { Claim, ReimbursementClaimService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SelectModule } from "../../select/select.module";
import { SnackbarService } from "../../services/snackbar.service";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { AuthState, GroupState } from "../../store";
import { trimmedRequiredValidator } from "../../validators/text-validators";
import {
  applyClaimApiErrors,
  canUpdateClaims,
  CLAIM_STATUS_OPTIONS,
  claimStatusLabel,
  DEFAULT_CLAIM_STATUS,
  receiptGroupIds,
  receiptNoun,
  unmatchedApiErrorMessage,
} from "../claim.utils";
import { reportClaimReceiptsResult } from "../claim-receipts-result-dialog/report-claim-receipts-result";

/** The select value that means "create a claim from this selection". */
export const NEW_CLAIM_OPTION = "new";

export interface AddToClaimDialogReceipt {
  id: number;
  name: string;
  groupId: number;
}

export interface AddToClaimDialogData {
  receipts: AddToClaimDialogReceipt[];
}

interface ClaimTargetOption {
  value: number | typeof NEW_CLAIM_OPTION;
  displayValue: string;
}

/**
 * "Add selected receipts to a claim", from the receipts table. Picks an existing
 * claim of the selection's group, or creates one (create, then add), and reports
 * the per-receipt result itself, so the table only has to open it.
 *
 * A claim holds receipts of exactly one group. A selection made on the All view
 * can span groups; rather than guess which group the user meant and have the
 * server refuse the rest, the dialog says so and offers nothing to submit.
 */
@Component({
  selector: "app-add-to-claim-dialog",
  templateUrl: "./add-to-claim-dialog.component.html",
  styleUrls: ["./add-to-claim-dialog.component.scss"],
  standalone: true,
  imports: [
    InputModule,
    MatProgressSpinnerModule,
    PipesModule,
    ReactiveFormsModule,
    SelectModule,
    SharedUiModule,
  ],
})
export class AddToClaimDialogComponent {
  private readonly claimService = inject(ReimbursementClaimService);
  private readonly formBuilder = inject(FormBuilder);
  private readonly matDialog = inject(MatDialog);
  private readonly snackbarService = inject(SnackbarService);
  private readonly store = inject(Store);

  public readonly matDialogRef = inject<MatDialogRef<AddToClaimDialogComponent, Claim>>(MatDialogRef);

  public readonly data = inject<AddToClaimDialogData>(MAT_DIALOG_DATA);

  public readonly receiptCount = this.data.receipts.length;

  public readonly receiptNoun = receiptNoun(this.receiptCount);

  public readonly groupIds = receiptGroupIds(this.data.receipts);

  /** The one group the selection lives in, or undefined when it spans several. */
  public readonly groupId: number | undefined =
    this.groupIds.length === 1 ? this.groupIds[0] : undefined;

  public readonly groupNames = this.groupIds.map(
    (id) => this.store.selectSnapshot(GroupState.getGroupById(id.toString()))?.name ?? `Group #${id}`
  );

  public readonly canUpdate = canUpdateClaims(
    this.store.selectSnapshot(AuthState.groupPermissions),
    this.groupId
  );

  public readonly claims = rxResource({
    params: () => (this.groupId !== undefined && this.canUpdate ? { groupId: this.groupId } : undefined),
    stream: ({ params }) =>
      this.claimService.getClaimsForGroup(params.groupId).pipe(map((list) => list.claims)),
  });

  public readonly targetOptions = computed<ClaimTargetOption[]>(() => [
    { value: NEW_CLAIM_OPTION, displayValue: "New claim from these receipts" },
    ...(this.claims.hasValue() ? this.claims.value() : []).map((claim) => ({
      value: claim.id,
      displayValue: `${claim.name} (${claimStatusLabel(claim.status)})`,
    })),
  ]);

  public readonly statusOptions = CLAIM_STATUS_OPTIONS;

  public readonly form = this.formBuilder.group({
    // Blank rather than defaulting to an existing claim: adding to the wrong
    // claim is the mistake this dialog most needs to prevent.
    claimId: this.formBuilder.control<number | typeof NEW_CLAIM_OPTION | "">("", Validators.required),
    name: this.formBuilder.control<string>(""),
    status: this.formBuilder.control<string>(DEFAULT_CLAIM_STATUS, Validators.required),
  });

  public readonly isNewClaim = toSignal(
    this.form.controls.claimId.valueChanges.pipe(
      startWith(this.form.controls.claimId.value),
      map((value) => value === NEW_CLAIM_OPTION)
    ),
    { initialValue: false }
  );

  public readonly submitting = signal(false);

  /**
   * Set once a new claim exists, so a retry after the add step failed reuses it
   * rather than creating a second, empty claim.
   */
  private createdClaim?: Claim;

  constructor() {
    this.form.controls.claimId.valueChanges.pipe(takeUntilDestroyed()).subscribe((value) => {
      const name = this.form.controls.name;
      name.setValidators(value === NEW_CLAIM_OPTION ? trimmedRequiredValidator() : null);
      name.updateValueAndValidity();
    });
  }

  public cancel(): void {
    this.matDialogRef.close(undefined);
  }

  public submit(): void {
    if (this.groupId === undefined || this.submitting()) {
      return;
    }

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.submitting.set(true);
    const receiptIds = this.data.receipts.map((receipt) => receipt.id);

    this.resolveClaim(this.groupId)
      .pipe(
        take(1),
        switchMap((claim) =>
          this.claimService
            .addReceiptsToClaim(claim.id, { receiptIds })
            .pipe(map((result) => ({ claim, result })))
        ),
        finalize(() => this.submitting.set(false))
      )
      .subscribe({
        next: ({ claim, result }) => {
          this.matDialogRef.close(claim);
          reportClaimReceiptsResult(this.matDialog, this.snackbarService, {
            action: "add",
            claimName: claim.name,
            requested: this.data.receipts,
            result,
          });
        },
        error: (error) => {
          if (!applyClaimApiErrors(this.form, error)) {
            const message = unmatchedApiErrorMessage(this.form, error);
            if (message) {
              this.snackbarService.error(message);
            }
          }
        },
      });
  }

  private resolveClaim(groupId: number): Observable<Claim> {
    const target = this.form.controls.claimId.value;
    if (target !== NEW_CLAIM_OPTION) {
      const existing = (this.claims.hasValue() ? this.claims.value() : []).find(
        (claim) => claim.id === target
      );
      return existing ? of(existing) : of();
    }

    if (this.createdClaim) {
      return of(this.createdClaim);
    }

    return this.claimService
      .createClaim(groupId, {
        name: (this.form.controls.name.value ?? "").trim(),
        status: this.form.controls.status.value ?? DEFAULT_CLAIM_STATUS,
      })
      .pipe(
        map((claim) => {
          this.createdClaim = claim;
          return claim;
        })
      );
  }
}
