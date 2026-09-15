import { Component, ElementRef, provideZonelessChangeDetection, signal, viewChild } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { PANE_RATIO_PROPERTY, PaneSplitterComponent } from "./pane-splitter.component";

const CONTAINER_WIDTH = 1000;

@Component({
  standalone: true,
  imports: [PaneSplitterComponent],
  template: `
    <div class="panes" #panes>
      <div class="panes__left"></div>
      <app-pane-splitter
        [container]="panes"
        [min]="min()"
        [max]="max()"
        [defaultRatio]="defaultRatio()"
        [minLeftPaneWidth]="minLeftPaneWidth()"
        [(ratio)]="ratio"
      ></app-pane-splitter>
      <div class="panes__right"></div>
    </div>
  `,
})
class HostComponent {
  public readonly panes = viewChild.required<ElementRef<HTMLElement>>("panes");

  public readonly ratio = signal(50);

  public readonly min = signal(20);

  public readonly max = signal(80);

  public readonly defaultRatio = signal(50);

  // Defaulted off so the percentage bounds are what most cases exercise; the
  // pixel floor gets its own test.
  public readonly minLeftPaneWidth = signal(0);
}

/**
 * jsdom implements neither `PointerEvent` nor pointer capture, so pointer events
 * are faked from `MouseEvent` (which carries every property the component reads)
 * and the capture calls are stubbed per element.
 */
function pointerEvent(type: string, clientX: number, pointerId = 1, button = 0): MouseEvent {
  const event = new MouseEvent(type, { bubbles: true, cancelable: true, clientX, button });
  Object.defineProperty(event, "pointerId", { value: pointerId });

  return event;
}

