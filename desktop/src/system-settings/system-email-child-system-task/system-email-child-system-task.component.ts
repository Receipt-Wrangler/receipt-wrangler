import {AfterViewInit, Component, Input, TemplateRef, ViewChild} from '@angular/core';
import {ReceiptProcessingSettings, SystemTask, SystemTaskType} from "../../open-api";
import {AccordionPanel} from "../../shared-ui/accordion/accordion-panel.interface";

@Component({
    selector: 'app-system-email-child-system-task',
    templateUrl: './system-email-child-system-task.component.html',
    styleUrl: './system-email-child-system-task.component.scss',
    standalone: false
})
export class SystemEmailChildSystemTaskComponent implements AfterViewInit {
  @ViewChild('emailUploadDetails') public emailUploadDetails!: TemplateRef<any>;

  @Input() public childTasks: SystemTask[] = [];

  @Input() public allReceiptProcessingSettings: ReceiptProcessingSettings[] = [];

  public accordionPanels: AccordionPanel[] = [];

  public ngAfterViewInit(): void {
    this.initAccordionPanels();
  }

  private initAccordionPanels(): void {
    this.childTasks.forEach(task => {
      if (task.type === SystemTaskType.EmailUpload) {
        const settings = this.allReceiptProcessingSettings.find(s => s.id === task.associatedEntityId);
        let description = "";

        if (settings) {
          description = `Used ${settings?.name} to process Receipt`;
        }


        this.accordionPanels.push({
          title: 'Email Upload Details',
          description: description,
          content: this.emailUploadDetails,
        });
      }
    })
  }
}
