// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'notification.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Notification extends Notification {
  @override
  final String body;
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final int id;
  @override
  final String title;
  @override
  final String type;
  @override
  final String? updatedAt;
  @override
  final int userId;

  factory _$Notification([void Function(NotificationBuilder)? updates]) =>
      (new NotificationBuilder()..update(updates))._build();

  _$Notification._(
      {required this.body,
      this.createdAt,
      this.createdBy,
      required this.id,
      required this.title,
      required this.type,
      this.updatedAt,
      required this.userId})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(body, r'Notification', 'body');
    BuiltValueNullFieldError.checkNotNull(id, r'Notification', 'id');
    BuiltValueNullFieldError.checkNotNull(title, r'Notification', 'title');
    BuiltValueNullFieldError.checkNotNull(type, r'Notification', 'type');
    BuiltValueNullFieldError.checkNotNull(userId, r'Notification', 'userId');
  }

  @override
  Notification rebuild(void Function(NotificationBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  NotificationBuilder toBuilder() => new NotificationBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Notification &&
        body == other.body &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        id == other.id &&
        title == other.title &&
        type == other.type &&
        updatedAt == other.updatedAt &&
        userId == other.userId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, body.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, title.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Notification')
          ..add('body', body)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('id', id)
          ..add('title', title)
          ..add('type', type)
          ..add('updatedAt', updatedAt)
          ..add('userId', userId))
        .toString();
  }
}

class NotificationBuilder
    implements Builder<Notification, NotificationBuilder> {
  _$Notification? _$v;

  String? _body;
  String? get body => _$this._body;
  set body(String? body) => _$this._body = body;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _title;
  String? get title => _$this._title;
  set title(String? title) => _$this._title = title;

  String? _type;
  String? get type => _$this._type;
  set type(String? type) => _$this._type = type;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  NotificationBuilder() {
    Notification._defaults(this);
  }

  NotificationBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _body = $v.body;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _id = $v.id;
      _title = $v.title;
      _type = $v.type;
      _updatedAt = $v.updatedAt;
      _userId = $v.userId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Notification other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Notification;
  }

  @override
  void update(void Function(NotificationBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Notification build() => _build();

  _$Notification _build() {
    final _$result = _$v ??
        new _$Notification._(
            body: BuiltValueNullFieldError.checkNotNull(
                body, r'Notification', 'body'),
            createdAt: createdAt,
            createdBy: createdBy,
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'Notification', 'id'),
            title: BuiltValueNullFieldError.checkNotNull(
                title, r'Notification', 'title'),
            type: BuiltValueNullFieldError.checkNotNull(
                type, r'Notification', 'type'),
            updatedAt: updatedAt,
            userId: BuiltValueNullFieldError.checkNotNull(
                userId, r'Notification', 'userId'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
