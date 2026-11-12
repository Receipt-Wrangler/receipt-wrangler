import { inject } from "@angular/core";
import { RedirectCommand, ResolveFn, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { catchError, map, of } from "rxjs";
import { FormMode } from "../enums/form-mode.enum";
import { Claim, ReimbursementClaimService } from "../open-api";
import { SnackbarService } from "../services/snackbar.service";
import { AuthState, GroupState } from "../store";
import { canUpdateClaims } from "./claim.utils";

export const CLAIM_NOT_FOUND_MESSAGE = "That claim could not be found, or you do not have access to it.";

/**
 * Loads the claim for `/claims/:id/view|edit`, and is also where those routes are
 * gated, because the claim's group — which every permission is judged in — is
 * only known once the claim has loaded.
 *
 * - A missing claim and a forbidden one both answer 403 with the same message
 *   (the server never reveals that a claim exists). Either way the user is sent
 *   to the claims list of the group they are browsing, with a toast: the HTTP
 *   interceptor deliberately stays quiet for a GET, so without one the page
 *   would just never appear.
 * - Edit mode without `group.receipts.update` in the claim's group falls back
 *   to the read-only view rather than rendering a form whose save can only 403.
 */
export const claimResolverFn: ResolveFn<Claim> = (route) => {
  const claimService = inject(ReimbursementClaimService);
  const router = inject(Router);
  const snackbarService = inject(SnackbarService);
  const store = inject(Store);

  const claimId = Number(route.paramMap.get("id"));
  const listUrl = () =>
    router.parseUrl(`/claims/group/${store.selectSnapshot(GroupState.selectedGroupId)}`);

  if (!Number.isInteger(claimId) || claimId <= 0) {
    snackbarService.error(CLAIM_NOT_FOUND_MESSAGE);
    return of(new RedirectCommand(listUrl()));
  }

  return claimService.getClaimById(claimId).pipe(
    map((claim) => {
      const wantsEdit = route.data["mode"] === FormMode.edit;
      const groupPermissions = store.selectSnapshot(AuthState.groupPermissions);
      if (wantsEdit && !canUpdateClaims(groupPermissions, claim.groupId)) {
        return new RedirectCommand(router.parseUrl(`/claims/${claim.id}/view`));
      }

      return claim;
    }),
    catchError(() => {
      snackbarService.error(CLAIM_NOT_FOUND_MESSAGE);
      return of(new RedirectCommand(listUrl()));
    })
  );
};
