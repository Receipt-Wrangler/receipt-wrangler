import { Component, OnChanges, OnInit, signal, SimpleChanges, input } from "@angular/core";
import { ActivatedRoute, NavigationEnd, Router, RouterLinkActive, RouterLink, RouterOutlet } from "@angular/router";
import { tap } from "rxjs";
import { TabConfig } from "./tab-config.interface";
import { MatTabNav, MatTabLink, MatTabNavPanel } from "@angular/material/tabs";

@Component({
    selector: "app-tabs",
    templateUrl: "./tabs.component.html",
    styleUrls: ["./tabs.component.scss"],
    imports: [MatTabNav, MatTabLink, RouterLinkActive, RouterLink, MatTabNavPanel, RouterOutlet]
})
export class TabsComponent implements OnInit, OnChanges {
  public readonly tabs = input<TabConfig[]>([]);

  public activeName = signal("");

  constructor(private router: Router, private activatedRoute: ActivatedRoute) {}

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["tabs"]) {
      this.setActiveTab();
    }
  }

  public ngOnInit(): void {
    this.listenToRouteChanges();
  }

  private setActiveTab(): void {
    const activeTabName = this.activatedRoute.snapshot.queryParams["tab"];

    const tabs = this.tabs();
    this.activeName.set(tabs.find(t => t.name === activeTabName)?.name || tabs[0].name);
  }

  private listenToRouteChanges(): void {
    this.router.events
      .pipe(
        tap((event) => {
          if (event instanceof NavigationEnd) {
            this.setActiveTab();
          }
        })
      )
      .subscribe();
  }
}
