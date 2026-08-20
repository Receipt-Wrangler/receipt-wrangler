import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTableDataSource } from "@angular/material/table";
import { By } from "@angular/platform-browser";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { ButtonModule } from "../../button";
import { DirectivesModule } from "../../directives/directives.module";
import { FormMode } from "../../enums/form-mode.enum";
import { ApiModule, CustomField, PagedDataDataInner, Permission } from "../../open-api";
import { PipesModule } from "../../pipes";
import { AuthState } from "../../store";
import { SetPermissions } from "../../store/auth.state.actions";
import { CustomFieldTableState } from "../../store/custom-field-table.state";
import { TableModule } from "../../table/table.module";

import { CustomFieldTableComponent } from "./custom-field-table.component";

describe("CustomFieldTableComponent", () => {
  let component: CustomFieldTableComponent;
  let fixture: ComponentFixture<CustomFieldTableComponent>;
  let store: Store;

  const customField = {
    id: 3,
    name: "Cost Centre",
    type: "TEXT",
    description: "Which cost centre paid",
  } as unknown as CustomField;

  const queryTestId = (id: string) =>
    fixture.debugElement.query(By.css(`[data-testid="${id}"]`));

  // Row actions only exist once the table has columns AND a row. Permissions are
  // dispatched BEFORE the first change-detection pass because setColumns() reads
  // them with a one-shot selectSnapshot from ngAfterViewInit -- which matches the
  // real app, where the route guard has already loaded AppData before the table
  // renders.
  const renderRowWith = (appPermissions: string[]) => {
    store.dispatch(new SetPermissions(appPermissions, {}));
    fixture.detectChanges(false);
    component.dataSource.set(
      new MatTableDataSource([customField as unknown as PagedDataDataInner])
    );
    fixture.detectChanges(false);
    TestBed.flushEffects();
    fixture.detectChanges(false);
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CustomFieldTableComponent],
      imports: [
        ApiModule,
        ButtonModule,
        DirectivesModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, CustomFieldTableState]),
        NoopAnimationsModule,
        PipesModule,
        SharedUiModule,
        TableModule,
      ],
      providers: [
        provideZonelessChangeDetection(),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
        {
          provide: ActivatedRoute,
          useValue: {}
        }
      ]
    })
      .compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(CustomFieldTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  describe("row actions", () => {
    it("shows the edit action to a holder of the update permission", () => {
      renderRowWith([Permission.AppCustomFieldsUpdate]);

      expect(queryTestId("custom-field-edit")).toBeTruthy();
    });

    it("hides the edit action from a caller who can only read custom fields", () => {
      renderRowWith([Permission.AppCustomFieldsRead]);

      expect(queryTestId("custom-field-edit")).toBeFalsy();
    });

    // The actions column used to be gated on delete alone, which would have
    // hidden the whole column -- and with it the new edit button -- from a
    // caller who may update but not delete.
    it("renders the actions column for an update-only holder", () => {
      renderRowWith([Permission.AppCustomFieldsUpdate]);

      expect(component.displayedColumns).toContain("actions");
    });

    it("renders the actions column for a delete-only holder", () => {
      renderRowWith([Permission.AppCustomFieldsDelete]);

      expect(component.displayedColumns).toContain("actions");
    });

    it("omits the actions column when the caller can neither update nor delete", () => {
      renderRowWith([Permission.AppCustomFieldsRead]);

      expect(component.displayedColumns).not.toContain("actions");
    });
  });

  describe("openCustomFieldDialog", () => {
    const openDialogAndReadInstance = (
      appPermissions: string[],
      field?: CustomField,
      mode?: FormMode
    ) => {
      const componentInstance: any = {};
      const matDialog = TestBed.inject(MatDialog);
      jest.spyOn(matDialog, "open").mockReturnValue({
        componentInstance,
        afterClosed: () => of(false),
      } as any);

      store.dispatch(new SetPermissions(appPermissions, {}));
      component.openCustomFieldDialog(field, mode);

      return componentInstance;
    };

    it("opens in edit mode from the name link for an update holder", () => {
      const instance = openDialogAndReadInstance(
        [Permission.AppCustomFieldsUpdate],
        customField
      );

      expect(instance.mode).toBe(FormMode.edit);
      expect(instance.headerText).toBe("Edit Cost Centre");
      expect(instance.customField).toBe(customField);
    });

    it("opens read-only from the name link without the update permission", () => {
      const instance = openDialogAndReadInstance(
        [Permission.AppCustomFieldsRead],
        customField
      );

      expect(instance.mode).toBe(FormMode.view);
      expect(instance.headerText).toBe("View Custom Field");
    });

    it("opens in add mode when no custom field is passed", () => {
      const instance = openDialogAndReadInstance([Permission.AppCustomFieldsCreate]);

      expect(instance.mode).toBe(FormMode.add);
      expect(instance.headerText).toBe("Add Custom Field");
    });

    it("honors an explicitly requested mode", () => {
      const instance = openDialogAndReadInstance(
        [Permission.AppCustomFieldsUpdate],
        customField,
        FormMode.edit
      );

      expect(instance.mode).toBe(FormMode.edit);
    });
  });
});
