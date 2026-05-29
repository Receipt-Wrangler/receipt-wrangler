import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';

/// Test-only interceptor that fails any request whose path matches a
/// pattern. Lets a test simulate a server-side failure for one
/// endpoint while leaving everything else alone (so we can still
/// `createReceipt` successfully and only the image-upload step
/// fails).
class _FailRequestInterceptor extends Interceptor {
  _FailRequestInterceptor(this._pathPattern);

  final RegExp _pathPattern;

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) {
    if (_pathPattern.hasMatch(options.path)) {
      handler.reject(
        DioException(
          requestOptions: options,
          type: DioExceptionType.badResponse,
          response: Response(
            requestOptions: options,
            statusCode: 500,
            data: <String, String>{
              'errorMsg': 'Simulated failure (test injector)',
            },
          ),
        ),
      );
      return;
    }
    handler.next(options);
  }
}

/// Forces every receipt-image upload (POST `/receiptImage/...`) to fail
/// with HTTP 500. Used to verify bug-fix #2 (`addReceipt`'s partial
/// upload reporting): the receipt itself is created server-side, but
/// the image upload Future.wait throws, exercising the partial-failure
/// snackbar path.
///
/// Must be called AFTER `loginAsAdmin` -- the production
/// AuthModel._updateDefaultApiClient rebuilds the dio instance on
/// login, so any earlier interceptor would be discarded.
///
/// Cleans itself up via [addTearDown] -- the interceptor is removed
/// from the active dio when the test finishes.
void installFailReceiptImageUpload() {
  final interceptor =
      _FailRequestInterceptor(RegExp(r'/receiptImage(/|$)'));
  // Insert at index 0 so it runs BEFORE the AuthInterceptor (which
  // adds the Bearer header). Order doesn't strictly matter for the
  // reject path, but front-of-chain makes the failure unambiguous.
  OpenApiClient.client.dio.interceptors.insert(0, interceptor);
  addTearDown(() {
    OpenApiClient.client.dio.interceptors.remove(interceptor);
  });
}
