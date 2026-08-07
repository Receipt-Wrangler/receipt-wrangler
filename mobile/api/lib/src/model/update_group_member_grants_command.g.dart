// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'update_group_member_grants_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpdateGroupMemberGrantsCommand extends UpdateGroupMemberGrantsCommand {
  @override
  final BuiltList<int>? categoryGrants;
  @override
  final BuiltList<int>? tagGrants;

  factory _$UpdateGroupMemberGrantsCommand(
          [void Function(UpdateGroupMemberGrantsCommandBuilder)? updates]) =>
      (UpdateGroupMemberGrantsCommandBuilder()..update(updates))._build();

  _$UpdateGroupMemberGrantsCommand._({this.categoryGrants, this.tagGrants})
      : super._();
  @override
  UpdateGroupMemberGrantsCommand rebuild(
          void Function(UpdateGroupMemberGrantsCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpdateGroupMemberGrantsCommandBuilder toBuilder() =>
      UpdateGroupMemberGrantsCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpdateGroupMemberGrantsCommand &&
        categoryGrants == other.categoryGrants &&
        tagGrants == other.tagGrants;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, categoryGrants.hashCode);
    _$hash = $jc(_$hash, tagGrants.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpdateGroupMemberGrantsCommand')
          ..add('categoryGrants', categoryGrants)
          ..add('tagGrants', tagGrants))
        .toString();
  }
}

class UpdateGroupMemberGrantsCommandBuilder
    implements
        Builder<UpdateGroupMemberGrantsCommand,
            UpdateGroupMemberGrantsCommandBuilder> {
  _$UpdateGroupMemberGrantsCommand? _$v;

  ListBuilder<int>? _categoryGrants;
  ListBuilder<int> get categoryGrants =>
      _$this._categoryGrants ??= ListBuilder<int>();
  set categoryGrants(ListBuilder<int>? categoryGrants) =>
      _$this._categoryGrants = categoryGrants;

  ListBuilder<int>? _tagGrants;
  ListBuilder<int> get tagGrants => _$this._tagGrants ??= ListBuilder<int>();
  set tagGrants(ListBuilder<int>? tagGrants) => _$this._tagGrants = tagGrants;

  UpdateGroupMemberGrantsCommandBuilder() {
    UpdateGroupMemberGrantsCommand._defaults(this);
  }

  UpdateGroupMemberGrantsCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _categoryGrants = $v.categoryGrants?.toBuilder();
      _tagGrants = $v.tagGrants?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpdateGroupMemberGrantsCommand other) {
    _$v = other as _$UpdateGroupMemberGrantsCommand;
  }

  @override
  void update(void Function(UpdateGroupMemberGrantsCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpdateGroupMemberGrantsCommand build() => _build();

  _$UpdateGroupMemberGrantsCommand _build() {
    _$UpdateGroupMemberGrantsCommand _$result;
    try {
      _$result = _$v ??
          _$UpdateGroupMemberGrantsCommand._(
            categoryGrants: _categoryGrants?.build(),
            tagGrants: _tagGrants?.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'categoryGrants';
        _categoryGrants?.build();
        _$failedField = 'tagGrants';
        _tagGrants?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'UpdateGroupMemberGrantsCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
