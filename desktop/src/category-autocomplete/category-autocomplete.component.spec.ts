import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl } from "@angular/forms";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";

import { CategoryAutocompleteComponent } from "./category-autocomplete.component";

describe("CategoryAutocompleteComponent", () => {
  let component: CategoryAutocompleteComponent;
  let fixture: ComponentFixture<CategoryAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CategoryAutocompleteComponent, NoopAnimationsModule],
    })
      .compileComponents();

    fixture = TestBed.createComponent(CategoryAutocompleteComponent);
    component = fixture.componentInstance;
    component.inputFormControl = new FormControl();
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
