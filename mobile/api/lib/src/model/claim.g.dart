// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Claim extends Claim {
  @override
  final int id;
  @override
  final int groupId;
  @override
  final String name;
  @override
  final String? referenceNumber;
  @override
  final String? type;
  @override
  final String status;
  @override
  final bool awaitingSettlement;
  @override
  final String? submittedDate;
  @override
  final String? dueDate;
  @override
  final String? claimedAmount;
  @override
  final String? notes;
  @override
  final String createdAt;
  @override
  final String updatedAt;
  @override
  final int? createdBy;
  @override
  final BuiltList<int> receiptIds;
  @override
  final int receiptCount;
  @override
  final String receiptsTotal;
  @override
  final String effectiveClaimedAmount;
  @override
  final String paymentsTotal;
  @override
  final String recoupmentsTotal;
  @override
  final String paidAmount;
  @override
  final String outstandingAmount;
  @override
  final int? daysWaiting;
  @override
  final BuiltList<ClaimLedgerEntry> ledgerEntries;

  factory _$Claim([void Function(ClaimBuilder)? updates]) =>
      (ClaimBuilder()..update(updates))._build();

  _$Claim._(
      {required this.id,
      required this.groupId,
      required this.name,
      this.referenceNumber,
      this.type,
      required this.status,
      required this.awaitingSettlement,
      this.submittedDate,
      this.dueDate,
      this.claimedAmount,
      this.notes,
      required this.createdAt,
      required this.updatedAt,
      this.createdBy,
      required this.receiptIds,
      required this.receiptCount,
      required this.receiptsTotal,
      required this.effectiveClaimedAmount,
      required this.paymentsTotal,
      required this.recoupmentsTotal,
      required this.paidAmount,
      required this.outstandingAmount,
      this.daysWaiting,
      required this.ledgerEntries})
      : super._();
  @override
  Claim rebuild(void Function(ClaimBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimBuilder toBuilder() => ClaimBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Claim &&
        id == other.id &&
        groupId == other.groupId &&
        name == other.name &&
        referenceNumber == other.referenceNumber &&
        type == other.type &&
        status == other.status &&
        awaitingSettlement == other.awaitingSettlement &&
        submittedDate == other.submittedDate &&
        dueDate == other.dueDate &&
        claimedAmount == other.claimedAmount &&
        notes == other.notes &&
        createdAt == other.createdAt &&
        updatedAt == other.updatedAt &&
        createdBy == other.createdBy &&
        receiptIds == other.receiptIds &&
        receiptCount == other.receiptCount &&
        receiptsTotal == other.receiptsTotal &&
        effectiveClaimedAmount == other.effectiveClaimedAmount &&
        paymentsTotal == other.paymentsTotal &&
        recoupmentsTotal == other.recoupmentsTotal &&
        paidAmount == other.paidAmount &&
        outstandingAmount == other.outstandingAmount &&
        daysWaiting == other.daysWaiting &&
        ledgerEntries == other.ledgerEntries;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, referenceNumber.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, awaitingSettlement.hashCode);
    _$hash = $jc(_$hash, submittedDate.hashCode);
    _$hash = $jc(_$hash, dueDate.hashCode);
    _$hash = $jc(_$hash, claimedAmount.hashCode);
    _$hash = $jc(_$hash, notes.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, receiptIds.hashCode);
    _$hash = $jc(_$hash, receiptCount.hashCode);
    _$hash = $jc(_$hash, receiptsTotal.hashCode);
    _$hash = $jc(_$hash, effectiveClaimedAmount.hashCode);
    _$hash = $jc(_$hash, paymentsTotal.hashCode);
    _$hash = $jc(_$hash, recoupmentsTotal.hashCode);
    _$hash = $jc(_$hash, paidAmount.hashCode);
    _$hash = $jc(_$hash, outstandingAmount.hashCode);
    _$hash = $jc(_$hash, daysWaiting.hashCode);
    _$hash = $jc(_$hash, ledgerEntries.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Claim')
          ..add('id', id)
          ..add('groupId', groupId)
          ..add('name', name)
          ..add('referenceNumber', referenceNumber)
          ..add('type', type)
          ..add('status', status)
          ..add('awaitingSettlement', awaitingSettlement)
          ..add('submittedDate', submittedDate)
          ..add('dueDate', dueDate)
          ..add('claimedAmount', claimedAmount)
          ..add('notes', notes)
          ..add('createdAt', createdAt)
          ..add('updatedAt', updatedAt)
          ..add('createdBy', createdBy)
          ..add('receiptIds', receiptIds)
          ..add('receiptCount', receiptCount)
          ..add('receiptsTotal', receiptsTotal)
          ..add('effectiveClaimedAmount', effectiveClaimedAmount)
          ..add('paymentsTotal', paymentsTotal)
          ..add('recoupmentsTotal', recoupmentsTotal)
          ..add('paidAmount', paidAmount)
          ..add('outstandingAmount', outstandingAmount)
          ..add('daysWaiting', daysWaiting)
          ..add('ledgerEntries', ledgerEntries))
        .toString();
  }
}

