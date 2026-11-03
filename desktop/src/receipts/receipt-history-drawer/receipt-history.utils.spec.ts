import { ReceiptHistoryChange, ReceiptHistoryEntry, User } from "../../open-api";
import {
  DEFAULT_ACTION_VERB,
  EMPTY_VALUE_TEXT,
  HistoryValueFormatters,
  UNKNOWN_USER_NAME,
  UserLookup,
  buildHistoryChangeRow,
  formatHistoryValue,
  groupHistoryByDay,
  historyActionVerb,
  historyActorDisplay,
  historyBeforeStateNotice,
  historyFieldLabel,
  historySnapshotsJson,
  historySummary,
  SAVE_BAR_MIN_WIDTH,
  saveBarInsetBesideDrawer,
  titleCaseCode,
} from "./receipt-history.utils";

const maya = { id: 2, displayName: "Maya Chen", username: "maya" } as User;
const findUser: UserLookup = (id) => (id === 2 ? maya : undefined);

const formatters: HistoryValueFormatters = {
  currency: (value) => `$${value}`,
  calendarDate: (value) => `day:${value}`,
  instant: (value) => `instant:${value}`,
  findUser,
};

function entry(overrides: Partial<ReceiptHistoryEntry> = {}): ReceiptHistoryEntry {
  return {
    systemTaskId: 1,
    action: "UPDATED",
    recordedAt: "2026-10-13T10:02:51Z",
    actor: { type: "USER", userId: 2 },
    hasSnapshots: true,
    changes: [],
    ...overrides,
  };
}

