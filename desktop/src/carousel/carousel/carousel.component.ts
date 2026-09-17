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

  public readonly stageHeight = input<string>("60vh");

  public readonly initialIndex = input<number>(-1);

  public readonly removeButtonClicked = output<number>();

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
    this.activeViewer()?.zoomOut();
  }

  public zoomIn() {
    this.activeViewer()?.zoomIn();
  }

  /**
   * Each image owns its own zoom and pan, so the header buttons act on the slide
   * being looked at rather than on one scale shared by all of them.
   *
   * This indexes a viewChildren query by SLIDE index, so every slide must render
   * exactly one viewer or the two fall out of step - which is why the template
   * renders app-image-viewer unconditionally and lets the viewer decide whether
   * it has anything to show.
   */
  private activeViewer(): ImageViewerComponent | undefined {
    return this.viewers()[this.currentlyShownImageIndex];
  }

  public updateCurrentlyShownImage(index: number): void {
    this.currentlyShownImageIndex = index;
  }
}
