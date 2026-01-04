// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'about.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$About extends About {
  @override
  final String buildDate;
  @override
  final String version;

  factory _$About([void Function(AboutBuilder)? updates]) =>
      (new AboutBuilder()..update(updates))._build();

  _$About._({required this.buildDate, required this.version}) : super._() {
    BuiltValueNullFieldError.checkNotNull(buildDate, r'About', 'buildDate');
    BuiltValueNullFieldError.checkNotNull(version, r'About', 'version');
  }

  @override
  About rebuild(void Function(AboutBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  AboutBuilder toBuilder() => new AboutBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is About &&
        buildDate == other.buildDate &&
        version == other.version;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, buildDate.hashCode);
    _$hash = $jc(_$hash, version.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'About')
          ..add('buildDate', buildDate)
          ..add('version', version))
        .toString();
  }
}

class AboutBuilder implements Builder<About, AboutBuilder> {
  _$About? _$v;

  String? _buildDate;
  String? get buildDate => _$this._buildDate;
  set buildDate(String? buildDate) => _$this._buildDate = buildDate;

  String? _version;
  String? get version => _$this._version;
  set version(String? version) => _$this._version = version;

  AboutBuilder() {
    About._defaults(this);
  }

  AboutBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _buildDate = $v.buildDate;
      _version = $v.version;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(About other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$About;
  }

  @override
  void update(void Function(AboutBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  About build() => _build();

  _$About _build() {
    final _$result = _$v ??
        new _$About._(
            buildDate: BuiltValueNullFieldError.checkNotNull(
                buildDate, r'About', 'buildDate'),
            version: BuiltValueNullFieldError.checkNotNull(
                version, r'About', 'version'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
