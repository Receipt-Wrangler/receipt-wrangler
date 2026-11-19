import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { ApiModule, OidcProviderService, Permission } from "../../open-api";
import { SnackbarService } from "../../services/snackbar.service";
import { AuthState } from "../../store/auth.state";
import { SetPermissions } from "../../store/auth.state.actions";
import { OidcProviderTableState } from "../../store/oidc-provider-table.state";
import { OidcProviderTableComponent } from "./oidc-provider-table.component";

describe("OidcProviderTableComponent", () => {
  let component: OidcProviderTableComponent;
  let fixture: ComponentFixture<OidcProviderTableComponent>;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        ApiModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, OidcProviderTableState]),
        NoopAnimationsModule,
        RouterTestingModule,
        OidcProviderTableComponent,
      ],
      providers: [
        SnackbarService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    store = TestBed.inject(Store);

    jest
      .spyOn(TestBed.inject(OidcProviderService), "getPagedOidcProviders")
      .mockReturnValue(of({ data: [], totalCount: 0 }) as any);

    fixture = TestBed.createComponent(OidcProviderTableComponent);
    component = fixture.componentInstance;
  });

  afterEach(() => jest.restoreAllMocks());

  it("hides the actions column when the user can neither edit nor delete", () => {
    store.dispatch(new SetPermissions([Permission.AppOidcProvidersRead], {}));
    fixture.detectChanges(false);

    expect(component.displayedColumns).not.toContain("actions");
  });

  // Gating on delete alone would hide the whole column, and with it the edit
  // button, from an update-only holder.
  it("shows the actions column for an update-only holder", () => {
    store.dispatch(new SetPermissions([Permission.AppOidcProvidersUpdate], {}));
    fixture.detectChanges(false);

    expect(component.displayedColumns).toContain("actions");
  });

  it("shows the actions column for a delete-only holder", () => {
    store.dispatch(new SetPermissions([Permission.AppOidcProvidersDelete], {}));
    fixture.detectChanges(false);

    expect(component.displayedColumns).toContain("actions");
  });
});
