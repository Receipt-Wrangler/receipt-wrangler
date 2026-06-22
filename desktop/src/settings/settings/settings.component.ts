import { Component, computed } from "@angular/core";
import { Store } from "@ngxs/store";
import { DEFAULT_HOST_CLASS } from "src/constants";
import { Permission } from "../../open-api";
import { AuthState } from "../../store";
import { TabConfig } from "../../shared-ui/tabs/tab-config.interface";

@Component({
    selector: "app-settings",
    templateUrl: "./settings.component.html",
    styleUrls: ["./settings.component.scss"],
    host: DEFAULT_HOST_CLASS,
    standalone: false
})
export class SettingsComponent {
  private readonly canReadProfile = this.store.selectSignal(
    AuthState.hasAppPermission(Permission.AppAccountRead)
  );

  private readonly canReadPreferences = this.store.selectSignal(
    AuthState.hasAppPermission(Permission.AppUserPreferencesRead)
  );

  private readonly canReadApiKeys = this.store.selectSignal(
    AuthState.hasAppPermission(Permission.AppApiKeysRead)
  );

  public readonly tabs = computed<TabConfig[]>(() => {
    const tabs: TabConfig[] = [];

    if (this.canReadProfile()) {
      tabs.push({
        label: "User Profile",
        routerLink: "user-profile/view",
        name: "user-profile",
      });
    }

    if (this.canReadPreferences()) {
      tabs.push({
        label: "User Preferences",
        routerLink: "user-preferences/view",
        name: "user-preferences",
      });
    }

    if (this.canReadApiKeys()) {
      tabs.push({
        label: "API Keys",
        routerLink: "api-keys/view",
        name: "api-keys",
      });
    }

    return tabs;
  });

  constructor(private store: Store) {}
}
