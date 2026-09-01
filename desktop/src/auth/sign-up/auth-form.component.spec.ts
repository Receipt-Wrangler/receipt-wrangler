import { provideHttpClientTesting } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { ActivatedRoute } from '@angular/router';
import { NgxsModule, Store } from '@ngxs/store';
import { of } from 'rxjs';
import * as QRCode from 'qrcode';
import { ApiModule } from '../../open-api';
import { SetFeatureConfig } from '../../store/feature-config.state.actions';

jest.mock('qrcode', () => ({
  __esModule: true,
  toDataURL: jest.fn(),
}));
import { ButtonModule } from '../../button';
import { FeatureDirective } from '../../directives/feature.directive';
import { InputModule } from '../../input';
import { PipesModule } from '../../pipes/pipes.module';
import { AppInitService } from '../../services/app-init.service';
import { SnackbarService } from '../../services/snackbar.service';
import { AuthForm, buildOidcLoginUrl } from './auth-form.component';
import { AuthState } from '../../store/auth.state';
import { FeatureConfigState } from '../../store/feature-config.state';
import { AuthFormUtil } from './auth-form.util';
import { LoginQrComponent } from '../../shared-ui/login-qr/login-qr.component';
import { RouterTestingModule } from '@angular/router/testing';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

describe('AuthForm', () => {
  let component: AuthForm;
  let fixture: ComponentFixture<AuthForm>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [AuthForm, FeatureDirective],
    imports: [ButtonModule,
        InputModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, FeatureConfigState]),
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule,
        ApiModule,
        LoginQrComponent,
        RouterTestingModule],
    providers: [
        SnackbarService,
        AppInitService,
        AuthFormUtil,
        {
            provide: ActivatedRoute,
            useValue: {
                data: of(undefined),
                snapshot: { queryParams: {} },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    fixture = TestBed.createComponent(AuthForm);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  // The QR itself is unit-tested in shared-ui/login-qr; these two only pin that
  // the login page still wires the shared component up.
  it('does not render the login QR when loginQrUrl is empty', () => {
    // Default featureConfig carries an empty loginQrUrl.
    expect(
      fixture.nativeElement.querySelector('img.login-qr-code')
    ).toBeNull();
  });

  it('renders no OIDC buttons when no providers are configured', () => {
    expect(
      fixture.nativeElement.querySelector('.oidc-section')
    ).toBeNull();
  });

  it('renders one login button per configured provider', () => {
    const store = TestBed.inject(Store);

    store.dispatch(
      new SetFeatureConfig({
        enableLocalSignUp: false,
        aiPoweredReceipts: false,
        oidcProviders: [
          { name: 'google', displayName: 'Google' },
          { name: 'keycloak', displayName: 'Keycloak' },
        ],
      })
    );
    fixture.detectChanges();

    const buttons = fixture.nativeElement.querySelectorAll('.oidc-section app-button');
    expect(buttons.length).toBe(2);
    expect(fixture.nativeElement.textContent).toContain('Log in with Google');
    expect(fixture.nativeElement.textContent).toContain('Log in with Keycloak');
  });

  // A full-page navigation, not an HttpClient call: the browser has to follow
  // the redirect chain to the identity provider itself.
  it('navigates the whole page to the provider login URL', () => {
    const navigateSpy = jest
      .spyOn(component as any, 'navigateToUrl')
      .mockImplementation(() => undefined);

    component.loginWithOidc({ name: 'google', displayName: 'Google' });

    expect(navigateSpy).toHaveBeenCalledWith('/api/oidc/google/login');
  });

  it('escapes a provider name in the login URL', () => {
    expect(buildOidcLoginUrl('my provider/../evil')).toBe(
      '/api/oidc/my%20provider%2F..%2Fevil/login'
    );
  });

  it('surfaces an oidcError query param as a snackbar', () => {
    const snackbarService = TestBed.inject(SnackbarService);
    const errorSpy = jest.spyOn(snackbarService, 'error').mockImplementation(() => undefined as any);
    const route = TestBed.inject(ActivatedRoute) as any;
    route.snapshot.queryParams = { oidcError: 'no_account' };

    const localFixture = TestBed.createComponent(AuthForm);
    localFixture.detectChanges();

    expect(errorSpy).toHaveBeenCalledWith(
      expect.stringContaining('No Receipt Wrangler account is linked')
    );
  });

  it('renders the login QR when loginQrUrl is set', async () => {
    (QRCode.toDataURL as jest.Mock).mockResolvedValue(
      'data:image/png;base64,QR'
    );
    const store = TestBed.inject(Store);
    const loginQrUrl = 'https://receiptwrangler.io/app/setup#url=x';

    store.dispatch(
      new SetFeatureConfig({
        enableLocalSignUp: false,
        aiPoweredReceipts: false,
        loginQrUrl,
      })
    );
    // Effect regenerates the data URL (async) from the new store value.
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(QRCode.toDataURL).toHaveBeenCalledWith(
      loginQrUrl,
      expect.anything()
    );
    expect(
      fixture.nativeElement.querySelector('img.login-qr-code')
    ).toBeTruthy();
    expect(
      fixture.nativeElement.querySelector('.login-qr-divider').textContent
    ).toContain('Set up the mobile app');
  });
});
