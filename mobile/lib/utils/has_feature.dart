import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';

bool hasAiPoweredReceipts(BuildContext context) {
  return Provider.of<AuthModel>(context, listen: false)
      .featureConfig
      .aiPoweredReceipts;
}
