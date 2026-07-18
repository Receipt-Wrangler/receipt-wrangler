// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_preview_response.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReportPreviewResponse extends ReportPreviewResponse {
  @override
  final String html;
  @override
  final int receiptCount;
  @override
  final BuiltList<String>? allowedActions;

  factory _$ReportPreviewResponse(
          [void Function(ReportPreviewResponseBuilder)? updates]) =>
      (ReportPreviewResponseBuilder()..update(updates))._build();

  _$ReportPreviewResponse._(
      {required this.html, required this.receiptCount, this.allowedActions})
      : super._();
  @override
  ReportPreviewResponse rebuild(
          void Function(ReportPreviewResponseBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportPreviewResponseBuilder toBuilder() =>
      ReportPreviewResponseBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportPreviewResponse &&
        html == other.html &&
        receiptCount == other.receiptCount &&
        allowedActions == other.allowedActions;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, html.hashCode);
    _$hash = $jc(_$hash, receiptCount.hashCode);
    _$hash = $jc(_$hash, allowedActions.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportPreviewResponse')
          ..add('html', html)
          ..add('receiptCount', receiptCount)
          ..add('allowedActions', allowedActions))
        .toString();
  }
}

class ReportPreviewResponseBuilder
    implements Builder<ReportPreviewResponse, ReportPreviewResponseBuilder> {
  _$ReportPreviewResponse? _$v;

  String? _html;
  String? get html => _$this._html;
  set html(String? html) => _$this._html = html;

  int? _receiptCount;
  int? get receiptCount => _$this._receiptCount;
  set receiptCount(int? receiptCount) => _$this._receiptCount = receiptCount;

  ListBuilder<String>? _allowedActions;
  ListBuilder<String> get allowedActions =>
      _$this._allowedActions ??= ListBuilder<String>();
  set allowedActions(ListBuilder<String>? allowedActions) =>
      _$this._allowedActions = allowedActions;

  ReportPreviewResponseBuilder() {
    ReportPreviewResponse._defaults(this);
  }

  ReportPreviewResponseBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _html = $v.html;
      _receiptCount = $v.receiptCount;
      _allowedActions = $v.allowedActions?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportPreviewResponse other) {
    _$v = other as _$ReportPreviewResponse;
  }

  @override
  void update(void Function(ReportPreviewResponseBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportPreviewResponse build() => _build();

  _$ReportPreviewResponse _build() {
    _$ReportPreviewResponse _$result;
    try {
      _$result = _$v ??
          _$ReportPreviewResponse._(
            html: BuiltValueNullFieldError.checkNotNull(
                html, r'ReportPreviewResponse', 'html'),
            receiptCount: BuiltValueNullFieldError.checkNotNull(
                receiptCount, r'ReportPreviewResponse', 'receiptCount'),
            allowedActions: _allowedActions?.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'allowedActions';
        _allowedActions?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReportPreviewResponse', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
