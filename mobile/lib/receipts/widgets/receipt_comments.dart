import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_delete_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

import '../../client/client.dart';
import '../../models/receipt_model.dart';
import '../../models/user_model.dart';
import '../../shared/widgets/user_avatar.dart';
import '../../utils/date.dart';

class ReceiptComments extends StatefulWidget {
  const ReceiptComments({super.key, required this.comments, required this.formState});

  final List<api.Comment> comments;
  final WranglerFormState formState;

  @override
  State<ReceiptComments> createState() => _ReceiptComments();
}

class _ReceiptComments extends State<ReceiptComments> {
  late final formState = widget.formState;

  Widget buildWidgetList() {
    var slideEnabled = formState == WranglerFormState.edit ||
        formState == WranglerFormState.add;

    return ListView.builder(
      itemCount: widget.comments.length,
      padding: formState == WranglerFormState.view
          ? null
          : EdgeInsets.only(bottom: 60),
      itemBuilder: (context, index) {
        return SlidableWidget(
            slideEnabled: slideEnabled,
            endActionPaneChildren: [buildDeleteButton(index)],
            slidableChild: Column(
              children: [
                buildCommentRow(widget.comments[index], index),
                SizedBox(height: 10),
              ],
            ));
      },
    );
  }

  Widget buildDeleteButton(int index) {
    return SlidableDeleteButton(onPressed: () => deleteComment(index));
  }

  void deleteComment(int index) {
    if (formState == WranglerFormState.add) {
      var comments = new List<api.Comment>.from(widget.comments);
      comments.removeAt(index);

      Provider.of<ReceiptModel>(context, listen: false).setComments(comments);
    }

    if (formState == WranglerFormState.edit) {
      var commentId = widget.comments[index].id;
      OpenApiClient.client
          .getCommentApi()
          .deleteComment(commentId: commentId)
          .then((value) {
        var receiptModel = Provider.of<ReceiptModel>(context, listen: false);
        var comments = new List<api.Comment>.from(receiptModel.comments);
        comments.removeAt(index);

        receiptModel.setComments(comments);
      }).catchError((error) {
        showApiErrorSnackbar(context, error);
      });
    }
  }

  Widget buildCommentRow(api.Comment comment, int index) {
    var lastCommentHasSameUser = false;
    var commentDate = formatDate(
        "MMM d, h:mm a", DateTime.parse(comment.createdAt ?? "").toLocal());
    var spacerWidget = SizedBox.shrink();
    var user = Provider.of<UserModel>(context, listen: false)
        .getUserById(comment.userId.toString());

    var isLoggedInUsersComment = user?.id ==
        Provider.of<AuthModel>(context, listen: false).claims?.userId;

    if (index > 0 && widget.comments[index - 1].userId == comment.userId) {
      lastCommentHasSameUser = true;
    }

    if (lastCommentHasSameUser && !isLoggedInUsersComment) {
      spacerWidget = SizedBox(width: 55);
    }

    return Row(
      mainAxisAlignment: isLoggedInUsersComment
          ? MainAxisAlignment.end
          : MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        lastCommentHasSameUser || isLoggedInUsersComment
            ? spacerWidget
            : Center(
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Center(
                    child: UserAvatar(
                      userId: user?.id.toString(),
                    ),
                  ),
                ),
              ),
        SizedBox(width: 10),
        Container(
          width: MediaQuery.of(context).size.width * 0.50,
          padding: EdgeInsets.all(10),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(5),
            color: isLoggedInUsersComment
                ? Theme.of(context).primaryColorLight
                : Colors.grey[200],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              lastCommentHasSameUser
                  ? SizedBox.shrink()
                  : Text(user?.displayName ?? "",
                      style: TextStyle(fontWeight: FontWeight.bold)),
              lastCommentHasSameUser ? SizedBox.shrink() : SizedBox(height: 5),
              Text(comment.comment.trim()),
              Text(commentDate, style: TextStyle(fontSize: 10)),
            ],
          ),
        ),
      ],
    );
  }

  Widget buildCommentWidgets(bool hasComments) {
    return hasComments
        ? buildWidgetList()
        : Center(child: Text("No comments found"));
  }

  @override
  Widget build(BuildContext context) {
    var hasComments = widget.comments.length > 0;
    if (hasComments) {
      return buildCommentWidgets(hasComments);
    } else {
      return buildCommentWidgets(hasComments);
    }
  }
}
