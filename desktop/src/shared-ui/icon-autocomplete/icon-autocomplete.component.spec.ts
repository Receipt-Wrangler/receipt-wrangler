import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl } from "@angular/forms";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { NgxsModule } from "@ngxs/store";
import { AuthState } from "../../store/index";

import { IconAutocompleteComponent } from "./icon-autocomplete.component";

describe("IconAutocompleteComponent", () => {
  let component: IconAutocompleteComponent;
  let fixture: ComponentFixture<IconAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState]), NoopAnimationsModule, IconAutocompleteComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    })
      .compileComponents();

    fixture = TestBed.createComponent(IconAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('inputFormControl', new FormControl(""));
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
