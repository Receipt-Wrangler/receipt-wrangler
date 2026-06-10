import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { GroupRole, PermissionScope, Role, RoleService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { AuthState } from "../../store";
import { GroupMemberFormComponent } from "./group-member-form.component";

describe("GroupMemberFormComponent", () => {
  const defaultGroupRole: Role = {
    id: 10,
    name: "Legacy Owner",
    description: "Group owner",
    scope: PermissionScope.Group,
    isDefault: true,
    isSystem: true,
    permissions: [],
  };

  // An app role is included to confirm the selector filters to group roles only.
  const appRole: Role = {
    id: 1,
    name: "Legacy Admin",
    description: "App admin",
    scope: PermissionScope.App,
    isDefault: true,
    isSystem: true,
    permissions: [],
  };

  let component: GroupMemberFormComponent;
  let fixture: ComponentFixture<GroupMemberFormComponent>;
  let getRolesMock: jest.Mock;
  let dialogRefMock: { close: jest.Mock };

  beforeEach(async () => {
    getRolesMock = jest.fn().mockReturnValue(of([] as Role[]));
    dialogRefMock = { close: jest.fn() };
    await TestBed.configureTestingModule({
      declarations: [GroupMemberFormComponent],
      imports: [
        NgxsModule.forRoot([AuthState]),
        PipesModule,
        ReactiveFormsModule,
        MatDialogModule,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        provideZonelessChangeDetection(),
        { provide: MatDialogRef, useValue: dialogRefMock },
        { provide: RoleService, useValue: { getRoles: getRolesMock } },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    TestBed.inject(Store).reset({ auth: { userId: "1" } });
  });

  // Creates and initialises a fresh component after the role mock is configured;
  // the roles toSignal subscribes at construction, so the mock must be set first.
  async function createComponent(): Promise<void> {
    fixture = TestBed.createComponent(GroupMemberFormComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  }

  it("should create", async () => {
    await createComponent();
    expect(component).toBeTruthy();
  });

  it("defaults to the configured default group role on add", async () => {
    getRolesMock.mockReturnValue(of([appRole, defaultGroupRole]));
    await createComponent();

    expect(component.groupRoleOptions().map((role) => role.id)).toEqual([10]);
    expect(component.form.get("groupRoleId")?.value).toBe(10);
    expect(component.selectedRole()).toEqual(defaultGroupRole);
  });

  it("leaves the selector empty when roles fail to load", async () => {
    getRolesMock.mockReturnValue(throwError(() => new Error("forbidden")));
    await createComponent();

    expect(component).toBeTruthy();
    expect(component.groupRoleOptions()).toEqual([]);
    expect(component.form.get("groupRoleId")?.value).toBeNull();
  });

  it("opens the preview dialog for the selected role", async () => {
    getRolesMock.mockReturnValue(of([defaultGroupRole]));
    await createComponent();

    const openSpy = jest
      .spyOn(TestBed.inject(MatDialog), "open")
      .mockReturnValue({ afterClosed: () => of(undefined) } as any);

    component.previewRole();

    expect(openSpy).toHaveBeenCalledTimes(1);
    expect(openSpy.mock.calls[0][1]?.data?.role).toEqual(defaultGroupRole);
  });

  it("syncs the legacy group role and returns the form on submit", async () => {
    getRolesMock.mockReturnValue(of([defaultGroupRole]));
    await createComponent();

    component.form.get("userId")?.setValue("2");
    component.form.get("groupRoleId")?.setValue(10);

    component.submit();

    // Legacy Owner → OWNER, carried so the parent group form's legacy readers work.
    expect(component.form.get("groupRole")?.value).toBe(GroupRole.Owner);
    expect(dialogRefMock.close).toHaveBeenCalledWith(component.form);
  });
});
