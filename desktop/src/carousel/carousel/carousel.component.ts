import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges, ViewEncapsulation, } from "@angular/core";
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
  @Input() public images: FileDataView[] = [];

  @Input() public imagePreviews: ReceiptFileUploadCommand[] = [];

  @Input() public disabled: boolean = false;

  @Input() public mode!: FormMode;

  @Input() public hideButtonControls: boolean = false;

  @Input() public initialIndex: number = -1;

  @Output() public removeButtonClicked: EventEmitter<number> =
    new EventEmitter<number>();

  public scale: number = 1;

  public currentlyShownImageIndex: number = 0;

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["initialIndex"]) {
      this.currentlyShownImageIndex = this.initialIndex;
    }
  }

  public emitRemoveButtonClicked(index: number): void {
    this.removeButtonClicked.emit(index);
  }

  public zoomOut() {
    this.adjustScale(-0.1);
  }

  public zoomIn() {
    this.adjustScale(0.1);
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
