//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_paged_request_filter.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/report_period.dart';
import 'package:openapi/src/model/report_document.dart';
import 'package:openapi/src/model/report_column.dart';
import 'package:openapi/src/model/report_detail.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_request_command.g.dart';

/// ReportRequestCommand
///
/// Properties:
/// * [name] - Report name; used for the download filename
/// * [groupIds] - Ids of the groups the report covers
/// * [period] 
/// * [filter] - Which receipts go into the report
/// * [groupBy] - Ordered engine field keys to nest the report by
/// * [groupByLabels] - Column-heading overrides for the grouping levels, keyed by the groupBy field key. A key that is absent, blank, or not present in groupBy falls back to the field catalog's label.
/// * [detail] 
/// * [columns] 
/// * [subtotals] - Emit a subtotal row at each grouping level
/// * [grandTotals] - Emit one grand-total row across everything
/// * [document] 
/// * [formats] - One or more output formats; multiple are zipped together
@BuiltValue()
abstract class ReportRequestCommand implements Built<ReportRequestCommand, ReportRequestCommandBuilder> {
  /// Report name; used for the download filename
  @BuiltValueField(wireName: r'name')
  String? get name;

  /// Ids of the groups the report covers
  @BuiltValueField(wireName: r'groupIds')
  BuiltList<String> get groupIds;

  @BuiltValueField(wireName: r'period')
  ReportPeriod get period;

  /// Which receipts go into the report
  @BuiltValueField(wireName: r'filter')
  ReceiptPagedRequestFilter? get filter;

  /// Ordered engine field keys to nest the report by
  @BuiltValueField(wireName: r'groupBy')
  BuiltList<String>? get groupBy;

  /// Column-heading overrides for the grouping levels, keyed by the groupBy field key. A key that is absent, blank, or not present in groupBy falls back to the field catalog's label.
  @BuiltValueField(wireName: r'groupByLabels')
  BuiltMap<String, String>? get groupByLabels;

  @BuiltValueField(wireName: r'detail')
  ReportDetail get detail;

  @BuiltValueField(wireName: r'columns')
  BuiltList<ReportColumn> get columns;

  /// Emit a subtotal row at each grouping level
  @BuiltValueField(wireName: r'subtotals')
  bool? get subtotals;

  /// Emit one grand-total row across everything
  @BuiltValueField(wireName: r'grandTotals')
  bool? get grandTotals;

  @BuiltValueField(wireName: r'document')
  ReportDocument? get document;

  /// One or more output formats; multiple are zipped together
  @BuiltValueField(wireName: r'formats')
  BuiltList<ReportRequestCommandFormatsEnum> get formats;
  // enum formatsEnum {  csv,  xlsx,  pdf,  };

  ReportRequestCommand._();

  factory ReportRequestCommand([void updates(ReportRequestCommandBuilder b)]) = _$ReportRequestCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportRequestCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportRequestCommand> get serializer => _$ReportRequestCommandSerializer();
}

class _$ReportRequestCommandSerializer implements PrimitiveSerializer<ReportRequestCommand> {
  @override
  final Iterable<Type> types = const [ReportRequestCommand, _$ReportRequestCommand];

