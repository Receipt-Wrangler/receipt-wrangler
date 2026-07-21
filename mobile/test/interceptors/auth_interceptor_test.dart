import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/interceptors/auth_interceptor.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/services/token_refresh_service.dart';

import '../helpers/auth_test_helpers.dart';

/// Tracks which handler method was called by the interceptor.
/// Uses noSuchMethod to satisfy _BaseHandler mixin requirements
/// without triggering Dio's internal Completer plumbing.
class TestErrorHandler with _NoopBaseHandler
    implements ErrorInterceptorHandler {
  final Completer<String> _actionCompleter = Completer<String>();
  Response? resolvedResponse;
  DioException? nextError;

  Future<String> get result => _actionCompleter.future;

  @override
  void next(DioException err) {
    nextError = err;
    if (!_actionCompleter.isCompleted) _actionCompleter.complete('next');
  }

  @override
  void resolve(Response response) {
    resolvedResponse = response;
    if (!_actionCompleter.isCompleted) _actionCompleter.complete('resolve');
  }

  @override
  void reject(DioException err, [bool callFollowingErrorInterceptor = false]) {
    if (!_actionCompleter.isCompleted) _actionCompleter.complete('reject');
  }
}

/// Mixin that provides noSuchMethod to satisfy unimplemented members
/// of Dio's _BaseHandler (future, isCompleted).
mixin _NoopBaseHandler {
  @override
  dynamic noSuchMethod(Invocation invocation) => null;
}

/// Sets up a mock environment where the interceptor can successfully
/// refresh tokens and retry requests. Returns the [Dio] instance so
/// tests can add custom interceptors for assertions.
Dio _setUpSuccessfulRetryScenario({
  required MockAuthModel mockAuthModel,
  required MockGroupModel mockGroupModel,
}) {
  final newJwt = validJwt;
  final newRefresh = validJwt;

  final mockOpenapi = MockOpenapi();
  final mockAuthApi = MockAuthApi();
  final dio = Dio(BaseOptions(baseUrl: 'http://localhost:8081/api'));

  when(() => mockOpenapi.getAuthApi()).thenReturn(mockAuthApi);
  when(() => mockOpenapi.getUserApi()).thenReturn(MockUserApi());
  when(() => mockOpenapi.dio).thenReturn(dio);

  when(() => mockAuthApi.getNewRefreshToken(
          logoutCommand: any(named: 'logoutCommand')))
      .thenAnswer(
          (_) async => createTokenRefreshResponse(newJwt, newRefresh));

  when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
  when(() => mockAuthModel.getRefreshToken())
      .thenAnswer((_) async => validJwt);
  when(() => mockAuthModel.setTokens(any(), any()))
      .thenAnswer((_) async {});
  when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

  var refreshed = false;
  when(() => mockAuthModel.setTokens(any(), any())).thenAnswer((_) async {
    refreshed = true;
  });
  when(() => mockAuthModel.getJwt()).thenAnswer((_) async {
    return refreshed ? newJwt : expiredJwt;
  });

  OpenApiClient.client = mockOpenapi;

  return dio;
}

