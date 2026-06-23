import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { Store } from "@ngxs/store";
import { map, Observable, of } from "rxjs";
import { CustomField, CustomFieldService, Permission } from "../open-api/index";
import { AuthState } from "../store";

export const customFieldResolverFn: ResolveFn<Observable<CustomField[]>> = (route, state) => {
  const customFieldService = inject(CustomFieldService);
  const store = inject(Store);

  // The custom field catalog is admin-only (app.custom-fields.read). Skip the call
  // for users without it so opening a receipt doesn't 403 (and log them out); the
  // receipt itself carries the definitions for its own custom fields so values
  // still render.
  if (!store.selectSnapshot(AuthState.hasAppPermission(Permission.AppCustomFieldsRead))) {
    return of([]);
  }

  return customFieldService.getPagedCustomFields({
    page: 1,
    pageSize: -1,
    orderBy: "name",
    sortDirection: "desc"
  }).pipe(
    map((data) => data.data)
  );
};
