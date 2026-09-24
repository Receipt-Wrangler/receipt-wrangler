import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl } from "@angular/forms";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { NgxsModule } from "@ngxs/store";
import { AuthState } from "../store/auth.state";

import { CategoryAutocompleteComponent } from "./category-autocomplete.component";

describe("CategoryAutocompleteComponent", () => {
  let component: CategoryAutocompleteComponent;
  let fixture: ComponentFixture<CategoryAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CategoryAutocompleteComponent, NgxsModule.forRoot([AuthState]), NoopAnimationsModule],
    })
      .compileComponents();

    fixture = TestBed.createComponent(CategoryAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('inputFormControl', new FormControl());
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
