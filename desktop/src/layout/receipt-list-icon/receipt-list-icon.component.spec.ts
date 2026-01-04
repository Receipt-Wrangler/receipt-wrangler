import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReceiptListIconComponent } from './receipt-list-icon.component';

describe('ReceiptListIconComponent', () => {
  let component: ReceiptListIconComponent;
  let fixture: ComponentFixture<ReceiptListIconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ReceiptListIconComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ReceiptListIconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
