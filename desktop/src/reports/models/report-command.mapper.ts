import { format } from "date-fns";
import {
  ReceiptPagedRequestFilter,
  ReportColumn,
  ReportDetail,
  ReportPeriod,
  ReportRequestCommand,
} from "../../open-api";
import { ReportRequestFormat } from "./report-catalog.constants";

/** A single builder column, before mapping to the engine's ReportColumn. */
export interface ReportColumnValue {
  // A stable client id used for @for tracking and formula references in the UI.
  id: string;
  kind: ReportColumn.KindEnum;
  // The engine machine name (a plain identifier); formulas reference it.
  name: string;
  label: string;
  field?: string;
  aggFunc?: ReportColumn.AggFuncEnum;
  measure?: string;
  expr?: string;
}

/** The builder form's resolved value, shaped for mapping to a request command. */
export interface ReportBuilderValue {
  name: string;
  scope: string[];
  period: {
    preset: ReportPeriod.PresetEnum;
    startDate: Date | null;
    endDate: Date | null;
  };
  filter: ReceiptPagedRequestFilter;
  groupBy: string[];
  detail: { mode: ReportDetail.ModeEnum; by: string };
  columns: ReportColumnValue[];
  subtotals: boolean;
  grandTotals: boolean;
  document: { title: string; intro: string; footer: string };
  formats: Record<ReportRequestFormat, boolean>;
}

/**
 * A dimension column is "disabled" in aggregate mode when it reads a field that
 * isn't the aggregate-by dimension or one of the grouping levels — an aggregated
 * row is a summary of many receipts, so only those fields have a single value on
 * it (the engine rejects anything else). It's a derived state: records mode and
 * aggregate/formula columns are never disabled, so changing the config recomputes
 * it. Disabled columns are shown greyed in the builder and left out of the request.
 */
export function isDimensionColumnDisabled(
  column: ReportColumnValue,
  mode: ReportDetail.ModeEnum,
  detailBy: string,
  groupBy: string[]
): boolean {
  if (mode !== ReportDetail.ModeEnum.Aggregate || column.kind !== ReportColumn.KindEnum.Dimension) {
    return false;
  }
  const field = column.field ?? "";
  return field !== detailBy && !groupBy.includes(field);
}

/** The columns actually sent to the engine — every column minus the disabled ones. */
export function enabledReportColumns(value: ReportBuilderValue): ReportColumnValue[] {
  return value.columns.filter(
    (column) => !isDimensionColumnDisabled(column, value.detail.mode, value.detail.by, value.groupBy)
  );
}

/** Formats the picked date range into the YYYY-MM-DD strings the API expects. */
function toApiDate(date: Date | null): string {
  return date ? format(date, "yyyy-MM-dd") : "";
}

function toPeriod(period: ReportBuilderValue["period"]): ReportPeriod {
  if (period.preset === ReportPeriod.PresetEnum.Custom) {
    return {
      preset: period.preset,
      startDate: toApiDate(period.startDate),
      endDate: toApiDate(period.endDate),
    };
  }
  return { preset: period.preset };
}

/** Maps a builder column to an engine column, carrying only the fields its kind uses. */
function toColumn(column: ReportColumnValue): ReportColumn {
  const mapped: ReportColumn = {
    kind: column.kind,
    name: column.name,
    label: column.label,
  };
  switch (column.kind) {
    case ReportColumn.KindEnum.Dimension:
      mapped.field = column.field;
      break;
    case ReportColumn.KindEnum.Aggregate:
      mapped.aggFunc = column.aggFunc;
      if (column.aggFunc !== ReportColumn.AggFuncEnum.Count) {
        mapped.measure = column.measure;
      }
      break;
    case ReportColumn.KindEnum.Formula:
      mapped.expr = column.expr;
      break;
  }
  return mapped;
}

/** Turns the checkbox map into the ordered format list the endpoint expects. */
function toFormats(formats: Record<ReportRequestFormat, boolean>): ReportRequestCommand.FormatsEnum[] {
  const order: ReportRequestFormat[] = ["csv", "xlsx", "pdf"];
  return order
    .filter((format) => formats[format])
    .map((format) => format as ReportRequestCommand.FormatsEnum);
}

/**
 * Maps the builder form's value onto the ReportRequestCommand the endpoint
 * consumes. Detail `by` is only sent in aggregate mode; the document is omitted
 * when entirely empty. The same command feeds both preview and generate.
 */
export function toReportRequestCommand(value: ReportBuilderValue): ReportRequestCommand {
  const detail: ReportDetail = { mode: value.detail.mode };
  if (value.detail.mode === ReportDetail.ModeEnum.Aggregate) {
    detail.by = value.detail.by;
  }

  const command: ReportRequestCommand = {
    name: value.name,
    groupIds: value.scope,
    period: toPeriod(value.period),
    filter: value.filter,
    groupBy: value.groupBy,
    detail,
    columns: enabledReportColumns(value).map(toColumn),
    subtotals: value.subtotals,
    grandTotals: value.grandTotals,
    formats: toFormats(value.formats),
  };

  const { title, intro, footer } = value.document;
  if (title || intro || footer) {
    command.document = { title, intro, footer };
  }

  return command;
}
