import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { DirectivesModule } from "../../directives";
import { PipesModule } from "../../pipes";
import { StoreModule } from "../../store/store.module";

import { QuickScanButtonComponent } from "./quick-scan-button.component";

describe("QuickScanButtonComponent", () => {
  let component: QuickScanButtonComponent;
  let fixture: ComponentFixture<QuickScanButtonComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [QuickScanButtonComponent],
      imports: [MatDialogModule, PipesModule, DirectivesModule, StoreModule],
      providers: [
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
