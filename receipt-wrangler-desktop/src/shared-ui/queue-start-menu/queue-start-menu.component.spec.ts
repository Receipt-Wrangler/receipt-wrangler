import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatMenuModule } from "@angular/material/menu";
import { GroupRolePipe } from "../../pipes/group-role.pipe";
import { QueueMode, ReceiptQueueService } from "../../services/receipt-queue.service";
import { StoreModule } from "../../store/store.module";

import { QueueStartMenuComponent } from "./queue-start-menu.component";

describe("QueueStartMenuComponent", () => {
  let component: QueueStartMenuComponent;
  let fixture: ComponentFixture<QueueStartMenuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [QueueStartMenuComponent, GroupRolePipe],
      imports: [
        StoreModule,
        MatMenuModule,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        GroupRolePipe,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(QueueStartMenuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init queue with receipt ids", () => {
    const receiptQueueServiceSpy = jest.spyOn(TestBed.inject(ReceiptQueueService), "initQueueAndNavigate");

    component.receiptIds = [1, 2, 3];
    component.initQueue(QueueMode.VIEW);

    expect(receiptQueueServiceSpy).toHaveBeenCalledWith(["1", "2", "3"], QueueMode.VIEW);
  });

  it("should init queue with receipt ids from receipts", () => {
    const receiptQueueServiceSpy = jest.spyOn(TestBed.inject(ReceiptQueueService), "initQueueAndNavigate");

    component.receipts = [
      { id: 1 } as any,
      { id: 2 } as any,
      { id: 3 } as any,
    ];
    component.initQueue(QueueMode.VIEW);

    expect(receiptQueueServiceSpy).toHaveBeenCalledWith(["1", "2", "3"], QueueMode.VIEW);
  });
});
