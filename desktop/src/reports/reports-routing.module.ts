import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { ReportBuilderComponent } from "./report-builder/report-builder.component";

const routes: Routes = [{ path: "", component: ReportBuilderComponent }];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ReportsRoutingModule {}
