import { Component, OnInit, TemplateRef, computed, input, viewChild } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Prompt, SystemTask, SystemTaskStatus, SystemTaskType } from "../../open-api";
import { AccordionPanel } from "../../shared-ui/accordion/accordion-panel.interface";

@Component({
    selector: "app-receipt-processing-settings-child-system-task-accordion",
    templateUrl: "./receipt-processing-settings-child-system-task-accordion.component.html",
    styleUrl: "./receipt-processing-settings-child-system-task-accordion.component.scss",
    standalone: false
})
export class ReceiptProcessingSettingsChildSystemTaskAccordionComponent implements OnInit {

  public readonly ocrProcessingDetails = viewChild.required<TemplateRef<any>>("ocrProcessingDetails");

  public readonly promptGenerationDetails = viewChild.required<TemplateRef<any>>("promptGenerationDetails");

  public readonly chatCompletionDetails = viewChild.required<TemplateRef<any>>("chatCompletionDetails");

  public readonly receiptUploadedDetails = viewChild.required<TemplateRef<any>>("receiptUploadedDetails");

  public readonly statusIcon = viewChild.required<TemplateRef<any>>("statusIcon");

  public readonly childTasks = input<SystemTask[]>([]);
  protected readonly SystemTaskType = SystemTaskType;

  protected readonly SystemTaskStatus = SystemTaskStatus;

  /**
   * The children that get a panel, in panel order. The templates read the task for panel `i`
   * from HERE, not from `childTasks()`: the accordion passes the PANEL index, so a skipped task
   * type ahead of a panel's own task would hand the panel the wrong task.
   */
  public readonly panelTasks = computed(() =>
    this.childTasks().filter(task => this.panelFor(task) !== undefined)
  );

  /**
   * Derived rather than pushed in ngAfterViewInit: an array filled after the accordion had
   * already been checked only rendered on a later change-detection pass, which @for reports
   * as NG0100 in dev mode. `prompts` is a plain field, but ngOnInit sets it before first render.
   */
  public readonly accordionPanels = computed<AccordionPanel[]>(() =>
    this.panelTasks().map(task => ({
      ...this.panelFor(task)!,
      descriptionTemplate: this.statusIcon(),
    }))
  );

  public prompts: Prompt[] = [];

  constructor(private activatedRoute: ActivatedRoute) {}

  public ngOnInit(): void {
    this.prompts = this.activatedRoute.snapshot.data["prompts"];
  }

  /**
   * How a task's panel reads, or undefined for a task type that gets none. The one place that
   * decides which types get a panel, so `panelTasks` and `accordionPanels` cannot disagree.
   */
  private panelFor(task: SystemTask): Omit<AccordionPanel, "descriptionTemplate"> | undefined {
    switch (task.type) {
      case SystemTaskType.OcrProcessing:
        return { title: "Raw OCR Processing Details", content: this.ocrProcessingDetails() };
      case SystemTaskType.PromptGenerated:
        return {
          title: `Prompt Used: ${this.prompts.find(p => p.id === task.associatedEntityId)?.name}`,
          content: this.promptGenerationDetails(),
        };
      case SystemTaskType.ChatCompletion:
        return { title: "Raw Chat Completion Details", content: this.chatCompletionDetails() };
      case SystemTaskType.ReceiptUploaded:
        return { title: "Receipt Uploaded", content: this.receiptUploadedDetails() };
      default:
        return undefined;
    }
  }
}
