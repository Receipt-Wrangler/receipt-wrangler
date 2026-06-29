import 'package:dio/dio.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/services/token_refresh_service.dart';
import 'package:receipt_wrangler_mobile/utils/auth.dart';

/// Dio interceptor that recovers an *expired-session* request by refreshing the
/// token via [TokenRefreshService] and retrying once.
///
/// The backend returns **403 for every access denial** — an expired session and
/// a genuine permission denial both surface as 403 (it never sends 401). So the
/// status code alone can't tell the two apart; we distinguish by **token
/// validity**, mirroring the desktop's http-interceptor.ts: a 403 with a
/// still-valid token is a permission denial and must pass through untouched
/// (refreshing can't grant a missing permission, and force-refreshing here would
/// burn the one-time refresh token and risk logging the user out). Only an
/// expired/invalid token is an auth failure worth a refresh + retry. Token
/// freshness is otherwise kept current proactively (the 15-min timer in
/// `main.dart` and the auth guard on navigation).
class AuthInterceptor extends Interceptor {
  static const _retryHeader = 'X-Token-Retry';

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    final statusCode = err.response?.statusCode;

    // Don't intercept token refresh requests — let TokenRefreshService handle those errors
    if (err.requestOptions.path.contains('/token/')) {
      return handler.next(err);
    }

    // Don't retry if we already retried this request
    if (err.requestOptions.headers.containsKey(_retryHeader)) {
      return handler.next(err);
    }

    if (statusCode == 401 || statusCode == 403) {
      // A 403 with a still-valid token is a permission denial, not an auth
      // failure — pass it through so the caller can handle it, without burning
      // the one-time refresh token or risking a logout. Only an expired/invalid
      // token warrants a refresh + retry.
      final currentJwt = await TokenRefreshService().getCurrentJwt();
      if (!isTokenValid(currentJwt)) {
        try {
          final success =
              await TokenRefreshService().refreshTokens(force: true);
          if (success) {
            final jwt = await TokenRefreshService().getCurrentJwt();
            final opts = err.requestOptions;
            opts.headers[_retryHeader] = 'true';
            if (jwt != null) {
              opts.headers['Authorization'] = 'Bearer $jwt';
            }

            // Retry the request using the current client's Dio instance
            final retryResponse = await OpenApiClient.client.dio.fetch(opts);
            return handler.resolve(retryResponse);
          }
        } catch (_) {
          // Refresh failed — fall through to original error
        }
      }
    }

    return handler.next(err);
  }
}
