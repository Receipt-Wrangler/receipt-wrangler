import { provideZonelessChangeDetection } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AvatarComponent } from './avatar.component';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { User } from '../../open-api';

describe('AvatarComponent', () => {
  let component: AvatarComponent;
  let fixture: ComponentFixture<AvatarComponent>;

  const user = { id: 1, displayName: 'maya chen', defaultAvatarColor: '#a78bfa' } as User;

  function circle(): HTMLElement {
    return fixture.nativeElement.querySelector('.avatar');
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AvatarComponent],
      imports: [MatIconModule, MatTooltipModule],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(AvatarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  // Every existing caller passes only a user or a group; none of them may pick up
  // inline sizing, so the stylesheet alone keeps sizing them exactly as before.
  it('renders the default avatar sized by its stylesheet alone', async () => {
    fixture.componentRef.setInput('user', user);
    await fixture.whenStable();

    expect(circle().style.width).toBe('');
    expect(circle().style.height).toBe('');
    expect(circle().style.fontSize).toBe('');
    expect(circle().style.backgroundColor).toBe('rgb(167, 139, 250)');
    expect(circle().textContent?.trim()).toBe('M');
    expect(fixture.nativeElement.querySelector('mat-icon')).toBeNull();
  });

  it('falls back to the brand blue without a user colour', async () => {
    fixture.componentRef.setInput('group', { id: 2, name: 'home' });
    await fixture.whenStable();

    expect(circle().style.backgroundColor).toBe('rgb(39, 177, 255)');
    expect(circle().textContent?.trim()).toBe('H');
  });

  it('sizes the circle and its initial for a custom size', async () => {
    fixture.componentRef.setInput('user', user);
    fixture.componentRef.setInput('size', 32);
    await fixture.whenStable();

    expect(circle().style.width).toBe('32px');
    expect(circle().style.height).toBe('32px');
    expect(circle().style.fontSize).toBe('14px');
  });

  it('renders an icon instead of the initial, in the given colour', async () => {
    fixture.componentRef.setInput('icon', 'mail');
    fixture.componentRef.setInput('color', '#64748b');
    fixture.componentRef.setInput('size', 32);
    await fixture.whenStable();

    const icon: HTMLElement = fixture.nativeElement.querySelector('mat-icon');
    expect(icon.textContent?.trim()).toBe('mail');
    expect(icon.style.fontSize).toBe('18px');
    expect(circle().style.backgroundColor).toBe('rgb(100, 116, 139)');
  });

  it('lets an explicit colour win over the user colour', async () => {
    fixture.componentRef.setInput('user', user);
    fixture.componentRef.setInput('color', '#94a3b8');
    await fixture.whenStable();

    expect(circle().style.backgroundColor).toBe('rgb(148, 163, 184)');
    expect(circle().textContent?.trim()).toBe('M');
  });
});
