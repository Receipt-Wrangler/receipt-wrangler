import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddReceiptIconComponent } from './add-receipt-icon.component';

describe('AddReceiptIconComponent', () => {
  let component: AddReceiptIconComponent;
  let fixture: ComponentFixture<AddReceiptIconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AddReceiptIconComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AddReceiptIconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
