import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NgxsModule } from '@ngxs/store';
import { SystemSettingsState } from '../../../store/system-settings.state';
import { DateBlockComponent } from './date-block.component';

describe('DateBlockComponent', () => {
  let component: DateBlockComponent;
  let fixture: ComponentFixture<DateBlockComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DateBlockComponent, NgxsModule.forRoot([SystemSettingsState])]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DateBlockComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('date', new Date());
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
