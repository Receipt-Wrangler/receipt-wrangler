import { HttpErrorResponse } from "@angular/common/http";
import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { ActivatedRouteSnapshot, convertToParamMap, provideRouter, RedirectCommand, RouterStateSnapshot } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { firstValueFrom, Observable, of, throwError } from "rxjs";
import { FormMode } from "../enums/form-mode.enum";
import { Claim, Permission, ReimbursementClaimService } from "../open-api";
import { SnackbarService } from "../services/snackbar.service";
import { AuthState, GroupState, SetPermissions, SetSelectedGroupId } from "../store";
import { CLAIM_NOT_FOUND_MESSAGE, claimResolverFn } from "./claim.resolver";

describe("claimResolverFn", () => {
  let store: Store;
  const claimService = { getClaimById: jest.fn() };
  const snackbarService = { error: jest.fn() };

  const claim = { id: 12, groupId: 5, name: "Voucher 14" } as Claim;

  const resolve = (id: string, mode: FormMode) => {
    const route = {
      paramMap: convertToParamMap({ id }),
      data: { mode },
    } as unknown as ActivatedRouteSnapshot;

    return TestBed.runInInjectionContext(() =>
      firstValueFrom(claimResolverFn(route, {} as RouterStateSnapshot) as Observable<Claim | RedirectCommand>)
    );
  };

  const redirectUrl = (result: Claim | RedirectCommand) =>
    (result as RedirectCommand).redirectTo.toString();

  beforeEach(() => {
    jest.resetAllMocks();
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState, GroupState])],
      providers: [
        provideZonelessChangeDetection(),
        provideRouter([]),
        { provide: ReimbursementClaimService, useValue: claimService },
        { provide: SnackbarService, useValue: snackbarService },
      ],
    });
    store = TestBed.inject(Store);
    store.dispatch(new SetSelectedGroupId("99"));
  });

  it("resolves the claim for view", async () => {
    claimService.getClaimById.mockReturnValue(of(claim));

    expect(await resolve("12", FormMode.view)).toBe(claim);
    expect(claimService.getClaimById).toHaveBeenCalledWith(12);
  });

  it("resolves the claim for edit when the claim's own group grants update", async () => {
    store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
    claimService.getClaimById.mockReturnValue(of(claim));

    expect(await resolve("12", FormMode.edit)).toBe(claim);
  });

  // The browsed group (99, e.g. All) granting update must not count.
  it("falls back to the view when the claim's group does not grant update", async () => {
    store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsRead], 99: [Permission.GroupReceiptsUpdate] }));
    claimService.getClaimById.mockReturnValue(of(claim));

    const result = await resolve("12", FormMode.edit);

    expect(result).toBeInstanceOf(RedirectCommand);
    expect(redirectUrl(result)).toBe("/claims/12/view");
    expect(snackbarService.error).not.toHaveBeenCalled();
  });

  // A missing claim and a forbidden one are the same 403.
  it("sends a missing or forbidden claim back to the list with a message", async () => {
    claimService.getClaimById.mockReturnValue(
      throwError(() => new HttpErrorResponse({ status: 403 }))
    );

    const result = await resolve("12", FormMode.view);

    expect(result).toBeInstanceOf(RedirectCommand);
    expect(redirectUrl(result)).toBe("/claims/group/99");
    expect(snackbarService.error).toHaveBeenCalledWith(CLAIM_NOT_FOUND_MESSAGE);
  });

  it("does not ask the server about an id that is not one", async () => {
    const result = await resolve("abc", FormMode.view);

    expect(result).toBeInstanceOf(RedirectCommand);
    expect(claimService.getClaimById).not.toHaveBeenCalled();
    expect(snackbarService.error).toHaveBeenCalledWith(CLAIM_NOT_FOUND_MESSAGE);
  });
});
