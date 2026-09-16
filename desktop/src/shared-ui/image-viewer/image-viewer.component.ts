import { Component, HostListener, Input, OnChanges, OnDestroy, SimpleChanges, input, output, signal, viewChild } from "@angular/core";
import { ImageCanvasComponent } from "../image-canvas/image-canvas.component";

@Component({
    selector: "app-image-viewer",
    templateUrl: "./image-viewer.component.html",
    styleUrl: "./image-viewer.component.scss",
    standalone: false
})
export class ImageViewerComponent implements OnChanges, OnDestroy {
  @HostListener("wheel", ["$event"])
  public onWheel(event: WheelEvent) {
    // The canvas owns the wheel in direct-manipulation mode; re-emitting would
    // also drive the carousel's shared scale, which nothing is rendering then.
    if (this.directManipulation()) {
      return;
    }

    this.wheel.emit(event);
  }

  @Input() public imageBase64?: string = "";

  public readonly imageFile = input<File>();

  public readonly scale = input<number>(1);

  /**
   * Renders the image on an interactive canvas — corner handles, pan, zoom —
   * instead of the plain scaled image. Opt-in so the fullscreen dialog, which
   * shares this component, keeps the behaviour it has always had.
   */
  public readonly directManipulation = input<boolean>(false);

  /** Height of the canvas stage. Ignored unless directManipulation is on. */
  public readonly stageHeight = input<string>("60vh");

  public readonly wheel = output<WheelEvent>();

  public imageFileUrl = signal("");

  private activeReader?: FileReader;

  private readonly canvas = viewChild(ImageCanvasComponent);

  /** No-ops unless directManipulation is on; the canvas owns the zoom then. */
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


