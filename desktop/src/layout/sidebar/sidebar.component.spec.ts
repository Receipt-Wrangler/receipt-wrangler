import { OverlayContainer } from "@angular/cdk/overlay";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatMenuModule, MatMenuTrigger } from "@angular/material/menu";
import { MatSidenavModule } from "@angular/material/sidenav";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { By } from "@angular/platform-browser";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { LayoutState } from "src/store/layout.state";
import { DirectivesModule } from "../../directives/directives.module";
import { ApiModule } from "../../open-api";
import { AuthState, FeatureConfigState, GroupState } from "../../store";
import { SetAuthState, SetPermissions } from "../../store/auth.state.actions";
import { SidebarComponent } from "./sidebar.component";

describe("SidebarComponent", () => {
  let component: SidebarComponent;
  let fixture: ComponentFixture<SidebarComponent>;
  let store: Store;
  let overlayContainer: OverlayContainer;

  const login = () =>
    store.dispatch(
      new SetAuthState({ exp: Math.floor(Date.now() / 1000) + 3600 } as any)
    );

  const openMenu = async () => {
    fixture.detectChanges();
    await fixture.whenStable();
    const trigger = fixture.debugElement
      .query(By.directive(MatMenuTrigger))
      .injector.get(MatMenuTrigger);
    trigger.openMenu();
    fixture.detectChanges();
    await fixture.whenStable();
  };

  const menuText = (): string =>
    overlayContainer.getContainerElement().textContent ?? "";

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [SidebarComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [NgxsModule.forRoot([AuthState, LayoutState, GroupState, FeatureConfigState]),
        DirectivesModule,
        MatMenuModule,
        MatSnackBarModule,
        MatSidenavModule,
        NoopAnimationsModule,
        RouterTestingModule,
        SharedUiModule,
        ApiModule],
    providers: [provideHttpClient(withInterceptorsFromDi())]
}).compileComponents();

    store = TestBed.inject(Store);
    overlayContainer = TestBed.inject(OverlayContainer);
    fixture = TestBed.createComponent(SidebarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("hides the admin manage menu items without the matching read permissions", async () => {
    login();
    await openMenu();

    const text = menuText();
    expect(text).not.toContain("Manage Categories");
    expect(text).not.toContain("Manage Tags");
    expect(text).not.toContain("Manage Groups");
    expect(text).not.toContain("Manage Custom Fields");
  });

  it("shows the admin manage menu items with the matching read permissions", async () => {
    login();
    store.dispatch(
      new SetPermissions(
        [
          "app.categories.read",
          "app.tags.read",
          "app.groups.read",
          "app.custom-fields.read",
        ],
        {}
      )
    );
    await openMenu();

    const text = menuText();
    expect(text).toContain("Manage Categories");
    expect(text).toContain("Manage Tags");
    expect(text).toContain("Manage Groups");
    expect(text).toContain("Manage Custom Fields");
  });

  it("hides User Settings without any settings read permission", async () => {
    login();
    await openMenu();

    expect(menuText()).not.toContain("User Settings");
  });

  it("shows User Settings with a settings read permission", async () => {
    login();
    store.dispatch(new SetPermissions(["app.account.read"], {}));
    await openMenu();

    expect(menuText()).toContain("User Settings");
  });
});
