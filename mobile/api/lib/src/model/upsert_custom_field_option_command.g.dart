// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_custom_field_option_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertCustomFieldOptionCommand extends UpsertCustomFieldOptionCommand {
  @override
  final String? value;
  @override
  final int customFieldId;

  factory _$UpsertCustomFieldOptionCommand(
          [void Function(UpsertCustomFieldOptionCommandBuilder)? updates]) =>
      (new UpsertCustomFieldOptionCommandBuilder()..update(updates))._build();

  _$UpsertCustomFieldOptionCommand._({this.value, required this.customFieldId})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        customFieldId, r'UpsertCustomFieldOptionCommand', 'customFieldId');
  }

  @override
  UpsertCustomFieldOptionCommand rebuild(
          void Function(UpsertCustomFieldOptionCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertCustomFieldOptionCommandBuilder toBuilder() =>
      new UpsertCustomFieldOptionCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertCustomFieldOptionCommand &&
        value == other.value &&
        customFieldId == other.customFieldId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, value.hashCode);
    _$hash = $jc(_$hash, customFieldId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertCustomFieldOptionCommand')
          ..add('value', value)
          ..add('customFieldId', customFieldId))
        .toString();
  }
}

class UpsertCustomFieldOptionCommandBuilder
    implements
        Builder<UpsertCustomFieldOptionCommand,
            UpsertCustomFieldOptionCommandBuilder> {
  _$UpsertCustomFieldOptionCommand? _$v;

  String? _value;
  String? get value => _$this._value;
  set value(String? value) => _$this._value = value;

  int? _customFieldId;
  int? get customFieldId => _$this._customFieldId;
  set customFieldId(int? customFieldId) =>
      _$this._customFieldId = customFieldId;

  UpsertCustomFieldOptionCommandBuilder() {
    UpsertCustomFieldOptionCommand._defaults(this);
  }

  UpsertCustomFieldOptionCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _value = $v.value;
      _customFieldId = $v.customFieldId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertCustomFieldOptionCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertCustomFieldOptionCommand;
  }

  @override
  void update(void Function(UpsertCustomFieldOptionCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertCustomFieldOptionCommand build() => _build();

  _$UpsertCustomFieldOptionCommand _build() {
    final _$result = _$v ??
        new _$UpsertCustomFieldOptionCommand._(
            value: value,
            customFieldId: BuiltValueNullFieldError.checkNotNull(customFieldId,
                r'UpsertCustomFieldOptionCommand', 'customFieldId'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
