import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { ButtonComponent } from "../../button/index";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SystemTaskTableState } from "../../store/system-task-table.state";
import { TableModule } from "../../table/table.module";

import { SystemTaskTableComponent } from "./system-task-table.component";

describe("SystemTaskTableComponent", () => {
  let component: SystemTaskTableComponent;
  let fixture: ComponentFixture<SystemTaskTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SystemTaskTableComponent, ButtonComponent],
      imports: [SharedUiModule,
        NgxsModule.forRoot([SystemTaskTableState]),
        TableModule,
        NoopAnimationsModule],
      providers: [{
        provide: ActivatedRoute,
        useValue: {
          snapshot: {
            data: {
              prompts: [],
              allReceiptProcessingSettings: []
            }
          }
        }
      }, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
    })
      .compileComponents();

    fixture = TestBed.createComponent(SystemTaskTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
