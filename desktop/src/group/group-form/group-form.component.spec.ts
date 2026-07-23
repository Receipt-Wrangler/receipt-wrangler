import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl, FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatCardModule } from "@angular/material/card";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatSort } from "@angular/material/sort";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute, Router, provideRouter } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { SelectModule } from "src/select/select.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { UserAutocompleteModule } from "src/user-autocomplete/user-autocomplete.module";
import { CheckboxModule } from "../../checkbox/checkbox.module";
import { ButtonModule } from "../../button";
import { InputModule } from "../../input";
import { ApiModule, Group, GroupsService, GroupStatus, Permission } from "../../open-api";
import { AddGroup, AuthState, UpdateGroup } from "../../store";
import { SetPermissions } from "../../store/auth.state.actions";
import { AppInitService } from "../../services";
import { GroupMemberFormComponent } from "../group-member-form/group-member-form.component";
import { buildGroupMemberForm } from "../utils/group-member.utils";
import { GroupFormComponent } from "./group-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("GroupFormComponent", () => {
  let component: GroupFormComponent;
  let fixture: ComponentFixture<GroupFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [GroupFormComponent, GroupMemberFormComponent],
    imports: [ApiModule,
        ButtonModule,
        CheckboxModule,
        PipesModule,
        InputModule,
        MatCardModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState]),
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule,
        SelectModule,
        SharedUiModule,
        TableModule,
        UserAutocompleteModule],
    providers: [
        provideRouter([]),
        {
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        group: {},
                        formConfig: { mode: FormMode.add },
                    },
                },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    fixture = TestBed.createComponent(GroupFormComponent);
    component = fixture.componentInstance;
    Object.defineProperty(component, 'table', {
      value: () => ({
        sort: () => new MatSort(),
      }),
      configurable: true,
    });
    // fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should add user to group member on add", () => {
    const matDialog = TestBed.inject(MatDialog);
    const formGroup = buildGroupMemberForm();
    formGroup.patchValue({
      userId: "2",
      groupId: "1",
      groupRoleId: 10,
    });
    jest.spyOn(matDialog, "open").mockReturnValue({
      afterClosed: () => of(formGroup),
      componentInstance: { currentGroupMembers: [] },
    } as any);
    component.ngOnInit();
    component.ngAfterViewInit();

    component.addGroupMemberClicked();

    expect(component.groupMembers.value).toEqual([formGroup.value]);
    expect(component.dataSource().data).toEqual([formGroup.value]);
  });

  it("should update form on edit", () => {
    const result = [
      {
        userId: 2,
        groupRoleId: null,
        groupId: 1,
      },
      {
        userId: 3,
        groupRoleId: null,
        groupId: 1,
      },
    ];
    const matDialog = TestBed.inject(MatDialog);
    const formGroup = buildGroupMemberForm();
    formGroup.patchValue({
      userId: 3,
      groupId: 1,
    });
    jest.spyOn(matDialog, "open").mockReturnValue({
      afterClosed: () => of(formGroup),
      componentInstance: { currentGroupMembers: [] },
    } as any);
    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      group: {
        id: 1,
        name: "test",
        isDefault: true,
        groupMembers: [
          {
            userId: 2,
            groupId: 1,
          },
          {
            userId: 1,
            groupId: 1,
          },
        ],
      },
      formConfig: {},
    };

    component.ngOnInit();
    component.ngAfterViewInit();

    component.editGroupMemberClicked(1);

    expect(component.groupMembers.value).toEqual(result);
    expect(component.dataSource().data).toEqual(result as any);
  });

  it("should remove user to group member on remove", () => {
    const route = TestBed.inject(ActivatedRoute);
    const result = {
      userId: 1,
      groupRoleId: null,
      groupId: 1,
    };

    route.snapshot.data = {
      group: {
        id: 1,
        name: "test",
        isDefault: true,
        groupMembers: [
          {
            userId: 2,
            groupId: 1,
          },
          result,
        ],
      },
      formConfig: {},
    };

    component.ngOnInit();
    component.ngAfterViewInit();

    component.removeGroupMember(0);

    expect(component.groupMembers.value).toEqual([result]);
    expect(component.dataSource().data).toEqual([result] as any);
  });

  it("should create group", () => {
    const createSpy = jest.spyOn(TestBed.inject(GroupsService), "createGroup");
    const storeSpy = jest.spyOn(TestBed.inject(Store), "dispatch");
    const routerSpy = jest.spyOn(TestBed.inject(Router), "navigate").mockResolvedValue(true);
    jest.spyOn(TestBed.inject(AppInitService), "getAppData").mockReturnValue(of([]));

    const group: Group = {
      id: 1,
      name: "test",
      isDefault: true,
      groupMembers: [],
      status: GroupStatus.Active,
      isAllGroup: false,
      groupReceiptSettings: {} as any,
    };

    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      formConfig: {
        mode: FormMode.add,
      },
    };

    component.ngOnInit();
    component.ngAfterViewInit();

    component.form.patchValue({
      name: group.name,
      isDefault: group.isDefault,
    });

    const returnValue = {
      ...component.form.value,
      id: 1,
    };

    createSpy.mockReturnValue(of(returnValue));

    component.submit();

    expect(createSpy).toHaveBeenCalledWith({
      name: "test",
      status: GroupStatus.Active,
      groupMembers: [],
      isolateMembers: false,
    } as any);
    expect(storeSpy).toHaveBeenCalledWith(new AddGroup(returnValue));
    expect(routerSpy).toHaveBeenCalledWith(["/groups/1/details/view"], {
      queryParams: {
        tab: "details",
      }
    });
  });

  it("should update group", () => {
    const updateSpy = jest.spyOn(TestBed.inject(GroupsService), "updateGroup");
    const storeSpy = jest.spyOn(TestBed.inject(Store), "dispatch");
    const routerSpy = jest.spyOn(TestBed.inject(Router), "navigate").mockResolvedValue(true);
    jest.spyOn(TestBed.inject(AppInitService), "getAppData").mockReturnValue(of([]));

    const group: Group = {
      id: 1,
      name: "test",
      isDefault: true,
      status: GroupStatus.Active,
      groupMembers: [
        {
          userId: 2,
          groupRoleId: 10,
          groupId: 1,
        },
        {
          userId: 1,
          groupRoleId: 11,
          groupId: 1,
        },
      ],
      isAllGroup: false,
      groupReceiptSettings: {} as any,
    };

    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      group: group,
      formConfig: {
        mode: FormMode.edit,
      },
    };

    component.ngOnInit();
    component.ngAfterViewInit();

    component.form.patchValue({
      name: "new name",
    });

    component.groupMembers.push(
      new FormGroup({
        userId: new FormControl(3),
        groupRoleId: new FormControl(12),
        groupId: new FormControl(1),
      })
    );

    const returnValue = {
      ...component.form.value,
      id: 1,
    };

    updateSpy.mockReturnValue(of(returnValue));

    component.submit();

    expect(updateSpy).toHaveBeenCalledWith(
      component.originalGroup?.id as number,
      {
        name: "new name",
        status: GroupStatus.Active,
        groupMembers: [
          {
            userId: 2,
            groupRoleId: 10,
            groupId: 1,
          },
          {
            userId: 1,
            groupRoleId: 11,
            groupId: 1,
          },
          {
            userId: 3,
            groupRoleId: 12,
            groupId: 1,
          },
        ],
        isolateMembers: false,
      } as Group
    );
    expect(storeSpy).toHaveBeenCalledWith(new UpdateGroup(returnValue));
    expect(routerSpy).toHaveBeenCalledWith(["/groups/1/details/view"], {
      queryParams: {
        tab: "details",
      }
    });
  });

  it("includes isolateMembers in the create payload when enabled", () => {
    const createSpy = jest.spyOn(TestBed.inject(GroupsService), "createGroup");
    jest.spyOn(TestBed.inject(Router), "navigate").mockResolvedValue(true);
    jest.spyOn(TestBed.inject(AppInitService), "getAppData").mockReturnValue(of([]));

    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      formConfig: {
        mode: FormMode.add,
      },
    };

    component.ngOnInit();
    component.ngAfterViewInit();

    component.form.patchValue({
      name: "test",
      isolateMembers: true,
    });

    createSpy.mockReturnValue(of({ ...component.form.value, id: 1 }));

    component.submit();

    expect(createSpy).toHaveBeenCalledWith(
      expect.objectContaining({ isolateMembers: true })
    );
  });

  it("hydrates isolateMembers from the group on edit", () => {
    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      group: {
        id: 1,
        name: "test",
        status: GroupStatus.Active,
        isolateMembers: true,
        groupMembers: [],
      },
      formConfig: { mode: FormMode.edit },
    };

    component.ngOnInit();

    expect(component.form.get("isolateMembers")?.value).toBe(true);
  });

  it("submits in edit mode even when no member holds an owner role", () => {
    // The legacy "must have at least one owner" guard was removed with the
    // group-role enums; the backend no longer enforces an owner, so the client
    // must save a non-owner-only membership without blocking.
    const updateSpy = jest.spyOn(TestBed.inject(GroupsService), "updateGroup");
    jest.spyOn(TestBed.inject(Router), "navigate").mockResolvedValue(true);
    jest.spyOn(TestBed.inject(AppInitService), "getAppData").mockReturnValue(of([]));

    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      group: {
        id: 1,
        name: "test",
        status: GroupStatus.Active,
        groupMembers: [{ userId: 2, groupRoleId: 11, groupId: 1 }],
      },
      formConfig: { mode: FormMode.edit },
    };

    component.ngOnInit();
    component.ngAfterViewInit();

    updateSpy.mockReturnValue(of({ id: 1 } as Group));

    component.submit();

    expect(updateSpy).toHaveBeenCalled();
  });

  it("keeps member controls enabled in create mode (the creator owns the new group)", () => {
    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = { formConfig: { mode: FormMode.add } }; // no group => create

    component.ngOnInit();

    expect(component.canCreateGroupMembers()).toBe(true);
    expect(component.canUpdateGroupMembers()).toBe(true);
    expect(component.canDeleteGroupMembers()).toBe(true);
  });

  it("hides member controls in edit mode without group.members.* permissions", () => {
    TestBed.inject(Store).dispatch(new SetPermissions([], {}));
    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      group: { id: 1, name: "test", status: GroupStatus.Active, groupMembers: [] },
      formConfig: { mode: FormMode.edit },
    };

    component.ngOnInit();

    expect(component.canCreateGroupMembers()).toBe(false);
    expect(component.canUpdateGroupMembers()).toBe(false);
    expect(component.canDeleteGroupMembers()).toBe(false);
  });

  it("shows member controls in edit mode for holders of group.members.*", () => {
    TestBed.inject(Store).dispatch(
      new SetPermissions([], {
        1: [
          Permission.GroupMembersCreate,
          Permission.GroupMembersUpdate,
          Permission.GroupMembersDelete,
        ],
      })
    );
    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.data = {
      group: { id: 1, name: "test", status: GroupStatus.Active, groupMembers: [] },
      formConfig: { mode: FormMode.edit },
    };

    component.ngOnInit();

    expect(component.canCreateGroupMembers()).toBe(true);
    expect(component.canUpdateGroupMembers()).toBe(true);
    expect(component.canDeleteGroupMembers()).toBe(true);
  });
});
