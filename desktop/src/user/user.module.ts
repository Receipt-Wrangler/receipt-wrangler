import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatTableModule } from "@angular/material/table";
import { CheckboxModule } from "src/checkbox/checkbox.module";
import { PipesModule } from "src/pipes/pipes.module";
import { SelectModule } from "src/select/select.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { ButtonModule } from "../button";
import { InputModule } from "../input";
import { DummyUserConversionDialogComponent } from "./dummy-user-conversion-dialog/dummy-user-conversion-dialog.component";
import { ResetPasswordComponent } from "./reset-password/reset-password.component";
import { UserFormComponent } from "./user-form/user-form.component";
import { UserListComponent } from "./user-list/user-list.component";
import { UserRoutingModule } from "./user-routing.module";

@NgModule({
  declarations: [
    UserListComponent,
    UserFormComponent,
    ResetPasswordComponent,
    DummyUserConversionDialogComponent,
  ],
  imports: [
    ButtonModule,
    CheckboxModule,
    CommonModule,
    InputModule,
    MatDialogModule,
    MatTableModule,
    PipesModule,
    ReactiveFormsModule,
    SelectModule,
    SharedUiModule,
    TableModule,
    UserRoutingModule,
  ],
})
export class UserModule {}
