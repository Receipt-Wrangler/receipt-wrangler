import { AfterViewInit, Component, computed, inject, OnInit, signal, TemplateRef, viewChild } from "@angular/core";
import { rxResource } from "@angular/core/rxjs-interop";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { MatTableDataSource } from "@angular/material/table";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { filter, finalize, map, switchMap, take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG, DEFAULT_HOST_CLASS } from "../../constants";
import { FormMode } from "../../enums/form-mode.enum";
import {
  Claim,
  ClaimLedgerEntry,
  Receipt,
  ReceiptService,
  ReimbursementClaimService,
  SortDirection,
} from "../../open-api";
import { ReceiptExportService } from "../../services/receipt-export.service";
import { SnackbarService } from "../../services/snackbar.service";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { AuthState, GroupState } from "../../store";
import { TableColumn } from "../../table/table-column.interface";
import { calendarDayToLocalDate } from "../../utils/app-time-zone";
import { trimmedRequiredValidator } from "../../validators/text-validators";
import { CLAIM_NOT_FOUND_MESSAGE } from "../claim.resolver";
import { CLAIM_PACKET_DESCRIPTION, ClaimPacketDownloadService } from "../claim-packet-download.service";
import {
  applyClaimApiErrors,
  buildUpsertClaimCommand,
  canUpdateClaims,
  CLAIM_LEDGER_NOTE,
  CLAIM_LEDGER_PAYMENT,
  CLAIM_LEDGER_RECOUPMENT,
  claimLedgerTypeOption,
  claimStatusLabel,
  claimStatusOptionsFor,
  claimStatusTone,
  DEFAULT_CLAIM_STATUS,
  isInOtherClaim,
  unmatchedApiErrorMessage,
} from "../claim.utils";
import {
  ClaimLedgerEntryDialogComponent,
  ClaimLedgerEntryDialogData,
} from "../claim-ledger-entry-dialog/claim-ledger-entry-dialog.component";
import {
  ClaimReceiptPickerDialogComponent,
  ClaimReceiptPickerDialogData,
} from "../claim-receipt-picker-dialog/claim-receipt-picker-dialog.component";
import { reportClaimReceiptsResult } from "../claim-receipts-result-dialog/report-claim-receipts-result";

/** A receipt row on the claim, or a placeholder for one the group list did not return. */
type ClaimReceiptRow = Pick<Receipt, "id" | "name"> & Partial<Receipt>;

/**
 * One claim: its own fields (a form, read-only in view mode), its receipts, and
 * its ledger of payments, recoupments and notes, under the server's figures.
 *
 * The receipts and the ledger are live in view AND edit mode — each change is
 * its own request, followed by a refetch of the claim because those writes
 * return only what they wrote. A refetch never rebuilds the form, so an edit in
 * progress on the claim's own fields survives a ledger change.
 *
 * The claim arrives through `claimResolverFn`, which also redirects a missing or
 * forbidden claim and an edit the viewer may not make. Writes are gated on
 * `group.receipts.update` in the claim's OWN group, never the group being
 * browsed, which may be the All group.
 */
@Component({
  selector: "app-claim-form",
  templateUrl: "./claim-form.component.html",
  styleUrls: ["./claim-form.component.scss"],
  host: DEFAULT_HOST_CLASS,
  standalone: false,
})
export class ClaimFormComponent implements OnInit, AfterViewInit {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly claimService = inject(ReimbursementClaimService);
  private readonly claimPacketDownloadService = inject(ClaimPacketDownloadService);
  private readonly formBuilder = inject(FormBuilder);
  private readonly matDialog = inject(MatDialog);
  private readonly receiptExportService = inject(ReceiptExportService);
  private readonly receiptService = inject(ReceiptService);
  private readonly router = inject(Router);
  private readonly snackbarService = inject(SnackbarService);
  private readonly store = inject(Store);

  private readonly receiptDateCell = viewChild.required<TemplateRef<any>>("receiptDateCell");
  private readonly receiptNameCell = viewChild.required<TemplateRef<any>>("receiptNameCell");
  private readonly receiptPaidByCell = viewChild.required<TemplateRef<any>>("receiptPaidByCell");
  private readonly receiptAmountCell = viewChild.required<TemplateRef<any>>("receiptAmountCell");
  private readonly receiptStatusCell = viewChild.required<TemplateRef<any>>("receiptStatusCell");
  private readonly receiptActionsCell = viewChild.required<TemplateRef<any>>("receiptActionsCell");
  private readonly ledgerDateCell = viewChild.required<TemplateRef<any>>("ledgerDateCell");
  private readonly ledgerTypeCell = viewChild.required<TemplateRef<any>>("ledgerTypeCell");
  private readonly ledgerAmountCell = viewChild.required<TemplateRef<any>>("ledgerAmountCell");
  private readonly ledgerNoteCell = viewChild.required<TemplateRef<any>>("ledgerNoteCell");
  private readonly ledgerActionsCell = viewChild.required<TemplateRef<any>>("ledgerActionsCell");

