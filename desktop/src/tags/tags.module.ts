import { CommonModule } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { DirectivesModule } from "../directives";
import { InputModule } from "../input";
import { PipesModule } from "../pipes";
import { TagFormComponent } from "./tag-form/tag-form.component";
import { TagTableComponent } from "./tag-table/tag-table.component";
import { TagsRoutingModule } from "./tags-routing.module";

@NgModule({ declarations: [TagTableComponent, TagFormComponent], imports: [CommonModule,
        DirectivesModule,
        InputModule,
        PipesModule,
        ReactiveFormsModule,
        SharedUiModule,
        TableModule,
        TagsRoutingModule], providers: [provideHttpClient(withInterceptorsFromDi())] })
export class TagsModule {}
