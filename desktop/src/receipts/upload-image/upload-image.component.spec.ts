import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { Subject } from "rxjs";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { ApiModule, ReceiptImageService } from "../../open-api";
import { UploadImageComponent } from "./upload-image.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

/**
 * jsdom's FileReader is asynchronous; this one finishes as soon as it is asked
 * to read, so a test controls exactly when each file is done. `failNextRead`
 * makes the next read fire onerror instead of onload.
 */
class FakeFileReader {
  public static failNextRead = false;

  public onload: (() => void) | null = null;

  public onerror: (() => void) | null = null;

  public readAsArrayBuffer(): void {
    if (FakeFileReader.failNextRead) {
      FakeFileReader.failNextRead = false;
      this.onerror?.();
      return;
    }
    this.onload?.();
  }
}

describe("UploadImageComponent", () => {
  let component: UploadImageComponent;
  let fixture: ComponentFixture<UploadImageComponent>;
  let receiptImageService: ReceiptImageService;
  let processing: boolean[];
  let loaded: ReceiptFileUploadCommand[];
  const originalFileReader = window.FileReader;

  const pick = (...files: File[]): void => {
    component.onFileChange({ target: { files } });
  };

  const pdf = (name = "receipt.pdf"): File =>
    new File(["%PDF"], name, { type: "application/pdf" });

  const png = (name = "receipt.png"): File =>
    new File(["png"], name, { type: "image/png" });

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule, MatSnackBarModule, UploadImageComponent],
    providers: [
        { provide: ActivatedRoute, useValue: { snapshot: { data: {} } } },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    fixture = TestBed.createComponent(UploadImageComponent);
    component = fixture.componentInstance;
    receiptImageService = TestBed.inject(ReceiptImageService);
    fixture.detectChanges();

    processing = [];
    loaded = [];
    component.processingChange.subscribe((value) => processing.push(value));
    component.fileLoaded.subscribe((command) => loaded.push(command));

    FakeFileReader.failNextRead = false;
    (window as any).FileReader = FakeFileReader;
  });

  afterEach(() => {
    window.FileReader = originalFileReader;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("stays processing while a PDF converts and stops once it is emitted", () => {
    const conversion = new Subject<any>();
    jest.spyOn(receiptImageService, "convertToJpg").mockReturnValue(conversion as any);

    pick(pdf());

    expect(processing).toEqual([true]);
    expect(loaded).toEqual([]);

    conversion.next({ encodedImage: "data:image/jpeg;base64,abc" });

    expect(loaded.length).toBe(1);
    expect(loaded[0].encodedImage).toBe("data:image/jpeg;base64,abc");
    expect(processing).toEqual([true, false]);
  });

  it("stops processing when the conversion fails", () => {
    const conversion = new Subject<any>();
    jest.spyOn(receiptImageService, "convertToJpg").mockReturnValue(conversion as any);

    pick(pdf());
    conversion.error(new Error("conversion failed"));

    expect(loaded).toEqual([]);
    expect(processing).toEqual([true, false]);
  });

  it("stops processing when a file cannot be read", () => {
    FakeFileReader.failNextRead = true;

    pick(png());

    expect(loaded).toEqual([]);
    expect(processing).toEqual([true, false]);
  });

  it("reports a plain image as processing until it is emitted", () => {
    const convertSpy = jest.spyOn(receiptImageService, "convertToJpg");

    pick(png());

    expect(convertSpy).not.toHaveBeenCalled();
    expect(loaded.length).toBe(1);
    expect(processing).toEqual([true, false]);
  });

  it("emits once each way for a batch, after the last file is done", () => {
    const first = new Subject<any>();
    const second = new Subject<any>();
    jest
      .spyOn(receiptImageService, "convertToJpg")
      .mockReturnValueOnce(first as any)
      .mockReturnValueOnce(second as any);

    pick(pdf("a.pdf"), png(), pdf("b.pdf"));

    expect(loaded.length).toBe(1);
    expect(processing).toEqual([true]);

    second.next({ encodedImage: "b" });
    expect(processing).toEqual([true]);

    first.next({ encodedImage: "a" });
    expect(loaded.length).toBe(3);
    expect(processing).toEqual([true, false]);
  });

  it("does not report processing for files it rejects", () => {
    pick(new File(["text"], "notes.txt", { type: "text/plain" }));

    expect(loaded).toEqual([]);
    expect(processing).toEqual([]);
  });
});
