import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { SystemSettingsState } from "./system-settings.state";
import { SetTimeZone } from "./system-settings.state.actions";

describe("SystemSettingsState time zone", () => {
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [NgxsModule.forRoot([SystemSettingsState])] });
    store = TestBed.inject(Store);
  });

  it("defaults to UTC", () => {
    expect(store.selectSnapshot(SystemSettingsState.timeZone)).toBe("UTC");
  });

  it("is set from AppData", () => {
    store.dispatch(new SetTimeZone("America/New_York"));

    expect(store.selectSnapshot(SystemSettingsState.timeZone)).toBe("America/New_York");
  });

  it("treats an empty value as UTC", () => {
    store.dispatch(new SetTimeZone(""));

    expect(store.selectSnapshot(SystemSettingsState.timeZone)).toBe("UTC");
  });

  // Defaults never run for a slice hydrated from localStorage, so a session
  // persisted before the key existed has none.
  it("falls back to UTC for a persisted slice that predates the key", () => {
    store.reset({ systemSettings: { currencyDisplay: "$" } });

    expect(store.selectSnapshot(SystemSettingsState.timeZone)).toBe("UTC");
  });
});
