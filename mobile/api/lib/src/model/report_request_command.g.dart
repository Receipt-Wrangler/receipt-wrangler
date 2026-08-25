// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_request_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ReportRequestCommandFormatsEnum _$reportRequestCommandFormatsEnum_csv =
    const ReportRequestCommandFormatsEnum._('csv');
const ReportRequestCommandFormatsEnum _$reportRequestCommandFormatsEnum_xlsx =
    const ReportRequestCommandFormatsEnum._('xlsx');
const ReportRequestCommandFormatsEnum _$reportRequestCommandFormatsEnum_pdf =
    const ReportRequestCommandFormatsEnum._('pdf');

ReportRequestCommandFormatsEnum _$reportRequestCommandFormatsEnumValueOf(
    String name) {
  switch (name) {
    case 'csv':
      return _$reportRequestCommandFormatsEnum_csv;
    case 'xlsx':
      return _$reportRequestCommandFormatsEnum_xlsx;
    case 'pdf':
      return _$reportRequestCommandFormatsEnum_pdf;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<ReportRequestCommandFormatsEnum>
    _$reportRequestCommandFormatsEnumValues = BuiltSet<
        ReportRequestCommandFormatsEnum>(const <ReportRequestCommandFormatsEnum>[
  _$reportRequestCommandFormatsEnum_csv,
  _$reportRequestCommandFormatsEnum_xlsx,
  _$reportRequestCommandFormatsEnum_pdf,
]);

Serializer<ReportRequestCommandFormatsEnum>
    _$reportRequestCommandFormatsEnumSerializer =
    _$ReportRequestCommandFormatsEnumSerializer();

class _$ReportRequestCommandFormatsEnumSerializer
    implements PrimitiveSerializer<ReportRequestCommandFormatsEnum> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'csv': 'csv',
    'xlsx': 'xlsx',
    'pdf': 'pdf',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'csv': 'csv',
    'xlsx': 'xlsx',
    'pdf': 'pdf',
  };

  @override
  final Iterable<Type> types = const <Type>[ReportRequestCommandFormatsEnum];
  @override
  final String wireName = 'ReportRequestCommandFormatsEnum';

  @override
  Object serialize(
          Serializers serializers, ReportRequestCommandFormatsEnum object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ReportRequestCommandFormatsEnum deserialize(
          Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ReportRequestCommandFormatsEnum.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

class _$ReportRequestCommand extends ReportRequestCommand {
  @override
  final String? name;
  @override
  final BuiltList<String> groupIds;
  @override
  final ReportPeriod period;
  @override
  final ReceiptPagedRequestFilter? filter;
  @override
  final BuiltList<String>? groupBy;
  @override
  final BuiltMap<String, String>? groupByLabels;
  @override
  final ReportDetail detail;
  @override
  final BuiltList<ReportColumn> columns;
  @override
  final bool? subtotals;
  @override
  final bool? grandTotals;
  @override
  final ReportDocument? document;
  @override
  final BuiltList<ReportRequestCommandFormatsEnum> formats;

  factory _$ReportRequestCommand(
          [void Function(ReportRequestCommandBuilder)? updates]) =>
      (ReportRequestCommandBuilder()..update(updates))._build();

  _$ReportRequestCommand._(
      {this.name,
      required this.groupIds,
      required this.period,
      this.filter,
      this.groupBy,
      this.groupByLabels,
      required this.detail,
      required this.columns,
      this.subtotals,
      this.grandTotals,
      this.document,
      required this.formats})
      : super._();
  @override
  ReportRequestCommand rebuild(
          void Function(ReportRequestCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportRequestCommandBuilder toBuilder() =>
      ReportRequestCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportRequestCommand &&
        name == other.name &&
        groupIds == other.groupIds &&
        period == other.period &&
        filter == other.filter &&
        groupBy == other.groupBy &&
        groupByLabels == other.groupByLabels &&
        detail == other.detail &&
        columns == other.columns &&
        subtotals == other.subtotals &&
        grandTotals == other.grandTotals &&
        document == other.document &&
        formats == other.formats;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, groupIds.hashCode);
    _$hash = $jc(_$hash, period.hashCode);
    _$hash = $jc(_$hash, filter.hashCode);
    _$hash = $jc(_$hash, groupBy.hashCode);
    _$hash = $jc(_$hash, groupByLabels.hashCode);
    _$hash = $jc(_$hash, detail.hashCode);
    _$hash = $jc(_$hash, columns.hashCode);
    _$hash = $jc(_$hash, subtotals.hashCode);
    _$hash = $jc(_$hash, grandTotals.hashCode);
    _$hash = $jc(_$hash, document.hashCode);
    _$hash = $jc(_$hash, formats.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportRequestCommand')
          ..add('name', name)
          ..add('groupIds', groupIds)
          ..add('period', period)
          ..add('filter', filter)
          ..add('groupBy', groupBy)
          ..add('groupByLabels', groupByLabels)
          ..add('detail', detail)
          ..add('columns', columns)
          ..add('subtotals', subtotals)
          ..add('grandTotals', grandTotals)
          ..add('document', document)
          ..add('formats', formats))
        .toString();
  }
}

class ReportRequestCommandBuilder
    implements Builder<ReportRequestCommand, ReportRequestCommandBuilder> {
  _$ReportRequestCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  ListBuilder<String>? _groupIds;
  ListBuilder<String> get groupIds =>
      _$this._groupIds ??= ListBuilder<String>();
  set groupIds(ListBuilder<String>? groupIds) => _$this._groupIds = groupIds;

  ReportPeriodBuilder? _period;
  ReportPeriodBuilder get period => _$this._period ??= ReportPeriodBuilder();
  set period(ReportPeriodBuilder? period) => _$this._period = period;

  ReceiptPagedRequestFilterBuilder? _filter;
  ReceiptPagedRequestFilterBuilder get filter =>
      _$this._filter ??= ReceiptPagedRequestFilterBuilder();
  set filter(ReceiptPagedRequestFilterBuilder? filter) =>
      _$this._filter = filter;

  ListBuilder<String>? _groupBy;
  ListBuilder<String> get groupBy => _$this._groupBy ??= ListBuilder<String>();
  set groupBy(ListBuilder<String>? groupBy) => _$this._groupBy = groupBy;

  MapBuilder<String, String>? _groupByLabels;
  MapBuilder<String, String> get groupByLabels =>
      _$this._groupByLabels ??= MapBuilder<String, String>();
  set groupByLabels(MapBuilder<String, String>? groupByLabels) =>
      _$this._groupByLabels = groupByLabels;

  ReportDetailBuilder? _detail;
  ReportDetailBuilder get detail => _$this._detail ??= ReportDetailBuilder();
  set detail(ReportDetailBuilder? detail) => _$this._detail = detail;

  ListBuilder<ReportColumn>? _columns;
  ListBuilder<ReportColumn> get columns =>
      _$this._columns ??= ListBuilder<ReportColumn>();
  set columns(ListBuilder<ReportColumn>? columns) => _$this._columns = columns;

  bool? _subtotals;
  bool? get subtotals => _$this._subtotals;
  set subtotals(bool? subtotals) => _$this._subtotals = subtotals;

  bool? _grandTotals;
  bool? get grandTotals => _$this._grandTotals;
  set grandTotals(bool? grandTotals) => _$this._grandTotals = grandTotals;

  ReportDocumentBuilder? _document;
  ReportDocumentBuilder get document =>
      _$this._document ??= ReportDocumentBuilder();
  set document(ReportDocumentBuilder? document) => _$this._document = document;

  ListBuilder<ReportRequestCommandFormatsEnum>? _formats;
  ListBuilder<ReportRequestCommandFormatsEnum> get formats =>
      _$this._formats ??= ListBuilder<ReportRequestCommandFormatsEnum>();
  set formats(ListBuilder<ReportRequestCommandFormatsEnum>? formats) =>
      _$this._formats = formats;

  ReportRequestCommandBuilder() {
    ReportRequestCommand._defaults(this);
  }

  ReportRequestCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _groupIds = $v.groupIds.toBuilder();
      _period = $v.period.toBuilder();
      _filter = $v.filter?.toBuilder();
      _groupBy = $v.groupBy?.toBuilder();
      _groupByLabels = $v.groupByLabels?.toBuilder();
      _detail = $v.detail.toBuilder();
      _columns = $v.columns.toBuilder();
      _subtotals = $v.subtotals;
      _grandTotals = $v.grandTotals;
      _document = $v.document?.toBuilder();
      _formats = $v.formats.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportRequestCommand other) {
    _$v = other as _$ReportRequestCommand;
  }

  @override
  void update(void Function(ReportRequestCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportRequestCommand build() => _build();

  _$ReportRequestCommand _build() {
    _$ReportRequestCommand _$result;
    try {
      _$result = _$v ??
          _$ReportRequestCommand._(
            name: name,
            groupIds: groupIds.build(),
            period: period.build(),
            filter: _filter?.build(),
            groupBy: _groupBy?.build(),
            groupByLabels: _groupByLabels?.build(),
            detail: detail.build(),
            columns: columns.build(),
            subtotals: subtotals,
            grandTotals: grandTotals,
            document: _document?.build(),
            formats: formats.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'groupIds';
        groupIds.build();
        _$failedField = 'period';
        period.build();
        _$failedField = 'filter';
        _filter?.build();
        _$failedField = 'groupBy';
        _groupBy?.build();
        _$failedField = 'groupByLabels';
        _groupByLabels?.build();
        _$failedField = 'detail';
        detail.build();
        _$failedField = 'columns';
        columns.build();

        _$failedField = 'document';
        _document?.build();
        _$failedField = 'formats';
        formats.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReportRequestCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
