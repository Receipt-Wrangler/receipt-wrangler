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

  it("resolves the role name for a matching id", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(2, roles)).toEqual("Legacy Owner");
  });

  it("returns empty string for an unknown id", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(99, roles)).toEqual("");
  });

  it("returns empty string for a null/undefined id", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(null, roles)).toEqual("");
    expect(pipe.transform(undefined, roles)).toEqual("");
  });

  it("returns empty string when roles are not yet loaded", () => {
    const pipe = new RoleNamePipe();
    expect(pipe.transform(1, [])).toEqual("");
  });
});
