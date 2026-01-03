import {TemplateRef} from "@angular/core";

export interface AccordionPanel {
  title: string;
  description?: string;
  descriptionTemplate?: TemplateRef<any>;
  content: TemplateRef<any>;
}
