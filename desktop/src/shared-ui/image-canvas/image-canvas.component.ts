import {
  AfterViewInit,
  Component,
  ElementRef,
  OnDestroy,
  Renderer2,
  effect,
  inject,
  input,
  viewChild,
} from "@angular/core";

/** The corner a resize handle is pinned to. */
export type ImageCanvasCorner = "nw" | "ne" | "se" | "sw";

const CORNERS: ImageCanvasCorner[] = ["nw", "ne", "se", "sw"];

/** Largest the image may be drawn, as a multiple of its natural size. */
const MAX_SCALE = 8;

/** Applied per wheel notch and per zoom-button press. */
const ZOOM_STEP = 1.15;

/** Where the image sits on the stage: its top-left in stage px, and its scale. */
interface Viewport {
  scale: number;
  x: number;
  y: number;
}

type Drag =
  | {
      kind: "pan";
      pointerId: number;
      startX: number;
      startY: number;
      origin: Viewport;
    }
  | {
      kind: "resize";
      pointerId: number;
      corner: ImageCanvasCorner;
      anchorX: number;
      anchorY: number;
    };

/**
 * An image on a fixed stage, manipulated directly: pull a corner handle to
 * resize it, drag to pan, wheel to zoom, double-click to fit.
 *
 * The stage never changes size — the image is clipped by it and you pan around,
 * the way a canvas viewport behaves — so nothing here can move the surrounding
 * page layout.
 *
 * Scale and pan are a **single transform on a single element**. Keeping them on
 * separate elements (a scaled wrapper around a translated image) makes drag
 * distance desynchronise from the cursor by a factor of the scale, which is what
 * the viewer this replaces did.
 *
 * @example
 * <app-image-canvas [src]="dataUrl" stageHeight="60vh"></app-image-canvas>
 */
@Component({
  selector: "app-image-canvas",
  templateUrl: "./image-canvas.component.html",
  styleUrl: "./image-canvas.component.scss",
  standalone: true,
  host: {
    class: "rw-image-canvas",
    "[style.height]": "stageHeight()",
    "(pointerdown)": "onPanStart($event)",
    "(wheel)": "onWheel($event)",
    "(dblclick)": "fit()",
  },
})
export class ImageCanvasComponent implements AfterViewInit, OnDestroy {
  public readonly src = input.required<string>();

  public readonly stageHeight = input<string>("60vh");

  public readonly alt = input<string>("Receipt image");

  protected readonly corners = CORNERS;

  protected readonly image = viewChild<ElementRef<HTMLImageElement>>("image");

  protected readonly frame = viewChild<ElementRef<HTMLElement>>("frame");

  private readonly host = inject<ElementRef<HTMLElement>>(ElementRef);

  private readonly renderer = inject(Renderer2);

  /**
   * The live viewport, and deliberately a plain field rather than a signal:
   * nothing in the template reads it — the image and the handle frame are both
   * positioned imperatively — so a 120Hz drag runs no change detection at all
   * over the host page, which here is the receipt form's large template.
   */
  private view: Viewport = { scale: 1, x: 0, y: 0 };

  private natural = { width: 0, height: 0 };

  private drag?: Drag;

  private observer?: ResizeObserver;

  public constructor() {
    // A new image invalidates the measurements; (load) re-establishes them.
    effect(() => {
      this.src();
      this.natural = { width: 0, height: 0 };
    });
  }

  public ngAfterViewInit(): void {
    if (typeof ResizeObserver === "undefined") {
      return;
    }

    this.observer = new ResizeObserver(() => this.onStageResized());
    this.observer.observe(this.host.nativeElement);
  }

  public ngOnDestroy(): void {
    this.observer?.disconnect();
    this.endDrag();
  }

  /** Sizes the image to the stage and centres it. */
  public fit(): void {
    this.commit({ scale: this.fitScale(), x: 0, y: 0 });
  }

