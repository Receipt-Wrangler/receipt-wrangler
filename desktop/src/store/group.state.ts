import { Injectable } from '@angular/core';
import {
  Action,
  createSelector,
  Selector,
  State,
  StateContext,
} from '@ngxs/store';
import { Group } from '../open-api/model/group';
import {
  AddGroup,
  RemoveGroup,
  SetGroups,
  SetSelectedDashboardId,
  SetSelectedGroupId,
  UpdateGroup,
} from './group.state.actions';
import { GroupMember } from '../open-api/model/groupMember';

export interface GroupStateInterface {
  groups: Group[];
  selectedGroupId: string;
  selectedDashboardId: string;
}

@State<GroupStateInterface>({
  name: 'groups',
  defaults: {
    groups: [],
    selectedGroupId: '',
    selectedDashboardId: '',
  },
})
@Injectable()
export class GroupState {
  @Selector()
  static groups(state: GroupStateInterface): Group[] {
    return state.groups;
  }

  @Selector()
  static allGroupMembers(state: GroupStateInterface): GroupMember[] {
    return state.groups.map((g) => g.groupMembers).flat();
  }

  @Selector()
  static groupsWithoutAll(state: GroupStateInterface): Group[] {
    return state.groups.filter((g) => !g.isAllGroup);
  }

  @Selector()
  static groupsWithoutSelectedGroup(state: GroupStateInterface): Group[] {
    return state.groups.filter(
      (g) => g.id.toString() !== state.selectedGroupId
    );
  }

  @Selector()
  static selectedDashboardId(state: GroupStateInterface): string {
    return state.selectedDashboardId;
  }

  @Selector()
  static selectedGroupId(state: GroupStateInterface): string {
    return state.selectedGroupId;
  }

  @Selector()
  static receiptListLink(state: GroupStateInterface): string {
    return `/receipts/group/${state.selectedGroupId}`;
  }

  @Selector()
  static dashboardLink(state: GroupStateInterface): string {
    return `/dashboard/group/${state.selectedGroupId}`;
  }

  @Selector()
  static settingsLinkBase(state: GroupStateInterface): string {
    return `/groups/${state.selectedGroupId}/settings`;
  }

  static getGroupById(groupId: string) {
    return createSelector([GroupState], (state: GroupStateInterface) => {
      return state.groups.find((g) => g.id.toString() === groupId.toString());
    });
  }

  @Action(AddGroup)
  addGroup(
    { getState, patchState }: StateContext<GroupStateInterface>,
    payload: AddGroup
  ) {
    const groups = Array.from(getState().groups);
    groups.push(payload.group);

    patchState({
      groups: groups,
    });
  }

  @Action(RemoveGroup)
  removeGroup(
    { getState, patchState }: StateContext<GroupStateInterface>,
    payload: RemoveGroup
  ) {
    const state = getState();
    const group = GroupState.getGroupById(payload.groupId)(state);
    if (group) {
      const index = state.groups.findIndex((g) => g === group);
      if (index >= 0) {
        const newInterface = {} as GroupStateInterface;
        const newGroups = Array.from(state.groups).filter(
          (g) => g.id !== group.id
        );
        newInterface.groups = newGroups;
        if (group.id.toString() === state.selectedGroupId.toString()) {
          newInterface.selectedGroupId = state.groups[0].id.toString();
        }
        patchState(newInterface);
      }
    }
  }

  @Action(SetGroups)
  setGroups(
    { patchState }: StateContext<GroupStateInterface>,
    payload: SetGroups
  ) {
    patchState({
      groups: payload.groups,
    });
  }

  @Action(UpdateGroup)
  updateGroup(
    { getState, patchState }: StateContext<GroupStateInterface>,
    payload: UpdateGroup
  ) {
    const groupIndex = getState().groups.findIndex(
      (g) => g.id?.toString() === payload?.group?.id?.toString()
    );
    if (groupIndex > -1) {
      const newGroups = Array.from(getState().groups);
      newGroups[groupIndex] = payload.group;

      patchState({
        groups: newGroups,
      });
    }
  }

  @Action(SetSelectedDashboardId)
  setSelectedDashboardId(
    { getState, patchState }: StateContext<GroupStateInterface>,
    payload: SetSelectedDashboardId
  ) {
    patchState({
      selectedDashboardId: payload.dashboardId,
    });
  }

  @Action(SetSelectedGroupId)
  setSelectedGroupId(
    { getState, patchState }: StateContext<GroupStateInterface>,
    payload: SetSelectedGroupId
  ) {
    let groupId = '';
    let dashboardId = '';

    if (payload?.groupId) {
      groupId = payload.groupId;
    } else {
      const groups = getState().groups;
      groupId = groups[0].id.toString();
    }

    if (payload.groupId === getState().selectedGroupId) {
      dashboardId = getState().selectedDashboardId;
    }

    patchState({
      selectedGroupId: groupId,
      selectedDashboardId: dashboardId,
    });
  }
}
