import { Injectable } from '@angular/core';
import { Action, Selector, State, StateContext } from '@ngxs/store';
import {
  HideProgressBar,
  ShowProgressBar,
  ToggleIsSidebarOpen,
  ToggleShowProgressBar,
} from './layout.state.actions';

export interface LayoutStateInterface {
  isSidebarOpen: boolean;
  showProgressBar: boolean;
}

@State<LayoutStateInterface>({
  name: 'layout',
  defaults: {
    isSidebarOpen: false,
    showProgressBar: false,
  },
})
@Injectable()
export class LayoutState {
  @Selector()
  static isSidebarOpen(state: LayoutStateInterface): boolean {
    return state.isSidebarOpen;
  }

  @Selector()
  static showProgressBar(state: LayoutStateInterface): boolean {
    return state.showProgressBar;
  }

  @Action(ToggleIsSidebarOpen)
  setIsSidebarOpen({
    getState,
    patchState,
  }: StateContext<LayoutStateInterface>) {
    patchState({
      isSidebarOpen: !getState().isSidebarOpen,
    });
  }

  @Action(ToggleShowProgressBar)
  toggleShowProgressBar({
    getState,
    patchState,
  }: StateContext<LayoutStateInterface>) {
    patchState({
      showProgressBar: !getState().showProgressBar,
    });
  }

  @Action(HideProgressBar)
  hideProgressBar({ patchState }: StateContext<LayoutStateInterface>) {
    patchState({
      showProgressBar: false,
    });
  }

  @Action(ShowProgressBar)
  showProgressBar({ patchState }: StateContext<LayoutStateInterface>) {
    patchState({
      showProgressBar: true,
    });
  }
}
