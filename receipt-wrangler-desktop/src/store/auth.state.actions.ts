import { Claims, Icon } from "../open-api";
import { UserPreferences } from "../open-api/model/userPreferences";

export class SetAuthState {
  static readonly type = "[Auth] Set Auth State";

  constructor(public userClaims: Claims) {}
}

export class SetUserPreferences {
  static readonly type = "[Auth] Set User PReferences";

  constructor(public userPreferences: UserPreferences) {}
}

export class SetIcons {
  static readonly type = "[Auth] Set Icons";

  constructor(public icons: Icon[]) {}
}

export class Logout {
  static readonly type = "[Auth] Logout";
}
