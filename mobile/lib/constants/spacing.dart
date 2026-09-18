import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

const textFieldSpacing = SizedBox(height: 20);

const lastFieldSpacing = SizedBox(height: 40);

const headerSpacing = SizedBox(height: 30);

/// Clears a floating `BottomSubmitButton` at the end of a scrollable, with a
/// little slack. Derived from the bar's own height so the two cannot drift.
const submitButtonSpacing = SizedBox(height: bottomSubmitBarHeight + 16);
