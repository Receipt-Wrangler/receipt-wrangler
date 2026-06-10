import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { MAT_DIALOG_DATA } from "@angular/material/dialog";
import { of, throwError } from "rxjs";
import {
  PermissionDescriptor,
  PermissionScope,
  PermissionService,
  Role,
} from "../../open-api";
import {
  RolePreviewDialogComponent,
  RolePreviewDialogData,
} from "./role-preview-dialog.component";

describe("RolePreviewDialogComponent", () => {
  const role: Role = {
    id: 1,
    name: "Auditor",
    description: "Read-only auditor",
    scope: PermissionScope.App,
    isDefault: false,
    isSystem: false,
    permissions: ["app.users.read", "app.users.create"],
  };

  const descriptors: PermissionDescriptor[] = [
    {
      key: "app.users.read" as any,
      label: "Read",
      description: "Read users",
      category: "Users",
      scope: PermissionScope.App,
    },
    {
      key: "app.users.create" as any,
      label: "Create",
      description: "Create users",
      category: "Users",
      scope: PermissionScope.App,
    },
  ];

  let fixture: ComponentFixture<RolePreviewDialogComponent>;
  let component: RolePreviewDialogComponent;
  let getPermissionsMock: jest.Mock;

  beforeEach(async () => {
    getPermissionsMock = jest.fn().mockReturnValue(of(descriptors));
    await TestBed.configureTestingModule({
      imports: [RolePreviewDialogComponent, NoopAnimationsModule],
      providers: [
        provideZonelessChangeDetection(),
        {
          provide: MAT_DIALOG_DATA,
          useValue: { role } as RolePreviewDialogData,
        },
        {
          provide: PermissionService,
          useValue: { getPermissions: getPermissionsMock },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(RolePreviewDialogComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("groups the role's permissions by resource", () => {
    const groups = component.groups();
    expect(groups.length).toBe(1);
    expect(groups[0].resourceKey).toBe("app.users");
    expect(groups[0].rows.length).toBe(2);
    expect(groups[0].rows.map((row) => row.label)).toContain("Create");
  });

  it("counts the role's permissions", () => {
    expect(component.permissionCount()).toBe(2);
  });

  it("labels an app role's scope", () => {
    expect(component.scopeLabel()).toBe("Application role");
  });

  it("renders the role name", async () => {
    await fixture.whenStable();
    const text = (fixture.nativeElement as HTMLElement).textContent ?? "";
    expect(text).toContain("Auditor");
  });

  it("still renders the role's permissions when the registry fails to load", async () => {
    // The role carries its own permission keys; only the human-readable
    // labels/descriptions come from the registry, so a failed load must
    // degrade gracefully rather than blank the dialog or throw.
    getPermissionsMock.mockReturnValue(throwError(() => new Error("boom")));

    const failedFixture = TestBed.createComponent(RolePreviewDialogComponent);
    const failedComponent = failedFixture.componentInstance;
    await failedFixture.whenStable();

    expect(failedComponent).toBeTruthy();
    expect(failedComponent.permissionCount()).toBe(2);

    const groups = failedComponent.groups();
    expect(groups.length).toBe(1);
    expect(groups[0].rows.length).toBe(2);
    // Labels fall back to the humanised action; registry-supplied descriptions
    // are absent because the registry never loaded.
    expect(groups[0].rows.map((row) => row.label)).toEqual(
      expect.arrayContaining(["Read", "Create"])
    );
    expect(groups[0].rows.every((row) => row.description === "")).toBe(true);
  });
});
