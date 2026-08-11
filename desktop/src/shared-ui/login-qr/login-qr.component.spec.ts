import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideStore, Store } from "@ngxs/store";
import * as QRCode from "qrcode";

jest.mock("qrcode", () => ({
  __esModule: true,
  toDataURL: jest.fn(),
}));

import { FeatureConfigState } from "../../store/feature-config.state";
import { SetFeatureConfig } from "../../store/feature-config.state.actions";
import { LoginQrComponent } from "./login-qr.component";

describe("LoginQrComponent", () => {
  let component: LoginQrComponent;
  let fixture: ComponentFixture<LoginQrComponent>;
  let store: Store;

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

  beforeEach(async () => {
    (QRCode.toDataURL as jest.Mock).mockReset();

    await TestBed.configureTestingModule({
      imports: [LoginQrComponent],
      providers: [provideStore([FeatureConfigState])],
    }).compileComponents();

    fixture = TestBed.createComponent(LoginQrComponent);
    component = fixture.componentInstance;
    store = TestBed.inject(Store);
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("renders nothing when loginQrUrl is empty", () => {
    // Default featureConfig carries an empty loginQrUrl.
    expect(component.qrDataUrl()).toBeNull();
    expect(fixture.nativeElement.querySelector("img.login-qr-code")).toBeNull();
  });

  it("renders a locally-generated QR when loginQrUrl is set", async () => {
    (QRCode.toDataURL as jest.Mock).mockResolvedValue(
      "data:image/png;base64,QR"
    );
    const loginQrUrl = "https://receiptwrangler.io/app/setup#url=x";

    setLoginQrUrl(loginQrUrl);
    // Effect regenerates the data URL (async) from the new store value.
    await fixture.whenStable();
    fixture.detectChanges();

    expect(QRCode.toDataURL).toHaveBeenCalledWith(loginQrUrl, {
      margin: 2,
      width: 220,
    });
    expect(component.qrDataUrl()).toBe("data:image/png;base64,QR");
    expect(
      fixture.nativeElement.querySelector("img.login-qr-code").getAttribute("src")
    ).toBe("data:image/png;base64,QR");
  });

  it("renders the divider only when headerText is provided", async () => {
    (QRCode.toDataURL as jest.Mock).mockResolvedValue(
      "data:image/png;base64,QR"
    );

    setLoginQrUrl("https://receiptwrangler.io/app/setup#url=x");
    await fixture.whenStable();
    fixture.detectChanges();

    expect(
      fixture.nativeElement.querySelector(".login-qr-divider")
    ).toBeNull();

    fixture.componentRef.setInput("headerText", "Set up the mobile app");
    fixture.detectChanges();

    expect(
      fixture.nativeElement.querySelector(".login-qr-divider").textContent
    ).toContain("Set up the mobile app");
  });

  it("hides the QR again when the feature is turned off", async () => {
    (QRCode.toDataURL as jest.Mock).mockResolvedValue(
      "data:image/png;base64,QR"
    );

    setLoginQrUrl("https://receiptwrangler.io/app/setup#url=x");
    await fixture.whenStable();
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector("img.login-qr-code")).toBeTruthy();

    setLoginQrUrl("");
    await fixture.whenStable();
    fixture.detectChanges();

    expect(component.qrDataUrl()).toBeNull();
    expect(fixture.nativeElement.querySelector("img.login-qr-code")).toBeNull();
  });

  it("ignores a stale QR generation that resolves after a newer one", async () => {
    // Hand out manually-resolved promises so the two generations can be
    // completed out of order.
    const resolvers: Array<(dataUrl: string) => void> = [];
    (QRCode.toDataURL as jest.Mock).mockImplementation(
      () => new Promise<string>((resolve) => resolvers.push(resolve))
    );

    setLoginQrUrl("https://receiptwrangler.io/app/setup#url=stale");
    setLoginQrUrl("https://receiptwrangler.io/app/setup#url=current");
    expect(resolvers.length).toBe(2);

    // Reverse order: the current generation finishes first, then the stale one.
    resolvers[1]("data:image/png;base64,CURRENT");
    resolvers[0]("data:image/png;base64,STALE");
    await fixture.whenStable();

    expect(component.qrDataUrl()).toBe("data:image/png;base64,CURRENT");
  });
});
