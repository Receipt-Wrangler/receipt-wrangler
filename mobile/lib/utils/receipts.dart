import 'dart:convert';
import 'dart:typed_data';

import 'package:built_collection/built_collection.dart';
import 'package:built_value/json_object.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/utils/forms.dart';

String? getReceiptId(BuildContext context) {
  return GoRouterState.of(context).pathParameters["receiptId"];
}

Receipt getDefaultReceipt() {
  return (ReceiptBuilder()
        ..id = 0
        ..date = DateTime.now().toIso8601String()
        ..groupId = 0
        ..paidByUserId = 0
        ..amount = "0"
        ..name = ""
        ..status = ReceiptStatus.OPEN
        ..categories = ListBuilder<Category>([])
        ..tags = ListBuilder<Tag>([]))
      .build();
}

String getTitleText(WranglerFormState formState, String receiptName) {
  return "${getFormStateHeader(formState)} $receiptName Receipt";
}

String getLeadingArrowRedirect(String groupId) {
  return "/groups/$groupId/receipts";
}

double getImagePreviewWidth(BuildContext context) {
  return MediaQuery.of(context).size.width;
}

double getImagePreviewHeight(BuildContext context) {
  return MediaQuery.of(context).size.height * .5;
}

EdgeInsets getImageDataPadding() {
  return const EdgeInsets.all(26);
}

Uint8List getBytesFromEncodedImage(String encodedImage) {
  var base64Image = encodedImage.split(",").last;
  if (base64Image == "") {
    return Uint8List(0);
  }

  var bytes = base64Decode(base64Image);
  return bytes;
}

ReceiptPagedRequestFilterBuilder dashboardConfigurationToFilter(
    BuiltMap<String, JsonObject?>? configuration) {
  var filter = (ReceiptPagedRequestFilterBuilder());
  if (configuration == null) {
    return filter;
  }

  configuration.forEach((key, value) {
    if (value.toString().length == 0) {
      return;
    }

    var valueMap =
        (value?.asMap as Map<String, dynamic>) ?? Map<String, dynamic>();
    var filterObject = JsonObject({
      "operation": valueMap["operation"] ?? "",
      "value": valueMap["value"] ?? null,
    });

    switch (key) {
      case 'date':
        filter.date = filterObject;
        break;
      case 'amount':
        filter.amount = filterObject;
        break;
      case 'name':
        filter.name = filterObject;
        break;
      case 'paidBy':
        filter.paidBy = filterObject;
        break;
      case 'categories':
        filter.categories = filterObject;
        break;
      case 'tags':
        filter.tags = filterObject;
        break;
      case 'status':
        filter.status = filterObject;
        break;
      case 'resolvedDate':
        filter.resolvedDate = filterObject;
        break;
      case 'createdAt':
        filter.createdAt = filterObject;
        break;
    }
  });

  return filter;
}
