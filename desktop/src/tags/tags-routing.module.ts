import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { TagTableComponent } from "./tag-table/tag-table.component";

const routes: Routes = [
  {
    path: "",
    component: TagTableComponent,
  },
  {
    path: "",
    redirectTo: "tags",
    pathMatch: "full",
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class TagsRoutingModule {}
