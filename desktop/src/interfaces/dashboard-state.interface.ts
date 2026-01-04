import { Dashboard } from "../open-api";

export interface DashboardStateInterface {
  dashboards: { [groupId: string]: Dashboard[] };
}
