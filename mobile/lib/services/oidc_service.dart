import 'package:flutter/services.dart';
import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/utils/pkce.dart';

/// The app's private-use URL scheme, reverse-DNS of a domain the project owns
/// per RFC 8252 section 7.1, and matching the Android applicationId.
///
/// It carries a one-time, PKCE-bound code -- never a token. A private-use scheme
/// is unverifiable on Android, so any installed app can register it and read the
/// redirect; the verifier this app kept is what makes that useless.
const oidcCallbackScheme = 'io.receiptwrangler';

/// The full callback URL the backend redirects to.
const oidcCallbackUrl = '$oidcCallbackScheme://oidc';

/// Raised when a sign-in attempt fails for a reason worth showing the user.
class OidcSignInException implements Exception {
  OidcSignInException(this.message);

  final String message;

  @override
  String toString() => message;
}

/// Raised when the user backs out of the browser. Callers stay silent on this.
class OidcSignInCancelled implements Exception {}

/// Builds the URL that starts a sign-in for [providerName].
///
/// [basePath] already ends in `/api` (it is the server URL the user entered on
/// the Connect screen), so no extra path juggling is needed.
String buildOidcLoginUrl(
  String basePath,
  String providerName,
  String codeChallenge,
) {
  final normalizedBase =
      basePath.endsWith('/') ? basePath.substring(0, basePath.length - 1) : basePath;

  return '$normalizedBase/oidc/${Uri.encodeComponent(providerName)}/login'
      '?client=mobile&codeChallenge=${Uri.encodeQueryComponent(codeChallenge)}';
}

/// Extracts the one-time code from the callback URL, or throws with the
/// backend's error if it sent one.
String extractCodeFromCallback(String callbackUrl) {
  final uri = Uri.parse(callbackUrl);

  final error = uri.queryParameters['error'];
  if (error != null && error.isNotEmpty) {
    throw OidcSignInException(oidcErrorMessage(error));
  }

  final code = uri.queryParameters['code'];
  if (code == null || code.isEmpty) {
    throw OidcSignInException('Sign in failed. Please try again.');
  }

  return code;
}

/// Copy for the error codes the backend can redirect back with. It only ever
/// sends codes from a small fixed vocabulary -- an identity provider's own error
/// text is never echoed through -- so an unknown code still gets a sane message.
String oidcErrorMessage(String code) {
  switch (code) {
    case 'unknown_provider':
      return 'That sign-in provider is not available.';
    case 'no_account':
      return 'No Receipt Wrangler account is linked to that identity. '
          'Sign in and connect it from your profile, or ask an administrator.';
    case 'account_exists':
      return 'An account with that username already exists. Sign in with your '
          'password and connect the provider from your profile.';
    case 'already_linked':
      return 'That identity is already connected to another account.';
    case 'invalid_state':
    case 'nonce_mismatch':
      return 'That sign-in attempt expired. Please try again.';
    case 'no_id_token':
      return 'The provider did not return an identity token.';
    case 'provider_error':
      return 'The sign-in provider reported an error.';
    default:
      return 'Sign in failed. Please try again.';
  }
}

/// Runs a full OIDC sign-in and returns the resulting [api.AppData].
///
/// The whole OIDC exchange happens on the server; this app only opens a browser
/// and redeems the code that comes back. The response carries `jwt` and
/// `refreshToken` exactly as `login?tokensInBody=true` does, so the caller hands
/// it straight to the existing storeAppData path.
Future<api.AppData> signInWithOidc({
  required String basePath,
  required String providerName,
}) async {
  final verifier = generateCodeVerifier();
  final challenge = codeChallengeS256(verifier);

  String callbackUrl;
  try {
    callbackUrl = await FlutterWebAuth2.authenticate(
      url: buildOidcLoginUrl(basePath, providerName, challenge),
      callbackUrlScheme: oidcCallbackScheme,
    );
  } on PlatformException catch (e) {
    // The user dismissed the browser. Not an error worth a snackbar.
    if (e.code == 'CANCELED' || e.code == 'canceled') {
      throw OidcSignInCancelled();
    }

    rethrow;
  }

  final code = extractCodeFromCallback(callbackUrl);

  final response = await OpenApiClient.client.getOidcApi().oidcExchange(
        oidcExchangeCommand: (api.OidcExchangeCommandBuilder()
              ..code = code
              ..codeVerifier = verifier)
            .build(),
      );

  final appData = response.data;
  if (appData == null) {
    throw OidcSignInException('Sign in failed. Please try again.');
  }

  return appData;
}
