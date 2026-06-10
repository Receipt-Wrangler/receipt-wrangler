import { computed, Directive, effect, inject, input, TemplateRef, ViewContainerRef } from "@angular/core";
import { Store } from "@ngxs/store";
import { hasAll, hasAny } from "../utils/permission.utils";
import { AuthState } from "../store/auth.state";

export interface HasGroupPermissionInput {
  groupId: number;
  permission: string;
  /**
   * App-scoped permissions that, if any are held, bypass the group check
   * (mirrors the backend `OrAppPermissions` admin-not-a-member pattern).
   */
  orApp?: string[];
}

/**
 * Structural directive that renders its template only while the current user
 * holds the given group-scoped permission for the given group (or one of the
 * optional `orApp` app-scoped fallbacks).
 *
 * Signal/effect-driven so it re-renders when the permission lists arrive after
 * first paint (e.g. AppData lands post-login).
 */
@Directive({
  selector: "[hasGroupPermission]",
  standalone: false,
})
export class HasGroupPermissionDirective {
  private readonly templateRef = inject(TemplateRef<unknown>);
  private readonly viewContainer = inject(ViewContainerRef);
  private readonly store = inject(Store);

  private readonly appPermissions = this.store.selectSignal(AuthState.appPermissions);
  private readonly groupPermissions = this.store.selectSignal(AuthState.groupPermissions);

  public readonly config = input.required<HasGroupPermissionInput>({
    alias: "hasGroupPermission",
  });

  private readonly granted = computed(() => {
    const { groupId, permission, orApp = [] } = this.config();
    if (orApp.length > 0 && hasAny(this.appPermissions(), ...orApp)) {
      return true;
    }
    return hasAll(this.groupPermissions()?.[groupId] ?? [], permission);
  });

  constructor() {
    effect(() => {
      this.viewContainer.clear();
      if (this.granted()) {
        this.viewContainer.createEmbeddedView(this.templateRef);
      }
    });
  }
}
