import { Component, OnChanges, SimpleChanges, ViewEncapsulation, input, output, viewChildren } from "@angular/core";
import { ImageViewerComponent } from "src/shared-ui/image-viewer/image-viewer.component";
import { UntilDestroy } from "@ngneat/until-destroy";
import { FormMode } from "src/enums/form-mode.enum";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { FileDataView } from "../../open-api";

@UntilDestroy()
@Component({
    selector: "app-carousel",
    templateUrl: "./carousel.component.html",
    styleUrls: ["./carousel.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class CarouselComponent implements OnChanges {
  public readonly images = input<FileDataView[]>([]);

  public readonly imagePreviews = input<ReceiptFileUploadCommand[]>([]);

  public readonly disabled = input<boolean>(false);

  public readonly mode = input.required<FormMode>();

  public readonly hideButtonControls = input<boolean>(false);

  /**
   * Renders each slide on an interactive canvas instead of a plain scaled image.
   * Opt-in, following the hideButtonControls precedent, so the fullscreen dialog
   * — which renders this same component — is left exactly as it was.
   */
  public readonly directManipulation = input<boolean>(false);

  public readonly stageHeight = input<string>("60vh");

  public readonly initialIndex = input<number>(-1);

  public readonly removeButtonClicked = output<number>();

  public scale: number = 1;

  public currentlyShownImageIndex: number = 0;

  private readonly viewers = viewChildren(ImageViewerComponent);

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["initialIndex"]) {
      this.currentlyShownImageIndex = this.initialIndex();
    }
  }

  public emitRemoveButtonClicked(index: number): void {
    this.removeButtonClicked.emit(index);
  }

  public zoomOut() {
    if (this.directManipulation()) {
      this.activeViewer()?.zoomOut();
      return;
    }

    this.adjustScale(-0.1);
  }

  public zoomIn() {
    if (this.directManipulation()) {
      this.activeViewer()?.zoomIn();
      return;
    }

    this.adjustScale(0.1);
  }

  /**
   * On a canvas each image owns its own zoom and pan, so the header buttons act
   * on the slide being looked at rather than on one scale shared by all of them.
   */
  private activeViewer(): ImageViewerComponent | undefined {
    return this.viewers()[this.currentlyShownImageIndex];
  }

  public onScroll(event: WheelEvent): void {
    event.preventDefault();
    let value = event.deltaY * -0.000001;
    this.adjustScale(value);
  }

  public updateCurrentlyShownImage(index: number): void {
    this.currentlyShownImageIndex = index;
  }

  public adjustScale(amount: number): void {
    const newScale = this.scale + amount;
    this.scale = Math.max(newScale, 0.1);
  }
}
