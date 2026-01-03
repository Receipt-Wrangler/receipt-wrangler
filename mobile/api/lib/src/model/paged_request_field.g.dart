// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_request_field.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PagedRequestField extends PagedRequestField {
  @override
  final FilterOperation? operation;
  @override
  final PagedRequestFieldValue? value;

  factory _$PagedRequestField(
          [void Function(PagedRequestFieldBuilder)? updates]) =>
      (new PagedRequestFieldBuilder()..update(updates))._build();

  _$PagedRequestField._({this.operation, this.value}) : super._();

  @override
  PagedRequestField rebuild(void Function(PagedRequestFieldBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PagedRequestFieldBuilder toBuilder() =>
      new PagedRequestFieldBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PagedRequestField &&
        operation == other.operation &&
        value == other.value;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, operation.hashCode);
    _$hash = $jc(_$hash, value.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PagedRequestField')
          ..add('operation', operation)
          ..add('value', value))
        .toString();
  }
}

class PagedRequestFieldBuilder
    implements Builder<PagedRequestField, PagedRequestFieldBuilder> {
  _$PagedRequestField? _$v;

  FilterOperation? _operation;
  FilterOperation? get operation => _$this._operation;
  set operation(FilterOperation? operation) => _$this._operation = operation;

  PagedRequestFieldValueBuilder? _value;
  PagedRequestFieldValueBuilder get value =>
      _$this._value ??= new PagedRequestFieldValueBuilder();
  set value(PagedRequestFieldValueBuilder? value) => _$this._value = value;

  PagedRequestFieldBuilder() {
    PagedRequestField._defaults(this);
  }

  PagedRequestFieldBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _operation = $v.operation;
      _value = $v.value?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PagedRequestField other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PagedRequestField;
  }

  @override
  void update(void Function(PagedRequestFieldBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PagedRequestField build() => _build();

  _$PagedRequestField _build() {
    _$PagedRequestField _$result;
    try {
      _$result = _$v ??
          new _$PagedRequestField._(
              operation: operation, value: _value?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'value';
        _value?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'PagedRequestField', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
