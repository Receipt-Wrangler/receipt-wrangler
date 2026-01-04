import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'name',
    standalone: false
})
export class NamePipe implements PipeTransform {
  public transform(value: any[]): string {
    return value.map((v) => v['name'] ?? '').join(', ');
  }
}
