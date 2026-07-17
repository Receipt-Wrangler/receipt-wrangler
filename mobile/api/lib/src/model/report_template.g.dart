// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_template.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReportTemplate extends ReportTemplate {
  @override
  final int configurationVersion;
  @override
  final ReportRequestCommand configuration;
  @override
  final String name;
  @override
  final BuiltList<String>? allowedActions;
  @override
  final int id;
  @override
  final String createdAt;
  @override
  final int? createdBy;
  @override
  final String? createdByString;
  @override
  final String? updatedAt;

  factory _$ReportTemplate([void Function(ReportTemplateBuilder)? updates]) =>
      (ReportTemplateBuilder()..update(updates))._build();

  _$ReportTemplate._(
      {required this.configurationVersion,
      required this.configuration,
      required this.name,
      this.allowedActions,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._();
  @override
  ReportTemplate rebuild(void Function(ReportTemplateBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportTemplateBuilder toBuilder() => ReportTemplateBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportTemplate &&
        configurationVersion == other.configurationVersion &&
        configuration == other.configuration &&
        name == other.name &&
        allowedActions == other.allowedActions &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, configurationVersion.hashCode);
    _$hash = $jc(_$hash, configuration.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, allowedActions.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, createdByString.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportTemplate')
          ..add('configurationVersion', configurationVersion)
          ..add('configuration', configuration)
          ..add('name', name)
          ..add('allowedActions', allowedActions)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class ReportTemplateBuilder
    implements
        Builder<ReportTemplate, ReportTemplateBuilder>,
        BaseModelBuilder {
  _$ReportTemplate? _$v;

  int? _configurationVersion;
  int? get configurationVersion => _$this._configurationVersion;
  set configurationVersion(covariant int? configurationVersion) =>
      _$this._configurationVersion = configurationVersion;

  ReportRequestCommandBuilder? _configuration;
  ReportRequestCommandBuilder get configuration =>
      _$this._configuration ??= ReportRequestCommandBuilder();
  set configuration(covariant ReportRequestCommandBuilder? configuration) =>
      _$this._configuration = configuration;

  String? _name;
  String? get name => _$this._name;
  set name(covariant String? name) => _$this._name = name;

  ListBuilder<String>? _allowedActions;
  ListBuilder<String> get allowedActions =>
      _$this._allowedActions ??= ListBuilder<String>();
  set allowedActions(covariant ListBuilder<String>? allowedActions) =>
      _$this._allowedActions = allowedActions;

  int? _id;
  int? get id => _$this._id;
  set id(covariant int? id) => _$this._id = id;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(covariant String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(covariant int? createdBy) => _$this._createdBy = createdBy;

  String? _createdByString;
  String? get createdByString => _$this._createdByString;
  set createdByString(covariant String? createdByString) =>
      _$this._createdByString = createdByString;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(covariant String? updatedAt) => _$this._updatedAt = updatedAt;

  ReportTemplateBuilder() {
    ReportTemplate._defaults(this);
  }

  ReportTemplateBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _configurationVersion = $v.configurationVersion;
      _configuration = $v.configuration.toBuilder();
      _name = $v.name;
      _allowedActions = $v.allowedActions?.toBuilder();
      _id = $v.id;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _createdByString = $v.createdByString;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant ReportTemplate other) {
    _$v = other as _$ReportTemplate;
  }

  @override
  void update(void Function(ReportTemplateBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportTemplate build() => _build();

  _$ReportTemplate _build() {
    _$ReportTemplate _$result;
    try {
      _$result = _$v ??
          _$ReportTemplate._(
            configurationVersion: BuiltValueNullFieldError.checkNotNull(
                configurationVersion,
                r'ReportTemplate',
                'configurationVersion'),
            configuration: configuration.build(),
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'ReportTemplate', 'name'),
            allowedActions: _allowedActions?.build(),
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'ReportTemplate', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'ReportTemplate', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'configuration';
        configuration.build();

        _$failedField = 'allowedActions';
        _allowedActions?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReportTemplate', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
