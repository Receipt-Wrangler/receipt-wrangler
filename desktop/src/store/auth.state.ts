import { Injectable } from "@angular/core";
import { Action, createSelector, Selector, State, StateContext } from "@ngxs/store";

import { hasAll, hasAny } from "../utils/permission.utils";
import { Category, Icon, Tag, UserPreferences } from "../open-api";
import { User } from "../open-api/model/user";
import { AuthStateInterface } from "./auth-state.interface";
import { Logout, SetAuthState, SetGroupCatalog, SetIcons, SetPermissions, SetUserPreferences } from "./auth.state.actions";

@State<AuthStateInterface>({
  name: "auth",
  defaults: {},
})
@Injectable()
export class AuthState {

  @Selector()
  static userPreferences(
    state: AuthStateInterface
  ): UserPreferences | undefined {
    return state.userPreferences;
  }

  @Selector()
  static icons(state: AuthStateInterface): Icon[] {
    return state.icons ?? [];
  }

  @Selector()
  static isLoggedIn(state: AuthStateInterface): boolean {
    return !AuthState.isTokenExpired(state);
  }

  @Selector()
  static userId(state: AuthStateInterface): string {
    return state.userId ?? "";
  }

  @Selector()
  static isTokenExpired(state: AuthStateInterface): boolean {
    if (state.expirationDate) {
      return new Date() >= new Date(Number(state.expirationDate) * 1000);
    } else {
      return true;
    }
  }

  @Selector()
  static loggedInUser(state: AuthStateInterface): User {
    return {
      defaultAvatarColor: state.defaultAvatarColor ?? "",
      displayName: state.displayname ?? "",
      id: Number(state.userId) ?? "",
      username: state.username ?? "",
    } as User;
  }

  @Selector()
  static appPermissions(state: AuthStateInterface): string[] {
    return state.appPermissions ?? [];
  }

  @Selector()
  static groupPermissions(state: AuthStateInterface): {
    [groupId: number]: string[];
  } {
    return state.groupPermissions ?? {};
  }

  static groupCategories(groupId: number) {
    return createSelector([AuthState], (state: AuthStateInterface): Category[] => {
      return state.groupCategories?.[groupId] ?? [];
    });
  }

  static groupTags(groupId: number) {
    return createSelector([AuthState], (state: AuthStateInterface): Tag[] => {
      return state.groupTags?.[groupId] ?? [];
    });
  }

  static hasAppPermission(permission: string) {
    return createSelector([AuthState], (state: AuthStateInterface) => {
      return hasAll(state.appPermissions ?? [], permission);
    });
  }

  static hasAnyAppPermission(permissions: string[]) {
    return createSelector([AuthState], (state: AuthStateInterface) => {
      return hasAny(state.appPermissions ?? [], ...permissions);
    });
  }

  static hasGroupPermission(
    groupId: number,
    permission: string,
    orAppPermissions: string[] = []
  ) {
    return createSelector([AuthState], (state: AuthStateInterface) => {
      if (
        orAppPermissions.length > 0 &&
        hasAny(state.appPermissions ?? [], ...orAppPermissions)
      ) {
        return true;
      }
      return hasAll(state.groupPermissions?.[groupId] ?? [], permission);
    });
  }

  @Action(SetAuthState)
  setAuthState(
    { getState, patchState }: StateContext<AuthStateInterface>,
    payload: SetAuthState
  ) {
    const claims = payload.userClaims;

    patchState({
      defaultAvatarColor: claims.defaultAvatarColor,
      displayname: claims.displayName,
      expirationDate: claims?.exp?.toString(),
      userId: claims?.userId?.toString(),
      username: claims?.username,
    });
  }

  @Action(SetPermissions)
  setPermissions(
    { patchState }: StateContext<AuthStateInterface>,
    { appPermissions, groupPermissions }: SetPermissions
  ) {
    patchState({
      appPermissions,
      groupPermissions,
    });
  }

  @Action(SetGroupCatalog)
  setGroupCatalog(
    { patchState }: StateContext<AuthStateInterface>,
    { groupCategories, groupTags }: SetGroupCatalog
  ) {
    patchState({
      groupCategories,
      groupTags,
    });
  }

  @Action(Logout)
  logout({ getState, patchState }: StateContext<AuthStateInterface>) {
    patchState({
      defaultAvatarColor: "",
      displayname: "",
      expirationDate: "",
      userId: "",
      username: "",
      userPreferences: undefined,
      appPermissions: undefined,
      groupPermissions: undefined,
      groupCategories: undefined,
      groupTags: undefined,
    });
  }

  @Action(SetUserPreferences)
  setUserPreferences(
    { patchState }: StateContext<AuthStateInterface>,
    payload: SetUserPreferences
  ) {
    patchState({
      userPreferences: payload.userPreferences,
    });
  }

  @Action(SetIcons)
  setIcons(
    { patchState }: StateContext<AuthStateInterface>,
    payload: SetIcons
  ) {
    patchState({
      icons: payload.icons,
    });
  }
}
