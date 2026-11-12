import { CurrencyPipe } from "@angular/common";
import { HttpErrorResponse } from "@angular/common/http";
import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute, convertToParamMap } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { Claim, ClaimList, Group, Permission, ReimbursementClaimService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { AuthState, GroupState, SetGroups, SetPermissions, SetSelectedGroupId } from "../../store";
import { SystemSettingsState } from "../../store/system-settings.state";
import { ClaimListComponent } from "./claim-list.component";

describe("ClaimListComponent", () => {
  let fixture: ComponentFixture<ClaimListComponent>;
  let component: ClaimListComponent;
  let store: Store;

  const claimService = { getClaimsForGroup: jest.fn() };

  const household = { id: 5, name: "Household", isAllGroup: false } as Group;
  const work = { id: 6, name: "Work", isAllGroup: false } as Group;
  const all = { id: 9, name: "All", isAllGroup: true } as Group;

  const claim = {
    id: 12,
    groupId: 5,
    name: "Voucher 14",
    status: "SUBMITTED",
    receiptIds: [],
    ledgerEntries: [],
  } as unknown as Claim;

  const list: ClaimList = {
    claims: [claim],
    summary: {
      asOfDate: "2026-11-11",
      claimCount: 3,
      awaitingCount: 2,
      outstandingTotal: "1250.40",
      oldestDaysWaiting: 47,
    },
  };

  const setUp = async (groupId: number, groups: Group[], permissions: { [groupId: number]: string[] }) => {
    await TestBed.configureTestingModule({
      declarations: [ClaimListComponent],
      imports: [PipesModule, NgxsModule.forRoot([AuthState, GroupState, SystemSettingsState])],
      providers: [
        provideZonelessChangeDetection(),
        CurrencyPipe,
        { provide: ReimbursementClaimService, useValue: claimService },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({ groupId: String(groupId) })),
            snapshot: { paramMap: convertToParamMap({ groupId: String(groupId) }) },
          },
        },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    store = TestBed.inject(Store);
    store.dispatch(new SetGroups(groups));
    store.dispatch(new SetSelectedGroupId(String(groupId)));
    store.dispatch(new SetPermissions([], permissions));

    fixture = TestBed.createComponent(ClaimListComponent);
    component = fixture.componentInstance;
    fixture.autoDetectChanges();
    await fixture.whenStable();
  };

  const text = (testId: string): string =>
    fixture.nativeElement.querySelector(`[data-testid='${testId}']`)?.textContent?.replace(/\s+/g, " ").trim();

  const value = (testId: string): string =>
    fixture.nativeElement
      .querySelector(`[data-testid='${testId}'] .claim-list-summary__value`)
      ?.textContent?.replace(/\s+/g, " ")
      .trim();

  beforeEach(() => {
    jest.resetAllMocks();
    claimService.getClaimsForGroup.mockReturnValue(of(list));
  });

  describe("on a real group", () => {
    beforeEach(async () => {
      await setUp(5, [all, household, work], { 5: [Permission.GroupReceiptsUpdate] });
    });

    it("loads the route group's claims", () => {
      expect(claimService.getClaimsForGroup).toHaveBeenCalledWith(5);
      expect(component.headerText()).toBe("Household Claims");
      expect(component.dataSource().data).toEqual([claim]);
    });

    it("shows the server's summary: awaiting, still owed and the longest wait", () => {
      expect(value("claim-list-summary-awaiting")).toBe("2 of 3 claims");
      expect(value("claim-list-summary-outstanding")).toBe("$1,250.40");
      expect(value("claim-list-summary-oldest")).toBe("47 days");
      expect(text("claim-list-summary")).toContain("As of Nov 11, 2026");
    });

    it("has no Group column", () => {
      expect(component.displayedColumns()).not.toContain("group");
      expect(component.displayedColumns()).toEqual([
        "name",
        "referenceNumber",
        "type",
        "status",
        "submittedDate",
        "dueDate",
        "claimed",
        "paid",
        "outstanding",
        "daysWaiting",
      ]);
    });

    it("offers Add Claim when the group grants update", () => {
      expect(component.canAdd()).toBe(true);
      expect(fixture.nativeElement.querySelector("[data-testid='claim-add']")).not.toBeNull();
    });
  });

  it("withholds Add Claim from a viewer who cannot update the group", async () => {
    await setUp(6, [all, household, work], { 6: [Permission.GroupReceiptsRead] });

    expect(component.canAdd()).toBe(false);
    expect(fixture.nativeElement.querySelector("[data-testid='claim-add']")).toBeNull();
  });

  describe("on the All group", () => {
    it("adds a Group column and offers no Add Claim when there is no single real group", async () => {
      // The All membership granting update must not count: a claim needs a real group.
      await setUp(9, [all, household, work], { 9: [Permission.GroupReceiptsUpdate], 5: [Permission.GroupReceiptsUpdate] });

      expect(claimService.getClaimsForGroup).toHaveBeenCalledWith(9);
      expect(component.displayedColumns()).toContain("group");
      expect(component.canAdd()).toBe(false);
      expect(component.subtitle()).toContain("Open a group to add a claim to it.");
    });

    // The same add target the receipt form seeds for a sole-group user.
    it("offers Add Claim for a member of exactly one real group", async () => {
      await setUp(9, [all, household], { 5: [Permission.GroupReceiptsUpdate] });

      expect(component.canAdd()).toBe(true);
    });
  });

  it("says so when there are no claims yet", async () => {
    claimService.getClaimsForGroup.mockReturnValue(
      of({ claims: [], summary: { asOfDate: "2026-11-11", claimCount: 0, awaitingCount: 0, outstandingTotal: "0.00" } })
    );
    await setUp(5, [all, household], { 5: [Permission.GroupReceiptsUpdate] });

    expect(text("claim-list-empty")).toContain("No claims yet.");
    expect(value("claim-list-summary-oldest")).toBe("—");
  });

  it("says so when the claims cannot be loaded", async () => {
    claimService.getClaimsForGroup.mockReturnValue(throwError(() => new HttpErrorResponse({ status: 403 })));
    await setUp(5, [all, household], { 5: [] });

    expect(text("claim-list-error")).toBe("Claims could not be loaded for this group.");
  });
});
