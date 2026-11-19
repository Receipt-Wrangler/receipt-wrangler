import { Routes } from "@angular/router";

import { CategoryTableComponent } from "./category-table/category-table.component";

export const categoriesRoutes: Routes = [
  {
    path: "",
    component: CategoryTableComponent,
  },
  {
    path: "",
    redirectTo: "categories",
    pathMatch: "full",
  },
];
