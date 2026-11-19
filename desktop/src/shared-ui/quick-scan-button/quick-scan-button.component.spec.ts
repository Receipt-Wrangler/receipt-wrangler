import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { provideTestStore } from "../../store/test-store";

import { QuickScanButtonComponent } from "./quick-scan-button.component";

describe("QuickScanButtonComponent", () => {
  let component: QuickScanButtonComponent;
  let fixture: ComponentFixture<QuickScanButtonComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MatDialogModule, QuickScanButtonComponent],
      providers: [
        provideTestStore(),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });
    fixture = TestBed.createComponent(QuickScanButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
