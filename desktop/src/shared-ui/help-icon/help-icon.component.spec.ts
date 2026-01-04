import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HelpIconComponent } from './help-icon.component';

describe('HelpIconComponent', () => {
  let component: HelpIconComponent;
  let fixture: ComponentFixture<HelpIconComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [HelpIconComponent]
    });
    fixture = TestBed.createComponent(HelpIconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
