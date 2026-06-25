import { PermissionScope, Role } from "../open-api";
import { RoleNamePipe } from "./role-name.pipe";

describe("RoleNamePipe", () => {
  const roles: Role[] = [
    {
      id: 1,
      name: "Legacy Admin",
      scope: PermissionScope.App,
      isDefault: false,
      isSystem: true,
      permissions: [],
    },
    {
      id: 2,
      name: "Legacy Owner",
      scope: PermissionScope.Group,
      isDefault: true,
      isSystem: true,
      permissions: [],
    },
  ];

  it("creates an instance", () => {
    expect(new RoleNamePipe()).toBeTruthy();
  });

  it("resolves the role name for a matching id and scope", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(2, roles, PermissionScope.Group)).toEqual("Legacy Owner");
  });

  it("returns empty string for an unknown id", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(99, roles, PermissionScope.App)).toEqual("");
  });

  it("returns empty string when the id exists only in the other scope", () => {
    const pipe = new RoleNamePipe();
    // id 1 is an app role; resolving it as a group role must not match.
    expect(pipe.transform(1, roles, PermissionScope.Group)).toEqual("");
  });

  it("resolves by scope when app and group role ids collide", () => {
    const pipe = new RoleNamePipe();
    // App and group roles have independent id sequences, so the same id can
    // exist in both scopes. The scope argument disambiguates them.
    const collidingRoles: Role[] = [
      {
        id: 1,
        name: "Legacy Admin",
        scope: PermissionScope.App,
        isDefault: false,
        isSystem: true,
        permissions: [],
      },
      {
        id: 1,
        name: "Legacy Viewer",
        scope: PermissionScope.Group,
        isDefault: false,
        isSystem: true,
        permissions: [],
      },
    ];
    expect(pipe.transform(1, collidingRoles, PermissionScope.Group)).toEqual("Legacy Viewer");
    expect(pipe.transform(1, collidingRoles, PermissionScope.App)).toEqual("Legacy Admin");
  });

  it("returns empty string for a null/undefined id", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(null, roles, PermissionScope.App)).toEqual("");
    expect(pipe.transform(undefined, roles, PermissionScope.App)).toEqual("");
  });

  it("returns empty string when roles are not yet loaded", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(1, [], PermissionScope.App)).toEqual("");
  });
});