describe("PaneSplitterComponent", () => {
  let fixture: ComponentFixture<HostComponent>;
  let host: HostComponent;
  let splitter: HTMLElement;
  let container: HTMLElement;

  const ratioProperty = (): string =>
    container.style.getPropertyValue(PANE_RATIO_PROPERTY).trim();

  const drag = async (from: number, to: number, release = true) => {
    splitter.dispatchEvent(pointerEvent("pointerdown", from));
    splitter.dispatchEvent(pointerEvent("pointermove", to));
    if (release) {
      splitter.dispatchEvent(pointerEvent("pointerup", to));
    }
    await fixture.whenStable();
  };

  const press = async (key: string) => {
    splitter.dispatchEvent(new KeyboardEvent("keydown", { key, bubbles: true, cancelable: true }));
    await fixture.whenStable();
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HostComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(HostComponent);
    host = fixture.componentInstance;
    await fixture.whenStable();

    container = fixture.nativeElement.querySelector(".panes") as HTMLElement;
    splitter = fixture.nativeElement.querySelector("app-pane-splitter") as HTMLElement;

    // jsdom lays nothing out, so the measurements the component reads are stubbed.
    Object.defineProperty(container, "clientWidth", {
      value: CONTAINER_WIDTH,
      configurable: true,
    });
    splitter.setPointerCapture = jest.fn();
    splitter.releasePointerCapture = jest.fn();
    splitter.hasPointerCapture = jest.fn().mockReturnValue(true);
  });

  it("should create", () => {
    expect(splitter).toBeTruthy();
  });

  it("writes the initial ratio to the container", () => {
    expect(ratioProperty()).toEqual("50%");
  });

  // A separator defaults to horizontal, which is the wrong axis for a divider
  // between a left and a right pane.
  it("describes itself as a vertical separator", () => {
    expect(splitter.getAttribute("role")).toEqual("separator");
    expect(splitter.getAttribute("aria-orientation")).toEqual("vertical");
    expect(splitter.getAttribute("aria-valuenow")).toEqual("50");
    expect(splitter.getAttribute("aria-valuemin")).toEqual("20");
    expect(splitter.getAttribute("aria-valuemax")).toEqual("80");
    expect(splitter.getAttribute("tabindex")).toEqual("0");
  });

  // Dragging left grows the pane on the splitter's right.
  it("grows the right pane when dragged left", async () => {
    await drag(500, 400);

    expect(host.ratio()).toEqual(60);
    expect(ratioProperty()).toEqual("60%");
  });

  it("shrinks the right pane when dragged right", async () => {
    await drag(500, 650);

    expect(host.ratio()).toEqual(35);
  });

  // The ratio is recomputed from where the drag started rather than accumulated
  // from the previous move, so a drag out past a bound and back is not clipped.
  it("recomputes from the drag origin rather than accumulating", async () => {
    splitter.dispatchEvent(pointerEvent("pointerdown", 500));
    // Far past the max, then back to a legal position.
    splitter.dispatchEvent(pointerEvent("pointermove", 100));
    splitter.dispatchEvent(pointerEvent("pointermove", 550));
    splitter.dispatchEvent(pointerEvent("pointerup", 550));
    await fixture.whenStable();

    expect(host.ratio()).toEqual(45);
  });

  it("clamps to min and max", async () => {
    await drag(500, -2000);
    expect(host.ratio()).toEqual(80);

    await drag(500, 2000);
    expect(host.ratio()).toEqual(20);
  });

  // A signal write per pointermove would re-render the whole host page, so the
  // drag writes the container's style directly and commits once, on release.
  it("updates the container during the drag but commits the ratio only on release", async () => {
    await drag(500, 400, false);

    expect(ratioProperty()).toEqual("60%");
    expect(host.ratio()).toEqual(50);

    splitter.dispatchEvent(pointerEvent("pointerup", 400));
    await fixture.whenStable();

    expect(host.ratio()).toEqual(60);
  });

  it("ignores a non-primary button", async () => {
    splitter.dispatchEvent(pointerEvent("pointerdown", 500, 1, 2));
    splitter.dispatchEvent(pointerEvent("pointermove", 400));
    await fixture.whenStable();

    expect(host.ratio()).toEqual(50);
  });

  it("stops tracking once the drag ends", async () => {
    await drag(500, 400);
    splitter.dispatchEvent(pointerEvent("pointermove", 200));
    await fixture.whenStable();

    expect(host.ratio()).toEqual(60);
  });

  it("moves on the arrow keys and jumps to the bounds on Home and End", async () => {
    await press("ArrowLeft");
    expect(host.ratio()).toEqual(52);

    await press("ArrowRight");
    expect(host.ratio()).toEqual(50);

    await press("Home");
    expect(host.ratio()).toEqual(80);

    await press("End");
    expect(host.ratio()).toEqual(20);
  });

  it("leaves other keys alone", async () => {
    await press("Enter");

    expect(host.ratio()).toEqual(50);
  });

  it("returns to the default ratio on double-click", async () => {
    host.defaultRatio.set(65);
    await drag(500, 700);
    expect(host.ratio()).not.toEqual(65);

    splitter.dispatchEvent(new MouseEvent("dblclick", { bubbles: true }));
    await fixture.whenStable();

    expect(host.ratio()).toEqual(65);
  });

  // A percentage floor alone lets the left pane collapse under the intrinsic
  // width of the controls inside it, which then overflow instead of pushing back.
  it("lowers the effective max to keep the left pane above its pixel floor", async () => {
    host.minLeftPaneWidth.set(400);
    await fixture.whenStable();

    await drag(500, -2000);

    // 400px of 1000px has to stay on the left, so the right pane stops at 60%
    // rather than the configured 80%.
    expect(host.ratio()).toEqual(60);
  });

  it("never inverts the bounds when the container cannot fit the floor", async () => {
    host.minLeftPaneWidth.set(5000);
    await fixture.whenStable();

    await drag(500, -2000);

    expect(host.ratio()).toEqual(20);
  });

  it("does not start a drag when the container has no width", async () => {
    Object.defineProperty(container, "clientWidth", { value: 0, configurable: true });

    await drag(500, 400);

    expect(host.ratio()).toEqual(50);
  });
});
