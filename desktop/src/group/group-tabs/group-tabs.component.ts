import { Component } from "@angular/core";
import { Store } from "@ngxs/store";
import { TabConfig } from "src/shared-ui/tabs/tab-config.interface";
import { FeatureConfigState } from "../../store";

@Component({
    selector: "app-group-tabs",
    templateUrl: "./group-tabs.component.html",
    styleUrls: ["./group-tabs.component.scss"],
    standalone: false
})
export class GroupTabsComponent {
  public tabs: TabConfig[] = [];

  public ngOnInit(): void {
    this.initTabs();
  }

  constructor(private store: Store) {}

  private initTabs(): void {
    this.tabs = [
      {
        label: "Group Details",
        routerLink: "details/view",
        name: "details",
      },
      {
        label: "Group Receipt Settings",
        routerLink: "receipt-settings/view",
        name: "receipt-settings",
      },
    ];

    const hasAiPoweredReceipts = this.store.selectSnapshot(
      FeatureConfigState.hasFeature("aiPoweredReceipts")
    );
    if (hasAiPoweredReceipts) {
      this.tabs.push({
        label: "Group AI Settings",
        routerLink: "settings/view",
        name: "settings",
      });
    }
  }
}
