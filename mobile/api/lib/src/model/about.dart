//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'about.g.dart';

/// About
///
/// Properties:
/// * [buildDate] - Build date
/// * [version] - Version
@BuiltValue()
abstract class About implements Built<About, AboutBuilder> {
  /// Build date
  @BuiltValueField(wireName: r'buildDate')
  String get buildDate;

  /// Version
  @BuiltValueField(wireName: r'version')
  String get version;

  About._();

  factory About([void updates(AboutBuilder b)]) = _$About;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(AboutBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<About> get serializer => _$AboutSerializer();
}

class _$AboutSerializer implements PrimitiveSerializer<About> {
  @override
  final Iterable<Type> types = const [About, _$About];

  @override
  final String wireName = r'About';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    About object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'buildDate';
    yield serializers.serialize(
      object.buildDate,
      specifiedType: const FullType(String),
    );
    yield r'version';
    yield serializers.serialize(
      object.version,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    About object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required AboutBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'buildDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.buildDate = valueDes;
          break;
        case r'version':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.version = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  About deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = AboutBuilder();
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

