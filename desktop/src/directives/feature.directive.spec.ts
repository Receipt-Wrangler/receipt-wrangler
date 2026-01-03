import { TemplateRef, ViewContainerRef } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { NgxsModule, Store } from '@ngxs/store';
import { FeatureConfigState } from '../store/feature-config.state';
import { FeatureDirective } from './feature.directive';

describe('FeatureDirective', () => {
  let store: Store;
  let templateRef: jest.Mocked<TemplateRef<any>>;
  let viewContainer: jest.Mocked<ViewContainerRef>;

  beforeEach(() => {
    const templateSpy = { createEmbeddedView: jest.fn() };
    const viewContainerSpy = { createEmbeddedView: jest.fn(), clear: jest.fn() };

    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([FeatureConfigState])],
      providers: [
        { provide: TemplateRef, useValue: templateSpy },
        { provide: ViewContainerRef, useValue: viewContainerSpy }
      ]
    });

    store = TestBed.inject(Store);
    templateRef = TestBed.inject(TemplateRef) as jest.Mocked<TemplateRef<any>>;
    viewContainer = TestBed.inject(ViewContainerRef) as jest.Mocked<ViewContainerRef>;
  });

  it('should create an instance', () => {
    const directive = new FeatureDirective(templateRef, viewContainer, store);
    expect(directive).toBeTruthy();
  });
});
