import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { AccordionModule } from "ngx-bootstrap/accordion";
import { SystemTask, SystemTaskType } from "../../open-api";

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
      imports: [AccordionModule],
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
});
