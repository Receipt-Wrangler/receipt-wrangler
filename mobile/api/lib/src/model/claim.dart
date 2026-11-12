//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/claim_ledger_entry.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim.g.dart';

/// A reimbursement claim: a bundle of one group's receipts submitted together, with a ledger of what happened after submission. Every computed figure and list covers only what the caller may see: linked receipts hidden from them (paid-by visibility, member isolation) are left out of receiptIds and of every receipt figure, and in a group with isolated members, ledger entries recorded by a member the caller cannot see are left out too. The response never says that anything was left out.
///
/// Properties:
/// * [id] 
/// * [groupId] - The group the claim and all its receipts belong to
/// * [name] 
/// * [referenceNumber] - A travel authorization or voucher number. Absent when not set.
/// * [type] - Free text, e.g. \"Travel voucher\" or \"TQSA\". Absent when not set.
/// * [status] - One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED. Set by hand; any status may follow any other. A plain string, not an enum, so a client must tolerate values it does not know.
/// * [awaitingSettlement] - Whether the status is one the server treats as waiting on the other party (SUBMITTED, PARTIALLY_PAID, DISPUTED). Drives daysWaiting and the list summary, so a client need not hard-code the set.
/// * [submittedDate] - Calendar day the claim was submitted (yyyy-MM-dd). Absent when not set.
/// * [dueDate] - Calendar day a response or payment is due (yyyy-MM-dd). Absent when not set.
/// * [claimedAmount] - The amount explicitly entered as claimed, as a decimal string. Absent when not set, in which case effectiveClaimedAmount falls back to receiptsTotal.
/// * [notes] - Free-text notes. Absent when not set.
/// * [createdAt] - RFC 3339 instant
/// * [updatedAt] - RFC 3339 instant
/// * [createdBy] - The user who created the claim
/// * [receiptIds] - Ids of the linked receipts the caller may see, ascending. Fetch the receipts themselves through the receipt endpoints. Always present; empty when none.
/// * [receiptCount] - Number of entries in receiptIds
/// * [receiptsTotal] - Sum of the amounts of the receipts in receiptIds
/// * [effectiveClaimedAmount] - claimedAmount when set, otherwise receiptsTotal
/// * [paymentsTotal] - Sum of the amounts of the PAYMENT ledger entries
/// * [recoupmentsTotal] - Sum of the amounts of the RECOUPMENT ledger entries
/// * [paidAmount] - paymentsTotal minus recoupmentsTotal; negative when more was taken back than paid
/// * [outstandingAmount] - effectiveClaimedAmount minus paidAmount. Negative when more was received than claimed.
/// * [daysWaiting] - Calendar days from submittedDate to today in the app time zone, when awaitingSettlement is true and submittedDate is set; otherwise absent. Never negative.
/// * [ledgerEntries] - The claim's ledger, oldest first (by date, then id). Always present; empty when none.
@BuiltValue()
abstract class Claim implements Built<Claim, ClaimBuilder> {
  @BuiltValueField(wireName: r'id')
  int get id;

  /// The group the claim and all its receipts belong to
  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  @BuiltValueField(wireName: r'name')
  String get name;

  /// A travel authorization or voucher number. Absent when not set.
  @BuiltValueField(wireName: r'referenceNumber')
  String? get referenceNumber;

  /// Free text, e.g. \"Travel voucher\" or \"TQSA\". Absent when not set.
  @BuiltValueField(wireName: r'type')
  String? get type;

  /// One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED. Set by hand; any status may follow any other. A plain string, not an enum, so a client must tolerate values it does not know.
  @BuiltValueField(wireName: r'status')
  String get status;

  /// Whether the status is one the server treats as waiting on the other party (SUBMITTED, PARTIALLY_PAID, DISPUTED). Drives daysWaiting and the list summary, so a client need not hard-code the set.
  @BuiltValueField(wireName: r'awaitingSettlement')
  bool get awaitingSettlement;

  /// Calendar day the claim was submitted (yyyy-MM-dd). Absent when not set.
  @BuiltValueField(wireName: r'submittedDate')
  String? get submittedDate;

  /// Calendar day a response or payment is due (yyyy-MM-dd). Absent when not set.
  @BuiltValueField(wireName: r'dueDate')
  String? get dueDate;

  /// The amount explicitly entered as claimed, as a decimal string. Absent when not set, in which case effectiveClaimedAmount falls back to receiptsTotal.
  @BuiltValueField(wireName: r'claimedAmount')
  String? get claimedAmount;

  /// Free-text notes. Absent when not set.
  @BuiltValueField(wireName: r'notes')
  String? get notes;

  /// RFC 3339 instant
  @BuiltValueField(wireName: r'createdAt')
  String get createdAt;

  /// RFC 3339 instant
  @BuiltValueField(wireName: r'updatedAt')
  String get updatedAt;

