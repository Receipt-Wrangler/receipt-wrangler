import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/user_avatar.dart';
import 'package:receipt_wrangler_mobile/utils/currency.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../client/client.dart';
import '../constants/text_styles.dart';

class GroupSummary extends StatefulWidget {
  const GroupSummary({super.key, required this.dashboardWidget});

  final api.Widget dashboardWidget;

  @override
  State<GroupSummary> createState() => _GroupSummary();
}

class _GroupSummary extends State<GroupSummary> {
  String _getUserOwesText(
      MapEntry<String, String> mapEntry, UserModel userModel) {
    var user = userModel.getUserById(mapEntry.key);
    var formattedAmount = formatCurrency(context, mapEntry.value);
    var value = formattedAmount.toString().replaceAll("-", "");

    if (mapEntry.value.contains("-") || mapEntry.value == "0") {
      return "${user!.displayName} owes you: $formattedAmount";
    } else {
      return "You owe ${user!.displayName}: $value";
    }
  }

  List<Widget> buildSummaryLineWidgets(Map<String, String>? userData) {
    var nothingOwed = [const Text("Phew, you're all caught up!")];
    var widgets = <Widget>[];
    if (userData == null || userData.isEmpty) {
      return nothingOwed;
    }
    var userModel = Provider.of<UserModel>(context, listen: false);

    userData.entries.forEach((element) {
      widgets.add(Wrap(
        children: [
          UserAvatar(userId: element.key),
          const SizedBox(width: 10),
          Text(_getUserOwesText(element, userModel)),
        ],
      ));
      widgets.add(const SizedBox(height: 10));
    });

    return widgets;
  }

  @override
  Widget build(BuildContext context) {
    var groupId = int.tryParse(getGroupId(context) ?? "");
    var groupSummaryFuture =
        OpenApiClient.client.getUserApi().getAmountOwedForUser(
              groupId: groupId ?? 0,
            );

    return FutureBuilder(
        future: groupSummaryFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            return Column(
              mainAxisAlignment: MainAxisAlignment.start,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const SizedBox(height: 10),
                Text(
                  widget.dashboardWidget.name ?? "",
                  style: dashboardWidgetNameStyle,
                ),
                const SizedBox(height: 10),
                ...buildSummaryLineWidgets(snapshot.data?.data?.toMap() ?? {}),
              ],
            );
          }
          return const CircularProgressIndicator();
        });
  }
}
