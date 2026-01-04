import { groupResolverFn } from "./group-resolver.service";

describe("GroupResolverService", () => {
  it("should export groupResolverFn", () => {
    expect(groupResolverFn).toBeDefined();
    expect(typeof groupResolverFn).toBe('function');
  });
});
