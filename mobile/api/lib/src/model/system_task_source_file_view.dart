//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'system_task_source_file_view.g.dart';

/// SystemTaskSourceFileView
///
/// Properties:
/// * [name] - The upload's original file name
/// * [encodedImage] - Base64 encoded image, converted for display
@BuiltValue()
abstract class SystemTaskSourceFileView implements Built<SystemTaskSourceFileView, SystemTaskSourceFileViewBuilder> {
  /// The upload's original file name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Base64 encoded image, converted for display
  @BuiltValueField(wireName: r'encodedImage')
  String get encodedImage;

  SystemTaskSourceFileView._();

  factory SystemTaskSourceFileView([void updates(SystemTaskSourceFileViewBuilder b)]) = _$SystemTaskSourceFileView;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(SystemTaskSourceFileViewBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<SystemTaskSourceFileView> get serializer => _$SystemTaskSourceFileViewSerializer();
}

class _$SystemTaskSourceFileViewSerializer implements PrimitiveSerializer<SystemTaskSourceFileView> {
  @override
  final Iterable<Type> types = const [SystemTaskSourceFileView, _$SystemTaskSourceFileView];

  @override
  final String wireName = r'SystemTaskSourceFileView';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    SystemTaskSourceFileView object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'encodedImage';
    yield serializers.serialize(
      object.encodedImage,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    SystemTaskSourceFileView object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required SystemTaskSourceFileViewBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
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
  SystemTaskSourceFileView deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = SystemTaskSourceFileViewBuilder();
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

