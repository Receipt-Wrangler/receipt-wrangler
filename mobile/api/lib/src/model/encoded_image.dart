//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'encoded_image.g.dart';

/// EncodedImage
///
/// Properties:
/// * [encodedImage] - base64 encoded jpg
@BuiltValue()
abstract class EncodedImage implements Built<EncodedImage, EncodedImageBuilder> {
  /// base64 encoded jpg
  @BuiltValueField(wireName: r'encodedImage')
  String get encodedImage;

  EncodedImage._();

  factory EncodedImage([void updates(EncodedImageBuilder b)]) = _$EncodedImage;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(EncodedImageBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<EncodedImage> get serializer => _$EncodedImageSerializer();
}

class _$EncodedImageSerializer implements PrimitiveSerializer<EncodedImage> {
  @override
  final Iterable<Type> types = const [EncodedImage, _$EncodedImage];

  @override
  final String wireName = r'EncodedImage';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    EncodedImage object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'encodedImage';
    yield serializers.serialize(
      object.encodedImage,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    EncodedImage object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required EncodedImageBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'encodedImage':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.encodedImage = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  EncodedImage deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = EncodedImageBuilder();
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

