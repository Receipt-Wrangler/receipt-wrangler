import { Injectable } from "@angular/core";
import { ActivatedRouteSnapshot, RouterStateSnapshot } from "@angular/router";
import { Store } from "@ngxs/store";
import { FeatureConfigState } from "../store";

@Injectable({
  providedIn: "root",
})
export class FeatureGuard {
  constructor(private store: Store) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean {
    return this.store.selectSnapshot(
      FeatureConfigState.hasFeature(route.data["feature"])
    );
  }
}
