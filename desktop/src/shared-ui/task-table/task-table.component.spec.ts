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

  it("omits the source file column by default, so other hosts are unchanged", () => {
    fixture.componentRef.setInput("expandedRowTemplate", {} as any);
    // setColumns() runs in ngAfterViewInit and mutates template-bound state, so
    // the dev-mode check-no-changes pass would report NG0100.
    fixture.detectChanges(false);

    expect(component.displayedColumns).not.toContain("source_file");
    expect(component.columns.some((column) => column.matColumnDef === "source_file")).toBe(false);
  });

  it("adds the source file column before expand when the host opts in", () => {
    fixture.componentRef.setInput("showSourceFileActions", true);
    fixture.componentRef.setInput("expandedRowTemplate", {} as any);
    fixture.detectChanges(false);

    // mat-table throws on a displayed id with no matching column definition, and
    // "expand" has to stay last.
    expect(component.columns.some((column) => column.matColumnDef === "source_file")).toBe(true);
    expect(component.displayedColumns.indexOf("source_file")).toBeGreaterThan(-1);
    expect(component.displayedColumns.indexOf("source_file"))
      .toBeLessThan(component.displayedColumns.indexOf("expand"));
    expect(component.displayedColumns[component.displayedColumns.length - 1]).toBe("expand");
  });
});
