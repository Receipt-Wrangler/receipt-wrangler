import { CommonModule, NgOptimizedImage } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatProgressSpinner } from "@angular/material/progress-spinner";
import { ButtonModule } from "../button/button.module";
import { DirectivesModule } from "../directives/directives.module";
import { InputModule } from "../input/input.module";
import { PipesModule } from "../pipes/pipes.module";
import { AuthRoutingModule } from "./auth-routing.module";
import { AuthForm } from "./sign-up/auth-form.component";

@NgModule({
  declarations: [AuthForm],
  imports: [
    AuthRoutingModule,
    ButtonModule,
    CommonModule,
    DirectivesModule,
    InputModule,
    PipesModule,
    ReactiveFormsModule,
    NgOptimizedImage,
    MatProgressSpinner,
  ],
  exports: [AuthForm],
})
export class AuthModule {
}
