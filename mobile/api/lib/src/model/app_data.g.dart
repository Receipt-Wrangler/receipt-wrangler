// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'app_data.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$AppData extends AppData {
  @override
  final About about;
  @override
  final Claims claims;
  @override
  final BuiltList<Group> groups;
  @override
  final BuiltList<UserView> users;
  @override
  final UserPreferences userPreferences;
  @override
  final FeatureConfig featureConfig;
  @override
  final BuiltList<Category> categories;
  @override
  final BuiltList<Tag> tags;
  @override
  final String? jwt;
  @override
  final String? refreshToken;
  @override
  final String currencyDisplay;
  @override
  final CurrencySeparator? currencyThousandthsSeparator;
  @override
  final CurrencySeparator? currencyDecimalSeparator;
  @override
  final CurrencySymbolPosition? currencySymbolPosition;
  @override
  final bool? currencyHideDecimalPlaces;
  @override
  final BuiltList<Icon> icons;
  @override
  final BuiltList<Permission> appPermissions;
  @override
  final BuiltMap<String, BuiltList<Permission>> groupPermissions;
  @override
  final BuiltMap<String, BuiltList<Category>>? groupCategories;
  @override
  final BuiltMap<String, BuiltList<Tag>>? groupTags;

  factory _$AppData([void Function(AppDataBuilder)? updates]) =>
      (AppDataBuilder()..update(updates))._build();

  _$AppData._(
      {required this.about,
      required this.claims,
      required this.groups,
      required this.users,
      required this.userPreferences,
      required this.featureConfig,
      required this.categories,
      required this.tags,
      this.jwt,
      this.refreshToken,
      required this.currencyDisplay,
      this.currencyThousandthsSeparator,
      this.currencyDecimalSeparator,
      this.currencySymbolPosition,
      this.currencyHideDecimalPlaces,
      required this.icons,
      required this.appPermissions,
      required this.groupPermissions,
      this.groupCategories,
      this.groupTags})
      : super._();
  @override
  AppData rebuild(void Function(AppDataBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  AppDataBuilder toBuilder() => AppDataBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is AppData &&
        about == other.about &&
        claims == other.claims &&
        groups == other.groups &&
        users == other.users &&
        userPreferences == other.userPreferences &&
        featureConfig == other.featureConfig &&
        categories == other.categories &&
        tags == other.tags &&
        jwt == other.jwt &&
        refreshToken == other.refreshToken &&
        currencyDisplay == other.currencyDisplay &&
        currencyThousandthsSeparator == other.currencyThousandthsSeparator &&
        currencyDecimalSeparator == other.currencyDecimalSeparator &&
        currencySymbolPosition == other.currencySymbolPosition &&
        currencyHideDecimalPlaces == other.currencyHideDecimalPlaces &&
        icons == other.icons &&
        appPermissions == other.appPermissions &&
        groupPermissions == other.groupPermissions &&
        groupCategories == other.groupCategories &&
        groupTags == other.groupTags;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, about.hashCode);
    _$hash = $jc(_$hash, claims.hashCode);
    _$hash = $jc(_$hash, groups.hashCode);
    _$hash = $jc(_$hash, users.hashCode);
    _$hash = $jc(_$hash, userPreferences.hashCode);
    _$hash = $jc(_$hash, featureConfig.hashCode);
    _$hash = $jc(_$hash, categories.hashCode);
    _$hash = $jc(_$hash, tags.hashCode);
    _$hash = $jc(_$hash, jwt.hashCode);
    _$hash = $jc(_$hash, refreshToken.hashCode);
    _$hash = $jc(_$hash, currencyDisplay.hashCode);
    _$hash = $jc(_$hash, currencyThousandthsSeparator.hashCode);
    _$hash = $jc(_$hash, currencyDecimalSeparator.hashCode);
    _$hash = $jc(_$hash, currencySymbolPosition.hashCode);
    _$hash = $jc(_$hash, currencyHideDecimalPlaces.hashCode);
    _$hash = $jc(_$hash, icons.hashCode);
    _$hash = $jc(_$hash, appPermissions.hashCode);
    _$hash = $jc(_$hash, groupPermissions.hashCode);
    _$hash = $jc(_$hash, groupCategories.hashCode);
    _$hash = $jc(_$hash, groupTags.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'AppData')
          ..add('about', about)
          ..add('claims', claims)
          ..add('groups', groups)
          ..add('users', users)
          ..add('userPreferences', userPreferences)
          ..add('featureConfig', featureConfig)
          ..add('categories', categories)
          ..add('tags', tags)
          ..add('jwt', jwt)
          ..add('refreshToken', refreshToken)
          ..add('currencyDisplay', currencyDisplay)
          ..add('currencyThousandthsSeparator', currencyThousandthsSeparator)
          ..add('currencyDecimalSeparator', currencyDecimalSeparator)
          ..add('currencySymbolPosition', currencySymbolPosition)
          ..add('currencyHideDecimalPlaces', currencyHideDecimalPlaces)
          ..add('icons', icons)
          ..add('appPermissions', appPermissions)
          ..add('groupPermissions', groupPermissions)
          ..add('groupCategories', groupCategories)
          ..add('groupTags', groupTags))
        .toString();
  }
}

