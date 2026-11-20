import { Directive, TemplateRef, ViewContainerRef, effect, input, untracked } from '@angular/core';
import { Store } from '@ngxs/store';
import { FeatureConfigState } from '../store/feature-config.state';

/**
 * Add the template content to the DOM unless the condition is true.
 */
@Directive({ selector: '[appFeature]' })
export class FeatureDirective {
  public readonly appFeature = input.required<string>();

  private hasView = false;

  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
    private store: Store
  ) {
    effect(() => {
      const feature = this.appFeature();
      // Re-evaluated when the feature name changes, not on every store update.
      const hasFeature = untracked(() =>
        this.store.selectSnapshot(FeatureConfigState.hasFeature(feature))
      );

      if (hasFeature && !this.hasView) {
        this.viewContainer.createEmbeddedView(this.templateRef);
        this.hasView = true;
      } else if (!hasFeature && this.hasView) {
        this.viewContainer.clear();
        this.hasView = false;
      }
    });
  }
}
