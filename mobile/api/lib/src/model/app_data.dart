//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/claims.dart';
import 'package:openapi/src/model/user_preferences.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/about.dart';
import 'package:openapi/src/model/category.dart';
import 'package:openapi/src/model/tag.dart';
import 'package:openapi/src/model/currency_symbol_position.dart';
import 'package:openapi/src/model/icon.dart';
import 'package:openapi/src/model/user_view.dart';
import 'package:openapi/src/model/currency_separator.dart';
import 'package:openapi/src/model/group.dart';
import 'package:openapi/src/model/feature_config.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'app_data.g.dart';

/// AppData
///
/// Properties:
/// * [about] 
/// * [claims] 
/// * [groups] - Groups in the system
/// * [users] - Users in the system
/// * [userPreferences] 
/// * [featureConfig] 
/// * [categories] - Categories in the system
/// * [tags] - Tags in the system
/// * [jwt] - JWT token
/// * [refreshToken] - Refresh token
/// * [currencyDisplay] - Currency display
/// * [currencyThousandthsSeparator] 
/// * [currencyDecimalSeparator] 
/// * [currencySymbolPosition] 
/// * [currencyHideDecimalPlaces] - Whether to hide decimal places
/// * [icons] - Icons in the system
@BuiltValue()
abstract class AppData implements Built<AppData, AppDataBuilder> {
  @BuiltValueField(wireName: r'about')
  About get about;

  @BuiltValueField(wireName: r'claims')
  Claims get claims;

  /// Groups in the system
  @BuiltValueField(wireName: r'groups')
  BuiltList<Group> get groups;

  /// Users in the system
  @BuiltValueField(wireName: r'users')
  BuiltList<UserView> get users;

  @BuiltValueField(wireName: r'userPreferences')
  UserPreferences get userPreferences;

  @BuiltValueField(wireName: r'featureConfig')
  FeatureConfig get featureConfig;

  /// Categories in the system
  @BuiltValueField(wireName: r'categories')
  BuiltList<Category> get categories;

  /// Tags in the system
  @BuiltValueField(wireName: r'tags')
  BuiltList<Tag> get tags;

  /// JWT token
  @BuiltValueField(wireName: r'jwt')
  String? get jwt;

  /// Refresh token
  @BuiltValueField(wireName: r'refreshToken')
  String? get refreshToken;

  /// Currency display
  @BuiltValueField(wireName: r'currencyDisplay')
  String get currencyDisplay;

  @BuiltValueField(wireName: r'currencyThousandthsSeparator')
  CurrencySeparator? get currencyThousandthsSeparator;
  // enum currencyThousandthsSeparatorEnum {  ,,  .,  };

  @BuiltValueField(wireName: r'currencyDecimalSeparator')
  CurrencySeparator? get currencyDecimalSeparator;
  // enum currencyDecimalSeparatorEnum {  ,,  .,  };

  @BuiltValueField(wireName: r'currencySymbolPosition')
  CurrencySymbolPosition? get currencySymbolPosition;
  // enum currencySymbolPositionEnum {  START,  END,  };

  /// Whether to hide decimal places
  @BuiltValueField(wireName: r'currencyHideDecimalPlaces')
  bool? get currencyHideDecimalPlaces;

  /// Icons in the system
  @BuiltValueField(wireName: r'icons')
  BuiltList<Icon> get icons;

  AppData._();

  factory AppData([void updates(AppDataBuilder b)]) = _$AppData;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(AppDataBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<AppData> get serializer => _$AppDataSerializer();
}

class _$AppDataSerializer implements PrimitiveSerializer<AppData> {
  @override
  final Iterable<Type> types = const [AppData, _$AppData];

