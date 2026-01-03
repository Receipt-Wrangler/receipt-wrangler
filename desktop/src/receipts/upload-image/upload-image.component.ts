import { Component, EventEmitter, Input, Output, ViewChild } from "@angular/core";
import { take, tap } from "rxjs";
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
  @ViewChild("uploadInput") uploadInput!: any;

  @Input() public receiptId?: string = "";

  @Input() public multiple: boolean = true;

  @Input() public acceptFileTypes: string[] = [
    "image/*",
    "application/pdf",
    "image/heic",
  ];

  @Output() public fileLoaded: EventEmitter<ReceiptFileUploadCommand> =
    new EventEmitter();

  public formMode = FormMode;
  
  constructor(private receiptImageService: ReceiptImageService) {}

  public clickInput(): void {
    this.uploadInput.nativeElement.click();
  }

  public async onFileChange(event: any): Promise<void> {
    const files: File[] = Array.from(event?.target?.files ?? []);
    const acceptedFiles = files.filter((f) =>
      this.acceptFileTypes.some((t) => new RegExp(t).test(f.type))
    );

    for (let i = 0; i < acceptedFiles.length; i++) {
      const reader = new FileReader();
      const f = acceptedFiles[i];

      reader.onload = () => {
        const command: ReceiptFileUploadCommand = {
          file: f,
          receiptId: Number(this.receiptId),
        };

        if (f.type === "application/pdf" || f.type === "image/heic") {
          this.receiptImageService
            .convertToJpg(f)
            .pipe(
              take(1),
              tap((encodedImage) => {
                command.encodedImage = encodedImage.encodedImage;
                this.fileLoaded.emit(command);
              })
            )
            .subscribe();
        } else {
          this.fileLoaded.emit(command);
        }
      };

      reader.readAsArrayBuffer(f);
    }
  }
}
