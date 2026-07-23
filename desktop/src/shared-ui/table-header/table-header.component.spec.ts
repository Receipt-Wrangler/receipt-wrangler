import { provideZonelessChangeDetection } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TableHeaderComponent } from './table-header.component';

describe('TableHeaderComponent', () => {
  let component: TableHeaderComponent;
  let fixture: ComponentFixture<TableHeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TableHeaderComponent ],
      providers: [ provideZonelessChangeDetection() ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TableHeaderComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('renders no subtitle by default', async () => {
    await fixture.whenStable();

    expect(fixture.nativeElement.querySelector('.table-header__subtitle')).toBeFalsy();
  });

  it('renders the subtitle when provided', async () => {
    fixture.componentRef.setInput('subtitle', 'A helpful description');
    await fixture.whenStable();

    const subtitle = fixture.nativeElement.querySelector('.table-header__subtitle');
    expect(subtitle).toBeTruthy();
    expect(subtitle.textContent).toContain('A helpful description');
  });
});
