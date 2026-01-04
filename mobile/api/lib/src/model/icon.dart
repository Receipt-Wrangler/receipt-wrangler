//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'icon.g.dart';

/// Icon
///
/// Properties:
/// * [value] - Icon value
/// * [displayValue] - Icon display value
@BuiltValue()
abstract class Icon implements Built<Icon, IconBuilder> {
  /// Icon value
  @BuiltValueField(wireName: r'value')
  String get value;

  /// Icon display value
  @BuiltValueField(wireName: r'displayValue')
  String get displayValue;

  Icon._();

  factory Icon([void updates(IconBuilder b)]) = _$Icon;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(IconBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Icon> get serializer => _$IconSerializer();
}

class _$IconSerializer implements PrimitiveSerializer<Icon> {
  @override
  final Iterable<Type> types = const [Icon, _$Icon];

  @override
  final String wireName = r'Icon';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Icon object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'value';
    yield serializers.serialize(
      object.value,
      specifiedType: const FullType(String),
    );
    yield r'displayValue';
    yield serializers.serialize(
      object.displayValue,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    Icon object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required IconBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'value':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.value = valueDes;
          break;
        case r'displayValue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.displayValue = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  Icon deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = IconBuilder();
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

