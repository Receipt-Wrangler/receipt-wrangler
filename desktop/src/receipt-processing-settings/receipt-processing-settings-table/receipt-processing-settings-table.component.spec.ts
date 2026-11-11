import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";

import { ReceiptProcessingSettingsTableComponent } from "./receipt-processing-settings-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptProcessingSettingsTableComponent", () => {
  let component: ReceiptProcessingSettingsTableComponent;
  let fixture: ComponentFixture<ReceiptProcessingSettingsTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    imports: [NgxsModule.forRoot(), ReceiptProcessingSettingsTableComponent],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {}
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
    ]
})
      .compileComponents();

    fixture = TestBed.createComponent(ReceiptProcessingSettingsTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
