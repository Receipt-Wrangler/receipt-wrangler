import { Component, provideZonelessChangeDetection, signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideRouter } from "@angular/router";
import { ReceiptClaimRef } from "../claim.utils";
import { ReceiptClaimLinkComponent } from "./receipt-claim-link.component";

@Component({
  standalone: true,
  imports: [ReceiptClaimLinkComponent],
  template: `<app-receipt-claim-link [receipt]="receipt()" [labelled]="labelled()"></app-receipt-claim-link>`,
})
class HostComponent {
  public readonly receipt = signal<ReceiptClaimRef | undefined>(undefined);
  public readonly labelled = signal(false);
}

describe("ReceiptClaimLinkComponent", () => {
  let fixture: ComponentFixture<HostComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HostComponent],
      providers: [provideZonelessChangeDetection(), provideRouter([])],
    }).compileComponents();

    fixture = TestBed.createComponent(HostComponent);
  });

  const render = async (receipt: ReceiptClaimRef | undefined, labelled = false) => {
    fixture.componentInstance.receipt.set(receipt);
    fixture.componentInstance.labelled.set(labelled);
    await fixture.whenStable();
  };

  const link = (): HTMLAnchorElement | null =>
    fixture.nativeElement.querySelector("[data-testid='receipt-claim-link']");

  it("links the claim's name to its page", async () => {
    await render({ claimId: 12, claimName: "Voucher 14" });

    expect(link()?.textContent?.trim()).toBe("Voucher 14");
    expect(link()?.getAttribute("href")).toBe("/claims/12/view");
    expect(fixture.nativeElement.textContent).not.toContain("In claim:");
  });

  it("prefixes the label when asked", async () => {
    await render({ claimId: 12, claimName: "Voucher 14" }, true);

    expect(fixture.nativeElement.textContent.replace(/\s+/g, " ")).toContain("In claim: Voucher 14");
  });

  // Absent means "nothing to show": no claim, or one hidden from the viewer.
  it("renders nothing for a receipt without a claim", async () => {
    await render({}, true);

    expect(link()).toBeNull();
    expect(fixture.nativeElement.textContent.trim()).toBe("");
  });

  it("renders nothing with no receipt", async () => {
    await render(undefined);

    expect(link()).toBeNull();
  });
});
