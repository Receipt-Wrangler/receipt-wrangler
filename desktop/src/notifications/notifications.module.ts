import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatListModule } from "@angular/material/list";
import { RouterModule } from "@angular/router";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { ButtonModule } from "../button";
import { NotificationComponent } from "./notification/notification.component";
import { NotificationsListComponent } from "./notifications-list/notifications-list.component";

@NgModule({
  declarations: [NotificationsListComponent, NotificationComponent],
  imports: [
    ButtonModule,
    CommonModule,
    MatButtonModule,
    MatListModule,
    RouterModule,
    SharedUiModule,
  ],
  exports: [NotificationsListComponent],
})
export class NotificationsModule {}
