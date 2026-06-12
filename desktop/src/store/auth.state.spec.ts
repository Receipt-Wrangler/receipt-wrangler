import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { Claims } from "../open-api";
import { AuthState } from "./auth.state";
import { Logout, SetAuthState, SetPermissions } from "./auth.state.actions";

describe("AuthState", () => {
  let store: Store;

  const APP_PERMISSIONS = ["app.users.read", "app.roles.read"];
  const GROUP_PERMISSIONS = { 1: ["group.view", "group.receipts.read"] };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState])],
    }).compileComponents();

    store = TestBed.inject(Store);
  });

  it("SetPermissions stores both app and group permissions", () => {
    store.dispatch(new SetPermissions(APP_PERMISSIONS, GROUP_PERMISSIONS));

    expect(store.selectSnapshot(AuthState.appPermissions)).toEqual(APP_PERMISSIONS);
    expect(store.selectSnapshot(AuthState.hasAppPermission("app.users.read"))).toBe(true);
    expect(store.selectSnapshot(AuthState.hasGroupPermission(1, "group.view"))).toBe(true);
  });

  it("Logout clears both permission maps", () => {
    store.dispatch(new SetPermissions(APP_PERMISSIONS, GROUP_PERMISSIONS));
    store.dispatch(new Logout());

    expect(store.selectSnapshot(AuthState.appPermissions)).toEqual([]);
    expect(store.selectSnapshot(AuthState.hasAppPermission("app.users.read"))).toBe(false);
    expect(store.selectSnapshot(AuthState.hasGroupPermission(1, "group.view"))).toBe(false);
  });

  it("SetAuthState (claims only) does NOT wipe permissions", () => {
    store.dispatch(new SetPermissions(APP_PERMISSIONS, GROUP_PERMISSIONS));

    const claims = {
      userId: 5,
      displayName: "Tester",
      defaultAvatarColor: "#fff",
      username: "tester",
      exp: 9999999999,
    } as Claims;
    store.dispatch(new SetAuthState(claims));

    expect(store.selectSnapshot(AuthState.appPermissions)).toEqual(APP_PERMISSIONS);
    expect(store.selectSnapshot(AuthState.hasAppPermission("app.users.read"))).toBe(true);
  });

  describe("selectors", () => {
    beforeEach(() => {
      store.dispatch(new SetPermissions(APP_PERMISSIONS, GROUP_PERMISSIONS));
    });

    it("hasAppPermission denies a permission the user does not hold", () => {
      expect(store.selectSnapshot(AuthState.hasAppPermission("app.users.delete"))).toBe(false);
    });

    it("hasAnyAppPermission returns true when at least one matches", () => {
      expect(
        store.selectSnapshot(
          AuthState.hasAnyAppPermission(["app.users.delete", "app.roles.read"])
        )
      ).toBe(true);
    });

    it("hasAnyAppPermission returns false when none match", () => {
      expect(
        store.selectSnapshot(
          AuthState.hasAnyAppPermission(["app.users.delete", "app.users.create"])
        )
      ).toBe(false);
    });

    it("hasGroupPermission checks the group's permission list", () => {
      expect(store.selectSnapshot(AuthState.hasGroupPermission(1, "group.receipts.read"))).toBe(true);
      expect(store.selectSnapshot(AuthState.hasGroupPermission(1, "group.receipts.update"))).toBe(false);
      // Non-member group resolves to deny.
      expect(store.selectSnapshot(AuthState.hasGroupPermission(2, "group.view"))).toBe(false);
    });

    it("hasGroupPermission orApp override grants access to a non-member holding the app perm", () => {
      // group 2 has no permissions, but the app fallback is held → allowed.
      expect(
        store.selectSnapshot(
          AuthState.hasGroupPermission(2, "group.view", ["app.users.read"])
        )
      ).toBe(true);
    });

    it("hasGroupPermission orApp override does not grant when the app perm is absent", () => {
      expect(
        store.selectSnapshot(
          AuthState.hasGroupPermission(2, "group.view", ["app.groups.read"])
        )
      ).toBe(false);
    });

    it("hasGroupPermission coerces string JSON keys to numeric index access", () => {
      // Persisted/JSON-revived state has string keys; numeric-id lookup must still resolve.
      store.dispatch(
        new SetPermissions(APP_PERMISSIONS, { ["3" as any]: ["group.view"] })
      );
      expect(store.selectSnapshot(AuthState.hasGroupPermission(3, "group.view"))).toBe(true);
    });
  });
});
