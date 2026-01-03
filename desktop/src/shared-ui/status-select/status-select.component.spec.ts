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
      declarations: [StatusSelectComponent],
      imports: [SelectModule, NoopAnimationsModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    });
    fixture = TestBed.createComponent(StatusSelectComponent);
    component = fixture.componentInstance;
    component.inputFormControl = new FormControl();
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should add blank option', () => {
    component.addBlankOption = true;
    component.ngOnChanges({ addBlankOption: { currentValue: true } } as any);
    expect(component.receiptStatusOptions).toEqual([
      {
        value: null,
        displayValue: '',
      } as any,
      ...RECEIPT_STATUS_OPTIONS,
    ]);
  });
});
