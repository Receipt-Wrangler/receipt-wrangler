import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ButtonComponent } from './button.component';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterTestingModule } from '@angular/router/testing';
import { MatBadgeModule } from '@angular/material/badge';

describe('ButtonComponent', () => {
  let component: ButtonComponent;
  let fixture: ComponentFixture<ButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        MatBadgeModule,
        MatButtonModule,
        MatTooltipModule,
        RouterTestingModule,
        ButtonComponent,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('icon button accessible name', () => {
    const renderIconButton = (ariaLabel?: string): HTMLButtonElement => {
      fixture.componentRef.setInput('matButtonType', 'iconButton');
      if (ariaLabel !== undefined) {
        fixture.componentRef.setInput('ariaLabel', ariaLabel);
      }
      fixture.detectChanges();
      return fixture.nativeElement.querySelector('button');
    };

    it('takes its aria-label from ariaLabel', () => {
      expect(renderIconButton('Close').getAttribute('aria-label')).toBe('Close');
    });

    it('carries no aria-label by default', () => {
      expect(renderIconButton().hasAttribute('aria-label')).toBe(false);
    });
  });
});
