//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'oidc_link_start_view.g.dart';

/// The authorization URL for a mobile \"connect account\" flow. Safe to hand to the caller - its state belongs to a session already bound to that caller's user id, so holding it grants nothing the bearer token did not already.
///
/// Properties:
/// * [authorizationUrl] 
@BuiltValue()
abstract class OidcLinkStartView implements Built<OidcLinkStartView, OidcLinkStartViewBuilder> {
  @BuiltValueField(wireName: r'authorizationUrl')
  String get authorizationUrl;

  OidcLinkStartView._();

  factory OidcLinkStartView([void updates(OidcLinkStartViewBuilder b)]) = _$OidcLinkStartView;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(OidcLinkStartViewBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<OidcLinkStartView> get serializer => _$OidcLinkStartViewSerializer();
}

class _$OidcLinkStartViewSerializer implements PrimitiveSerializer<OidcLinkStartView> {
  @override
  final Iterable<Type> types = const [OidcLinkStartView, _$OidcLinkStartView];

  @override
  final String wireName = r'OidcLinkStartView';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    OidcLinkStartView object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'authorizationUrl';
    yield serializers.serialize(
      object.authorizationUrl,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    OidcLinkStartView object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required OidcLinkStartViewBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'authorizationUrl':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.authorizationUrl = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  OidcLinkStartView deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = OidcLinkStartViewBuilder();
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

