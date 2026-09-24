import { ScrollingModule } from "@angular/cdk/scrolling";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog } from "@angular/material/dialog";
import { MatListModule } from "@angular/material/list";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { DirectivesModule } from "../../directives/directives.module";
import { Permission, SystemTaskService, SystemTaskStatus, SystemTaskType } from "../../open-api";
import { PipesModule } from "../../pipes/pipes.module";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SystemTaskTypePipe } from "../../shared-ui/task-table/system-task-type.pipe";
import { AuthState } from "../../store/auth.state";
import { SetPermissions } from "../../store/auth.state.actions";
import { GroupState } from "../../store/group.state";
import { downloadFile } from "../../utils/file";
import { DashboardListComponent } from "../dashboard-list/dashboard-list.component";

import { ActivityComponent } from "./activity.component";

jest.mock("../../utils/file", () => ({
  ...jest.requireActual("../../utils/file"),
  downloadFile: jest.fn(),
}));

// The widget renders activities for one group, but on the "All" dashboard that
// group is the synthetic aggregate. ACTIVITY_GROUP_ID is the group an activity
// actually belongs to and is what every gate must key on.
const ACTIVITY_GROUP_ID = 7;
const WIDGET_GROUP_ID = 99;

describe("ActivityComponent", () => {
  let component: ActivityComponent;
  let fixture: ComponentFixture<ActivityComponent>;
  let store: Store;
  let systemTaskService: SystemTaskService;

  function buildActivity(overrides: Record<string, unknown> = {}): any {
    return {
      id: 1,
      type: SystemTaskType.QuickScan,
      status: SystemTaskStatus.Failed,
      startedAt: new Date().toISOString(),
      endedAt: new Date().toISOString(),
      groupId: ACTIVITY_GROUP_ID,
      canBeRestarted: false,
      hasSourceFile: false,
      ...overrides,
    };
  }

  async function renderActivities(activities: any[]): Promise<void> {
    component.activities.set(activities);
    fixture.detectChanges();
    await fixture.whenStable();
  }

  function grantIn(groupId: number, ...permissions: Permission[]): void {
    store.dispatch(new SetPermissions([], { [groupId]: permissions }));
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ActivityComponent, DashboardListComponent],
      imports: [
        SharedUiModule,
        ScrollingModule,
        MatListModule,
        DirectivesModule,
        PipesModule,
        // Standalone, and not re-exported by SharedUiModule, so the spec imports it
        // the same way DashboardModule does.
        SystemTaskTypePipe,
        NgxsModule.forRoot([AuthState, GroupState]),
      ],
      providers: [
        provideZonelessChangeDetection(),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ActivityComponent);
    component = fixture.componentInstance;
    store = TestBed.inject(Store);
    systemTaskService = TestBed.inject(SystemTaskService);

    jest.spyOn(systemTaskService, "getPagedActivities").mockReturnValue(of({ data: [], totalCount: 0 } as any));

    fixture.componentRef.setInput("widget", { name: "" } as any);
    fixture.componentRef.setInput("groupId", WIDGET_GROUP_ID);
    fixture.detectChanges();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("offers no source file controls when the activity has none", async () => {
    grantIn(ACTIVITY_GROUP_ID, Permission.GroupActivitiesRead);
    await renderActivities([buildActivity({ hasSourceFile: false })]);

    expect(fixture.nativeElement.querySelector("[data-testid='activity-source-file-preview']")).toBeNull();
    expect(fixture.nativeElement.querySelector("[data-testid='activity-source-file-download']")).toBeNull();
  });

  it("offers preview and download when the activity still has its upload", async () => {
    grantIn(ACTIVITY_GROUP_ID, Permission.GroupActivitiesRead);
    await renderActivities([buildActivity({ hasSourceFile: true })]);

    expect(fixture.nativeElement.querySelector("[data-testid='activity-source-file-preview']")).not.toBeNull();
    expect(fixture.nativeElement.querySelector("[data-testid='activity-source-file-download']")).not.toBeNull();
  });

  it("hides the source file controls without the group permission", async () => {
    grantIn(ACTIVITY_GROUP_ID);
    await renderActivities([buildActivity({ hasSourceFile: true })]);

    expect(fixture.nativeElement.querySelector("[data-testid='activity-source-file-preview']")).toBeNull();
  });

  // The regression this guards: gating on the widget's own groupId() resolves to
  // the synthetic "All" group on the all-groups dashboard, which holds no
  // permissions, so every control silently disappears for activities the user is
  // entitled to act on. Mobile fixed this long ago; desktop had not.
  it("gates on the activity's own group, not the widget's", async () => {
    grantIn(ACTIVITY_GROUP_ID, Permission.GroupActivitiesRead, Permission.GroupActivitiesRerun);
    await renderActivities([buildActivity({ hasSourceFile: true, canBeRestarted: true })]);

    expect(fixture.nativeElement.querySelector("[data-testid='activity-source-file-preview']")).not.toBeNull();
    expect(fixture.nativeElement.querySelector("[data-testid='activity-rerun']")).not.toBeNull();
  });

  it("opens the preview dialog with the returned image", () => {
    const dialog = TestBed.inject(MatDialog);
    const openSpy = jest.spyOn(dialog, "open").mockReturnValue({} as any);
    jest.spyOn(systemTaskService, "getSystemTaskSourceFile").mockReturnValue(
      of({ name: "receipt.png", encodedImage: "data:image/png;base64,AAA" } as any)
    );

    component.onPreviewButtonClick(1);

    expect(openSpy).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({
        data: { name: "receipt.png", encodedImage: "data:image/png;base64,AAA" },
      })
    );
  });

  it("downloads under the name the server sent", () => {
    const blob = new Blob(["bytes"]);
    jest.spyOn(systemTaskService, "downloadSystemTaskSourceFile").mockReturnValue(
      of({
        body: blob,
        headers: { get: () => 'attachment; filename="receipt.png"' },
      } as any)
    );

    component.onDownloadButtonClick(1);

    expect(downloadFile).toHaveBeenCalledWith(blob, "receipt.png");
  });
});
