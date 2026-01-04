import { CdkVirtualScrollViewport } from "@angular/cdk/scrolling";
import { AfterViewInit, Component, EventEmitter, Input, Output, TemplateRef, ViewChild, ViewEncapsulation } from "@angular/core";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { tap } from "rxjs";

@UntilDestroy()
@Component({
    selector: "app-dashboard-list",
    templateUrl: "./dashboard-list.component.html",
    styleUrl: "./dashboard-list.component.scss",
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class DashboardListComponent implements AfterViewInit {
  @ViewChild(CdkVirtualScrollViewport) public cdkVirtualScrollViewport!: CdkVirtualScrollViewport;

  @Input() public itemHeaderTemplate!: TemplateRef<any>;

  @Input() public itemLineTemplate!: TemplateRef<any>;

  @Input() public itemLineTemplate2!: TemplateRef<any>;

  @Input() public itemAvatarTemplate!: TemplateRef<any>;

  @Input() public itemMetaTemplate!: TemplateRef<any>;

  @Input() public items: any[] = [];

  @Input() public noItemFoundText = "";

  @Input() public itemSize = 67;

  @Input() public buildRouterLinkString: (item: any) => string = (item: any) => "";

  @Output() public endOfListReached = new EventEmitter<void>();

  public ngAfterViewInit(): void {
    this.listenForEndOfList();
  }

  private listenForEndOfList(): void {
    this.cdkVirtualScrollViewport.renderedRangeStream
      .pipe(
        untilDestroyed(this),
        tap((range) => {
          if (range.end === this.items.length) {
            this.endOfListReached.emit();
          }
        })
      )
      .subscribe();
  }
}
