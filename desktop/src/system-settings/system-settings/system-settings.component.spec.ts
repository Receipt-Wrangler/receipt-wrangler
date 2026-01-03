import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

import { SystemSettingsComponent } from "./system-settings.component";

describe("SystemSettingsComponent", () => {
  let component: SystemSettingsComponent;
  let fixture: ComponentFixture<SystemSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SystemSettingsComponent],
      imports: [SharedUiModule],
      providers: [

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

    fixture = TestBed.createComponent(SystemSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
