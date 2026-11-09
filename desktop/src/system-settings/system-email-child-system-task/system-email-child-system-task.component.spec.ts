import {ComponentFixture, TestBed} from '@angular/core/testing';

import {SystemEmailChildSystemTaskComponent} from './system-email-child-system-task.component';
import {CUSTOM_ELEMENTS_SCHEMA} from "@angular/core";
import {ReceiptProcessingSettings, SystemTask, SystemTaskType} from "../../open-api";

function task(id: number, type: SystemTaskType, extra: Partial<SystemTask> = {}): SystemTask {
  return {id, type, ...extra} as SystemTask;
}

const settings = (id: number, name: string) => ({id, name}) as ReceiptProcessingSettings;

describe('SystemEmailChildSystemTaskComponent', () => {
  let component: SystemEmailChildSystemTaskComponent;
  let fixture: ComponentFixture<SystemEmailChildSystemTaskComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SystemEmailChildSystemTaskComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();

    fixture = TestBed.createComponent(SystemEmailChildSystemTaskComponent);
    component = fixture.componentInstance;
  });

  // Inputs are set BEFORE the first detectChanges, so the template's accordionPanels() read is
  // also the proof that the viewChild.required template resolves on first render.
  function render(childTasks: SystemTask[], allSettings: ReceiptProcessingSettings[] = []) {
    fixture.componentRef.setInput('childTasks', childTasks);
    fixture.componentRef.setInput('allReceiptProcessingSettings', allSettings);
    fixture.detectChanges();
  }

  it('should create', () => {
    render([]);
    expect(component).toBeTruthy();
    expect(component.accordionPanels()).toEqual([]);
  });

  it('builds an Email Upload Details panel naming the settings that processed it', () => {
    render([task(1, SystemTaskType.EmailUpload, {associatedEntityId: 5})], [settings(5, 'Primary AI')]);

    expect(component.accordionPanels()).toEqual([
      {
        title: 'Email Upload Details',
        description: 'Used Primary AI to process Receipt',
        content: component.emailUploadDetails(),
      },
    ]);
  });

  it('gives no panel to a child that is not an email upload', () => {
    render([task(1, SystemTaskType.HtmlToPdf), task(2, SystemTaskType.EmailUpload, {associatedEntityId: 5})]);

    expect(component.accordionPanels().length).toBe(1);
    expect(component.emailUploadTasks().map(t => t.id)).toEqual([2]);
  });

  it('leaves the description empty when no settings match', () => {
    render([task(1, SystemTaskType.EmailUpload, {associatedEntityId: 99})], [settings(5, 'Primary AI')]);

    expect(component.accordionPanels()[0].description).toBe('');
  });

  it('re-derives the description when the settings input changes', () => {
    render([task(1, SystemTaskType.EmailUpload, {associatedEntityId: 5})], []);
    expect(component.accordionPanels()[0].description).toBe('');

    fixture.componentRef.setInput('allReceiptProcessingSettings', [settings(5, 'Primary AI')]);
    fixture.detectChanges();

    expect(component.accordionPanels()[0].description).toBe('Used Primary AI to process Receipt');
  });

  it("shows each upload's own processing tasks even when another child comes first", () => {
    // An EMAIL_READ that processed the email body carries an HTML_TO_PDF child, created ahead
    // of the uploads. Panel 0 is the upload, so its content must read the upload's children.
    const uploadChildren = [task(11, SystemTaskType.PromptGenerated), task(12, SystemTaskType.ChatCompletion)];
    render([
      task(1, SystemTaskType.HtmlToPdf, {childSystemTasks: []}),
      task(2, SystemTaskType.EmailUpload, {associatedEntityId: 5, childSystemTasks: uploadChildren}),
    ]);

    // Render panel 0's content the way the accordion does: its template, with the panel index.
    const view = component.accordionPanels()[0].content.createEmbeddedView({index: 0});
    view.detectChanges();
    const nested = view.rootNodes.find(
      (node: Node) => node.nodeName === 'APP-RECEIPT-PROCESSING-SETTINGS-CHILD-SYSTEM-TASK-ACCORDION'
    );

    expect(nested.childTasks).toEqual(uploadChildren);
    view.destroy();
  });
});
