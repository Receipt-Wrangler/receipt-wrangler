// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'icon.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Icon extends Icon {
  @override
  final String value;
  @override
  final String displayValue;

  factory _$Icon([void Function(IconBuilder)? updates]) =>
      (new IconBuilder()..update(updates))._build();

  _$Icon._({required this.value, required this.displayValue}) : super._() {
    BuiltValueNullFieldError.checkNotNull(value, r'Icon', 'value');
    BuiltValueNullFieldError.checkNotNull(
        displayValue, r'Icon', 'displayValue');
  }

  @override
  Icon rebuild(void Function(IconBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  IconBuilder toBuilder() => new IconBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Icon &&
        value == other.value &&
        displayValue == other.displayValue;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, value.hashCode);
    _$hash = $jc(_$hash, displayValue.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Icon')
          ..add('value', value)
          ..add('displayValue', displayValue))
        .toString();
  }
}

class IconBuilder implements Builder<Icon, IconBuilder> {
  _$Icon? _$v;

  String? _value;
  String? get value => _$this._value;
  set value(String? value) => _$this._value = value;

  String? _displayValue;
  String? get displayValue => _$this._displayValue;
  set displayValue(String? displayValue) => _$this._displayValue = displayValue;

  IconBuilder() {
    Icon._defaults(this);
  }

  IconBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _value = $v.value;
      _displayValue = $v.displayValue;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Icon other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Icon;
  }

  @override
  void update(void Function(IconBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Icon build() => _build();

  _$Icon _build() {
    final _$result = _$v ??
        new _$Icon._(
            value:
                BuiltValueNullFieldError.checkNotNull(value, r'Icon', 'value'),
            displayValue: BuiltValueNullFieldError.checkNotNull(
                displayValue, r'Icon', 'displayValue'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
