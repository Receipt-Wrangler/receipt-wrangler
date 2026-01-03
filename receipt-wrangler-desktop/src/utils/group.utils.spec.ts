import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { GroupRole } from "../open-api";
import { AuthState, GroupState } from "../store";
import { GroupUtil } from "./group.utils";

describe("GroupUtil", () => {
  let groupUtil: GroupUtil;
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState, GroupState])],
      providers: [GroupUtil],
    });

    groupUtil = TestBed.inject(GroupUtil);
    store = TestBed.inject(Store);
  });

  it("should be created", () => {
    expect(groupUtil).toBeTruthy();
  });

  describe("hasGroupAccess", () => {
    const testGroupId = 1;
    const testGroupRole = GroupRole.Editor;

    it("should return true when groupId is undefined", () => {
      const result = groupUtil.hasGroupAccess(undefined, testGroupRole, false);
      expect(result).toBe(true);
    });

    it("should return false when group is not found in store", () => {
      const result = groupUtil.hasGroupAccess(testGroupId, testGroupRole, false);
      expect(result).toBe(false);
    });

    it("should return false when user is not a member of group", () => {
      const userId = "2";
      store.reset({
        auth: { userId: userId },
        groups: {
          groups: [
            {
              id: testGroupId,
              groupMembers: [{ userId: "3", groupRole: GroupRole.Editor }],
            },
          ],
        },
      });

      const result = groupUtil.hasGroupAccess(testGroupId, testGroupRole, false);
      expect(result).toBe(false);
    });

    it("should return false when user has lower role than required", () => {
      const userId = "1";
      store.reset({
        auth: { userId: userId },
        groups: {
          groups: [
            {
              id: testGroupId,
              groupMembers: [{ userId: userId, groupRole: GroupRole.Viewer }],
            },
          ],
        },
      });

      const result = groupUtil.hasGroupAccess(testGroupId, testGroupRole, false);
      expect(result).toBe(false);
    });

    it("should return true when user has same or higher role than required", () => {
      const userId = "1";
      store.reset({
        auth: { userId: userId },
        groups: {
          groups: [
            {
              id: testGroupId,
              groupMembers: [{ userId: userId, groupRole: GroupRole.Owner }],
            },
          ],
        },
      });

      const result = groupUtil.hasGroupAccess(testGroupId, testGroupRole, false);
      expect(result).toBe(true);
    });
  });


  it("should return true when allowAdminOverride is true and user is admin", () => {
    jest.spyOn(store, "selectSnapshot")
      .mockReturnValueOnce(true)
      .mockReturnValueOnce("1")
      .mockReturnValueOnce(null);
    const result = groupUtil.hasGroupAccess(1, GroupRole.Viewer, true);
    expect(result).toBe(true);
  });

  it("should return false when allowAdminOverride is false and user is admin but not in group", () => {
    const userId = "1";
    const group = {
      id: "1",
      groupMembers: [{ userId: "2" }]
    };
    store.reset({
      auth: { userId: userId },
      groups: {
        groups: [group]
      },
    });
    const result = groupUtil.hasGroupAccess(1, GroupRole.Viewer, false);
    expect(result).toBe(false);
  });
});