  @override
  final String wireName = r'ReportRequestCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
        specifiedType: const FullType(String),
      );
    }
    yield r'groupIds';
    yield serializers.serialize(
      object.groupIds,
      specifiedType: const FullType(BuiltList, [FullType(String)]),
    );
    yield r'period';
    yield serializers.serialize(
      object.period,
      specifiedType: const FullType(ReportPeriod),
    );
    if (object.filter != null) {
      yield r'filter';
      yield serializers.serialize(
        object.filter,
        specifiedType: const FullType(ReceiptPagedRequestFilter),
      );
    }
    if (object.groupBy != null) {
      yield r'groupBy';
      yield serializers.serialize(
        object.groupBy,
        specifiedType: const FullType(BuiltList, [FullType(String)]),
      );
    }
    if (object.groupByLabels != null) {
      yield r'groupByLabels';
      yield serializers.serialize(
        object.groupByLabels,
        specifiedType: const FullType(BuiltMap, [FullType(String), FullType(String)]),
      );
    }
    yield r'detail';
    yield serializers.serialize(
      object.detail,
      specifiedType: const FullType(ReportDetail),
    );
    yield r'columns';
    yield serializers.serialize(
      object.columns,
      specifiedType: const FullType(BuiltList, [FullType(ReportColumn)]),
    );
    if (object.subtotals != null) {
      yield r'subtotals';
      yield serializers.serialize(
        object.subtotals,
        specifiedType: const FullType(bool),
      );
    }
    if (object.grandTotals != null) {
      yield r'grandTotals';
      yield serializers.serialize(
        object.grandTotals,
        specifiedType: const FullType(bool),
      );
    }
    if (object.document != null) {
      yield r'document';
      yield serializers.serialize(
        object.document,
        specifiedType: const FullType(ReportDocument),
      );
    }
    yield r'formats';
    yield serializers.serialize(
      object.formats,
      specifiedType: const FullType(BuiltList, [FullType(ReportRequestCommandFormatsEnum)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportRequestCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportRequestCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'groupIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(String)]),
          ) as BuiltList<String>;
          result.groupIds.replace(valueDes);
          break;
        case r'period':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportPeriod),
          ) as ReportPeriod;
          result.period.replace(valueDes);
          break;
        case r'filter':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptPagedRequestFilter),
          ) as ReceiptPagedRequestFilter;
          result.filter.replace(valueDes);
          break;
        case r'groupBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(String)]),
          ) as BuiltList<String>;
          result.groupBy.replace(valueDes);
          break;
        case r'groupByLabels':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltMap, [FullType(String), FullType(String)]),
          ) as BuiltMap<String, String>;
          result.groupByLabels.replace(valueDes);
          break;
        case r'detail':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportDetail),
          ) as ReportDetail;
          result.detail.replace(valueDes);
          break;
        case r'columns':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReportColumn)]),
          ) as BuiltList<ReportColumn>;
          result.columns.replace(valueDes);
          break;
        case r'subtotals':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.subtotals = valueDes;
          break;
        case r'grandTotals':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.grandTotals = valueDes;
          break;
        case r'document':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportDocument),
          ) as ReportDocument;
          result.document.replace(valueDes);
          break;
        case r'formats':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReportRequestCommandFormatsEnum)]),
          ) as BuiltList<ReportRequestCommandFormatsEnum>;
          result.formats.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportRequestCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportRequestCommandBuilder();
    final serializedList = (serialized as Iterable<Object?>).toList();
    final unhandled = <Object?>[];
    _deserializeProperties(
      serializers,
      serialized,
      specifiedType: specifiedType,
      serializedList: serializedList,
      unhandled: unhandled,
      result: result,
    );
    return result.build();
  }
}

class ReportRequestCommandFormatsEnum extends EnumClass {

  @BuiltValueEnumConst(wireName: r'csv')
  static const ReportRequestCommandFormatsEnum csv = _$reportRequestCommandFormatsEnum_csv;
  @BuiltValueEnumConst(wireName: r'xlsx')
  static const ReportRequestCommandFormatsEnum xlsx = _$reportRequestCommandFormatsEnum_xlsx;
  @BuiltValueEnumConst(wireName: r'pdf')
  static const ReportRequestCommandFormatsEnum pdf = _$reportRequestCommandFormatsEnum_pdf;

  static Serializer<ReportRequestCommandFormatsEnum> get serializer => _$reportRequestCommandFormatsEnumSerializer;

  const ReportRequestCommandFormatsEnum._(String name): super(name);

  static BuiltSet<ReportRequestCommandFormatsEnum> get values => _$reportRequestCommandFormatsEnumValues;
  static ReportRequestCommandFormatsEnum valueOf(String name) => _$reportRequestCommandFormatsEnumValueOf(name);
}

