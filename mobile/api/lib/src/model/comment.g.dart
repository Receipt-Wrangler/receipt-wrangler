// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'comment.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Comment extends Comment {
  @override
  final String? additionalInfo;
  @override
  final String comment;
  @override
  final int? commentId;
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final int id;
  @override
  final int receiptId;
  @override
  final String? updatedAt;
  @override
  final int userId;

  factory _$Comment([void Function(CommentBuilder)? updates]) =>
      (new CommentBuilder()..update(updates))._build();

  _$Comment._(
      {this.additionalInfo,
      required this.comment,
      this.commentId,
      this.createdAt,
      this.createdBy,
      required this.id,
      required this.receiptId,
      this.updatedAt,
      required this.userId})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(comment, r'Comment', 'comment');
    BuiltValueNullFieldError.checkNotNull(id, r'Comment', 'id');
    BuiltValueNullFieldError.checkNotNull(receiptId, r'Comment', 'receiptId');
    BuiltValueNullFieldError.checkNotNull(userId, r'Comment', 'userId');
  }

  @override
  Comment rebuild(void Function(CommentBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CommentBuilder toBuilder() => new CommentBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Comment &&
        additionalInfo == other.additionalInfo &&
        comment == other.comment &&
        commentId == other.commentId &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        id == other.id &&
        receiptId == other.receiptId &&
        updatedAt == other.updatedAt &&
        userId == other.userId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, additionalInfo.hashCode);
    _$hash = $jc(_$hash, comment.hashCode);
    _$hash = $jc(_$hash, commentId.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Comment')
          ..add('additionalInfo', additionalInfo)
          ..add('comment', comment)
          ..add('commentId', commentId)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('id', id)
          ..add('receiptId', receiptId)
          ..add('updatedAt', updatedAt)
          ..add('userId', userId))
        .toString();
  }
}

class CommentBuilder implements Builder<Comment, CommentBuilder> {
  _$Comment? _$v;

  String? _additionalInfo;
  String? get additionalInfo => _$this._additionalInfo;
  set additionalInfo(String? additionalInfo) =>
      _$this._additionalInfo = additionalInfo;

  String? _comment;
  String? get comment => _$this._comment;
  set comment(String? comment) => _$this._comment = comment;

  int? _commentId;
  int? get commentId => _$this._commentId;
  set commentId(int? commentId) => _$this._commentId = commentId;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  CommentBuilder() {
    Comment._defaults(this);
  }

  CommentBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _additionalInfo = $v.additionalInfo;
      _comment = $v.comment;
      _commentId = $v.commentId;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _id = $v.id;
      _receiptId = $v.receiptId;
      _updatedAt = $v.updatedAt;
      _userId = $v.userId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Comment other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Comment;
  }

  @override
  void update(void Function(CommentBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Comment build() => _build();

  _$Comment _build() {
    final _$result = _$v ??
        new _$Comment._(
            additionalInfo: additionalInfo,
            comment: BuiltValueNullFieldError.checkNotNull(
                comment, r'Comment', 'comment'),
            commentId: commentId,
            createdAt: createdAt,
            createdBy: createdBy,
            id: BuiltValueNullFieldError.checkNotNull(id, r'Comment', 'id'),
            receiptId: BuiltValueNullFieldError.checkNotNull(
                receiptId, r'Comment', 'receiptId'),
            updatedAt: updatedAt,
            userId: BuiltValueNullFieldError.checkNotNull(
                userId, r'Comment', 'userId'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
