import { Component, OnInit, output, signal } from "@angular/core";
import { take, tap } from "rxjs";
import { Notification, NotificationsService } from "../../open-api";
import { ButtonComponent } from "../../button/button/button.component";
import { NotificationComponent } from "../notification/notification.component";

@Component({
    selector: "app-notifications-list",
    templateUrl: "./notifications-list.component.html",
    styleUrls: ["./notifications-list.component.scss"],
    imports: [ButtonComponent, NotificationComponent]
})
export class NotificationsListComponent implements OnInit {
  public notifications = signal<Notification[]>([]);

  public readonly notificationCountChanged = output<number | undefined>();

  constructor(private notificationsService: NotificationsService) {}

  public ngOnInit(): void {
    this.getNotifications();
  }

  private getNotifications(): void {
    this.notificationsService
      .getNotificationsForuser()
      .pipe(
        take(1),
        tap((notifications) => {
          this.notifications.set(notifications);
          this.emitCount();
        })
      )
      .subscribe();
  }

  public deleteAllNotifications(): void {
    this.notificationsService
      .deleteAllNotificationsForUser()
      .pipe(
        take(1),
        tap(() => {
          this.notifications.set([]);
          this.emitCount();
        })
      )
      .subscribe();
  }

  public notificationDeleted(id: number): void {
    this.notifications.set(this.notifications().filter((n) => n.id !== id));
    this.emitCount();
  }

  private emitCount(): void {
    const count = this.notifications().length;
    this.notificationCountChanged.emit(count > 0 ? count : undefined);
  }
}
