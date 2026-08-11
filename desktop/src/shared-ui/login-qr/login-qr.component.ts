import { Component, effect, inject, input, signal } from "@angular/core";
import { Store } from "@ngxs/store";
import * as QRCode from "qrcode";
import { FeatureConfigState } from "../../store/feature-config.state";

/**
 * Renders the mobile app setup QR for `featureConfig.loginQrUrl`, or nothing at
 * all when the feature is disabled/unset. Used by the login page and the About
 * dialog so both share one generation path.
 */
@Component({
  selector: "app-login-qr",
  templateUrl: "./login-qr.component.html",
  styleUrls: ["./login-qr.component.scss"],
  standalone: true,
})
export class LoginQrComponent {
  /** Optional divider label rendered above the QR. Omit to render just the code. */
  public readonly headerText = input<string | undefined>(undefined);

  public readonly caption = input<string>(
    "Scan with your phone to open the app and connect it to this server."
  );

  // Rendered locally from featureConfig.loginQrUrl; null when the login QR is
  // disabled/unset, which hides the whole block.
  public qrDataUrl = signal<string | null>(null);

  private store = inject(Store);

  constructor() {
    const loginQrUrl = this.store.selectSignal(FeatureConfigState.loginQrUrl);
    effect((onCleanup) => {
      const url = loginQrUrl();
      if (!url) {
        this.qrDataUrl.set(null);
        return;
      }

      // Generation runs async, so a URL change can leave an earlier call in
      // flight. Cleanup runs before the next execution (and on destroy), so
      // only the latest generation is allowed to write the signal — otherwise
      // a late-resolving older QR could overwrite the current one.
      let cancelled = false;
      onCleanup(() => (cancelled = true));

      QRCode.toDataURL(url, { margin: 2, width: 220 })
        .then((dataUrl) => {
          if (!cancelled) {
            this.qrDataUrl.set(dataUrl);
          }
        })
        .catch(() => {
          if (!cancelled) {
            this.qrDataUrl.set(null);
          }
        });
    });
  }
}
