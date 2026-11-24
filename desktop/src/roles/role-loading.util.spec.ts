import type { Mock } from "vitest";
import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { firstValueFrom, of, throwError } from "rxjs";
import { PermissionScope, Role, RoleService } from "../open-api";
import { AuthState } from "../store";
import { SetPermissions } from "../store/auth.state.actions";
import { loadAssignableRoles } from "./role-loading.util";

describe("loadAssignableRoles", () => {
  let store: Store;
  let getRoles: Mock;
  let roleService: RoleService;

  const roles: Role[] = [
    {
      id: 1,
      name: "Legacy Admin",
      scope: PermissionScope.App,
      isDefault: false,
      isSystem: true,
      permissions: [],
    },
  ];

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState])],
    });
    store = TestBed.inject(Store);
    getRoles = vi.fn().mockReturnValue(of(roles));
    roleService = { getRoles } as unknown as RoleService;
  });

  function load(): Promise<Role[]> {
    return firstValueFrom(loadAssignableRoles(store, roleService));
  }

  it("does NOT call getRoles when app.roles.read is absent", async () => {
    const result = await load();

    expect(result).toEqual([]);
    expect(getRoles).not.toHaveBeenCalled();
  });

  it("loads roles when app.roles.read is held", async () => {
    store.dispatch(new SetPermissions(["app.roles.read"], {}));

    const result = await load();

    expect(result).toEqual(roles);
    expect(getRoles).toHaveBeenCalledTimes(1);
  });

  it("degrades to an empty list when getRoles errors", async () => {
    store.dispatch(new SetPermissions(["app.roles.read"], {}));
    getRoles.mockReturnValue(throwError(() => new Error("boom")));

    expect(await load()).toEqual([]);
  });
});
