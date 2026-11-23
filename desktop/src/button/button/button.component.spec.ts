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

  // The host has no role, so popup state only counts on the native button.
  describe.each(['basic', 'matRaisedButton', 'iconButton'] as const)(
    'popup state on a %s button',
    (matButtonType) => {
      const render = (inputs: { ariaHasPopup?: string; ariaExpanded?: boolean } = {}): HTMLButtonElement => {
        fixture.componentRef.setInput('matButtonType', matButtonType);
        for (const [name, value] of Object.entries(inputs)) {
          fixture.componentRef.setInput(name, value);
        }
        fixture.detectChanges();
        return fixture.nativeElement.querySelector('button');
      };

      it('forwards aria-haspopup and aria-expanded to the native button', () => {
        const button = render({ ariaHasPopup: 'dialog', ariaExpanded: true });

        expect(button.getAttribute('aria-haspopup')).toBe('dialog');
        expect(button.getAttribute('aria-expanded')).toBe('true');
      });

      it('renders a closed popup as aria-expanded="false", not as absent', () => {
        expect(render({ ariaExpanded: false }).getAttribute('aria-expanded')).toBe('false');
      });

      it('carries neither attribute by default', () => {
        const button = render();

        expect(button.hasAttribute('aria-haspopup')).toBe(false);
        expect(button.hasAttribute('aria-expanded')).toBe(false);
      });
    },
  );
});
