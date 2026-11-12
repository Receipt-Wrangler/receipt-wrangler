import { Component, inject, signal } from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { finalize, map, Observable, startWith, take } from "rxjs";
import { DatepickerModule } from "../../datepicker/datepicker.module";
import { InputModule } from "../../input";
import { ClaimLedgerEntry, ReimbursementClaimService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SelectModule } from "../../select/select.module";
import { SnackbarService } from "../../services/snackbar.service";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SystemSettingsState } from "../../store/system-settings.state";
import { TextareaModule } from "../../textarea/textarea.module";
import { calendarDayToLocalDate, todayInZone } from "../../utils/app-time-zone";
import { trimmedRequiredValidator } from "../../validators/text-validators";
import {
  applyClaimApiErrors,
  buildUpsertClaimLedgerEntryCommand,
  CLAIM_LEDGER_NOTE,
  CLAIM_LEDGER_PAYMENT,
  claimLedgerTypeOptionsFor,
  unmatchedApiErrorMessage,
} from "../claim.utils";

export interface ClaimLedgerEntryDialogData {
  claimId: number;
  /** Absent to add an entry. */
  entry?: ClaimLedgerEntry;
}

/**
 * Adds or edits one ledger entry: a payment received, a recoupment (money taken
 * back), or a note. Closes with `true` once the write lands, so the claim page
 * knows to refetch the claim for its new figures — the write returns only the
 * entry.
 *
 * Amount is required for a payment or a recoupment and optional for a note,
 * whose amount never counts; a note's text is required. The server enforces the
 * same rules and its field-keyed 400s land on the matching controls.
 */
@Component({
  selector: "app-claim-ledger-entry-dialog",
  templateUrl: "./claim-ledger-entry-dialog.component.html",
  styleUrls: ["./claim-ledger-entry-dialog.component.scss"],
  standalone: true,
  imports: [
    DatepickerModule,
    InputModule,
    PipesModule,
    ReactiveFormsModule,
    SelectModule,
    SharedUiModule,
    TextareaModule,
  ],
})
export class ClaimLedgerEntryDialogComponent {
  private readonly claimService = inject(ReimbursementClaimService);
  private readonly formBuilder = inject(FormBuilder);
  private readonly snackbarService = inject(SnackbarService);
  private readonly store = inject(Store);

  public readonly matDialogRef = inject<MatDialogRef<ClaimLedgerEntryDialogComponent, boolean>>(MatDialogRef);

  public readonly data = inject<ClaimLedgerEntryDialogData>(MAT_DIALOG_DATA);

  public readonly headerText = this.data.entry ? "Edit ledger entry" : "Add ledger entry";

  public readonly typeOptions = claimLedgerTypeOptionsFor(this.data.entry?.type);

  public readonly form = this.formBuilder.group({
    type: this.formBuilder.control<string>(this.data.entry?.type ?? CLAIM_LEDGER_PAYMENT, Validators.required),
    // A calendar day: the datepicker holds the local midnight of the day, and a
    // new entry defaults to today in the app time zone, not the browser's.
    date: this.formBuilder.control<Date | null>(
      calendarDayToLocalDate(this.data.entry?.date) ??
        todayInZone(this.store.selectSnapshot(SystemSettingsState.timeZone)),
      Validators.required
    ),
    amount: this.formBuilder.control<string | number | null>(this.data.entry?.amount ?? ""),
    note: this.formBuilder.control<string>(this.data.entry?.note ?? ""),
  });

  public readonly isNote = toSignal(
    this.form.controls.type.valueChanges.pipe(
      startWith(this.form.controls.type.value),
      map((type) => type === CLAIM_LEDGER_NOTE)
    ),
    { initialValue: false }
  );

  public readonly submitting = signal(false);

  constructor() {
    this.form.controls.type.valueChanges
      .pipe(startWith(this.form.controls.type.value), takeUntilDestroyed())
      .subscribe((type) => this.applyTypeValidators(type));
  }

  private applyTypeValidators(type: string | null): void {
    const { amount, note } = this.form.controls;
    const isNote = type === CLAIM_LEDGER_NOTE;

    amount.setValidators(isNote ? null : Validators.required);
    note.setValidators(isNote ? trimmedRequiredValidator() : null);
    amount.updateValueAndValidity({ emitEvent: false });
    note.updateValueAndValidity({ emitEvent: false });
  }

  public cancel(): void {
    this.matDialogRef.close(false);
  }

  public submit(): void {
    if (this.submitting()) {
      return;
    }

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.submitting.set(true);
    this.save()
      .pipe(
        take(1),
        finalize(() => this.submitting.set(false))
      )
      .subscribe({
        next: () => {
          this.snackbarService.success(this.data.entry ? "Ledger entry updated" : "Ledger entry added");
          this.matDialogRef.close(true);
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

  private save(): Observable<unknown> {
    const command = buildUpsertClaimLedgerEntryCommand(this.form.getRawValue() as any);
    const entry = this.data.entry;

    return entry
      ? this.claimService.updateClaimLedgerEntry(this.data.claimId, entry.id, command)
      : this.claimService.createClaimLedgerEntry(this.data.claimId, command);
  }
}
