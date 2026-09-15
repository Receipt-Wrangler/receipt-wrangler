import { CustomField, CustomFieldType } from "../open-api";
import { DEFAULT_RECEIPT_TABLE_COLUMNS } from "../interfaces/receipt-table-column-config.interface";
import {
  columnDisplayName,
  customFieldColumnDef,
  mergeCustomFieldColumns,
  parseCustomFieldColumnDef,
} from "./receipt-table-columns";

function customField(id: number, name: string): CustomField {
  return { id, name, type: CustomFieldType.Text } as CustomField;
}

function config(matColumnDef: string, visible: boolean, order: number) {
  return { matColumnDef, visible, order };
}

/** The built-in columns as persisted, so a case can focus on the custom ones. */
function builtIns() {
  return DEFAULT_RECEIPT_TABLE_COLUMNS.map((column, index) =>
    config(column.matColumnDef, column.visible, index)
  );
}

describe("receipt-table-columns", () => {
  describe("parseCustomFieldColumnDef", () => {
    it("parses a custom field column", () => {
      expect(parseCustomFieldColumnDef("custom_7")).toBe(7);
      expect(parseCustomFieldColumnDef(customFieldColumnDef(42))).toBe(42);
    });

    it("rejects anything that is not the prefix followed by digits alone", () => {
      // Mirrors the backend's receiptsource.ParseCustomFieldKey, which the same
      // string is sent to as an orderBy.
      ["date", "custom_", "custom_abc", "custom_1_month", "custom_+1", ""].forEach(
        (value) => expect(parseCustomFieldColumnDef(value)).toBeUndefined()
      );
    });
  });

  describe("columnDisplayName", () => {
    it("names a built-in column and a custom field", () => {
      expect(columnDisplayName("resolved_date", [])).toBe("Resolved Date");
      expect(columnDisplayName("custom_7", [customField(7, "Vendor")])).toBe("Vendor");
    });

    it("falls back to the raw column def for an unresolvable custom field", () => {
      expect(columnDisplayName("custom_7", [])).toBe("custom_7");
    });
  });

  describe("mergeCustomFieldColumns", () => {
    it("appends a newly created custom field hidden", () => {
      const merged = mergeCustomFieldColumns(builtIns(), [customField(7, "Vendor")], true);

      const appended = merged.find((column) => column.matColumnDef === "custom_7");
      expect(appended).toBeDefined();
      expect(appended?.visible).toBe(false);
    });

    it("drops a custom field the catalog no longer lists", () => {
      const persisted = [...builtIns(), config("custom_7", true, 99)];

      const merged = mergeCustomFieldColumns(persisted, [customField(3, "Other")], true);

      expect(merged.some((column) => column.matColumnDef === "custom_7")).toBe(false);
    });

    it("preserves the persisted order and re-derives it contiguously", () => {
      const persisted = [
        config("custom_7", true, 0),
        ...builtIns().map((column, index) => ({ ...column, order: index + 1 })),
      ];

      const merged = mergeCustomFieldColumns(persisted, [customField(7, "Vendor")], true);

      expect(merged[0].matColumnDef).toBe("custom_7");
      expect(merged.map((column) => column.order)).toEqual(merged.map((_, i) => i));
    });

    it("heals a missing built-in column", () => {
      const persisted = builtIns().filter((column) => column.matColumnDef !== "status");

      const merged = mergeCustomFieldColumns(persisted, [], true);

      expect(merged.some((column) => column.matColumnDef === "status")).toBe(true);
    });

    describe("when the catalog is unavailable", () => {
      // An empty catalog then means "this user may not read custom fields", not
      // "there are none" - and the configuration is shared by every account on the
      // browser, so dropping against it discards an administrator's saved layout.
      it("keeps persisted custom field columns, with their visibility and order", () => {
        const persisted = [
          config("custom_7", true, 0),
          ...builtIns().map((column, index) => ({ ...column, order: index + 1 })),
          config("custom_9", false, 99),
        ];

        const merged = mergeCustomFieldColumns(persisted, [], false);

        const kept = merged.find((column) => column.matColumnDef === "custom_7");
        expect(kept).toBeDefined();
        expect(kept?.visible).toBe(true);
        expect(merged[0].matColumnDef).toBe("custom_7");
        expect(merged.some((column) => column.matColumnDef === "custom_9")).toBe(true);
      });

      it("still heals missing built-ins, which need no catalog", () => {
        const persisted = builtIns().filter((column) => column.matColumnDef !== "status");

        const merged = mergeCustomFieldColumns(persisted, [], false);

        expect(merged.some((column) => column.matColumnDef === "status")).toBe(true);
      });

      it("still drops a column that is not a custom field and not built in", () => {
        const persisted = [...builtIns(), config("bogus_column", true, 99)];

        const merged = mergeCustomFieldColumns(persisted, [], false);

        expect(merged.some((column) => column.matColumnDef === "bogus_column")).toBe(false);
      });
    });
  });
});
