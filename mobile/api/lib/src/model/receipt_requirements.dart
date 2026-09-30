//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_requirements.g.dart';

/// ReceiptRequirements
///
/// Properties:
/// * [commentRequired] - At least one comment is required on the group's receipts
/// * [imageRequired] - At least one image is required on the group's receipts
@BuiltValue()
abstract class ReceiptRequirements implements Built<ReceiptRequirements, ReceiptRequirementsBuilder> {
  /// At least one comment is required on the group's receipts
  @BuiltValueField(wireName: r'commentRequired')
  bool get commentRequired;

  /// At least one image is required on the group's receipts
  @BuiltValueField(wireName: r'imageRequired')
  bool get imageRequired;

  ReceiptRequirements._();

  factory ReceiptRequirements([void updates(ReceiptRequirementsBuilder b)]) = _$ReceiptRequirements;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptRequirementsBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptRequirements> get serializer => _$ReceiptRequirementsSerializer();
}

class _$ReceiptRequirementsSerializer implements PrimitiveSerializer<ReceiptRequirements> {
  @override
  final Iterable<Type> types = const [ReceiptRequirements, _$ReceiptRequirements];

  @override
  final String wireName = r'ReceiptRequirements';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptRequirements object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'commentRequired';
    yield serializers.serialize(
      object.commentRequired,
      specifiedType: const FullType(bool),
    );
    yield r'imageRequired';
    yield serializers.serialize(
      object.imageRequired,
      specifiedType: const FullType(bool),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptRequirements object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptRequirementsBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'commentRequired':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.commentRequired = valueDes;
          break;
        case r'imageRequired':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.imageRequired = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptRequirements deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptRequirementsBuilder();
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

