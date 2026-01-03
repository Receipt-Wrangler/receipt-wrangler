// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_dashboard_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertDashboardCommand extends UpsertDashboardCommand {
  @override
  final String name;
  @override
  final String groupId;
  @override
  final BuiltList<UpsertWidgetCommand>? widgets;

  factory _$UpsertDashboardCommand(
          [void Function(UpsertDashboardCommandBuilder)? updates]) =>
      (new UpsertDashboardCommandBuilder()..update(updates))._build();

  _$UpsertDashboardCommand._(
      {required this.name, required this.groupId, this.widgets})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        name, r'UpsertDashboardCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        groupId, r'UpsertDashboardCommand', 'groupId');
  }

  @override
  UpsertDashboardCommand rebuild(
          void Function(UpsertDashboardCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertDashboardCommandBuilder toBuilder() =>
      new UpsertDashboardCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertDashboardCommand &&
        name == other.name &&
        groupId == other.groupId &&
        widgets == other.widgets;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, widgets.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertDashboardCommand')
          ..add('name', name)
          ..add('groupId', groupId)
          ..add('widgets', widgets))
        .toString();
  }
}

class UpsertDashboardCommandBuilder
    implements Builder<UpsertDashboardCommand, UpsertDashboardCommandBuilder> {
  _$UpsertDashboardCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _groupId;
  String? get groupId => _$this._groupId;
  set groupId(String? groupId) => _$this._groupId = groupId;

  ListBuilder<UpsertWidgetCommand>? _widgets;
  ListBuilder<UpsertWidgetCommand> get widgets =>
      _$this._widgets ??= new ListBuilder<UpsertWidgetCommand>();
  set widgets(ListBuilder<UpsertWidgetCommand>? widgets) =>
      _$this._widgets = widgets;

  UpsertDashboardCommandBuilder() {
    UpsertDashboardCommand._defaults(this);
  }

  UpsertDashboardCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _groupId = $v.groupId;
      _widgets = $v.widgets?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertDashboardCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertDashboardCommand;
  }

  @override
  void update(void Function(UpsertDashboardCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertDashboardCommand build() => _build();

  _$UpsertDashboardCommand _build() {
    _$UpsertDashboardCommand _$result;
    try {
      _$result = _$v ??
          new _$UpsertDashboardCommand._(
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'UpsertDashboardCommand', 'name'),
              groupId: BuiltValueNullFieldError.checkNotNull(
                  groupId, r'UpsertDashboardCommand', 'groupId'),
              widgets: _widgets?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'widgets';
        _widgets?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'UpsertDashboardCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
