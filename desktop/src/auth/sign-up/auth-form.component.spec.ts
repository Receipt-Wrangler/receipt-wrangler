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
import { AuthForm } from './auth-form.component';
import { AuthState } from '../../store/auth.state';
import { FeatureConfigState } from '../../store/feature-config.state';
import { AuthFormUtil } from './auth-form.util';
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
        RouterTestingModule],
    providers: [
        SnackbarService,
        AppInitService,
        AuthFormUtil,
        {
            provide: ActivatedRoute,
            useValue: {
                data: of(undefined),
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

  it('does not render the login QR when loginQrUrl is empty', () => {
    // Default featureConfig carries an empty loginQrUrl.
    expect(component.qrDataUrl()).toBeNull();
    expect(
      fixture.nativeElement.querySelector('img.login-qr-code')
    ).toBeNull();
  });

  it('renders a locally-generated QR when loginQrUrl is set', async () => {
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
    expect(component.qrDataUrl()).toBe('data:image/png;base64,QR');
    expect(
      fixture.nativeElement.querySelector('img.login-qr-code')
    ).toBeTruthy();
  });

  it('ignores a stale QR generation that resolves after a newer one', async () => {
    // Hand out manually-resolved promises so the two generations can be
    // completed out of order.
    const resolvers: Array<(dataUrl: string) => void> = [];
    (QRCode.toDataURL as jest.Mock).mockImplementation(
      () => new Promise<string>((resolve) => resolvers.push(resolve))
    );
    const store = TestBed.inject(Store);

    const setLoginQrUrl = (loginQrUrl: string) => {
      store.dispatch(
        new SetFeatureConfig({
          enableLocalSignUp: false,
          aiPoweredReceipts: false,
          loginQrUrl,
        })
      );
      fixture.detectChanges();
    };

    setLoginQrUrl('https://receiptwrangler.io/app/setup#url=stale');
    setLoginQrUrl('https://receiptwrangler.io/app/setup#url=current');
    expect(resolvers.length).toBe(2);

    // Reverse order: the current generation finishes first, then the stale one.
    resolvers[1]('data:image/png;base64,CURRENT');
    resolvers[0]('data:image/png;base64,STALE');
    await fixture.whenStable();

    expect(component.qrDataUrl()).toBe('data:image/png;base64,CURRENT');
  });
});
