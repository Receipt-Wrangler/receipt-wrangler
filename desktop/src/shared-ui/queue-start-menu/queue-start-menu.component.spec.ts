import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatMenuModule } from "@angular/material/menu";
import { QueueMode, ReceiptQueueService } from "../../services/receipt-queue.service";
import { provideTestStore } from "../../store/test-store";

import { QueueStartMenuComponent } from "./queue-start-menu.component";
import { provideRouter } from "@angular/router";

describe("QueueStartMenuComponent", () => {
  let component: QueueStartMenuComponent;
  let fixture: ComponentFixture<QueueStartMenuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        MatMenuModule,
        QueueStartMenuComponent,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        provideTestStore(),
        provideRouter([]),
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
    const receiptQueueServiceSpy = jest.spyOn(TestBed.inject(ReceiptQueueService), "initQueueAndNavigate").mockImplementation(() => {});

    fixture.componentRef.setInput('receiptIds', [1, 2, 3]);
    component.initQueue(QueueMode.VIEW);

    expect(receiptQueueServiceSpy).toHaveBeenCalledWith(["1", "2", "3"], QueueMode.VIEW);
  });

  it("should init queue with receipt ids from receipts", () => {
    const receiptQueueServiceSpy = jest.spyOn(TestBed.inject(ReceiptQueueService), "initQueueAndNavigate").mockImplementation(() => {});

    fixture.componentRef.setInput('receipts', [
      { id: 1 } as any,
      { id: 2 } as any,
      { id: 3 } as any,
    ]);
    component.initQueue(QueueMode.VIEW);

    expect(receiptQueueServiceSpy).toHaveBeenCalledWith(["1", "2", "3"], QueueMode.VIEW);
  });
});
