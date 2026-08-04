import { forkJoin, Observable, of } from "rxjs";
import { Group, GroupMember, GroupsService, Role } from "../../open-api";
import { GrantCeiling, GrantSelection } from "./grant-picker.component";

/**
 * One editable row of per-member category/tag assignment: the group, the
 * membership inside it, and the ceiling its group role imposes.
 *
 * Shared by the user form (which shows one row per group the user belongs to)
 * and the group-member form (which shows the single row for that membership),
 * so the two entry points cannot drift in how they resolve the ceiling or write
 * the grants.
 */
export interface MemberGrantRow {
  groupId: number;
  groupName: string;
  roleName: string;
  ceiling: GrantCeiling | null;
  current: GrantSelection;

  /**
   * Whether the group role makes an individual assignment MANDATORY, per resource.
   *
   * These invert what an empty selection means: normally empty is "no narrowing"
   * and the member gets everything the role allows, but under these flags an
   * unconfigured member sees NOTHING (the backend's fail-closed rule — see
   * `api/CLAUDE.md` → "Category/tag grant resolution"). The hint text has to say
   * so, which is why they are carried on the row rather than folded into the
   * ceiling: a ceiling bounds what may be picked, these change what picking
   * nothing means.
   */
  requiresIndividualCategories: boolean;
  requiresIndividualTags: boolean;
}

/**
 * The sentence explaining what leaving a picker empty does for this membership.
 *
 * Shared by both hosts so the two entry points cannot describe the same rule
 * differently — the wording is the only thing telling an admin that an empty
 * selection can mean "nothing" rather than "everything".
 */
export function emptySelectionHint(row: MemberGrantRow): string {
  const categories = row.requiresIndividualCategories;
  const tags = row.requiresIndividualTags;

  if (categories && tags) {
    return "Their role requires an individual assignment, so leaving these empty gives them no categories and no tags.";
  }
  if (categories) {
    return "Their role requires an individual category assignment, so leaving categories empty gives them none. Leaving tags empty gives them every tag the role allows.";
  }
  if (tags) {
    return "Their role requires an individual tag assignment, so leaving tags empty gives them none. Leaving categories empty gives them every category the role allows.";
  }

  return "Leave empty to give them everything their role allows.";
}

/**
 * Builds the ceiling a group role imposes on its members' assignments.
 *
 * Returns null when the role imposes none — either the member has no role, or
 * the role grants nothing (the backend's "empty grants = unrestricted" rule).
 * The label names the role so a narrowed picker explains itself.
 */
export function ceilingForRole(role: Role | undefined): GrantCeiling | null {
  if (!role) {
    return null;
  }

  const categoryIds = role.categoryGrants ?? [];
  const tagIds = role.tagGrants ?? [];
  if (!categoryIds.length && !tagIds.length) {
    return null;
  }

  return { categoryIds, tagIds, label: `role ${role.name}` };
}

/** Builds a row for one membership, resolving its role from the loaded role list. */
export function buildMemberGrantRow(
  group: Group,
  member: GroupMember,
  groupRoles: Role[],
): MemberGrantRow {
  const role = groupRoles.find((candidate) => candidate.id === member.groupRoleId);

  return {
    groupId: group.id,
    groupName: group.name,
    roleName: role?.name ?? "",
    ceiling: ceilingForRole(role),
    current: {
      categoryIds: member.categoryGrants ?? [],
      tagIds: member.tagGrants ?? [],
    },
    requiresIndividualCategories: role?.requiresIndividualCategoryGrants ?? false,
    requiresIndividualTags: role?.requiresIndividualTagGrants ?? false,
  };
}

/** Every group the user is a member of, as assignment rows. */
export function buildMemberGrantRows(
  userId: number,
  groups: Group[],
  groupRoles: Role[],
): MemberGrantRow[] {
  const rows: MemberGrantRow[] = [];

  for (const group of groups) {
    const member = group.groupMembers?.find((candidate) => candidate.userId === userId);
    if (member) {
      rows.push(buildMemberGrantRow(group, member, groupRoles));
    }
  }

  return rows;
}

/**
 * Persists the assignments that actually changed, one request per group.
 *
 * Grants are written through their own endpoint rather than the group-member
 * upsert, so they are saved independently of whatever form hosts the picker.
 * Unchanged rows are skipped so opening and closing a form never rewrites
 * assignments (and never trips the endpoint's ceiling check on a role that was
 * narrowed since the assignment was made).
 */
export function saveChangedMemberGrants(
  groupsService: GroupsService,
  userId: number,
  rows: MemberGrantRow[],
  edited: Map<number, GrantSelection>,
): Observable<unknown> {
  const requests = rows
    .filter((row) => {
      const selection = edited.get(row.groupId);
      return selection !== undefined && !sameSelection(selection, row.current);
    })
    .map((row) =>
      groupsService.updateGroupMemberGrants(row.groupId, userId, {
        categoryGrants: edited.get(row.groupId)!.categoryIds,
        tagGrants: edited.get(row.groupId)!.tagIds,
      }),
    );

  return requests.length ? forkJoin(requests) : of(undefined);
}

function sameSelection(left: GrantSelection, right: GrantSelection): boolean {
  return sameIds(left.categoryIds, right.categoryIds) && sameIds(left.tagIds, right.tagIds);
}

function sameIds(left: number[], right: number[]): boolean {
  if (left.length !== right.length) {
    return false;
  }
  const sortedLeft = [...left].sort((a, b) => a - b);
  const sortedRight = [...right].sort((a, b) => a - b);
  return sortedLeft.every((id, index) => id === sortedRight[index]);
}