  public zoomIn(): void {
    this.zoomAroundCentre(ZOOM_STEP);
  }

  public zoomOut(): void {
    this.zoomAroundCentre(1 / ZOOM_STEP);
  }

  protected onImageLoad(): void {
    const image = this.image()?.nativeElement;
    if (!image) {
      return;
    }

    this.natural = { width: image.naturalWidth, height: image.naturalHeight };
    this.fit();
  }

  protected onWheel(event: WheelEvent): void {
    if (!this.natural.width) {
      return;
    }

    event.preventDefault();

    const rect = this.host.nativeElement.getBoundingClientRect();
    this.zoomAround(
      event.deltaY < 0 ? ZOOM_STEP : 1 / ZOOM_STEP,
      event.clientX - rect.left,
      event.clientY - rect.top,
    );
  }

  protected onPanStart(event: PointerEvent): void {
    if (event.button !== 0 || !this.natural.width) {
      return;
    }

    event.preventDefault();
    this.beginDrag(event, {
      kind: "pan",
      pointerId: event.pointerId,
      startX: event.clientX,
      startY: event.clientY,
      origin: { ...this.view },
    });
  }

  protected onResizeStart(event: PointerEvent, corner: ImageCanvasCorner): void {
    if (event.button !== 0 || !this.natural.width) {
      return;
    }

    event.preventDefault();
    // The handles sit inside the stage, which starts a pan on pointerdown.
    event.stopPropagation();

    // The corner opposite the one being pulled stays put for the whole gesture.
    const box = this.imageBox();
    this.beginDrag(event, {
      kind: "resize",
      pointerId: event.pointerId,
      corner,
      anchorX: corner.includes("w") ? box.x + box.width : box.x,
      anchorY: corner.includes("n") ? box.y + box.height : box.y,
    });
  }

  /**
   * Bound as native listeners rather than host bindings: a host binding schedules
   * change detection on every event, which for a pointermove would re-render the
   * whole host page. See the note on {@link view}.
   */
  private readonly onPointerMove = (event: PointerEvent): void => {
    if (this.drag?.pointerId !== event.pointerId) {
      return;
    }

    this.view = this.clamp(this.viewportAt(event));
    this.paint();
  };

  private readonly onPointerUp = (event: PointerEvent): void => {
    if (this.drag?.pointerId !== event.pointerId) {
      return;
    }

    this.view = this.clamp(this.viewportAt(event));
    this.endDrag();
    this.paint();
  };

  /**
   * The viewport the given pointer position implies. Measured from where the
   * gesture started rather than accumulated from the previous move, so it cannot
   * drift and re-clamping every move is idempotent.
   */
  private viewportAt(event: PointerEvent): Viewport {
    const drag = this.drag;
    if (!drag) {
      return this.view;
    }

    if (drag.kind === "pan") {
      return {
        scale: drag.origin.scale,
        x: drag.origin.x + (event.clientX - drag.startX),
        y: drag.origin.y + (event.clientY - drag.startY),
      };
    }

    // Re-read the rect each move: the page can scroll mid-drag.
    const rect = this.host.nativeElement.getBoundingClientRect();
    const pointerX = event.clientX - rect.left;
    const pointerY = event.clientY - rect.top;

    // Aspect-locked — a distorted receipt is never wanted — following whichever
    // axis was pulled further.
    const scale = this.clampScale(
      Math.max(
        Math.abs(pointerX - drag.anchorX) / this.natural.width,
        Math.abs(pointerY - drag.anchorY) / this.natural.height,
      ),
    );

    const width = this.natural.width * scale;
    const height = this.natural.height * scale;

    return {
      scale,
      x: drag.corner.includes("w") ? drag.anchorX - width : drag.anchorX,
      y: drag.corner.includes("n") ? drag.anchorY - height : drag.anchorY,
    };
  }

