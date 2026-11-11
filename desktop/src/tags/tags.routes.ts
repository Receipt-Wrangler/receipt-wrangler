import { Routes } from "@angular/router";
import { TagTableComponent } from "./tag-table/tag-table.component";

export const tagsRoutes: Routes = [
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
