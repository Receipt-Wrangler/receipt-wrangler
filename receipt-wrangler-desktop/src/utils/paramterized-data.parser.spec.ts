import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { GroupState, UserState } from "../store";
import { ParameterizedDataParser } from "./paramterterized-data-parser";

describe("ParameterizedDataParser", () => {
  let parameterizedDataParser: ParameterizedDataParser;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [NgxsModule.forRoot([GroupState, UserState])],
      providers: [],
    }).compileComponents();

    parameterizedDataParser = TestBed.inject(ParameterizedDataParser);
    store = TestBed.inject(Store);
  });

  it("create an instance", () => {
    expect(parameterizedDataParser).toBeTruthy();
  });

  it("should resolve userId", () => {
    const user = { id: "1", name: "John" };
    store.reset({
      users: {
        users: [user],
      },
    });

    const result = parameterizedDataParser.parse("${userId:1.name}");

    expect(result).toBe("John");
  });

  it("should resolve groupId and display it if it is type string", () => {
    const group = { id: "1", name: "Group A" };
    store.reset({
      groups: {
        groups: [group],
      },
    });

    const result = parameterizedDataParser.parse("${groupId:1.name:string}");

    expect(result).toBe("Group A");
  });

  it("should resolve groupId and not display it if it is type link", () => {
    const group = { id: "1", name: "Group A" };
    store.reset({
      groups: {
        groups: [group],
      },
    });

    const result = parameterizedDataParser.parse("${groupId:1.name:link}");

    expect(result).toBe("");
    expect(parameterizedDataParser.link).toEqual("/receipts/group/1");
  });

  it("should resolve receiptId and not display it if it is type link", () => {
    const result = parameterizedDataParser.parse("${receiptId:1.noop:link}");

    expect(result).toBe("");
    expect(parameterizedDataParser.link).toEqual("/receipts/1/view");
  });

  it("should resolve multiple pieces of paramterized data", () => {
    const group = { id: "1", name: "Hello" };
    const group2 = { id: "2", name: "World" };
    store.reset({
      groups: {
        groups: [group, group2],
      },
    });

    const result = parameterizedDataParser.parse(
      "${groupId:1.name:string} ${groupId:2.name:string}"
    );

    expect(result).toBe("Hello World");
  });

  it("should return empty string when id does not exist", () => {
    const result = parameterizedDataParser.parse("${userId:999.name}");

    expect(result).toBe("");
  });
});