  public readonly FormMode = FormMode;

  public readonly mode: FormMode = this.activatedRoute.snapshot.data["mode"] ?? FormMode.view;

  public readonly claim = signal<Claim | undefined>(this.activatedRoute.snapshot.data["claim"]);

  /** The claim's group, or — when adding — the group the claim would be created in. */
  public readonly groupId = computed(
    () => this.claim()?.groupId ?? this.store.selectSnapshot(GroupState.addTargetGroupId)
  );

  private readonly groupPermissions = this.store.selectSignal(AuthState.groupPermissions);

  private readonly realGroups = this.store.selectSignal(GroupState.groupsWithoutAll);

  public readonly canUpdate = computed(() => canUpdateClaims(this.groupPermissions(), this.groupId()));

  /** Only worth naming when the viewer belongs to more than one group. */
  public readonly groupName = computed(() => {
    const groups = this.realGroups();
    if (groups.length <= 1) {
      return "";
    }

    return groups.find((group) => group.id === this.groupId())?.name ?? "";
  });

  public readonly headerText = computed(() => {
    const name = this.claim()?.name ?? "";
    switch (this.mode) {
      case FormMode.add:
        return "Add Claim";
      case FormMode.edit:
        return `Edit ${name}`;
      default:
        return name;
    }
  });

  public readonly statusOptions = claimStatusOptionsFor(this.claim()?.status);

  public form: FormGroup = new FormGroup({});

  private readonly claimGroupId = computed(() => this.claim()?.groupId);

  /**
   * The claim group's receipts, loaded once with the receipt list's own endpoint
   * (unpaged). One request serves both halves of the receipts section: the rows
   * of the receipts on the claim, and the picker's candidates.
   */
  public readonly groupReceipts = rxResource({
    // Keyed on the primitive claimGroupId, not on claim(): every receipt or
    // ledger change refetches the claim, and a params object rebuilt from it
    // would reload the whole group's receipts each time.
    params: () => {
      const groupId = this.claimGroupId();
      return groupId ? { groupId } : undefined;
    },
    stream: ({ params }) =>
      this.receiptService
        .getReceiptsForGroup(params.groupId, {
          page: 1,
          pageSize: -1,
          orderBy: "date",
          sortDirection: SortDirection.Desc,
        })
        // PagedData is a union over every paged type; this endpoint returns receipts.
        .pipe(map((paged) => paged.data as unknown as Receipt[])),
  });

  public readonly claimReceipts = computed<ClaimReceiptRow[]>(() => {
    const receipts = this.groupReceipts.hasValue() ? this.groupReceipts.value() : [];
    const byId = new Map(receipts.map((receipt) => [receipt.id, receipt]));

    return (this.claim()?.receiptIds ?? []).map(
      (id) => byId.get(id) ?? { id, name: `Receipt #${id}` }
    );
  });

  /**
   * What the picker offers: the group's receipts not on this claim and not in
   * another claim the viewer can see. The group list is loaded once, so a
   * receipt's `claimId` there can be stale; both directions are harmless â€” one
   * added here since is excluded by `receiptIds`, and one removed from here
   * still reads as this claim's, which is not "another". One in a hidden claim
   * carries no `claimId`, is offered, and the server's refusal reports it.
   */
  private readonly pickerReceipts = computed(() => {
    const claim = this.claim();
    const onClaim = new Set(claim?.receiptIds ?? []);
    const receipts = this.groupReceipts.hasValue() ? this.groupReceipts.value() : [];
    const notOnClaim = receipts.filter((receipt) => !onClaim.has(receipt.id));
    const candidates = notOnClaim.filter((receipt) => !isInOtherClaim(receipt, claim?.id));

    return { candidates, inOtherClaimsCount: notOnClaim.length - candidates.length };
  });

  public readonly receiptsDataSource = computed(
    () => new MatTableDataSource<ClaimReceiptRow>(this.claimReceipts())
  );

  public readonly ledgerDataSource = computed(
    () => new MatTableDataSource<ClaimLedgerEntry>(this.claim()?.ledgerEntries ?? [])
  );

  private readonly viewReady = signal(false);

  public readonly receiptColumns = computed<TableColumn[]>(() => {
    if (!this.viewReady()) {
      return [];
    }

    return [
      { columnHeader: "Date", matColumnDef: "date", template: this.receiptDateCell(), sortable: false },
      { columnHeader: "Name", matColumnDef: "name", template: this.receiptNameCell(), sortable: false },
      { columnHeader: "Paid By", matColumnDef: "paidBy", template: this.receiptPaidByCell(), sortable: false },
      { columnHeader: "Amount", matColumnDef: "amount", template: this.receiptAmountCell(), sortable: false },
      { columnHeader: "Status", matColumnDef: "status", template: this.receiptStatusCell(), sortable: false },
      ...(this.canUpdate()
        ? [{ columnHeader: "", matColumnDef: "actions", template: this.receiptActionsCell(), sortable: false }]
        : []),
    ];
  });

