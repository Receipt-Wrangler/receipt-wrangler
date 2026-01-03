import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DateBlockComponent } from './date-block.component';

describe('DateBlockComponent', () => {
  let component: DateBlockComponent;
  let fixture: ComponentFixture<DateBlockComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DateBlockComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DateBlockComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
