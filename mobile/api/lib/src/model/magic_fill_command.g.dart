// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'magic_fill_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$MagicFillCommand extends MagicFillCommand {
  @override
  final BuiltList<int> imageData;
  @override
  final String filename;

  factory _$MagicFillCommand(
          [void Function(MagicFillCommandBuilder)? updates]) =>
      (new MagicFillCommandBuilder()..update(updates))._build();

  _$MagicFillCommand._({required this.imageData, required this.filename})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        imageData, r'MagicFillCommand', 'imageData');
    BuiltValueNullFieldError.checkNotNull(
        filename, r'MagicFillCommand', 'filename');
  }

  @override
  MagicFillCommand rebuild(void Function(MagicFillCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  MagicFillCommandBuilder toBuilder() =>
      new MagicFillCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is MagicFillCommand &&
        imageData == other.imageData &&
        filename == other.filename;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, imageData.hashCode);
    _$hash = $jc(_$hash, filename.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'MagicFillCommand')
          ..add('imageData', imageData)
          ..add('filename', filename))
        .toString();
  }
}

class MagicFillCommandBuilder
    implements Builder<MagicFillCommand, MagicFillCommandBuilder> {
  _$MagicFillCommand? _$v;

  ListBuilder<int>? _imageData;
  ListBuilder<int> get imageData =>
      _$this._imageData ??= new ListBuilder<int>();
  set imageData(ListBuilder<int>? imageData) => _$this._imageData = imageData;

  String? _filename;
  String? get filename => _$this._filename;
  set filename(String? filename) => _$this._filename = filename;

  MagicFillCommandBuilder() {
    MagicFillCommand._defaults(this);
  }

  MagicFillCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _imageData = $v.imageData.toBuilder();
      _filename = $v.filename;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(MagicFillCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$MagicFillCommand;
  }

  @override
  void update(void Function(MagicFillCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  MagicFillCommand build() => _build();

  _$MagicFillCommand _build() {
    _$MagicFillCommand _$result;
    try {
      _$result = _$v ??
          new _$MagicFillCommand._(
              imageData: imageData.build(),
              filename: BuiltValueNullFieldError.checkNotNull(
                  filename, r'MagicFillCommand', 'filename'));
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'imageData';
        imageData.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'MagicFillCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
