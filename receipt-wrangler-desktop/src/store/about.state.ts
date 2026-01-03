import { Injectable } from "@angular/core";
import { Action, Selector, State, StateContext } from "@ngxs/store";

import { About } from "../open-api";
import { AboutStateInterface } from "./about-state.interface";
import { SetAbout } from "./about.state.actions";

@State<AboutStateInterface>({
  name: "about",
  defaults: {
    about: {
      buildDate: "",
      version: "",
    } as About,
  },
})
@Injectable()
export class AboutState {
  @Selector()
  static about(state: AboutStateInterface): About {
    return state.about;
  }


  @Action(SetAbout)
  setAuthState(
    { patchState }: StateContext<AboutStateInterface>,
    payload: SetAbout
  ) {
    patchState({
      about: payload.about,
    });
  }
}
