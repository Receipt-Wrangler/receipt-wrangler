import { Injectable } from "@angular/core";
import { Action, createSelector, Selector, State, StateContext } from "@ngxs/store";

import { Icon, UserPreferences } from "../open-api";
import { User } from "../open-api/model/user";
import { AuthStateInterface } from "./auth-state.interface";
import { Logout, SetAuthState, SetIcons, SetUserPreferences } from "./auth.state.actions";

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
  static userRole(state: AuthStateInterface): string {
    return state.userRole ?? "";
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

  static hasRole(role: string) {
    return createSelector([AuthState], (state: AuthStateInterface) => {
      return state.userRole === role;
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
      userRole: claims?.userRole,
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
      userRole: undefined,
      userPreferences: undefined,
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
