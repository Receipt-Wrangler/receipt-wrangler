import {Component, TemplateRef, computed, input, viewChild} from '@angular/core';
import {ReceiptProcessingSettings, SystemTask, SystemTaskType} from "../../open-api";
import {AccordionPanel} from "../../shared-ui/accordion/accordion-panel.interface";

@Component({
    selector: 'app-system-email-child-system-task',
    templateUrl: './system-email-child-system-task.component.html',
    styleUrl: './system-email-child-system-task.component.scss',
    standalone: false
})
export class SystemEmailChildSystemTaskComponent {
  public readonly emailUploadDetails = viewChild.required<TemplateRef<any>>('emailUploadDetails');

  public readonly childTasks = input<SystemTask[]>([]);

  public readonly allReceiptProcessingSettings = input<ReceiptProcessingSettings[]>([]);

  /**
   * Derived rather than pushed in ngAfterViewInit: an array filled after the accordion had
   * already been checked only rendered on a later change-detection pass, which @for reports
   * as NG0100 in dev mode.
   */
  public readonly accordionPanels = computed<AccordionPanel[]>(() => {
    const panels: AccordionPanel[] = [];

    this.childTasks().forEach(task => {
      if (task.type === SystemTaskType.EmailUpload) {
        const settings = this.allReceiptProcessingSettings().find(s => s.id === task.associatedEntityId);
        let description = "";

        if (settings) {
          description = `Used ${settings?.name} to process Receipt`;
        }


        panels.push({
          title: 'Email Upload Details',
          description: description,
          content: this.emailUploadDetails(),
        });
      }
    })

    return panels;
  });
}
