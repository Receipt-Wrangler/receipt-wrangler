import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormButtonBarComponent } from './form-button-bar.component';

describe('FormButtonBarComponent', () => {
  let component: FormButtonBarComponent;
  let fixture: ComponentFixture<FormButtonBarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ FormButtonBarComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FormButtonBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
