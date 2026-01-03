import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { Store } from '@ngxs/store';
import { AuthState } from '../store/auth.state';

@Directive({
    selector: '[appRole]',
    standalone: false
})
export class RoleDirective {
  private hasView = false;

  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
    private store: Store
  ) {}

  @Input() set appRole(role: string) {
    const hasRole = this.store.selectSnapshot(AuthState.hasRole(role));

    if (hasRole) {
      this.viewContainer.createEmbeddedView(this.templateRef);
      this.hasView = true;
    } else if (!hasRole) {
      this.viewContainer.clear();
      this.hasView = false;
    }
  }
}