  /// The user who created the claim
  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  /// Ids of the linked receipts the caller may see, ascending. Fetch the receipts themselves through the receipt endpoints. Always present; empty when none.
  @BuiltValueField(wireName: r'receiptIds')
  BuiltList<int> get receiptIds;

  /// Number of entries in receiptIds
  @BuiltValueField(wireName: r'receiptCount')
  int get receiptCount;

  /// Sum of the amounts of the receipts in receiptIds
  @BuiltValueField(wireName: r'receiptsTotal')
  String get receiptsTotal;

  /// claimedAmount when set, otherwise receiptsTotal
  @BuiltValueField(wireName: r'effectiveClaimedAmount')
  String get effectiveClaimedAmount;

  /// Sum of the amounts of the PAYMENT ledger entries
  @BuiltValueField(wireName: r'paymentsTotal')
  String get paymentsTotal;

  /// Sum of the amounts of the RECOUPMENT ledger entries
  @BuiltValueField(wireName: r'recoupmentsTotal')
  String get recoupmentsTotal;

  /// paymentsTotal minus recoupmentsTotal; negative when more was taken back than paid
  @BuiltValueField(wireName: r'paidAmount')
  String get paidAmount;

  /// effectiveClaimedAmount minus paidAmount. Negative when more was received than claimed.
  @BuiltValueField(wireName: r'outstandingAmount')
  String get outstandingAmount;

  /// Calendar days from submittedDate to today in the app time zone, when awaitingSettlement is true and submittedDate is set; otherwise absent. Never negative.
  @BuiltValueField(wireName: r'daysWaiting')
  int? get daysWaiting;

  /// The claim's ledger, oldest first (by date, then id). Always present; empty when none.
  @BuiltValueField(wireName: r'ledgerEntries')
  BuiltList<ClaimLedgerEntry> get ledgerEntries;

  Claim._();

  factory Claim([void updates(ClaimBuilder b)]) = _$Claim;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Claim> get serializer => _$ClaimSerializer();
}

class _$ClaimSerializer implements PrimitiveSerializer<Claim> {
  @override
  final Iterable<Type> types = const [Claim, _$Claim];

  @override
  final String wireName = r'Claim';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Claim object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(int),
    );
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
    yield r'awaitingSettlement';
    yield serializers.serialize(
      object.awaitingSettlement,
      specifiedType: const FullType(bool),
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
    yield r'receiptIds';
    yield serializers.serialize(
      object.receiptIds,
      specifiedType: const FullType(BuiltList, [FullType(int)]),
    );
    yield r'receiptCount';
    yield serializers.serialize(
      object.receiptCount,
      specifiedType: const FullType(int),
    );
    yield r'receiptsTotal';
    yield serializers.serialize(
      object.receiptsTotal,
      specifiedType: const FullType(String),
    );
    yield r'effectiveClaimedAmount';
    yield serializers.serialize(
      object.effectiveClaimedAmount,
      specifiedType: const FullType(String),
    );
    yield r'paymentsTotal';
    yield serializers.serialize(
      object.paymentsTotal,
      specifiedType: const FullType(String),
    );
    yield r'recoupmentsTotal';
    yield serializers.serialize(
      object.recoupmentsTotal,
      specifiedType: const FullType(String),
    );
    yield r'paidAmount';
    yield serializers.serialize(
      object.paidAmount,
      specifiedType: const FullType(String),
    );
    yield r'outstandingAmount';
    yield serializers.serialize(
      object.outstandingAmount,
      specifiedType: const FullType(String),
    );
    if (object.daysWaiting != null) {
      yield r'daysWaiting';
      yield serializers.serialize(
        object.daysWaiting,
        specifiedType: const FullType(int),
      );
    }
    yield r'ledgerEntries';
    yield serializers.serialize(
      object.ledgerEntries,
      specifiedType: const FullType(BuiltList, [FullType(ClaimLedgerEntry)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    Claim object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimBuilder result,
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
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
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
        case r'awaitingSettlement':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.awaitingSettlement = valueDes;
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
        case r'receiptIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.receiptIds.replace(valueDes);
          break;
        case r'receiptCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptCount = valueDes;
          break;
        case r'receiptsTotal':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.receiptsTotal = valueDes;
          break;
        case r'effectiveClaimedAmount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.effectiveClaimedAmount = valueDes;
          break;
        case r'paymentsTotal':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.paymentsTotal = valueDes;
          break;
        case r'recoupmentsTotal':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.recoupmentsTotal = valueDes;
          break;
        case r'paidAmount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.paidAmount = valueDes;
          break;
        case r'outstandingAmount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.outstandingAmount = valueDes;
          break;
        case r'daysWaiting':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.daysWaiting = valueDes;
          break;
        case r'ledgerEntries':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ClaimLedgerEntry)]),
          ) as BuiltList<ClaimLedgerEntry>;
          result.ledgerEntries.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  Claim deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimBuilder();
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

