//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_claim_ledger_entry_command.g.dart';

/// A ledger entry's fields, for create and for a full replace on update. An optional field that is omitted, null or \"\" is stored as not set.
///
/// Properties:
/// * [type] - One of PAYMENT, RECOUPMENT, NOTE; anything else is a 400 keyed type
/// * [date] - Calendar day (yyyy-MM-dd)
/// * [amount] - Decimal string, greater than zero, at most two decimal places. Required for PAYMENT and RECOUPMENT; optional for NOTE.
/// * [note] - At most 2000 characters. Required (non-blank) for NOTE.
@BuiltValue()
abstract class UpsertClaimLedgerEntryCommand implements Built<UpsertClaimLedgerEntryCommand, UpsertClaimLedgerEntryCommandBuilder> {
  /// One of PAYMENT, RECOUPMENT, NOTE; anything else is a 400 keyed type
  @BuiltValueField(wireName: r'type')
  String get type;

  /// Calendar day (yyyy-MM-dd)
  @BuiltValueField(wireName: r'date')
  String get date;

  /// Decimal string, greater than zero, at most two decimal places. Required for PAYMENT and RECOUPMENT; optional for NOTE.
  @BuiltValueField(wireName: r'amount')
  String? get amount;

  /// At most 2000 characters. Required (non-blank) for NOTE.
  @BuiltValueField(wireName: r'note')
  String? get note;

  UpsertClaimLedgerEntryCommand._();

  factory UpsertClaimLedgerEntryCommand([void updates(UpsertClaimLedgerEntryCommandBuilder b)]) = _$UpsertClaimLedgerEntryCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertClaimLedgerEntryCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertClaimLedgerEntryCommand> get serializer => _$UpsertClaimLedgerEntryCommandSerializer();
}

class _$UpsertClaimLedgerEntryCommandSerializer implements PrimitiveSerializer<UpsertClaimLedgerEntryCommand> {
  @override
  final Iterable<Type> types = const [UpsertClaimLedgerEntryCommand, _$UpsertClaimLedgerEntryCommand];

  @override
  final String wireName = r'UpsertClaimLedgerEntryCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertClaimLedgerEntryCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'type';
    yield serializers.serialize(
      object.type,
      specifiedType: const FullType(String),
    );
    yield r'date';
    yield serializers.serialize(
      object.date,
      specifiedType: const FullType(String),
    );
    if (object.amount != null) {
      yield r'amount';
      yield serializers.serialize(
        object.amount,
        specifiedType: const FullType(String),
      );
    }
    if (object.note != null) {
      yield r'note';
      yield serializers.serialize(
        object.note,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertClaimLedgerEntryCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertClaimLedgerEntryCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.type = valueDes;
          break;
        case r'date':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.date = valueDes;
          break;
        case r'amount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.amount = valueDes;
          break;
        case r'note':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.note = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertClaimLedgerEntryCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertClaimLedgerEntryCommandBuilder();
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

