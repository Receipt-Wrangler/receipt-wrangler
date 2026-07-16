import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { ReportBuilderComponent } from "./report-builder/report-builder.component";
import { reportTemplateResolver } from "./report-builder/report-template.resolver";
import { ReportTemplateListComponent } from "./report-template-list/report-template-list.component";

const routes: Routes = [
  // The list is the landing page; the builder is reached via "New Report" or a
  // template's edit/open action. Only the builder routes opt into fullHeight (the
  // shell's bounded, no-padding two-pane frame — see SidebarComponent); the list is
  // a normal padded, scrolling page.
  { path: "", component: ReportTemplateListComponent },
  { path: "new", component: ReportBuilderComponent, data: { fullHeight: true } },
  {
    path: ":id/edit",
    component: ReportBuilderComponent,
    data: { fullHeight: true },
    resolve: { template: reportTemplateResolver },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ReportsRoutingModule {}
