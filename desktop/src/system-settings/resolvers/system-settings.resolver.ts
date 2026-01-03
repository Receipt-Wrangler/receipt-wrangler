import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { take } from "rxjs";
import { SystemSettings, SystemSettingsService } from "../../open-api";

export const systemSettingsResolver: ResolveFn<SystemSettings> = (route, state) => {
  const service = inject(SystemSettingsService);
  return service.getSystemSettings().pipe(take(1));
};
