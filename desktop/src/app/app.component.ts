import { Component, OnInit, DestroyRef, inject } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { RouterOutlet } from "@angular/router";
import { Store } from "@ngxs/store";
import { interval, take, tap } from "rxjs";
import { HideProgressBar } from "src/store/layout.state.actions";
import { TokenRefreshService } from "../services";
import { AuthState } from "../store";

@Component({
    selector: "app-root",
    templateUrl: "./app.component.html",
    styleUrls: ["./app.component.scss"],
    imports: [RouterOutlet],
})
export class AppComponent implements OnInit {
  private readonly destroyRef = inject(DestroyRef);

  constructor(
    private tokenRefreshService: TokenRefreshService,
    private store: Store
  ) {}

  public ngOnInit(): void {
    this.store.dispatch(new HideProgressBar());
    this.refreshTokens();
  }

  private refreshTokens(): void {
    const fifteenMinutes = 1000 * 60 * 15;
    interval(fifteenMinutes)
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        tap(() => {
          if (this.store.selectSnapshot(AuthState.isLoggedIn)) {
            this.tokenRefreshService.refreshToken().pipe(take(1)).subscribe();
          }
        })
      ).subscribe();
  }
}
