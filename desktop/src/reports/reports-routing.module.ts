import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { ReportBuilderComponent } from "./report-builder/report-builder.component";

const routes: Routes = [
  // fullHeight opts this route into the shell's bounded, no-padding frame so the
  // builder's two panes scroll independently (see SidebarComponent).
  { path: "", component: ReportBuilderComponent, data: { fullHeight: true } },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ReportsRoutingModule {}
