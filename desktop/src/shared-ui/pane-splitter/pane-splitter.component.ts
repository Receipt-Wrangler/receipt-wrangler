import {
  Component,
  ElementRef,
  Renderer2,
  RendererStyleFlags2,
  computed,
  effect,
  inject,
  input,
  model,
} from "@angular/core";

/**
 * The CSS custom property the splitter writes its ratio to. The container sizes
 * its right-hand pane from it, e.g. `flex: 0 0 var(--rw-pane-ratio, 50%)`.
 */
export const PANE_RATIO_PROPERTY = "--rw-pane-ratio";

/** How far one arrow-key press moves the splitter, in percentage points. */
const KEYBOARD_STEP = 2;

/**
 * A draggable vertical divider between two panes of a flex container.
 *
 * Presentational and gesture-owning: it reports the width of the pane to its
 * right as a percentage of the container, and writes that percentage to the
 * container as {@link PANE_RATIO_PROPERTY}. The container's own stylesheet
 * decides what to do with it, so nothing about the panes is baked in here.
 *
 * Deliberately NOT built on `cdkDrag`. cdkDrag resizes nothing — it *translates
 * the element it sits on*, which fights a splitter that the layout is already
 * moving as the ratio changes. Pointer events with `setPointerCapture` keep
 * tracking when the cursor leaves the handle and cover mouse, touch and pen in
 * one path.
 *
 * @example
 * <div class="panes" #panes>
 *   <section class="panes__left"></section>
 *   <app-pane-splitter [container]="panes" [(ratio)]="rightPaneRatio" />
 *   <section class="panes__right"></section>
 * </div>
 */
@Component({
  selector: "app-pane-splitter",
  templateUrl: "./pane-splitter.component.html",
  styleUrl: "./pane-splitter.component.scss",
  standalone: true,
  host: {
    class: "rw-pane-splitter",
    role: "separator",
    // `role="separator"` defaults to horizontal, which is the wrong axis for a
    // divider between a left and a right pane.
    "aria-orientation": "vertical",
    tabindex: "0",
    "[attr.aria-label]": "label()",
    "[attr.aria-valuenow]": "ariaValueNow()",
    "[attr.aria-valuemin]": "min()",
    "[attr.aria-valuemax]": "max()",
    "(pointerdown)": "onPointerDown($event)",
    "(keydown)": "onKeyDown($event)",
    "(dblclick)": "reset()",
  },
})
export class PaneSplitterComponent {
  /** The flex container whose right-hand pane this splitter sizes. */
  public readonly container = input.required<HTMLElement>();

  /** Width of the right-hand pane, as a percentage of the container. */
  public readonly ratio = model<number>(50);

  public readonly min = input<number>(20);

  public readonly max = input<number>(80);

  /** Where a double-click puts the splitter back to. */
  public readonly defaultRatio = input<number>(50);

  /**
   * The smallest the LEFT pane may become, in pixels. A percentage floor alone
   * lets it collapse under the intrinsic width of the controls inside it, which
   * then overflow instead of pushing back.
   */
  public readonly minLeftPaneWidth = input<number>(360);

  public readonly label = input<string>("Resize panel");

  protected readonly ariaValueNow = computed(() => Math.round(this.ratio()));

  private readonly host = inject<ElementRef<HTMLElement>>(ElementRef);

  private readonly renderer = inject(Renderer2);

  /** State captured at `pointerdown`; undefined whenever no drag is running. */
  private drag?: {
    pointerId: number;
    startX: number;
    startRatio: number;
    containerWidth: number;
  };

  public constructor() {
    // Mirrors the ratio onto the container. An effect rather than a template
    // binding because the container is an input, not part of this component's
    // own view — this is the "sync a signal to an imperative DOM API" case.
    effect(() => this.writeRatio(this.ratio()));
  }

