import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';

class MultiSelectField<T> extends StatefulWidget {
  const MultiSelectField(
      {super.key,
      required this.name,
      required this.label,
      required this.itemDisplayName,
      required this.itemName,
      this.initialValue,
      this.onTap,
      this.onRemove,
      this.required});

  final String name;

  final String label;

  final String Function(T) itemDisplayName;

  final String itemName;

  final List<T>? initialValue;

  final Function()? onTap;

  /// Called when a chip's remove button is tapped, with the list that should
  /// replace the current value (the removed item already dropped). Null
  /// installs no remove button at all -- view mode -- mirroring desktop's
  /// `*ngIf="!readonly"` on its `matChipRemove` button.
  ///
  /// The new list goes out through here rather than straight into the field so
  /// that removal takes the exact same path the picker's result does: the
  /// owning forms write the value back with `setValue` (and Quick Scan also
  /// mirrors it into its `QuickScanImage`), so mutating the field internally
  /// would leave those callers holding removed items.
  final void Function(List<T>)? onRemove;

  final bool? required;

  @override
  State<MultiSelectField> createState() => _MultiSelectField<T>();
}

class _MultiSelectField<T> extends State<MultiSelectField<T>> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return FormBuilderField<List<T>>(
      name: widget.name,
      validator:
          (widget.required ?? false) ? FormBuilderValidators.required() : null,
      initialValue: widget!.initialValue as dynamic,
      builder: (FormFieldState<List<T>?> field) {
        Widget buildChipLabel(T thing) {
          return Text(widget.itemDisplayName(thing));
        }

        // By index, not by value: the items carry no equality, so removing
        // `thing` itself would drop the wrong chip whenever a list holds two
        // equal entries. Desktop removes by index for the same reason.
        void removeAt(int index) {
          final remaining = List<T>.from(field.value ?? const [])
            ..removeAt(index);
          widget.onRemove!(remaining);
        }

        // An InputChip rather than a ChoiceChip: ChoiceChip does not implement
        // DeletableChipAttributes, so it has no remove button to offer. Every
        // other property carries over, and under the app's M3 theme the two
        // render identically in the selected state -- the X is the only
        // visible difference.
        //
        // deleteIcon is set explicitly because InputChip's M3 default is
        // Icons.clear (a bare X); desktop's matChipRemove button renders
        // <mat-icon>cancel</mat-icon>, the filled circle. 18 is the size the
        // M3 default carries.
        Widget buildChip(int index, T thing) {
          return InputChip(
            label: buildChipLabel(thing),
            // The selected fill and its white label come from the theme's
            // chipTheme, so every chip in the app agrees.
            showCheckmark: false,
            selected: true,
            onSelected: (bool selected) {
              if (widget.onTap != null) {
                widget!.onTap!();
              }
            },
            deleteIcon: const Icon(Icons.cancel, size: 18),
            onDeleted: widget.onRemove == null ? null : () => removeAt(index),
            deleteButtonTooltipMessage:
                "Remove ${widget.itemDisplayName(thing)}",
          );
        }

        List<Widget> buildChipList() {
          if (field.value != null && field.value!.isNotEmpty) {
            List<Widget> widgets = [];
            for (var index = 0; index < field.value!.length; index++) {
              const space = SizedBox(width: 5);
              widgets.add(buildChip(index, field.value![index]));
              widgets.add(space);
            }
            return widgets;
          } else {
            return [Text("No ${widget.itemName} selected")];
          }
        }

        final decorated = InputDecorator(
          decoration: InputDecoration(labelText: widget.label),
          child: Wrap(
            // The remove button makes each chip wider, so a selection wraps
            // onto a second run sooner; without runSpacing the runs touch.
            runSpacing: 5,
            children: buildChipList(),
          ),
        );

        // View mode has nothing to open, so no tap surface is installed at all
        // -- an opaque detector there would swallow pointers for a no-op.
        if (widget.onTap == null) {
          return decorated;
        }

        // The detector wraps the InputDecorator (not the other way around) and
        // is opaque, so the label, the border and the padding gutters are one
        // tap target. Wrapping only the inner Wrap leaves it shrink-wrapped to
        // its text, and deferToChild drops the gaps between chips. A chip's
        // remove button sits deeper in the tree, so it wins the gesture arena
        // against this detector.
        return GestureDetector(
          behavior: HitTestBehavior.opaque,
          onTap: widget.onTap,
          child: decorated,
        );
      },
    );
  }
}
