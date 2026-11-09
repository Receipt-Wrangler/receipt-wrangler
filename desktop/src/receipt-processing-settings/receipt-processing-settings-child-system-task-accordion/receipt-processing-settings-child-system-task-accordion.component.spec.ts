import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { AccordionModule } from "ngx-bootstrap/accordion";
import { SystemTask, SystemTaskStatus, SystemTaskType } from "../../open-api";
import { AppDatePipe } from "../../pipes/app-date.pipe";
import { SystemSettingsState } from "../../store/system-settings.state";

import {
  ReceiptProcessingSettingsChildSystemTaskAccordionComponent
} from "./receipt-processing-settings-child-system-task-accordion.component";

function task(id: number, type: SystemTaskType, extra: Partial<SystemTask> = {}): SystemTask {
  return { id, type, ...extra } as SystemTask;
}

describe("ReceiptProcessingSettingsChildSystemTaskAccordionComponent", () => {
  let component: ReceiptProcessingSettingsChildSystemTaskAccordionComponent;
  let fixture: ComponentFixture<ReceiptProcessingSettingsChildSystemTaskAccordionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReceiptProcessingSettingsChildSystemTaskAccordionComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: {
                prompts: [{ id: 7, name: "Receipt prompt" }],
              }
            }
          }
        }
      ],
      // The pipe and store let the panel templates themselves render (their dates use appDate).
      imports: [AccordionModule, AppDatePipe, NgxsModule.forRoot([SystemSettingsState])],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();

    fixture = TestBed.createComponent(ReceiptProcessingSettingsChildSystemTaskAccordionComponent);
    component = fixture.componentInstance;
  });

  // Inputs are set BEFORE the first detectChanges, so the template's accordionPanels() read is
  // also the proof that the viewChild.required templates resolve on first render.
  function render(childTasks: SystemTask[]) {
    fixture.componentRef.setInput("childTasks", childTasks);
    fixture.detectChanges();
  }

  it("should create", () => {
    render([]);
    expect(component).toBeTruthy();
    expect(component.accordionPanels()).toEqual([]);
  });

  it("builds one panel per processing step, in task order, and skips other task types", () => {
    render([
      task(1, SystemTaskType.OcrProcessing),
      task(2, SystemTaskType.PromptGenerated, { associatedEntityId: 7 }),
      task(3, SystemTaskType.HtmlToPdf),
      task(4, SystemTaskType.ChatCompletion),
      task(5, SystemTaskType.ReceiptUploaded),
    ]);

    const statusIcon = component.statusIcon();
    expect(component.accordionPanels()).toEqual([
      { title: "Raw OCR Processing Details", content: component.ocrProcessingDetails(), descriptionTemplate: statusIcon },
      { title: "Prompt Used: Receipt prompt", content: component.promptGenerationDetails(), descriptionTemplate: statusIcon },
      { title: "Raw Chat Completion Details", content: component.chatCompletionDetails(), descriptionTemplate: statusIcon },
      { title: "Receipt Uploaded", content: component.receiptUploadedDetails(), descriptionTemplate: statusIcon },
    ]);
  });

  it("re-derives the panels when the child tasks change", () => {
    render([task(1, SystemTaskType.OcrProcessing)]);
    expect(component.accordionPanels().map(panel => panel.title)).toEqual(["Raw OCR Processing Details"]);

    fixture.componentRef.setInput("childTasks", [task(2, SystemTaskType.ChatCompletion)]);
    fixture.detectChanges();

    expect(component.accordionPanels().map(panel => panel.title)).toEqual(["Raw Chat Completion Details"]);
  });

  it("renders each panel from its own task when a skipped task type comes first", () => {
    render([
      task(1, SystemTaskType.OcrProcessing),
      task(2, SystemTaskType.HtmlToPdf, {
        resultDescription: "pdf-description",
        status: SystemTaskStatus.Succeeded,
        startedAt: "2020-01-02T10:00:00Z",
        endedAt: "2020-01-02T10:00:01Z",
      }),
      task(3, SystemTaskType.ChatCompletion, {
        resultDescription: "chat-json",
        status: SystemTaskStatus.Failed,
        startedAt: "2031-03-04T10:00:00Z",
        endedAt: "2031-03-04T10:00:01Z",
      }),
    ]);

    // Panel 1 is the chat completion. Render its content and status the way the accordion
    // does: the panel's templates, given the PANEL index.
    const panel = component.accordionPanels()[1];
    expect(panel.title).toBe("Raw Chat Completion Details");

    const content = panel.content.createEmbeddedView({ index: 1 });
    content.detectChanges();
    const contentNodes = content.rootNodes as Node[];
    const prettyJson = contentNodes.find(node => node.nodeName === "APP-PRETTY-JSON") as any;
    const contentText = contentNodes.map(node => node.textContent ?? "").join(" ");

    const status = panel.descriptionTemplate!.createEmbeddedView({ index: 1 });
    status.detectChanges();
    const statusIcon = (status.rootNodes as Node[]).find(node => node.nodeName === "APP-STATUS-ICON") as any;

    expect(prettyJson.json).toBe("chat-json");
    expect(contentText).toContain("2031");
    expect(contentText).not.toContain("2020");
    expect(statusIcon.taskStatus).toBe(SystemTaskStatus.Failed);

    content.destroy();
    status.destroy();
  });
});
