import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { TableModule } from "src/table/table.module";
import { DirectivesModule } from "../../directives";
import { ApiModule, UserService } from "../../open-api";
import { AuthState } from "../../store";
import { UserTableState } from "../../store/user-table.state";
import { UserListComponent } from "./user-list.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

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

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should fetch a page from the backend and set datasource + total count", () => {
    const serviceSpy = jest.spyOn(TestBed.inject(UserService), "getPagedUsers");
    serviceSpy.mockReturnValue(
      of({
        data: [{ id: 1, username: "admin" }],
        totalCount: 1,
      } as any)
    );

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
});
