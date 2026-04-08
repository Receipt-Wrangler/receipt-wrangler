import { Component, OnInit, TemplateRef, viewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { AssociatedEntityType, Prompt, ReceiptProcessingSettings } from "../../open-api";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { SystemTaskTableService } from "../../services/system-task-table.service";
import { SystemTaskFilterComponent } from "../../shared-ui/system-task-filter/system-task-filter.component";
import { TaskTableComponent } from "../../shared-ui/task-table/task-table.component";
import { SystemTaskTableState } from "../../store/system-task-table.state";
import { ResetSystemTaskFilter, SetPage } from "../../store/system-task-table.state.actions";
import { applyFormCommand } from "../../utils/index";
import { buildSystemTaskFilterForm } from "../../utils/system-task-filter";

@UntilDestroy()
@Component({
  selector: "app-system-task-table",
  templateUrl: "./system-task-table.component.html",
  styleUrl: "./system-task-table.component.scss",
  providers: [
    {
      provide: TABLE_SERVICE_INJECTION_TOKEN,
      useClass: SystemTaskTableService
    },
  ],
  standalone: false
})
export class SystemTaskTableComponent implements OnInit {
  public readonly expandedRowTemplate = viewChild.required<TemplateRef<any>>("expandedRowTemplate");
  public readonly taskTableComponent = viewChild.required(TaskTableComponent);

  public prompts: Prompt[] = [];
  public allReceiptProcessingSettings: ReceiptProcessingSettings[] = [];
  protected readonly AssociatedEntityType = AssociatedEntityType;

  public numFiltersApplied = this.store.selectSignal(SystemTaskTableState.numFiltersApplied);

  constructor(
    private activatedRoute: ActivatedRoute,
    private matDialog: MatDialog,
    private store: Store,
  ) {}

  public ngOnInit(): void {
    this.prompts = this.activatedRoute.snapshot.data["prompts"] || [];
    this.allReceiptProcessingSettings = this.activatedRoute.snapshot.data["allReceiptProcessingSettings"] || [];
  }

  public refresh(): void {
    this.taskTableComponent().getTableData();
  }

  public filterButtonClicked(): void {
    const filter = this.store.selectSnapshot(SystemTaskTableState.filterData).filter as any;

    const dialogRef = this.matDialog.open(SystemTaskFilterComponent, {
      minWidth: "75%",
      maxWidth: "100%",
    });

    dialogRef.componentInstance.parentForm = buildSystemTaskFilterForm(filter, this);
    dialogRef.componentInstance.headerText = "Filter System Tasks";
    const formCommandSubscription = dialogRef.componentInstance.formCommand.subscribe((formCommand) => {
      applyFormCommand(dialogRef.componentInstance.parentForm, formCommand);
    });

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((applyFilter) => {
          if (applyFilter) {
            this.store.dispatch(new SetPage(1));
            this.taskTableComponent().getTableData();
          }

          formCommandSubscription.unsubscribe();
        })
      )
      .subscribe();
  }

  public resetFilterButtonClicked(): void {
    this.store.dispatch(new ResetSystemTaskFilter());
    this.taskTableComponent().getTableData();
  }
}
