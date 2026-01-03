import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/custom_field_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/receipts/nav/receipt_app_bar.dart';
import 'package:receipt_wrangler_mobile/receipts/nav/receipt_app_bar_action_builder.dart';
import 'package:receipt_wrangler_mobile/receipts/nav/receipt_bottom_sheet_builder.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_form.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/circular_loading_progress.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

class ReceiptFormScreen extends StatefulWidget {
  const ReceiptFormScreen({super.key});

  @override
  State<ReceiptFormScreen> createState() => _ReceiptFormScreen();
}

class _ReceiptFormScreen extends State<ReceiptFormScreen> {
  late final Future<List<dynamic>> _coordinatedFuture = _loadData();

  Future<List<dynamic>> _loadData() async {
    var customFieldModel = Provider.of<CustomFieldModel>(context, listen: false);
    var state = GoRouterState.of(context);
    var receiptId = state.pathParameters['receiptId'] ?? "0";

    Future<void> customFieldsFuture = Future.value(null);
    if (!customFieldModel.isLoading) {
      customFieldsFuture = customFieldModel.loadCustomFields().catchError((error) {
        debugPrint('Custom fields loading failed: $error');
        return null; // Continue with empty custom fields
      });
    }

    Future<api.Receipt?> receiptFuture = Future.value(null);
    if (receiptId != "0") {
      receiptFuture = OpenApiClient.client
          .getReceiptApi()
          .getReceiptById(
              receiptId: int.parse(state.pathParameters['receiptId'] as String))
          .then((response) => response.data)
          .catchError((error) {
        debugPrint('Receipt loading failed: $error');
        throw error; // Receipt failure should trigger error state
      });
    }

    // Combine both futures - wait for both to complete
    // Custom field errors are handled gracefully, receipt errors bubble up
    return Future.wait([
      customFieldsFuture,
      receiptFuture,
    ]);
  }

  @override
  Widget build(BuildContext context) {
    EdgeInsets? padding;

    var contextModel = Provider.of<ContextModel>(context, listen: false);
    var receiptModel = Provider.of<ReceiptModel>(context, listen: false);
    var customFieldModel = Provider.of<CustomFieldModel>(context, listen: false);

    var state = GoRouterState.of(context);
    var actionBuilder = ReceiptAppBarActionBuilder(context, receiptModel);
    var bottomSheetBuilder = ReceiptBottomSheetBuilder(context, receiptModel);

    contextModel.setShellContext(context);

    return FutureBuilder<List<dynamic>>(
        future: _coordinatedFuture,
        builder: (context, snapshot) {
          var isReady = snapshot.connectionState == ConnectionState.done;

          // Extract receipt data from combined result
          api.Receipt? receiptData;
          if (isReady && snapshot.hasData && snapshot.data!.length > 1) {
            receiptData = snapshot.data![1] as api.Receipt?;
          }

          if (isReady && receiptData != null) {
            receiptModel.setReceipt(receiptData, false);
          }

          if (isReady && snapshot.hasError) {
            showErrorSnackbar(context, 'Failed to load receipt');
            context.go('/');
          }

          var showChild = isReady && !customFieldModel.isLoading;

          return ScreenWrapper(
            appBarWidget:
                ReceiptAppBar(actions: actionBuilder.buildAppBarMenu(state)),
            bodyPadding: padding,
            bottomSheetWidget: bottomSheetBuilder.buildBottomSheet(state),
            child: showChild
                ? SingleChildScrollView(child: const ReceiptForm())
                : const CircularLoadingProgress(),
          );
        });
  }
}
