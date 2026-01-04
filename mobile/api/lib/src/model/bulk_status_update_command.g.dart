// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bulk_status_update_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BulkStatusUpdateCommand extends BulkStatusUpdateCommand {
  @override
  final String? comment;
  @override
  final String status;
  @override
  final BuiltList<int> receiptIds;

  factory _$BulkStatusUpdateCommand(
          [void Function(BulkStatusUpdateCommandBuilder)? updates]) =>
      (new BulkStatusUpdateCommandBuilder()..update(updates))._build();

  _$BulkStatusUpdateCommand._(
      {this.comment, required this.status, required this.receiptIds})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        status, r'BulkStatusUpdateCommand', 'status');
    BuiltValueNullFieldError.checkNotNull(
        receiptIds, r'BulkStatusUpdateCommand', 'receiptIds');
  }

  @override
  BulkStatusUpdateCommand rebuild(
          void Function(BulkStatusUpdateCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BulkStatusUpdateCommandBuilder toBuilder() =>
      new BulkStatusUpdateCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BulkStatusUpdateCommand &&
        comment == other.comment &&
        status == other.status &&
        receiptIds == other.receiptIds;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, comment.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, receiptIds.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BulkStatusUpdateCommand')
          ..add('comment', comment)
          ..add('status', status)
          ..add('receiptIds', receiptIds))
        .toString();
  }
}

class BulkStatusUpdateCommandBuilder
    implements
        Builder<BulkStatusUpdateCommand, BulkStatusUpdateCommandBuilder> {
  _$BulkStatusUpdateCommand? _$v;

  String? _comment;
  String? get comment => _$this._comment;
  set comment(String? comment) => _$this._comment = comment;

  String? _status;
  String? get status => _$this._status;
  set status(String? status) => _$this._status = status;

  ListBuilder<int>? _receiptIds;
  ListBuilder<int> get receiptIds =>
      _$this._receiptIds ??= new ListBuilder<int>();
  set receiptIds(ListBuilder<int>? receiptIds) =>
      _$this._receiptIds = receiptIds;

  BulkStatusUpdateCommandBuilder() {
    BulkStatusUpdateCommand._defaults(this);
  }

  BulkStatusUpdateCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _comment = $v.comment;
      _status = $v.status;
      _receiptIds = $v.receiptIds.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BulkStatusUpdateCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$BulkStatusUpdateCommand;
  }

  @override
  void update(void Function(BulkStatusUpdateCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BulkStatusUpdateCommand build() => _build();

  _$BulkStatusUpdateCommand _build() {
    _$BulkStatusUpdateCommand _$result;
    try {
      _$result = _$v ??
          new _$BulkStatusUpdateCommand._(
              comment: comment,
              status: BuiltValueNullFieldError.checkNotNull(
                  status, r'BulkStatusUpdateCommand', 'status'),
              receiptIds: receiptIds.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'receiptIds';
        receiptIds.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'BulkStatusUpdateCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