class ClaimBuilder implements Builder<Claim, ClaimBuilder> {
  _$Claim? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _referenceNumber;
  String? get referenceNumber => _$this._referenceNumber;
  set referenceNumber(String? referenceNumber) =>
      _$this._referenceNumber = referenceNumber;

  String? _type;
  String? get type => _$this._type;
  set type(String? type) => _$this._type = type;

  String? _status;
  String? get status => _$this._status;
  set status(String? status) => _$this._status = status;

  bool? _awaitingSettlement;
  bool? get awaitingSettlement => _$this._awaitingSettlement;
  set awaitingSettlement(bool? awaitingSettlement) =>
      _$this._awaitingSettlement = awaitingSettlement;

  String? _submittedDate;
  String? get submittedDate => _$this._submittedDate;
  set submittedDate(String? submittedDate) =>
      _$this._submittedDate = submittedDate;

  String? _dueDate;
  String? get dueDate => _$this._dueDate;
  set dueDate(String? dueDate) => _$this._dueDate = dueDate;

  String? _claimedAmount;
  String? get claimedAmount => _$this._claimedAmount;
  set claimedAmount(String? claimedAmount) =>
      _$this._claimedAmount = claimedAmount;

  String? _notes;
  String? get notes => _$this._notes;
  set notes(String? notes) => _$this._notes = notes;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  ListBuilder<int>? _receiptIds;
  ListBuilder<int> get receiptIds => _$this._receiptIds ??= ListBuilder<int>();
  set receiptIds(ListBuilder<int>? receiptIds) =>
      _$this._receiptIds = receiptIds;

  int? _receiptCount;
  int? get receiptCount => _$this._receiptCount;
  set receiptCount(int? receiptCount) => _$this._receiptCount = receiptCount;

  String? _receiptsTotal;
  String? get receiptsTotal => _$this._receiptsTotal;
  set receiptsTotal(String? receiptsTotal) =>
      _$this._receiptsTotal = receiptsTotal;

  String? _effectiveClaimedAmount;
  String? get effectiveClaimedAmount => _$this._effectiveClaimedAmount;
  set effectiveClaimedAmount(String? effectiveClaimedAmount) =>
      _$this._effectiveClaimedAmount = effectiveClaimedAmount;

  String? _paymentsTotal;
  String? get paymentsTotal => _$this._paymentsTotal;
  set paymentsTotal(String? paymentsTotal) =>
      _$this._paymentsTotal = paymentsTotal;

  String? _recoupmentsTotal;
  String? get recoupmentsTotal => _$this._recoupmentsTotal;
  set recoupmentsTotal(String? recoupmentsTotal) =>
      _$this._recoupmentsTotal = recoupmentsTotal;

  String? _paidAmount;
  String? get paidAmount => _$this._paidAmount;
  set paidAmount(String? paidAmount) => _$this._paidAmount = paidAmount;

  String? _outstandingAmount;
  String? get outstandingAmount => _$this._outstandingAmount;
  set outstandingAmount(String? outstandingAmount) =>
      _$this._outstandingAmount = outstandingAmount;

