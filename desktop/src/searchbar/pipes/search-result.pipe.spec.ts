import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { SearchResult } from "../../open-api";
import { GroupState } from "../../store";
import { SystemSettingsState } from "../../store/system-settings.state";
import { SearchResultPipe } from "./search-result.pipe";

describe("SearchResultPipe", () => {
  let pipe: SearchResultPipe;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SearchResultPipe],
      imports: [NgxsModule.forRoot([GroupState, SystemSettingsState])],
      providers: [SearchResultPipe],
    }).compileComponents();

    pipe = TestBed.inject(SearchResultPipe);
    store = TestBed.inject(Store);
  });

  it("should create the pipe", () => {
    expect(pipe).toBeTruthy();
  });

  it("should return search result string", () => {
    store.reset({
      groups: {
        groups: [
          {
            id: 1,
            name: "group name",
            isDefault: true,
            groupMembers: [],
          },
        ],
      },
    });

    const searchResult: SearchResult = {
      id: 1,
      name: "my result",
      type: "receipt",
      groupId: 1,
      date: "2022-12-12",
      createdAt: "2022-12-12",
    };

    const result = pipe.transform(searchResult);

    expect(result).toEqual("Dec 12, 2022 - my result (group name)");
  });

  it("shows a stored midnight-UTC receipt date as its own calendar day", () => {
    store.reset({
      groups: { groups: [{ id: 1, name: "group name", isDefault: true, groupMembers: [] }] },
      systemSettings: { timeZone: "Pacific/Honolulu" },
    });

    const result = pipe.transform({
      id: 1,
      name: "my result",
      type: "receipt",
      groupId: 1,
      date: "2026-09-30T00:00:00Z",
      createdAt: "2026-09-30T00:00:00Z",
    });

    expect(result).toEqual("Sep 30, 2026 - my result (group name)");
  });
});
