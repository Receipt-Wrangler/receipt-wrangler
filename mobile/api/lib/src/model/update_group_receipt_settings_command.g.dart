// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'update_group_receipt_settings_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpdateGroupReceiptSettingsCommand
    extends UpdateGroupReceiptSettingsCommand {
  @override
  final bool? hideImages;
  @override
  final bool? hideReceiptCategories;
  @override
  final bool? hideReceiptTags;
  @override
  final bool? hideItemCategories;
  @override
  final bool? hideItemTags;
  @override
  final bool? hideComments;
  @override
  final bool? hideShareCategories;
  @override
  final bool? hideShareTags;

  factory _$UpdateGroupReceiptSettingsCommand(
          [void Function(UpdateGroupReceiptSettingsCommandBuilder)? updates]) =>
      (new UpdateGroupReceiptSettingsCommandBuilder()..update(updates))
          ._build();

  _$UpdateGroupReceiptSettingsCommand._(
      {this.hideImages,
      this.hideReceiptCategories,
      this.hideReceiptTags,
      this.hideItemCategories,
      this.hideItemTags,
      this.hideComments,
      this.hideShareCategories,
      this.hideShareTags})
      : super._();

  @override
  UpdateGroupReceiptSettingsCommand rebuild(
          void Function(UpdateGroupReceiptSettingsCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpdateGroupReceiptSettingsCommandBuilder toBuilder() =>
      new UpdateGroupReceiptSettingsCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpdateGroupReceiptSettingsCommand &&
        hideImages == other.hideImages &&
        hideReceiptCategories == other.hideReceiptCategories &&
        hideReceiptTags == other.hideReceiptTags &&
        hideItemCategories == other.hideItemCategories &&
        hideItemTags == other.hideItemTags &&
        hideComments == other.hideComments &&
        hideShareCategories == other.hideShareCategories &&
        hideShareTags == other.hideShareTags;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, hideImages.hashCode);
    _$hash = $jc(_$hash, hideReceiptCategories.hashCode);
    _$hash = $jc(_$hash, hideReceiptTags.hashCode);
    _$hash = $jc(_$hash, hideItemCategories.hashCode);
    _$hash = $jc(_$hash, hideItemTags.hashCode);
    _$hash = $jc(_$hash, hideComments.hashCode);
    _$hash = $jc(_$hash, hideShareCategories.hashCode);
    _$hash = $jc(_$hash, hideShareTags.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpdateGroupReceiptSettingsCommand')
          ..add('hideImages', hideImages)
          ..add('hideReceiptCategories', hideReceiptCategories)
          ..add('hideReceiptTags', hideReceiptTags)
          ..add('hideItemCategories', hideItemCategories)
          ..add('hideItemTags', hideItemTags)
          ..add('hideComments', hideComments)
          ..add('hideShareCategories', hideShareCategories)
          ..add('hideShareTags', hideShareTags))
        .toString();
  }
}

class UpdateGroupReceiptSettingsCommandBuilder
    implements
        Builder<UpdateGroupReceiptSettingsCommand,
            UpdateGroupReceiptSettingsCommandBuilder> {
  _$UpdateGroupReceiptSettingsCommand? _$v;

  bool? _hideImages;
  bool? get hideImages => _$this._hideImages;
  set hideImages(bool? hideImages) => _$this._hideImages = hideImages;

  bool? _hideReceiptCategories;
  bool? get hideReceiptCategories => _$this._hideReceiptCategories;
  set hideReceiptCategories(bool? hideReceiptCategories) =>
      _$this._hideReceiptCategories = hideReceiptCategories;

  bool? _hideReceiptTags;
  bool? get hideReceiptTags => _$this._hideReceiptTags;
  set hideReceiptTags(bool? hideReceiptTags) =>
      _$this._hideReceiptTags = hideReceiptTags;

  bool? _hideItemCategories;
  bool? get hideItemCategories => _$this._hideItemCategories;
  set hideItemCategories(bool? hideItemCategories) =>
      _$this._hideItemCategories = hideItemCategories;

  bool? _hideItemTags;
  bool? get hideItemTags => _$this._hideItemTags;
  set hideItemTags(bool? hideItemTags) => _$this._hideItemTags = hideItemTags;

  bool? _hideComments;
  bool? get hideComments => _$this._hideComments;
  set hideComments(bool? hideComments) => _$this._hideComments = hideComments;

  bool? _hideShareCategories;
  bool? get hideShareCategories => _$this._hideShareCategories;
  set hideShareCategories(bool? hideShareCategories) =>
      _$this._hideShareCategories = hideShareCategories;

  bool? _hideShareTags;
  bool? get hideShareTags => _$this._hideShareTags;
  set hideShareTags(bool? hideShareTags) =>
      _$this._hideShareTags = hideShareTags;

  UpdateGroupReceiptSettingsCommandBuilder() {
    UpdateGroupReceiptSettingsCommand._defaults(this);
  }

  UpdateGroupReceiptSettingsCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _hideImages = $v.hideImages;
      _hideReceiptCategories = $v.hideReceiptCategories;
      _hideReceiptTags = $v.hideReceiptTags;
      _hideItemCategories = $v.hideItemCategories;
      _hideItemTags = $v.hideItemTags;
      _hideComments = $v.hideComments;
      _hideShareCategories = $v.hideShareCategories;
      _hideShareTags = $v.hideShareTags;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpdateGroupReceiptSettingsCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpdateGroupReceiptSettingsCommand;
  }

  @override
  void update(
      void Function(UpdateGroupReceiptSettingsCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpdateGroupReceiptSettingsCommand build() => _build();

  _$UpdateGroupReceiptSettingsCommand _build() {
    final _$result = _$v ??
        new _$UpdateGroupReceiptSettingsCommand._(
            hideImages: hideImages,
            hideReceiptCategories: hideReceiptCategories,
            hideReceiptTags: hideReceiptTags,
            hideItemCategories: hideItemCategories,
            hideItemTags: hideItemTags,
            hideComments: hideComments,
            hideShareCategories: hideShareCategories,
            hideShareTags: hideShareTags);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