  protected onPointerDown(event: PointerEvent): void {
    // Primary button only; touch and pen both report button 0.
    if (event.button !== 0) {
      return;
    }

    const containerWidth = this.containerWidth();
    if (containerWidth <= 0) {
      return;
    }

    // Suppresses the text selection a drag across both panes would otherwise
    // start. Touch scrolling is handled by `touch-action: none` in the styles.
    event.preventDefault();

    this.drag = {
      pointerId: event.pointerId,
      startX: event.clientX,
      startRatio: this.ratio(),
      containerWidth,
    };

    const element = this.host.nativeElement;
    // preventDefault above suppresses the focus a click would normally give the
    // handle, and losing it would mean the arrow keys only work after a Tab.
    element.focus();
    element.setPointerCapture(event.pointerId);
    element.addEventListener("pointermove", this.onPointerMove);
    element.addEventListener("pointerup", this.onPointerUp);
    element.addEventListener("pointercancel", this.onPointerUp);
    this.renderer.addClass(element, "rw-pane-splitter--dragging");
  }

  /**
   * Bound as native listeners rather than Angular host bindings on purpose: a
   * host binding schedules change detection on every event, and a pointermove
   * at 120Hz would then re-render the whole host page. The drag writes the
   * container's style directly and commits to the signal once, on release.
   */
  private readonly onPointerMove = (event: PointerEvent): void => {
    if (this.drag?.pointerId !== event.pointerId) {
      return;
    }

    this.writeRatio(this.ratioAt(event.clientX));
  };

  private readonly onPointerUp = (event: PointerEvent): void => {
    const drag = this.drag;
    if (drag?.pointerId !== event.pointerId) {
      return;
    }

    const ratio = this.ratioAt(event.clientX);
    this.endDrag(drag.pointerId);
    this.ratio.set(ratio);
  };

  protected onKeyDown(event: KeyboardEvent): void {
    const containerWidth = this.containerWidth();
    let next: number;

    switch (event.key) {
      // Moving the splitter left grows the pane on its right.
      case "ArrowLeft":
        next = this.ratio() + KEYBOARD_STEP;
        break;
      case "ArrowRight":
        next = this.ratio() - KEYBOARD_STEP;
        break;
      case "Home":
        next = this.max();
        break;
      case "End":
        next = this.min();
        break;
      default:
        return;
    }

    event.preventDefault();
    this.ratio.set(this.clamp(next, containerWidth));
  }

  protected reset(): void {
    this.ratio.set(this.clamp(this.defaultRatio(), this.containerWidth()));
  }

  /**
   * The ratio the given pointer position implies, measured from where the drag
   * started rather than accumulated from the last move — so it cannot drift,
   * and re-clamping it every move is idempotent.
   */
  private ratioAt(clientX: number): number {
    const drag = this.drag;
    if (!drag) {
      return this.ratio();
    }

    const delta = ((clientX - drag.startX) / drag.containerWidth) * 100;
    return this.clamp(drag.startRatio - delta, drag.containerWidth);
  }

  private clamp(ratio: number, containerWidth: number): number {
    return Math.min(Math.max(ratio, this.min()), this.effectiveMax(containerWidth));
  }

  /** {@link max}, lowered as needed to keep the left pane above its pixel floor. */
  private effectiveMax(containerWidth: number): number {
    if (containerWidth <= 0) {
      return this.max();
    }

    const room =
      ((containerWidth - this.minLeftPaneWidth() - this.host.nativeElement.offsetWidth) /
        containerWidth) *
      100;

    // Never below `min`, or a narrow container would invert the two bounds.
    return Math.max(this.min(), Math.min(this.max(), room));
  }

  /**
   * The container's content-box width — what a percentage flex-basis resolves
   * against. `getBoundingClientRect()` reports the border box instead.
   */
  private containerWidth(): number {
    const element = this.container();
    const styles = getComputedStyle(element);

    return (
      element.clientWidth -
      Number.parseFloat(styles.paddingLeft || "0") -
      Number.parseFloat(styles.paddingRight || "0")
    );
  }

  private writeRatio(ratio: number): void {
    this.renderer.setStyle(
      this.container(),
      PANE_RATIO_PROPERTY,
      `${ratio}%`,
      RendererStyleFlags2.DashCase,
    );
  }

  private endDrag(pointerId: number): void {
    const element = this.host.nativeElement;

    element.removeEventListener("pointermove", this.onPointerMove);
    element.removeEventListener("pointerup", this.onPointerUp);
    element.removeEventListener("pointercancel", this.onPointerUp);

    if (element.hasPointerCapture(pointerId)) {
      element.releasePointerCapture(pointerId);
    }

    this.renderer.removeClass(element, "rw-pane-splitter--dragging");
    this.drag = undefined;
  }
}
