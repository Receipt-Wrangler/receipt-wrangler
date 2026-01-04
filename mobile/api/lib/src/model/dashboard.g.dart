// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'dashboard.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Dashboard extends Dashboard {
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final int id;
  @override
  final String name;
  @override
  final int? groupId;
  @override
  final int userId;
  @override
  final String? updatedAt;
  @override
  final BuiltList<Widget>? widgets;

  factory _$Dashboard([void Function(DashboardBuilder)? updates]) =>
      (new DashboardBuilder()..update(updates))._build();

  _$Dashboard._(
      {this.createdAt,
      this.createdBy,
      required this.id,
      required this.name,
      this.groupId,
      required this.userId,
      this.updatedAt,
      this.widgets})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'Dashboard', 'id');
    BuiltValueNullFieldError.checkNotNull(name, r'Dashboard', 'name');
    BuiltValueNullFieldError.checkNotNull(userId, r'Dashboard', 'userId');
  }

  @override
  Dashboard rebuild(void Function(DashboardBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  DashboardBuilder toBuilder() => new DashboardBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Dashboard &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        id == other.id &&
        name == other.name &&
        groupId == other.groupId &&
        userId == other.userId &&
        updatedAt == other.updatedAt &&
        widgets == other.widgets;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, widgets.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Dashboard')
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('id', id)
          ..add('name', name)
          ..add('groupId', groupId)
          ..add('userId', userId)
          ..add('updatedAt', updatedAt)
          ..add('widgets', widgets))
        .toString();
  }
}

class DashboardBuilder implements Builder<Dashboard, DashboardBuilder> {
  _$Dashboard? _$v;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  ListBuilder<Widget>? _widgets;
  ListBuilder<Widget> get widgets =>
      _$this._widgets ??= new ListBuilder<Widget>();
  set widgets(ListBuilder<Widget>? widgets) => _$this._widgets = widgets;

  DashboardBuilder() {
    Dashboard._defaults(this);
  }

  DashboardBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _id = $v.id;
      _name = $v.name;
      _groupId = $v.groupId;
      _userId = $v.userId;
      _updatedAt = $v.updatedAt;
      _widgets = $v.widgets?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Dashboard other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Dashboard;
  }

  @override
  void update(void Function(DashboardBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Dashboard build() => _build();

  _$Dashboard _build() {
    _$Dashboard _$result;
    try {
      _$result = _$v ??
          new _$Dashboard._(
              createdAt: createdAt,
              createdBy: createdBy,
              id: BuiltValueNullFieldError.checkNotNull(id, r'Dashboard', 'id'),
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'Dashboard', 'name'),
              groupId: groupId,
              userId: BuiltValueNullFieldError.checkNotNull(
                  userId, r'Dashboard', 'userId'),
              updatedAt: updatedAt,
              widgets: _widgets?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'widgets';
        _widgets?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'Dashboard', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
