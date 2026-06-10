import { hasAll, hasAny, matches } from "./permission.utils";

describe("permission.utils", () => {
  describe("matches", () => {
    const cases: { name: string; granted: string; required: string; want: boolean }[] = [
      { name: "exact match", granted: "app.users.read", required: "app.users.read", want: true },
      { name: "exact mismatch", granted: "app.users.read", required: "app.users.create", want: false },
      { name: "global wildcard", granted: "*", required: "group.receipts.read", want: true },
      { name: "scope wildcard matches deeper", granted: "group.*", required: "group.receipts.read", want: true },
      { name: "scope wildcard matches single", granted: "group.*", required: "group.view", want: true },
      { name: "scope wildcard does not cross scope", granted: "group.*", required: "app.users.read", want: false },
      { name: "scope wildcard does not match bare prefix", granted: "group.*", required: "group", want: false },
      { name: "domain wildcard matches", granted: "group.receipts.*", required: "group.receipts.read", want: true },
      { name: "domain wildcard does not match sibling", granted: "group.receipts.*", required: "group.dashboards.read", want: false },
      { name: "domain wildcard covers sub-scope suffix", granted: "group.receipts.*", required: "group.receipts.read:any", want: true },
      { name: "mid wildcard matches one segment", granted: "group.*.read", required: "group.receipts.read", want: true },
      { name: "mid wildcard requires the trailing literal", granted: "group.*.read", required: "group.receipts.create", want: false },
      { name: "mid wildcard matches exactly one segment", granted: "group.*.read", required: "group.a.b.read", want: false },
      { name: "granted more specific than required", granted: "group.receipts.read", required: "group.receipts", want: false },
      { name: "required more specific than granted", granted: "group.receipts", required: "group.receipts.read", want: false },
      { name: "literal suffix exact", granted: "group.receipts.read:any", required: "group.receipts.read:any", want: true },
      { name: "literal suffix no superset", granted: "group.receipts.read:any", required: "group.receipts.read:paid_by", want: false },
    ];

    cases.forEach(({ name, granted, required, want }) => {
      it(`${name}: matches("${granted}", "${required}") === ${want}`, () => {
        expect(matches(granted, required)).toBe(want);
      });
    });
  });

  describe("hasAll", () => {
    const cases: { name: string; granted: string[]; required: string[]; want: boolean }[] = [
      { name: "single granted", granted: ["app.users.read"], required: ["app.users.read"], want: true },
      { name: "single not granted", granted: ["app.users.read"], required: ["app.users.create"], want: false },
      { name: "all granted", granted: ["app.users.read", "app.users.create"], required: ["app.users.read", "app.users.create"], want: true },
      { name: "one missing fails AND", granted: ["app.users.read"], required: ["app.users.read", "app.users.create"], want: false },
      { name: "wildcard grant satisfies all", granted: ["group.*"], required: ["group.receipts.read", "group.view"], want: true },
      { name: "empty granted denies", granted: [], required: ["app.users.read"], want: false },
      { name: "no required denies", granted: ["app.users.read"], required: [], want: false },
    ];

    cases.forEach(({ name, granted, required, want }) => {
      it(name, () => {
        expect(hasAll(granted, ...required)).toBe(want);
      });
    });
  });

  describe("hasAny", () => {
    const cases: { name: string; granted: string[]; required: string[]; want: boolean }[] = [
      { name: "one of many granted", granted: ["app.users.read"], required: ["app.users.create", "app.users.read"], want: true },
      { name: "none granted", granted: ["app.users.read"], required: ["app.users.create", "app.users.delete"], want: false },
      { name: "wildcard satisfies one", granted: ["group.receipts.*"], required: ["group.view", "group.receipts.read"], want: true },
      { name: "empty granted denies", granted: [], required: ["app.users.read"], want: false },
      { name: "no required denies", granted: ["app.users.read"], required: [], want: false },
    ];

    cases.forEach(({ name, granted, required, want }) => {
      it(name, () => {
        expect(hasAny(granted, ...required)).toBe(want);
      });
    });
  });
});
