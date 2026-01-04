//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'magic_fill_command.g.dart';

/// MagicFillCommand
///
/// Properties:
/// * [imageData] - Image data
/// * [filename] - Name of file
@BuiltValue()
abstract class MagicFillCommand implements Built<MagicFillCommand, MagicFillCommandBuilder> {
  /// Image data
  @BuiltValueField(wireName: r'imageData')
  BuiltList<int> get imageData;

  /// Name of file
  @BuiltValueField(wireName: r'filename')
  String get filename;

  MagicFillCommand._();

  factory MagicFillCommand([void updates(MagicFillCommandBuilder b)]) = _$MagicFillCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(MagicFillCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<MagicFillCommand> get serializer => _$MagicFillCommandSerializer();
}

class _$MagicFillCommandSerializer implements PrimitiveSerializer<MagicFillCommand> {
  @override
  final Iterable<Type> types = const [MagicFillCommand, _$MagicFillCommand];

  @override
  final String wireName = r'MagicFillCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    MagicFillCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'imageData';
    yield serializers.serialize(
      object.imageData,
      specifiedType: const FullType(BuiltList, [FullType(int)]),
    );
    yield r'filename';
    yield serializers.serialize(
      object.filename,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    MagicFillCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required MagicFillCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'imageData':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.imageData.replace(valueDes);
          break;
        case r'filename':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.filename = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  MagicFillCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = MagicFillCommandBuilder();
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

