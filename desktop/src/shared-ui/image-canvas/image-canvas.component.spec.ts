import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ImageCanvasComponent } from "./image-canvas.component";

const STAGE_WIDTH = 600;
const STAGE_HEIGHT = 400;
const NATURAL_WIDTH = 1200;
const NATURAL_HEIGHT = 1600;

/** Fit is limited by height here: 400 / 1600 = 0.25. */
const FIT_SCALE = STAGE_HEIGHT / NATURAL_HEIGHT;

/** Mirrors the component's own MAX_SCALE, which it does not export. */
const MAX_SCALE = 8;

/** Likewise MIN_VISIBLE_PX: the sliver that must stay on the stage. */
const MIN_VISIBLE_PX = 48;

/** The scale leaving the image's shorter side at MIN_VISIBLE_PX: 48 / 1200. */
const MIN_SCALE = MIN_VISIBLE_PX / Math.min(NATURAL_WIDTH, NATURAL_HEIGHT);

/**
 * jsdom implements neither `PointerEvent` nor pointer capture, so pointer events
 * are faked from `MouseEvent` — which carries every property the component reads
 * — and the capture calls are stubbed per element.
 */
function pointerEvent(
  type: string,
  clientX: number,
  clientY: number,
  pointerId = 1,
  button = 0,
): MouseEvent {
  const event = new MouseEvent(type, { bubbles: true, cancelable: true, clientX, clientY, button });
  Object.defineProperty(event, "pointerId", { value: pointerId });

  return event;
}

