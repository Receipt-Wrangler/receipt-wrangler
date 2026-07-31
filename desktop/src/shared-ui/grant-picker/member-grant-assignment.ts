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
