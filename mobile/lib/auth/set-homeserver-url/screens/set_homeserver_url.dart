import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

class SetHomeserverUrl extends StatefulWidget {
  const SetHomeserverUrl({super.key});

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

  @override
  Widget build(BuildContext context) {
    var serverModel = Provider.of<AuthModel>(context);

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
              decoration: const InputDecoration(
                  hintText: "https://demo.receiptwrangler.io/api",
                  labelText: "Server URL",
                  border: OutlineInputBorder()),
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
