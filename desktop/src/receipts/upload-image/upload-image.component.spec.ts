import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { ApiModule } from "../../open-api";
import { UploadImageComponent } from "./upload-image.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("UploadImageComponent", () => {
  let component: UploadImageComponent;
  let fixture: ComponentFixture<UploadImageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [UploadImageComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule, MatSnackBarModule],
    providers: [
        { provide: ActivatedRoute, useValue: { snapshot: { data: {} } } },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    fixture = TestBed.createComponent(UploadImageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
