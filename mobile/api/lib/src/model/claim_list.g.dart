// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim_list.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ClaimList extends ClaimList {
  @override
  final BuiltList<Claim> claims;
  @override
  final ClaimListSummary summary;

  factory _$ClaimList([void Function(ClaimListBuilder)? updates]) =>
      (ClaimListBuilder()..update(updates))._build();

  _$ClaimList._({required this.claims, required this.summary}) : super._();
  @override
  ClaimList rebuild(void Function(ClaimListBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimListBuilder toBuilder() => ClaimListBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ClaimList &&
        claims == other.claims &&
        summary == other.summary;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, claims.hashCode);
    _$hash = $jc(_$hash, summary.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ClaimList')
          ..add('claims', claims)
          ..add('summary', summary))
        .toString();
  }
}

class ClaimListBuilder implements Builder<ClaimList, ClaimListBuilder> {
  _$ClaimList? _$v;

  ListBuilder<Claim>? _claims;
  ListBuilder<Claim> get claims => _$this._claims ??= ListBuilder<Claim>();
  set claims(ListBuilder<Claim>? claims) => _$this._claims = claims;

  ClaimListSummaryBuilder? _summary;
  ClaimListSummaryBuilder get summary =>
      _$this._summary ??= ClaimListSummaryBuilder();
  set summary(ClaimListSummaryBuilder? summary) => _$this._summary = summary;

  ClaimListBuilder() {
    ClaimList._defaults(this);
  }

  ClaimListBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _claims = $v.claims.toBuilder();
      _summary = $v.summary.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ClaimList other) {
    _$v = other as _$ClaimList;
  }

  @override
  void update(void Function(ClaimListBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ClaimList build() => _build();

  _$ClaimList _build() {
    _$ClaimList _$result;
    try {
      _$result = _$v ??
          _$ClaimList._(
            claims: claims.build(),
            summary: summary.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'claims';
        claims.build();
        _$failedField = 'summary';
        summary.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ClaimList', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
