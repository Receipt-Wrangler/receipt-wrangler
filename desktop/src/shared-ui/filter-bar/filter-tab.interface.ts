/**
 * A single tab in a {@link FilterBarComponent}.
 */
export interface FilterTab {
  /** Identifier emitted when this tab is selected. */
  value: string;
  /** Visible label. */
  label: string;
  /** Optional Material icon ligature shown before the label. */
  icon?: string;
  /** Optional count rendered as a pill; hidden when undefined. */
  count?: number;
}
