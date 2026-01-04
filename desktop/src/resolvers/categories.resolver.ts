import {inject, Injectable} from "@angular/core";
import {ActivatedRouteSnapshot, ResolveFn, RouterStateSnapshot} from "@angular/router";
import { Observable } from "rxjs";
import { Category, CategoryService } from "../open-api";

export const categoryResolverFn: ResolveFn<Category[]> = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const categoryService = inject(CategoryService);
  return categoryService.getAllCategories();
}
