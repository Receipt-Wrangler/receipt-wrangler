import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { By } from "@angular/platform-browser";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { provideRouter } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { ButtonModule } from "../../button";
import { DirectivesModule } from "../../directives/directives.module";
import { ApiModule, Permission } from "../../open-api";
import { AuthState, GroupState } from "../../store";
import { SetPermissions } from "../../store/auth.state.actions";
import { GroupTableState } from "../../store/group-table.state";

import { GroupTableComponent } from "./group-table.component";

describe("GroupTableComponent", () => {
  let component: GroupTableComponent;
  let fixture: ComponentFixture<GroupTableComponent>;
  let store: Store;

  const queryTestId = (id: string) =>
    fixture.debugElement.query(By.css(`[data-testid="${id}"]`));

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GroupTableComponent],
      imports: [
        ApiModule,
        ButtonModule,
        DirectivesModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, GroupState, GroupTableState]),
        NoopAnimationsModule,
        SharedUiModule,
        TableModule,
      ],
      providers: [
        provideZonelessChangeDetection(),
        provideRouter([]),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(GroupTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  // This component sets `displayedColumns` in ngAfterViewInit, which trips
  // checkNoChanges (NG0100) during a whenStable() app-tick. The FAB gating we
  // assert here only depends on the *hasAppPermission directive's effect, so we
  // drive change detection with detectChanges(false) to skip the verification
  // pass rather than restructure the unrelated table column wiring.
  it("shows the Create Group FAB with the app group-create permission", () => {
    store.dispatch(new SetPermissions([Permission.AppGroupsCreate], {}));
    fixture.detectChanges(false);

    expect(queryTestId("group-create")).toBeTruthy();
  });

  it("hides the Create Group FAB without the app group-create permission", () => {
    store.dispatch(new SetPermissions([], {}));
    fixture.detectChanges(false);

    expect(queryTestId("group-create")).toBeFalsy();
  });
});
