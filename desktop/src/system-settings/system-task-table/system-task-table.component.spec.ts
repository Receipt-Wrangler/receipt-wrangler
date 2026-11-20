import { firstValueFrom } from "rxjs";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormGroup } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { MatChipsModule } from "@angular/material/chips";
import { MatIconModule } from "@angular/material/icon";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { ButtonComponent } from "../../button/index";
import { FilterOperation } from "../../open-api";
import { UserState } from "../../store";
import { buildDefaultSystemTaskFilter, SystemTaskTableState } from "../../store/system-task-table.state";
import {
  SetPage,
  SetSystemTaskFilter,
  SetSystemTaskFilterField
} from "../../store/system-task-table.state.actions";

import { SystemTaskFilterComponent } from "../system-task-filter/system-task-filter.component";
import { SystemTaskTableComponent } from "./system-task-table.component";

describe("SystemTaskTableComponent", () => {
  let component: SystemTaskTableComponent;
  let fixture: ComponentFixture<SystemTaskTableComponent>;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [ButtonComponent,
        MatChipsModule,
        MatIconModule,
        NgxsModule.forRoot([SystemTaskTableState, UserState]),
        NoopAnimationsModule, SystemTaskTableComponent],
      providers: [{
        provide: ActivatedRoute,
        useValue: {
          snapshot: {
            data: {
              prompts: [],
              allReceiptProcessingSettings: []
            }
          }
        }
      }, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
    })
      .compileComponents();

    store = TestBed.inject(Store);

    fixture = TestBed.createComponent(SystemTaskTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("starts with no chips and an empty badge", () => {
    expect(component.filterChips()).toEqual([]);
    expect(component.numFiltersApplied()).toBe(0);
  });

  it("builds one chip per active condition off the stored filter", () => {
    store.dispatch(new SetSystemTaskFilter({
      ...buildDefaultSystemTaskFilter(),
      type: { operation: FilterOperation.Contains, value: ["QUICK_SCAN"] },
    } as any));

    expect(component.filterChips().map((chip) => chip.key)).toEqual(["type"]);
    expect(component.numFiltersApplied()).toBe(1);
  });

  // Clearing a chip from page 7 must not leave the user on a page the narrowed
  // result set no longer has.
  it("clears a chip and resets to page one", () => {
    store.dispatch(new SetPage(7));
    store.dispatch(new SetSystemTaskFilter({
      ...buildDefaultSystemTaskFilter(),
      type: { operation: FilterOperation.Contains, value: ["QUICK_SCAN"] },
    } as any));

    const dispatchSpy = jest.spyOn(store, "dispatch");
    jest.spyOn(component, "refresh").mockImplementation(() => {});

    component.filterChipCleared("type");

    expect(dispatchSpy).toHaveBeenCalledWith(new SetSystemTaskFilterField("type", null));
    expect(dispatchSpy).toHaveBeenCalledWith(new SetPage(1));
    expect(store.selectSnapshot(SystemTaskTableState.page)).toBe(1);
    expect(component.filterChips()).toEqual([]);
  });

  // #750: the filter form is built per open, so its subscriptions must end with
  // the dialog, not live on until the table is destroyed.
  it("ends the filter form's subscriptions when the dialog closes", async () => {
    component.filterButtonClicked();
    const dialogRef = TestBed.inject(MatDialog).openDialogs[0];
    const form: FormGroup = (dialogRef.componentInstance as SystemTaskFilterComponent).parentForm();
    const watched = [form.get("type.value")!, form.get("startedAt.operation")!];

    expect(watched.every((control) => (control.valueChanges as any).observed)).toBe(true);

    const closed = firstValueFrom(dialogRef.afterClosed());
    dialogRef.close();
    await closed;

    expect(dialogRef.componentRef?.hostView.destroyed).toBe(true);
    expect(watched.some((control) => (control.valueChanges as any).observed)).toBe(false);
  });
});
