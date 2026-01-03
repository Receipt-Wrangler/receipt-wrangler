import { About } from "../open-api";

export class SetAbout {
  static readonly type = "[About] Set About State";

  constructor(public about: About) {}
}
