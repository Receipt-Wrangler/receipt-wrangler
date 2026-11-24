import { defineConfig } from "vitest/config";

// Runner settings for `ng test` (`@angular/build:unit-test`, see angular.json). Only what the
// builder cannot express lives here; everything else, coverage included, is in angular.json.
export default defineConfig({
  test: {
    // The builder defaults to sharing one module graph across spec files. Specs here patch globals
    // (fake timers, navigator.clipboard, localStorage via the NGXS storage plugin), so each file
    // gets a fresh worker, as it did under Jest.
    isolate: true,
    // The builder silently picks happy-dom whenever it can resolve it; the suite targets jsdom.
    environment: "jsdom",
    // Node 25+ ships its own localStorage global, which Vitest's jsdom environment leaves in place
    // of jsdom's, so specs (and the NGXS storage plugin) would get Node's, unusable without a
    // backing file. Turn Node's off in the test workers. Older Nodes have no such global, and
    // Node 20 would reject the flag, so it is only passed when it is needed.
    execArgv: "localStorage" in globalThis ? ["--no-experimental-webstorage"] : [],
    coverage: {
      // CI's coverage summary and Codecov steps read cobertura-coverage.xml from this folder.
      // angular.json has no option for it, and the builder's default is coverage/<project>.
      reportsDirectory: "coverage/receipt-wrangler-desktop",
    },
  },
});
