import { Injectable } from "@angular/core";
import { Action, createSelector, Selector, State, StateContext, } from "@ngxs/store";
import { FeatureConfig } from "../open-api";
import { SetFeatureConfig } from "./feature-config.state.actions";

@State<FeatureConfig>({
  name: "featureConfig",
  defaults: { enableLocalSignUp: true, aiPoweredReceipts: false },
})
@Injectable()
export class FeatureConfigState {
  @Selector()
  static enableLocalSignUp(state: FeatureConfig): boolean {
    return state.enableLocalSignUp as boolean;
  }

  @Selector()
  static aiPoweredReceipts(state: FeatureConfig): boolean {
    return state.aiPoweredReceipts as boolean;
  }

  @Selector()
  static featureConfig(state: FeatureConfig): FeatureConfig {
    return state;
  }

  static hasFeature(feature: string) {
    return createSelector([FeatureConfigState], (state: FeatureConfig) => {
      return !!(state as any)[feature];
    });
  }

  @Action(SetFeatureConfig)
  setFeatureConfig(
    { patchState }: StateContext<FeatureConfig>,
    payload: SetFeatureConfig
  ) {
    patchState({
      aiPoweredReceipts: payload.config?.aiPoweredReceipts,
      enableLocalSignUp: payload.config?.enableLocalSignUp,
    });
  }
}
