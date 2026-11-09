import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatIconModule } from '@angular/material/icon';

import { StatusIconComponent } from './status-icon.component';

describe('StatusIconComponent', () => {
  let component: StatusIconComponent;
  let fixture: ComponentFixture<StatusIconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MatIconModule, StatusIconComponent],
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(StatusIconComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('taskStatus', 'SUCCEEDED');
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
