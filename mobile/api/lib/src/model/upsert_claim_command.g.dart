// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_claim_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertClaimCommand extends UpsertClaimCommand {
  @override
  final String name;
  @override
  final String? referenceNumber;
  @override
  final String? type;
  @override
  final String status;
  @override
  final String? submittedDate;
  @override
  final String? dueDate;
  @override
  final String? claimedAmount;
  @override
  final String? notes;

  factory _$UpsertClaimCommand(
          [void Function(UpsertClaimCommandBuilder)? updates]) =>
      (UpsertClaimCommandBuilder()..update(updates))._build();

  _$UpsertClaimCommand._(
      {required this.name,
      this.referenceNumber,
      this.type,
      required this.status,
      this.submittedDate,
      this.dueDate,
      this.claimedAmount,
      this.notes})
      : super._();
  @override
  UpsertClaimCommand rebuild(
          void Function(UpsertClaimCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertClaimCommandBuilder toBuilder() =>
      UpsertClaimCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertClaimCommand &&
        name == other.name &&
        referenceNumber == other.referenceNumber &&
        type == other.type &&
        status == other.status &&
        submittedDate == other.submittedDate &&
        dueDate == other.dueDate &&
        claimedAmount == other.claimedAmount &&
        notes == other.notes;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, referenceNumber.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, submittedDate.hashCode);
    _$hash = $jc(_$hash, dueDate.hashCode);
    _$hash = $jc(_$hash, claimedAmount.hashCode);
    _$hash = $jc(_$hash, notes.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertClaimCommand')
          ..add('name', name)
          ..add('referenceNumber', referenceNumber)
          ..add('type', type)
          ..add('status', status)
          ..add('submittedDate', submittedDate)
          ..add('dueDate', dueDate)
          ..add('claimedAmount', claimedAmount)
          ..add('notes', notes))
        .toString();
  }
}

class UpsertClaimCommandBuilder
    implements Builder<UpsertClaimCommand, UpsertClaimCommandBuilder> {
  _$UpsertClaimCommand? _$v;

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

  UpsertClaimCommandBuilder() {
    UpsertClaimCommand._defaults(this);
  }

  UpsertClaimCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _referenceNumber = $v.referenceNumber;
      _type = $v.type;
      _status = $v.status;
      _submittedDate = $v.submittedDate;
      _dueDate = $v.dueDate;
      _claimedAmount = $v.claimedAmount;
      _notes = $v.notes;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertClaimCommand other) {
    _$v = other as _$UpsertClaimCommand;
  }

  @override
  void update(void Function(UpsertClaimCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertClaimCommand build() => _build();

  _$UpsertClaimCommand _build() {
    final _$result = _$v ??
        _$UpsertClaimCommand._(
          name: BuiltValueNullFieldError.checkNotNull(
              name, r'UpsertClaimCommand', 'name'),
          referenceNumber: referenceNumber,
          type: type,
          status: BuiltValueNullFieldError.checkNotNull(
              status, r'UpsertClaimCommand', 'status'),
          submittedDate: submittedDate,
          dueDate: dueDate,
          claimedAmount: claimedAmount,
          notes: notes,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
