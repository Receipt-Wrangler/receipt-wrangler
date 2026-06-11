// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'permission_descriptor.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PermissionDescriptor extends PermissionDescriptor {
  @override
  final Permission key;
  @override
  final String label;
  @override
  final String description;
  @override
  final String category;
  @override
  final PermissionScope scope;

  factory _$PermissionDescriptor(
          [void Function(PermissionDescriptorBuilder)? updates]) =>
      (PermissionDescriptorBuilder()..update(updates))._build();

  _$PermissionDescriptor._(
      {required this.key,
      required this.label,
      required this.description,
      required this.category,
      required this.scope})
      : super._();
  @override
  PermissionDescriptor rebuild(
          void Function(PermissionDescriptorBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PermissionDescriptorBuilder toBuilder() =>
      PermissionDescriptorBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PermissionDescriptor &&
        key == other.key &&
        label == other.label &&
        description == other.description &&
        category == other.category &&
        scope == other.scope;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, key.hashCode);
    _$hash = $jc(_$hash, label.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, category.hashCode);
    _$hash = $jc(_$hash, scope.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PermissionDescriptor')
          ..add('key', key)
          ..add('label', label)
          ..add('description', description)
          ..add('category', category)
          ..add('scope', scope))
        .toString();
  }
}

class PermissionDescriptorBuilder
    implements Builder<PermissionDescriptor, PermissionDescriptorBuilder> {
  _$PermissionDescriptor? _$v;

  Permission? _key;
  Permission? get key => _$this._key;
  set key(Permission? key) => _$this._key = key;

  String? _label;
  String? get label => _$this._label;
  set label(String? label) => _$this._label = label;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

  String? _category;
  String? get category => _$this._category;
  set category(String? category) => _$this._category = category;

  PermissionScope? _scope;
  PermissionScope? get scope => _$this._scope;
  set scope(PermissionScope? scope) => _$this._scope = scope;

  PermissionDescriptorBuilder() {
    PermissionDescriptor._defaults(this);
  }

  PermissionDescriptorBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _key = $v.key;
      _label = $v.label;
      _description = $v.description;
      _category = $v.category;
      _scope = $v.scope;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PermissionDescriptor other) {
    _$v = other as _$PermissionDescriptor;
  }

  @override
  void update(void Function(PermissionDescriptorBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PermissionDescriptor build() => _build();

  _$PermissionDescriptor _build() {
    final _$result = _$v ??
        _$PermissionDescriptor._(
          key: BuiltValueNullFieldError.checkNotNull(
              key, r'PermissionDescriptor', 'key'),
          label: BuiltValueNullFieldError.checkNotNull(
              label, r'PermissionDescriptor', 'label'),
          description: BuiltValueNullFieldError.checkNotNull(
              description, r'PermissionDescriptor', 'description'),
          category: BuiltValueNullFieldError.checkNotNull(
              category, r'PermissionDescriptor', 'category'),
          scope: BuiltValueNullFieldError.checkNotNull(
              scope, r'PermissionDescriptor', 'scope'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
