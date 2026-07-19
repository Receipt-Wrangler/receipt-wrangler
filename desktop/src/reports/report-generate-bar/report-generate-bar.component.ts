import { ChangeDetectionStrategy, Component, computed, input, output } from "@angular/core";
import { FormGroup } from "@angular/forms";

interface FormatChip {
  key: string;
  label: string;
  icon: string;
}

/**
 * The bottom action bar: output-format selection and the (synchronous) Generate
 * button. There is no progress bar or cancel — generation streams the file back,
 * so the button just shows an in-flight state while the request is running.
 */
@Component({
  selector: "app-report-generate-bar",
  templateUrl: "./report-generate-bar.component.html",
  styleUrls: ["./report-generate-bar.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportGenerateBarComponent {
  public readonly form = input.required<FormGroup>();
  public readonly generating = input<boolean>(false);
  public readonly canGenerate = input<boolean>(false);
  public readonly canSaveTemplate = input<boolean>(false);
  public readonly saving = input<boolean>(false);
  // Editing an opened template updates it in place; the blank builder creates a new
  // one — isEditMode drives the Save action's label. The parent resolves the Save and
  // Generate permission gates (each honoring the base + "*All" variant) into these
  // booleans, so this bar just renders the actions it is told the user may perform.
  public readonly isEditMode = input<boolean>(false);
  public readonly saveTemplateAllowed = input<boolean>(false);
  public readonly generateAllowed = input<boolean>(false);
  public readonly generate = output<void>();
  public readonly saveTemplate = output<void>();

  /** The Save button's label, reflecting create-vs-update and the in-flight state. */
  public readonly saveTemplateText = computed<string>(() => {
    if (this.isEditMode()) {
      return this.saving() ? "Updating…" : "Update Template";
    }
    return this.saving() ? "Saving…" : "Save Template";
  });

  public readonly formats: FormatChip[] = [
    { key: "csv", label: "CSV", icon: "description" },
    { key: "xlsx", label: "XLSX", icon: "table_chart" },
    { key: "pdf", label: "PDF", icon: "picture_as_pdf" },
  ];

  public isSelected(key: string): boolean {
    return !!this.form().get("formats." + key)?.value;
  }

  public toggle(key: string): void {
    const control = this.form().get("formats." + key);
    control?.setValue(!control.value);
  }

  public formatSummary(): string {
    const selected = this.formats.filter((format) => this.isSelected(format.key)).map((format) => format.label);
    if (selected.length === 0) {
      return "Pick at least one format";
    }
    if (selected.length > 1) {
      return selected.join(" + ") + " → zipped";
    }
    return "Single " + selected[0] + " file";
  }

  public onGenerate(): void {
    this.generate.emit();
  }

  public onSaveTemplate(): void {
    this.saveTemplate.emit();
  }
}
