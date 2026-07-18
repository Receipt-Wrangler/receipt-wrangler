import { WidgetType } from "../open-api";
import { WidgetTypePipe } from './widget-type.pipe';

describe('WidgetTypePipe', () => {
  const pipe = new WidgetTypePipe();

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  it('maps each widget type to its display label', () => {
    expect(pipe.transform(WidgetType.GroupSummary)).toBe("Group Summary");
    expect(pipe.transform(WidgetType.FilteredReceipts)).toBe("Filtered Receipts");
    expect(pipe.transform(WidgetType.GroupActivity)).toBe("Activity");
    expect(pipe.transform(WidgetType.PieChart)).toBe("Pie Chart");
    expect(pipe.transform(WidgetType.Report)).toBe("Report");
  });
});