  public readonly ledgerColumns = computed<TableColumn[]>(() => {
    if (!this.viewReady()) {
      return [];
    }

    return [
      { columnHeader: "Date", matColumnDef: "date", template: this.ledgerDateCell(), sortable: false },
      { columnHeader: "Type", matColumnDef: "type", template: this.ledgerTypeCell(), sortable: false },
      { columnHeader: "Amount", matColumnDef: "amount", template: this.ledgerAmountCell(), sortable: false },
      { columnHeader: "Note", matColumnDef: "note", template: this.ledgerNoteCell(), sortable: false },
      ...(this.canUpdate()
        ? [{ columnHeader: "", matColumnDef: "actions", template: this.ledgerActionsCell(), sortable: false }]
        : []),
    ];
  });

  public readonly receiptDisplayedColumns = computed(() =>
    this.receiptColumns().map((column) => column.matColumnDef)
  );

  public readonly ledgerDisplayedColumns = computed(() =>
    this.ledgerColumns().map((column) => column.matColumnDef)
  );

  public readonly statusLabel = claimStatusLabel;

  public readonly statusTone = claimStatusTone;

  public readonly ledgerTypeOption = claimLedgerTypeOption;

  public readonly PAYMENT = CLAIM_LEDGER_PAYMENT;

  public readonly RECOUPMENT = CLAIM_LEDGER_RECOUPMENT;

  public readonly NOTE = CLAIM_LEDGER_NOTE;

  public readonly packetDescription = CLAIM_PACKET_DESCRIPTION;

  /** True while the packet builds: the server sends nothing until the zip is complete. */
  public readonly packetBusy = signal(false);

  public ngOnInit(): void {
    if (this.mode === FormMode.add && !this.groupId()) {
      // Only reachable by URL: the list offers "Add Claim" only when there is a
      // real group to add it to.
      this.snackbarService.error("Choose a group before adding a claim.");
      this.router.navigate([this.store.selectSnapshot(GroupState.dashboardLink)]);
      return;
    }

    this.initForm();
  }

  public ngAfterViewInit(): void {
    this.viewReady.set(true);
  }

  private initForm(): void {
    const claim = this.claim();

    this.form = this.formBuilder.group({
      name: [claim?.name ?? "", trimmedRequiredValidator()],
      referenceNumber: [claim?.referenceNumber ?? ""],
      type: [claim?.type ?? ""],
      status: [claim?.status ?? DEFAULT_CLAIM_STATUS, Validators.required],
      // Calendar days: the datepicker works in local midnights, so the stored
      // day becomes the local midnight of that same day — never shifted by zone.
      submittedDate: [calendarDayToLocalDate(claim?.submittedDate)],
      dueDate: [calendarDayToLocalDate(claim?.dueDate)],
      claimedAmount: [claim?.claimedAmount ?? ""],
      notes: [claim?.notes ?? ""],
    });
  }

  public isNegative(amount: string | undefined): boolean {
    return !!amount && amount.trim().startsWith("-");
  }

  public submit(): void {
    if (this.mode === FormMode.view) {
      return;
    }

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const command = buildUpsertClaimCommand(this.form.getRawValue());
    const claim = this.claim();
    const request = claim
      ? this.claimService.updateClaim(claim.id, command)
      : this.claimService.createClaim(this.groupId() as number, command);

    request.pipe(take(1)).subscribe({
      next: (saved) => {
        this.snackbarService.success(claim ? "Claim updated" : "Claim added");
        this.router.navigate(["/claims", saved.id, "view"]);
      },
      error: (error) => this.reportFormError(error),
    });
  }

  private reportFormError(error: unknown): void {
    if (!applyClaimApiErrors(this.form, error)) {
      const message = unmatchedApiErrorMessage(this.form, error);
      if (message) {
        this.snackbarService.error(message);
      }
    }
  }

  /**
   * Re-reads the claim after a receipt or ledger change, since those writes
   * return only what they wrote. If the claim has gone in the meantime (deleted
   * in another session), leave rather than show stale figures.
   */
  public refresh(): void {
    const claim = this.claim();
    if (!claim) {
      return;
    }

    this.claimService
      .getClaimById(claim.id)
      .pipe(take(1))
      .subscribe({
        next: (fresh) => this.claim.set(fresh),
        error: () => {
          this.snackbarService.error(CLAIM_NOT_FOUND_MESSAGE);
          this.navigateToList();
        },
      });
  }