  @override
  final String wireName = r'AppData';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    AppData object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'about';
    yield serializers.serialize(
      object.about,
      specifiedType: const FullType(About),
    );
    yield r'claims';
    yield serializers.serialize(
      object.claims,
      specifiedType: const FullType(Claims),
    );
    yield r'groups';
    yield serializers.serialize(
      object.groups,
      specifiedType: const FullType(BuiltList, [FullType(Group)]),
    );
    yield r'users';
    yield serializers.serialize(
      object.users,
      specifiedType: const FullType(BuiltList, [FullType(UserView)]),
    );
    yield r'userPreferences';
    yield serializers.serialize(
      object.userPreferences,
      specifiedType: const FullType(UserPreferences),
    );
    yield r'featureConfig';
    yield serializers.serialize(
      object.featureConfig,
      specifiedType: const FullType(FeatureConfig),
    );
    yield r'categories';
    yield serializers.serialize(
      object.categories,
      specifiedType: const FullType(BuiltList, [FullType(Category)]),
    );
    yield r'tags';
    yield serializers.serialize(
      object.tags,
      specifiedType: const FullType(BuiltList, [FullType(Tag)]),
    );
    if (object.jwt != null) {
      yield r'jwt';
      yield serializers.serialize(
        object.jwt,
        specifiedType: const FullType(String),
      );
    }
    if (object.refreshToken != null) {
      yield r'refreshToken';
      yield serializers.serialize(
        object.refreshToken,
        specifiedType: const FullType(String),
      );
    }
    yield r'currencyDisplay';
    yield serializers.serialize(
      object.currencyDisplay,
      specifiedType: const FullType(String),
    );
    if (object.currencyThousandthsSeparator != null) {
      yield r'currencyThousandthsSeparator';
      yield serializers.serialize(
        object.currencyThousandthsSeparator,
        specifiedType: const FullType(CurrencySeparator),
      );
    }
    if (object.currencyDecimalSeparator != null) {
      yield r'currencyDecimalSeparator';
      yield serializers.serialize(
        object.currencyDecimalSeparator,
        specifiedType: const FullType(CurrencySeparator),
      );
    }
    if (object.currencySymbolPosition != null) {
      yield r'currencySymbolPosition';
      yield serializers.serialize(
        object.currencySymbolPosition,
        specifiedType: const FullType(CurrencySymbolPosition),
      );
    }
    if (object.currencyHideDecimalPlaces != null) {
      yield r'currencyHideDecimalPlaces';
      yield serializers.serialize(
        object.currencyHideDecimalPlaces,
        specifiedType: const FullType(bool),
      );
    }
    yield r'icons';
    yield serializers.serialize(
      object.icons,
      specifiedType: const FullType(BuiltList, [FullType(Icon)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    AppData object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required AppDataBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'about':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(About),
          ) as About;
          result.about.replace(valueDes);
          break;
        case r'claims':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(Claims),
          ) as Claims;
          result.claims.replace(valueDes);
          break;
        case r'groups':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Group)]),
          ) as BuiltList<Group>;
          result.groups.replace(valueDes);
          break;
        case r'users':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(UserView)]),
          ) as BuiltList<UserView>;
          result.users.replace(valueDes);
          break;
        case r'userPreferences':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(UserPreferences),
          ) as UserPreferences;
          result.userPreferences.replace(valueDes);
          break;
        case r'featureConfig':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(FeatureConfig),
          ) as FeatureConfig;
          result.featureConfig.replace(valueDes);
          break;
        case r'categories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Category)]),
          ) as BuiltList<Category>;
          result.categories.replace(valueDes);
          break;
        case r'tags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Tag)]),
          ) as BuiltList<Tag>;
          result.tags.replace(valueDes);
          break;
        case r'jwt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.jwt = valueDes;
          break;
        case r'refreshToken':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.refreshToken = valueDes;
          break;
        case r'currencyDisplay':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.currencyDisplay = valueDes;
          break;
        case r'currencyThousandthsSeparator':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(CurrencySeparator),
          ) as CurrencySeparator;
          result.currencyThousandthsSeparator = valueDes;
          break;
        case r'currencyDecimalSeparator':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(CurrencySeparator),
          ) as CurrencySeparator;
          result.currencyDecimalSeparator = valueDes;
          break;
        case r'currencySymbolPosition':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(CurrencySymbolPosition),
          ) as CurrencySymbolPosition;
          result.currencySymbolPosition = valueDes;
          break;
        case r'currencyHideDecimalPlaces':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.currencyHideDecimalPlaces = valueDes;
          break;
        case r'icons':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Icon)]),
          ) as BuiltList<Icon>;
          result.icons.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  AppData deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = AppDataBuilder();
    final serializedList = (serialized as Iterable<Object?>).toList();
    final unhandled = <Object?>[];
    _deserializeProperties(
      serializers,
      serialized,
      specifiedType: specifiedType,
      serializedList: serializedList,
      unhandled: unhandled,
      result: result,
    );
    return result.build();
  }
}