class AppDataBuilder implements Builder<AppData, AppDataBuilder> {
  _$AppData? _$v;

  AboutBuilder? _about;
  AboutBuilder get about => _$this._about ??= AboutBuilder();
  set about(AboutBuilder? about) => _$this._about = about;

  ClaimsBuilder? _claims;
  ClaimsBuilder get claims => _$this._claims ??= ClaimsBuilder();
  set claims(ClaimsBuilder? claims) => _$this._claims = claims;

  ListBuilder<Group>? _groups;
  ListBuilder<Group> get groups => _$this._groups ??= ListBuilder<Group>();
  set groups(ListBuilder<Group>? groups) => _$this._groups = groups;

  ListBuilder<UserView>? _users;
  ListBuilder<UserView> get users => _$this._users ??= ListBuilder<UserView>();
  set users(ListBuilder<UserView>? users) => _$this._users = users;

  UserPreferencesBuilder? _userPreferences;
  UserPreferencesBuilder get userPreferences =>
      _$this._userPreferences ??= UserPreferencesBuilder();
  set userPreferences(UserPreferencesBuilder? userPreferences) =>
      _$this._userPreferences = userPreferences;

  FeatureConfigBuilder? _featureConfig;
  FeatureConfigBuilder get featureConfig =>
      _$this._featureConfig ??= FeatureConfigBuilder();
  set featureConfig(FeatureConfigBuilder? featureConfig) =>
      _$this._featureConfig = featureConfig;

  ListBuilder<Category>? _categories;
  ListBuilder<Category> get categories =>
      _$this._categories ??= ListBuilder<Category>();
  set categories(ListBuilder<Category>? categories) =>
      _$this._categories = categories;

  ListBuilder<Tag>? _tags;
  ListBuilder<Tag> get tags => _$this._tags ??= ListBuilder<Tag>();
  set tags(ListBuilder<Tag>? tags) => _$this._tags = tags;

  String? _jwt;
  String? get jwt => _$this._jwt;
  set jwt(String? jwt) => _$this._jwt = jwt;

  String? _refreshToken;
  String? get refreshToken => _$this._refreshToken;
  set refreshToken(String? refreshToken) => _$this._refreshToken = refreshToken;

  String? _currencyDisplay;
  String? get currencyDisplay => _$this._currencyDisplay;
  set currencyDisplay(String? currencyDisplay) =>
      _$this._currencyDisplay = currencyDisplay;

  CurrencySeparator? _currencyThousandthsSeparator;
  CurrencySeparator? get currencyThousandthsSeparator =>
      _$this._currencyThousandthsSeparator;
  set currencyThousandthsSeparator(
          CurrencySeparator? currencyThousandthsSeparator) =>
      _$this._currencyThousandthsSeparator = currencyThousandthsSeparator;

  CurrencySeparator? _currencyDecimalSeparator;
  CurrencySeparator? get currencyDecimalSeparator =>
      _$this._currencyDecimalSeparator;
  set currencyDecimalSeparator(CurrencySeparator? currencyDecimalSeparator) =>
      _$this._currencyDecimalSeparator = currencyDecimalSeparator;

  CurrencySymbolPosition? _currencySymbolPosition;
  CurrencySymbolPosition? get currencySymbolPosition =>
      _$this._currencySymbolPosition;
  set currencySymbolPosition(CurrencySymbolPosition? currencySymbolPosition) =>
      _$this._currencySymbolPosition = currencySymbolPosition;

  bool? _currencyHideDecimalPlaces;
  bool? get currencyHideDecimalPlaces => _$this._currencyHideDecimalPlaces;
  set currencyHideDecimalPlaces(bool? currencyHideDecimalPlaces) =>
      _$this._currencyHideDecimalPlaces = currencyHideDecimalPlaces;

  ListBuilder<Icon>? _icons;
  ListBuilder<Icon> get icons => _$this._icons ??= ListBuilder<Icon>();
  set icons(ListBuilder<Icon>? icons) => _$this._icons = icons;

