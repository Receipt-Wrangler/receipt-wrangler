import { Component, provideZonelessChangeDetection } from '@angular/core';
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

// Host component exercising the [table-header-subtitle] content-projection slot,
// used for rich subtitles that need markup (e.g. an inline link) rather than the
// plain-string `subtitle` input.
@Component({
  template: `
    <app-table-header headerText="Hosted">
      <p table-header-subtitle>
        Projected <a href="#">Manage Users</a>
      </p>
    </app-table-header>
  `,
  standalone: false,
})
class TableHeaderHostComponent {}

describe('TableHeaderComponent — projected subtitle', () => {
  let fixture: ComponentFixture<TableHeaderHostComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TableHeaderComponent, TableHeaderHostComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(TableHeaderHostComponent);
    await fixture.whenStable();
  });

  it('renders content projected into the [table-header-subtitle] slot', async () => {
    await fixture.whenStable();

    const projected = fixture.nativeElement.querySelector('[table-header-subtitle]');
    expect(projected).toBeTruthy();
    expect(projected.textContent).toContain('Projected');
    expect(projected.querySelector('a')?.textContent).toContain('Manage Users');
    // The plain-string subtitle path is not used here, so its element is absent.
    expect(fixture.nativeElement.querySelector('.table-header__subtitle')).toBeFalsy();
  });
});
