import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { PermissionScope, Role, RoleService } from "../open-api";
import { AuthState } from "../store";
import { SetPermissions } from "../store/auth.state.actions";
import { loadAssignableRoles } from "./role-loading.util";

describe("loadAssignableRoles", () => {
  let store: Store;
  let getRoles: jest.Mock;
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
    getRoles = jest.fn().mockReturnValue(of(roles));
    roleService = { getRoles } as unknown as RoleService;
  });

  function collect(done: (value: Role[]) => void): void {
    loadAssignableRoles(store, roleService).subscribe(done);
  }

  it("does NOT call getRoles when app.roles.read is absent", (done) => {
    collect((result) => {
      expect(result).toEqual([]);
      expect(getRoles).not.toHaveBeenCalled();
      done();
    });
  });

  it("loads roles when app.roles.read is held", (done) => {
    store.dispatch(new SetPermissions(["app.roles.read"], {}));
    collect((result) => {
      expect(result).toEqual(roles);
      expect(getRoles).toHaveBeenCalledTimes(1);
      done();
    });
  });

  it("degrades to an empty list when getRoles errors", (done) => {
    store.dispatch(new SetPermissions(["app.roles.read"], {}));
    getRoles.mockReturnValue(throwError(() => new Error("boom")));
    collect((result) => {
      expect(result).toEqual([]);
      done();
    });
  });
});
