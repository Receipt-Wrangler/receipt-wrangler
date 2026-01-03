import { EmbeddedViewRef, Injectable, TemplateRef } from "@angular/core";
import { MatSnackBar, MatSnackBarConfig, MatSnackBarRef, } from "@angular/material/snack-bar";
import { DEFAULT_SNACKBAR_ACTION, DEFAULT_SNACKBAR_CONFIG, } from "../constants/snackbar.constant";
import { SnackbarServiceInterface } from "../interfaces/snackbar.interface";

@Injectable({
  providedIn: "root",
})
export class SnackbarService implements SnackbarServiceInterface {
  constructor(private snackbar: MatSnackBar) {}

  public info(message: string): void {
    this.snackbar.open(message, DEFAULT_SNACKBAR_ACTION, {
      ...DEFAULT_SNACKBAR_CONFIG,
      duration: undefined,
      panelClass: ["info-snackbar"],
    });
  }

  public error(message: string): void {
    this.snackbar.open(message, DEFAULT_SNACKBAR_ACTION, {
      ...DEFAULT_SNACKBAR_CONFIG,
      panelClass: ["error-snackbar"],
    });
  }

  public success(
    message: string,
    configOverrides?: MatSnackBarConfig<any>
  ): void {
    this.snackbar.open(message, DEFAULT_SNACKBAR_ACTION, {
      ...DEFAULT_SNACKBAR_CONFIG,
      ...configOverrides,
      panelClass: ["success-snackbar"],
    });
  }

  public successFromTemplate(
    template: TemplateRef<any>,
    configOverrides?: MatSnackBarConfig<any>
  ): MatSnackBarRef<EmbeddedViewRef<any>> {
    return this.snackbar.openFromTemplate(template, {
      ...DEFAULT_SNACKBAR_CONFIG,
      ...configOverrides,
      panelClass: ["success-snackbar"],
    });
  }
}
