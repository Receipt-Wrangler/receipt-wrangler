// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'group_filter.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$GroupFilter extends GroupFilter {
  @override
  final AssociatedGroup? associatedGroup;

  factory _$GroupFilter([void Function(GroupFilterBuilder)? updates]) =>
      (new GroupFilterBuilder()..update(updates))._build();

  _$GroupFilter._({this.associatedGroup}) : super._();

  @override
  GroupFilter rebuild(void Function(GroupFilterBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GroupFilterBuilder toBuilder() => new GroupFilterBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GroupFilter && associatedGroup == other.associatedGroup;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, associatedGroup.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'GroupFilter')
          ..add('associatedGroup', associatedGroup))
        .toString();
  }
}

class GroupFilterBuilder implements Builder<GroupFilter, GroupFilterBuilder> {
  _$GroupFilter? _$v;

  AssociatedGroup? _associatedGroup;
  AssociatedGroup? get associatedGroup => _$this._associatedGroup;
  set associatedGroup(AssociatedGroup? associatedGroup) =>
      _$this._associatedGroup = associatedGroup;

  GroupFilterBuilder() {
    GroupFilter._defaults(this);
  }

  GroupFilterBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _associatedGroup = $v.associatedGroup;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(GroupFilter other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$GroupFilter;
  }

  @override
  void update(void Function(GroupFilterBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GroupFilter build() => _build();

  _$GroupFilter _build() {
    final _$result =
        _$v ?? new _$GroupFilter._(associatedGroup: associatedGroup);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
