// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_template_option.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReportTemplateOption extends ReportTemplateOption {
  @override
  final int id;
  @override
  final String name;
  @override
  final BuiltList<int> groupIds;

  factory _$ReportTemplateOption(
          [void Function(ReportTemplateOptionBuilder)? updates]) =>
      (ReportTemplateOptionBuilder()..update(updates))._build();

  _$ReportTemplateOption._(
      {required this.id, required this.name, required this.groupIds})
      : super._();
  @override
  ReportTemplateOption rebuild(
          void Function(ReportTemplateOptionBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportTemplateOptionBuilder toBuilder() =>
      ReportTemplateOptionBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportTemplateOption &&
        id == other.id &&
        name == other.name &&
        groupIds == other.groupIds;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, groupIds.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportTemplateOption')
          ..add('id', id)
          ..add('name', name)
          ..add('groupIds', groupIds))
        .toString();
  }
}

class ReportTemplateOptionBuilder
    implements Builder<ReportTemplateOption, ReportTemplateOptionBuilder> {
  _$ReportTemplateOption? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  ListBuilder<int>? _groupIds;
  ListBuilder<int> get groupIds => _$this._groupIds ??= ListBuilder<int>();
  set groupIds(ListBuilder<int>? groupIds) => _$this._groupIds = groupIds;

  ReportTemplateOptionBuilder() {
    ReportTemplateOption._defaults(this);
  }

  ReportTemplateOptionBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _name = $v.name;
      _groupIds = $v.groupIds.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportTemplateOption other) {
    _$v = other as _$ReportTemplateOption;
  }

  @override
  void update(void Function(ReportTemplateOptionBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportTemplateOption build() => _build();

  _$ReportTemplateOption _build() {
    _$ReportTemplateOption _$result;
    try {
      _$result = _$v ??
          _$ReportTemplateOption._(
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'ReportTemplateOption', 'id'),
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'ReportTemplateOption', 'name'),
            groupIds: groupIds.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'groupIds';
        groupIds.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReportTemplateOption', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
