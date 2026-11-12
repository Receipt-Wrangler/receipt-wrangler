// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_claim_ledger_entry_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertClaimLedgerEntryCommand extends UpsertClaimLedgerEntryCommand {
  @override
  final String type;
  @override
  final String date;
  @override
  final String? amount;
  @override
  final String? note;

  factory _$UpsertClaimLedgerEntryCommand(
          [void Function(UpsertClaimLedgerEntryCommandBuilder)? updates]) =>
      (UpsertClaimLedgerEntryCommandBuilder()..update(updates))._build();

  _$UpsertClaimLedgerEntryCommand._(
      {required this.type, required this.date, this.amount, this.note})
      : super._();
  @override
  UpsertClaimLedgerEntryCommand rebuild(
          void Function(UpsertClaimLedgerEntryCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertClaimLedgerEntryCommandBuilder toBuilder() =>
      UpsertClaimLedgerEntryCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertClaimLedgerEntryCommand &&
        type == other.type &&
        date == other.date &&
        amount == other.amount &&
        note == other.note;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, note.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertClaimLedgerEntryCommand')
          ..add('type', type)
          ..add('date', date)
          ..add('amount', amount)
          ..add('note', note))
        .toString();
  }
}

class UpsertClaimLedgerEntryCommandBuilder
    implements
        Builder<UpsertClaimLedgerEntryCommand,
            UpsertClaimLedgerEntryCommandBuilder> {
  _$UpsertClaimLedgerEntryCommand? _$v;

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

  UpsertClaimLedgerEntryCommandBuilder() {
    UpsertClaimLedgerEntryCommand._defaults(this);
  }

  UpsertClaimLedgerEntryCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _type = $v.type;
      _date = $v.date;
      _amount = $v.amount;
      _note = $v.note;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertClaimLedgerEntryCommand other) {
    _$v = other as _$UpsertClaimLedgerEntryCommand;
  }

  @override
  void update(void Function(UpsertClaimLedgerEntryCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertClaimLedgerEntryCommand build() => _build();

  _$UpsertClaimLedgerEntryCommand _build() {
    final _$result = _$v ??
        _$UpsertClaimLedgerEntryCommand._(
          type: BuiltValueNullFieldError.checkNotNull(
              type, r'UpsertClaimLedgerEntryCommand', 'type'),
          date: BuiltValueNullFieldError.checkNotNull(
              date, r'UpsertClaimLedgerEntryCommand', 'date'),
          amount: amount,
          note: note,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
