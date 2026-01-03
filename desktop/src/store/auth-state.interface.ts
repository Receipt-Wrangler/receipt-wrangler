import { Icon, UserRole } from "../open-api";
import { UserPreferences } from "../open-api/model/userPreferences";

export interface AuthStateInterface {
  userId?: string;
  displayname?: string;
  username?: string;
  expirationDate?: string;
  userRole?: UserRole;
  defaultAvatarColor?: string;
  userPreferences?: UserPreferences;
  icons?: Icon[];
}
