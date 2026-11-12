//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_claim_command.g.dart';

/// A claim's own fields, for create and for a full replace on update. An optional field that is omitted, null or \"\" is stored as not set. Receipts and ledger entries are managed through their own endpoints and are not part of this command.
///
/// Properties:
/// * [name] - Required, non-blank after trimming, at most 255 characters
/// * [referenceNumber] - At most 100 characters
/// * [type] - Free text, at most 100 characters
/// * [status] - One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED; anything else is a 400 keyed status
/// * [submittedDate] - Calendar day (yyyy-MM-dd)
/// * [dueDate] - Calendar day (yyyy-MM-dd)
/// * [claimedAmount] - Decimal string, zero or more, at most two decimal places
/// * [notes] - At most 10000 characters
@BuiltValue()
abstract class UpsertClaimCommand implements Built<UpsertClaimCommand, UpsertClaimCommandBuilder> {
  /// Required, non-blank after trimming, at most 255 characters
  @BuiltValueField(wireName: r'name')
  String get name;

  /// At most 100 characters
  @BuiltValueField(wireName: r'referenceNumber')
  String? get referenceNumber;

  /// Free text, at most 100 characters
  @BuiltValueField(wireName: r'type')
  String? get type;

  /// One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED; anything else is a 400 keyed status
  @BuiltValueField(wireName: r'status')
  String get status;

  /// Calendar day (yyyy-MM-dd)
  @BuiltValueField(wireName: r'submittedDate')
  String? get submittedDate;

  /// Calendar day (yyyy-MM-dd)
  @BuiltValueField(wireName: r'dueDate')
  String? get dueDate;

  /// Decimal string, zero or more, at most two decimal places
  @BuiltValueField(wireName: r'claimedAmount')
  String? get claimedAmount;

  /// At most 10000 characters
  @BuiltValueField(wireName: r'notes')
  String? get notes;

  UpsertClaimCommand._();

  factory UpsertClaimCommand([void updates(UpsertClaimCommandBuilder b)]) = _$UpsertClaimCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertClaimCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertClaimCommand> get serializer => _$UpsertClaimCommandSerializer();
}

class _$UpsertClaimCommandSerializer implements PrimitiveSerializer<UpsertClaimCommand> {
  @override
  final Iterable<Type> types = const [UpsertClaimCommand, _$UpsertClaimCommand];

  @override
  final String wireName = r'UpsertClaimCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertClaimCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    if (object.referenceNumber != null) {
      yield r'referenceNumber';
      yield serializers.serialize(
        object.referenceNumber,
        specifiedType: const FullType(String),
      );
    }
    if (object.type != null) {
      yield r'type';
      yield serializers.serialize(
        object.type,
        specifiedType: const FullType(String),
      );
    }
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(String),
    );
    if (object.submittedDate != null) {
      yield r'submittedDate';
      yield serializers.serialize(
        object.submittedDate,
        specifiedType: const FullType(String),
      );
    }
    if (object.dueDate != null) {
      yield r'dueDate';
      yield serializers.serialize(
        object.dueDate,
        specifiedType: const FullType(String),
      );
    }
    if (object.claimedAmount != null) {
      yield r'claimedAmount';
      yield serializers.serialize(
        object.claimedAmount,
        specifiedType: const FullType(String),
      );
    }
    if (object.notes != null) {
      yield r'notes';
      yield serializers.serialize(
        object.notes,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertClaimCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertClaimCommandBuilder result,
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
        case r'referenceNumber':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.referenceNumber = valueDes;
          break;
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.type = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.status = valueDes;
          break;
        case r'submittedDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.submittedDate = valueDes;
          break;
        case r'dueDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.dueDate = valueDes;
          break;
        case r'claimedAmount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.claimedAmount = valueDes;
          break;
        case r'notes':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.notes = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertClaimCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertClaimCommandBuilder();
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

