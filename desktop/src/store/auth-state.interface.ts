import { Icon } from "../open-api";
import { UserPreferences } from "../open-api/model/userPreferences";

export interface AuthStateInterface {
  userId?: string;
  displayname?: string;
  username?: string;
  expirationDate?: string;
  defaultAvatarColor?: string;
  userPreferences?: UserPreferences;
  icons?: Icon[];
  appPermissions?: string[];
  groupPermissions?: { [groupId: number]: string[] };
}
