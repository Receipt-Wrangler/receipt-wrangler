// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claims.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Claims extends Claims {
  @override
  final int userId;
  @override
  final UserRole userRole;
  @override
  final String displayName;
  @override
  final String defaultAvatarColor;
  @override
  final String username;
  @override
  final String iss;
  @override
  final String? sub;
  @override
  final BuiltList<String>? aud;
  @override
  final int exp;
  @override
  final int? nbf;
  @override
  final int? iat;
  @override
  final String? jti;

  factory _$Claims([void Function(ClaimsBuilder)? updates]) =>
      (new ClaimsBuilder()..update(updates))._build();

  _$Claims._(
      {required this.userId,
      required this.userRole,
      required this.displayName,
      required this.defaultAvatarColor,
      required this.username,
      required this.iss,
      this.sub,
      this.aud,
      required this.exp,
      this.nbf,
      this.iat,
      this.jti})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(userId, r'Claims', 'userId');
    BuiltValueNullFieldError.checkNotNull(userRole, r'Claims', 'userRole');
    BuiltValueNullFieldError.checkNotNull(
        displayName, r'Claims', 'displayName');
    BuiltValueNullFieldError.checkNotNull(
        defaultAvatarColor, r'Claims', 'defaultAvatarColor');
    BuiltValueNullFieldError.checkNotNull(username, r'Claims', 'username');
    BuiltValueNullFieldError.checkNotNull(iss, r'Claims', 'iss');
    BuiltValueNullFieldError.checkNotNull(exp, r'Claims', 'exp');
  }

  @override
  Claims rebuild(void Function(ClaimsBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimsBuilder toBuilder() => new ClaimsBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Claims &&
        userId == other.userId &&
        userRole == other.userRole &&
        displayName == other.displayName &&
        defaultAvatarColor == other.defaultAvatarColor &&
        username == other.username &&
        iss == other.iss &&
        sub == other.sub &&
        aud == other.aud &&
        exp == other.exp &&
        nbf == other.nbf &&
        iat == other.iat &&
        jti == other.jti;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jc(_$hash, userRole.hashCode);
    _$hash = $jc(_$hash, displayName.hashCode);
    _$hash = $jc(_$hash, defaultAvatarColor.hashCode);
    _$hash = $jc(_$hash, username.hashCode);
    _$hash = $jc(_$hash, iss.hashCode);
    _$hash = $jc(_$hash, sub.hashCode);
    _$hash = $jc(_$hash, aud.hashCode);
    _$hash = $jc(_$hash, exp.hashCode);
    _$hash = $jc(_$hash, nbf.hashCode);
    _$hash = $jc(_$hash, iat.hashCode);
    _$hash = $jc(_$hash, jti.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Claims')
          ..add('userId', userId)
          ..add('userRole', userRole)
          ..add('displayName', displayName)
          ..add('defaultAvatarColor', defaultAvatarColor)
          ..add('username', username)
          ..add('iss', iss)
          ..add('sub', sub)
          ..add('aud', aud)
          ..add('exp', exp)
          ..add('nbf', nbf)
          ..add('iat', iat)
          ..add('jti', jti))
        .toString();
  }
}

class ClaimsBuilder implements Builder<Claims, ClaimsBuilder> {
  _$Claims? _$v;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  UserRole? _userRole;
  UserRole? get userRole => _$this._userRole;
  set userRole(UserRole? userRole) => _$this._userRole = userRole;

  String? _displayName;
  String? get displayName => _$this._displayName;
  set displayName(String? displayName) => _$this._displayName = displayName;

  String? _defaultAvatarColor;
  String? get defaultAvatarColor => _$this._defaultAvatarColor;
  set defaultAvatarColor(String? defaultAvatarColor) =>
      _$this._defaultAvatarColor = defaultAvatarColor;

  String? _username;
  String? get username => _$this._username;
  set username(String? username) => _$this._username = username;

  String? _iss;
  String? get iss => _$this._iss;
  set iss(String? iss) => _$this._iss = iss;

  String? _sub;
  String? get sub => _$this._sub;
  set sub(String? sub) => _$this._sub = sub;

  ListBuilder<String>? _aud;
  ListBuilder<String> get aud => _$this._aud ??= new ListBuilder<String>();
  set aud(ListBuilder<String>? aud) => _$this._aud = aud;

  int? _exp;
  int? get exp => _$this._exp;
  set exp(int? exp) => _$this._exp = exp;

  int? _nbf;
  int? get nbf => _$this._nbf;
  set nbf(int? nbf) => _$this._nbf = nbf;

  int? _iat;
  int? get iat => _$this._iat;
  set iat(int? iat) => _$this._iat = iat;

  String? _jti;
  String? get jti => _$this._jti;
  set jti(String? jti) => _$this._jti = jti;

  ClaimsBuilder() {
    Claims._defaults(this);
  }

  ClaimsBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _userId = $v.userId;
      _userRole = $v.userRole;
      _displayName = $v.displayName;
      _defaultAvatarColor = $v.defaultAvatarColor;
      _username = $v.username;
      _iss = $v.iss;
      _sub = $v.sub;
      _aud = $v.aud?.toBuilder();
      _exp = $v.exp;
      _nbf = $v.nbf;
      _iat = $v.iat;
      _jti = $v.jti;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Claims other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Claims;
  }

  @override
  void update(void Function(ClaimsBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Claims build() => _build();

  _$Claims _build() {
    _$Claims _$result;
    try {
      _$result = _$v ??
          new _$Claims._(
              userId: BuiltValueNullFieldError.checkNotNull(
                  userId, r'Claims', 'userId'),
              userRole: BuiltValueNullFieldError.checkNotNull(
                  userRole, r'Claims', 'userRole'),
              displayName: BuiltValueNullFieldError.checkNotNull(
                  displayName, r'Claims', 'displayName'),
              defaultAvatarColor: BuiltValueNullFieldError.checkNotNull(
                  defaultAvatarColor, r'Claims', 'defaultAvatarColor'),
              username: BuiltValueNullFieldError.checkNotNull(
                  username, r'Claims', 'username'),
              iss: BuiltValueNullFieldError.checkNotNull(iss, r'Claims', 'iss'),
              sub: sub,
              aud: _aud?.build(),
              exp: BuiltValueNullFieldError.checkNotNull(exp, r'Claims', 'exp'),
              nbf: nbf,
              iat: iat,
              jti: jti);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'aud';
        _aud?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'Claims', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
