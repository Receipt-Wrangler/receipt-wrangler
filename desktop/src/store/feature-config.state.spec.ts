import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { FeatureConfig } from "../open-api";
import { FeatureConfigState } from "./feature-config.state";
import { SetFeatureConfig } from "./feature-config.state.actions";

describe("FeatureConfigState", () => {
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([FeatureConfigState])],
    });
    store = TestBed.inject(Store);
  });

  it("defaults loginQrUrl to an empty string", () => {
    expect(store.selectSnapshot(FeatureConfigState.loginQrUrl)).toBe("");
  });

  // Guards the SetFeatureConfig patchState block: it lists each field
  // explicitly, so a new field is silently dropped unless added there.
  it("patches loginQrUrl from SetFeatureConfig", () => {
    const config: FeatureConfig = {
      enableLocalSignUp: false,
      aiPoweredReceipts: false,
      loginQrUrl: "https://receiptwrangler.io/app/setup#url=x",
    };

    store.dispatch(new SetFeatureConfig(config));

    expect(store.selectSnapshot(FeatureConfigState.loginQrUrl)).toBe(
      "https://receiptwrangler.io/app/setup#url=x"
    );
  });

  it("keeps the other feature flags when patching", () => {
    store.dispatch(
      new SetFeatureConfig({
        enableLocalSignUp: true,
        aiPoweredReceipts: true,
        loginQrUrl: "https://receiptwrangler.io/app/setup#url=x",
      })
    );

    expect(store.selectSnapshot(FeatureConfigState.enableLocalSignUp)).toBe(true);
    expect(store.selectSnapshot(FeatureConfigState.aiPoweredReceipts)).toBe(true);
  });
});
