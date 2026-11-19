import { Component, computed, input } from "@angular/core";
import { Group, User } from "../../open-api";
import { MatTooltip } from "@angular/material/tooltip";
import { MatIcon } from "@angular/material/icon";
import { SlicePipe } from "@angular/common";

/** The diameter every existing caller renders at; its sizing lives in the stylesheet. */
export const DEFAULT_AVATAR_SIZE = 50;

@Component({
    selector: "app-avatar",
    templateUrl: "./avatar.component.html",
    styleUrls: ["./avatar.component.scss"],
    imports: [MatTooltip, MatIcon, SlicePipe]
})
export class AvatarComponent {
  public readonly user = input<User>();
  public readonly group = input<Group>();

  /** Diameter in px. */
  public readonly size = input(DEFAULT_AVATAR_SIZE);

  /** A material icon shown in place of the initial, for an avatar that stands for no user or group. */
  public readonly icon = input<string>();

  /** Background colour; falls back to the user's own colour, then the brand blue. */
  public readonly color = input<string>();

  /**
   * The inline sizing for a non-default `size`, or null. Null at the default so a
   * default avatar renders exactly as it always has, sized by its stylesheet alone.
   * The initial and icon scale with the circle (14px / 18px at 32px).
   */
  public readonly sizing = computed(() => {
    const size = this.size();
    if (size === DEFAULT_AVATAR_SIZE) {
      return null;
    }

    return {
      size,
      fontSize: Math.round(size * 0.44),
      iconSize: Math.round(size * 0.56),
    };
  });

  /** The icon's glyph size: scaled for a custom size, 28px at the default. */
  public readonly iconSize = computed(
    () => this.sizing()?.iconSize ?? Math.round(DEFAULT_AVATAR_SIZE * 0.56)
  );
}
