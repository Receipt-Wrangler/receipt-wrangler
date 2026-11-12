//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim_ledger_entry.g.dart';

/// One dated event recorded on a claim after submission
///
/// Properties:
/// * [id] 
/// * [claimId] 
/// * [type] - One of PAYMENT (money received; adds to paidAmount), RECOUPMENT (money taken back; subtracts from paidAmount) or NOTE (no effect on the figures). A plain string, not an enum; an entry of a type the server does not recognize counts toward no figure.
/// * [date] - Calendar day the event happened (yyyy-MM-dd)
/// * [amount] - Decimal string, always positive; the type gives the direction. Present on every PAYMENT and RECOUPMENT; optional, and never counted, on a NOTE.
/// * [note] - Free text. Absent when not set.
/// * [createdAt] - RFC 3339 instant the entry was recorded
/// * [updatedAt] - RFC 3339 instant the entry was last changed
/// * [createdBy] - The user who recorded the entry
@BuiltValue()
abstract class ClaimLedgerEntry implements Built<ClaimLedgerEntry, ClaimLedgerEntryBuilder> {
  @BuiltValueField(wireName: r'id')
  int get id;

  @BuiltValueField(wireName: r'claimId')
  int get claimId;

  /// One of PAYMENT (money received; adds to paidAmount), RECOUPMENT (money taken back; subtracts from paidAmount) or NOTE (no effect on the figures). A plain string, not an enum; an entry of a type the server does not recognize counts toward no figure.
  @BuiltValueField(wireName: r'type')
  String get type;

  /// Calendar day the event happened (yyyy-MM-dd)
  @BuiltValueField(wireName: r'date')
  String get date;

  /// Decimal string, always positive; the type gives the direction. Present on every PAYMENT and RECOUPMENT; optional, and never counted, on a NOTE.
  @BuiltValueField(wireName: r'amount')
  String? get amount;

  /// Free text. Absent when not set.
  @BuiltValueField(wireName: r'note')
  String? get note;

  /// RFC 3339 instant the entry was recorded
  @BuiltValueField(wireName: r'createdAt')
  String get createdAt;

  /// RFC 3339 instant the entry was last changed
  @BuiltValueField(wireName: r'updatedAt')
  String get updatedAt;

  /// The user who recorded the entry
  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  ClaimLedgerEntry._();

  factory ClaimLedgerEntry([void updates(ClaimLedgerEntryBuilder b)]) = _$ClaimLedgerEntry;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimLedgerEntryBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ClaimLedgerEntry> get serializer => _$ClaimLedgerEntrySerializer();
}

class _$ClaimLedgerEntrySerializer implements PrimitiveSerializer<ClaimLedgerEntry> {
  @override
  final Iterable<Type> types = const [ClaimLedgerEntry, _$ClaimLedgerEntry];

  @override
  final String wireName = r'ClaimLedgerEntry';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ClaimLedgerEntry object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'claimId';
    yield serializers.serialize(
      object.claimId,
      specifiedType: const FullType(int),
    );
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
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(String),
    );
    yield r'updatedAt';
    yield serializers.serialize(
      object.updatedAt,
      specifiedType: const FullType(String),
    );
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ClaimLedgerEntry object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimLedgerEntryBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'claimId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.claimId = valueDes;
          break;
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
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ClaimLedgerEntry deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimLedgerEntryBuilder();
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

