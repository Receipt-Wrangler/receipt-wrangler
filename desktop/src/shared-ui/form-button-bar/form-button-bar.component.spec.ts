import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormMode } from 'src/enums/form-mode.enum';
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

  describe('inset', () => {
    function bar(): HTMLElement {
      return fixture.nativeElement.querySelector('[data-testid="form-button-bar"]');
    }

    beforeEach(() => {
      fixture.componentRef.setInput('mode', FormMode.edit);
    });

    it('spans the viewport, as the stylesheet sets it, without an inset', () => {
      fixture.detectChanges();

      expect(bar().style.left).toBe('');
      expect(bar().style.right).toBe('');
    });

    it('keeps clear of both edges by the inset', () => {
      fixture.componentRef.setInput('inset', { left: 82, right: 480 });
      fixture.detectChanges();

      expect(bar().style.left).toBe('82px');
      expect(bar().style.right).toBe('480px');
    });

    it('goes back to the full width when the inset is cleared', () => {
      fixture.componentRef.setInput('inset', { left: 82, right: 480 });
      fixture.detectChanges();
      fixture.componentRef.setInput('inset', null);
      fixture.detectChanges();

      expect(bar().style.left).toBe('');
      expect(bar().style.right).toBe('');
    });
  });
});
