import { FilterOperation } from "../open-api";

export const listOperationOptions = Object.values(FilterOperation).filter(
  (k) => k === "CONTAINS" && k !== FilterOperation.WithinCurrentMonth && !!k
);

export const dateOperationOptions = Object.values(FilterOperation).filter(
  (k) => !k.includes("CONTAINS") && !!k
);

export const numberOperationOptions = Object.values(FilterOperation).filter(
  (k) => !k.includes("CONTAINS") && k !== FilterOperation.WithinCurrentMonth && !!k
);

export const textOperationOptions = Object.values(FilterOperation).filter(
  (k) => !k.includes("THAN") && k !== FilterOperation.WithinCurrentMonth && !k.includes("BETWEEN") && !!k
);

export const usersOperationOptions = Object.values(FilterOperation).filter(
  (k) => k === "CONTAINS" && k !== FilterOperation.WithinCurrentMonth && !!k
);
