import 'package:built_collection/built_collection.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/custom_field_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_form.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import 'widget_test_helpers.dart';

// Builders + a pump helper for widget tests of `receipt_form.dart`. Modelled on
// `test/widgets/quick_scan_form_test.dart`: real models seeded through their
// public setters and real generated API builders, no mocks.

/// Builds a [api.CustomFieldOption] for a SELECT custom field.
api.CustomFieldOption buildCustomFieldOption({
  required int id,
  required int customFieldId,
  required String value,
}) =>
    (api.CustomFieldOptionBuilder()
          ..id = id
          ..createdAt = ''
          ..customFieldId = customFieldId
          ..value = value)
        .build();

/// Builds a custom field *template* — one entry of the catalog
/// [CustomFieldModel] holds.
api.CustomField buildCustomField({
  required int id,
  required String name,
  api.CustomFieldType type = api.CustomFieldType.TEXT,
  String? description,
  List<api.CustomFieldOption>? options,
}) =>
    (api.CustomFieldBuilder()
          ..id = id
          ..createdAt = ''
          ..name = name
          ..type = type
          ..description = description
          ..options = options == null
              ? null
              : ListBuilder<api.CustomFieldOption>(options))
        .build();

/// Builds a custom field *value* — one entry of a receipt's `customFields`.
/// Every value column defaults to null, i.e. an **empty** value, which is a
/// first-class state (see `lib/shared/functions/custom_field_values.dart`).
api.CustomFieldValue buildCustomFieldValue({
  required int customFieldId,
  int id = 0,
  int receiptId = 0,
  String? stringValue,
  String? dateValue,
  int? selectValue,
  String? currencyValue,
  bool? booleanValue,
}) =>
    (api.CustomFieldValueBuilder()
          ..id = id
          ..createdAt = ''
          ..customFieldId = customFieldId
          ..receiptId = receiptId
          ..stringValue = stringValue
          ..dateValue = dateValue
          ..selectValue = selectValue
          ..currencyValue = currencyValue
          ..booleanValue = booleanValue)
        .build();

/// Builds a [api.GroupReceiptSettings]. Categories/tags are shown by default
/// (`hide* == false`) so the pickers mount, matching the backend defaults.
///
/// [defaultCustomFieldIds] is the group's default custom field set. Leave it
/// null to model a group that has none — the backend always serializes `[]`,
/// but an absent value has to behave identically (a released build predating
/// the field would see null).
api.GroupReceiptSettings buildGroupReceiptSettings({
  required int groupId,
  int? id,
  bool hideReceiptCategories = false,
  bool hideReceiptTags = false,
  bool hideComments = false,
  bool hideImages = false,
  List<int>? defaultCustomFieldIds,
  bool? applyDefaultCustomFieldsOnIngest,
  bool? receiptSummaryEnabled,
  api.ReceiptSummaryPosition? receiptSummaryPosition,
  List<int>? receiptSummaryCustomFieldIds,
  List<api.ReceiptStatus>? receiptSummaryStatuses,
}) =>
    (api.GroupReceiptSettingsBuilder()
          ..id = id ?? groupId
          ..createdAt = ''
          ..groupId = groupId
          ..hideReceiptCategories = hideReceiptCategories
          ..hideReceiptTags = hideReceiptTags
          ..hideComments = hideComments
          ..hideImages = hideImages
          ..applyDefaultCustomFieldsOnIngest = applyDefaultCustomFieldsOnIngest
          ..receiptSummaryEnabled = receiptSummaryEnabled
          ..receiptSummaryPosition = receiptSummaryPosition
          ..receiptSummaryCustomFieldIds = receiptSummaryCustomFieldIds == null
              ? null
              : ListBuilder<int>(receiptSummaryCustomFieldIds)
          ..receiptSummaryStatuses = receiptSummaryStatuses == null
              ? null
              : ListBuilder<api.ReceiptStatus>(receiptSummaryStatuses)
          ..defaultCustomFieldIds = defaultCustomFieldIds == null
              ? null
              : ListBuilder<int>(defaultCustomFieldIds))
        .build();

