import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'mapKey' })
export class MapKeyPipe implements PipeTransform {
  public transform(value: Map<any, any>): any[] {
    return Array.from(value.keys());
  }
}