  private beginDrag(event: PointerEvent, drag: Drag): void {
    this.drag = drag;

    const element = this.host.nativeElement;
    element.setPointerCapture(event.pointerId);
    element.addEventListener("pointermove", this.onPointerMove);
    element.addEventListener("pointerup", this.onPointerUp);
    element.addEventListener("pointercancel", this.onPointerUp);
    this.renderer.addClass(element, `rw-image-canvas--${drag.kind}`);
  }

  private endDrag(): void {
    const drag = this.drag;
    if (!drag) {
      return;
    }

    const element = this.host.nativeElement;
    element.removeEventListener("pointermove", this.onPointerMove);
    element.removeEventListener("pointerup", this.onPointerUp);
    element.removeEventListener("pointercancel", this.onPointerUp);

    if (element.hasPointerCapture(drag.pointerId)) {
      element.releasePointerCapture(drag.pointerId);
    }

    this.renderer.removeClass(element, `rw-image-canvas--${drag.kind}`);
    this.drag = undefined;
  }

  private onStageResized(): void {
    if (!this.natural.width) {
      return;
    }

    // clampScale floors at the fit, so a stage that grew pulls an image that was
    // showing in full back up to filling it — without undoing a deliberate zoom.
    this.commit({ ...this.view, scale: this.clampScale(this.view.scale) });
  }

  private zoomAroundCentre(factor: number): void {
    const { width, height } = this.stageSize();
    this.zoomAround(factor, width / 2, height / 2);
  }

  /** Zooms while holding the image point under (originX, originY) still. */
  private zoomAround(factor: number, originX: number, originY: number): void {
    const scale = this.clampScale(this.view.scale * factor);
    if (scale === this.view.scale) {
      return;
    }

    const ratio = scale / this.view.scale;
    this.commit({
      scale,
      x: originX - (originX - this.view.x) * ratio,
      y: originY - (originY - this.view.y) * ratio,
    });
  }

  private commit(view: Viewport): void {
    this.view = this.clamp(view);
    this.paint();
  }

  private clamp(view: Viewport): Viewport {
    const stage = this.stageSize();

    return {
      scale: view.scale,
      x: this.clampAxis(view.x, this.natural.width * view.scale, stage.width),
      y: this.clampAxis(view.y, this.natural.height * view.scale, stage.height),
    };
  }

  private clampAxis(value: number, content: number, stage: number): number {
    // Smaller than the stage: centre it, so a shrunk image cannot drift into a
    // corner. Larger: stop either edge coming inside the stage.
    if (content <= stage) {
      return (stage - content) / 2;
    }

    return Math.min(Math.max(value, stage - content), 0);
  }

  private clampScale(scale: number): number {
    return Math.min(Math.max(scale, this.fitScale()), MAX_SCALE);
  }

  /** The scale at which the whole image is visible; never magnifies past 1:1. */
  private fitScale(): number {
    const stage = this.stageSize();
    if (!this.natural.width || !this.natural.height || !stage.width || !stage.height) {
      return 1;
    }

    return Math.min(stage.width / this.natural.width, stage.height / this.natural.height, 1);
  }

  private stageSize(): { width: number; height: number } {
    const element = this.host.nativeElement;

    return { width: element.clientWidth, height: element.clientHeight };
  }

  private imageBox(): { x: number; y: number; width: number; height: number } {
    return {
      x: this.view.x,
      y: this.view.y,
      width: this.natural.width * this.view.scale,
      height: this.natural.height * this.view.scale,
    };
  }

  private paint(): void {
    const image = this.image()?.nativeElement;
    const frame = this.frame()?.nativeElement;
    if (!image || !frame) {
      return;
    }

    const { x, y, width, height } = this.imageBox();

    this.renderer.setStyle(image, "transform", `translate(${x}px, ${y}px) scale(${this.view.scale})`);
    this.renderer.setStyle(frame, "transform", `translate(${x}px, ${y}px)`);
    this.renderer.setStyle(frame, "width", `${width}px`);
    this.renderer.setStyle(frame, "height", `${height}px`);
  }
}
