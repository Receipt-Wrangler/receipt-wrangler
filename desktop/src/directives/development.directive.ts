import { Directive, TemplateRef, ViewContainerRef, } from "@angular/core";
import { EnvironmentService } from "src/services/environment.service";

@Directive({
    selector: "[appDevelopment]",
    standalone: false
})
export class DevelopmentDirective {
  public hasView = false;

  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
    private environmentService: EnvironmentService
  ) {
    const isProd = this.environmentService.isProduction();

    if (!isProd) {
      this.viewContainer.createEmbeddedView(this.templateRef);
      this.hasView = true;
    } else if (isProd) {
      this.viewContainer.clear();
      this.hasView = false;
    }
  }
}
