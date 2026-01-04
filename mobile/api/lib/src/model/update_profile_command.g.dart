// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'update_profile_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpdateProfileCommand extends UpdateProfileCommand {
  @override
  final String displayName;
  @override
  final String defaultAvatarColor;

  factory _$UpdateProfileCommand(
          [void Function(UpdateProfileCommandBuilder)? updates]) =>
      (new UpdateProfileCommandBuilder()..update(updates))._build();

  _$UpdateProfileCommand._(
      {required this.displayName, required this.defaultAvatarColor})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        displayName, r'UpdateProfileCommand', 'displayName');
    BuiltValueNullFieldError.checkNotNull(
        defaultAvatarColor, r'UpdateProfileCommand', 'defaultAvatarColor');
  }

  @override
  UpdateProfileCommand rebuild(
          void Function(UpdateProfileCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpdateProfileCommandBuilder toBuilder() =>
      new UpdateProfileCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpdateProfileCommand &&
        displayName == other.displayName &&
        defaultAvatarColor == other.defaultAvatarColor;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, displayName.hashCode);
    _$hash = $jc(_$hash, defaultAvatarColor.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpdateProfileCommand')
          ..add('displayName', displayName)
          ..add('defaultAvatarColor', defaultAvatarColor))
        .toString();
  }
}

class UpdateProfileCommandBuilder
    implements Builder<UpdateProfileCommand, UpdateProfileCommandBuilder> {
  _$UpdateProfileCommand? _$v;

  String? _displayName;
  String? get displayName => _$this._displayName;
  set displayName(String? displayName) => _$this._displayName = displayName;

  String? _defaultAvatarColor;
  String? get defaultAvatarColor => _$this._defaultAvatarColor;
  set defaultAvatarColor(String? defaultAvatarColor) =>
      _$this._defaultAvatarColor = defaultAvatarColor;

  UpdateProfileCommandBuilder() {
    UpdateProfileCommand._defaults(this);
  }

  UpdateProfileCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _displayName = $v.displayName;
      _defaultAvatarColor = $v.defaultAvatarColor;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpdateProfileCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpdateProfileCommand;
  }

  @override
  void update(void Function(UpdateProfileCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpdateProfileCommand build() => _build();

  _$UpdateProfileCommand _build() {
    final _$result = _$v ??
        new _$UpdateProfileCommand._(
            displayName: BuiltValueNullFieldError.checkNotNull(
                displayName, r'UpdateProfileCommand', 'displayName'),
            defaultAvatarColor: BuiltValueNullFieldError.checkNotNull(
                defaultAvatarColor,
                r'UpdateProfileCommand',
                'defaultAvatarColor'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
