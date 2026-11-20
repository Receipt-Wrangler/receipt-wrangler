import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { SystemSettingsState } from "../../store/system-settings.state";

import { AuditDetailSectionComponent } from "./audit-detail-section.component";

describe("AuditDetailSectionComponent", () => {
  let component: AuditDetailSectionComponent;
  let fixture: ComponentFixture<AuditDetailSectionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([SystemSettingsState]), AuditDetailSectionComponent]
    })
      .compileComponents();

    fixture = TestBed.createComponent(AuditDetailSectionComponent);
    component = fixture.componentInstance;
    // Every template that renders this section binds [data].
    fixture.componentRef.setInput("data", {});
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
