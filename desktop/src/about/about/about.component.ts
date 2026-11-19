import { CommonModule } from "@angular/common";
import { Component, inject } from "@angular/core";
import { Store } from "@ngxs/store";
import { About } from "../../open-api/index";
import { AppDatePipe } from "../../pipes/app-date.pipe";
import { AboutState } from "../../store/about.state";
import { FeatureConfigState } from "../../store/feature-config.state";
import { DialogComponent } from "../../shared-ui/dialog/dialog.component";
import { FormSectionComponent } from "../../shared-ui/form-section/form-section.component";
import { LoginQrComponent } from "../../shared-ui/login-qr/login-qr.component";

interface Link {
  url: string;
  label: string;
}

@Component({
    selector: "app-about",
    imports: [
        AppDatePipe,
        CommonModule,
        DialogComponent,
        FormSectionComponent,
        LoginQrComponent
    ],
    templateUrl: "./about.component.html",
    styleUrl: "./about.component.scss"
})
export class AboutComponent {
  private store = inject(Store);
  public about = this.store.selectSignal(AboutState.about);
  public loginQrUrl = this.store.selectSignal(FeatureConfigState.loginQrUrl);

  public links: Link[] = [
    {
      label: "Documentation",
      url: "https://receiptwrangler.io"
    },
    {
      label: "Source Code",
      url: "https://github.com/Receipt-Wrangler"
    },
    {
      label: "Reddit",
      url: "https://reddit.com/r/receiptwrangler"
    }
  ];

}
