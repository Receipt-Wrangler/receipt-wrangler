import { User } from '../open-api/model/user';

export class SetUsers {
  static readonly type = '[User] Set Users';
  constructor(public users: User[]) {}
}

export class UpdateUser {
  static readonly type = '[User] Update User';
  constructor(public userId: string, public user: User) {}
}
export class AddUser {
  static readonly type = '[User] Add User';
  constructor(public user: User) {}
}

export class RemoveUser {
  static readonly type = '[User] Remove User';
  constructor(public userId: string) {}
}

export class RemoveUsers {
  static readonly type = '[User] Remove Users';
  constructor(public userIds: string[]) {}
}
