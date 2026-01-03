import {AfterViewInit, Component, Input, OnInit, TemplateRef, ViewChild} from "@angular/core";
import {FormArray, FormBuilder, FormGroup, Validators} from "@angular/forms";
import {MatDialog} from "@angular/material/dialog";
import {Sort} from "@angular/material/sort";
import {MatTableDataSource} from "@angular/material/table";
import {ActivatedRoute, Router} from "@angular/router";
import {UntilDestroy, untilDestroyed} from "@ngneat/until-destroy";
import {Store} from "@ngxs/store";
import {startWith, take, tap} from "rxjs";
import {DEFAULT_HOST_CLASS} from "src/constants";
import {GROUP_STATUS_OPTIONS} from "src/constants/receipt-status-options";
import {FormMode} from "src/enums/form-mode.enum";
import {FormConfig} from "src/interfaces/form-config.interface";
import {TableColumn} from "src/table/table-column.interface";
import {TableComponent} from "src/table/table/table.component";
import {SortByDisplayName} from "src/utils/sort-by-displayname";
import {Group, GroupMember, GroupRole, GroupsService, GroupStatus} from "../../open-api";
import {SnackbarService} from "../../services";
import {AddGroup, UpdateGroup} from "../../store";
import {GroupMemberFormComponent} from "../group-member-form/group-member-form.component";
import {buildGroupMemberForm} from "../utils/group-member.utils";

@UntilDestroy()
@Component({
    selector: "app-create-group-form",
    templateUrl: "./group-form.component.html",
    styleUrls: ["./group-form.component.scss"],
    host: DEFAULT_HOST_CLASS,
    standalone: false
})
export class GroupFormComponent implements OnInit, AfterViewInit {
  @ViewChild("nameCell") public nameCell!: TemplateRef<any>;

  @ViewChild("groupRoleCell") public groupRoleCell!: TemplateRef<any>;

  @ViewChild("actionsCell") public actionsCell!: TemplateRef<any>;

  @ViewChild(TableComponent) public table!: TableComponent;

  @Input() public canEdit = true;

  public form: FormGroup = new FormGroup({});

  public get groupMembers(): FormArray {
    return this.form.get("groupMembers") as FormArray;
  }

  public formConfig!: FormConfig;

  public originalGroup: Group | undefined = undefined;

  public displayedColumns: string[] = [];

  public columns: TableColumn[] = [];

  public disableDeleteButton: boolean = false;

  public editLink: string = "";

  public groupRole = GroupRole;

  public groupStatusOptions = GROUP_STATUS_OPTIONS;

  public dataSource: MatTableDataSource<GroupMember> =
    new MatTableDataSource<GroupMember>([]);

  constructor(
    private formBuilder: FormBuilder,
    private groupsService: GroupsService,
    private snackbarService: SnackbarService,
    private router: Router,
    private store: Store,
    private activatedRoute: ActivatedRoute,
    private matDialog: MatDialog,
    private sortByDisplayName: SortByDisplayName
  ) {
  }

  public ngOnInit(): void {
    this.originalGroup = this.activatedRoute.snapshot.data["group"];
    this.formConfig = this.activatedRoute.snapshot.data["formConfig"];
    if (this.originalGroup) {
      this.editLink = `/groups/${this.originalGroup.id}/details/edit`;
    }
    this.initForm();
  }

  public ngAfterViewInit(): void {
    this.initTable();
  }

  private initTable(): void {
    this.setColumns();
    this.setDataSource();
    this.listenForGroupMemberChanges();
  }

  private listenForGroupMemberChanges(): void {
    this.groupMembers.valueChanges
      .pipe(
        untilDestroyed(this),
        tap((v) => {
          this.dataSource.data = Array.from(v);
        })
      )
      .subscribe();
  }

  private setColumns(): void {
    this.columns = [
      {
        columnHeader: "Name",
        matColumnDef: "name",
        template: this.nameCell,
        sortable: true,
      },
      {
        columnHeader: "Group Role",
        matColumnDef: "groupRole",
        template: this.groupRoleCell,
        sortable: true,
      },
      {
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell,
        sortable: true,
      },
    ];
    this.displayedColumns = ["name", "groupRole"];

    if (this.formConfig.mode !== FormMode.view) {
      this.displayedColumns.push("actions");
    }
  }

