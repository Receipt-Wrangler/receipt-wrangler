import { computed, Directive, effect, inject, input, TemplateRef, ViewContainerRef } from "@angular/core";
import { Store } from "@ngxs/store";
import { hasAll } from "../utils/permission.utils";
import { AuthState } from "../store/auth.state";

/**
 * Structural directive that renders its template only while the current user
 * holds the given app-scoped permission.
 *
 * Signal/effect-driven so it re-renders when the permission list arrives after
 * first paint (e.g. AppData lands post-login) — unlike a one-shot
 * `selectSnapshot` in an input setter, which never updates.
 */
@Directive({
  selector: "[hasAppPermission]",
  standalone: false,
})
export class HasAppPermissionDirective {
  private readonly templateRef = inject(TemplateRef<unknown>);
  private readonly viewContainer = inject(ViewContainerRef);
  private readonly store = inject(Store);

  private readonly appPermissions = this.store.selectSignal(AuthState.appPermissions);

  public readonly permission = input.required<string>({ alias: "hasAppPermission" });

  private readonly granted = computed(() =>
    hasAll(this.appPermissions(), this.permission())
  );

  constructor() {
    effect(() => {
      this.viewContainer.clear();
      if (this.granted()) {
        this.viewContainer.createEmbeddedView(this.templateRef);
      }
    });
  }
}
