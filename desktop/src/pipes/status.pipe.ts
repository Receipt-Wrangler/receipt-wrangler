import { formatStatus } from "src/utils";

import { Pipe, PipeTransform } from "@angular/core";

@Pipe({ name: 'status' })
export class StatusPipe implements PipeTransform {
  public transform(status: string): string {
    return formatStatus(status);
  }
}
