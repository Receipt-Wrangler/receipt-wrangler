import { Component, computed, input } from "@angular/core";
import { RouterLink } from "@angular/router";
import { ReceiptClaimRef } from "../claim.utils";

/**
 * The claim a receipt is in, as a link to the claim's page: the receipts
 * table's Claim column and the receipt form's audit details both render it.
 * Renders nothing when the receipt carries no claim, which the server also
 * sends for a claim the viewer may not see.
 */
@Component({
  selector: "app-receipt-claim-link",
  templateUrl: "./receipt-claim-link.component.html",
  standalone: true,
  imports: [RouterLink],
})
export class ReceiptClaimLinkComponent {
  public readonly receipt = input<ReceiptClaimRef | undefined>(undefined);

  /** Prefixes "In claim:", for where the link stands on its own line. */
  public readonly labelled = input<boolean>(false);

  public readonly claim = computed(() => {
    const receipt = this.receipt();
    return receipt?.claimId && receipt.claimName
      ? { id: receipt.claimId, name: receipt.claimName }
      : undefined;
  });
}