describe("receipt history utils", () => {
  describe("historyActionVerb", () => {
    it.each([
      ["CREATED", "created this receipt"],
      ["UPDATED", "edited this receipt"],
      ["STATUS_CHANGED", "changed the status"],
      ["RESOLVED", "resolved this receipt"],
      ["IMAGES_COMMENTS_UPDATED", "updated images and comments"],
      ["FAILED", "tried to save changes"],
    ])("words %s", (action, verb) => {
      expect(historyActionVerb(action)).toBe(verb);
    });

    it("falls back for an action this client does not know", () => {
      expect(historyActionVerb("ARCHIVED")).toBe(DEFAULT_ACTION_VERB);
      expect(historyActionVerb(undefined)).toBe(DEFAULT_ACTION_VERB);
    });
  });

  describe("historyFieldLabel", () => {
    const change = (field: string, extra: Partial<ReceiptHistoryChange> = {}): ReceiptHistoryChange => ({
      field,
      kind: "CHANGED",
      ...extra,
    });

    it.each([
      ["NAME", "Name"],
      ["AMOUNT", "Amount"],
      ["DATE", "Receipt Date"],
      ["GROUP", "Group"],
      ["PAID_BY", "Paid By"],
      ["STATUS", "Status"],
      ["RESOLVED_DATE", "Resolved Date"],
      ["CATEGORIES", "Categories"],
      ["TAGS", "Tags"],
      ["COMMENT", "Comment"],
      ["IMAGE", "Image"],
      ["EMAIL_INBOX", "Inbox"],
      ["EMAIL_SENDER", "Sender"],
    ])("labels %s", (field, label) => {
      expect(historyFieldLabel(change(field), findUser)).toBe(label);
    });

    it("names the subject of item, share and custom field rows", () => {
      expect(historyFieldLabel(change("CUSTOM_FIELD", { label: "Project" }), findUser)).toBe("Project");
      expect(historyFieldLabel(change("ITEM", { label: "Produce" }), findUser)).toBe("Item: Produce");
      expect(historyFieldLabel(change("ITEM_CATEGORIES", { label: "Produce" }), findUser)).toBe(
        "Item categories: Produce"
      );
      expect(historyFieldLabel(change("ITEM_TAGS", { label: "Produce" }), findUser)).toBe("Item tags: Produce");
      expect(historyFieldLabel(change("SHARE", { userId: 2 }), findUser)).toBe("Share: Maya Chen");
    });

    it("degrades when the subject is missing or unknown", () => {
      expect(historyFieldLabel(change("ITEM"), findUser)).toBe("Item");
      expect(historyFieldLabel(change("CUSTOM_FIELD"), findUser)).toBe("Custom field");
      expect(historyFieldLabel(change("SHARE", { userId: 99 }), findUser)).toBe(`Share: ${UNKNOWN_USER_NAME}`);
    });

    it("title-cases a field code it does not know", () => {
      expect(historyFieldLabel(change("RECEIPT_NOTES"), findUser)).toBe("Receipt Notes");
      expect(titleCaseCode("")).toBe("");
    });
  });

  describe("historyActorDisplay", () => {
    it("shows a known user with their own avatar", () => {
      expect(historyActorDisplay({ type: "USER", userId: 2 }, findUser)).toEqual({
        kind: "user",
        name: "Maya Chen",
        user: maya,
      });
    });

    it("shows an unknown user as such, with an icon", () => {
      const actor = historyActorDisplay({ type: "USER", userId: 99 }, findUser);
      expect(actor.name).toBe(UNKNOWN_USER_NAME);
      expect(actor.icon).toBe("person");
      expect(actor.user).toBeUndefined();
    });

    it.each([
      ["QUICK_SCAN", "Quick Scan", "photo_camera", "#009efa"],
      ["EMAIL", "Email Upload", "mail", "#64748b"],
      ["SYSTEM", "System", "settings", "#94a3b8"],
      ["ROBOT", "System", "settings", "#94a3b8"],
    ])("shows a %s actor", (type, name, icon, color) => {
      expect(historyActorDisplay({ type }, findUser)).toEqual(
        expect.objectContaining({ name, icon, color })
      );
    });

    it("treats a missing actor as the system", () => {
      expect(historyActorDisplay(undefined, findUser).name).toBe("System");
    });
  });

  describe("groupHistoryByDay", () => {
    const now = new Date("2026-11-03T12:00:00Z");

    it("groups by the app-zone day, keeping the server's newest-first order", () => {
      const entries = [
        entry({ systemTaskId: 3, recordedAt: "2026-10-13T18:00:00Z" }),
        entry({ systemTaskId: 2, recordedAt: "2026-10-13T09:00:00Z" }),
        entry({ systemTaskId: 1, recordedAt: "2026-10-12T09:14:08Z" }),
      ];

      const days = groupHistoryByDay(entries, "UTC", "en-US", now);

      expect(days.map((day) => day.label)).toEqual(["Tuesday, Oct 13", "Monday, Oct 12"]);
      expect(days[0].entries.map((e) => e.systemTaskId)).toEqual([3, 2]);
      expect(days[1].key).toBe("2026-10-12");
    });

    // 01:30 UTC on Oct 13 is still the evening of Oct 12 in New York.
    it("reads the day in the app zone, not UTC", () => {
      const days = groupHistoryByDay(
        [entry({ recordedAt: "2026-10-13T01:30:00Z" })],
        "America/New_York",
        "en-US",
        now
      );

      expect(days[0].label).toBe("Monday, Oct 12");
    });

    it("names the year for a day outside the current one", () => {
      const days = groupHistoryByDay([entry({ recordedAt: "2025-10-13T10:00:00Z" })], "UTC", "en-US", now);

      expect(days[0].label).toBe("Monday, Oct 13, 2025");
    });

    it("keeps an unparseable timestamp rather than dropping the entry", () => {
      const days = groupHistoryByDay([entry({ recordedAt: "not a date" })], "UTC", "en-US", now);

      expect(days).toHaveLength(1);
      expect(days[0].label).toBe("Unknown date");
    });
  });

  describe("historySummary", () => {
    it("pluralizes each count", () => {
      expect(historySummary({ changeCount: 9, saveCount: 6, contributorCount: 3 })).toBe(
        "9 changes · 6 saves · 3 contributors"
      );
    });

    it("uses the singular for one", () => {
      expect(historySummary({ changeCount: 1, saveCount: 1, contributorCount: 1 })).toBe(
        "1 change · 1 save · 1 contributor"
      );
    });
  });

  describe("historyBeforeStateNotice", () => {
    it("needs no notice for a complete or unknown state", () => {
      expect(historyBeforeStateNotice({ beforeState: "COMPLETE" }, "UTC", "en-US")).toBeNull();
      expect(historyBeforeStateNotice({}, "UTC", "en-US")).toBeNull();
      expect(historyBeforeStateNotice({ beforeState: "SOMETHING_NEW" }, "UTC", "en-US")).toBeNull();
    });

    it("names the previous update a rebuilt entry was compared against", () => {
      const notice = historyBeforeStateNotice(
        {
          beforeState: "REBUILT",
          beforeSourceType: "RECEIPT_UPDATED",
          beforeSourceRecordedAt: "2026-10-16T08:03:17Z",
        },
        "UTC",
        "en-US"
      );

      expect(notice).toEqual({
        state: "REBUILT",
        text:
          "Recorded by an older version — compared against the previous update (Oct 16, 2026). " +
          "Bulk status changes can also show here.",
      });
    });

    it("names the created copy, and drops a missing date", () => {
      const notice = historyBeforeStateNotice(
        { beforeState: "REBUILT", beforeSourceType: "RECEIPT_UPLOADED" },
        "UTC",
        "en-US"
      );

      expect(notice?.text).toBe(
        "Recorded by an older version — compared against the copy saved when it was created. " +
          "Bulk status changes can also show here."
      );
    });

    it("warns that an incomplete entry may be inaccurate", () => {
      expect(historyBeforeStateNotice({ beforeState: "INCOMPLETE" }, "UTC", "en-US")).toEqual({
        state: "INCOMPLETE",
        text: "Recorded by an older version with no earlier complete copy. Some changes may be inaccurate.",
      });
    });
  });

  describe("formatHistoryValue", () => {
    it.each([
      [{ type: "TEXT", value: "Costco" }, "Costco"],
      [{ type: "CURRENCY", value: "187.42" }, "$187.42"],
      [{ type: "DATE", value: "2026-10-11" }, "day:2026-10-11"],
      [{ type: "INSTANT", value: "2026-10-20T16:52:31Z" }, "instant:2026-10-20T16:52:31Z"],
      [{ type: "USER", userId: 2 }, "Maya Chen"],
      [{ type: "USER", userId: 99 }, UNKNOWN_USER_NAME],
      [{ type: "STATUS", value: "NEEDS_ATTENTION" }, "Needs Attention"],
      [{ type: "BOOLEAN", value: "true" }, "Yes"],
      [{ type: "BOOLEAN", value: "false" }, "No"],
      [{ type: "LIST", value: "a, b" }, "a, b"],
    ])("formats %o", (value, text) => {
      expect(formatHistoryValue(value, formatters)).toBe(text);
    });

    it("leaves an unparseable date as the server sent it", () => {
      expect(
        formatHistoryValue({ type: "DATE", value: "soon" }, { ...formatters, calendarDate: () => null })
      ).toBe("soon");
    });
  });

  describe("buildHistoryChangeRow", () => {
    it("shows both sides of a change with the arrow between", () => {
      const row = buildHistoryChangeRow(
        { field: "NAME", kind: "CHANGED", from: { type: "TEXT", value: "A" }, to: { type: "TEXT", value: "B" } },
        formatters
      );

      expect(row).toEqual({
        field: "NAME",
        kind: "CHANGED",
        label: "Name",
        from: { text: "A" },
        to: { text: "B" },
        arrow: true,
      });
    });

    it("shows only the new value of an addition and the old value of a removal", () => {
      const both = { from: { type: "TEXT", value: "old" }, to: { type: "TEXT", value: "new" } };

      const added = buildHistoryChangeRow({ field: "TAGS", kind: "ADDED", ...both }, formatters);
      expect(added.from).toBeUndefined();
      expect(added.to?.text).toBe("new");
      expect(added.arrow).toBe(false);

      const removed = buildHistoryChangeRow({ field: "TAGS", kind: "REMOVED", ...both }, formatters);
      expect(removed.from?.text).toBe("old");
      expect(removed.to).toBeUndefined();
      expect(removed.arrow).toBe(false);
    });

    it("never renders a blank chip", () => {
      const row = buildHistoryChangeRow(
        { field: "CUSTOM_FIELD", label: "Notes", kind: "CHANGED", from: { type: "TEXT", value: "" }, to: { type: "TEXT", value: "x" } },
        formatters
      );

      expect(row.from?.text).toBe(EMPTY_VALUE_TEXT);
    });

    it("prefixes a comment with its author and clips long text, keeping it whole in the title", () => {
      const text = "Took the bag fee off, we brought our own. My share is the chicken + half the paper towels.";
      const row = buildHistoryChangeRow(
        { field: "COMMENT", kind: "ADDED", userId: 2, to: { type: "TEXT", value: text } },
        formatters
      );

      expect(row.to?.text).toBe(`Maya Chen: “${text.slice(0, 57)}…”`);
      expect(row.to?.title).toBe(`Maya Chen: “${text}”`);
    });

    it("leaves a short comment unclipped and untitled", () => {
      const row = buildHistoryChangeRow(
        { field: "COMMENT", kind: "REMOVED", from: { type: "TEXT", value: "Thanks!" } },
        formatters
      );

      expect(row.from).toEqual({ text: "“Thanks!”", title: undefined });
    });
  });

  describe("historySnapshotsJson", () => {
    it("prints the created receipt alone for a creation entry", () => {
      expect(historySnapshotsJson({ after: '{"name":"Costco"}' })).toBe('{\n  "name": "Costco"\n}');
    });

    it("prints both sides of a change", () => {
      const json = JSON.parse(historySnapshotsJson({ before: '{"name":"A"}', after: '{"name":"B"}' }));

      expect(json).toEqual({ before: { name: "A" }, after: { name: "B" } });
    });

    it("keeps text that is not JSON instead of failing", () => {
      expect(JSON.parse(historySnapshotsJson({ before: "{oops", after: '{"name":"B"}' })).before).toBe("{oops");
    });
  });

  describe("saveBarInsetBesideDrawer", () => {
    it("keeps the bar between the sidebar and the drawer", () => {
      expect(saveBarInsetBesideDrawer({ contentLeft: 82, width: 480, viewportWidth: 1440 })).toEqual({
        left: 82,
        right: 480,
      });
    });

    it("insets right up to the minimum width the bar needs", () => {
      const viewportWidth = 82 + 480 + SAVE_BAR_MIN_WIDTH;

      expect(saveBarInsetBesideDrawer({ contentLeft: 82, width: 480, viewportWidth })).toEqual({
        left: 82,
        right: 480,
      });
      expect(saveBarInsetBesideDrawer({ contentLeft: 82, width: 480, viewportWidth: viewportWidth - 1 })).toBeNull();
    });

    it("leaves the bar alone before the drawer has been measured", () => {
      expect(saveBarInsetBesideDrawer(null)).toBeNull();
      expect(saveBarInsetBesideDrawer({ contentLeft: 0, width: 0, viewportWidth: 1440 })).toBeNull();
    });
  });
});
