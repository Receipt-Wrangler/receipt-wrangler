import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

import 'receipt_form_test_helpers.dart';

/// The models every filter surface reads, seeded with one All group and two
/// real ones so the group-scoped rules are exercisable.
class ReceiptFilterHarness {
  ReceiptFilterHarness({
    required this.receiptListModel,
    required this.groupModel,
    required this.categoryModel,
    required this.tagModel,
    required this.userModel,
  });

  final ReceiptListModel receiptListModel;
  final GroupModel groupModel;
  final CategoryModel categoryModel;
  final TagModel tagModel;
  final UserModel userModel;

  static const allGroupId = 1;
  static const householdId = 2;
  static const officeId = 3;
}

ReceiptFilterHarness buildReceiptFilterHarness({
  List<api.Category> householdCategories = const [],
  List<api.Tag> householdTags = const [],
}) {
  final groupModel = GroupModel();
  final categoryModel = CategoryModel();
  final tagModel = TagModel();
  final userModel = UserModel();

  userModel.setUsers([
    buildUserView(id: 10, displayName: "Noah Hall"),
    buildUserView(id: 11, displayName: "Dana Kim"),
  ]);

  groupModel.setGroups([
    buildGroup(
      id: ReceiptFilterHarness.allGroupId,
      name: "All",
      isAllGroup: true,
      members: [
        buildGroupMember(
            userId: 10, groupId: ReceiptFilterHarness.allGroupId),
      ],
    ),
    buildGroup(
      id: ReceiptFilterHarness.householdId,
      name: "Household",
      members: [
        buildGroupMember(
            userId: 10, groupId: ReceiptFilterHarness.householdId),
        buildGroupMember(
            userId: 11, groupId: ReceiptFilterHarness.householdId),
      ],
    ),
    buildGroup(id: ReceiptFilterHarness.officeId, name: "Office"),
  ]);

  categoryModel.setGroupCategories({
    ReceiptFilterHarness.householdId: householdCategories,
  });
  tagModel.setGroupTags({
    ReceiptFilterHarness.householdId: householdTags,
  });

  return ReceiptFilterHarness(
    receiptListModel: ReceiptListModel(),
    groupModel: groupModel,
    categoryModel: categoryModel,
    tagModel: tagModel,
    userModel: userModel,
  );
}

/// Pumps [child] under every provider the filter surfaces read.
Future<void> pumpWithFilterHarness(
  WidgetTester tester,
  ReceiptFilterHarness harness,
  Widget child,
) async {
  await tester.pumpWidget(MultiProvider(
    providers: [
      ChangeNotifierProvider<ReceiptListModel>.value(
          value: harness.receiptListModel),
      ChangeNotifierProvider<GroupModel>.value(value: harness.groupModel),
      ChangeNotifierProvider<CategoryModel>.value(value: harness.categoryModel),
      ChangeNotifierProvider<TagModel>.value(value: harness.tagModel),
      ChangeNotifierProvider<UserModel>.value(value: harness.userModel),
      ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
      ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
      ChangeNotifierProvider<SystemSettingsModel>(
          create: (_) => SystemSettingsModel()),
    ],
    child: MaterialApp(home: child),
  ));
  await tester.pump();
}
