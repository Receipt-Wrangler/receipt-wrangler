import { TestBed } from '@angular/core/testing';
import { CanActivateFn } from '@angular/router';

import { developmentGuard } from './development.guard';
import { EnvironmentService } from 'src/services/environment.service';

describe('developmentGuard', () => {
  let environmentService: EnvironmentService;
  const executeGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => developmentGuard(...guardParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  beforeEach(() => {
    environmentService = TestBed.inject(EnvironmentService);
  });

  it('should be created', () => {
    expect(executeGuard).toBeTruthy();
  });

  it('should return false if is production', () => {
    jest.spyOn(environmentService, 'isProduction').mockReturnValue(true);
    expect(executeGuard({} as any, {} as any)).toEqual(false);
  });

  it('should return true if is not production', () => {
    jest.spyOn(environmentService, 'isProduction').mockReturnValue(false);
    expect(executeGuard({} as any, {} as any)).toEqual(true);
  });
});
