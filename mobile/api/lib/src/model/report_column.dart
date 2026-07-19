//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'report_column.g.dart';

/// ReportColumn
///
/// Properties:
/// * [kind] 
/// * [name] - Machine identifier; formula expressions reference this
/// * [label] 
/// * [field] - Field key the column displays (dimension columns)
/// * [aggFunc] - Aggregate function (aggregate columns)
/// * [measure] - Measure field key (aggregate columns, omitted for COUNT)
/// * [expr] - Expression over other column names (formula columns)
@BuiltValue()
abstract class ReportColumn implements Built<ReportColumn, ReportColumnBuilder> {
  @BuiltValueField(wireName: r'kind')
  ReportColumnKindEnum get kind;
  // enum kindEnum {  dimension,  aggregate,  formula,  };

  /// Machine identifier; formula expressions reference this
  @BuiltValueField(wireName: r'name')
  String get name;

  @BuiltValueField(wireName: r'label')
  String? get label;

  /// Field key the column displays (dimension columns)
  @BuiltValueField(wireName: r'field')
  String? get field;

  /// Aggregate function (aggregate columns)
  @BuiltValueField(wireName: r'aggFunc')
  ReportColumnAggFuncEnum? get aggFunc;
  // enum aggFuncEnum {  SUM,  COUNT,  AVG,  MIN,  MAX,  };

  /// Measure field key (aggregate columns, omitted for COUNT)
  @BuiltValueField(wireName: r'measure')
  String? get measure;

  /// Expression over other column names (formula columns)
  @BuiltValueField(wireName: r'expr')
  String? get expr;

  ReportColumn._();

  factory ReportColumn([void updates(ReportColumnBuilder b)]) = _$ReportColumn;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReportColumnBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReportColumn> get serializer => _$ReportColumnSerializer();
}

class _$ReportColumnSerializer implements PrimitiveSerializer<ReportColumn> {
  @override
  final Iterable<Type> types = const [ReportColumn, _$ReportColumn];

  @override
  final String wireName = r'ReportColumn';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReportColumn object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'kind';
    yield serializers.serialize(
      object.kind,
      specifiedType: const FullType(ReportColumnKindEnum),
    );
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    if (object.label != null) {
      yield r'label';
      yield serializers.serialize(
        object.label,
        specifiedType: const FullType(String),
      );
    }
    if (object.field != null) {
      yield r'field';
      yield serializers.serialize(
        object.field,
        specifiedType: const FullType(String),
      );
    }
    if (object.aggFunc != null) {
      yield r'aggFunc';
      yield serializers.serialize(
        object.aggFunc,
        specifiedType: const FullType(ReportColumnAggFuncEnum),
      );
    }
    if (object.measure != null) {
      yield r'measure';
      yield serializers.serialize(
        object.measure,
        specifiedType: const FullType(String),
      );
    }
    if (object.expr != null) {
      yield r'expr';
      yield serializers.serialize(
        object.expr,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ReportColumn object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReportColumnBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'kind':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportColumnKindEnum),
          ) as ReportColumnKindEnum;
          result.kind = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'label':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.label = valueDes;
          break;
        case r'field':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.field = valueDes;
          break;
        case r'aggFunc':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReportColumnAggFuncEnum),
          ) as ReportColumnAggFuncEnum;
          result.aggFunc = valueDes;
          break;
        case r'measure':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.measure = valueDes;
          break;
        case r'expr':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.expr = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReportColumn deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReportColumnBuilder();
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

class ReportColumnKindEnum extends EnumClass {

  @BuiltValueEnumConst(wireName: r'dimension')
  static const ReportColumnKindEnum dimension = _$reportColumnKindEnum_dimension;
  @BuiltValueEnumConst(wireName: r'aggregate')
  static const ReportColumnKindEnum aggregate = _$reportColumnKindEnum_aggregate;
  @BuiltValueEnumConst(wireName: r'formula')
  static const ReportColumnKindEnum formula = _$reportColumnKindEnum_formula;

  static Serializer<ReportColumnKindEnum> get serializer => _$reportColumnKindEnumSerializer;

  const ReportColumnKindEnum._(String name): super(name);

  static BuiltSet<ReportColumnKindEnum> get values => _$reportColumnKindEnumValues;
  static ReportColumnKindEnum valueOf(String name) => _$reportColumnKindEnumValueOf(name);
}

class ReportColumnAggFuncEnum extends EnumClass {

  /// Aggregate function (aggregate columns)
  @BuiltValueEnumConst(wireName: r'SUM')
  static const ReportColumnAggFuncEnum SUM = _$reportColumnAggFuncEnum_SUM;
  /// Aggregate function (aggregate columns)
  @BuiltValueEnumConst(wireName: r'COUNT')
  static const ReportColumnAggFuncEnum COUNT = _$reportColumnAggFuncEnum_COUNT;
  /// Aggregate function (aggregate columns)
  @BuiltValueEnumConst(wireName: r'AVG')
  static const ReportColumnAggFuncEnum AVG = _$reportColumnAggFuncEnum_AVG;
  /// Aggregate function (aggregate columns)
  @BuiltValueEnumConst(wireName: r'MIN')
  static const ReportColumnAggFuncEnum MIN = _$reportColumnAggFuncEnum_MIN;
  /// Aggregate function (aggregate columns)
  @BuiltValueEnumConst(wireName: r'MAX')
  static const ReportColumnAggFuncEnum MAX = _$reportColumnAggFuncEnum_MAX;

  static Serializer<ReportColumnAggFuncEnum> get serializer => _$reportColumnAggFuncEnumSerializer;

  const ReportColumnAggFuncEnum._(String name): super(name);

  static BuiltSet<ReportColumnAggFuncEnum> get values => _$reportColumnAggFuncEnumValues;
  static ReportColumnAggFuncEnum valueOf(String name) => _$reportColumnAggFuncEnumValueOf(name);
}

