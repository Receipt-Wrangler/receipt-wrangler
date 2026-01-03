import { CdkDragDrop } from "@angular/cdk/drag-drop";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { NgxsModule } from "@ngxs/store";
import { DEFAULT_RECEIPT_TABLE_COLUMNS, ReceiptTableColumnConfig } from "../../interfaces";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { ReceiptTableState } from "../../store/receipt-table.state";
import { ColumnConfigurationDialogComponent } from "./column-configuration-dialog.component";

describe("ColumnConfigurationDialogComponent", () => {
  let component: ColumnConfigurationDialogComponent;
  let fixture: ComponentFixture<ColumnConfigurationDialogComponent>;
  let mockDialogRef: jest.Mocked<MatDialogRef<ColumnConfigurationDialogComponent>>;
  let mockDialogData: { currentColumns?: ReceiptTableColumnConfig[] };

  beforeEach(async () => {
    // Create spy for MatDialogRef
    mockDialogRef = { close: jest.fn() } as unknown as jest.Mocked<MatDialogRef<ColumnConfigurationDialogComponent>>;

    // Initialize mock dialog data
    mockDialogData = {
      currentColumns: [
        { matColumnDef: "created_at", visible: true, order: 0 },
        { matColumnDef: "name", visible: true, order: 1 },
        { matColumnDef: "amount", visible: false, order: 2 }
      ]
    };

    await TestBed.configureTestingModule({
      declarations: [ColumnConfigurationDialogComponent],
      imports: [SharedUiModule, NgxsModule.forRoot([ReceiptTableState])],
      providers: [
        { provide: MatDialogRef, useValue: mockDialogRef },
        { provide: MAT_DIALOG_DATA, useValue: mockDialogData }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ColumnConfigurationDialogComponent);
    component = fixture.componentInstance;
  });

  describe("Component Initialization", () => {
    it("should create", () => {
      expect(component).toBeTruthy();
    });

    it("should have correct columnDisplayNames mapping", () => {
      const expectedMapping = {
        "created_at": "Added At",
        "date": "Receipt Date",
        "name": "Name",
        "paid_by_user_id": "Paid By",
        "amount": "Amount",
        "categories": "Categories",
        "tags": "Tags",
        "status": "Status",
        "resolved_date": "Resolved Date"
      };

      expect(component["columnDisplayNames"]).toEqual(expectedMapping);
    });

    it("should initialize columns on ngOnInit", () => {
      jest.spyOn(component as any, "initializeColumns");

      component.ngOnInit();

      expect(component["initializeColumns"]).toHaveBeenCalled();
    });
  });

  describe("initializeColumns", () => {
    it("should initialize columns from provided data", () => {
      component.ngOnInit();

      expect(component.columns.length).toEqual(3);
      expect(component.columns[0]).toEqual({
        matColumnDef: "created_at",
        visible: true,
        order: 0,
        displayName: "Added At"
      });
      expect(component.columns[1]).toEqual({
        matColumnDef: "name",
        visible: true,
        order: 1,
        displayName: "Name"
      });
      expect(component.columns[2]).toEqual({
        matColumnDef: "amount",
        visible: false,
        order: 2,
        displayName: "Amount"
      });
    });

    it("should sort columns by order before mapping", () => {
      mockDialogData.currentColumns = [
        { matColumnDef: "amount", visible: true, order: 2 },
        { matColumnDef: "created_at", visible: true, order: 0 },
        { matColumnDef: "name", visible: true, order: 1 }
      ];

      component.ngOnInit();

      expect(component.columns[0].matColumnDef).toBe("created_at");
      expect(component.columns[1].matColumnDef).toBe("name");
      expect(component.columns[2].matColumnDef).toBe("amount");
    });

    it("should use DEFAULT_RECEIPT_TABLE_COLUMNS when no currentColumns provided", () => {
      mockDialogData.currentColumns = undefined;

      component.ngOnInit();

      expect(component.columns.length).toEqual(DEFAULT_RECEIPT_TABLE_COLUMNS.length);
      expect(component.columns[0].matColumnDef).toBe(DEFAULT_RECEIPT_TABLE_COLUMNS[0].matColumnDef);
    });

    it("should use matColumnDef as displayName when no mapping exists", () => {
      mockDialogData.currentColumns = [
        { matColumnDef: "unknown_column", visible: true, order: 0 }
      ];

      component.ngOnInit();

      expect(component.columns[0].displayName).toBe("unknown_column");
    });
  });

  describe("toggleColumnVisibility", () => {
    beforeEach(() => {
      component.ngOnInit();
    });

    it("should toggle visible column to invisible", () => {
      const visibleColumn = component.columns.find(col => col.visible)!;
      const originalVisibility = visibleColumn.visible;

      component.toggleColumnVisibility(visibleColumn);

      expect(visibleColumn.visible).toBe(!originalVisibility);
    });

    it("should toggle invisible column to visible", () => {
      const invisibleColumn = component.columns.find(col => !col.visible)!;
      const originalVisibility = invisibleColumn.visible;

      component.toggleColumnVisibility(invisibleColumn);

      expect(invisibleColumn.visible).toBe(!originalVisibility);
    });

    it("should not affect other columns when toggling one column", () => {
      const targetColumn = component.columns[0];
      const otherColumns = component.columns.slice(1);
      const originalStates = otherColumns.map(col => col.visible);

      component.toggleColumnVisibility(targetColumn);

      otherColumns.forEach((col, index) => {
        expect(col.visible).toBe(originalStates[index]);
      });
    });
  });

  describe("drop", () => {
    beforeEach(() => {
      component.ngOnInit();
    });

    it("should reorder columns when dropped", () => {
      const originalOrder = [...component.columns];
      const mockEvent: CdkDragDrop<any[]> = {
        previousIndex: 0,
        currentIndex: 2,
        item: null as any,
        container: null as any,
        previousContainer: null as any,
        isPointerOverContainer: false,
        distance: { x: 0, y: 0 },
        dropPoint: { x: 0, y: 0 },
        event: null as any
      };

      component.drop(mockEvent);

      expect(component.columns[0]).toEqual({ ...originalOrder[1], order: 0 });
      expect(component.columns[1]).toEqual({ ...originalOrder[2], order: 1 });
      expect(component.columns[2]).toEqual({ ...originalOrder[0], order: 2 });
    });

    it("should update order property for all columns after drop", () => {
      const mockEvent: CdkDragDrop<any[]> = {
        previousIndex: 1,
        currentIndex: 0,
        item: null as any,
        container: null as any,
        previousContainer: null as any,
        isPointerOverContainer: false,
        distance: { x: 0, y: 0 },
        dropPoint: { x: 0, y: 0 },
        event: null as any
      };

      component.drop(mockEvent);

      component.columns.forEach((column, index) => {
        expect(column.order).toBe(index);
      });
    });

    it("should handle drop from last to first position", () => {
      const originalLastColumn = component.columns[component.columns.length - 1];
      const mockEvent: CdkDragDrop<any[]> = {
        previousIndex: component.columns.length - 1,
        currentIndex: 0,
        item: null as any,
        container: null as any,
        previousContainer: null as any,
        isPointerOverContainer: false,
        distance: { x: 0, y: 0 },
        dropPoint: { x: 0, y: 0 },
        event: null as any
      };

      component.drop(mockEvent);

      expect(component.columns[0].matColumnDef).toBe(originalLastColumn.matColumnDef);
      expect(component.columns[0].order).toBe(0);
    });
  });

  describe("resetToDefaults", () => {
    beforeEach(() => {
      component.ngOnInit();
    });

    it("should reset columns to DEFAULT_RECEIPT_TABLE_COLUMNS", () => {
      // Modify columns first
      component.columns[0].visible = false;
      component.columns[1].order = 999;

      component.resetToDefaults();

      expect(component.columns.length).toEqual(DEFAULT_RECEIPT_TABLE_COLUMNS.length);
      DEFAULT_RECEIPT_TABLE_COLUMNS.forEach((defaultCol, index) => {
        expect(component.columns[index].matColumnDef).toBe(defaultCol.matColumnDef);
        expect(component.columns[index].visible).toBe(defaultCol.visible);
        expect(component.columns[index].order).toBe(defaultCol.order);
      });
    });

    it("should add display names to default columns", () => {
      component.resetToDefaults();

      component.columns.forEach(column => {
        const expectedDisplayName = component["columnDisplayNames"][column.matColumnDef] || column.matColumnDef;
        expect(column.displayName).toBe(expectedDisplayName);
      });
    });

    it("should preserve displayName mapping when resetting", () => {
      component.resetToDefaults();

      const createdAtColumn = component.columns.find(col => col.matColumnDef === "created_at");
      expect(createdAtColumn?.displayName).toBe("Added At");

      const nameColumn = component.columns.find(col => col.matColumnDef === "name");
      expect(nameColumn?.displayName).toBe("Name");
    });
  });

  describe("saveConfiguration", () => {
    beforeEach(() => {
      component.ngOnInit();
    });

    it("should close dialog with configuration result", () => {
      component.saveConfiguration();

      expect(mockDialogRef.close).toHaveBeenCalledWith(
        expect.arrayContaining([
          expect.objectContaining({
            matColumnDef: expect.any(String),
            visible: expect.any(Boolean),
            order: expect.any(Number)
          })
        ])
      );
    });

    it("should remove displayName property from result", () => {
      component.saveConfiguration();

      const result = mockDialogRef.close.mock.calls[mockDialogRef.close.mock.calls.length - 1][0];
      result.forEach((col: any) => {
        expect(col.displayName).toBeUndefined();
        expect(col.matColumnDef).toBeDefined();
        expect(col.visible).toBeDefined();
        expect(col.order).toBeDefined();
      });
    });

    it("should preserve all other column properties in result", () => {
      component.saveConfiguration();

      const result = mockDialogRef.close.mock.calls[mockDialogRef.close.mock.calls.length - 1][0];
      expect(result.length).toEqual(component.columns.length);

      result.forEach((resultCol: ReceiptTableColumnConfig, index: number) => {
        const originalCol = component.columns[index];
        expect(resultCol.matColumnDef).toBe(originalCol.matColumnDef);
        expect(resultCol.visible).toBe(originalCol.visible);
        expect(resultCol.order).toBe(originalCol.order);
      });
    });
  });

  describe("cancel", () => {
    it("should close dialog with null result", () => {
      component.cancel();

      expect(mockDialogRef.close).toHaveBeenCalledWith(null);
    });
  });

  describe("visibleColumnsCount", () => {
    beforeEach(() => {
      component.ngOnInit();
    });

    it("should return correct count of visible columns", () => {
      const visibleCount = component.columns.filter(col => col.visible).length;

      expect(component.visibleColumnsCount).toBe(visibleCount);
    });

    it("should update when column visibility changes", () => {
      const initialCount = component.visibleColumnsCount;
      const invisibleColumn = component.columns.find(col => !col.visible);

      if (invisibleColumn) {
        component.toggleColumnVisibility(invisibleColumn);
        expect(component.visibleColumnsCount).toBe(initialCount + 1);
      }
    });

    it("should return 0 when all columns are invisible", () => {
      component.columns.forEach(col => col.visible = false);

      expect(component.visibleColumnsCount).toBe(0);
    });

    it("should return total count when all columns are visible", () => {
      component.columns.forEach(col => col.visible = true);

      expect(component.visibleColumnsCount).toBe(component.columns.length);
    });
  });

  describe("canToggleOff", () => {
    beforeEach(() => {
      component.ngOnInit();
    });

    it("should return true when more than one column is visible", () => {
      // Ensure at least 2 columns are visible
      component.columns[0].visible = true;
      component.columns[1].visible = true;
      component.columns[2].visible = false;

      expect(component.canToggleOff(component.columns[0])).toBe(true);
    });

    it("should return false when only one column is visible and trying to toggle it off", () => {
      // Set only one column as visible
      component.columns.forEach((col, index) => {
        col.visible = index === 0;
      });

      expect(component.canToggleOff(component.columns[0])).toBe(false);
    });

    it("should return true for invisible columns regardless of visible count", () => {
      // Set only one column as visible
      component.columns.forEach((col, index) => {
        col.visible = index === 0;
      });

      const invisibleColumn = component.columns.find(col => !col.visible)!;
      expect(component.canToggleOff(invisibleColumn)).toBe(true);
    });

    it("should return true when trying to toggle off visible column with multiple visible columns", () => {
      component.columns.forEach(col => col.visible = true);

      expect(component.canToggleOff(component.columns[0])).toBe(true);
    });

    it("should handle edge case with no visible columns", () => {
      component.columns.forEach(col => col.visible = false);

      expect(component.canToggleOff(component.columns[0])).toBe(true);
    });
  });

  describe("Component Integration", () => {
    it("should maintain column integrity during multiple operations", () => {
      component.ngOnInit();

      // Toggle visibility
      component.toggleColumnVisibility(component.columns[0]);

      // Drop operation
      const mockEvent: CdkDragDrop<any[]> = {
        previousIndex: 0,
        currentIndex: 1,
        item: null as any,
        container: null as any,
        previousContainer: null as any,
        isPointerOverContainer: false,
        distance: { x: 0, y: 0 },
        dropPoint: { x: 0, y: 0 },
        event: null as any
      };
      component.drop(mockEvent);

      // Verify columns still have required properties
      component.columns.forEach((column, index) => {
        expect(column.matColumnDef).toBeDefined();
        expect(typeof column.visible).toBe("boolean");
        expect(column.order).toBe(index);
        expect(column.displayName).toBeDefined();
      });
    });

    it("should handle empty data gracefully", () => {
      mockDialogData.currentColumns = [];

      component.ngOnInit();

      expect(component.columns).toEqual([]);
      expect(component.visibleColumnsCount).toBe(0);
    });
  });
});
