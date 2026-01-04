import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormListComponent } from './form-list.component';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { MatListModule } from '@angular/material/list';

describe('FormListComponent', () => {
  let component: FormListComponent;
  let fixture: ComponentFixture<FormListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormListComponent],
      imports: [MatListModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    });
    fixture = TestBed.createComponent(FormListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
