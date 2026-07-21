import { ChangeDetectionStrategy, Component, input } from "@angular/core";

/**
 * The card shell every report-builder config section shares: an accent icon, a
 * title and subtitle, an optional actions slot ([sectionActions]), and projected
 * body content. Extracted so the section-card pattern is defined once.
 */
@Component({
  selector: "app-report-section",
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="report-section">
      <div class="report-section__header">
        <span class="report-section__icon">
          <mat-icon>{{ icon() }}</mat-icon>
        </span>
        <div class="report-section__heading">
          <div class="report-section__title">{{ headerText() }}</div>
          <div class="report-section__subtitle">{{ subtitle() }}</div>
        </div>
        <span class="report-section__actions">
          <ng-content select="[sectionActions]"></ng-content>
        </span>
      </div>
      <ng-content></ng-content>
    </div>
  `,
  styles: [
    `
      .report-section {
        background: #fff;
        border: 1px solid rgba(0, 0, 0, 0.06);
        border-radius: 12px;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
        padding: 15px 17px;
        margin-bottom: 13px;
      }
      .report-section__header {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-bottom: 13px;
      }
      .report-section__icon {
        width: 30px;
        height: 30px;
        border-radius: 8px;
        background: rgba(39, 177, 255, 0.1);
        display: inline-flex;
        align-items: center;
        justify-content: center;
        flex: none;
      }
      .report-section__icon mat-icon {
        font-size: 18px;
        width: 18px;
        height: 18px;
        color: #009efa;
      }
      .report-section__heading {
        flex: 1;
        min-width: 0;
      }
      .report-section__title {
        font-size: 13.5px;
        font-weight: 600;
        color: #0f172a;
      }
      .report-section__subtitle {
        font-size: 11.5px;
        color: #94a3b8;
      }
    `,
  ],
})
export class ReportSectionComponent {
  public readonly icon = input<string>("");
  public readonly headerText = input<string>("");
  public readonly subtitle = input<string>("");
}
