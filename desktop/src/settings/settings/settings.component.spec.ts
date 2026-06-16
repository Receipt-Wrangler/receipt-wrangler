import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SettingsComponent } from './settings.component';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { NgxsModule, Store } from '@ngxs/store';
import { Permission } from '../../open-api';
import { AuthState } from '../../store';
import { SetPermissions } from '../../store/auth.state.actions';

describe('SettingsComponent', () => {
  let component: SettingsComponent;
  let fixture: ComponentFixture<SettingsComponent>;
  let store: Store;

  const tabNames = () => component.tabs.map((t) => t.name);

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SettingsComponent],
      imports: [
        MatTabsModule,
        MatTooltipModule,
        NgxsModule.forRoot([AuthState]),
        ReactiveFormsModule,
        RouterModule,
      ],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {},
        },
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(SettingsComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('hides the API Keys tab without app.api-keys.read', () => {
    store.dispatch(new SetPermissions([], {}));
    fixture.detectChanges();

    expect(tabNames()).not.toContain('api-keys');
    expect(tabNames()).toEqual(['user-profile', 'user-preferences']);
  });

  it('shows the API Keys tab with app.api-keys.read', () => {
    store.dispatch(new SetPermissions([Permission.AppApiKeysRead], {}));
    fixture.detectChanges();

    expect(tabNames()).toContain('api-keys');
  });
});
