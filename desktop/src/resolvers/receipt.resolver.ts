import { inject } from "@angular/core";
import { ResolveFn } from "@angular/router";
import { Observable } from "rxjs";
import { Receipt, ReceiptService } from "../open-api";

export const receiptResolverFn: ResolveFn<Observable<Receipt>> = (
  route,
  state
) => {
  const receiptService = inject(ReceiptService);
  return receiptService.getReceiptById(route.params["id"]);
};
