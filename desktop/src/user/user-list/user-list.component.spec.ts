import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { PageEvent } from "@angular/material/paginator";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { NgxsModule } from "@ngxs/store";
import { config, of, throwError } from "rxjs";
import { TableModule } from "src/table/table.module";
import { DirectivesModule } from "../../directives";
import { ApiModule, User, UserService } from "../../open-api";
import { AuthState } from "../../store";
import { UserTableState } from "../../store/user-table.state";
import { UserListComponent } from "./user-list.component";

describe("UserListComponent", () => {
  let component: UserListComponent;
  let fixture: ComponentFixture<UserListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [UserListComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        DirectivesModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, UserTableState]),
        TableModule],
    providers: [provideZonelessChangeDetection(), provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
}).compileComponents();

    fixture = TestBed.createComponent(UserListComponent);
    component = fixture.componentInstance;
    // fixture.detectChanges();
  });

  const mockPagedUsers = () =>
    jest.spyOn(TestBed.inject(UserService), "getPagedUsers").mockReturnValue(
      of({ data: [{ id: 1, username: "admin" }], totalCount: 1 } as any)
    );

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should fetch a page from the backend and set datasource + total count", () => {
    const serviceSpy = mockPagedUsers();

    component.getTableData();

    expect(serviceSpy).toHaveBeenCalledWith({
      page: 1,
      pageSize: 50,
      orderBy: "username",
      sortDirection: "asc",
    });
    expect(component.totalCount()).toBe(1);
    expect(component.dataSource().data.length).toBe(1);
  });

  it("should refetch with the new 1-based page and page size on page change", () => {
    const serviceSpy = mockPagedUsers();

    component.pageChanged({ pageIndex: 2, pageSize: 15 } as PageEvent);

    expect(serviceSpy).toHaveBeenCalledWith({
      page: 3, // pageIndex + 1
      pageSize: 15,
      orderBy: "username",
      sortDirection: "asc",
    });
  });

  it("should refetch with the requested column and direction on sort", () => {
    const serviceSpy = mockPagedUsers();

    component.sorted({ active: "display_name", direction: "desc" } as Sort);

    expect(serviceSpy).toHaveBeenCalledWith(
      expect.objectContaining({ orderBy: "display_name", sortDirection: "desc" })
    );
  });

  it("should refetch with a cleared sort direction", () => {
    const serviceSpy = mockPagedUsers();

    component.sorted({ active: "username", direction: "" } as Sort);

    expect(serviceSpy).toHaveBeenCalledWith(
      expect.objectContaining({ orderBy: "username", sortDirection: "" })
    );
  });

  it("should refetch the current page after deleting a user", () => {
    const serviceSpy = mockPagedUsers();
    const deleteSpy = jest
      .spyOn(TestBed.inject(UserService), "deleteUserById")
      .mockReturnValue(of({} as any));

    // A user other than the (undefined) current user, so the delete guard passes.
    component.dataSource.set(
      new MatTableDataSource<User>([{ id: 5, username: "victim" } as User])
    );

    jest
      .spyOn(TestBed.inject(MatDialog), "open")
      .mockReturnValue({ afterClosed: () => of(true), componentInstance: {} } as any);

    component.deleteUser(0);

    expect(deleteSpy).toHaveBeenCalledWith(5);
    expect(serviceSpy).toHaveBeenCalled(); // refetch
  });

  it("should not corrupt datasource/total count when the fetch errors", () => {
    // getTableData delegates errors to the global HTTP interceptor (no component
    // handler), so RxJS reports the error asynchronously. Swallow that rethrow with
    // fake timers + a no-op onUnhandledError so the assertion stays deterministic.
    jest.useFakeTimers();
    const previousHandler = config.onUnhandledError;
    config.onUnhandledError = () => {};

    try {
      const serviceSpy = jest
        .spyOn(TestBed.inject(UserService), "getPagedUsers")
        .mockReturnValue(throwError(() => new Error("boom")) as any);

      component.getTableData();

      expect(serviceSpy).toHaveBeenCalled();
      expect(component.totalCount()).toBe(0);
      expect(component.dataSource().data.length).toBe(0);

      jest.runOnlyPendingTimers();
    } finally {
      config.onUnhandledError = previousHandler;
      jest.useRealTimers();
    }
  });
});
