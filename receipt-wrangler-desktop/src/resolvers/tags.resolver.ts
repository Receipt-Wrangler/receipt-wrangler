import { inject } from "@angular/core";
import { ActivatedRouteSnapshot, RouterStateSnapshot } from "@angular/router";
import { Observable } from "rxjs";
import { Tag, TagService } from "../open-api";

export const tagResolverFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
): Tag[] | Observable<Tag[]> | Promise<Tag[]> => {
  const tagService = inject(TagService);
  return tagService.getAllTags();
};
