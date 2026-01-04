// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_comment_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertCommentCommand extends UpsertCommentCommand {
  @override
  final String comment;
  @override
  final int receiptId;
  @override
  final int? userId;

  factory _$UpsertCommentCommand(
          [void Function(UpsertCommentCommandBuilder)? updates]) =>
      (new UpsertCommentCommandBuilder()..update(updates))._build();

  _$UpsertCommentCommand._(
      {required this.comment, required this.receiptId, this.userId})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        comment, r'UpsertCommentCommand', 'comment');
    BuiltValueNullFieldError.checkNotNull(
        receiptId, r'UpsertCommentCommand', 'receiptId');
  }

  @override
  UpsertCommentCommand rebuild(
          void Function(UpsertCommentCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertCommentCommandBuilder toBuilder() =>
      new UpsertCommentCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertCommentCommand &&
        comment == other.comment &&
        receiptId == other.receiptId &&
        userId == other.userId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, comment.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertCommentCommand')
          ..add('comment', comment)
          ..add('receiptId', receiptId)
          ..add('userId', userId))
        .toString();
  }
}

class UpsertCommentCommandBuilder
    implements Builder<UpsertCommentCommand, UpsertCommentCommandBuilder> {
  _$UpsertCommentCommand? _$v;

  String? _comment;
  String? get comment => _$this._comment;
  set comment(String? comment) => _$this._comment = comment;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  UpsertCommentCommandBuilder() {
    UpsertCommentCommand._defaults(this);
  }

  UpsertCommentCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _comment = $v.comment;
      _receiptId = $v.receiptId;
      _userId = $v.userId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertCommentCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertCommentCommand;
  }

  @override
  void update(void Function(UpsertCommentCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertCommentCommand build() => _build();

  _$UpsertCommentCommand _build() {
    final _$result = _$v ??
        new _$UpsertCommentCommand._(
            comment: BuiltValueNullFieldError.checkNotNull(
                comment, r'UpsertCommentCommand', 'comment'),
            receiptId: BuiltValueNullFieldError.checkNotNull(
                receiptId, r'UpsertCommentCommand', 'receiptId'),
            userId: userId);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
