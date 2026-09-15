//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_paged_request_filter.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_summary_command.g.dart';

/// ReceiptSummaryCommand
///
/// Properties:
/// * [filter] 
/// * [configurationGroupId] - The group whose receipt settings shape the breakdown. Exists for the synthetic \"All\" group, which spans every group the caller belongs to and so has no meaningful settings of its own - the client picks which member group's configuration to apply. OMIT it for a real group, where the group in the path is used. The DATA is always the path group's filtered set; this only chooses the shape of the breakdown. A group the caller cannot read is a 403.
@BuiltValue()
abstract class ReceiptSummaryCommand implements Built<ReceiptSummaryCommand, ReceiptSummaryCommandBuilder> {
  @BuiltValueField(wireName: r'filter')
  ReceiptPagedRequestFilter? get filter;

  /// The group whose receipt settings shape the breakdown. Exists for the synthetic \"All\" group, which spans every group the caller belongs to and so has no meaningful settings of its own - the client picks which member group's configuration to apply. OMIT it for a real group, where the group in the path is used. The DATA is always the path group's filtered set; this only chooses the shape of the breakdown. A group the caller cannot read is a 403.
  @BuiltValueField(wireName: r'configurationGroupId')
  int? get configurationGroupId;

  ReceiptSummaryCommand._();

  factory ReceiptSummaryCommand([void updates(ReceiptSummaryCommandBuilder b)]) = _$ReceiptSummaryCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptSummaryCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptSummaryCommand> get serializer => _$ReceiptSummaryCommandSerializer();
}

class _$ReceiptSummaryCommandSerializer implements PrimitiveSerializer<ReceiptSummaryCommand> {
  @override
  final Iterable<Type> types = const [ReceiptSummaryCommand, _$ReceiptSummaryCommand];

  @override
  final String wireName = r'ReceiptSummaryCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptSummaryCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.filter != null) {
      yield r'filter';
      yield serializers.serialize(
        object.filter,
        specifiedType: const FullType(ReceiptPagedRequestFilter),
      );
    }
    if (object.configurationGroupId != null) {
      yield r'configurationGroupId';
      yield serializers.serialize(
        object.configurationGroupId,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptSummaryCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptSummaryCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'filter':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptPagedRequestFilter),
          ) as ReceiptPagedRequestFilter;
          result.filter.replace(valueDes);
          break;
        case r'configurationGroupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.configurationGroupId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptSummaryCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptSummaryCommandBuilder();
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