  int? _daysWaiting;
  int? get daysWaiting => _$this._daysWaiting;
  set daysWaiting(int? daysWaiting) => _$this._daysWaiting = daysWaiting;

  ListBuilder<ClaimLedgerEntry>? _ledgerEntries;
  ListBuilder<ClaimLedgerEntry> get ledgerEntries =>
      _$this._ledgerEntries ??= ListBuilder<ClaimLedgerEntry>();
  set ledgerEntries(ListBuilder<ClaimLedgerEntry>? ledgerEntries) =>
      _$this._ledgerEntries = ledgerEntries;

  ClaimBuilder() {
    Claim._defaults(this);
  }

  ClaimBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _groupId = $v.groupId;
      _name = $v.name;
      _referenceNumber = $v.referenceNumber;
      _type = $v.type;
      _status = $v.status;
      _awaitingSettlement = $v.awaitingSettlement;
      _submittedDate = $v.submittedDate;
      _dueDate = $v.dueDate;
      _claimedAmount = $v.claimedAmount;
      _notes = $v.notes;
      _createdAt = $v.createdAt;
      _updatedAt = $v.updatedAt;
      _createdBy = $v.createdBy;
      _receiptIds = $v.receiptIds.toBuilder();
      _receiptCount = $v.receiptCount;
      _receiptsTotal = $v.receiptsTotal;
      _effectiveClaimedAmount = $v.effectiveClaimedAmount;
      _paymentsTotal = $v.paymentsTotal;
      _recoupmentsTotal = $v.recoupmentsTotal;
      _paidAmount = $v.paidAmount;
      _outstandingAmount = $v.outstandingAmount;
      _daysWaiting = $v.daysWaiting;
      _ledgerEntries = $v.ledgerEntries.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Claim other) {
    _$v = other as _$Claim;
  }

  @override
  void update(void Function(ClaimBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Claim build() => _build();

  _$Claim _build() {
    _$Claim _$result;
    try {
      _$result = _$v ??
          _$Claim._(
            id: BuiltValueNullFieldError.checkNotNull(id, r'Claim', 'id'),
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'Claim', 'groupId'),
            name: BuiltValueNullFieldError.checkNotNull(name, r'Claim', 'name'),
            referenceNumber: referenceNumber,
            type: type,
            status: BuiltValueNullFieldError.checkNotNull(
                status, r'Claim', 'status'),
            awaitingSettlement: BuiltValueNullFieldError.checkNotNull(
                awaitingSettlement, r'Claim', 'awaitingSettlement'),
            submittedDate: submittedDate,
            dueDate: dueDate,
            claimedAmount: claimedAmount,
            notes: notes,
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'Claim', 'createdAt'),
            updatedAt: BuiltValueNullFieldError.checkNotNull(
                updatedAt, r'Claim', 'updatedAt'),
            createdBy: createdBy,
            receiptIds: receiptIds.build(),
            receiptCount: BuiltValueNullFieldError.checkNotNull(
                receiptCount, r'Claim', 'receiptCount'),
            receiptsTotal: BuiltValueNullFieldError.checkNotNull(
                receiptsTotal, r'Claim', 'receiptsTotal'),
            effectiveClaimedAmount: BuiltValueNullFieldError.checkNotNull(
                effectiveClaimedAmount, r'Claim', 'effectiveClaimedAmount'),
            paymentsTotal: BuiltValueNullFieldError.checkNotNull(
                paymentsTotal, r'Claim', 'paymentsTotal'),
            recoupmentsTotal: BuiltValueNullFieldError.checkNotNull(
                recoupmentsTotal, r'Claim', 'recoupmentsTotal'),
            paidAmount: BuiltValueNullFieldError.checkNotNull(
                paidAmount, r'Claim', 'paidAmount'),
            outstandingAmount: BuiltValueNullFieldError.checkNotNull(
                outstandingAmount, r'Claim', 'outstandingAmount'),
            daysWaiting: daysWaiting,
            ledgerEntries: ledgerEntries.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'receiptIds';
        receiptIds.build();

        _$failedField = 'ledgerEntries';
        ledgerEntries.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(r'Claim', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
