import { TestBed } from "@angular/core/testing";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { QueueMode, ReceiptQueueService } from "./receipt-queue.service";

describe("ReceiptQueueService", () => {
  let service: ReceiptQueueService;
  let router: Router;
  let navigateSpy: jest.SpyInstance;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [RouterTestingModule],
      providers: [ReceiptQueueService]
    });
    service = TestBed.inject(ReceiptQueueService);
    router = TestBed.inject(Router);
    navigateSpy = jest.spyOn(router, "navigate");
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });

  it("should navigate to the correct URL in initQueueAndNavigate", () => {
    const receiptIds = ["id1", "id2", "id3"];
    const mode = QueueMode.VIEW;
    const indexToStartAt = 1;

    service.initQueueAndNavigate(receiptIds, mode, indexToStartAt);

    expect(navigateSpy).toHaveBeenCalledWith(
      ["receipts", "id2", mode],
      {
        queryParams: {
          ids: receiptIds,
          queueMode: mode,
        },
        onSameUrlNavigation: "reload"
      }
    );
  });

  it("should navigate to the next item in queueNext", () => {
    const receiptIds = ["id1", "id2", "id3"];
    const mode = QueueMode.EDIT;
    const currentIndex = 1;

    service.queueNext(currentIndex, receiptIds, mode);

    expect(navigateSpy).toHaveBeenCalledWith(
      ["receipts", "id3", mode],
      {
        queryParams: {
          ids: receiptIds,
          queueMode: mode,
        },
        onSameUrlNavigation: "reload"
      }
    );
  });

  it("should navigate to the previous item in queuePrevious", () => {
    const receiptIds = ["id1", "id2", "id3"];
    const mode = QueueMode.VIEW;
    const currentIndex = 1;

    service.queuePrevious(currentIndex, receiptIds, mode);

    expect(navigateSpy).toHaveBeenCalledWith(
      ["receipts", "id1", mode],
      {
        queryParams: {
          ids: receiptIds,
          queueMode: mode,
        },
        onSameUrlNavigation: "reload"
      }
    );
  });
});
