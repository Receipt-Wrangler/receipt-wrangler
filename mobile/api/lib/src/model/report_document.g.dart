// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_document.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReportDocument extends ReportDocument {
  @override
  final String? title;
  @override
  final String? intro;
  @override
  final String? footer;

  factory _$ReportDocument([void Function(ReportDocumentBuilder)? updates]) =>
      (ReportDocumentBuilder()..update(updates))._build();

  _$ReportDocument._({this.title, this.intro, this.footer}) : super._();
  @override
  ReportDocument rebuild(void Function(ReportDocumentBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportDocumentBuilder toBuilder() => ReportDocumentBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportDocument &&
        title == other.title &&
        intro == other.intro &&
        footer == other.footer;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, title.hashCode);
    _$hash = $jc(_$hash, intro.hashCode);
    _$hash = $jc(_$hash, footer.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportDocument')
          ..add('title', title)
          ..add('intro', intro)
          ..add('footer', footer))
        .toString();
  }
}

class ReportDocumentBuilder
    implements Builder<ReportDocument, ReportDocumentBuilder> {
  _$ReportDocument? _$v;

  String? _title;
  String? get title => _$this._title;
  set title(String? title) => _$this._title = title;

  String? _intro;
  String? get intro => _$this._intro;
  set intro(String? intro) => _$this._intro = intro;

  String? _footer;
  String? get footer => _$this._footer;
  set footer(String? footer) => _$this._footer = footer;

  ReportDocumentBuilder() {
    ReportDocument._defaults(this);
  }

  ReportDocumentBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _title = $v.title;
      _intro = $v.intro;
      _footer = $v.footer;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportDocument other) {
    _$v = other as _$ReportDocument;
  }

  @override
  void update(void Function(ReportDocumentBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportDocument build() => _build();

  _$ReportDocument _build() {
    final _$result = _$v ??
        _$ReportDocument._(
          title: title,
          intro: intro,
          footer: footer,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
