import { Component, OnInit } from "@angular/core";
import { DEFAULT_HOST_CLASS } from "src/constants";
import { TabConfig } from "../../shared-ui/tabs/tab-config.interface";

@Component({
    selector: "app-settings",
    templateUrl: "./settings.component.html",
    styleUrls: ["./settings.component.scss"],
    host: DEFAULT_HOST_CLASS,
    standalone: false
})
export class SettingsComponent implements OnInit {
  public tabs: TabConfig[] = [];

  public ngOnInit(): void {
    this.initTabs();
  }

  private initTabs(): void {
    this.tabs = [
      {
        label: "User Profile",
        routerLink: "user-profile/view",
        name: "user-profile",
      },
      {
        label: "User Preferences",
        routerLink: "user-preferences/view",
        name: "user-preferences",
      },
      {
        label: "API Keys",
        routerLink: "api-keys/view",
        name: "api-keys",
      },
    ];
  }
}
