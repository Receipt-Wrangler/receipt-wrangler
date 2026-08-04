import { Component, computed, effect, input, output } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormArray, FormControl } from "@angular/forms";
import { CategoryAutocompleteComponent } from "../../category-autocomplete/category-autocomplete.component";
import { Category, Tag } from "../../open-api";
import { TagAutocompleteComponent } from "../../tag-autocomplete/tag-autocomplete.component";

/**
 * The ids a picker is allowed to offer, and the name of whatever imposes the
 * limit (shown to the admin so a shrunken pool is explained rather than
 * mysterious).
 *
 * An EMPTY id array means unrestricted for that resource — the same convention
 * the backend uses, where a group role with no grants sees every category/tag.
 */
export interface GrantCeiling {
  categoryIds: number[];
  tagIds: number[];
  label: string;
}

/** The picker's current selection, as ids. */
export interface GrantSelection {
  categoryIds: number[];
  tagIds: number[];
}

// The user form renders one picker per group, so the underlying autocompletes
// would otherwise all derive the same DOM id from their label. Duplicate ids
// break <label for> association (every field after the first loses its
// accessible name) and misdirect the base autocomplete's getElementById-based
// filter clear to the first instance's input. A per-instance counter keeps them
// unique.
let grantPickerInstanceCount = 0;

/**
 * Shared category/tag grant picker.
 *
 * Used in two places with different meanings: on the role form it edits a group
 * role's grants (the ceiling for everyone holding that role), and on the user /
 * group-member forms it edits one member's individual assignment, which narrows
 * WITHIN that ceiling. Passing a `ceiling` filters the offered pool to it, so an
 * admin can only pick something the member could actually end up seeing (the
 * backend rejects out-of-ceiling ids with a 400 regardless).
 */
@Component({
  selector: "app-grant-picker",
  standalone: true,
  imports: [CategoryAutocompleteComponent, TagAutocompleteComponent],
  templateUrl: "./grant-picker.component.html",
  styleUrl: "./grant-picker.component.scss",
})
export class GrantPickerComponent {
  public readonly categoryPool = input<Category[]>([]);
  public readonly tagPool = input<Tag[]>([]);

  /** Currently granted ids. Resolved to pool objects once the pool arrives. */
  public readonly selectedCategoryIds = input<number[]>([]);
  public readonly selectedTagIds = input<number[]>([]);

  public readonly readonly = input(false);

  /** Optional upper bound on what may be picked. Null means no limit. */
  public readonly ceiling = input<GrantCeiling | null>(null);

  public readonly grantsChange = output<GrantSelection>();

  public readonly grantedCategories = new FormArray<FormControl<Category>>([]);
  public readonly grantedTags = new FormArray<FormControl<Tag>>([]);

  // Set once the user touches a picker, after which the seeding effect stops
  // overwriting their selection (see the constructor).
  private userEditedCategories = false;
  private userEditedTags = false;

  private readonly instanceId = ++grantPickerInstanceCount;
  public readonly categoryInputId = `grant-categories-${this.instanceId}`;
  public readonly tagInputId = `grant-tags-${this.instanceId}`;

  /**
   * The pools actually offered, narrowed to the ceiling. Filtering rather than
   * disabling individual options is deliberate: the shared `app-autocomlete` has
   * no per-option disable support, and adding one would touch a component used
   * across the whole app. The hint lines below explain the narrowing instead.
   */
  public readonly availableCategories = computed<Category[]>(() =>
    this.withinCeiling(this.categoryPool(), this.ceiling()?.categoryIds),
  );

  public readonly availableTags = computed<Tag[]>(() =>
    this.withinCeiling(this.tagPool(), this.ceiling()?.tagIds),
  );

  public readonly categoryCeilingHint = computed<string>(() =>
    this.ceilingHint("categories", this.categoryPool(), this.availableCategories()),
  );

  public readonly tagCeilingHint = computed<string>(() =>
    this.ceilingHint("tags", this.tagPool(), this.availableTags()),
  );

  constructor() {
    // Resolve granted ids to pool objects once the pool lands (pool and ids load
    // independently). Written with emitEvent: false so seeding the form never
    // echoes back out as a user edit.
    //
    // Seeding STOPS once the user has picked something. The effect's inputs keep
    // changing after mount — the pool arrives asynchronously, and the ceiling
    // changes again when the host's role list resolves — so without this guard a
    // late re-run would silently discard the user's selection and hand the host
    // back the original value.
    effect(() => {
      const pool = this.availableCategories();
      const ids = this.selectedCategoryIds();
      if (this.userEditedCategories) {
        return;
      }
      this.setGrantArray(
        this.grantedCategories,
        pool.filter((category) => category.id !== undefined && ids.includes(category.id)),
      );
    });

    effect(() => {
      const pool = this.availableTags();
      const ids = this.selectedTagIds();
      if (this.userEditedTags) {
        return;
      }
      this.setGrantArray(
        this.grantedTags,
        pool.filter((tag) => tag.id !== undefined && ids.includes(tag.id)),
      );
    });

    // valueChanges only fires for real user edits — seeding uses emitEvent: false.
    this.grantedCategories.valueChanges.pipe(takeUntilDestroyed()).subscribe(() => {
      this.userEditedCategories = true;
      this.emitSelection();
    });
    this.grantedTags.valueChanges.pipe(takeUntilDestroyed()).subscribe(() => {
      this.userEditedTags = true;
      this.emitSelection();
    });
  }

  /** The current selection, for callers that read on submit rather than on change. */
  public selection(): GrantSelection {
    return {
      categoryIds: toIds(this.grantedCategories.value as Category[]),
      tagIds: toIds(this.grantedTags.value as Tag[]),
    };
  }

  private emitSelection(): void {
    this.grantsChange.emit(this.selection());
  }

  private withinCeiling<T extends { id?: number }>(pool: T[], allowedIds?: number[]): T[] {
    // Empty/absent means unrestricted, matching the backend's grant convention.
    if (!allowedIds?.length) {
      return pool;
    }
    return pool.filter((entry) => entry.id !== undefined && allowedIds.includes(entry.id));
  }

  private ceilingHint(resource: string, pool: unknown[], available: unknown[]): string {
    const ceiling = this.ceiling();
    if (!ceiling || available.length === pool.length) {
      return "";
    }
    return `Limited to the ${available.length} ${resource} allowed by ${ceiling.label}.`;
  }

  private setGrantArray<T>(array: FormArray, values: T[]): void {
    array.clear({ emitEvent: false });
    for (const value of values) {
      array.push(new FormControl(value), { emitEvent: false });
    }
  }
}

function toIds(entries: { id?: number }[]): number[] {
  return entries.map((entry) => entry.id).filter((id): id is number => id !== undefined);
}
