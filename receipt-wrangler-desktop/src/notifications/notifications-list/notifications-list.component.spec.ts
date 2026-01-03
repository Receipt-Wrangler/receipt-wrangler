import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { of } from "rxjs";
import { ApiModule, Notification, NotificationsService } from "../../open-api";
import { NotificationsListComponent } from "./notifications-list.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("NotificationsListComponent", () => {
  let component: NotificationsListComponent;
  let fixture: ComponentFixture<NotificationsListComponent>;
  let service: NotificationsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
    declarations: [NotificationsListComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});

    service = TestBed.inject(NotificationsService);
    fixture = TestBed.createComponent(NotificationsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should load notifications on init", () => {
    const mockNotifications: Notification[] = [
      {
        id: 1,
        body: "Body 1",
        userId: 10,
        type: "Type 1",
        title: "Title 1",
        createdAt: "2023-07-03",
      },
      {
        id: 2,
        body: "Body 2",
        userId: 20,
        type: "Type 2",
        title: "Title 2",
        createdAt: "2023-07-04",
      },
    ];

    jest.spyOn(service, "getNotificationsForuser").mockReturnValue(
      of(mockNotifications as any)
    );

    component.ngOnInit();

    expect(component.notifications).toEqual(mockNotifications);
    expect(service.getNotificationsForuser).toHaveBeenCalled();
  });

  it("should delete all notifications", () => {
    jest.spyOn(service, "deleteAllNotificationsForUser").mockReturnValue(
      of(undefined as any)
    );
    component.notifications = [
      {
        id: 1,
        body: "Body 1",
        userId: 10,
        type: "Type 1",
        title: "Title 1",
        createdAt: "2023-07-03",
      },
      {
        id: 2,
        body: "Body 2",
        userId: 20,
        type: "Type 2",
        title: "Title 2",
        createdAt: "2023-07-04",
      },
    ];

    component.deleteAllNotifications();

    expect(component.notifications.length).toEqual(0);
    expect(service.deleteAllNotificationsForUser).toHaveBeenCalled();
  });

  it("should delete a notification by id", () => {
    component.notifications = [
      {
        id: 1,
        body: "Body 1",
        userId: 10,
        type: "Type 1",
        title: "Title 1",
        createdAt: "2023-07-03",
      },
      {
        id: 2,
        body: "Body 2",
        userId: 20,
        type: "Type 2",
        title: "Title 2",
        createdAt: "2023-07-04",
      },
    ];

    component.notificationDeleted(1);

    expect(component.notifications).toEqual([
      {
        id: 2,
        body: "Body 2",
        userId: 20,
        type: "Type 2",
        title: "Title 2",
        createdAt: "2023-07-04",
      },
    ]);
  });
});
