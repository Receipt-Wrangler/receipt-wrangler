import { Component, OnInit, ViewEncapsulation, computed, input, signal } from "@angular/core";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import {
  Activity,
  Group,
  PagedActivityRequestCommand,
  PagedDataDataInner,
  Permission,
  SystemTaskService,
  SystemTaskStatus,
  Widget
} from "../../open-api/index";
import { SnackbarService } from "../../services/index";
import { GroupState } from "../../store/index";

@Component({
  selector: "app-activity",
  templateUrl: "./activity.component.html",
  styleUrl: "./activity.component.scss",
  encapsulation: ViewEncapsulation.None,
  standalone: false
})
export class ActivityComponent implements OnInit {
  public readonly widget = input.required<Widget>();

  public readonly groupId = input<number>();

  public group = computed(() => {
    const id = this.groupId();
    return id ? this.store.selectSnapshot(GroupState.getGroupById(id.toString())) : undefined;
  });

  public page: number = 1;

  public pageSize: number = 25;

  public activities = signal<PagedDataDataInner[]>([]);

  public ranActivities = signal<{ [key: number]: boolean }>({});

  protected readonly SystemTaskStatus = SystemTaskStatus;

  protected readonly Permission = Permission;

  constructor(
    private systemTaskService: SystemTaskService,
    private snackbarService: SnackbarService,
    private store: Store
  ) {}

  public ngOnInit(): void {
    this.getData();
  }

  public endOfListReached(): void {
    this.page++;
    this.getData();
  }

  public onRefreshButtonClick(id: number): void {
    this.systemTaskService
      .rerunActivity(id)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Activity has been successfully queued.");
          this.ranActivities.update(prev => ({ ...prev, [id]: true }));
        })
      ).subscribe();
  }

  private getData(): void {
    if (!this.groupId()) {
      return;
    }

    const command: PagedActivityRequestCommand = {
      groupIds: this.getGroupIds(),
      orderBy: "started_at",
      page: this.page,
      pageSize: this.pageSize,
      sortDirection: "desc"
    };
    this.systemTaskService.getPagedActivities(command)
      .pipe(
        take(1),
        tap((response) => {
          this.activities.update(prev => [...prev, ...response.data]);
        })
      )
      .subscribe();
  }

  private getGroupIds(): number[] {
    if (this.group()?.isAllGroup) {
      return this.store.selectSnapshot(GroupState.groupsWithoutAll).map((group) => group.id);
    } else {
      return [this.groupId() ?? 0];
    }
  }

  public buildItemRouterLinkString(item: Activity): string {
    if (!item?.receiptId) {
      return "";
    }
    return `/receipts/${item.receiptId}/view`;
  }
}
