import { Component, signal } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NgxsModule, Store } from '@ngxs/store';
import { FeatureConfigState } from '../store/feature-config.state';
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
});
