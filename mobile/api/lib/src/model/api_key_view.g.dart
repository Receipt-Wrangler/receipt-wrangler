// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_key_view.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ApiKeyView extends ApiKeyView {
  @override
  final String? id;
  @override
  final DateTime? createdAt;
  @override
  final DateTime? updatedAt;
  @override
  final int? createdBy;
  @override
  final String? createdByString;
  @override
  final String? name;
  @override
  final String? description;
  @override
  final int? userId;
  @override
  final String? scope;
  @override
  final DateTime? lastUsedAt;

  factory _$ApiKeyView([void Function(ApiKeyViewBuilder)? updates]) =>
      (new ApiKeyViewBuilder()..update(updates))._build();

  _$ApiKeyView._(
      {this.id,
      this.createdAt,
      this.updatedAt,
      this.createdBy,
      this.createdByString,
      this.name,
      this.description,
      this.userId,
      this.scope,
      this.lastUsedAt})
      : super._();

  @override
  ApiKeyView rebuild(void Function(ApiKeyViewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ApiKeyViewBuilder toBuilder() => new ApiKeyViewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ApiKeyView &&
        id == other.id &&
        createdAt == other.createdAt &&
        updatedAt == other.updatedAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        name == other.name &&
        description == other.description &&
        userId == other.userId &&
        scope == other.scope &&
        lastUsedAt == other.lastUsedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, createdByString.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jc(_$hash, scope.hashCode);
    _$hash = $jc(_$hash, lastUsedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ApiKeyView')
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('updatedAt', updatedAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('name', name)
          ..add('description', description)
          ..add('userId', userId)
          ..add('scope', scope)
          ..add('lastUsedAt', lastUsedAt))
        .toString();
  }
}

class ApiKeyViewBuilder implements Builder<ApiKeyView, ApiKeyViewBuilder> {
  _$ApiKeyView? _$v;

  String? _id;
  String? get id => _$this._id;
  set id(String? id) => _$this._id = id;

  DateTime? _createdAt;
  DateTime? get createdAt => _$this._createdAt;
  set createdAt(DateTime? createdAt) => _$this._createdAt = createdAt;

  DateTime? _updatedAt;
  DateTime? get updatedAt => _$this._updatedAt;
  set updatedAt(DateTime? updatedAt) => _$this._updatedAt = updatedAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  String? _createdByString;
  String? get createdByString => _$this._createdByString;
  set createdByString(String? createdByString) =>
      _$this._createdByString = createdByString;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  String? _scope;
  String? get scope => _$this._scope;
  set scope(String? scope) => _$this._scope = scope;

  DateTime? _lastUsedAt;
  DateTime? get lastUsedAt => _$this._lastUsedAt;
  set lastUsedAt(DateTime? lastUsedAt) => _$this._lastUsedAt = lastUsedAt;

  ApiKeyViewBuilder() {
    ApiKeyView._defaults(this);
  }

  ApiKeyViewBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _createdAt = $v.createdAt;
      _updatedAt = $v.updatedAt;
      _createdBy = $v.createdBy;
      _createdByString = $v.createdByString;
      _name = $v.name;
      _description = $v.description;
      _userId = $v.userId;
      _scope = $v.scope;
      _lastUsedAt = $v.lastUsedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ApiKeyView other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ApiKeyView;
  }

  @override
  void update(void Function(ApiKeyViewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ApiKeyView build() => _build();

  _$ApiKeyView _build() {
    final _$result = _$v ??
        new _$ApiKeyView._(
            id: id,
            createdAt: createdAt,
            updatedAt: updatedAt,
            createdBy: createdBy,
            createdByString: createdByString,
            name: name,
            description: description,
            userId: userId,
            scope: scope,
            lastUsedAt: lastUsedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
