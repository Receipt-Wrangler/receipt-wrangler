// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_preferences.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UserPreferences extends UserPreferences {
  @override
  final ReceiptStatus? quickScanDefaultStatus;
  @override
  final BuiltList<UserShortcut>? userShortcuts;
  @override
  final bool? showLargeImagePreviews;
  @override
  final int userId;
  @override
  final int? quickScanDefaultGroupId;
  @override
  final int? quickScanDefaultPaidById;
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

  factory _$UserPreferences([void Function(UserPreferencesBuilder)? updates]) =>
      (new UserPreferencesBuilder()..update(updates))._build();

  _$UserPreferences._(
      {this.quickScanDefaultStatus,
      this.userShortcuts,
      this.showLargeImagePreviews,
      required this.userId,
      this.quickScanDefaultGroupId,
      this.quickScanDefaultPaidById,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(userId, r'UserPreferences', 'userId');
    BuiltValueNullFieldError.checkNotNull(id, r'UserPreferences', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'UserPreferences', 'createdAt');
  }

  @override
  UserPreferences rebuild(void Function(UserPreferencesBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UserPreferencesBuilder toBuilder() =>
      new UserPreferencesBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UserPreferences &&
        quickScanDefaultStatus == other.quickScanDefaultStatus &&
        userShortcuts == other.userShortcuts &&
        showLargeImagePreviews == other.showLargeImagePreviews &&
        userId == other.userId &&
        quickScanDefaultGroupId == other.quickScanDefaultGroupId &&
        quickScanDefaultPaidById == other.quickScanDefaultPaidById &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, quickScanDefaultStatus.hashCode);
    _$hash = $jc(_$hash, userShortcuts.hashCode);
    _$hash = $jc(_$hash, showLargeImagePreviews.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jc(_$hash, quickScanDefaultGroupId.hashCode);
    _$hash = $jc(_$hash, quickScanDefaultPaidById.hashCode);
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
    return (newBuiltValueToStringHelper(r'UserPreferences')
          ..add('quickScanDefaultStatus', quickScanDefaultStatus)
          ..add('userShortcuts', userShortcuts)
          ..add('showLargeImagePreviews', showLargeImagePreviews)
          ..add('userId', userId)
          ..add('quickScanDefaultGroupId', quickScanDefaultGroupId)
          ..add('quickScanDefaultPaidById', quickScanDefaultPaidById)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class UserPreferencesBuilder
    implements
        Builder<UserPreferences, UserPreferencesBuilder>,
        BaseModelBuilder {
  _$UserPreferences? _$v;

  ReceiptStatus? _quickScanDefaultStatus;
  ReceiptStatus? get quickScanDefaultStatus => _$this._quickScanDefaultStatus;
  set quickScanDefaultStatus(covariant ReceiptStatus? quickScanDefaultStatus) =>
      _$this._quickScanDefaultStatus = quickScanDefaultStatus;

  ListBuilder<UserShortcut>? _userShortcuts;
  ListBuilder<UserShortcut> get userShortcuts =>
      _$this._userShortcuts ??= new ListBuilder<UserShortcut>();
  set userShortcuts(covariant ListBuilder<UserShortcut>? userShortcuts) =>
      _$this._userShortcuts = userShortcuts;

  bool? _showLargeImagePreviews;
  bool? get showLargeImagePreviews => _$this._showLargeImagePreviews;
  set showLargeImagePreviews(covariant bool? showLargeImagePreviews) =>
      _$this._showLargeImagePreviews = showLargeImagePreviews;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(covariant int? userId) => _$this._userId = userId;

  int? _quickScanDefaultGroupId;
  int? get quickScanDefaultGroupId => _$this._quickScanDefaultGroupId;
  set quickScanDefaultGroupId(covariant int? quickScanDefaultGroupId) =>
      _$this._quickScanDefaultGroupId = quickScanDefaultGroupId;

  int? _quickScanDefaultPaidById;
  int? get quickScanDefaultPaidById => _$this._quickScanDefaultPaidById;
  set quickScanDefaultPaidById(covariant int? quickScanDefaultPaidById) =>
      _$this._quickScanDefaultPaidById = quickScanDefaultPaidById;

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

  UserPreferencesBuilder() {
    UserPreferences._defaults(this);
  }

  UserPreferencesBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _quickScanDefaultStatus = $v.quickScanDefaultStatus;
      _userShortcuts = $v.userShortcuts?.toBuilder();
      _showLargeImagePreviews = $v.showLargeImagePreviews;
      _userId = $v.userId;
      _quickScanDefaultGroupId = $v.quickScanDefaultGroupId;
      _quickScanDefaultPaidById = $v.quickScanDefaultPaidById;
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
  void replace(covariant UserPreferences other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UserPreferences;
  }

  @override
  void update(void Function(UserPreferencesBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UserPreferences build() => _build();

  _$UserPreferences _build() {
    _$UserPreferences _$result;
    try {
      _$result = _$v ??
          new _$UserPreferences._(
              quickScanDefaultStatus: quickScanDefaultStatus,
              userShortcuts: _userShortcuts?.build(),
              showLargeImagePreviews: showLargeImagePreviews,
              userId: BuiltValueNullFieldError.checkNotNull(
                  userId, r'UserPreferences', 'userId'),
              quickScanDefaultGroupId: quickScanDefaultGroupId,
              quickScanDefaultPaidById: quickScanDefaultPaidById,
              id: BuiltValueNullFieldError.checkNotNull(
                  id, r'UserPreferences', 'id'),
              createdAt: BuiltValueNullFieldError.checkNotNull(
                  createdAt, r'UserPreferences', 'createdAt'),
              createdBy: createdBy,
              createdByString: createdByString,
              updatedAt: updatedAt);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'userShortcuts';
        _userShortcuts?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'UserPreferences', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
