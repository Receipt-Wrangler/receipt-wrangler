import { Routes } from "@angular/router";
import { CustomFieldTableComponent } from "./custom-field-table/custom-field-table.component";

export const customFieldsRoutes: Routes = [
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
