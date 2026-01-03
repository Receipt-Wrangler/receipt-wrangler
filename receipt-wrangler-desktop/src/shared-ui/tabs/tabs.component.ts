import { Component, Input, OnChanges, OnInit, SimpleChanges, } from "@angular/core";
import { ActivatedRoute, NavigationEnd, Router } from "@angular/router";
import { tap } from "rxjs";
import { TabConfig } from "./tab-config.interface";

@Component({
    selector: "app-tabs",
    templateUrl: "./tabs.component.html",
    styleUrls: ["./tabs.component.scss"],
    standalone: false
})
export class TabsComponent implements OnInit, OnChanges {
  @Input() public tabs: TabConfig[] = [];

  public activeName: string = "";

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

    this.activeName = this.tabs.find(t => t.name === activeTabName)?.name || this.tabs[0].name;
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
