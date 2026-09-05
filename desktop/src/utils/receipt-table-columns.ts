import { CustomField } from "../open-api";
import {
  DEFAULT_RECEIPT_TABLE_COLUMNS,
  ReceiptTableColumnConfig,
} from "../interfaces/receipt-table-column-config.interface";

/**
 * Prefix a custom field column is identified by, shared with the backend's
 * `receiptsource.CustomFieldKey`. The same string is the column's `matColumnDef`
 * and the `orderBy` the API is asked to sort on, so both ends must agree.
 */
const CUSTOM_FIELD_COLUMN_PREFIX = "custom_";

/** The header shown for each built-in column, in the table and in the dialog. */
export const RECEIPT_COLUMN_DISPLAY_NAMES: Readonly<Record<string, string>> = {
  created_at: "Added At",
  date: "Receipt Date",
  name: "Name",
  paid_by_user_id: "Paid By",
  amount: "Amount",
  categories: "Categories",
  tags: "Tags",
  status: "Status",
  resolved_date: "Resolved Date",
};

export function customFieldColumnDef(customFieldId: number): string {
  return `${CUSTOM_FIELD_COLUMN_PREFIX}${customFieldId}`;
}

/**
 * The custom field id a column refers to, or undefined when the column is not a
 * custom field one. Deliberately strict — digits only — so it matches the id the
 * backend will parse out of the same string.
 */
export function parseCustomFieldColumnDef(
  matColumnDef: string
): number | undefined {
  if (!matColumnDef.startsWith(CUSTOM_FIELD_COLUMN_PREFIX)) {
    return undefined;
  }

  const digits = matColumnDef.slice(CUSTOM_FIELD_COLUMN_PREFIX.length);
  if (!/^\d+$/.test(digits)) {
    return undefined;
  }

  return Number(digits);
}

export function columnDisplayName(
  matColumnDef: string,
  customFields: CustomField[]
): string {
  const builtInName = RECEIPT_COLUMN_DISPLAY_NAMES[matColumnDef];
  if (builtInName) {
    return builtInName;
  }

  const customFieldId = parseCustomFieldColumnDef(matColumnDef);
  const customField = customFields.find((field) => field.id === customFieldId);

  return customField?.name ?? matColumnDef;
}

/**
 * Reconciles a persisted column configuration with the custom fields that exist
 * right now.
 *
 * The configuration lives in localStorage, so it outlives the catalog it was
 * written against — a custom field can be deleted, a new one added, and the
 * whole slice is shared by every account using the same browser. Anything that
 * no longer resolves is dropped (leaving it in place would make the table ask
 * the API to sort on a column that no longer exists), and anything new is
 * appended **hidden**, so creating a custom field never silently widens
 * everyone's table.
 *
 * Order is otherwise preserved exactly as persisted, including a custom field
 * the user dragged above a built-in column.
 */
export function mergeCustomFieldColumns(
  persisted: ReceiptTableColumnConfig[],
  customFields: CustomField[]
): ReceiptTableColumnConfig[] {
  const knownCustomFieldIds = new Set(customFields.map((field) => field.id));
  const builtInColumnDefs = DEFAULT_RECEIPT_TABLE_COLUMNS.map(
    (column) => column.matColumnDef
  );

  const kept = [...persisted]
    .sort((a, b) => a.order - b.order)
    .filter((column) => {
      const customFieldId = parseCustomFieldColumnDef(column.matColumnDef);
      return customFieldId === undefined
        ? builtInColumnDefs.includes(column.matColumnDef)
        : knownCustomFieldIds.has(customFieldId);
    });

  const present = new Set(kept.map((column) => column.matColumnDef));

  const missingBuiltIns = DEFAULT_RECEIPT_TABLE_COLUMNS.filter(
    (column) => !present.has(column.matColumnDef)
  );

  const missingCustomFields = [...customFields]
    .sort((a, b) => a.name.localeCompare(b.name))
    .filter((field) => !present.has(customFieldColumnDef(field.id)))
    .map((field) => ({
      matColumnDef: customFieldColumnDef(field.id),
      visible: false,
    }));

  return [...kept, ...missingBuiltIns, ...missingCustomFields].map(
    (column, index) => ({
      matColumnDef: column.matColumnDef,
      visible: column.visible,
      order: index,
    })
  );
}
