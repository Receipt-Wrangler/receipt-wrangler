import { Store } from "@ngxs/store";
import { forkJoin, Observable, of } from "rxjs";
import { AppData, CurrencySeparator, CurrencySymbolPosition } from "../open-api";
import {
  GroupState,
  SetAuthState,
  SetFeatureConfig,
  SetGroups,
  SetIcons,
  SetSelectedGroupId,
  SetUserPreferences,
  SetUsers,
} from "../store";
import { SetAbout } from "../store/about.state.actions";
import { SetCurrencyData, SetCurrencyDisplay } from "../store/system-settings.state.actions";

export function setAppData(store: Store, appData: AppData): Observable<any[]> {
  let selectedGroupIdObservable: Observable<void> = of(undefined);
  if (!store.selectSnapshot(GroupState.selectedGroupId)) {
    selectedGroupIdObservable = store.dispatch(
      new SetSelectedGroupId(appData.groups[0].id.toString())
    );
  }

  return forkJoin([
    store.dispatch(new SetAuthState(appData.claims)),
    store.dispatch(new SetFeatureConfig(appData.featureConfig)),
    store.dispatch(new SetGroups(appData.groups)),
    store.dispatch(new SetUserPreferences(appData.userPreferences)),
    store.dispatch(new SetUsers(appData.users)),
    store.dispatch(new SetCurrencyDisplay(appData.currencyDisplay)),
    store.dispatch(new SetCurrencyData(
      appData.currencySymbolPosition ?? CurrencySymbolPosition.Start,
      appData.currencyDecimalSeparator ?? CurrencySeparator.Period,
      appData.currencyThousandthsSeparator ?? CurrencySeparator.Comma,
      appData.currencyHideDecimalPlaces ?? false
    )),
    store.dispatch(new SetIcons(appData.icons)),
    store.dispatch(new SetAbout(appData.about)),
    selectedGroupIdObservable,
  ]);
}

