import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormGroup, ReactiveFormsModule } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { ButtonModule } from "../../button";
import { StoreModule } from "../../store/store.module";
import { SharedUiModule } from "../shared-ui.module";
import { FormComponent } from "./form.component";

describe("FormComponent", () => {
  let component: FormComponent;
  let fixture: ComponentFixture<FormComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormComponent],
      imports: [
        PipesModule,
        ReactiveFormsModule,
        ButtonModule,
        SharedUiModule,
        StoreModule,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {},
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });
    fixture = TestBed.createComponent(FormComponent);
    component = fixture.componentInstance;
    component.formConfig = {
      headerText: "test",
      mode: FormMode.add,
    };
    component.form = new FormGroup({});
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
