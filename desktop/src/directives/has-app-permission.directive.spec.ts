import { Component, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { AuthState } from "../store/auth.state";
import { SetPermissions } from "../store/auth.state.actions";
import { HasAppPermissionDirective } from "./has-app-permission.directive";

@Component({
  template: `<span *hasAppPermission="permission" data-testid="gated">gated</span>`,
  standalone: false,
})
class HostComponent {
  permission = "app.users.read";
}

describe("HasAppPermissionDirective", () => {
  let fixture: ComponentFixture<HostComponent>;
  let store: Store;

  const gated = () => fixture.nativeElement.querySelector("[data-testid='gated']");

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [HostComponent, HasAppPermissionDirective],
      imports: [NgxsModule.forRoot([AuthState])],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(HostComponent);
  });

  it("renders the template when the permission is granted", async () => {
    store.dispatch(new SetPermissions(["app.users.read"], {}));
    await fixture.whenStable();

    expect(gated()).toBeTruthy();
  });

  it("does not render when the permission is not granted", async () => {
    store.dispatch(new SetPermissions(["app.roles.read"], {}));
    await fixture.whenStable();

    expect(gated()).toBeFalsy();
  });

  it("re-renders when permissions are dispatched after the initial render", async () => {
    // Initial paint with no permissions → hidden.
    await fixture.whenStable();
    expect(gated()).toBeFalsy();

    // Permissions arrive late (e.g. AppData lands post-login).
    store.dispatch(new SetPermissions(["app.users.read"], {}));
    TestBed.flushEffects();
    await fixture.whenStable();

    expect(gated()).toBeTruthy();
  });

  it("honors wildcard grants", async () => {
    store.dispatch(new SetPermissions(["*"], {}));
    await fixture.whenStable();

    expect(gated()).toBeTruthy();
  });
});
