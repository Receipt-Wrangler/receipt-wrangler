import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatusSelectComponent } from './status-select.component';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { SelectModule } from 'src/select/select.module';
import { FormControl } from '@angular/forms';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { RECEIPT_STATUS_OPTIONS } from 'src/constants';

describe('StatusSelectComponent', () => {
  let component: StatusSelectComponent;
  let fixture: ComponentFixture<StatusSelectComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [SelectModule, NoopAnimationsModule, StatusSelectComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    });
    fixture = TestBed.createComponent(StatusSelectComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('inputFormControl', new FormControl());
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should add blank option', () => {
    fixture.componentRef.setInput('addBlankOption', true);
    fixture.detectChanges();
    expect(component.receiptStatusOptions[0]).toEqual({
      value: null,
      displayValue: '',
    });
    expect(component.receiptStatusOptions.length).toBeGreaterThan(RECEIPT_STATUS_OPTIONS.length);
  });
});
