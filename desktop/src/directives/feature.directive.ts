import { Directive, TemplateRef, ViewContainerRef, effect, input } from '@angular/core';
import { Store } from '@ngxs/store';
import { FeatureConfigState, isFeatureEnabled } from '../store/feature-config.state';

/**
 * Renders the template content only while the named feature is enabled, and
 * follows the feature config live: an admin toggling a feature in System
 * Settings shows or hides always-mounted content without a reload.
 */
@Directive({ selector: '[appFeature]' })
export class FeatureDirective {
  public readonly appFeature = input.required<string>();

  // One signal for the directive's lifetime, tracked by the effect below, so a
  // store update re-evaluates without building a selector per run.
  private readonly featureConfig = this.store.selectSignal(FeatureConfigState.featureConfig);

  private hasView = false;

  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
    private store: Store
  ) {
    effect(() => {
      const hasFeature = isFeatureEnabled(this.featureConfig(), this.appFeature());

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
