import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";

import { ExportButtonComponent } from "./export-button.component";

describe("ExportButtonComponent", () => {
  let component: ExportButtonComponent;
  let fixture: ComponentFixture<ExportButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ExportButtonComponent],
      providers: [
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
        {
          provide: ActivatedRoute,
          useValue: {}
        }
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(ExportButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
