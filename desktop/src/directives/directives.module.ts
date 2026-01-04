import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { DevelopmentDirective } from "./development.directive";
import { FeatureDirective } from "./feature.directive";
import { RoleDirective } from "./role.directive";

@NgModule({
  declarations: [DevelopmentDirective, FeatureDirective, RoleDirective],
  exports: [DevelopmentDirective, FeatureDirective, RoleDirective],
  imports: [CommonModule],
})
export class DirectivesModule {}
