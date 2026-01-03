import { Component, OnInit } from "@angular/core";
import { TabConfig } from "../../shared-ui/tabs/tab-config.interface";

@Component({
    selector: "app-system-settings",
    templateUrl: "./system-settings.component.html",
    styleUrl: "./system-settings.component.scss",
    standalone: false
})
export class SystemSettingsComponent implements OnInit {
  public tabs: TabConfig[] = [];

  public ngOnInit(): void {
    this.initTabs();
  }

  private initTabs(): void {
    this.tabs = [
      {
        label: "System Settings",
        routerLink: "settings/view",
        name: "settings",
      },
      {
        label: "Receipt Processing Settings",
        routerLink: "receipt-processing-settings",
        name: "receipt-processing-settings",
      },
      {
        label: "Prompts",
        routerLink: "prompts",
        name: "prompts",
      },
      {
        label: "System Emails",
        routerLink: "system-emails",
        name: "system-emails",
      },
      {
        label: "System Tasks",
        routerLink: "system-tasks",
        name: "system-tasks",
      },
    ];
  }
}
