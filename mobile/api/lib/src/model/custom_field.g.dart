// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'custom_field.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CustomField extends CustomField {
  @override
  final String name;
  @override
  final BuiltList<CustomFieldOption>? options;
  @override
  final String? description;
  @override
  final CustomFieldType type;
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

  factory _$CustomField([void Function(CustomFieldBuilder)? updates]) =>
      (new CustomFieldBuilder()..update(updates))._build();

  _$CustomField._(
      {required this.name,
      this.options,
      this.description,
      required this.type,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(name, r'CustomField', 'name');
    BuiltValueNullFieldError.checkNotNull(type, r'CustomField', 'type');
    BuiltValueNullFieldError.checkNotNull(id, r'CustomField', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'CustomField', 'createdAt');
  }

  @override
  CustomField rebuild(void Function(CustomFieldBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CustomFieldBuilder toBuilder() => new CustomFieldBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CustomField &&
        name == other.name &&
        options == other.options &&
        description == other.description &&
        type == other.type &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, options.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
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
    return (newBuiltValueToStringHelper(r'CustomField')
          ..add('name', name)
          ..add('options', options)
          ..add('description', description)
          ..add('type', type)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class CustomFieldBuilder
    implements Builder<CustomField, CustomFieldBuilder>, BaseModelBuilder {
  _$CustomField? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(covariant String? name) => _$this._name = name;

  ListBuilder<CustomFieldOption>? _options;
  ListBuilder<CustomFieldOption> get options =>
      _$this._options ??= new ListBuilder<CustomFieldOption>();
  set options(covariant ListBuilder<CustomFieldOption>? options) =>
      _$this._options = options;

  String? _description;
  String? get description => _$this._description;
  set description(covariant String? description) =>
      _$this._description = description;

  CustomFieldType? _type;
  CustomFieldType? get type => _$this._type;
  set type(covariant CustomFieldType? type) => _$this._type = type;

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

  CustomFieldBuilder() {
    CustomField._defaults(this);
  }

  CustomFieldBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _options = $v.options?.toBuilder();
      _description = $v.description;
      _type = $v.type;
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
  void replace(covariant CustomField other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$CustomField;
  }

  @override
  void update(void Function(CustomFieldBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CustomField build() => _build();

  _$CustomField _build() {
    _$CustomField _$result;
    try {
      _$result = _$v ??
          new _$CustomField._(
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'CustomField', 'name'),
              options: _options?.build(),
              description: description,
              type: BuiltValueNullFieldError.checkNotNull(
                  type, r'CustomField', 'type'),
              id: BuiltValueNullFieldError.checkNotNull(
                  id, r'CustomField', 'id'),
              createdAt: BuiltValueNullFieldError.checkNotNull(
                  createdAt, r'CustomField', 'createdAt'),
              createdBy: createdBy,
              createdByString: createdByString,
              updatedAt: updatedAt);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'options';
        _options?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'CustomField', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
