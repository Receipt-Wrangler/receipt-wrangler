import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';
import 'package:receipt_wrangler_mobile/utils/auth.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

import '../../../models/system_settings_model.dart';
import '../../../utils/currency.dart';

class AuthForm extends StatefulWidget {
  const AuthForm({super.key});

  @override
  State<AuthForm> createState() => _Login();
}

class _Login extends State<AuthForm> {
  final _formKey = GlobalKey<FormBuilderState>();
  late final screenSize = MediaQuery.of(context).size;
  bool isSignUp = false;
  bool isLoading = false;

  Future<void> _submit() async {
    _formKey.currentState!.save();

    if (_formKey.currentState!.validate()) {
      if (isSignUp) {
        await signUp();
      } else {
        await login();
      }
    }
  }

  Future<void> login() async {
    var form = _formKey.currentState!.value;
    var command = (api.LoginCommandBuilder()
          ..username = form["username"]
          ..password = form["password"])
        .build();
    try {
      toggleIsLoading();
      var appDataResponse =
          await OpenApiClient.client.getAuthApi().login(loginCommand: command);
      await _onLoginSuccess(appDataResponse.data as api.AppData);
    } catch (e) {
      toggleIsLoading();
      print(e);
      showApiErrorSnackbar(context, e as dynamic);
    }
  }

  Future<void> signUp() async {
    var form = _formKey.currentState!.value;
    toggleIsLoading();
    OpenApiClient.client
        .getAuthApi()
        .signUp(
            signUpCommand: (api.SignUpCommandBuilder()
                  ..username = form["username"]
                  ..password = form["password"]
                  ..displayName = form["displayName"])
                .build())
        .then((data) async => {await login()})
        .catchError((err) => {
              toggleIsLoading(),
              print(err),
              showApiErrorSnackbar(context, err),
            });
  }

  void toggleIsLoading() {
    setState(() {
      isLoading = !isLoading;
    });
  }

  Future<void> _onLoginSuccess(api.AppData appData) async {
    var authModel = Provider.of<AuthModel>(context, listen: false);
    var groupModel = Provider.of<GroupModel>(context, listen: false);
    var userModel = Provider.of<UserModel>(context, listen: false);
    var userPreferencesModel =
        Provider.of<UserPreferencesModel>(context, listen: false);
    var categoryModel = Provider.of<CategoryModel>(context, listen: false);
    var tagModel = Provider.of<TagModel>(context, listen: false);
    var systemSettingsModel =
        Provider.of<SystemSettingsModel>(context, listen: false);

    await storeAppData(authModel, groupModel, userModel, userPreferencesModel,
        categoryModel, tagModel, systemSettingsModel, appData);
    registerCustomCurrency(context);
    context.go("/groups");
  }

  Widget _getDisplaynameField() {
    if (isSignUp) {
      return FormBuilderTextField(
          name: "displayName",
          decoration: const InputDecoration(
              labelText: "Displayname", border: OutlineInputBorder()),
          validator: FormBuilderValidators.compose([
            FormBuilderValidators.required(),
          ]));
    } else {
      return const SizedBox.shrink();
    }
  }

  Widget _getSignUpButton() {
    var buttonText = "";

    if (isSignUp) {
      buttonText = "Return to Login";
    } else {
      buttonText = "Create an Account";
    }

    return CupertinoButton(
        onPressed: () {
          setState(() {
            isSignUp = !isSignUp;
          });
        },
        child: Text(buttonText));
  }

  Widget _getSubmitButtonText() {
    if (isSignUp) {
      return const Text("Create an Account");
    } else {
      return const Text("Log In");
    }
  }

  Widget _getServerInfoText(AuthModel server) {
    if (isSignUp) {
      return Text('Signing up on: ${server.basePath}');
    }
    return Text('Logging into: ${server.basePath}');
  }

  Widget _getChangeServerButton() {
    return CupertinoButton(
        onPressed: () {
          context.go("/");
        },
        child: const Text("Change Server"));
  }

  @override
  Widget build(BuildContext context) {
    return AutofillGroup(
        child: FormBuilder(
      key: _formKey,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          if (isLoading) ...[
            const CircularProgressIndicator(),
          ] else ...[
            SvgPicture.asset(
              "assets/branding/logo-large.svg",
              width: screenSize.width * 0.25,
              height: screenSize.width * 0.25,
            ),
            const SizedBox(
              height: 10,
            ),
            Consumer<AuthModel>(
              builder: (context, auth, child) {
                return _getServerInfoText(auth);
              },
            ),
            headerSpacing,
            _getDisplaynameField(),
            isSignUp ? textFieldSpacing : const SizedBox.shrink(),
            FormBuilderTextField(
                name: "username",
                autofillHints: const [AutofillHints.username],
                decoration: const InputDecoration(
                    labelText: "Username", border: OutlineInputBorder()),
                validator: FormBuilderValidators.compose([
                  FormBuilderValidators.required(),
                ])),
            textFieldSpacing,
            FormBuilderTextField(
                name: "password",
                autofillHints: const [AutofillHints.password],
                obscureText: true,
                decoration: const InputDecoration(
                    labelText: "Password", border: OutlineInputBorder()),
                validator: FormBuilderValidators.compose([
                  FormBuilderValidators.required(),
                ])),
            lastFieldSpacing,
            Row(
              children: [
                Expanded(
                  child: CupertinoButton.filled(
                      onPressed: () {
                        _submit();
                      },
                      child: _getSubmitButtonText()),
                )
              ],
            ),
            Consumer<AuthModel>(
              builder: (context, auth, child) {
                if (auth.featureConfig.enableLocalSignUp) {
                  return _getSignUpButton();
                } else {
                  return const SizedBox.shrink();
                }
              },
            ),
            _getChangeServerButton(),
          ],
        ],
      ),
    ));
  }
}
