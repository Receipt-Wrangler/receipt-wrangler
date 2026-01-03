import { FeatureConfig } from '../open-api';

export class SetFeatureConfig {
  static readonly type = '[FeatureConfig] Set Feature Config';

  constructor(public config: FeatureConfig) {}
}
