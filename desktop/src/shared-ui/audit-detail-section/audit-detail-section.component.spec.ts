import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { SystemSettingsState } from "../../store/system-settings.state";
import { SharedUiModule } from "../shared-ui.module";

import { AuditDetailSectionComponent } from "./audit-detail-section.component";

describe("AuditDetailSectionComponent", () => {
  let component: AuditDetailSectionComponent;
  let fixture: ComponentFixture<AuditDetailSectionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AuditDetailSectionComponent],
      imports: [SharedUiModule, NgxsModule.forRoot([SystemSettingsState])]
    })
      .compileComponents();

    fixture = TestBed.createComponent(AuditDetailSectionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
