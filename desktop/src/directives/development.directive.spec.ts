import { TestBed } from '@angular/core/testing';
import { DevelopmentDirective } from './development.directive';
import { EnvironmentService } from 'src/services/environment.service';
import { TemplateRef, ViewContainerRef } from '@angular/core';

describe('DevelopmentDirective', () => {
  let directive: DevelopmentDirective;
  let environmentService: EnvironmentService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [DevelopmentDirective],
      providers: [
        EnvironmentService,
        DevelopmentDirective,
        {
          provide: TemplateRef,
          useValue: {},
        },
        {
          provide: ViewContainerRef,
          useValue: {
            clear: () => {},
            createEmbeddedView: () => {},
          },
        },
      ],
    });
  });

  beforeEach(() => {
    environmentService = TestBed.inject(EnvironmentService);
  });

  it('should create an instance', () => {
    directive = TestBed.inject(DevelopmentDirective);
    expect(directive).toBeTruthy();
  });

  it('should remove view if is prod', () => {
    jest.spyOn(environmentService, 'isProduction').mockReturnValue(true);
    directive = TestBed.inject(DevelopmentDirective);

    expect(directive.hasView).toEqual(false);
  });

  it('should display view if is not prod', () => {
    jest.spyOn(environmentService, 'isProduction').mockReturnValue(false);
    directive = TestBed.inject(DevelopmentDirective);

    expect(directive.hasView).toEqual(true);
  });
});
