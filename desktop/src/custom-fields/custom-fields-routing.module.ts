import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { CustomFieldTableComponent } from "./custom-field-table/custom-field-table.component";

const routes: Routes = [
  {
    path: "",
    component: CustomFieldTableComponent,
  },
  {
    path: "",
    redirectTo: "custom-fields",
    pathMatch: "full",
  },
];

@NgModule({
  declarations: [],
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CustomFieldsRoutingModule {}
