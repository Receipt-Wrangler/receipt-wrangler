import { Component, Input, OnChanges, OnDestroy, SimpleChanges, input, signal, viewChild } from "@angular/core";
import { ImageCanvasComponent } from "../image-canvas/image-canvas.component";

@Component({
    selector: "app-image-viewer",
    templateUrl: "./image-viewer.component.html",
    standalone: false
})
export class ImageViewerComponent implements OnChanges, OnDestroy {
  @Input() public imageBase64?: string = "";

  public readonly imageFile = input<File>();

  /** Height of the canvas stage. */
  public readonly stageHeight = input<string>("60vh");

  public imageFileUrl = signal("");

  private activeReader?: FileReader;

  private readonly canvas = viewChild(ImageCanvasComponent);

  /** No-ops until the image has loaded and the canvas exists. */
  public zoomIn(): void {
    this.canvas()?.zoomIn();
  }

  public zoomOut(): void {
    this.canvas()?.zoomOut();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["imageFile"] && changes["imageFile"].currentValue) {
      this.setImageFileUrl(changes["imageFile"].currentValue);
    }
  }

  public ngOnDestroy(): void {
    this.activeReader?.abort();
  }

  private setImageFileUrl(file: File): void {
    this.activeReader?.abort();

    const reader = new FileReader();
    this.activeReader = reader;

    reader.onload = (event) => {
      this.imageFileUrl.set((event?.target?.result ?? "") as string);
    };

    reader.readAsDataURL(file);
  }
}


