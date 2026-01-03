import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { UserState } from "../store";
import { UserPipe } from "./user.pipe";

describe("UserPipe", () => {
  let pipe: UserPipe;
  let store: Store;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [UserPipe],
      imports: [NgxsModule.forRoot([UserState])],
      providers: [UserPipe],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    store = TestBed.inject(Store);
    pipe = TestBed.inject(UserPipe);

    store.reset({
      users: {
        users: [{ id: 1 }, { id: 2 }, { id: 3 }],
      },
    });
  });

  it("create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("should return a user", () => {
    const result = pipe.transform("1");

    expect(result).toEqual({ id: 1 } as any);
  });

  it("should return a user", () => {
    const result = pipe.transform("not a user id");

    expect(result).toEqual(undefined);
  });
});
