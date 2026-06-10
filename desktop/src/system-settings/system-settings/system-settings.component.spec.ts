import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Permission } from "../../open-api";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { AuthState } from "../../store/auth.state";
import { SetPermissions } from "../../store/auth.state.actions";

import { SystemSettingsComponent } from "./system-settings.component";

describe("SystemSettingsComponent", () => {
  let component: SystemSettingsComponent;
  let fixture: ComponentFixture<SystemSettingsComponent>;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SystemSettingsComponent],
      imports: [SharedUiModule, NgxsModule.forRoot([AuthState])],
      providers: [
        provideZonelessChangeDetection(),
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              queryParams: {
                tab: "settings",
              },
            }
          }
        }
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    })
      .compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(SystemSettingsComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    store.dispatch(
      new SetPermissions([Permission.AppSystemEmailsRead], {})
    );
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it("only shows tabs the user can read", () => {
    store.dispatch(
      new SetPermissions(
        [Permission.AppPromptsRead, Permission.AppSystemTasksRead],
        {}
      )
    );

    fixture.detectChanges();

    expect(component.tabs.map((tab) => tab.name)).toEqual([
      "prompts",
      "system-tasks",
    ]);
  });

  it("shows all tabs when the user can read everything", () => {
    store.dispatch(
      new SetPermissions(
        [
          Permission.AppSystemSettingsRead,
          Permission.AppReceiptProcessingSettingsRead,
          Permission.AppPromptsRead,
          Permission.AppSystemEmailsRead,
          Permission.AppSystemTasksRead,
        ],
        {}
      )
    );

    fixture.detectChanges();

    expect(component.tabs).toHaveLength(5);
  });
});