  private setDataSource(): void {
    this.dataSource = new MatTableDataSource<GroupMember>(
      this.groupMembers.value ?? []
    );
    this.dataSource.sort = this.table.sort;
  }

  public sortName(sortState: Sort): void {
    if (sortState.active === "name") {
      if (sortState.direction === "") {
        this.dataSource.data = this.groupMembers.value;
        return;
      }

      const newData = this.sortByDisplayName.sort(
        this.dataSource.data,
        sortState,
        "userId"
      );

      this.dataSource.data = newData;
    }
  }

  private initForm(): void {
    let groupMembers: FormGroup[] = [];
    if (this.originalGroup?.groupMembers) {
      groupMembers = this.originalGroup.groupMembers.map((m) =>
        buildGroupMemberForm(m)
      );
    }
    this.form = this.formBuilder.group({
      name: [this.originalGroup?.name ?? "", Validators.required],
      groupMembers: this.formBuilder.array(groupMembers),
      status: this.originalGroup?.status ?? GroupStatus.Active,
    });

    this.groupMembers.valueChanges
      .pipe(
        untilDestroyed(this),
        startWith(this.groupMembers.value),
        tap((v) => {
          this.disableDeleteButton = v.length === 1;
        })
      )
      .subscribe();

    if (this.formConfig.mode === FormMode.view) {
      this.form.get("status")?.disable();
    }
  }

  public addGroupMemberClicked(): void {
    const dialogRef = this.matDialog.open(GroupMemberFormComponent);

    dialogRef.componentInstance.currentGroupMembers = this.groupMembers.value;
    dialogRef.componentInstance.headerText = "Add Group Member";

    dialogRef
      .afterClosed()
      .pipe(take(1))
      .subscribe((form) => {
        if (form) {
          this.groupMembers.push(form);
        }
      });
  }

  public editGroupMemberClicked(index: number): void {
    const groupMember = this.groupMembers.at(index).value;
    const dialogRef = this.matDialog.open(GroupMemberFormComponent);

    dialogRef.componentInstance.currentGroupMembers = this.groupMembers.value;
    dialogRef.componentInstance.groupMember = groupMember;
    dialogRef.componentInstance.headerText = "Edit Group Member";

    dialogRef
      .afterClosed()
      .pipe(take(1))
      .subscribe((form) => {
        if (form) {
          this.groupMembers.at(index).patchValue(form.value);
        }
      });
  }

  public removeGroupMember(index: number): void {
    this.groupMembers.removeAt(index);
  }

  public submit(): void {
    if (this.form.valid) {
      const owners = (this.groupMembers.value as GroupMember[]).filter(
        (gm) => gm.groupRole === GroupRole.Owner
      );
      if (owners.length === 0 && this.formConfig.mode !== FormMode.add) {
        this.snackbarService.error("Group must have at least one owner!");
        return;
      }
      switch (this.formConfig.mode) {
        case FormMode.add:
          this.createGroup();

          break;
        case FormMode.edit:
          this.updateGroup();
          break;

        default:
          break;
      }
    }
  }

  private createGroup(): void {
    this.groupsService
      .createGroup(this.form.value)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Group successfully created");
        }),
        tap((group: Group) => {
          this.store.dispatch(new AddGroup(group));
          this.navigateToGroupDetails(group.id);
        })
      )
      .subscribe();
  }

  private updateGroup(): void {
    this.groupsService
      .updateGroup(this.originalGroup?.id ?? 0, this.form.value)
      .pipe(
        take(1),
        tap((group: Group) => {
          this.snackbarService.success("Group successfully updated");
          this.store.dispatch(new UpdateGroup(group));
          this.navigateToGroupDetails(group.id);
        })
      )
      .subscribe();
  }

  private navigateToGroupDetails(groupId: number): void {
    this.router.navigate([`/groups/${groupId}/details/view`],
      {
        queryParams: {tab: "details"}
      });
  }
}
