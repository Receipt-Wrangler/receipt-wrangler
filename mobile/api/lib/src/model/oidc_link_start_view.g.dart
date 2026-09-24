// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'oidc_link_start_view.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$OidcLinkStartView extends OidcLinkStartView {
  @override
  final String authorizationUrl;

  factory _$OidcLinkStartView(
          [void Function(OidcLinkStartViewBuilder)? updates]) =>
      (OidcLinkStartViewBuilder()..update(updates))._build();

  _$OidcLinkStartView._({required this.authorizationUrl}) : super._();
  @override
  OidcLinkStartView rebuild(void Function(OidcLinkStartViewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  OidcLinkStartViewBuilder toBuilder() =>
      OidcLinkStartViewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is OidcLinkStartView &&
        authorizationUrl == other.authorizationUrl;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, authorizationUrl.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'OidcLinkStartView')
          ..add('authorizationUrl', authorizationUrl))
        .toString();
  }
}

class OidcLinkStartViewBuilder
    implements Builder<OidcLinkStartView, OidcLinkStartViewBuilder> {
  _$OidcLinkStartView? _$v;

  String? _authorizationUrl;
  String? get authorizationUrl => _$this._authorizationUrl;
  set authorizationUrl(String? authorizationUrl) =>
      _$this._authorizationUrl = authorizationUrl;

  OidcLinkStartViewBuilder() {
    OidcLinkStartView._defaults(this);
  }

  OidcLinkStartViewBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _authorizationUrl = $v.authorizationUrl;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(OidcLinkStartView other) {
    _$v = other as _$OidcLinkStartView;
  }

  @override
  void update(void Function(OidcLinkStartViewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  OidcLinkStartView build() => _build();

  _$OidcLinkStartView _build() {
    final _$result = _$v ??
        _$OidcLinkStartView._(
          authorizationUrl: BuiltValueNullFieldError.checkNotNull(
              authorizationUrl, r'OidcLinkStartView', 'authorizationUrl'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
