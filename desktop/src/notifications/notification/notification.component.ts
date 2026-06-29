import { Component, OnInit, input, output } from "@angular/core";
import { take, tap } from "rxjs";
import { ParameterizedDataParser } from "src/utils";
import { Notification, NotificationsService, Permission } from "../../open-api";

@Component({
    selector: "app-notification",
    templateUrl: "./notification.component.html",
    styleUrls: ["./notification.component.scss"],
    standalone: false
})
export class NotificationComponent implements OnInit {
  public readonly notification = input.required<Notification>();

  public readonly notificationDeleted = output<number>();

  protected readonly Permission = Permission;

  public link?: string;

  public parsedBody: string = "";

  constructor(
    private notificationsService: NotificationsService,
    private parameterizedDataParser: ParameterizedDataParser
  ) {}

  public ngOnInit(): void {
    this.parseBody();
  }

  private parseBody(): void {
    this.parsedBody = this.parameterizedDataParser.parse(
      this.notification().body
    );
    this.link = this.parameterizedDataParser.link;
  }

  public deleteNotification(): void {
    this.notificationsService
      .deleteNotificationById(this.notification().id)
      .pipe(
        take(1),
        tap(() => {
          this.notificationDeleted.emit(this.notification().id);
        })
      )
      .subscribe();
  }
}
