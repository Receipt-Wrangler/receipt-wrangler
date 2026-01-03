import { Component, EventEmitter, HostListener, Input, OnChanges, Output, SimpleChanges } from "@angular/core";

@Component({
    selector: "app-image-viewer",
    templateUrl: "./image-viewer.component.html",
    styleUrl: "./image-viewer.component.scss",
    standalone: false
})
export class ImageViewerComponent implements OnChanges {
  @HostListener("wheel", ["$event"])
  public onWheel(event: WheelEvent) {
    this.wheel.emit(event);
  }

  @Input() public imageBase64?: string = "";

  @Input() public imageFile?: File;

  @Input() public scale: number = 1;

  @Output() public wheel: EventEmitter<WheelEvent> = new EventEmitter<WheelEvent>();

  public imageFileUrl: string = "";

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["imageFile"] && changes["imageFile"].currentValue) {
      this.setImageFileUrl(changes["imageFile"].currentValue);
    }
  }

  private setImageFileUrl(file: File): void {
    const reader = new FileReader();
    reader.onload = (event) => {
      this.imageFileUrl = (event?.target?.result ?? "") as string;
    };

    reader.readAsDataURL(file);
  }
}


