//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_document.g.dart';

/// ReportDocument
///
/// Properties:
/// * [title] 
/// * [intro] 
/// * [footer] 
@BuiltValue()
abstract class ReportDocument implements Built<ReportDocument, ReportDocumentBuilder> {
  @BuiltValueField(wireName: r'title')
  String? get title;

  @BuiltValueField(wireName: r'intro')
  String? get intro;

  @BuiltValueField(wireName: r'footer')
  String? get footer;

  ReportDocument._();

  factory ReportDocument([void updates(ReportDocumentBuilder b)]) = _$ReportDocument;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportDocumentBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportDocument> get serializer => _$ReportDocumentSerializer();
}

class _$ReportDocumentSerializer implements PrimitiveSerializer<ReportDocument> {
  @override
  final Iterable<Type> types = const [ReportDocument, _$ReportDocument];

  @override
  final String wireName = r'ReportDocument';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportDocument object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.title != null) {
      yield r'title';
      yield serializers.serialize(
        object.title,
        specifiedType: const FullType(String),
      );
    }
    if (object.intro != null) {
      yield r'intro';
      yield serializers.serialize(
        object.intro,
        specifiedType: const FullType(String),
      );
    }
    if (object.footer != null) {
      yield r'footer';
      yield serializers.serialize(
        object.footer,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportDocument object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportDocumentBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'title':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.title = valueDes;
          break;
        case r'intro':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.intro = valueDes;
          break;
        case r'footer':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.footer = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportDocument deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportDocumentBuilder();
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

