import { ScrollingModule } from "@angular/cdk/scrolling";
import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatCardModule } from "@angular/material/card";
import { MatChipsModule } from "@angular/material/chips";
import { MatListModule } from "@angular/material/list";
import { CheckboxModule } from "src/checkbox/checkbox.module";
import { PipesModule } from "src/pipes/pipes.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { AvatarModule } from "../avatar/index";
import { ButtonModule } from "../button/index";
import { InputModule } from "../input";
import { SelectModule } from "../select/select.module";
import { SystemTaskTypePipe } from "../shared-ui/task-table/system-task-type.pipe";
import { DateBlockComponent } from "../standalone/components/date-block/date-block.component";
import { ExportButtonComponent } from "../standalone/components/export-button/export-button.component";
import { ActivityComponent } from "./activity/activity.component";
import { DashboardFormComponent } from "./dashboard-form/dashboard-form.component";
import { DashboardListComponent } from "./dashboard-list/dashboard-list.component";
import { DashboardRoutingModule } from "./dashboard-routing.module";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { FilteredReceiptsComponent } from "./filtered-receipts/filtered-receipts.component";
import { GroupDashboardsComponent } from "./group-dashboards/group-dashboards.component";
import { WidgetTypePipe } from "./widget-type.pipe";

@NgModule({
  declarations: [
    DashboardComponent,
    DashboardFormComponent,
    GroupDashboardsComponent,
    FilteredReceiptsComponent,
    WidgetTypePipe,
    ActivityComponent,
    DashboardListComponent,
  ],
  imports: [
    CheckboxModule,
    CommonModule,
    DashboardRoutingModule,
    InputModule,
    MatCardModule,
    MatChipsModule,
    MatListModule,
    PipesModule,
    PipesModule,
    ReactiveFormsModule,
    ScrollingModule,
    SharedUiModule,
    ButtonModule,
    SelectModule,
    SystemTaskTypePipe,
    AvatarModule,
    DateBlockComponent,
    ExportButtonComponent,
  ],
  exports: [
    WidgetTypePipe
  ],
})
export class DashboardModule {}
