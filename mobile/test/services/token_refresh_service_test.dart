import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/services/token_refresh_service.dart';

import '../helpers/auth_test_helpers.dart';

// --- Test-specific mocks (not shared) ---

class MockAppData extends Mock implements AppData {}

class MockClaims extends Mock implements Claims {}

class MockFeatureConfig extends Mock implements FeatureConfig {}

class MockUserPreferences extends Mock implements UserPreferences {}

void main() {
  late MockAuthModel mockAuthModel;
  late MockGroupModel mockGroupModel;
  late MockUserModel mockUserModel;
  late MockUserPreferencesModel mockUserPreferencesModel;
  late MockCategoryModel mockCategoryModel;
  late MockTagModel mockTagModel;
  late MockSystemSettingsModel mockSystemSettingsModel;
  // PermissionsModel has no I/O, so use a real instance (per mobile testing guide).
  late PermissionsModel permissionsModel;
  late MockOpenapi mockClient;
  late MockAuthApi mockAuthApi;
  late MockUserApi mockUserApi;
  late TokenRefreshService service;

  setUpAll(() {
    registerFallbackValue(FakeLogoutCommand());
    registerFallbackValue(MockClaims());
    registerFallbackValue(MockFeatureConfig());
    registerFallbackValue(MockUserPreferences());
    registerFallbackValue(<Group>[]);
    registerFallbackValue(<UserView>[]);
    registerFallbackValue(<Category>[]);
    registerFallbackValue(<Tag>[]);
    registerFallbackValue('');
    registerFallbackValue(CurrencySeparator.period);
    registerFallbackValue(CurrencySymbolPosition.END);
    registerFallbackValue(false);
  });

  setUp(() {
    mockAuthModel = MockAuthModel();
    mockGroupModel = MockGroupModel();
    mockUserModel = MockUserModel();
    mockUserPreferencesModel = MockUserPreferencesModel();
    mockCategoryModel = MockCategoryModel();
    mockTagModel = MockTagModel();
    mockSystemSettingsModel = MockSystemSettingsModel();
    permissionsModel = PermissionsModel();
    mockClient = MockOpenapi();
    mockAuthApi = MockAuthApi();
    mockUserApi = MockUserApi();

    when(() => mockClient.getAuthApi()).thenReturn(mockAuthApi);
    when(() => mockClient.getUserApi()).thenReturn(mockUserApi);

    OpenApiClient.client = mockClient;

    // Reset and re-initialize the singleton for each test.
    service = TokenRefreshService();
    service.resetForTesting();
    service.initialize(
      authModel: mockAuthModel,
      groupModel: mockGroupModel,
      userModel: mockUserModel,
      userPreferencesModel: mockUserPreferencesModel,
      categoryModel: mockCategoryModel,
      tagModel: mockTagModel,
      systemSettingsModel: mockSystemSettingsModel,
      permissionsModel: permissionsModel,
    );
  });

  group('TokenRefreshService', () {
    group('refreshTokens with force=false', () {
      test('returns true when JWT is still valid', () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

        final result = await service.refreshTokens();

        expect(result, true);
        verifyNever(() => mockAuthApi.getNewRefreshToken(
            logoutCommand: any(named: 'logoutCommand')));
      });

      test('refreshes token when JWT is expired but refresh token is valid',
          () async {
        final newJwt = validJwt;
        final newRefresh = validJwt;

        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenAnswer(
                (_) async => createTokenRefreshResponse(newJwt, newRefresh));
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});

        final result = await service.refreshTokens();

        expect(result, true);
        verify(() => mockAuthModel.setTokens(newJwt, newRefresh)).called(1);
      });

      test('purges tokens when both JWT and refresh token are expired',
          () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});

        final result = await service.refreshTokens();

        expect(result, false);
        verify(() => mockAuthModel.purgeTokens()).called(1);
      });

      test('treats corrupted JWT as invalid and attempts refresh', () async {
        final newJwt = validJwt;
        final newRefresh = validJwt;

        when(() => mockAuthModel.getJwt())
            .thenAnswer((_) async => 'not-a-valid-jwt');
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenAnswer(
                (_) async => createTokenRefreshResponse(newJwt, newRefresh));
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});

        final result = await service.refreshTokens();

        expect(result, true);
        verify(() => mockAuthApi.getNewRefreshToken(
            logoutCommand: any(named: 'logoutCommand'))).called(1);
      });

      test('purges tokens when both JWT and refresh token are corrupted',
          () async {
        when(() => mockAuthModel.getJwt())
            .thenAnswer((_) async => 'corrupted-jwt');
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => 'corrupted-refresh');
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});

        final result = await service.refreshTokens();

        expect(result, false);
        verify(() => mockAuthModel.purgeTokens()).called(1);
      });

      test('purges tokens when JWT is null', () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => null);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => null);
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});

        final result = await service.refreshTokens();

        expect(result, false);
        verify(() => mockAuthModel.purgeTokens()).called(1);
      });

      test('purges tokens when refresh endpoint throws', () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});
        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenThrow(DioException(
          requestOptions: RequestOptions(path: '/token/'),
          response: Response(
            statusCode: 500,
            requestOptions: RequestOptions(path: '/token/'),
          ),
        ));

        final result = await service.refreshTokens();

        expect(result, false);
        verify(() => mockAuthModel.purgeTokens()).called(1);
      });
    });

    group('refreshTokens with force=true', () {
      test('refreshes even when JWT is still valid', () async {
        final newJwt = validJwt;
        final newRefresh = validJwt;

        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenAnswer(
                (_) async => createTokenRefreshResponse(newJwt, newRefresh));
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});

        final result = await service.refreshTokens(force: true);

        expect(result, true);
        verify(() => mockAuthApi.getNewRefreshToken(
            logoutCommand: any(named: 'logoutCommand'))).called(1);
      });

      test('returns false when force=true but refresh token is expired',
          () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});

        final result = await service.refreshTokens(force: true);

        expect(result, false);
        verify(() => mockAuthModel.purgeTokens()).called(1);
        verifyNever(() => mockAuthApi.getNewRefreshToken(
            logoutCommand: any(named: 'logoutCommand')));
      });
    });

    group('serialization - concurrent calls share one Future', () {
      test('concurrent calls return the same result without duplicate requests',
          () async {
        var callCount = 0;
        final newJwt = validJwt;
        final newRefresh = validJwt;

        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenAnswer((_) async {
          callCount++;
          await Future.delayed(const Duration(milliseconds: 50));
          return createTokenRefreshResponse(newJwt, newRefresh);
        });
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});

        final results = await Future.wait([
          service.refreshTokens(),
          service.refreshTokens(),
          service.refreshTokens(),
          service.refreshTokens(),
          service.refreshTokens(),
        ]);

        expect(results, everyElement(true));
        expect(callCount, 1);
      });

      test('after completion, a new call makes a new request', () async {
        var callCount = 0;
        final newJwt = validJwt;
        final newRefresh = validJwt;

        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenAnswer((_) async {
          callCount++;
          return createTokenRefreshResponse(newJwt, newRefresh);
        });
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});

        await service.refreshTokens();
        expect(callCount, 1);

        await service.refreshTokens();
        expect(callCount, 2);
      });

      test('concurrent calls all get false when refresh fails', () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});
        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenAnswer((_) async {
          await Future.delayed(const Duration(milliseconds: 50));
          throw DioException(
            requestOptions: RequestOptions(path: '/token/'),
            response: Response(
              statusCode: 500,
              requestOptions: RequestOptions(path: '/token/'),
            ),
          );
        });

        final results = await Future.wait([
          service.refreshTokens(),
          service.refreshTokens(),
          service.refreshTokens(),
        ]);

        expect(results, everyElement(false));
      });
    });

    group('app data loading', () {
      test('loads app data when groups are empty and user is authenticated',
          () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([]);

        final mockAppData = MockAppData();
        when(() => mockAppData.jwt).thenReturn('');
        when(() => mockAppData.refreshToken).thenReturn('');
        when(() => mockAppData.claims).thenReturn(MockClaims());
        when(() => mockAppData.featureConfig).thenReturn(MockFeatureConfig());
        when(() => mockAppData.groups).thenReturn(BuiltList<Group>());
        when(() => mockAppData.users).thenReturn(BuiltList<UserView>());
        when(() => mockAppData.userPreferences)
            .thenReturn(MockUserPreferences());
        when(() => mockAppData.categories).thenReturn(BuiltList<Category>());
        when(() => mockAppData.tags).thenReturn(BuiltList<Tag>());
        when(() => mockAppData.currencyDisplay).thenReturn('');
        when(() => mockAppData.currencyDecimalSeparator)
            .thenReturn(CurrencySeparator.period);
        when(() => mockAppData.currencyThousandthsSeparator)
            .thenReturn(CurrencySeparator.comma);
        when(() => mockAppData.currencySymbolPosition)
            .thenReturn(CurrencySymbolPosition.END);
        when(() => mockAppData.currencyHideDecimalPlaces).thenReturn(false);
        when(() => mockAppData.appPermissions)
            .thenReturn(BuiltList<String>());
        when(() => mockAppData.groupPermissions)
            .thenReturn(BuiltMap<String, BuiltList<String>>());

        when(() => mockUserApi.getAppData()).thenAnswer((_) async => Response(
              data: mockAppData,
              requestOptions: RequestOptions(path: '/user/appData'),
              statusCode: 200,
            ));
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});
        when(() => mockAuthModel.setClaims(any())).thenReturn(null);
        when(() => mockAuthModel.setFeatureConfig(any())).thenReturn(null);
        when(() => mockGroupModel.setGroups(any())).thenReturn(null);
        when(() => mockUserModel.setUsers(any())).thenReturn(null);
        when(() => mockUserPreferencesModel.setUserPreferences(any()))
            .thenReturn(null);
        when(() => mockCategoryModel.setCategories(any())).thenReturn(null);
        when(() => mockTagModel.setTags(any())).thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyDisplay(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyDecimalSeparator(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyThousandSeparator(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencySymbolPosition(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyHideDecimalPlaces(any()))
            .thenReturn(null);

        final result = await service.refreshTokens();

        expect(result, true);
        verify(() => mockUserApi.getAppData()).called(1);
      });

      test(
          'stores all app data fields when jwt and refreshToken are null in response',
          () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([]);

        final mockAppData = MockAppData();
        when(() => mockAppData.jwt).thenReturn(null);
        when(() => mockAppData.refreshToken).thenReturn(null);
        when(() => mockAppData.claims).thenReturn(MockClaims());
        when(() => mockAppData.featureConfig).thenReturn(MockFeatureConfig());
        when(() => mockAppData.groups).thenReturn(BuiltList<Group>());
        when(() => mockAppData.users).thenReturn(BuiltList<UserView>());
        when(() => mockAppData.userPreferences)
            .thenReturn(MockUserPreferences());
        when(() => mockAppData.categories).thenReturn(BuiltList<Category>());
        when(() => mockAppData.tags).thenReturn(BuiltList<Tag>());
        when(() => mockAppData.currencyDisplay).thenReturn('');
        when(() => mockAppData.currencyDecimalSeparator)
            .thenReturn(CurrencySeparator.period);
        when(() => mockAppData.currencyThousandthsSeparator)
            .thenReturn(CurrencySeparator.comma);
        when(() => mockAppData.currencySymbolPosition)
            .thenReturn(CurrencySymbolPosition.END);
        when(() => mockAppData.currencyHideDecimalPlaces).thenReturn(false);
        when(() => mockAppData.appPermissions)
            .thenReturn(BuiltList<String>());
        when(() => mockAppData.groupPermissions)
            .thenReturn(BuiltMap<String, BuiltList<String>>());

        when(() => mockUserApi.getAppData()).thenAnswer((_) async => Response(
              data: mockAppData,
              requestOptions: RequestOptions(path: '/user/appData'),
              statusCode: 200,
            ));
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});
        when(() => mockAuthModel.setClaims(any())).thenReturn(null);
        when(() => mockAuthModel.setFeatureConfig(any())).thenReturn(null);
        when(() => mockGroupModel.setGroups(any())).thenReturn(null);
        when(() => mockUserModel.setUsers(any())).thenReturn(null);
        when(() => mockUserPreferencesModel.setUserPreferences(any()))
            .thenReturn(null);
        when(() => mockCategoryModel.setCategories(any())).thenReturn(null);
        when(() => mockTagModel.setTags(any())).thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyDisplay(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyDecimalSeparator(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyThousandSeparator(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencySymbolPosition(any()))
            .thenReturn(null);
        when(() => mockSystemSettingsModel.setCurrencyHideDecimalPlaces(any()))
            .thenReturn(null);

        final result = await service.refreshTokens();

        expect(result, true);
        // setTokens should NOT be called when both tokens are null
        verifyNever(() => mockAuthModel.setTokens(any(), any()));
        // But all other app data fields should still be stored
        verify(() => mockAuthModel.setClaims(any())).called(1);
        verify(() => mockGroupModel.setGroups(any())).called(1);
        verify(() => mockUserModel.setUsers(any())).called(1);
      });

      test(
          'returns true when app data loading fails after successful token refresh',
          () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => expiredJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([]);
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});

        final newJwt = validJwt;
        final newRefresh = validJwt;

        when(() => mockAuthApi.getNewRefreshToken(
                logoutCommand: any(named: 'logoutCommand')))
            .thenAnswer(
                (_) async => createTokenRefreshResponse(newJwt, newRefresh));
        when(() => mockAuthModel.setTokens(any(), any()))
            .thenAnswer((_) async {});

        // App data loading throws
        when(() => mockUserApi.getAppData()).thenThrow(DioException(
          requestOptions: RequestOptions(path: '/user/appData'),
          response: Response(
            statusCode: 500,
            requestOptions: RequestOptions(path: '/user/appData'),
          ),
        ));

        final result = await service.refreshTokens();

        // Should still return true because tokens were refreshed successfully
        expect(result, true);
        // Should NOT purge the freshly-obtained tokens
        verifyNever(() => mockAuthModel.purgeTokens());
      });

      test('returns true when app data loading fails with valid JWT',
          () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([]);
        when(() => mockAuthModel.purgeTokens()).thenAnswer((_) async {});

        when(() => mockUserApi.getAppData()).thenThrow(DioException(
          requestOptions: RequestOptions(path: '/user/appData'),
          response: Response(
            statusCode: 500,
            requestOptions: RequestOptions(path: '/user/appData'),
          ),
        ));

        final result = await service.refreshTokens();

        expect(result, true);
        verifyNever(() => mockAuthModel.purgeTokens());
      });

      test('skips app data loading when groups already exist', () async {
        when(() => mockAuthModel.getJwt()).thenAnswer((_) async => validJwt);
        when(() => mockAuthModel.getRefreshToken())
            .thenAnswer((_) async => validJwt);
        when(() => mockGroupModel.groups).thenReturn([MockGroup()]);

        final result = await service.refreshTokens();

        expect(result, true);
        verifyNever(() => mockUserApi.getAppData());
      });
    });
  });
}
