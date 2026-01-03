import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { map, Observable } from "rxjs";
import { CustomField, CustomFieldService } from "../open-api/index";

export const customFieldResolverFn: ResolveFn<Observable<CustomField[]>> = (route, state) => {
  const customFieldService = inject(CustomFieldService);

  return customFieldService.getPagedCustomFields({
    page: 1,
    pageSize: -1,
    orderBy: "name",
    sortDirection: "desc"
  }).pipe(
    map((data) => data.data)
  );
};
