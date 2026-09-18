import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/constants/colors.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';

/// Total height this bar occupies, including its hairline.
///
/// `submitButtonSpacing` (lib/constants/spacing.dart) is sized off this: where
/// the bar is a floating `Scaffold.bottomSheet` rather than a bottom bar, the
/// scrollable above it has to reserve at least this much or its last field sits
/// underneath and cannot be scrolled clear.
const double bottomSubmitBarHeight =
    _barVerticalPadding * 2 + _buttonHeight + 1;

const double _barVerticalPadding = 12;
const double _buttonHeight = 48;

class BottomSubmitButton extends StatefulWidget {
  const BottomSubmitButton(
      {super.key, required this.onPressed, this.buttonText, this.disabled});

  final String? buttonText;

  final disabled;

  final void Function() onPressed;

  @override
  State<BottomSubmitButton> createState() => _BottomSubmitButtonState();
}

class _BottomSubmitButtonState extends State<BottomSubmitButton> {
  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Container(
      // The bar is its own surface: a hairline separates it from the content it
      // scrolls over, and the button is inset rather than bled to the edges.
      decoration: BoxDecoration(
        color: colorScheme.surface,
        border: Border(top: BorderSide(color: colorScheme.outlineVariant)),
      ),
      padding: const EdgeInsets.symmetric(
          horizontal: 16, vertical: _barVerticalPadding),
      child: SizedBox(
        width: double.infinity,
        height: _buttonHeight,
        child: Consumer<LoadingModel>(
          builder: (context, loadingModel, child) {
            return FilledButton(
              onPressed: (loadingModel.isLoading || (widget.disabled ?? false))
                  ? null
                  : widget.onPressed,
              style: FilledButton.styleFrom(
                shape: const StadiumBorder(),
                backgroundColor: colorScheme.primary,
                foregroundColor: colorScheme.onPrimary,
                // Pinned rather than `outline`: a disabled fill as dark
                // as a border reads as enabled. Exempt from SC 1.4.11.
                disabledBackgroundColor: slate300,
                disabledForegroundColor: colorScheme.onPrimary,
                // styleFrom's textStyle REPLACES the theme's labelLarge, so
                // the font family has to be restated or the label falls back
                // to the platform default.
                textStyle: const TextStyle(
                  fontFamily: appFontFamily,
                  fontSize: 15,
                  fontWeight: FontWeight.w600,
                ),
              ),
              child: loadingModel.isLoading
                  ? SizedBox(
                      height: 22,
                      width: 22,
                      child: CircularProgressIndicator(
                        color: colorScheme.onPrimary,
                        strokeWidth: 2,
                      ),
                    )
                  : Text(widget.buttonText ?? "Submit"),
            );
          },
        ),
      ),
    );
  }
}
