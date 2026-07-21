import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/reports/widgets/report_list_item.dart';

import '../helpers/report_test_helpers.dart';

void main() {
  Future<void> pumpItem(WidgetTester tester, List<String>? allowedActions) async {
    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: ReportListItem(
            template: buildReportTemplate(id: 7, allowedActions: allowedActions),
            onChanged: () {},
          ),
        ),
      ),
    );
    await tester.pump();
  }

  final previewKey = const ValueKey('report-preview-7');
  final generateKey = const ValueKey('report-generate-7');
  final deleteKey = const ValueKey('report-delete-7');

  testWidgets('read allowedAction shows only the preview button', (tester) async {
    await pumpItem(tester, ['read']);
    expect(find.byKey(previewKey), findsOneWidget);
    expect(find.byKey(generateKey), findsNothing);
    expect(find.byKey(deleteKey), findsNothing);
  });

  testWidgets('generate allowedAction shows only the generate button',
      (tester) async {
    await pumpItem(tester, ['generate']);
    expect(find.byKey(previewKey), findsNothing);
    expect(find.byKey(generateKey), findsOneWidget);
    expect(find.byKey(deleteKey), findsNothing);
  });

  testWidgets('delete allowedAction shows only the delete button',
      (tester) async {
    await pumpItem(tester, ['delete']);
    expect(find.byKey(previewKey), findsNothing);
    expect(find.byKey(generateKey), findsNothing);
    expect(find.byKey(deleteKey), findsOneWidget);
  });

  testWidgets('all actions render all three buttons', (tester) async {
    await pumpItem(tester, ['read', 'generate', 'delete']);
    expect(find.byKey(previewKey), findsOneWidget);
    expect(find.byKey(generateKey), findsOneWidget);
    expect(find.byKey(deleteKey), findsOneWidget);
  });

  testWidgets('empty allowedActions renders no action buttons', (tester) async {
    await pumpItem(tester, []);
    expect(find.byKey(previewKey), findsNothing);
    expect(find.byKey(generateKey), findsNothing);
    expect(find.byKey(deleteKey), findsNothing);
  });

  testWidgets('null allowedActions renders no action buttons', (tester) async {
    await pumpItem(tester, null);
    expect(find.byKey(previewKey), findsNothing);
    expect(find.byKey(generateKey), findsNothing);
    expect(find.byKey(deleteKey), findsNothing);
  });

  testWidgets('shows the template name and formats subtitle', (tester) async {
    await pumpItem(tester, ['read']);
    expect(find.text('Test Report'), findsOneWidget);
    expect(find.text('CSV'), findsOneWidget);
  });
}
