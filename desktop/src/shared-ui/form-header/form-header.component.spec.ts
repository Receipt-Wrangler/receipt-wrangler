import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormHeaderComponent } from './form-header.component';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';

describe('FormHeaderComponent', () => {
  let component: FormHeaderComponent;
  let fixture: ComponentFixture<FormHeaderComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormHeaderComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    });
    fixture = TestBed.createComponent(FormHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
