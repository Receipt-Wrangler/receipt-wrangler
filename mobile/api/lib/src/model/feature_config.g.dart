// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'feature_config.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$FeatureConfig extends FeatureConfig {
  @override
  final bool aiPoweredReceipts;
  @override
  final bool enableLocalSignUp;

  factory _$FeatureConfig([void Function(FeatureConfigBuilder)? updates]) =>
      (new FeatureConfigBuilder()..update(updates))._build();

  _$FeatureConfig._(
      {required this.aiPoweredReceipts, required this.enableLocalSignUp})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        aiPoweredReceipts, r'FeatureConfig', 'aiPoweredReceipts');
    BuiltValueNullFieldError.checkNotNull(
        enableLocalSignUp, r'FeatureConfig', 'enableLocalSignUp');
  }

  @override
  FeatureConfig rebuild(void Function(FeatureConfigBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  FeatureConfigBuilder toBuilder() => new FeatureConfigBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is FeatureConfig &&
        aiPoweredReceipts == other.aiPoweredReceipts &&
        enableLocalSignUp == other.enableLocalSignUp;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, aiPoweredReceipts.hashCode);
    _$hash = $jc(_$hash, enableLocalSignUp.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'FeatureConfig')
          ..add('aiPoweredReceipts', aiPoweredReceipts)
          ..add('enableLocalSignUp', enableLocalSignUp))
        .toString();
  }
}

class FeatureConfigBuilder
    implements Builder<FeatureConfig, FeatureConfigBuilder> {
  _$FeatureConfig? _$v;

  bool? _aiPoweredReceipts;
  bool? get aiPoweredReceipts => _$this._aiPoweredReceipts;
  set aiPoweredReceipts(bool? aiPoweredReceipts) =>
      _$this._aiPoweredReceipts = aiPoweredReceipts;

  bool? _enableLocalSignUp;
  bool? get enableLocalSignUp => _$this._enableLocalSignUp;
  set enableLocalSignUp(bool? enableLocalSignUp) =>
      _$this._enableLocalSignUp = enableLocalSignUp;

  FeatureConfigBuilder() {
    FeatureConfig._defaults(this);
  }

  FeatureConfigBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _aiPoweredReceipts = $v.aiPoweredReceipts;
      _enableLocalSignUp = $v.enableLocalSignUp;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(FeatureConfig other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$FeatureConfig;
  }

  @override
  void update(void Function(FeatureConfigBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  FeatureConfig build() => _build();

  _$FeatureConfig _build() {
    final _$result = _$v ??
        new _$FeatureConfig._(
            aiPoweredReceipts: BuiltValueNullFieldError.checkNotNull(
                aiPoweredReceipts, r'FeatureConfig', 'aiPoweredReceipts'),
            enableLocalSignUp: BuiltValueNullFieldError.checkNotNull(
                enableLocalSignUp, r'FeatureConfig', 'enableLocalSignUp'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
