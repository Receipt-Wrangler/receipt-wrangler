import 'dart:typed_data';

import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_scan_form.dart';
import 'package:receipt_wrangler_mobile/shared/classes/quick_scan_image.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/category_select_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/tag_select_field.dart';

// QuickScanForm shows/hides + conditionally-requires each field per the selected
// group's quick-scan config (GroupReceiptSettings.quickScan*). When a field is
// hidden/optional the server backfills a default, so the form omits it.

const _groupId = 1;

api.GroupReceiptSettings _settings({
  bool paidByEnabled = true,
  bool paidByRequired = true,
  bool statusEnabled = true,
  bool statusRequired = true,
  bool categoriesEnabled = false,
  bool categoriesRequired = false,
  bool tagsEnabled = false,
  bool tagsRequired = false,
}) {
  return (api.GroupReceiptSettingsBuilder()
        ..id = _groupId
        ..createdAt = ''
        ..groupId = _groupId
        ..quickScanPaidByEnabled = paidByEnabled
        ..quickScanPaidByRequired = paidByRequired
        ..quickScanStatusEnabled = statusEnabled
        ..quickScanStatusRequired = statusRequired
        ..quickScanCategoriesEnabled = categoriesEnabled
        ..quickScanCategoriesRequired = categoriesRequired
        ..quickScanTagsEnabled = tagsEnabled
        ..quickScanTagsRequired = tagsRequired)
      .build();
}

api.Group _group(api.GroupReceiptSettings settings) => (api.GroupBuilder()
      ..id = _groupId
      ..createdAt = ''
      ..name = 'Test Group'
      ..isAllGroup = false
      ..status = api.GroupStatus.ACTIVE
      ..groupMembers = ListBuilder<api.GroupMember>()
      ..groupReceiptSettings.replace(settings))
    .build();

QuickScanImage _image(int? groupId) => QuickScanImage(
      multipartFile: MultipartFile.fromBytes(const <int>[]),
      bytes: Uint8List(0),
      formKey: GlobalKey<FormBuilderState>(),
      groupId: groupId,
    );

/// Pumps [QuickScanForm] with a real [GroupModel] holding one group configured
/// by [settings]. [imageGroupId] drives which group's config the form reads
/// (0 = no group selected → null settings fallback). Returns the form key.
Future<GlobalKey<FormBuilderState>> _pumpForm(
  WidgetTester tester, {
  required api.GroupReceiptSettings settings,
  int imageGroupId = _groupId,
}) async {
  final image = _image(imageGroupId);
  final groupModel = GroupModel()..setGroups([_group(settings)]);

  await tester.pumpWidget(
    MultiProvider(
      providers: [
        ChangeNotifierProvider<GroupModel>.value(value: groupModel),
        ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
        ChangeNotifierProvider<UserPreferencesModel>(
            create: (_) => UserPreferencesModel()),
        ChangeNotifierProvider<CategoryModel>(create: (_) => CategoryModel()),
        ChangeNotifierProvider<TagModel>(create: (_) => TagModel()),
        ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
      ],
      child: MaterialApp(
        home: Scaffold(
          body: SingleChildScrollView(
            child: QuickScanForm(
              formKey: image.formKey,
              image: image,
              index: 0,
              onFormChangeCallback: (_, __, ___, ____, _____) {},
            ),
          ),
        ),
      ),
    ),
  );
  await tester.pump();
  return image.formKey;
}

Finder _dropdown(String name) => find.byWidgetPredicate(
    (w) => w is FormBuilderDropdown && w.name == name);

void main() {
  testWidgets('shows/hides each field per the group config', (tester) async {
    // Paid By hidden, Status shown, Categories shown, Tags hidden.
    await _pumpForm(
      tester,
      settings: _settings(
        paidByEnabled: false,
        statusEnabled: true,
        categoriesEnabled: true,
        tagsEnabled: false,
      ),
    );

    expect(_dropdown('groupId'), findsOneWidget); // always present
    expect(_dropdown('paidByUserId'), findsNothing);
    expect(_dropdown('status'), findsOneWidget);
    expect(find.byType(CategorySelectField), findsOneWidget);
    expect(find.byType(TagSelectField), findsNothing);
  });

  testWidgets('shows every field when all are enabled', (tester) async {
    await _pumpForm(
      tester,
      settings: _settings(
        paidByEnabled: true,
        statusEnabled: true,
        categoriesEnabled: true,
        tagsEnabled: true,
      ),
    );

    expect(_dropdown('paidByUserId'), findsOneWidget);
    expect(_dropdown('status'), findsOneWidget);
    expect(find.byType(CategorySelectField), findsOneWidget);
    expect(find.byType(TagSelectField), findsOneWidget);
  });

  testWidgets('falls back to backend defaults when no group is selected',
      (tester) async {
    // imageGroupId 0 → getGroupReceiptSettings returns null → paid-by/status
    // default shown, categories/tags default hidden.
    await _pumpForm(
      tester,
      settings: _settings(),
      imageGroupId: 0,
    );

    expect(_dropdown('paidByUserId'), findsOneWidget);
    expect(_dropdown('status'), findsOneWidget);
    expect(find.byType(CategorySelectField), findsNothing);
    expect(find.byType(TagSelectField), findsNothing);
  });

  testWidgets('a shown+required field carries a required validator',
      (tester) async {
    await _pumpForm(
      tester,
      settings: _settings(paidByEnabled: true, paidByRequired: true),
    );

    final paidBy = tester.widget(_dropdown('paidByUserId')) as FormBuilderDropdown;
    expect(paidBy.validator, isNotNull);
    expect(paidBy.validator!(null), isNotNull, reason: 'empty fails validation');
  });

  testWidgets('a shown+optional field has no required validator',
      (tester) async {
    await _pumpForm(
      tester,
      settings: _settings(paidByEnabled: true, paidByRequired: false),
    );

    final paidBy = tester.widget(_dropdown('paidByUserId')) as FormBuilderDropdown;
    expect(paidBy.validator, isNull);
  });
}