/// Builds a [api.GroupMember] (backs the paid-by dropdown items).
api.GroupMember buildGroupMember({required int userId, required int groupId}) =>
    (api.GroupMemberBuilder()
          ..groupId = groupId
          ..userId = userId)
        .build();

/// Builds a [api.Group], defaulting to a settings row that shows everything.
api.Group buildGroup({
  required int id,
  String name = 'Test Group',
  api.GroupReceiptSettings? receiptSettings,
  List<api.GroupMember> members = const [],
  List<int>? defaultCustomFieldIds,
  bool isAllGroup = false,
  bool? receiptSummaryEnabled,
}) =>
    (api.GroupBuilder()
          ..id = id
          ..createdAt = ''
          ..name = name
          ..isAllGroup = isAllGroup
          ..status = api.GroupStatus.ACTIVE
          ..groupMembers = ListBuilder<api.GroupMember>(members)
          ..groupReceiptSettings.replace(receiptSettings ??
              buildGroupReceiptSettings(
                groupId: id,
                defaultCustomFieldIds: defaultCustomFieldIds,
                receiptSummaryEnabled: receiptSummaryEnabled,
              )))
        .build();

api.UserView buildUserView({required int id, String? displayName}) =>
    (api.UserViewBuilder()
          ..id = id
          ..username = 'u$id'
          ..displayName = displayName ?? 'User $id'
          ..isDummyUser = false)
        .build();

/// The models the pumped [ReceiptForm] is wired to, so a test can seed or read
/// them back without digging through the widget tree.
class ReceiptFormHarness {
  ReceiptFormHarness({
    required this.receiptModel,
    required this.groupModel,
    required this.userModel,
    required this.customFieldModel,
  });

  final ReceiptModel receiptModel;

  final GroupModel groupModel;

  final UserModel userModel;

  final CustomFieldModel customFieldModel;

  /// The form the widget attached itself to. Read through the model rather
  /// than cached, matching `receipt_form.dart`.
  GlobalKey<FormBuilderState> get formKey => receiptModel.receiptFormKey;

  /// The current value of the FormBuilder field a custom field renders under.
  dynamic customFieldValue(int customFieldId) =>
      formKey.currentState?.fields['customField_$customFieldId']?.value;
}

String _locationFor(WranglerFormState formState, int receiptId) {
  switch (formState) {
    case WranglerFormState.add:
      return '/receipts/add';
    case WranglerFormState.edit:
      return '/receipts/$receiptId/edit';
    case WranglerFormState.view:
      return '/receipts/$receiptId/view';
  }
}

