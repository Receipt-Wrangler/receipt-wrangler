import { Component, input, output, viewChild } from "@angular/core";
import { finalize, take, tap } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { ReceiptImageService } from "../../open-api";

@Component({
    selector: "app-upload-image",
    templateUrl: "./upload-image.component.html",
    styleUrls: ["./upload-image.component.scss"],
    standalone: false
})
export class UploadImageComponent {
  readonly uploadInput = viewChild.required<any>("uploadInput");

  public readonly receiptId = input<string | undefined>("");

  public readonly multiple = input<boolean>(true);

  public readonly acceptFileTypes = input<string[]>([
    "image/*",
    "application/pdf",
    "image/heic",
]);

  public readonly fileLoaded = output<ReceiptFileUploadCommand>();

  /**
   * Emits true when a picked file starts being read or converted, and false once
   * every one of them has been emitted or has failed. An output rather than a
   * public signal on purpose: hosts' specs render this element under
   * CUSTOM_ELEMENTS_SCHEMA, where a template-ref read of a signal would throw.
   */
  public readonly processingChange = output<boolean>();

  public formMode = FormMode;

  private pendingFileCount = 0;
  
  constructor(private receiptImageService: ReceiptImageService) {}

  public clickInput(): void {
    this.uploadInput().nativeElement.click();
  }

  public async onFileChange(event: any): Promise<void> {
    const files: File[] = Array.from(event?.target?.files ?? []);
    const acceptedFiles = files.filter((f) =>
      this.acceptFileTypes().some((t) => new RegExp(t).test(f.type))
    );

    this.adjustPending(acceptedFiles.length);

    for (let i = 0; i < acceptedFiles.length; i++) {
      const reader = new FileReader();
      const f = acceptedFiles[i];

      reader.onload = () => {
        const command: ReceiptFileUploadCommand = {
          file: f,
          receiptId: Number(this.receiptId()),
        };

        if (f.type === "application/pdf" || f.type === "image/heic") {
          this.receiptImageService
            .convertToJpg(f)
            .pipe(
              take(1),
              tap((encodedImage) => {
                command.encodedImage = encodedImage.encodedImage;
                this.fileLoaded.emit(command);
              }),
              finalize(() => this.adjustPending(-1))
            )
            .subscribe();
        } else {
          this.fileLoaded.emit(command);
          this.adjustPending(-1);
        }
      };
      reader.onerror = () => this.adjustPending(-1);

      reader.readAsArrayBuffer(f);
    }
  }

  private adjustPending(delta: number): void {
    const wasProcessing = this.pendingFileCount > 0;
    this.pendingFileCount = Math.max(0, this.pendingFileCount + delta);
    const isProcessing = this.pendingFileCount > 0;

    if (isProcessing !== wasProcessing) {
      this.processingChange.emit(isProcessing);
    }
  }
}
