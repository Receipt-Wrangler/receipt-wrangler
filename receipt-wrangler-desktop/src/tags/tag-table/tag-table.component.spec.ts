import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "src/constants";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { TagTableState } from "src/store/tag-table.state";
import { ApiModule, TagService } from "../../open-api";
import { TagFormComponent } from "../tag-form/tag-form.component";
import { TagTableComponent } from "./tag-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("TagsListComponent", () => {
  let component: TagTableComponent;
  let fixture: ComponentFixture<TagTableComponent>;
  let store: Store;
  let tagService: jest.Mocked<TagService>;

  beforeEach(() => {
    const tagServiceSpy = {
      getPagedTags: jest.fn(),
      deleteTag: jest.fn()
    };
    tagServiceSpy.getPagedTags.mockReturnValue(of({
      data: [],
      totalCount: 0
    } as any));

    TestBed.configureTestingModule({
    declarations: [TagTableComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        MatDialogModule,
        NgxsModule.forRoot([TagTableState]),
        MatSnackBarModule],
    providers: [
        { provide: TagService, useValue: tagServiceSpy },
        provideHttpClient(withInterceptorsFromDi()), 
        provideHttpClientTesting()
    ]
});
    fixture = TestBed.createComponent(TagTableComponent);
    store = TestBed.inject(Store);
    tagService = TestBed.inject(TagService) as jest.Mocked<TagService>;
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should attempt to get table data, set datasource and total count", () => {
    tagService.getPagedTags.mockReturnValue(
      of({
        data: [{}],
        totalCount: 1,
      } as any)
    );

    component.ngOnInit();

    expect(tagService.getPagedTags).toHaveBeenCalledWith({
      page: 1,
      pageSize: 50,
      orderBy: "name",
      sortDirection: "desc",
    });

    expect(component.totalCount).toEqual(1);
    expect(component.dataSource.data).toEqual([{} as any]);
  });

  it("should attempt to get table data, with new sorted direction and key", () => {
    tagService.getPagedTags.mockReturnValue(
      of({
        data: [{}],
        totalCount: 1,
      } as any)
    );

    component.sorted({
      active: "numberOfReceipts",
      direction: "asc",
    });

    expect(store.selectSnapshot(TagTableState.state)).toEqual({
      page: 1,
      pageSize: 50,
      orderBy: "numberOfReceipts",
      sortDirection: "asc",
    });
    expect(tagService.getPagedTags).toHaveBeenCalledWith({
      page: 1,
      pageSize: 50,
      orderBy: "numberOfReceipts",
      sortDirection: "asc",
    });
  });

  it("should attempt to get table data, with newpage and new page size", () => {
    tagService.getPagedTags.mockReturnValue(
      of({
        data: [{}],
        totalCount: 1,
      } as any)
    );

    component.updatePageData({
      pageIndex: 2,
      pageSize: 100,
    } as any);

    expect(store.selectSnapshot(TagTableState.state)).toEqual({
      page: 3,
      pageSize: 100,
      orderBy: "name",
      sortDirection: "desc",
    });
    expect(tagService.getPagedTags).toHaveBeenCalledWith({
      page: 3,
      pageSize: 100,
      orderBy: "name",
      sortDirection: "desc",
    });
  });

  it("should set columns", () => {
    component.ngAfterViewInit();

    expect(component.columns.length).toEqual(4);
    expect(component.displayedColumns).toEqual([
      "name",
      "description",
      "numberOfReceipts",
      "actions",
    ]);
  });

  it("should open edit dialog and refresh data when after closed with true", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open");
    dialogSpy.mockReturnValue({
      componentInstance: {
        tag: {},
        headerText: "",
      },
      afterClosed: () => of(true),
    } as any);

    const tagView: any = {};
    component.openEditDialog(tagView);

    expect(dialogSpy).toHaveBeenCalledWith(
      TagFormComponent,
      DEFAULT_DIALOG_CONFIG
    );
    expect(tagService.getPagedTags).toHaveBeenCalledTimes(1);
  });

  it("should open edit dialog and not refresh data when after closed with false", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open");
    dialogSpy.mockReturnValue({
      componentInstance: {
        tag: {},
        headerText: "",
      },
      afterClosed: () => of(false),
    } as any);

    const tagView: any = {};
    component.openEditDialog(tagView);

    expect(dialogSpy).toHaveBeenCalledWith(
      TagFormComponent,
      DEFAULT_DIALOG_CONFIG
    );
    expect(tagService.getPagedTags).toHaveBeenCalledTimes(0);
  });

  it("should open confirmation dialog and refresh data when after closed with true", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open");
    tagService.deleteTag.mockReturnValue(of(undefined as any));
    dialogSpy.mockReturnValue({
      componentInstance: {
        tag: {},
        headerText: "",
      },
      afterClosed: () => of(true),
    } as any);

    const tagView: any = { id: 1 };
    component.openDeleteConfirmationDialog(tagView);

    expect(dialogSpy).toHaveBeenCalledWith(
      ConfirmationDialogComponent,
      DEFAULT_DIALOG_CONFIG
    );
    expect(tagService.deleteTag).toHaveBeenCalledWith(1);
  });

  it("should open confirmation dialog and not refresh data when after closed with false", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open");
    dialogSpy.mockReturnValue({
      componentInstance: {
        tag: {},
        headerText: "",
      },
      afterClosed: () => of(false),
    } as any);

    const tagView: any = {};
    component.openDeleteConfirmationDialog(tagView);

    expect(dialogSpy).toHaveBeenCalledWith(
      ConfirmationDialogComponent,
      DEFAULT_DIALOG_CONFIG
    );
    expect(tagService.getPagedTags).toHaveBeenCalledTimes(0);
  });

  it("should open add dialog and refresh data when after closed with true", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open");
    dialogSpy.mockReturnValue({
      componentInstance: {
        tag: {},
        headerText: "",
      },
      afterClosed: () => of(true),
    } as any);

    component.openAddDialog();

    expect(dialogSpy).toHaveBeenCalledWith(
      TagFormComponent,
      DEFAULT_DIALOG_CONFIG
    );
    expect(tagService.getPagedTags).toHaveBeenCalledTimes(1);
  });

  it("should open add dialog and not refresh data when after closed with false", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open");
    dialogSpy.mockReturnValue({
      componentInstance: {
        tag: {},
        headerText: "",
      },
      afterClosed: () => of(false),
    } as any);

    component.openAddDialog();

    expect(dialogSpy).toHaveBeenCalledWith(
      TagFormComponent,
      DEFAULT_DIALOG_CONFIG
    );
    expect(tagService.getPagedTags).toHaveBeenCalledTimes(0);
  });
});
