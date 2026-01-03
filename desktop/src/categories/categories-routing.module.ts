import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";

import { CategoryTableComponent } from "./category-table/category-table.component";

const routes: Routes = [
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

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class CategoriesRoutingModule {}
