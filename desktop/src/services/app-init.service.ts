import { Injectable } from "@angular/core";
import { Store } from "@ngxs/store";
import { catchError, finalize, Observable, of, switchMap, take, tap, } from "rxjs";
import { AppData, AuthService, FeatureConfigService, UserService, } from "../open-api";
import { AuthState, SetFeatureConfig } from "../store";
import { setAppData } from "../utils";

@Injectable({
  providedIn: "root",
})
export class AppInitService {
  constructor(
    private authService: AuthService,
    private store: Store,
    private userService: UserService,
    private featureConfigService: FeatureConfigService
  ) {}

  public initAppData(): Promise<boolean> {

    return new Promise((resolve) => {
      const isLoggedIn = this.store.selectSnapshot(AuthState.isLoggedIn);

      if (!isLoggedIn) {
        this.featureConfigService.getFeatureConfig().pipe(
          take(1),
          switchMap((config) => this.store.dispatch(new SetFeatureConfig(config))),
          finalize(() => {
            resolve(true);
          })
        ).subscribe();
        return;
      }

      this.authService.getNewRefreshToken().pipe(
        take(1),
        switchMap(() => this.getAppData()),
        tap(() => resolve(true)),
        catchError((err) => {
          resolve(false);
          return of(err);
        })
      )
        .subscribe();
    });
  }

  public getAppData(): Observable<any[]> {
    return this.userService.getAppData().pipe(switchMap((appData: AppData) => setAppData(this.store, appData)));
  }

}

export function initAppData(appInitService: AppInitService) {
  return async () => await appInitService.initAppData();
}
