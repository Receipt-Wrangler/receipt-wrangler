import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/constants/font.dart';
import 'package:receipt_wrangler_mobile/shared/functions/activities.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

import '../../constants/colors.dart';
import '../../extensions/duration.dart';
import '../../models/auth_model.dart';
import '../../models/group_model.dart';
import '../../models/user_model.dart';
import '../../shared/functions/permissions.dart';
import '../../shared/widgets/list_item_lead.dart';
import '../../shared/widgets/list_item_trailing_status.dart';

class GroupActivityListItem extends StatefulWidget {
  const GroupActivityListItem(
      {super.key, required this.activity, required this.groupId});

  final api.Activity activity;

  final int groupId;

  @override
  State<GroupActivityListItem> createState() => _GroupActivityListItem();
}

class _GroupActivityListItem extends State<GroupActivityListItem> {
  late final authModel = Provider.of<AuthModel>(context, listen: false);
  late final userModel = Provider.of<UserModel>(context, listen: false);
  late final groupModel = Provider.of<GroupModel>(context, listen: false);
  var hasBeenRerun = false;

  Color getActivityColor() {
    switch (widget.activity.status) {
      case api.SystemTaskStatus.SUCCEEDED:
        return successGreen;
      case api.SystemTaskStatus.FAILED:
        return errorRed;
      default:
        return Colors.grey;
    }
  }

  Widget getLeadingWidget() {
    return ListItemLead(
      date: widget.activity.startedAt,
      color: getActivityColor(),
    );
  }

  Widget getTrailingWidget() {
    return Column(
      children: [
        Text(widget.activity.startedAt.toDurationString()),
        ListItemTrailingStatus(
          color: getActivityColor(),
          text: getActivityStatusDisplay(widget.activity.status),
          height: 40,
        )
      ],
    );
  }

  Widget buildSubtitleWidget() {
    var displayName = "System";
    if (widget.activity.ranByUserId != null) {
      var user = userModel.getUserById(widget.activity.ranByUserId.toString());
      displayName = user?.displayName ?? "Unknown";
    }

    var group = groupModel.getGroupById(widget.groupId.toString());
    if (group?.isAllGroup ?? false) {
      var taskGroup =
          groupModel.getGroupById(widget.activity.groupId.toString());

      return Text("Ran by ${displayName} in ${taskGroup?.name ?? 'Unknown'} ");
    }

    return Text("Ran by ${displayName}");
  }

  Widget buildListTile() {
    return ListTile(
      leading: getLeadingWidget(),
      subtitle: buildSubtitleWidget(),
      title: Text(
        getActivityTypeDisplay(widget.activity.type),
        style: boldText,
      ),
      trailing: getTrailingWidget(),
      onTap: () => {
        if (widget.activity.receiptId != null)
          {context.go('/receipts/${widget.activity.receiptId}/view')}
      },
    );
  }

  Widget buildSlideableAction() {
    return SlidableAction(
      icon: Icons.refresh,
      label: "Rerun",
      foregroundColor: Theme.of(context).colorScheme.primary,
      onPressed: (BuildContext context) async {
        await rerunActivity();
      },
    );
  }

  Future<void> rerunActivity() async {
    try {
      await OpenApiClient.client
          .getSystemTaskApi()
          .rerunActivity(id: widget.activity.id);

      showSuccessSnackbar(context, "Activity has been successfully queued.");
      setState(() {
        hasBeenRerun = true;
      });
    } catch (e) {
      showApiErrorSnackbar(context, e as dynamic);
    }
  }

  @override
  Widget build(BuildContext context) {
    var canEdit = canEditReceipt(authModel, groupModel, widget.groupId);
    var slideEnabled =
        canEdit && (widget.activity.canBeRestarted ?? false) && !hasBeenRerun;

    return SlidableWidget(
        slideEnabled: slideEnabled,
        endActionPaneChildren: [buildSlideableAction()],
        slidableChild: buildListTile());
  }
}
