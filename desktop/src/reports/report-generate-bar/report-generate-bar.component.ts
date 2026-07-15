import { ChangeDetectionStrategy, Component, input, output } from "@angular/core";
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
  public readonly generate = output<void>();

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
}
