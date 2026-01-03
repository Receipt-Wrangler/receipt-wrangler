import { AfterViewInit, Component, Input, OnInit, TemplateRef, ViewChild } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Prompt, SystemTask, SystemTaskStatus, SystemTaskType } from "../../open-api";
import { AccordionPanel } from "../../shared-ui/accordion/accordion-panel.interface";

@Component({
    selector: "app-receipt-processing-settings-child-system-task-accordion",
    templateUrl: "./receipt-processing-settings-child-system-task-accordion.component.html",
    styleUrl: "./receipt-processing-settings-child-system-task-accordion.component.scss",
    standalone: false
})
export class ReceiptProcessingSettingsChildSystemTaskAccordionComponent implements OnInit, AfterViewInit {

  @ViewChild("ocrProcessingDetails") public ocrProcessingDetails!: TemplateRef<any>;

  @ViewChild("promptGenerationDetails") public promptGenerationDetails!: TemplateRef<any>;

  @ViewChild("chatCompletionDetails") public chatCompletionDetails!: TemplateRef<any>;

  @ViewChild("receiptUploadedDetails") public receiptUploadedDetails!: TemplateRef<any>;

  @ViewChild("statusIcon") public statusIcon!: TemplateRef<any>;

  @Input() public childTasks: SystemTask[] = [];
  protected readonly SystemTaskType = SystemTaskType;

  protected readonly SystemTaskStatus = SystemTaskStatus;

  public accordionPanels: AccordionPanel[] = [];

  public prompts: Prompt[] = [];

  constructor(private activatedRoute: ActivatedRoute) {}

  public ngOnInit(): void {
    this.prompts = this.activatedRoute.snapshot.data["prompts"];
  }

  public ngAfterViewInit(): void {
    this.setAccordionPanels();
  }

  private setAccordionPanels(): void {
    this.childTasks.forEach(task => {
      if (task.type === SystemTaskType.OcrProcessing) {
        this.accordionPanels.push({
          title: "Raw OCR Processing Details",
          content: this.ocrProcessingDetails,
          descriptionTemplate: this.statusIcon,
        });
      }

      if (task.type === SystemTaskType.PromptGenerated) {
        const prompt = this.prompts.find(p => p.id === task.associatedEntityId);
        const title = `Prompt Used: ${prompt?.name}`;

        this.accordionPanels.push({
          title: title,
          content: this.promptGenerationDetails,
          descriptionTemplate: this.statusIcon,
        });
      }

      if (task.type === SystemTaskType.ChatCompletion) {
        this.accordionPanels.push({
          title: "Raw Chat Completion Details",
          content: this.chatCompletionDetails,
          descriptionTemplate: this.statusIcon,
        });
      }

      if (task.type === SystemTaskType.ReceiptUploaded) {
        this.accordionPanels.push({
          title: "Receipt Uploaded",
          content: this.receiptUploadedDetails,
          descriptionTemplate: this.statusIcon,
        });
      }
    });
  }
}
