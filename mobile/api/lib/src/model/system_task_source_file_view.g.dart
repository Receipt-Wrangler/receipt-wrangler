// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'system_task_source_file_view.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$SystemTaskSourceFileView extends SystemTaskSourceFileView {
  @override
  final String name;
  @override
  final String encodedImage;

  factory _$SystemTaskSourceFileView(
          [void Function(SystemTaskSourceFileViewBuilder)? updates]) =>
      (SystemTaskSourceFileViewBuilder()..update(updates))._build();

  _$SystemTaskSourceFileView._({required this.name, required this.encodedImage})
      : super._();
  @override
  SystemTaskSourceFileView rebuild(
          void Function(SystemTaskSourceFileViewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  SystemTaskSourceFileViewBuilder toBuilder() =>
      SystemTaskSourceFileViewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is SystemTaskSourceFileView &&
        name == other.name &&
        encodedImage == other.encodedImage;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, encodedImage.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'SystemTaskSourceFileView')
          ..add('name', name)
          ..add('encodedImage', encodedImage))
        .toString();
  }
}

class SystemTaskSourceFileViewBuilder
    implements
        Builder<SystemTaskSourceFileView, SystemTaskSourceFileViewBuilder> {
  _$SystemTaskSourceFileView? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _encodedImage;
  String? get encodedImage => _$this._encodedImage;
  set encodedImage(String? encodedImage) => _$this._encodedImage = encodedImage;

  SystemTaskSourceFileViewBuilder() {
    SystemTaskSourceFileView._defaults(this);
  }

  SystemTaskSourceFileViewBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _encodedImage = $v.encodedImage;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(SystemTaskSourceFileView other) {
    _$v = other as _$SystemTaskSourceFileView;
  }

  @override
  void update(void Function(SystemTaskSourceFileViewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  SystemTaskSourceFileView build() => _build();

  _$SystemTaskSourceFileView _build() {
    final _$result = _$v ??
        _$SystemTaskSourceFileView._(
          name: BuiltValueNullFieldError.checkNotNull(
              name, r'SystemTaskSourceFileView', 'name'),
          encodedImage: BuiltValueNullFieldError.checkNotNull(
              encodedImage, r'SystemTaskSourceFileView', 'encodedImage'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
