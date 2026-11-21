import { Component, signal } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NgxsModule, Store } from '@ngxs/store';
import { FeatureConfigState } from '../store/feature-config.state';
import { SetFeatureConfig } from '../store/feature-config.state.actions';
import { FeatureDirective } from './feature.directive';

@Component({
  imports: [FeatureDirective],
  template: `<span class="gated" *appFeature="feature()">gated</span>`,
})
class HostComponent {
  public readonly feature = signal('enabledOne');
}

describe('FeatureDirective', () => {
  let fixture: ComponentFixture<HostComponent>;

  function renderedCount(): number {
    return fixture.nativeElement.querySelectorAll('.gated').length;
  }

  async function setFeature(feature: string): Promise<void> {
    fixture.componentInstance.feature.set(feature);
    await fixture.whenStable();
  }

  beforeEach(async () => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([FeatureConfigState]), HostComponent],
    });

    const store = TestBed.inject(Store);
    store.reset({
      ...store.snapshot(),
      featureConfig: { enabledOne: true, enabledTwo: true, disabled: false },
    });

    fixture = TestBed.createComponent(HostComponent);
    await fixture.whenStable();
  });

  it('renders the content when the feature is enabled', () => {
    expect(renderedCount()).toEqual(1);
  });

  it('does not render the content when the feature is disabled', async () => {
    await setFeature('disabled');

    expect(renderedCount()).toEqual(0);
  });

  it('renders the content again when the feature changes back to an enabled one', async () => {
    await setFeature('disabled');
    await setFeature('enabledOne');

    expect(renderedCount()).toEqual(1);
  });

  it('keeps a single view when switching between two enabled features', async () => {
    await setFeature('enabledTwo');

    expect(renderedCount()).toEqual(1);
  });

  // #749: the config changes at runtime (an admin saving System Settings), and
  // content that is always mounted, like the sidebar Quick Scan button, has to
  // follow it without the feature name changing.
  describe('when the feature config changes', () => {
    async function setAiPoweredReceipts(enabled: boolean): Promise<void> {
      TestBed.inject(Store).dispatch(new SetFeatureConfig({ aiPoweredReceipts: enabled } as any));
      await fixture.whenStable();
    }

    beforeEach(async () => {
      await setAiPoweredReceipts(true);
      await setFeature('aiPoweredReceipts');
    });

    it('removes the content when the feature is turned off', async () => {
      await setAiPoweredReceipts(false);

      expect(renderedCount()).toEqual(0);
    });

    it('renders the content when the feature is turned on', async () => {
      await setAiPoweredReceipts(false);
      // Re-select the name so the view is hidden whichever way the directive
      // reads the store; only the config change below may bring it back.
      await setFeature('enabledOne');
      await setFeature('aiPoweredReceipts');
      expect(renderedCount()).toEqual(0);

      await setAiPoweredReceipts(true);

      expect(renderedCount()).toEqual(1);
    });

    it('keeps a single view across an update that leaves the feature on', async () => {
      await setAiPoweredReceipts(true);

      expect(renderedCount()).toEqual(1);
    });
  });
});
