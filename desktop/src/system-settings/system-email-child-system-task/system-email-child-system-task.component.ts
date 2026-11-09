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
   * The children that get a panel, in panel order. The template reads the task for panel `i`
   * from HERE, not from `childTasks()`: an EMAIL_READ also carries an HTML_TO_PDF child when the
   * group processes email bodies, created ahead of the uploads, so indexing the full list by
   * panel position handed an upload's panel the render task and showed it empty.
   */
  public readonly emailUploadTasks = computed(() =>
    this.childTasks().filter(task => task.type === SystemTaskType.EmailUpload)
  );

  /**
   * Derived rather than pushed in ngAfterViewInit: an array filled after the accordion had
   * already been checked only rendered on a later change-detection pass, which @for reports
   * as NG0100 in dev mode.
   */
  public readonly accordionPanels = computed<AccordionPanel[]>(() =>
    this.emailUploadTasks().map(task => {
      const settings = this.allReceiptProcessingSettings().find(s => s.id === task.associatedEntityId);

      return {
        title: 'Email Upload Details',
        description: settings ? `Used ${settings.name} to process Receipt` : "",
        content: this.emailUploadDetails(),
      };
    })
  );
}
