// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'report_detail.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const ReportDetailModeEnum _$reportDetailModeEnum_records =
    const ReportDetailModeEnum._('records');
const ReportDetailModeEnum _$reportDetailModeEnum_aggregate =
    const ReportDetailModeEnum._('aggregate');

ReportDetailModeEnum _$reportDetailModeEnumValueOf(String name) {
  switch (name) {
    case 'records':
      return _$reportDetailModeEnum_records;
    case 'aggregate':
      return _$reportDetailModeEnum_aggregate;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<ReportDetailModeEnum> _$reportDetailModeEnumValues =
    BuiltSet<ReportDetailModeEnum>(const <ReportDetailModeEnum>[
  _$reportDetailModeEnum_records,
  _$reportDetailModeEnum_aggregate,
]);

Serializer<ReportDetailModeEnum> _$reportDetailModeEnumSerializer =
    _$ReportDetailModeEnumSerializer();

class _$ReportDetailModeEnumSerializer
    implements PrimitiveSerializer<ReportDetailModeEnum> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'records': 'records',
    'aggregate': 'aggregate',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'records': 'records',
    'aggregate': 'aggregate',
  };

  @override
  final Iterable<Type> types = const <Type>[ReportDetailModeEnum];
  @override
  final String wireName = 'ReportDetailModeEnum';

  @override
  Object serialize(Serializers serializers, ReportDetailModeEnum object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  ReportDetailModeEnum deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      ReportDetailModeEnum.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

class _$ReportDetail extends ReportDetail {
  @override
  final ReportDetailModeEnum mode;
  @override
  final String? by;

  factory _$ReportDetail([void Function(ReportDetailBuilder)? updates]) =>
      (ReportDetailBuilder()..update(updates))._build();

  _$ReportDetail._({required this.mode, this.by}) : super._();
  @override
  ReportDetail rebuild(void Function(ReportDetailBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReportDetailBuilder toBuilder() => ReportDetailBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReportDetail && mode == other.mode && by == other.by;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, mode.hashCode);
    _$hash = $jc(_$hash, by.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReportDetail')
          ..add('mode', mode)
          ..add('by', by))
        .toString();
  }
}

class ReportDetailBuilder
    implements Builder<ReportDetail, ReportDetailBuilder> {
  _$ReportDetail? _$v;

  ReportDetailModeEnum? _mode;
  ReportDetailModeEnum? get mode => _$this._mode;
  set mode(ReportDetailModeEnum? mode) => _$this._mode = mode;

  String? _by;
  String? get by => _$this._by;
  set by(String? by) => _$this._by = by;

  ReportDetailBuilder() {
    ReportDetail._defaults(this);
  }

  ReportDetailBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _mode = $v.mode;
      _by = $v.by;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReportDetail other) {
    _$v = other as _$ReportDetail;
  }

  @override
  void update(void Function(ReportDetailBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReportDetail build() => _build();

  _$ReportDetail _build() {
    final _$result = _$v ??
        _$ReportDetail._(
          mode: BuiltValueNullFieldError.checkNotNull(
              mode, r'ReportDetail', 'mode'),
          by: by,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
