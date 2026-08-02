import { of } from "rxjs";
import { Group, GroupsService, Role } from "../../open-api";
import { GrantSelection } from "./grant-picker.component";
import {
  buildMemberGrantRows,
  ceilingForRole,
  emptySelectionHint,
  MemberGrantRow,
  saveChangedMemberGrants,
} from "./member-grant-assignment";

describe("member-grant-assignment", () => {
  const restrictedRole = {
    id: 1,
    name: "Foster Parent",
    categoryGrants: [1, 2],
    tagGrants: [],
  } as unknown as Role;

  const openRole = {
    id: 2,
    name: "Coordinator",
    categoryGrants: [],
    tagGrants: [],
  } as unknown as Role;

  describe("ceilingForRole", () => {
    it("returns null for a role that grants nothing", () => {
      expect(ceilingForRole(openRole)).toBeNull();
    });

    it("returns null when the member has no role", () => {
      expect(ceilingForRole(undefined)).toBeNull();
    });

    it("names the role so a narrowed picker explains itself", () => {
      const ceiling = ceilingForRole(restrictedRole);
      expect(ceiling?.categoryIds).toEqual([1, 2]);
      expect(ceiling?.label).toBe("role Foster Parent");
    });
  });

  describe("buildMemberGrantRows", () => {
    const groups = [
      {
        id: 100,
        name: "Agency",
        groupMembers: [
          { userId: 7, groupId: 100, groupRoleId: 1, categoryGrants: [1], tagGrants: [] },
          { userId: 8, groupId: 100, groupRoleId: 1, categoryGrants: [2], tagGrants: [] },
        ],
      },
      {
        id: 200,
        name: "Household",
        groupMembers: [{ userId: 9, groupId: 200, groupRoleId: 2 }],
      },
    ] as unknown as Group[];

    it("returns one row per group the user belongs to", () => {
      const rows = buildMemberGrantRows(7, groups, [restrictedRole, openRole]);

      expect(rows.length).toBe(1);
      expect(rows[0].groupId).toBe(100);
      expect(rows[0].groupName).toBe("Agency");
      expect(rows[0].roleName).toBe("Foster Parent");
      expect(rows[0].current.categoryIds).toEqual([1]);
    });

    it("skips groups the user is not a member of", () => {
      expect(buildMemberGrantRows(999, groups, [restrictedRole, openRole])).toEqual([]);
    });

    it("carries no ceiling when the member's role grants nothing", () => {
      const rows = buildMemberGrantRows(9, groups, [restrictedRole, openRole]);
      expect(rows[0].ceiling).toBeNull();
    });

    it("defaults a member with no stored grants to an empty selection", () => {
      const rows = buildMemberGrantRows(9, groups, [restrictedRole, openRole]);
      expect(rows[0].current).toEqual({ categoryIds: [], tagIds: [] });
    });
  });

  describe("require-individual flags", () => {
    // These invert what an empty picker means, so the row has to carry them or the
    // hint tells the admin the exact opposite of what the backend will do.
    const requiresBoth = {
      id: 3,
      name: "Strict",
      categoryGrants: [],
      tagGrants: [],
      requiresIndividualCategoryGrants: true,
      requiresIndividualTagGrants: true,
    } as unknown as Role;

    const groups = [
      {
        id: 100,
        name: "Agency",
        groupMembers: [{ userId: 7, groupId: 100, groupRoleId: 3 }],
      },
    ] as unknown as Group[];

    it("carries both flags from the role onto the row", () => {
      const [row] = buildMemberGrantRows(7, groups, [requiresBoth]);

      expect(row.requiresIndividualCategories).toBe(true);
      expect(row.requiresIndividualTags).toBe(true);
    });

    it("defaults both flags to false when the role omits them", () => {
      const [row] = buildMemberGrantRows(7, groups, [
        { id: 3, name: "Strict", categoryGrants: [], tagGrants: [] } as unknown as Role,
      ]);

      expect(row.requiresIndividualCategories).toBe(false);
      expect(row.requiresIndividualTags).toBe(false);
    });

    it("defaults both flags to false when the member has no role at all", () => {
      const roleless = [
        {
          id: 100,
          name: "Agency",
          groupMembers: [{ userId: 7, groupId: 100 }],
        },
      ] as unknown as Group[];

      const [row] = buildMemberGrantRows(7, roleless, [requiresBoth]);

      expect(row.requiresIndividualCategories).toBe(false);
      expect(row.requiresIndividualTags).toBe(false);
    });
  });

  describe("emptySelectionHint", () => {
    function rowWith(categories: boolean, tags: boolean): MemberGrantRow {
      return {
        groupId: 1,
        groupName: "Agency",
        roleName: "Strict",
        ceiling: null,
        current: { categoryIds: [], tagIds: [] },
        requiresIndividualCategories: categories,
        requiresIndividualTags: tags,
      };
    }

    it("promises everything the role allows when neither resource requires an assignment", () => {
      expect(emptySelectionHint(rowWith(false, false))).toBe(
        "Leave empty to give them everything their role allows."
      );
    });

    it("warns that empty means nothing when both resources require an assignment", () => {
      const hint = emptySelectionHint(rowWith(true, true));

      expect(hint).toContain("no categories and no tags");
      expect(hint).not.toContain("everything");
    });

    it("distinguishes the resources when only categories require an assignment", () => {
      const hint = emptySelectionHint(rowWith(true, false));

      expect(hint).toContain("leaving categories empty gives them none");
      expect(hint).toContain("every tag the role allows");
    });

    it("distinguishes the resources when only tags require an assignment", () => {
      const hint = emptySelectionHint(rowWith(false, true));

      expect(hint).toContain("leaving tags empty gives them none");
      expect(hint).toContain("every category the role allows");
    });
  });

  describe("saveChangedMemberGrants", () => {
    let groupsService: Pick<GroupsService, "updateGroupMemberGrants"> | any;

    beforeEach(() => {
      groupsService = {
        updateGroupMemberGrants: jest.fn().mockReturnValue(of({})),
      };
    });

    const row = {
      groupId: 100,
      groupName: "Agency",
      roleName: "Foster Parent",
      ceiling: null,
      current: { categoryIds: [1], tagIds: [] },
    };

    it("writes nothing when no edit was made", (done) => {
      saveChangedMemberGrants(groupsService, 7, [row], new Map()).subscribe(() => {
        expect(groupsService.updateGroupMemberGrants).not.toHaveBeenCalled();
        done();
      });
    });

    it("writes nothing when the edit matches what is already stored", (done) => {
      // Opening and closing a form must not rewrite assignments — a needless write
      // could also trip the endpoint's ceiling check on a since-narrowed role.
      const edited = new Map<number, GrantSelection>([
        [100, { categoryIds: [1], tagIds: [] }],
      ]);

      saveChangedMemberGrants(groupsService, 7, [row], edited).subscribe(() => {
        expect(groupsService.updateGroupMemberGrants).not.toHaveBeenCalled();
        done();
      });
    });

    it("ignores id ordering when deciding whether something changed", (done) => {
      const unordered = {
        ...row,
        current: { categoryIds: [2, 1], tagIds: [] },
      };
      const edited = new Map<number, GrantSelection>([
        [100, { categoryIds: [1, 2], tagIds: [] }],
      ]);

      saveChangedMemberGrants(groupsService, 7, [unordered], edited).subscribe(() => {
        expect(groupsService.updateGroupMemberGrants).not.toHaveBeenCalled();
        done();
      });
    });

    it("writes only the groups whose assignment actually changed", (done) => {
      const secondRow = {
        ...row,
        groupId: 200,
        groupName: "Household",
        current: { categoryIds: [], tagIds: [] },
      };
      const edited = new Map<number, GrantSelection>([
        [100, { categoryIds: [1], tagIds: [] }],
        [200, { categoryIds: [5], tagIds: [] }],
      ]);

      saveChangedMemberGrants(groupsService, 7, [row, secondRow], edited).subscribe(() => {
        expect(groupsService.updateGroupMemberGrants).toHaveBeenCalledTimes(1);
        expect(groupsService.updateGroupMemberGrants).toHaveBeenCalledWith(200, 7, {
          categoryGrants: [5],
          tagGrants: [],
        });
        done();
      });
    });

    it("writes a cleared assignment", (done) => {
      const edited = new Map<number, GrantSelection>([
        [100, { categoryIds: [], tagIds: [] }],
      ]);

      saveChangedMemberGrants(groupsService, 7, [row], edited).subscribe(() => {
        expect(groupsService.updateGroupMemberGrants).toHaveBeenCalledWith(100, 7, {
          categoryGrants: [],
          tagGrants: [],
        });
        done();
      });
    });
  });
});
