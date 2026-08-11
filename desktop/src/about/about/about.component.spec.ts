import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideStore, Store } from "@ngxs/store";
import * as QRCode from "qrcode";

jest.mock("qrcode", () => ({
  __esModule: true,
  toDataURL: jest.fn(),
}));

import { AboutComponent } from "./about.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { AboutState } from "../../store/about.state";
import { FeatureConfigState } from "../../store/feature-config.state";
import { SetFeatureConfig } from "../../store/feature-config.state.actions";

describe("AboutComponent", () => {
  let component: AboutComponent;
  let fixture: ComponentFixture<AboutComponent>;

  beforeEach(async () => {
    (QRCode.toDataURL as jest.Mock).mockReset();

    await TestBed.configureTestingModule({
      imports: [AboutComponent],
      providers: [
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
        provideStore([AboutState, FeatureConfigState])
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(AboutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("does not render the mobile app section when the login QR is not configured", () => {
    // Default featureConfig carries an empty loginQrUrl.
    expect(fixture.nativeElement.textContent).not.toContain("Mobile App");
    expect(fixture.nativeElement.querySelector("app-login-qr")).toBeNull();
  });

  it("renders the mobile app setup QR when loginQrUrl is set", async () => {
    (QRCode.toDataURL as jest.Mock).mockResolvedValue(
      "data:image/png;base64,QR"
    );

    TestBed.inject(Store).dispatch(
      new SetFeatureConfig({
        enableLocalSignUp: false,
        aiPoweredReceipts: false,
        loginQrUrl: "https://receiptwrangler.io/app/setup#url=x",
      })
    );
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain("Mobile App");
    expect(
      fixture.nativeElement.querySelector("img.login-qr-code")
    ).toBeTruthy();
  });
});
