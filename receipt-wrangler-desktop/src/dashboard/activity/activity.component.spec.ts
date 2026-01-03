import { ScrollingModule } from "@angular/cdk/scrolling";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatListModule } from "@angular/material/list";
import { NgxsModule } from "@ngxs/store";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { DashboardListComponent } from "../dashboard-list/dashboard-list.component";

import { ActivityComponent } from "./activity.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ActivityComponent", () => {
  let component: ActivityComponent;
  let fixture: ComponentFixture<ActivityComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [ActivityComponent, DashboardListComponent],
    imports: [SharedUiModule, ScrollingModule, MatListModule, NgxsModule.forRoot([])],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
})
      .compileComponents();

    fixture = TestBed.createComponent(ActivityComponent);
    component = fixture.componentInstance;
    component.widget = { name: "" } as any;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
