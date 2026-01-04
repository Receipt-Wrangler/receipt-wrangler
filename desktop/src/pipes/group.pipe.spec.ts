import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { Group } from "../open-api";
import { GroupState } from "../store";
import { GroupPipe } from "./group.pipe";

describe("GroupPipe", () => {
  let pipe: GroupPipe;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [GroupPipe],
      imports: [NgxsModule.forRoot([GroupState])],
      providers: [GroupPipe],
    }).compileComponents();

    const groups: Group[] = [
      {
        id: 1,
      } as Group,
      { id: 2 } as Group,
    ];

    pipe = TestBed.inject(GroupPipe);
    store = TestBed.inject(Store);

    store.reset({
      groups: {
        groups: groups,
      },
    });
  });

  it("create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("should call group state", () => {
    const spy = jest.spyOn(store, "selectSnapshot");
    pipe.transform("hello");

    expect(spy).toHaveBeenCalled();
  });
});