void main() {
  late MockAuthModel mockAuthModel;
  late MockGroupModel mockGroupModel;
  late AuthInterceptor interceptor;

  setUpAll(() {
    registerFallbackValue(FakeLogoutCommand());
  });

  setUp(() {
    mockAuthModel = MockAuthModel();
    mockGroupModel = MockGroupModel();

    TokenRefreshService().resetForTesting();
    TokenRefreshService().initialize(
      authModel: mockAuthModel,
      groupModel: mockGroupModel,
      userModel: MockUserModel(),
      userPreferencesModel: MockUserPreferencesModel(),
      categoryModel: MockCategoryModel(),
      tagModel: MockTagModel(),
      systemSettingsModel: MockSystemSettingsModel(),
      permissionsModel: PermissionsModel(),
    );

    interceptor = AuthInterceptor();
  });

  group('AuthInterceptor', () {
    test('passes through non-401/403 errors without interception', () async {
      final requestOptions = RequestOptions(path: '/receipts');
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 500,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      final action = await handler.result.timeout(const Duration(seconds: 2));
      expect(action, 'next');
    });

    test('passes through errors for /token/ endpoint without retry', () async {
      final requestOptions = RequestOptions(path: '/token/');
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 403,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      final action = await handler.result.timeout(const Duration(seconds: 2));
      expect(action, 'next');
    });

    test('passes through errors that already have X-Token-Retry header',
        () async {
      final requestOptions = RequestOptions(
        path: '/receipts',
        headers: {'X-Token-Retry': 'true'},
      );
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 403,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      final action = await handler.result.timeout(const Duration(seconds: 2));
      expect(action, 'next');
    });

    test('attempts refresh on 401 and retries request on success', () async {
      final dio = _setUpSuccessfulRetryScenario(
        mockAuthModel: mockAuthModel,
        mockGroupModel: mockGroupModel,
      );

      dio.interceptors.add(InterceptorsWrapper(
        onRequest: (options, handler) {
          if (options.headers.containsKey('X-Token-Retry')) {
            handler.resolve(Response(
              data: {'success': true},
              statusCode: 200,
              requestOptions: options,
            ));
          } else {
            handler.next(options);
          }
        },
      ));

      final requestOptions = RequestOptions(path: '/receipts');
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 401,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      final action = await handler.result.timeout(const Duration(seconds: 5));
      expect(action, 'resolve',
          reason: 'Should resolve with retried response');
      expect(handler.resolvedResponse?.statusCode, 200);
    });

    test('refreshes and retries on a 403 when the token is expired', () async {
      // The backend returns 403 for an expired session too (it never sends
      // 401). With an expired token in hand, the interceptor still refreshes
      // and retries — the recovery path is keyed on token validity, not status.
      final dio = _setUpSuccessfulRetryScenario(
        mockAuthModel: mockAuthModel,
        mockGroupModel: mockGroupModel,
      );

      dio.interceptors.add(InterceptorsWrapper(
        onRequest: (options, handler) {
          if (options.headers.containsKey('X-Token-Retry')) {
            handler.resolve(Response(
              data: {'success': true},
              statusCode: 200,
              requestOptions: options,
            ));
          } else {
            handler.next(options);
          }
        },
      ));

      final requestOptions = RequestOptions(path: '/receipts');
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 403,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      final action = await handler.result.timeout(const Duration(seconds: 5));
      expect(action, 'resolve',
          reason: 'Should resolve with retried response on an expired-token 403');
    });

    test(
        'does not refresh on a 403 when the token is still valid '
        '(permission denial)', () async {
      // A valid access token means the 403 is a permission denial, not an
      // expired session: the interceptor must surface it untouched — no token
      // refresh (which would burn the one-time refresh token / risk logout) and
      // no retry (a new token can't grant a missing permission).
      when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);

      final requestOptions = RequestOptions(path: '/receipts');
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 403,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      final action = await handler.result.timeout(const Duration(seconds: 2));
      expect(action, 'next',
          reason: 'A permission 403 should pass through without interception');
      verifyNever(() => mockAuthModel.setTokens(any(), any()));
      verifyNever(() => mockAuthModel.getRefreshToken());
    });

    test('falls through to next when token refresh fails', () async {
      when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
      when(() => mockAuthModel.getRefreshToken())
          .thenAnswer((_) async => expiredJwt);
      when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});

      final requestOptions = RequestOptions(path: '/receipts');
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 401,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      final action = await handler.result.timeout(const Duration(seconds: 5));
      expect(action, 'next',
          reason: 'Should pass error through when refresh fails');
    });

    test('sets Authorization header with new JWT on retry', () async {
      final dio = _setUpSuccessfulRetryScenario(
        mockAuthModel: mockAuthModel,
        mockGroupModel: mockGroupModel,
      );

      String? capturedAuthHeader;
      dio.interceptors.add(InterceptorsWrapper(
        onRequest: (options, handler) {
          if (options.headers.containsKey('X-Token-Retry')) {
            capturedAuthHeader = options.headers['Authorization'] as String?;
            handler.resolve(Response(
              data: {'success': true},
              statusCode: 200,
              requestOptions: options,
            ));
          } else {
            handler.next(options);
          }
        },
      ));

      final requestOptions = RequestOptions(path: '/receipts');
      final error = DioException(
        requestOptions: requestOptions,
        response: Response(
          statusCode: 401,
          requestOptions: requestOptions,
        ),
      );

      final handler = TestErrorHandler();
      interceptor.onError(error, handler);

      await handler.result.timeout(const Duration(seconds: 5));

      expect(capturedAuthHeader, isNotNull,
          reason: 'Retry should include Authorization header');
      expect(capturedAuthHeader, startsWith('Bearer '),
          reason: 'Retry should include new JWT in Authorization header');
    });
  });
}