describe("ImageCanvasComponent", () => {
  let fixture: ComponentFixture<ImageCanvasComponent>;
  let component: ImageCanvasComponent;
  let host: HTMLElement;
  let image: HTMLImageElement;

  /** The image's transform, parsed back into the viewport it encodes. */
  const viewport = (): { x: number; y: number; scale: number } => {
    const match = /translate\((-?[\d.]+)px, (-?[\d.]+)px\) scale\(([\d.]+)\)/.exec(
      image.style.transform,
    );

    return match
      ? { x: Number(match[1]), y: Number(match[2]), scale: Number(match[3]) }
      : { x: NaN, y: NaN, scale: NaN };
  };

  const handle = (corner: string): HTMLElement =>
    host.querySelector(`[data-corner="${corner}"]`) as HTMLElement;

  const drag = async (
    target: HTMLElement,
    from: [number, number],
    to: [number, number],
    release = true,
  ) => {
    target.dispatchEvent(pointerEvent("pointerdown", from[0], from[1]));
    host.dispatchEvent(pointerEvent("pointermove", to[0], to[1]));
    if (release) {
      host.dispatchEvent(pointerEvent("pointerup", to[0], to[1]));
    }
    await fixture.whenStable();
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ImageCanvasComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(ImageCanvasComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("src", "data:image/png;base64,abc");
    await fixture.whenStable();

    host = fixture.nativeElement as HTMLElement;
    image = host.querySelector("img") as HTMLImageElement;

    // jsdom lays nothing out and decodes no images, so every measurement the
    // component reads has to be stubbed.
    Object.defineProperty(host, "clientWidth", { value: STAGE_WIDTH, configurable: true });
    Object.defineProperty(host, "clientHeight", { value: STAGE_HEIGHT, configurable: true });
    host.getBoundingClientRect = () => ({ left: 0, top: 0, width: STAGE_WIDTH, height: STAGE_HEIGHT }) as DOMRect;
    host.setPointerCapture = jest.fn();
    host.releasePointerCapture = jest.fn();
    host.hasPointerCapture = jest.fn().mockReturnValue(true);
    Object.defineProperty(image, "naturalWidth", { value: NATURAL_WIDTH, configurable: true });
    Object.defineProperty(image, "naturalHeight", { value: NATURAL_HEIGHT, configurable: true });

    image.dispatchEvent(new Event("load"));
    await fixture.whenStable();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("fits the image to the stage on load and centres it", () => {
    const view = viewport();

    expect(view.scale).toBeCloseTo(FIT_SCALE);
    // 1200 * 0.25 = 300 wide in a 600 stage, so 150 either side; the height fills.
    expect(view.x).toBeCloseTo((STAGE_WIDTH - NATURAL_WIDTH * FIT_SCALE) / 2);
    expect(view.y).toBeCloseTo(0);
  });

  it("renders a handle at each corner", () => {
    expect(host.querySelectorAll("[data-corner]").length).toEqual(4);
    ["nw", "ne", "se", "sw"].forEach((corner) => expect(handle(corner)).toBeTruthy());
  });

  it("positions the handle frame over the image", () => {
    const frame = host.querySelector(".rw-image-canvas__frame") as HTMLElement;

    expect(frame.style.width).toEqual(`${NATURAL_WIDTH * FIT_SCALE}px`);
    expect(frame.style.height).toEqual(`${NATURAL_HEIGHT * FIT_SCALE}px`);
  });

  // Pulling a corner scales about the corner opposite it, which stays put for the
  // whole gesture - in BOTH axes, including the one narrower than the stage.
  it("scales about the opposite corner when a handle is pulled", async () => {
    const before = viewport();
    const anchorX = before.x + NATURAL_WIDTH * before.scale;
    const anchorY = before.y + NATURAL_HEIGHT * before.scale;

    // Drag the top-left handle up and to the left: the image grows.
    await drag(handle("nw"), [before.x, before.y], [before.x - 200, before.y - 200]);

    const after = viewport();
    expect(after.scale).toBeGreaterThan(before.scale);
    expect(after.x + NATURAL_WIDTH * after.scale).toBeCloseTo(anchorX, 0);
    expect(after.y + NATURAL_HEIGHT * after.scale).toBeCloseTo(anchorY, 0);
  });

  // The headline of a real handle: it shrinks as well as grows. The image opens
  // at the fit, so a handle that cannot go below it is a handle that only works
  // in one direction.
  it("shrinks the image below the fit when a handle is pulled inward", async () => {
    const before = viewport();

    // The se handle sits at the image's bottom-right; pull it in toward the nw
    // corner, which is the anchor.
    await drag(handle("se"), [450, 400], [250, 200]);

    const after = viewport();
    expect(after.scale).toBeLessThan(FIT_SCALE);
    expect(NATURAL_WIDTH * after.scale).toBeLessThan(STAGE_WIDTH);
    expect(NATURAL_HEIGHT * after.scale).toBeLessThan(STAGE_HEIGHT);

    // The nw anchor held in BOTH axes - nothing re-centred it.
    expect(after.x).toBeCloseTo(before.x);
    expect(after.y).toBeCloseTo(before.y);
  });

  // The axis narrower than the stage used to be auto-centred, which threw the
  // anchor away in exactly the axis a letterboxed receipt has spare room in. At
  // the fit scale this image is 300px wide in a 600px stage, so this grows it
  // while it is still narrower than the stage and checks the anchor survives.
  it("holds the anchor in an axis narrower than the stage", async () => {
    const before = viewport();
    const anchorX = before.x + NATURAL_WIDTH * before.scale;

    await drag(handle("nw"), [150, 0], [100, -100]);

    const after = viewport();
    expect(NATURAL_WIDTH * after.scale).toBeLessThan(STAGE_WIDTH);
    expect(after.x + NATURAL_WIDTH * after.scale).toBeCloseTo(anchorX, 0);
  });

  it("will not shrink past the minimum size when a handle is pulled inward", async () => {
    // Pull the se handle right onto its own anchor, which implies a scale of 0.
    await drag(handle("se"), [450, 400], [150, 0]);

    const after = viewport();
    expect(after.scale).toBeCloseTo(MIN_SCALE);
    expect(NATURAL_WIDTH * after.scale).toBeCloseTo(MIN_VISIBLE_PX);
  });

  // A distorted receipt is never wanted, so a corner drag is uniform.
  it("keeps the aspect ratio locked while resizing", async () => {
    const before = viewport();
    await drag(handle("se"), [0, 0], [600, 10]);

    const after = viewport();
    const frame = host.querySelector(".rw-image-canvas__frame") as HTMLElement;
    const ratio = Number.parseFloat(frame.style.width) / Number.parseFloat(frame.style.height);

    expect(after.scale).not.toEqual(before.scale);
    expect(ratio).toBeCloseTo(NATURAL_WIDTH / NATURAL_HEIGHT);
  });

  it("pans a zoomed image and stops once only a sliver of it is left", async () => {
    component.zoomIn();
    component.zoomIn();
    await fixture.whenStable();

    const before = viewport();
    await drag(host, [300, 200], [300, 120]);

    const after = viewport();
    expect(after.y).toBeLessThan(before.y);

    // Dragging far past the edge stops with MIN_VISIBLE_PX still on the stage.
    await drag(host, [300, 200], [300, 5000]);
    expect(viewport().y).toEqual(STAGE_HEIGHT - MIN_VISIBLE_PX);
  });

  // An image smaller than the stage is an object on it, not something pinned to
  // its middle: it goes where you put it.
  it("pans an axis narrower than the stage freely", async () => {
    const before = viewport();
    await drag(host, [300, 200], [40, 200]);

    expect(viewport().x).toBeCloseTo(before.x - 260);
  });

  it("lets the image be pushed into a corner, leaving a sliver on the stage", async () => {
    await drag(host, [300, 200], [-5000, -5000]);

    const after = viewport();
    expect(after.x).toBeCloseTo(MIN_VISIBLE_PX - NATURAL_WIDTH * after.scale);
    expect(after.y).toBeCloseTo(MIN_VISIBLE_PX - NATURAL_HEIGHT * after.scale);
  });

  it("will not zoom out past the minimum scale", async () => {
    // Well past the ~14 steps of 1.15 it takes to fall from the fit to the floor.
    for (let i = 0; i < 30; i += 1) {
      component.zoomOut();
    }
    await fixture.whenStable();

    expect(viewport().scale).toBeCloseTo(MIN_SCALE);
  });

  it("will not zoom in past the maximum scale", async () => {
    // Well past the ~25 steps it takes to climb from the fit to the cap.
    for (let i = 0; i < 40; i += 1) {
      component.zoomIn();
    }
    await fixture.whenStable();

    expect(viewport().scale).toBeCloseTo(MAX_SCALE);
  });

  it("zooms in and back out again", async () => {
    component.zoomIn();
    await fixture.whenStable();
    const zoomed = viewport().scale;
    expect(zoomed).toBeGreaterThan(FIT_SCALE);

    component.zoomOut();
    await fixture.whenStable();
    expect(viewport().scale).toBeCloseTo(FIT_SCALE);
  });

  // The viewer this replaced multiplied deltaY by -0.000001, so a notch moved the
  // scale by 0.0001 and wheel zoom did nothing at all.
  it("zooms a meaningful amount on the wheel, anchored at the cursor", async () => {
    const before = viewport();
    const cursorY = 200;

    host.dispatchEvent(
      new WheelEvent("wheel", { deltaY: -100, clientX: 300, clientY: cursorY, cancelable: true }),
    );
    await fixture.whenStable();

    const after = viewport();
    expect(after.scale).toBeGreaterThan(before.scale * 1.1);
    // The image point under the cursor did not move. Asserted vertically, the
    // axis that overflows the stage - see the centring rule above.
    expect((cursorY - before.y) / before.scale).toBeCloseTo((cursorY - after.y) / after.scale, 1);
  });

  // A horizontal wheel (a two-finger trackpad swipe) reports deltaY === 0, which
  // the zoom ternary would read as "not < 0" and treat as a zoom OUT. The event
  // must also keep its default, or the page loses the horizontal scroll.
  it("ignores a horizontal wheel instead of zooming out", async () => {
    const before = viewport();
    const event = new WheelEvent("wheel", {
      deltaY: 0,
      deltaX: 120,
      clientX: 300,
      clientY: 200,
      cancelable: true,
    });

    host.dispatchEvent(event);
    await fixture.whenStable();

    expect(viewport()).toEqual(before);
    expect(event.defaultPrevented).toEqual(false);
  });

  it("returns to the fit on double-click", async () => {
    component.zoomIn();
    component.zoomIn();
    await fixture.whenStable();
    expect(viewport().scale).toBeGreaterThan(FIT_SCALE);

    host.dispatchEvent(new MouseEvent("dblclick", { bubbles: true }));
    await fixture.whenStable();

    expect(viewport().scale).toBeCloseTo(FIT_SCALE);
  });

  it("ignores a non-primary button", async () => {
    const before = viewport();
    host.dispatchEvent(pointerEvent("pointerdown", 300, 200, 1, 2));
    host.dispatchEvent(pointerEvent("pointermove", 300, 100));
    await fixture.whenStable();

    expect(viewport()).toEqual(before);
  });

  it("stops tracking once the drag ends", async () => {
    component.zoomIn();
    component.zoomIn();
    await fixture.whenStable();

    await drag(host, [300, 200], [300, 150]);
    const after = viewport();

    host.dispatchEvent(pointerEvent("pointermove", 300, 50));
    await fixture.whenStable();

    expect(viewport()).toEqual(after);
  });

  // preventDefault on pointerdown suppresses the compatibility mouse events but
  // NOT dblclick, so without an explicit stop a quick double pull of a corner
  // bubbles to the host and throws away the size just dragged to.
  it("does not fit when a handle is double-clicked", async () => {
    // Pull the top-left handle out so the image is bigger than the fit.
    const start = viewport();
    await drag(handle("nw"), [start.x, start.y], [start.x - 200, start.y - 200]);
    const resized = viewport();
    expect(resized.scale).toBeGreaterThan(FIT_SCALE);

    handle("nw").dispatchEvent(new MouseEvent("dblclick", { bubbles: true }));
    await fixture.whenStable();

    expect(viewport()).toEqual(resized);
  });

  // A handle sits over the image, which starts a pan on pointerdown.
  it("does not start a pan when a handle is grabbed", async () => {
    const before = viewport();
    handle("se").dispatchEvent(pointerEvent("pointerdown", 100, 100));
    host.dispatchEvent(pointerEvent("pointermove", 140, 140));
    await fixture.whenStable();

    // Resized about the top-left, so the origin is unchanged - a pan would have
    // moved it by the full 40px in each axis.
    expect(viewport().scale).not.toEqual(before.scale);
    expect(viewport().x).toEqual(before.x);
    expect(viewport().y).toEqual(before.y);
  });
});
