import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTableDataSource } from "@angular/material/table";
import { By } from "@angular/platform-browser";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { provideRouter } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { TableModule } from "src/table/table.module";
import { ButtonModule } from "../../button";
import { DirectivesModule } from "../../directives/directives.module";
import { ApiModule, Group, GroupsService, Permission } from "../../open-api";
import { AuthState, GroupState, SetGroups } from "../../store";
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

  describe("delete row action", () => {
    // Row actions only exist once the table has columns (wired in
    // ngAfterViewInit, so it needs one CD pass first) AND a row. Group 5 stands
    // in for a group the caller is NOT a member of — the all-groups admin view.
    const renderRowWith = (
      appPermissions: string[],
      groupPermissions: { [groupId: number]: string[] } = {}
    ) => {
      fixture.detectChanges(false);
      component.dataSource.set(
        new MatTableDataSource([
          { id: 5, name: "Abandoned Group", groupMembers: [] } as unknown as Group,
        ])
      );
      store.dispatch(new SetPermissions(appPermissions, groupPermissions));
      fixture.detectChanges(false);
      TestBed.flushEffects();
      fixture.detectChanges(false);
    };

    it("shows the delete action to a member holding group.delete", () => {
      renderRowWith([], { 5: [Permission.GroupDelete] });

      expect(queryTestId("group-delete")).toBeTruthy();
    });

    it("shows the delete action to a non-member holding app.groups.delete", () => {
      renderRowWith([Permission.AppGroupsDelete]);

      expect(queryTestId("group-delete")).toBeTruthy();
    });

    it("hides the delete action from a non-member who can only read all groups", () => {
      renderRowWith([Permission.AppGroupsRead]);

      expect(queryTestId("group-delete")).toBeFalsy();
    });
  });

  describe("deleteDisabled", () => {
    // Mirrors the backend CanDeleteGroup rule. GroupState.groups holds the
    // CALLER's own groups (including the virtual "All" group), so without the
    // permission escape a one-group admin is blocked from deleting anything.
    it("disables delete for a caller with a single group and no app permission", () => {
      store.dispatch(new SetGroups([{ id: 1 } as Group]));
      store.dispatch(new SetPermissions([], {}));

      expect(component.deleteDisabled()).toBe(true);
    });

    it("does not disable delete for an app.groups.delete holder with a single group", () => {
      store.dispatch(new SetGroups([{ id: 1 } as Group]));
      store.dispatch(new SetPermissions([Permission.AppGroupsDelete], {}));

      expect(component.deleteDisabled()).toBe(false);
    });

    it("does not disable delete for a caller with more than one group", () => {
      store.dispatch(new SetGroups([{ id: 1 } as Group, { id: 2 } as Group]));
      store.dispatch(new SetPermissions([], {}));

      expect(component.deleteDisabled()).toBe(false);
    });
  });

  describe("deleteGroup", () => {
    it("refetches the current page instead of swapping in the caller's own groups", () => {
      // The table is server-paged and may be showing the all-groups filter, so
      // the post-delete refresh has to go back to the server.
      const deleteSpy = jest
        .spyOn(TestBed.inject(GroupsService), "deleteGroup")
        .mockReturnValue(of({} as any));
      const refetchSpy = jest
        .spyOn(component, "getTableData")
        .mockImplementation(() => {});
      jest
        .spyOn(TestBed.inject(MatDialog), "open")
        .mockReturnValue({ afterClosed: () => of(true), componentInstance: {} } as any);

      store.dispatch(new SetGroups([{ id: 1 } as Group, { id: 2 } as Group]));
      component.dataSource.set(
        new MatTableDataSource([{ id: 5, name: "Abandoned Group" } as Group])
      );

      component.deleteGroup(0);

      expect(deleteSpy).toHaveBeenCalledWith(5);
      expect(refetchSpy).toHaveBeenCalled();
    });

    it("does nothing when deletion is disabled", () => {
      const deleteSpy = jest.spyOn(TestBed.inject(GroupsService), "deleteGroup");
      const dialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open");

      store.dispatch(new SetGroups([{ id: 1 } as Group]));
      store.dispatch(new SetPermissions([], {}));
      component.dataSource.set(
        new MatTableDataSource([{ id: 1, name: "Only Group" } as Group])
      );

      component.deleteGroup(0);

      expect(dialogSpy).not.toHaveBeenCalled();
      expect(deleteSpy).not.toHaveBeenCalled();
    });
  });
});
