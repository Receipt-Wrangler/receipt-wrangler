import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/auth/set-homeserver-url/screens/qr_scanner_screen.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';
import 'package:receipt_wrangler_mobile/utils/url.dart';

class SetHomeserverUrl extends StatefulWidget {
  const SetHomeserverUrl({super.key, this.scanQrCode});

  /// Injectable for tests: returns the raw scanned string, or null if the user
  /// cancelled. Defaults to opening [QrScannerScreen].
  final Future<String?> Function(BuildContext context)? scanQrCode;

  @override
  State<SetHomeserverUrl> createState() => _SetHomeserverUrl();
}

class _SetHomeserverUrl extends State<SetHomeserverUrl> {
  final _formKey = GlobalKey<FormBuilderState>();

  Future<void> _submit() async {
    if (_formKey.currentState!.validate()) {
      _formKey.currentState!.save();
      var authModel = Provider.of<AuthModel>(context, listen: false);

      await authModel.setBasePath(_formKey.currentState!.value["url"]);
      try {
        var featureConfig =
            await OpenApiClient.client.getFeatureConfigApi().getFeatureConfig();
        authModel.setFeatureConfig(featureConfig.data);
        showSuccessSnackbar(context, "Successfully connected to server");
        context.go("/login");
      } catch (e) {
        showErrorSnackbar(context, "Failed to connect to server");
      }
    }
  }

  Future<void> _onScanPressed() async {
    final scan = widget.scanQrCode ?? _openScanner;
    final raw = await scan(context);
    if (!mounted || raw == null) {
      return; // null => user cancelled
    }

    // Resolve the deep-link QR (fragment-encoded server url) first, falling back
    // to treating the scan as a plain server-URL QR. Both fill the field only —
    // the user still reviews and taps Connect (no auto-connect).
    final normalized =
        extractDeepLinkServerUrl(raw) ?? normalizeServerUrl(raw);
    if (normalized == null) {
      showErrorSnackbar(
          context, "That QR code doesn't contain a valid server URL");
      return;
    }

    _formKey.currentState?.patchValue({"url": normalized});
  }

  Future<String?> _openScanner(BuildContext context) {
    return Navigator.of(context).push<String>(
      MaterialPageRoute(builder: (_) => const QrScannerScreen()),
    );
  }

  @override
  Widget build(BuildContext context) {
    var serverModel = Provider.of<AuthModel>(context);

    // A server URL arriving from a receiptwrangler.io/app/setup deep link is
    // stashed on AuthModel.pendingServerUrl by the deep-link handler in
    // main.dart. Consume it here to pre-fill the field. This covers both the
    // cold-start case (value already present when this widget first mounts) and
    // the warm case (value arrives later -> this listener fires -> rebuild). The
    // patch runs post-frame because the FormBuilder state isn't attached yet
    // during the first build. We never auto-connect (phishing mitigation).
    final pending = serverModel.pendingServerUrl;
    if (pending != null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted) {
          return;
        }
        _formKey.currentState?.patchValue({"url": pending});
        serverModel.clearPendingServerUrl();
      });
    }

    return FormBuilder(
      key: _formKey,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          const Text(
            "Connect to Server",
            style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
          ),
          headerSpacing,
          FormBuilderTextField(
              name: "url",
              decoration: InputDecoration(
                  hintText: "https://demo.receiptwrangler.io/api",
                  labelText: "Server URL",
                  border: const OutlineInputBorder(),
                  suffixIcon: IconButton(
                    key: const ValueKey("qr-scan-button"),
                    icon: const Icon(Icons.qr_code_scanner),
                    tooltip: "Scan server QR code",
                    onPressed: _onScanPressed,
                  )),
              initialValue: serverModel.basePath,
              validator: FormBuilderValidators.compose([
                FormBuilderValidators.required(),
              ])),
          lastFieldSpacing,
          Row(
            children: [
              Expanded(
                  child: CupertinoButton.filled(
                      onPressed: () async {
                        await _submit();
                      },
                      child: const Text("Connect")))
            ],
          ),
        ],
      ),
    );
  }
}
