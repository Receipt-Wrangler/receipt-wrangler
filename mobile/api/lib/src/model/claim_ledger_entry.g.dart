// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim_ledger_entry.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ClaimLedgerEntry extends ClaimLedgerEntry {
  @override
  final int id;
  @override
  final int claimId;
  @override
  final String type;
  @override
  final String date;
  @override
  final String? amount;
  @override
  final String? note;
  @override
  final String createdAt;
  @override
  final String updatedAt;
  @override
  final int? createdBy;

  factory _$ClaimLedgerEntry(
          [void Function(ClaimLedgerEntryBuilder)? updates]) =>
      (ClaimLedgerEntryBuilder()..update(updates))._build();

  _$ClaimLedgerEntry._(
      {required this.id,
      required this.claimId,
      required this.type,
      required this.date,
      this.amount,
      this.note,
      required this.createdAt,
      required this.updatedAt,
      this.createdBy})
      : super._();
  @override
  ClaimLedgerEntry rebuild(void Function(ClaimLedgerEntryBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimLedgerEntryBuilder toBuilder() =>
      ClaimLedgerEntryBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ClaimLedgerEntry &&
        id == other.id &&
        claimId == other.claimId &&
        type == other.type &&
        date == other.date &&
        amount == other.amount &&
        note == other.note &&
        createdAt == other.createdAt &&
        updatedAt == other.updatedAt &&
        createdBy == other.createdBy;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, claimId.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, note.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ClaimLedgerEntry')
          ..add('id', id)
          ..add('claimId', claimId)
          ..add('type', type)
          ..add('date', date)
          ..add('amount', amount)
          ..add('note', note)
          ..add('createdAt', createdAt)
          ..add('updatedAt', updatedAt)
          ..add('createdBy', createdBy))
        .toString();
  }
}

class ClaimLedgerEntryBuilder
    implements Builder<ClaimLedgerEntry, ClaimLedgerEntryBuilder> {
  _$ClaimLedgerEntry? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  int? _claimId;
  int? get claimId => _$this._claimId;
  set claimId(int? claimId) => _$this._claimId = claimId;

  String? _type;
  String? get type => _$this._type;
  set type(String? type) => _$this._type = type;

  String? _date;
  String? get date => _$this._date;
  set date(String? date) => _$this._date = date;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  String? _note;
  String? get note => _$this._note;
  set note(String? note) => _$this._note = note;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  ClaimLedgerEntryBuilder() {
    ClaimLedgerEntry._defaults(this);
  }

  ClaimLedgerEntryBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _claimId = $v.claimId;
      _type = $v.type;
      _date = $v.date;
      _amount = $v.amount;
      _note = $v.note;
      _createdAt = $v.createdAt;
      _updatedAt = $v.updatedAt;
      _createdBy = $v.createdBy;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ClaimLedgerEntry other) {
    _$v = other as _$ClaimLedgerEntry;
  }

  @override
  void update(void Function(ClaimLedgerEntryBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ClaimLedgerEntry build() => _build();

  _$ClaimLedgerEntry _build() {
    final _$result = _$v ??
        _$ClaimLedgerEntry._(
          id: BuiltValueNullFieldError.checkNotNull(
              id, r'ClaimLedgerEntry', 'id'),
          claimId: BuiltValueNullFieldError.checkNotNull(
              claimId, r'ClaimLedgerEntry', 'claimId'),
          type: BuiltValueNullFieldError.checkNotNull(
              type, r'ClaimLedgerEntry', 'type'),
          date: BuiltValueNullFieldError.checkNotNull(
              date, r'ClaimLedgerEntry', 'date'),
          amount: amount,
          note: note,
          createdAt: BuiltValueNullFieldError.checkNotNull(
              createdAt, r'ClaimLedgerEntry', 'createdAt'),
          updatedAt: BuiltValueNullFieldError.checkNotNull(
              updatedAt, r'ClaimLedgerEntry', 'updatedAt'),
          createdBy: createdBy,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
