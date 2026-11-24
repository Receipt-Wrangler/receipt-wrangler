import { ErrorHandler, Provider } from "@angular/core";

// Providers added to every spec's TestBed environment (angular.json `providersFile`).
//
// Rethrow instead of logging, so an error reported to Angular's ErrorHandler fails the test that
// caused it. TestBed already rethrows errors from change detection and listeners, but code that
// calls ErrorHandler.handleError directly (NGXS's unhandled action errors, for one) would otherwise
// only reach console.error. This is the behavior the specs were written against under
// jest-preset-angular.
const providers: Provider[] = [
  {
    provide: ErrorHandler,
    useValue: {
      handleError: (error: unknown) => {
        throw error;
      },
    },
  },
];

export default providers;
