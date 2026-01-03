import { Injectable } from "@angular/core";
import { ActivatedRouteSnapshot, RouterStateSnapshot } from "@angular/router";
import { Store } from "@ngxs/store";
import { AuthState } from "../store";

@Injectable({
  providedIn: "root",
})
export class RoleGuard {
  constructor(private store: Store) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean {
    return this.store.selectSnapshot(AuthState.hasRole(route.data["role"]));
  }
}
