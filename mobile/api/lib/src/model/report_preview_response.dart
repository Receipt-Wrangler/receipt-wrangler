//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_preview_response.g.dart';

/// ReportPreviewResponse
///
/// Properties:
/// * [html] - The rendered report preview as a self-contained HTML document.
/// * [receiptCount] - The number of receipts the current configuration covers.
@BuiltValue()
abstract class ReportPreviewResponse implements Built<ReportPreviewResponse, ReportPreviewResponseBuilder> {
  /// The rendered report preview as a self-contained HTML document.
  @BuiltValueField(wireName: r'html')
  String get html;

  /// The number of receipts the current configuration covers.
  @BuiltValueField(wireName: r'receiptCount')
  int get receiptCount;

  ReportPreviewResponse._();

  factory ReportPreviewResponse([void updates(ReportPreviewResponseBuilder b)]) = _$ReportPreviewResponse;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportPreviewResponseBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportPreviewResponse> get serializer => _$ReportPreviewResponseSerializer();
}

class _$ReportPreviewResponseSerializer implements PrimitiveSerializer<ReportPreviewResponse> {
  @override
  final Iterable<Type> types = const [ReportPreviewResponse, _$ReportPreviewResponse];

  @override
  final String wireName = r'ReportPreviewResponse';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportPreviewResponse object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'html';
    yield serializers.serialize(
      object.html,
      specifiedType: const FullType(String),
    );
    yield r'receiptCount';
    yield serializers.serialize(
      object.receiptCount,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportPreviewResponse object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportPreviewResponseBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'html':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.html = valueDes;
          break;
        case r'receiptCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptCount = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportPreviewResponse deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportPreviewResponseBuilder();
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

