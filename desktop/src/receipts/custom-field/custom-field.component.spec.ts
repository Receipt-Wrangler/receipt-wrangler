import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormGroup } from "@angular/forms";
import { CustomFieldPipe } from "../pipes/custom-field.pipe";

import { CustomFieldComponent } from "./custom-field.component";

describe("CustomFieldComponent", () => {
  let component: CustomFieldComponent;
  let fixture: ComponentFixture<CustomFieldComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CustomFieldComponent, CustomFieldPipe]
    })
      .compileComponents();


    fixture = TestBed.createComponent(CustomFieldComponent);
    component = fixture.componentInstance;
    component.formGroup = new FormGroup({}) as any;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
