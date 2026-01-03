import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { FormConfig } from "../../interfaces";
import { PipesModule } from "../../pipes";

import { GroupSettingsAiComponent } from "./group-settings-ai.component";

describe("GroupSettingsAiComponent", () => {
  let component: GroupSettingsAiComponent;
  let fixture: ComponentFixture<GroupSettingsAiComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GroupSettingsAiComponent],
      imports: [PipesModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: {
                group: {}
              }
            }
          }
        }
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();

    fixture = TestBed.createComponent(GroupSettingsAiComponent);
    component = fixture.componentInstance;
    component.formConfig = {} as FormConfig;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
