// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'group_receipt_settings.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$GroupReceiptSettings extends GroupReceiptSettings {
  @override
  final bool? hideReceiptTags;
  @override
  final bool? hideReceiptCategories;
  @override
  final int groupId;
  @override
  final bool? hideItemCategories;
  @override
  final bool? hideImages;
  @override
  final bool? hideItemTags;
  @override
  final bool? hideComments;
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

  factory _$GroupReceiptSettings(
          [void Function(GroupReceiptSettingsBuilder)? updates]) =>
      (new GroupReceiptSettingsBuilder()..update(updates))._build();

  _$GroupReceiptSettings._(
      {this.hideReceiptTags,
      this.hideReceiptCategories,
      required this.groupId,
      this.hideItemCategories,
      this.hideImages,
      this.hideItemTags,
      this.hideComments,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        groupId, r'GroupReceiptSettings', 'groupId');
    BuiltValueNullFieldError.checkNotNull(id, r'GroupReceiptSettings', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'GroupReceiptSettings', 'createdAt');
  }

  @override
  GroupReceiptSettings rebuild(
          void Function(GroupReceiptSettingsBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GroupReceiptSettingsBuilder toBuilder() =>
      new GroupReceiptSettingsBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GroupReceiptSettings &&
        hideReceiptTags == other.hideReceiptTags &&
        hideReceiptCategories == other.hideReceiptCategories &&
        groupId == other.groupId &&
        hideItemCategories == other.hideItemCategories &&
        hideImages == other.hideImages &&
        hideItemTags == other.hideItemTags &&
        hideComments == other.hideComments &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, hideReceiptTags.hashCode);
    _$hash = $jc(_$hash, hideReceiptCategories.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, hideItemCategories.hashCode);
    _$hash = $jc(_$hash, hideImages.hashCode);
    _$hash = $jc(_$hash, hideItemTags.hashCode);
    _$hash = $jc(_$hash, hideComments.hashCode);
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
    return (newBuiltValueToStringHelper(r'GroupReceiptSettings')
          ..add('hideReceiptTags', hideReceiptTags)
          ..add('hideReceiptCategories', hideReceiptCategories)
          ..add('groupId', groupId)
          ..add('hideItemCategories', hideItemCategories)
          ..add('hideImages', hideImages)
          ..add('hideItemTags', hideItemTags)
          ..add('hideComments', hideComments)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class GroupReceiptSettingsBuilder
    implements
        Builder<GroupReceiptSettings, GroupReceiptSettingsBuilder>,
        BaseModelBuilder {
  _$GroupReceiptSettings? _$v;

  bool? _hideReceiptTags;
  bool? get hideReceiptTags => _$this._hideReceiptTags;
  set hideReceiptTags(covariant bool? hideReceiptTags) =>
      _$this._hideReceiptTags = hideReceiptTags;

  bool? _hideReceiptCategories;
  bool? get hideReceiptCategories => _$this._hideReceiptCategories;
  set hideReceiptCategories(covariant bool? hideReceiptCategories) =>
      _$this._hideReceiptCategories = hideReceiptCategories;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(covariant int? groupId) => _$this._groupId = groupId;

  bool? _hideItemCategories;
  bool? get hideItemCategories => _$this._hideItemCategories;
  set hideItemCategories(covariant bool? hideItemCategories) =>
      _$this._hideItemCategories = hideItemCategories;

  bool? _hideImages;
  bool? get hideImages => _$this._hideImages;
  set hideImages(covariant bool? hideImages) => _$this._hideImages = hideImages;

  bool? _hideItemTags;
  bool? get hideItemTags => _$this._hideItemTags;
  set hideItemTags(covariant bool? hideItemTags) =>
      _$this._hideItemTags = hideItemTags;

  bool? _hideComments;
  bool? get hideComments => _$this._hideComments;
  set hideComments(covariant bool? hideComments) =>
      _$this._hideComments = hideComments;

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

  GroupReceiptSettingsBuilder() {
    GroupReceiptSettings._defaults(this);
  }

  GroupReceiptSettingsBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _hideReceiptTags = $v.hideReceiptTags;
      _hideReceiptCategories = $v.hideReceiptCategories;
      _groupId = $v.groupId;
      _hideItemCategories = $v.hideItemCategories;
      _hideImages = $v.hideImages;
      _hideItemTags = $v.hideItemTags;
      _hideComments = $v.hideComments;
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
  void replace(covariant GroupReceiptSettings other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$GroupReceiptSettings;
  }

  @override
  void update(void Function(GroupReceiptSettingsBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GroupReceiptSettings build() => _build();

  _$GroupReceiptSettings _build() {
    final _$result = _$v ??
        new _$GroupReceiptSettings._(
            hideReceiptTags: hideReceiptTags,
            hideReceiptCategories: hideReceiptCategories,
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'GroupReceiptSettings', 'groupId'),
            hideItemCategories: hideItemCategories,
            hideImages: hideImages,
            hideItemTags: hideItemTags,
            hideComments: hideComments,
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'GroupReceiptSettings', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'GroupReceiptSettings', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
