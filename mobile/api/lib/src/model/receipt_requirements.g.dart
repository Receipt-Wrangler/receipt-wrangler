// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_requirements.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptRequirements extends ReceiptRequirements {
  @override
  final bool commentRequired;
  @override
  final bool imageRequired;

  factory _$ReceiptRequirements(
          [void Function(ReceiptRequirementsBuilder)? updates]) =>
      (ReceiptRequirementsBuilder()..update(updates))._build();

  _$ReceiptRequirements._(
      {required this.commentRequired, required this.imageRequired})
      : super._();
  @override
  ReceiptRequirements rebuild(
          void Function(ReceiptRequirementsBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptRequirementsBuilder toBuilder() =>
      ReceiptRequirementsBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptRequirements &&
        commentRequired == other.commentRequired &&
        imageRequired == other.imageRequired;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, commentRequired.hashCode);
    _$hash = $jc(_$hash, imageRequired.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptRequirements')
          ..add('commentRequired', commentRequired)
          ..add('imageRequired', imageRequired))
        .toString();
  }
}

class ReceiptRequirementsBuilder
    implements Builder<ReceiptRequirements, ReceiptRequirementsBuilder> {
  _$ReceiptRequirements? _$v;

  bool? _commentRequired;
  bool? get commentRequired => _$this._commentRequired;
  set commentRequired(bool? commentRequired) =>
      _$this._commentRequired = commentRequired;

  bool? _imageRequired;
  bool? get imageRequired => _$this._imageRequired;
  set imageRequired(bool? imageRequired) =>
      _$this._imageRequired = imageRequired;

  ReceiptRequirementsBuilder() {
    ReceiptRequirements._defaults(this);
  }

  ReceiptRequirementsBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _commentRequired = $v.commentRequired;
      _imageRequired = $v.imageRequired;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptRequirements other) {
    _$v = other as _$ReceiptRequirements;
  }

  @override
  void update(void Function(ReceiptRequirementsBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptRequirements build() => _build();

  _$ReceiptRequirements _build() {
    final _$result = _$v ??
        _$ReceiptRequirements._(
          commentRequired: BuiltValueNullFieldError.checkNotNull(
              commentRequired, r'ReceiptRequirements', 'commentRequired'),
          imageRequired: BuiltValueNullFieldError.checkNotNull(
              imageRequired, r'ReceiptRequirements', 'imageRequired'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
