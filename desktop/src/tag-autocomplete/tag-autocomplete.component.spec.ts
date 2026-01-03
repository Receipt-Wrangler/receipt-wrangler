import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl } from "@angular/forms";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";

import { TagAutocompleteComponent } from "./tag-autocomplete.component";

describe("TagAutocompleteComponent", () => {
  let component: TagAutocompleteComponent;
  let fixture: ComponentFixture<TagAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TagAutocompleteComponent, NoopAnimationsModule]
    })
      .compileComponents();

    fixture = TestBed.createComponent(TagAutocompleteComponent);
    component = fixture.componentInstance;
    component.inputFormControl = new FormControl();
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
