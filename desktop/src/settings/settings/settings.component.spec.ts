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

  const tabNames = () => component.tabs().map((t) => t.name);

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

  it('renders no tabs when the user can read none of them', () => {
    store.dispatch(new SetPermissions([], {}));
    fixture.detectChanges();

    expect(tabNames()).toEqual([]);
  });

  it('gates the User Profile tab on app.account.read', () => {
    store.dispatch(new SetPermissions([Permission.AppAccountRead], {}));
    fixture.detectChanges();

    expect(tabNames()).toEqual(['user-profile']);
  });

  it('gates the User Preferences tab on app.user-preferences.read', () => {
    store.dispatch(
      new SetPermissions([Permission.AppUserPreferencesRead], {})
    );
    fixture.detectChanges();

    expect(tabNames()).toEqual(['user-preferences']);
  });

  it('renders the tabs in order when all read permissions are present', () => {
    store.dispatch(
      new SetPermissions(
        [
          Permission.AppAccountRead,
          Permission.AppUserPreferencesRead,
          Permission.AppApiKeysRead,
        ],
        {}
      )
    );
    fixture.detectChanges();

    expect(tabNames()).toEqual([
      'user-profile',
      'user-preferences',
      'api-keys',
    ]);
  });

  it('hides the API Keys tab without app.api-keys.read', () => {
    store.dispatch(
      new SetPermissions(
        [Permission.AppAccountRead, Permission.AppUserPreferencesRead],
        {}
      )
    );
    fixture.detectChanges();

    expect(tabNames()).not.toContain('api-keys');
  });
});
