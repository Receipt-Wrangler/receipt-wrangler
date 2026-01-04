import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardIconComponent } from './dashboard-icon.component';

describe('DashboardIconComponent', () => {
  let component: DashboardIconComponent;
  let fixture: ComponentFixture<DashboardIconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DashboardIconComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DashboardIconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
