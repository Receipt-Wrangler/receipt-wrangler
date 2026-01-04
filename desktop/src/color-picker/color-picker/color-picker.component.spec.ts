import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ColorPickerComponent } from './color-picker.component';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('ColorPickerComponent', () => {
  let component: ColorPickerComponent;
  let fixture: ComponentFixture<ColorPickerComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ColorPickerComponent],
      imports: [
        ReactiveFormsModule,
        MatFormFieldModule,
        MatInputModule,
        NoopAnimationsModule,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    });
    fixture = TestBed.createComponent(ColorPickerComponent);
    component = fixture.componentInstance;
    component.inputFormControl = new FormControl();
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set form control value when a color is selected', () => {
    const event = {
      target: {
        value: '#CD5C5C',
      },
    };
    component.colorSelected(event);

    expect(component.inputFormControl.value).toEqual('#CD5C5C');
  });

  it('should not prevent default if the field is not readonly', () => {
    const event: any = {
      preventDefault: () => {},
    };
    const eventSpy = jest.spyOn(event, 'preventDefault');

    component.handleClick(event);

    expect(eventSpy).toHaveBeenCalledTimes(0);
  });

  it('should  prevent default if the field is readonly', () => {
    const event: any = {
      preventDefault: () => {},
    };
    const eventSpy = jest.spyOn(event, 'preventDefault');
    component.readonly = true;
    component.handleClick(event);

    expect(eventSpy).toHaveBeenCalledTimes(1);
  });
});
