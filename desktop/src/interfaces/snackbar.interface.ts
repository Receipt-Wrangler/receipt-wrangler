import { TemplateRef } from '@angular/core';

export interface SnackbarServiceInterface {
  error(message: string): void;
  success(message: string, configOverrides?: any): void;
  successFromTemplate(template: TemplateRef<any>, configOverrides?: any): any;
}
