import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl } from "@angular/forms";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { NgxsModule } from "@ngxs/store";
import { AutocompleteModule } from "../../autocomplete/autocomplete.module";
import { AuthState } from "../../store/index";

import { IconAutocompleteComponent } from "./icon-autocomplete.component";

describe("IconAutocompleteComponent", () => {
  let component: IconAutocompleteComponent;
  let fixture: ComponentFixture<IconAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [IconAutocompleteComponent],
      imports: [AutocompleteModule, NgxsModule.forRoot([AuthState]), NoopAnimationsModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    })
      .compileComponents();

    fixture = TestBed.createComponent(IconAutocompleteComponent);
    component = fixture.componentInstance;
    component.inputFormControl = new FormControl("");
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
