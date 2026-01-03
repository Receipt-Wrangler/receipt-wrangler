import { TestBed } from "@angular/core/testing";
import { Group, GroupRole } from "../../open-api";
import { GroupUtil } from "../../utils/index";
import { GroupTableEditButtonPipe } from "./group-table-edit-button.pipe";

describe("GroupTableEditButtonPipe", () => {
  let pipe: GroupTableEditButtonPipe;
  let groupUtilMock: jest.Mocked<GroupUtil>;

  const mockGroup: Group = {
    id: "123",
    name: "Test Group",
  } as any;

  beforeEach(() => {
    groupUtilMock = { hasGroupAccess: jest.fn() } as unknown as jest.Mocked<GroupUtil>;

    TestBed.configureTestingModule({
      providers: [
        GroupTableEditButtonPipe,
        { provide: GroupUtil, useValue: groupUtilMock }
      ]
    });

    pipe = TestBed.inject(GroupTableEditButtonPipe);
  });

  it("should create the pipe", () => {
    expect(pipe).toBeTruthy();
  });

  describe("transform", () => {
    it("should return edit route when user is group owner", () => {
      groupUtilMock.hasGroupAccess.mockReturnValue(true);
      const isAdmin = false;

      const result = pipe.transform(mockGroup, isAdmin);

      expect(result).toEqual({
        routerLink: [`/groups/${mockGroup.id}/details/edit`],
        queryParams: { tab: "details" }
      });
      expect(groupUtilMock.hasGroupAccess).toHaveBeenCalledWith(
        mockGroup.id,
        GroupRole.Owner,
        false,
        false
      );
    });

    it("should return settings route when user is admin but not owner", () => {
      groupUtilMock.hasGroupAccess.mockReturnValue(false);
      const isAdmin = true;

      const result = pipe.transform(mockGroup, isAdmin);

      expect(result).toEqual({
        routerLink: [`/groups/${mockGroup.id}/settings/edit`],
        queryParams: { tab: "settings" }
      });
      expect(groupUtilMock.hasGroupAccess).toHaveBeenCalledWith(
        mockGroup.id,
        GroupRole.Owner,
        false,
        false
      );
    });

    it("should return view route when user is neither owner nor admin", () => {
      groupUtilMock.hasGroupAccess.mockReturnValue(false);
      const isAdmin = false;

      const result = pipe.transform(mockGroup, isAdmin);

      expect(result).toEqual({
        routerLink: [`/groups/${mockGroup.id}/details/view`],
        queryParams: { tab: "details" }
      });
      expect(groupUtilMock.hasGroupAccess).toHaveBeenCalledWith(
        mockGroup.id,
        GroupRole.Owner,
        false,
        false
      );
    });

    it("should handle undefined group id gracefully", () => {
      const undefinedGroup: Group = {
        ...mockGroup,
        id: undefined
      } as any;
      groupUtilMock.hasGroupAccess.mockReturnValue(false);
      const isAdmin = false;

      const result = pipe.transform(undefinedGroup, isAdmin);
      
      expect(result).toEqual({
        routerLink: ["/groups/undefined/details/view"],
        queryParams: { tab: "details" }
      });
    });
  });
});
