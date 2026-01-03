import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { SystemEmailTaskTableService } from "../../services/system-email-task-table.service";
import { SystemEmailTaskTableState } from "../../store/system-email-task-table.state";

import { TaskTableComponent } from "./task-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("TaskTableComponent", () => {
  let component: TaskTableComponent;
  let fixture: ComponentFixture<TaskTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [TaskTableComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [NgxsModule.forRoot([SystemEmailTaskTableState])],
    providers: [
        {
            provide: TABLE_SERVICE_INJECTION_TOKEN,
            useClass: SystemEmailTaskTableService
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
})
      .compileComponents();

    fixture = TestBed.createComponent(TaskTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
