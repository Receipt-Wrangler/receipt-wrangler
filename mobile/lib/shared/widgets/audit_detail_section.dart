import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/shared/functions/forms.dart';

import '../../models/user_model.dart';
import '../../utils/date.dart';

class AuditDetailSection extends StatefulWidget {
  const AuditDetailSection({
    super.key,
    required this.entity,
  });

  final dynamic entity;

  @override
  _AuditDetailSectionState createState() {
    return _AuditDetailSectionState();
  }
}

class _AuditDetailSectionState extends State<AuditDetailSection> {
  late final userModel = Provider.of<UserModel>(context, listen: true);
  final rowSpacer = SizedBox(height: 2);

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

  Widget buildAddedByRow() {
    var id = widget.entity.createdBy;
    if (id == null || id == 0) {
      return SizedBox.shrink();
    }

    var user = userModel.getUserById(widget.entity.createdBy.toString());

    return Row(
      children: [
        Text("Added by: ${user?.displayName ?? ""}"),
      ],
    );
  }

  Widget buildAddedAtRow() {
    var createdAt = widget.entity.createdAt;
    if (createdAt == null) {
      return SizedBox.shrink();
    }

    var formattedDate =
        formatDate(defaultLongDateFormat, DateTime.parse(createdAt));

    return Row(
      children: [
        Text("Added at: ${formattedDate}"),
      ],
    );
  }

  Widget buildUpdatedAtRow() {
    var updatedAt = widget.entity.updatedAt;
    if (updatedAt == null) {
      return SizedBox.shrink();
    }

    var formattedDate =
        formatDate(defaultLongDateFormat, DateTime.parse(updatedAt));

    return Row(
      children: [
        Text("Updated at: ${formattedDate}"),
      ],
    );
  }

  Widget buildResolvedAtRow() {
    var resolvedAt = widget.entity.resolvedDate;
    if (resolvedAt == null) {
      return SizedBox.shrink();
    }

    var formattedDate =
        formatDate(defaultLongDateFormat, DateTime.parse(resolvedAt));

    return Row(
      children: [
        Text("Resolved at: ${formattedDate}"),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Row(
          children: [
            buildHeaderText("Audit Details"),
          ],
        ),
        const SizedBox(height: 10),
        buildAddedByRow(),
        rowSpacer,
        buildAddedAtRow(),
        rowSpacer,
        buildUpdatedAtRow(),
        rowSpacer,
        buildResolvedAtRow(),
      ],
    );
  }
}
