import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl, FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatCardModule } from "@angular/material/card";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatSort } from "@angular/material/sort";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute, Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { SelectModule } from "src/select/select.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { UserAutocompleteModule } from "src/user-autocomplete/user-autocomplete.module";
import { ButtonModule } from "../../button";
import { InputModule } from "../../input";
import { ApiModule, Group, GroupRole, GroupsService, GroupStatus } from "../../open-api";
import { AddGroup, UpdateGroup } from "../../store";
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
        PipesModule,
        InputModule,
        MatCardModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([]),
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule,
        RouterTestingModule.withRoutes([]),
        SelectModule,
        SharedUiModule,
        TableModule,
        UserAutocompleteModule],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        group: {},
                        formConfig: {},
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
    component.table = {
      sort: new MatSort(),
    } as any;
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
      groupRole: GroupRole.Viewer,
    });
    jest.spyOn(matDialog, "open").mockReturnValue({
      afterClosed: () => of(formGroup),
      componentInstance: { currentGroupMembers: [] },
    } as any);
    component.ngOnInit();
    component.ngAfterViewInit();

    component.addGroupMemberClicked();

    expect(component.groupMembers.value).toEqual([formGroup.value]);
    expect(component.dataSource.data).toEqual([formGroup.value]);
  });

  it("should update form on edit", () => {
    const result = [
      {
        userId: 2,
        groupId: 1,
        groupRole: GroupRole.Owner,
      },
      {
        userId: 3,
        groupId: 1,
        groupRole: GroupRole.Owner,
      },
    ];
    const matDialog = TestBed.inject(MatDialog);
    const formGroup = buildGroupMemberForm();
    formGroup.patchValue({
      userId: 3,
      groupId: 1,
      groupRole: GroupRole.Owner,
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
            groupRole: GroupRole.Owner,
          },
          {
            userId: 1,
            groupId: 1,
            groupRole: GroupRole.Viewer,
          },
        ],
      },
      formConfig: {},
    };

    component.ngOnInit();
    component.ngAfterViewInit();

    component.editGroupMemberClicked(1);

    expect(component.groupMembers.value).toEqual(result);
    expect(component.dataSource.data).toEqual(result as any);
  });

  it("should remove user to group member on remove", () => {
    const route = TestBed.inject(ActivatedRoute);
    const result = {
      userId: 1,
      groupId: 1,
      groupRole: GroupRole.Viewer,
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
            groupRole: GroupRole.Owner,
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
    expect(component.dataSource.data).toEqual([result] as any);
  });

  it("should create group", () => {
    const createSpy = jest.spyOn(TestBed.inject(GroupsService), "createGroup");
    const storeSpy = jest.spyOn(TestBed.inject(Store), "dispatch");
    const routerSpy = jest.spyOn(TestBed.inject(Router), "navigate");

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
    const routerSpy = jest.spyOn(TestBed.inject(Router), "navigate");

    const group: Group = {
      id: 1,
      name: "test",
      isDefault: true,
      status: GroupStatus.Active,
      groupMembers: [
        {
          userId: 2,
          groupId: 1,
          groupRole: GroupRole.Owner,
        },
        {
          userId: 1,
          groupId: 1,
          groupRole: GroupRole.Viewer,
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
        groupId: new FormControl(1),
        groupRole: new FormControl(GroupRole.Editor),
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
            groupId: 1,
            groupRole: GroupRole.Owner,
          },
          {
            userId: 1,
            groupId: 1,
            groupRole: GroupRole.Viewer,
          },
          {
            userId: 3,
            groupId: 1,
            groupRole: GroupRole.Editor,
          },
        ],
      } as Group
    );
    expect(storeSpy).toHaveBeenCalledWith(new UpdateGroup(returnValue));
    expect(routerSpy).toHaveBeenCalledWith(["/groups/1/details/view"], {
      queryParams: {
        tab: "details",
      }
    });
  });
});
