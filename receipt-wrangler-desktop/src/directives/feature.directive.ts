import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { Store } from '@ngxs/store';
import { FeatureConfigState } from '../store/feature-config.state';

/**
 * Add the template content to the DOM unless the condition is true.
 */
@Directive({
    selector: '[appFeature]',
    standalone: false
})
export class FeatureDirective {
  private hasView = false;

  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
    private store: Store
  ) {}

  @Input() set appFeature(feature: string) {
    const hasFeature = this.store.selectSnapshot(
      FeatureConfigState.hasFeature(feature)
    );

    if (hasFeature) {
      this.viewContainer.createEmbeddedView(this.templateRef);
      this.hasView = true;
    } else if (!hasFeature) {
      this.viewContainer.clear();
      this.hasView = false;
    }
  }
}
