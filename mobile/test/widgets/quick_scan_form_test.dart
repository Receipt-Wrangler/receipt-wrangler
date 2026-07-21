import 'dart:typed_data';

import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_scan_form.dart';
import 'package:receipt_wrangler_mobile/shared/classes/quick_scan_image.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/category_select_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/tag_select_field.dart';

// QuickScanForm shows/hides + conditionally-requires each field per the selected
// group's quick-scan config (GroupReceiptSettings.quickScan*). When a field is
// hidden/optional the server backfills a default, so the form omits it. Field
// visibility/requirement is reactive to the group dropdown: changing the group
// re-reads the new group's config and clears the group-scoped fields.

const _groupId = 1;
const _group2Id = 2;

api.GroupReceiptSettings _settings({
  int id = _groupId,
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
        ..id = id
        ..createdAt = ''
        ..groupId = id
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

api.Group _group(
  api.GroupReceiptSettings settings, {
  String name = 'Test Group',
  List<api.GroupMember> members = const [],
}) =>
    (api.GroupBuilder()
          ..id = settings.groupId
          ..createdAt = ''
          ..name = name
          ..isAllGroup = false
          ..status = api.GroupStatus.ACTIVE
          ..groupMembers = ListBuilder<api.GroupMember>(members)
          ..groupReceiptSettings.replace(settings))
        .build();

api.GroupMember _member(int userId, int groupId) =>
    (api.GroupMemberBuilder()
          ..groupId = groupId
          ..userId = userId)
        .build();

api.UserView _user(int id, String displayName) => (api.UserViewBuilder()
      ..id = id
      ..username = 'u$id'
      ..displayName = displayName
      ..isDummyUser = false)
    .build();

api.Category _category() =>
    (api.CategoryBuilder()..id = 1..name = 'Food').build();

api.Tag _tag() => (api.TagBuilder()..id = 1..name = 'Trip').build();

QuickScanImage _image(
  int? groupId, {
  int? paidByUserId,
  api.ReceiptStatus? status,
}) =>
    QuickScanImage(
      multipartFile: MultipartFile.fromBytes(const <int>[]),
      bytes: Uint8List(0),
      formKey: GlobalKey<FormBuilderState>(),
      groupId: groupId,
      paidByUserId: paidByUserId,
      status: status,
    );

/// Pumps [QuickScanForm] with a real [GroupModel] holding [groups] and a
/// [UserModel] holding [users] (the latter backs the paid-by dropdown items).
/// [imageGroupId] drives which group's config the form initially reads
/// (0 = no group selected → null settings fallback). Returns the form key.
Future<GlobalKey<FormBuilderState>> _pumpFormGroups(
  WidgetTester tester, {
  required List<api.Group> groups,
  required int imageGroupId,
  List<api.UserView> users = const [],
  int? imagePaidByUserId,
  api.ReceiptStatus? imageStatus,
}) async {
  final image = _image(imageGroupId,
      paidByUserId: imagePaidByUserId, status: imageStatus);
  final groupModel = GroupModel()..setGroups(groups);
  final userModel = UserModel()..setUsers(users);

  await tester.pumpWidget(
    MultiProvider(
      providers: [
        ChangeNotifierProvider<GroupModel>.value(value: groupModel),
        ChangeNotifierProvider<UserModel>.value(value: userModel),
        ChangeNotifierProvider<UserPreferencesModel>(
            create: (_) => UserPreferencesModel()),
        ChangeNotifierProvider<CategoryModel>(create: (_) => CategoryModel()),
        ChangeNotifierProvider<TagModel>(create: (_) => TagModel()),
        ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
        // The category/tag pickers open a fullscreen sheet whose TopAppBar reads
        // these two; needed once a test taps a picker open.
        ChangeNotifierProvider<AuthModel>(create: (_) => AuthModel()),
        ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
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

/// Convenience wrapper for the single-group tests.
Future<GlobalKey<FormBuilderState>> _pumpForm(
  WidgetTester tester, {
  required api.GroupReceiptSettings settings,
  int imageGroupId = _groupId,
}) =>
    _pumpFormGroups(tester,
        groups: [_group(settings)], imageGroupId: imageGroupId);

Finder _dropdown(String name) => find.byWidgetPredicate(
    (w) => w is FormBuilderDropdown && w.name == name);

/// Opens the group dropdown and picks [name]. `.last` picks the open menu's
/// copy of the label (the closed-state selected value also renders it).
Future<void> _selectGroup(WidgetTester tester, String name) async {
  await tester.tap(_dropdown('groupId'));
  await tester.pumpAndSettle();
  await tester.tap(find.text(name).last);
  await tester.pumpAndSettle();
}

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

  testWidgets('a group with default settings shows paid-by + status, '
      'hides categories + tags', (tester) async {
    // Non-null settings carrying the backend defaults (distinct from the
    // null-fallback case above).
    await _pumpForm(tester, settings: _settings());

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

  testWidgets('status shown+required carries a required validator',
      (tester) async {
    await _pumpForm(
      tester,
      settings: _settings(statusEnabled: true, statusRequired: true),
    );

    final status = tester.widget(_dropdown('status')) as FormBuilderDropdown;
    expect(status.validator, isNotNull);
    expect(status.validator!(null), isNotNull, reason: 'empty fails validation');
  });

  testWidgets('status shown+optional has no required validator',
      (tester) async {
    await _pumpForm(
      tester,
      settings: _settings(statusEnabled: true, statusRequired: false),
    );

    final status = tester.widget(_dropdown('status')) as FormBuilderDropdown;
    expect(status.validator, isNull);
  });

  testWidgets('re-renders fields per config when the group changes',
      (tester) async {
    // Group A: paid-by/status shown, categories/tags hidden.
    // Group B: the exact inverse — proves visibility is reactive to the group.
    final groupA = _group(
      _settings(
        id: _groupId,
        paidByEnabled: true,
        statusEnabled: true,
        categoriesEnabled: false,
        tagsEnabled: false,
      ),
      name: 'Group A',
    );
    final groupB = _group(
      _settings(
        id: _group2Id,
        paidByEnabled: false,
        statusEnabled: false,
        categoriesEnabled: true,
        tagsEnabled: true,
      ),
      name: 'Group B',
    );
    await _pumpFormGroups(tester,
        groups: [groupA, groupB], imageGroupId: _groupId);

    // Group A's field set.
    expect(_dropdown('paidByUserId'), findsOneWidget);
    expect(_dropdown('status'), findsOneWidget);
    expect(find.byType(CategorySelectField), findsNothing);
    expect(find.byType(TagSelectField), findsNothing);

    await _selectGroup(tester, 'Group B');

    // Group B's field set — flipped.
    expect(_dropdown('paidByUserId'), findsNothing);
    expect(_dropdown('status'), findsNothing);
    expect(find.byType(CategorySelectField), findsOneWidget);
    expect(find.byType(TagSelectField), findsOneWidget);
  });

  testWidgets('clears paid-by / categories / tags when the group changes',
      (tester) async {
    // Both groups show all three group-scoped fields; group A carries the
    // member that backs a valid paid-by selection.
    final groupA = _group(
      _settings(
        id: _groupId,
        categoriesEnabled: true,
        tagsEnabled: true,
      ),
      name: 'Group A',
      members: [_member(42, _groupId)],
    );
    final groupB = _group(
      _settings(
        id: _group2Id,
        categoriesEnabled: true,
        tagsEnabled: true,
      ),
      name: 'Group B',
    );
    final formKey = await _pumpFormGroups(
      tester,
      groups: [groupA, groupB],
      imageGroupId: _groupId,
      users: [_user(42, 'Payer')],
    );

    formKey.currentState!.fields['paidByUserId']!.didChange(42);
    formKey.currentState!.fields['categories']!
        .didChange(<api.Category>[_category()]);
    formKey.currentState!.fields['tags']!.didChange(<api.Tag>[_tag()]);
    await tester.pump();
    expect(formKey.currentState!.fields['paidByUserId']!.value, 42);
    expect(formKey.currentState!.fields['categories']!.value, isNotEmpty);
    expect(formKey.currentState!.fields['tags']!.value, isNotEmpty);

    await _selectGroup(tester, 'Group B');

    expect(formKey.currentState!.fields['paidByUserId']!.value, isNull,
        reason: 'paid-by cleared on group change');
    expect(formKey.currentState!.fields['categories']!.value, isEmpty,
        reason: 'categories cleared on group change');
    expect(formKey.currentState!.fields['tags']!.value, isEmpty,
        reason: 'tags cleared on group change');
  });

  testWidgets('re-evaluates the required validator after a group switch',
      (tester) async {
    // Paid-by optional in group A, required in group B.
    final groupA = _group(
      _settings(id: _groupId, paidByEnabled: true, paidByRequired: false),
      name: 'Group A',
    );
    final groupB = _group(
      _settings(id: _group2Id, paidByEnabled: true, paidByRequired: true),
      name: 'Group B',
    );
    await _pumpFormGroups(tester,
        groups: [groupA, groupB], imageGroupId: _groupId);

    var paidBy = tester.widget(_dropdown('paidByUserId')) as FormBuilderDropdown;
    expect(paidBy.validator, isNull, reason: 'optional in group A');

    await _selectGroup(tester, 'Group B');

    paidBy = tester.widget(_dropdown('paidByUserId')) as FormBuilderDropdown;
    expect(paidBy.validator, isNotNull, reason: 'required in group B');
    expect(paidBy.validator!(null), isNotNull,
        reason: 'empty fails after switch');
  });

  testWidgets('tapping Categories opens the picker with no shell context',
      (tester) async {
    // ContextModel.shellContext is null in this harness — the same condition as
    // the real Quick Scan flow, which never mounts the receipt-form screen that
    // sets it. Pre-fix this tap crashed via Navigator.of(null); the fix falls
    // back to the field's own (mounted) context so the picker opens.
    await _pumpForm(tester, settings: _settings(categoriesEnabled: true));

    await tester.tap(find.text('No Categories selected'));
    await tester.pumpAndSettle();

    expect(find.text('Select Categories'), findsOneWidget);
  });

  testWidgets('tapping Tags opens the picker with no shell context',
      (tester) async {
    await _pumpForm(tester, settings: _settings(tagsEnabled: true));

    await tester.tap(find.text('No Tags selected'));
    await tester.pumpAndSettle();

    expect(find.text('Select Tags'), findsOneWidget);
  });

  // Prefill (from userPreferences.quickScanDefault*) vs the group's config: the
  // per-image form seeds paid-by/status from the image, but a field the group
  // hides must not appear -- the preset "falls off" (and _submitQuickScan sends
  // the sentinel for it).
  testWidgets('a prefilled paid-by shows when the group shows paid-by',
      (tester) async {
    final group = _group(
      _settings(id: _groupId, paidByEnabled: true),
      members: [_member(42, _groupId)],
    );
    await _pumpFormGroups(
      tester,
      groups: [group],
      imageGroupId: _groupId,
      users: [_user(42, 'Payer')],
      imagePaidByUserId: 42,
    );

    final paidBy = tester.widget(_dropdown('paidByUserId')) as FormBuilderDropdown;
    expect(paidBy.initialValue, 42, reason: 'prefill honored on a shown field');
  });

  testWidgets('a prefilled paid-by falls off when the group hides paid-by',
      (tester) async {
    final group = _group(_settings(id: _groupId, paidByEnabled: false));
    await _pumpFormGroups(
      tester,
      groups: [group],
      imageGroupId: _groupId,
      imagePaidByUserId: 42,
    );

    expect(_dropdown('paidByUserId'), findsNothing,
        reason: 'preset paid-by fell off the hidden field');
  });

  testWidgets('a prefilled status shows when the group shows status',
      (tester) async {
    final group = _group(_settings(id: _groupId, statusEnabled: true));
    await _pumpFormGroups(
      tester,
      groups: [group],
      imageGroupId: _groupId,
      imageStatus: api.ReceiptStatus.OPEN,
    );

    final status = tester.widget(_dropdown('status')) as FormBuilderDropdown;
    expect(status.initialValue, api.ReceiptStatus.OPEN,
        reason: 'prefill honored on a shown field');
  });

  testWidgets('a prefilled status falls off when the group hides status',
      (tester) async {
    final group = _group(_settings(id: _groupId, statusEnabled: false));
    await _pumpFormGroups(
      tester,
      groups: [group],
      imageGroupId: _groupId,
      imageStatus: api.ReceiptStatus.OPEN,
    );

    expect(_dropdown('status'), findsNothing,
        reason: 'preset status fell off the hidden field');
  });
}