/// Pumps [ReceiptForm] with the full provider tree its `build()` reaches:
/// [ReceiptModel] (receipt + form key), [GroupModel] / [UserModel] (the group
/// and paid-by dropdowns), [CustomFieldModel] (the custom field catalog),
/// [CategoryModel] / [TagModel] / [ContextModel] (the category + tag pickers),
/// [SystemSettingsModel] (every `AmountField`), plus [PermissionsModel],
/// [AuthModel] and [LoadingModel], which the sheets opened from the form read.
///
/// The form derives its mode from the route (`getFormStateFromContext`), so a
/// real [GoRouter] is mounted at the location matching [formState].
///
/// Pass [receiptModel] to pump against a model an earlier call returned: the
/// second pump replaces the tree, so `ReceiptForm` gets a fresh `State`
/// while the working copy and its auto-applied provenance carry over. That is
/// what a view -> edit navigation does in the app -- the app bar menu's Edit
/// entry is a plain `go`, and `ReceiptFormScreen` re-hydrates only for a
/// different receipt id -- so a remount is reproduced rather than simulated.
/// [receipt] is ignored when a model is supplied; the model already holds one.
Future<ReceiptFormHarness> pumpReceiptForm(
  WidgetTester tester, {
  required List<api.Group> groups,
  api.Receipt? receipt,
  ReceiptModel? receiptModel,
  List<api.CustomField> customFields = const [],
  List<api.UserView> users = const [],
  WranglerFormState formState = WranglerFormState.add,
  /// Mount under a real theme. Null keeps the bare `MaterialApp.router` every
  /// existing caller gets; pass `buildAppTheme()` when the test is about how the
  /// form is drawn rather than how it behaves.
  ThemeData? theme,
  /// Mirrors `ReceiptFormScreen`'s real structure — a [ScreenWrapper] whose
  /// `bottomSheetWidget` is the submit button — rather than the bare `Scaffold`
  /// every other case uses. `Scaffold.bottomSheet` *floats over* the body, so
  /// this is the only shape in which the form's tail can be caught sitting
  /// underneath the button. Defaults off, leaving every existing case's tree
  /// unchanged.
  bool pinnedSubmitButton = false,
  /// Wraps the whole pumped tree. Only the `tool/demo_capture/` screenshot
  /// harness uses it, to put the form inside a capture surface — this helper
  /// owns the `pumpWidget` call, so a caller cannot nest the result itself.
  Widget Function(Widget)? wrap,
  /// Off only for the screenshot harness, where the debug ribbon is just noise
  /// across a captured panel. Defaults to the framework's own behaviour so no
  /// existing test changes.
  bool showDebugBanner = true,
}) async {
  registerCustomCurrencyForTests();

  final seededReceipt = receiptModel?.receipt ?? receipt ?? getDefaultReceipt();

  // Seed the receipt before the first pump: `ReceiptForm` captures the model's
  // form key once (`late final`), and `setReceipt` regenerates that key
  // whenever the receipt identity changes.
  final model =
      receiptModel ?? (ReceiptModel()..setReceipt(seededReceipt, false));
  final groupModel = GroupModel()..setGroups(groups);
  final userModel = UserModel()..setUsers(users);
  final customFieldModel = CustomFieldModel()..setCustomFields(customFields);

  final router = GoRouter(
    initialLocation: _locationFor(formState, seededReceipt.id),
    routes: [
      for (final path in const [
        '/receipts/add',
        '/receipts/:receiptId/edit',
        '/receipts/:receiptId/view',
      ])
        GoRoute(
          path: path,
          builder: (_, __) => pinnedSubmitButton
              ? ScreenWrapper(
                  bottomSheetWidget: BottomSubmitButton(onPressed: () {}),
                  child: const SingleChildScrollView(child: ReceiptForm()),
                )
              : const Scaffold(
                  body: SingleChildScrollView(child: ReceiptForm()),
                ),
        ),
    ],
  );

  final tree = MultiProvider(
      providers: [
        ChangeNotifierProvider<ReceiptModel>.value(value: model),
        ChangeNotifierProvider<GroupModel>.value(value: groupModel),
        ChangeNotifierProvider<UserModel>.value(value: userModel),
        ChangeNotifierProvider<CustomFieldModel>.value(value: customFieldModel),
        ChangeNotifierProvider<CategoryModel>(create: (_) => CategoryModel()),
        ChangeNotifierProvider<TagModel>(create: (_) => TagModel()),
        ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
        ChangeNotifierProvider<SystemSettingsModel>(
          create: (_) => SystemSettingsModel(),
        ),
        ChangeNotifierProvider<PermissionsModel>(
          create: (_) => PermissionsModel(),
        ),
        ChangeNotifierProvider<AuthModel>(create: (_) => AuthModel()),
        ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
      ],
      child: MaterialApp.router(
        routerConfig: router,
        theme: theme,
        debugShowCheckedModeBanner: showDebugBanner,
      ),
    );

  await tester.pumpWidget(wrap == null ? tree : wrap(tree));
  await tester.pumpAndSettle();

  return ReceiptFormHarness(
    receiptModel: model,
    groupModel: groupModel,
    userModel: userModel,
    customFieldModel: customFieldModel,
  );
}