  public addReceipts(): void {
    const claim = this.claim();
    if (!claim) {
      return;
    }

    this.matDialog
      .open<ClaimReceiptPickerDialogComponent, ClaimReceiptPickerDialogData, Receipt[]>(
        ClaimReceiptPickerDialogComponent,
        {
          ...DEFAULT_DIALOG_CONFIG,
          data: { claimName: claim.name, ...this.pickerReceipts() },
        }
      )
      .afterClosed()
      .pipe(
        take(1),
        filter((picked): picked is Receipt[] => !!picked && picked.length > 0),
        switchMap((picked) =>
          this.claimService
            .addReceiptsToClaim(claim.id, { receiptIds: picked.map((receipt) => receipt.id) })
            .pipe(
              tap((result) =>
                reportClaimReceiptsResult(this.matDialog, this.snackbarService, {
                  action: "add",
                  claimName: claim.name,
                  requested: picked,
                  result,
                })
              )
            )
        )
      )
      .subscribe({ next: () => this.refresh(), error: () => this.refresh() });
  }

  /** No confirmation: putting a receipt back is one click on "Add receipts". */
  public removeReceipt(receipt: ClaimReceiptRow): void {
    const claim = this.claim();
    if (!claim) {
      return;
    }

    this.claimService
      .removeReceiptsFromClaim(claim.id, { receiptIds: [receipt.id] })
      .pipe(take(1))
      .subscribe({
        next: (result) => {
          reportClaimReceiptsResult(this.matDialog, this.snackbarService, {
            action: "remove",
            claimName: claim.name,
            requested: [receipt],
            result,
          });
          this.refresh();
        },
        error: () => this.refresh(),
      });
  }

  public exportReceipts(): void {
    const claim = this.claim();
    if (claim && claim.receiptIds.length > 0) {
      this.receiptExportService.exportReceiptsById(claim.receiptIds);
    }
  }

  /**
   * Downloads the claim's packet. Failures arrive as one message, already read
   * out of the Blob error body; the interceptor stays quiet for this request.
   * Not cancelled on leaving the page: the user asked for the file.
   */
  public downloadPacket(): void {
    const claim = this.claim();
    if (!claim || this.packetBusy()) {
      return;
    }

    this.packetBusy.set(true);
    this.claimPacketDownloadService
      .download(claim)
      .pipe(
        take(1),
        finalize(() => this.packetBusy.set(false))
      )
      .subscribe({
        next: (filename) => this.snackbarService.success(`Claim packet downloaded: ${filename}`),
        error: (error: Error) => this.snackbarService.error(error.message),
      });
  }

  public addLedgerEntry(): void {
    this.openLedgerEntryDialog();
  }

  public editLedgerEntry(entry: ClaimLedgerEntry): void {
    this.openLedgerEntryDialog(entry);
  }

  private openLedgerEntryDialog(entry?: ClaimLedgerEntry): void {
    const claim = this.claim();
    if (!claim) {
      return;
    }

    this.matDialog
      .open<ClaimLedgerEntryDialogComponent, ClaimLedgerEntryDialogData, boolean>(
        ClaimLedgerEntryDialogComponent,
        { ...DEFAULT_DIALOG_CONFIG, data: { claimId: claim.id, entry } }
      )
      .afterClosed()
      .pipe(take(1), filter(Boolean))
      .subscribe(() => this.refresh());
  }

  public deleteLedgerEntry(entry: ClaimLedgerEntry): void {
    const claim = this.claim();
    if (!claim) {
      return;
    }

    const dialogRef = this.matDialog.open(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);
    dialogRef.componentInstance.headerText = "Delete Ledger Entry";
    dialogRef.componentInstance.dialogContent =
      `Are you sure you would like to delete this ${this.ledgerTypeOption(entry.type).displayValue.toLowerCase()} entry? This action is irreversible.`;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        filter(Boolean),
        switchMap(() => this.claimService.deleteClaimLedgerEntry(claim.id, entry.id))
      )
      .subscribe(() => {
        this.snackbarService.success("Ledger entry deleted");
        this.refresh();
      });
  }

  public deleteClaim(): void {
    const claim = this.claim();
    if (!claim) {
      return;
    }

    const dialogRef = this.matDialog.open(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);
    dialogRef.componentInstance.headerText = "Delete Claim";
    dialogRef.componentInstance.dialogContent =
      `Are you sure you would like to delete the claim ${claim.name}? Its ledger goes with it; its receipts are kept. This action is irreversible.`;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        filter(Boolean),
        switchMap(() => this.claimService.deleteClaim(claim.id))
      )
      .subscribe(() => {
        this.snackbarService.success("Claim deleted");
        this.navigateToList();
      });
  }

  private navigateToList(): void {
    this.router.navigate(["/claims/group", this.store.selectSnapshot(GroupState.selectedGroupId)]);
  }
}
