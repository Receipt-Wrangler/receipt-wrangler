import { Injectable } from "@angular/core";
import { AbstractControl, AsyncValidatorFn, ValidationErrors, } from "@angular/forms";
import { map, Observable, of } from "rxjs";
import { CategoryService, TagService } from "../open-api";

type DuplicateValidatorType = "category" | "tag";

@Injectable()
export class DuplicateValidator {
  constructor(
    private categoryService: CategoryService,
    private tagService: TagService
  ) {}

  isUnique(
    type: DuplicateValidatorType,
    threshold: number,
    originalValue: string
  ): AsyncValidatorFn {
    return (
      control: AbstractControl
    ):
      | Promise<ValidationErrors | null>
      | Observable<ValidationErrors | null> => {
      const obsrevable = this.getObservable(type, control);
      return obsrevable.pipe(
        map((usernameCount) => {
          if (usernameCount > threshold && control.value !== originalValue) {
            return { duplicate: true };
          }
          return null;
        })
      );
    };
  }

  private getObservable(
    type: DuplicateValidatorType,
    control: AbstractControl
  ): Observable<number> {
    switch (type) {
      case "category":
        return this.categoryService.getCategoryCountByName(control.value);
      case "tag":
        return this.tagService.getTagCountByName(control.value);
      default:
        return of(0);
    }
  }
}
