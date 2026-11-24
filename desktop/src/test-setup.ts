// Runs before every spec file, after the builder has initialized TestBed.

// jsdom does not implement HTMLCanvasElement.getContext; ng2-charts calls it
// during component initialization, which spams console.error in otherwise
// passing test runs. Returning null matches the no-op behavior the charts
// already tolerate under jsdom.
HTMLCanvasElement.prototype.getContext = vi.fn(() => null) as any;
