// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_template_grant.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReportTemplateGrant extends ReportTemplateGrant {
  @override
  final int reportTemplateId;
  @override
  final BuiltList<String> permissions;

  factory _$ReportTemplateGrant(
          [void Function(ReportTemplateGrantBuilder)? updates]) =>
      (ReportTemplateGrantBuilder()..update(updates))._build();

  _$ReportTemplateGrant._(
      {required this.reportTemplateId, required this.permissions})
      : super._();
  @override
  ReportTemplateGrant rebuild(
          void Function(ReportTemplateGrantBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportTemplateGrantBuilder toBuilder() =>
      ReportTemplateGrantBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportTemplateGrant &&
        reportTemplateId == other.reportTemplateId &&
        permissions == other.permissions;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, reportTemplateId.hashCode);
    _$hash = $jc(_$hash, permissions.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportTemplateGrant')
          ..add('reportTemplateId', reportTemplateId)
          ..add('permissions', permissions))
        .toString();
  }
}

class ReportTemplateGrantBuilder
    implements Builder<ReportTemplateGrant, ReportTemplateGrantBuilder> {
  _$ReportTemplateGrant? _$v;

  int? _reportTemplateId;
  int? get reportTemplateId => _$this._reportTemplateId;
  set reportTemplateId(int? reportTemplateId) =>
      _$this._reportTemplateId = reportTemplateId;

  ListBuilder<String>? _permissions;
  ListBuilder<String> get permissions =>
      _$this._permissions ??= ListBuilder<String>();
  set permissions(ListBuilder<String>? permissions) =>
      _$this._permissions = permissions;

  ReportTemplateGrantBuilder() {
    ReportTemplateGrant._defaults(this);
  }

  ReportTemplateGrantBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _reportTemplateId = $v.reportTemplateId;
      _permissions = $v.permissions.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportTemplateGrant other) {
    _$v = other as _$ReportTemplateGrant;
  }

  @override
  void update(void Function(ReportTemplateGrantBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportTemplateGrant build() => _build();

  _$ReportTemplateGrant _build() {
    _$ReportTemplateGrant _$result;
    try {
      _$result = _$v ??
          _$ReportTemplateGrant._(
            reportTemplateId: BuiltValueNullFieldError.checkNotNull(
                reportTemplateId, r'ReportTemplateGrant', 'reportTemplateId'),
            permissions: permissions.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'permissions';
        permissions.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReportTemplateGrant', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
