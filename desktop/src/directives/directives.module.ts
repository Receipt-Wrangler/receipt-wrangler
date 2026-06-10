import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { DevelopmentDirective } from "./development.directive";
import { FeatureDirective } from "./feature.directive";
import { HasAppPermissionDirective } from "./has-app-permission.directive";
import { HasGroupPermissionDirective } from "./has-group-permission.directive";

@NgModule({
  declarations: [
    DevelopmentDirective,
    FeatureDirective,
    HasAppPermissionDirective,
    HasGroupPermissionDirective,
  ],
  exports: [
    DevelopmentDirective,
    FeatureDirective,
    HasAppPermissionDirective,
    HasGroupPermissionDirective,
  ],
  imports: [CommonModule],
})
export class DirectivesModule {}
