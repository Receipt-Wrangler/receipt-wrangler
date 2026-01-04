// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'widget.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Widget extends Widget {
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final int id;
  @override
  final String? name;
  @override
  final int dashboardId;
  @override
  final String? updatedAt;
  @override
  final WidgetType? widgetType;
  @override
  final BuiltMap<String, JsonObject?>? configuration;

  factory _$Widget([void Function(WidgetBuilder)? updates]) =>
      (new WidgetBuilder()..update(updates))._build();

  _$Widget._(
      {this.createdAt,
      this.createdBy,
      required this.id,
      this.name,
      required this.dashboardId,
      this.updatedAt,
      this.widgetType,
      this.configuration})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'Widget', 'id');
    BuiltValueNullFieldError.checkNotNull(
        dashboardId, r'Widget', 'dashboardId');
  }

  @override
  Widget rebuild(void Function(WidgetBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  WidgetBuilder toBuilder() => new WidgetBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Widget &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        id == other.id &&
        name == other.name &&
        dashboardId == other.dashboardId &&
        updatedAt == other.updatedAt &&
        widgetType == other.widgetType &&
        configuration == other.configuration;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, dashboardId.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, widgetType.hashCode);
    _$hash = $jc(_$hash, configuration.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Widget')
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('id', id)
          ..add('name', name)
          ..add('dashboardId', dashboardId)
          ..add('updatedAt', updatedAt)
          ..add('widgetType', widgetType)
          ..add('configuration', configuration))
        .toString();
  }
}

class WidgetBuilder implements Builder<Widget, WidgetBuilder> {
  _$Widget? _$v;

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

  int? _dashboardId;
  int? get dashboardId => _$this._dashboardId;
  set dashboardId(int? dashboardId) => _$this._dashboardId = dashboardId;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  WidgetType? _widgetType;
  WidgetType? get widgetType => _$this._widgetType;
  set widgetType(WidgetType? widgetType) => _$this._widgetType = widgetType;

  MapBuilder<String, JsonObject?>? _configuration;
  MapBuilder<String, JsonObject?> get configuration =>
      _$this._configuration ??= new MapBuilder<String, JsonObject?>();
  set configuration(MapBuilder<String, JsonObject?>? configuration) =>
      _$this._configuration = configuration;

  WidgetBuilder() {
    Widget._defaults(this);
  }

  WidgetBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _id = $v.id;
      _name = $v.name;
      _dashboardId = $v.dashboardId;
      _updatedAt = $v.updatedAt;
      _widgetType = $v.widgetType;
      _configuration = $v.configuration?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Widget other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Widget;
  }

  @override
  void update(void Function(WidgetBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Widget build() => _build();

  _$Widget _build() {
    _$Widget _$result;
    try {
      _$result = _$v ??
          new _$Widget._(
              createdAt: createdAt,
              createdBy: createdBy,
              id: BuiltValueNullFieldError.checkNotNull(id, r'Widget', 'id'),
              name: name,
              dashboardId: BuiltValueNullFieldError.checkNotNull(
                  dashboardId, r'Widget', 'dashboardId'),
              updatedAt: updatedAt,
              widgetType: widgetType,
              configuration: _configuration?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'configuration';
        _configuration?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'Widget', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
