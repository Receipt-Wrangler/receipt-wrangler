import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'mapGet' })
export class MapGetPipe implements PipeTransform {
  public transform(value: Map<any, any>, key: any): any | undefined {
    return value.get(key);
  }
}
