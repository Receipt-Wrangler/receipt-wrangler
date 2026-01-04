// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'encoded_image.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$EncodedImage extends EncodedImage {
  @override
  final String encodedImage;

  factory _$EncodedImage([void Function(EncodedImageBuilder)? updates]) =>
      (new EncodedImageBuilder()..update(updates))._build();

  _$EncodedImage._({required this.encodedImage}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        encodedImage, r'EncodedImage', 'encodedImage');
  }

  @override
  EncodedImage rebuild(void Function(EncodedImageBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  EncodedImageBuilder toBuilder() => new EncodedImageBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is EncodedImage && encodedImage == other.encodedImage;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, encodedImage.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'EncodedImage')
          ..add('encodedImage', encodedImage))
        .toString();
  }
}

class EncodedImageBuilder
    implements Builder<EncodedImage, EncodedImageBuilder> {
  _$EncodedImage? _$v;

  String? _encodedImage;
  String? get encodedImage => _$this._encodedImage;
  set encodedImage(String? encodedImage) => _$this._encodedImage = encodedImage;

  EncodedImageBuilder() {
    EncodedImage._defaults(this);
  }

  EncodedImageBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _encodedImage = $v.encodedImage;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(EncodedImage other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$EncodedImage;
  }

  @override
  void update(void Function(EncodedImageBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  EncodedImage build() => _build();

  _$EncodedImage _build() {
    final _$result = _$v ??
        new _$EncodedImage._(
            encodedImage: BuiltValueNullFieldError.checkNotNull(
                encodedImage, r'EncodedImage', 'encodedImage'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
