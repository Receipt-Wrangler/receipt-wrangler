import { inject } from '@angular/core';
import { CanActivateFn } from '@angular/router';
import { EnvironmentService } from 'src/services/environment.service';

export const developmentGuard: CanActivateFn = (route, state) => {
  const environmentService = inject(EnvironmentService);
  return !environmentService.isProduction();
};
