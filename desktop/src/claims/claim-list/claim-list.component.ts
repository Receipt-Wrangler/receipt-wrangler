import { AfterViewInit, Component, computed, inject, signal, TemplateRef, viewChild } from "@angular/core";
import { rxResource, toSignal } from "@angular/core/rxjs-interop";
import { MatTableDataSource } from "@angular/material/table";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { map } from "rxjs";
import { DEFAULT_HOST_CLASS } from "../../constants";
import { Claim, ReimbursementClaimService } from "../../open-api";
import { AuthState, GroupState } from "../../store";
import { TableColumn } from "../../table/table-column.interface";
import { canUpdateClaims, claimStatusLabel, claimStatusTone } from "../claim.utils";

/**
 * A group's claims, newest first as the server returns them, under the figures
 * that answer "what am I still owed, and how long has it been waiting?". Every
 * figure is the server's; nothing here adds anything up.
 *
 * On the synthetic All group the server lists every group the viewer can read,
 * so a Group column appears there. "Add Claim" is offered only when there is a
 * real group to add it to — the same add target the receipt form seeds — and the
 * viewer may update claims in it.
 */
@Component({
  selector: "app-claim-list",
  templateUrl: "./claim-list.component.html",
  styleUrls: ["./claim-list.component.scss"],
  host: DEFAULT_HOST_CLASS,
  standalone: false,
})
export class ClaimListComponent implements AfterViewInit {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly claimService = inject(ReimbursementClaimService);
  private readonly store = inject(Store);

  private readonly nameCell = viewChild.required<TemplateRef<any>>("nameCell");
  private readonly groupCell = viewChild.required<TemplateRef<any>>("groupCell");
  private readonly statusCell = viewChild.required<TemplateRef<any>>("statusCell");
  private readonly submittedDateCell = viewChild.required<TemplateRef<any>>("submittedDateCell");
  private readonly dueDateCell = viewChild.required<TemplateRef<any>>("dueDateCell");
  private readonly claimedCell = viewChild.required<TemplateRef<any>>("claimedCell");
  private readonly paidCell = viewChild.required<TemplateRef<any>>("paidCell");
  private readonly outstandingCell = viewChild.required<TemplateRef<any>>("outstandingCell");
  private readonly daysWaitingCell = viewChild.required<TemplateRef<any>>("daysWaitingCell");

  /** The route's group, kept live so a param change on a reused route refetches. */
  public readonly groupId = toSignal(
    this.activatedRoute.paramMap.pipe(map((params) => Number(params.get("groupId")))),
    { initialValue: Number(this.activatedRoute.snapshot.paramMap.get("groupId")) }
  );

  private readonly groups = this.store.selectSignal(GroupState.groups);

  private readonly groupPermissions = this.store.selectSignal(AuthState.groupPermissions);

  private readonly addTargetGroupId = this.store.selectSignal(GroupState.addTargetGroupId);

  public readonly group = computed(() => this.groups().find((group) => group.id === this.groupId()));

  public readonly isAllGroup = computed(() => !!this.group()?.isAllGroup);

  public readonly headerText = computed(() => {
    const name = this.group()?.name;
    if (!name) {
      return "Claims";
    }

    return name.toLowerCase().includes("claim") ? name : `${name} Claims`;
  });

  /** A new claim needs a real group: never the All group itself. */
  public readonly canAdd = computed(() =>
    canUpdateClaims(this.groupPermissions(), this.addTargetGroupId())
  );

  public readonly subtitle = computed(() => {
    if (!this.isAllGroup()) {
      return "What was submitted, what has been paid, and what is still owed.";
    }

    return this.canAdd()
      ? "Claims across every group you can read."
      : "Claims across every group you can read. Open a group to add a claim to it.";
  });

  public readonly claims = rxResource({
    params: () => ({ groupId: this.groupId() }),
    stream: ({ params }) => this.claimService.getClaimsForGroup(params.groupId),
  });

  public readonly summary = computed(() => (this.claims.hasValue() ? this.claims.value().summary : undefined));

  public readonly dataSource = computed(
    () => new MatTableDataSource<Claim>(this.claims.hasValue() ? this.claims.value().claims : [])
  );

  private readonly viewReady = signal(false);

  public readonly columns = computed<TableColumn[]>(() => {
    if (!this.viewReady()) {
      return [];
    }

    return [
      { columnHeader: "Name", matColumnDef: "name", template: this.nameCell(), sortable: false },
      ...(this.isAllGroup()
        ? [{ columnHeader: "Group", matColumnDef: "group", template: this.groupCell(), sortable: false }]
        : []),
      { columnHeader: "Reference", matColumnDef: "referenceNumber", elementKey: "referenceNumber", sortable: false },
      { columnHeader: "Type", matColumnDef: "type", elementKey: "type", sortable: false },
      { columnHeader: "Status", matColumnDef: "status", template: this.statusCell(), sortable: false },
      { columnHeader: "Submitted", matColumnDef: "submittedDate", template: this.submittedDateCell(), sortable: false },
      { columnHeader: "Due", matColumnDef: "dueDate", template: this.dueDateCell(), sortable: false },
      { columnHeader: "Claimed", matColumnDef: "claimed", template: this.claimedCell(), sortable: false },
      { columnHeader: "Paid", matColumnDef: "paid", template: this.paidCell(), sortable: false },
      { columnHeader: "Outstanding", matColumnDef: "outstanding", template: this.outstandingCell(), sortable: false },
      { columnHeader: "Days waiting", matColumnDef: "daysWaiting", template: this.daysWaitingCell(), sortable: false },
    ];
  });

  public readonly displayedColumns = computed(() => this.columns().map((column) => column.matColumnDef));

  public readonly statusLabel = claimStatusLabel;

  public readonly statusTone = claimStatusTone;

  public ngAfterViewInit(): void {
    this.viewReady.set(true);
  }

  public isNegative(amount: string | undefined): boolean {
    return !!amount && amount.trim().startsWith("-");
  }
}
