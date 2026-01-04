// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'subject_line_regex.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$SubjectLineRegex extends SubjectLineRegex {
  @override
  final int id;
  @override
  final int groupSettingsId;
  @override
  final String regex;
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final String? updatedAt;

  factory _$SubjectLineRegex(
          [void Function(SubjectLineRegexBuilder)? updates]) =>
      (new SubjectLineRegexBuilder()..update(updates))._build();

  _$SubjectLineRegex._(
      {required this.id,
      required this.groupSettingsId,
      required this.regex,
      this.createdAt,
      this.createdBy,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'SubjectLineRegex', 'id');
    BuiltValueNullFieldError.checkNotNull(
        groupSettingsId, r'SubjectLineRegex', 'groupSettingsId');
    BuiltValueNullFieldError.checkNotNull(regex, r'SubjectLineRegex', 'regex');
  }

  @override
  SubjectLineRegex rebuild(void Function(SubjectLineRegexBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  SubjectLineRegexBuilder toBuilder() =>
      new SubjectLineRegexBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is SubjectLineRegex &&
        id == other.id &&
        groupSettingsId == other.groupSettingsId &&
        regex == other.regex &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, groupSettingsId.hashCode);
    _$hash = $jc(_$hash, regex.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'SubjectLineRegex')
          ..add('id', id)
          ..add('groupSettingsId', groupSettingsId)
          ..add('regex', regex)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class SubjectLineRegexBuilder
    implements Builder<SubjectLineRegex, SubjectLineRegexBuilder> {
  _$SubjectLineRegex? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  int? _groupSettingsId;
  int? get groupSettingsId => _$this._groupSettingsId;
  set groupSettingsId(int? groupSettingsId) =>
      _$this._groupSettingsId = groupSettingsId;

  String? _regex;
  String? get regex => _$this._regex;
  set regex(String? regex) => _$this._regex = regex;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  SubjectLineRegexBuilder() {
    SubjectLineRegex._defaults(this);
  }

  SubjectLineRegexBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _groupSettingsId = $v.groupSettingsId;
      _regex = $v.regex;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(SubjectLineRegex other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$SubjectLineRegex;
  }

  @override
  void update(void Function(SubjectLineRegexBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  SubjectLineRegex build() => _build();

  _$SubjectLineRegex _build() {
    final _$result = _$v ??
        new _$SubjectLineRegex._(
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'SubjectLineRegex', 'id'),
            groupSettingsId: BuiltValueNullFieldError.checkNotNull(
                groupSettingsId, r'SubjectLineRegex', 'groupSettingsId'),
            regex: BuiltValueNullFieldError.checkNotNull(
                regex, r'SubjectLineRegex', 'regex'),
            createdAt: createdAt,
            createdBy: createdBy,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
