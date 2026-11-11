import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { CancelButtonComponent } from "./cancel-button.component";

describe("CancelButtonComponent", () => {
  let component: CancelButtonComponent;
  let fixture: ComponentFixture<CancelButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CancelButtonComponent],
      providers: [{ provide: ActivatedRoute, useValue: {} }],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(CancelButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  // The button is icon-only, so its tooltip text doubles as its accessible name.
  it("names the button after its tooltip", () => {
    const button = (): HTMLButtonElement => fixture.nativeElement.querySelector("button");
    expect(button().getAttribute("aria-label")).toBe("Cancel");

    fixture.componentRef.setInput("tooltip", "Close");
    fixture.detectChanges();

    expect(button().getAttribute("aria-label")).toBe("Close");
  });
});