  ListBuilder<Permission>? _appPermissions;
  ListBuilder<Permission> get appPermissions =>
      _$this._appPermissions ??= ListBuilder<Permission>();
  set appPermissions(ListBuilder<Permission>? appPermissions) =>
      _$this._appPermissions = appPermissions;

  MapBuilder<String, BuiltList<Permission>>? _groupPermissions;
  MapBuilder<String, BuiltList<Permission>> get groupPermissions =>
      _$this._groupPermissions ??= MapBuilder<String, BuiltList<Permission>>();
  set groupPermissions(
          MapBuilder<String, BuiltList<Permission>>? groupPermissions) =>
      _$this._groupPermissions = groupPermissions;

  MapBuilder<String, BuiltList<Category>>? _groupCategories;
  MapBuilder<String, BuiltList<Category>> get groupCategories =>
      _$this._groupCategories ??= MapBuilder<String, BuiltList<Category>>();
  set groupCategories(
          MapBuilder<String, BuiltList<Category>>? groupCategories) =>
      _$this._groupCategories = groupCategories;

  MapBuilder<String, BuiltList<Tag>>? _groupTags;
  MapBuilder<String, BuiltList<Tag>> get groupTags =>
      _$this._groupTags ??= MapBuilder<String, BuiltList<Tag>>();
  set groupTags(MapBuilder<String, BuiltList<Tag>>? groupTags) =>
      _$this._groupTags = groupTags;

  AppDataBuilder() {
    AppData._defaults(this);
  }

  AppDataBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _about = $v.about.toBuilder();
      _claims = $v.claims.toBuilder();
      _groups = $v.groups.toBuilder();
      _users = $v.users.toBuilder();
      _userPreferences = $v.userPreferences.toBuilder();
      _featureConfig = $v.featureConfig.toBuilder();
      _categories = $v.categories.toBuilder();
      _tags = $v.tags.toBuilder();
      _jwt = $v.jwt;
      _refreshToken = $v.refreshToken;
      _currencyDisplay = $v.currencyDisplay;
      _currencyThousandthsSeparator = $v.currencyThousandthsSeparator;
      _currencyDecimalSeparator = $v.currencyDecimalSeparator;
      _currencySymbolPosition = $v.currencySymbolPosition;
      _currencyHideDecimalPlaces = $v.currencyHideDecimalPlaces;
      _icons = $v.icons.toBuilder();
      _appPermissions = $v.appPermissions.toBuilder();
      _groupPermissions = $v.groupPermissions.toBuilder();
      _groupCategories = $v.groupCategories?.toBuilder();
      _groupTags = $v.groupTags?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(AppData other) {
    _$v = other as _$AppData;
  }

  @override
  void update(void Function(AppDataBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  AppData build() => _build();

  _$AppData _build() {
    _$AppData _$result;
    try {
      _$result = _$v ??
          _$AppData._(
            about: about.build(),
            claims: claims.build(),
            groups: groups.build(),
            users: users.build(),
            userPreferences: userPreferences.build(),
            featureConfig: featureConfig.build(),
            categories: categories.build(),
            tags: tags.build(),
            jwt: jwt,
            refreshToken: refreshToken,
            currencyDisplay: BuiltValueNullFieldError.checkNotNull(
                currencyDisplay, r'AppData', 'currencyDisplay'),
            currencyThousandthsSeparator: currencyThousandthsSeparator,
            currencyDecimalSeparator: currencyDecimalSeparator,
            currencySymbolPosition: currencySymbolPosition,
            currencyHideDecimalPlaces: currencyHideDecimalPlaces,
            icons: icons.build(),
            appPermissions: appPermissions.build(),
            groupPermissions: groupPermissions.build(),
            groupCategories: _groupCategories?.build(),
            groupTags: _groupTags?.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'about';
        about.build();
        _$failedField = 'claims';
        claims.build();
        _$failedField = 'groups';
        groups.build();
        _$failedField = 'users';
        users.build();
        _$failedField = 'userPreferences';
        userPreferences.build();
        _$failedField = 'featureConfig';
        featureConfig.build();
        _$failedField = 'categories';
        categories.build();
        _$failedField = 'tags';
        tags.build();

        _$failedField = 'icons';
        icons.build();
        _$failedField = 'appPermissions';
        appPermissions.build();
        _$failedField = 'groupPermissions';
        groupPermissions.build();
        _$failedField = 'groupCategories';
        _groupCategories?.build();
        _$failedField = 'groupTags';
        _groupTags?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'AppData', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
