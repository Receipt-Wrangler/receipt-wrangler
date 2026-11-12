import { HttpErrorResponse } from "@angular/common/http";
import { FormControl, FormGroup } from "@angular/forms";
import { Permission } from "../open-api";
import {
  applyClaimApiErrors,
  buildUpsertClaimCommand,
  buildUpsertClaimLedgerEntryCommand,
  canReadClaims,
  canUpdateClaims,
  claimedReceiptsByClaim,
  claimGroupMoveWarning,
  CLAIM_LEDGER_TYPE_OPTIONS,
  CLAIM_STATUS_OPTIONS,
  claimLedgerTypeOption,
  claimLedgerTypeOptionsFor,
  claimStatusLabel,
  claimStatusOptionsFor,
  claimStatusTone,
  isInOtherClaim,
  optionalCalendarDay,
  receiptGroupIds,
  unmatchedApiErrorMessage,
} from "./claim.utils";

describe("claim.utils", () => {
  describe("statuses", () => {
    it("labels every known status", () => {
      expect(CLAIM_STATUS_OPTIONS.map((option) => claimStatusLabel(option.value))).toEqual([
        "Draft",
        "Submitted",
        "Partially Paid",
        "Paid",
        "Disputed",
        "Closed",
      ]);
    });

    // The wire field is a plain string, so a status added on the server later
    // must render rather than break the page.
    it("shows the raw code for a status it does not know", () => {
      expect(claimStatusLabel("ESCALATED")).toBe("ESCALATED");
      expect(claimStatusTone("ESCALATED")).toBeUndefined();
    });

    it("colours the statuses that need attention", () => {
      expect(claimStatusTone("DISPUTED")).toBe("red");
      expect(claimStatusTone("PAID")).toBe("green");
      expect(claimStatusTone("SUBMITTED")).toBe("yellow");
      expect(claimStatusTone("DRAFT")).toBeUndefined();
    });

    it("offers an unknown current status so the select can still show it", () => {
      expect(claimStatusOptionsFor("PAID")).toBe(CLAIM_STATUS_OPTIONS);
      expect(claimStatusOptionsFor(undefined)).toBe(CLAIM_STATUS_OPTIONS);
      expect(claimStatusOptionsFor("ESCALATED").at(-1)).toEqual({
        value: "ESCALATED",
        displayValue: "ESCALATED",
      });
    });
  });

  describe("ledger types", () => {
    it("falls back to the raw code for an unknown type", () => {
      expect(claimLedgerTypeOption("PAYMENT").displayValue).toBe("Payment received");
      expect(claimLedgerTypeOption("ADJUSTMENT")).toEqual({
        value: "ADJUSTMENT",
        displayValue: "ADJUSTMENT",
        icon: "help_outline",
      });
    });

    it("gives a payment and a recoupment different icons", () => {
      expect(claimLedgerTypeOption("PAYMENT").icon).not.toBe(claimLedgerTypeOption("RECOUPMENT").icon);
    });

    it("offers an unknown current type alongside the known ones", () => {
      expect(claimLedgerTypeOptionsFor("NOTE")).toBe(CLAIM_LEDGER_TYPE_OPTIONS);
      expect(claimLedgerTypeOptionsFor("ADJUSTMENT").map((option) => option.value)).toEqual([
        "PAYMENT",
        "RECOUPMENT",
        "NOTE",
        "ADJUSTMENT",
      ]);
    });
  });

  describe("permissions", () => {
    const permissions = {
      5: [Permission.GroupReceiptsRead, Permission.GroupReceiptsUpdate],
      6: [Permission.GroupReceiptsRead],
    };

    it("judges each group on its own", () => {
      expect(canUpdateClaims(permissions, 5)).toBe(true);
      expect(canUpdateClaims(permissions, 6)).toBe(false);
      expect(canReadClaims(permissions, 6)).toBe(true);
      expect(canReadClaims(permissions, 7)).toBe(false);
    });

    it("grants nothing without a group", () => {
      expect(canUpdateClaims(permissions, undefined)).toBe(false);
      expect(canUpdateClaims(permissions, 0)).toBe(false);
    });
  });

  describe("buildUpsertClaimCommand", () => {
    it("sends calendar days as bare yyyy-MM-dd, read on the picked day", () => {
      const command = buildUpsertClaimCommand({
        name: "  Voucher 14  ",
        referenceNumber: " TA-2026-14 ",
        type: "Travel voucher",
        status: "SUBMITTED",
        // The datepicker's local midnight of Sep 30.
        submittedDate: new Date(2026, 8, 30),
        dueDate: "2026-10-31",
        claimedAmount: 412.5,
        notes: "Sent by email",
      });

      expect(command).toEqual({
        name: "Voucher 14",
        referenceNumber: "TA-2026-14",
        type: "Travel voucher",
        status: "SUBMITTED",
        submittedDate: "2026-09-30",
        dueDate: "2026-10-31",
        claimedAmount: "412.5",
        notes: "Sent by email",
      });
    });

    // PUT is a full replace, so an emptied field must go out absent, which clears it.
    it("omits every emptied optional field", () => {
      const command = buildUpsertClaimCommand({
        name: "Voucher",
        referenceNumber: "   ",
        type: null,
        status: "DRAFT",
        submittedDate: null,
        dueDate: "",
        claimedAmount: "",
        notes: "",
      });

      expect(command).toEqual({ name: "Voucher", status: "DRAFT" });
      expect(JSON.parse(JSON.stringify(command))).toEqual({ name: "Voucher", status: "DRAFT" });
    });

    it("keeps a claimed amount of zero", () => {
      expect(
        buildUpsertClaimCommand({
          name: "Voucher",
          referenceNumber: null,
          type: null,
          status: "DRAFT",
          submittedDate: null,
          dueDate: null,
          claimedAmount: 0,
          notes: null,
        }).claimedAmount
      ).toBe("0");
    });
  });

  describe("buildUpsertClaimLedgerEntryCommand", () => {
    it("builds a payment", () => {
      expect(
        buildUpsertClaimLedgerEntryCommand({
          type: "PAYMENT",
          date: new Date(2026, 9, 1),
          amount: "300.00",
          note: "",
        })
      ).toEqual({ type: "PAYMENT", date: "2026-10-01", amount: "300.00" });
    });

    it("builds a note without an amount", () => {
      expect(
        buildUpsertClaimLedgerEntryCommand({
          type: "NOTE",
          date: "2026-10-02",
          amount: "",
          note: " Called the travel office ",
        })
      ).toEqual({ type: "NOTE", date: "2026-10-02", note: "Called the travel office" });
    });
  });

  describe("optionalCalendarDay", () => {
    it("is undefined for no day", () => {
      expect(optionalCalendarDay(null)).toBeUndefined();
      expect(optionalCalendarDay(undefined)).toBeUndefined();
      expect(optionalCalendarDay("")).toBeUndefined();
    });
  });

  describe("server errors", () => {
    const form = () =>
      new FormGroup({
        name: new FormControl("x"),
        date: new FormControl(null),
      });

    const badRequest = (body: unknown) =>
      new HttpErrorResponse({ status: 400, error: body });

    it("puts a field-keyed 400 on the matching control and marks it touched", () => {
      const group = form();

      const applied = applyClaimApiErrors(group, badRequest({ name: "Name is required" }));

      expect(applied).toBe(true);
      expect(group.get("name")?.errors).toEqual({ error: "Name is required" });
      expect(group.get("name")?.touched).toBe(true);
      expect(group.get("date")?.errors).toBeNull();
    });

    // The interceptor already toasts an errorMsg.
    it("leaves an errorMsg body to the interceptor", () => {
      const group = form();

      expect(applyClaimApiErrors(group, badRequest({ errorMsg: "The request body is not valid" }))).toBe(false);
      expect(unmatchedApiErrorMessage(group, badRequest({ errorMsg: "x" }))).toBeUndefined();
      expect(group.get("name")?.errors).toBeNull();
    });

    it("ignores anything that is not a 400", () => {
      const group = form();
      const forbidden = new HttpErrorResponse({ status: 403, error: { name: "nope" } });

      expect(applyClaimApiErrors(group, forbidden)).toBe(false);
      expect(applyClaimApiErrors(group, new Error("boom"))).toBe(false);
      expect(group.get("name")?.errors).toBeNull();
    });

    it("reports a field-keyed message no control can show", () => {
      const group = form();
      const error = badRequest({ receiptIds: "Between 1 and 500 receipt ids are required" });

      expect(applyClaimApiErrors(group, error)).toBe(false);
      expect(unmatchedApiErrorMessage(group, error)).toBe("Between 1 and 500 receipt ids are required");
    });
  });

  it("lists the distinct groups of a selection in first-seen order", () => {
    expect(receiptGroupIds([{ groupId: 6 }, { groupId: 5 }, { groupId: 6 }])).toEqual([6, 5]);
  });

  describe("a receipt's own claim", () => {
    it("is another claim only when it is visible and not the one given", () => {
      expect(isInOtherClaim({ claimId: 40, claimName: "Voucher 13" }, 12)).toBe(true);
      expect(isInOtherClaim({ claimId: 12, claimName: "Voucher 14" }, 12)).toBe(false);
      expect(isInOtherClaim({ claimId: 40, claimName: "Voucher 13" }, undefined)).toBe(true);
      // Absent: in no claim, or in one hidden from the viewer.
      expect(isInOtherClaim({}, 12)).toBe(false);
    });

    it("groups the claimed receipts of a selection by claim, in first-seen order", () => {
      const receipts = [
        { name: "Taxi" },
        { name: "Fuel", claimId: 40, claimName: "Voucher 13" },
        { name: "Parking", claimId: 39, claimName: "Voucher 12" },
        { name: "Toll", claimId: 40, claimName: "Voucher 13" },
        { name: "Hotel", claimId: 12, claimName: "Voucher 14" },
      ];

      expect(claimedReceiptsByClaim(receipts, 12)).toEqual([
        { claimId: 40, claimName: "Voucher 13", receiptNames: ["Fuel", "Toll"] },
        { claimId: 39, claimName: "Voucher 12", receiptNames: ["Parking"] },
      ]);
      expect(claimedReceiptsByClaim(receipts).map((group) => group.claimId)).toEqual([40, 39, 12]);
    });

    it("warns about a group move only once a different group is picked", () => {
      const receipt = { groupId: 7, claimId: 12, claimName: "Voucher 14" };

      expect(claimGroupMoveWarning(receipt, 8)).toBe(
        "Moving this receipt to another group will remove it from claim 'Voucher 14'."
      );
      // The group control holds a number or a string, depending on how it was set.
      expect(claimGroupMoveWarning(receipt, "8")).toBeDefined();
      expect(claimGroupMoveWarning(receipt, 7)).toBeUndefined();
      expect(claimGroupMoveWarning(receipt, "7")).toBeUndefined();
      expect(claimGroupMoveWarning(receipt, "")).toBeUndefined();
      expect(claimGroupMoveWarning(receipt, null)).toBeUndefined();
      expect(claimGroupMoveWarning({ groupId: 7 }, 8)).toBeUndefined();
      expect(claimGroupMoveWarning(undefined, 8)).toBeUndefined();
    });
  });
});
