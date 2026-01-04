// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_widget_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertWidgetCommand extends UpsertWidgetCommand {
  @override
  final String? name;
  @override
  final WidgetType widgetType;
  @override
  final BuiltMap<String, JsonObject?>? configuration;

  factory _$UpsertWidgetCommand(
          [void Function(UpsertWidgetCommandBuilder)? updates]) =>
      (new UpsertWidgetCommandBuilder()..update(updates))._build();

  _$UpsertWidgetCommand._(
      {this.name, required this.widgetType, this.configuration})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        widgetType, r'UpsertWidgetCommand', 'widgetType');
  }

  @override
  UpsertWidgetCommand rebuild(
          void Function(UpsertWidgetCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertWidgetCommandBuilder toBuilder() =>
      new UpsertWidgetCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertWidgetCommand &&
        name == other.name &&
        widgetType == other.widgetType &&
        configuration == other.configuration;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, widgetType.hashCode);
    _$hash = $jc(_$hash, configuration.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertWidgetCommand')
          ..add('name', name)
          ..add('widgetType', widgetType)
          ..add('configuration', configuration))
        .toString();
  }
}

class UpsertWidgetCommandBuilder
    implements Builder<UpsertWidgetCommand, UpsertWidgetCommandBuilder> {
  _$UpsertWidgetCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  WidgetType? _widgetType;
  WidgetType? get widgetType => _$this._widgetType;
  set widgetType(WidgetType? widgetType) => _$this._widgetType = widgetType;

  MapBuilder<String, JsonObject?>? _configuration;
  MapBuilder<String, JsonObject?> get configuration =>
      _$this._configuration ??= new MapBuilder<String, JsonObject?>();
  set configuration(MapBuilder<String, JsonObject?>? configuration) =>
      _$this._configuration = configuration;

  UpsertWidgetCommandBuilder() {
    UpsertWidgetCommand._defaults(this);
  }

  UpsertWidgetCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _widgetType = $v.widgetType;
      _configuration = $v.configuration?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertWidgetCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertWidgetCommand;
  }

  @override
  void update(void Function(UpsertWidgetCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertWidgetCommand build() => _build();

  _$UpsertWidgetCommand _build() {
    _$UpsertWidgetCommand _$result;
    try {
      _$result = _$v ??
          new _$UpsertWidgetCommand._(
              name: name,
              widgetType: BuiltValueNullFieldError.checkNotNull(
                  widgetType, r'UpsertWidgetCommand', 'widgetType'),
              configuration: _configuration?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'configuration';
        _configuration?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'UpsertWidgetCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
