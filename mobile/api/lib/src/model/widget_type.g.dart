// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'widget_type.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const WidgetType _$GROUP_SUMMARY = const WidgetType._('GROUP_SUMMARY');
const WidgetType _$FILTERED_RECEIPTS = const WidgetType._('FILTERED_RECEIPTS');
const WidgetType _$GROUP_ACTIVITY = const WidgetType._('GROUP_ACTIVITY');

WidgetType _$valueOf(String name) {
  switch (name) {
    case 'GROUP_SUMMARY':
      return _$GROUP_SUMMARY;
    case 'FILTERED_RECEIPTS':
      return _$FILTERED_RECEIPTS;
    case 'GROUP_ACTIVITY':
      return _$GROUP_ACTIVITY;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<WidgetType> _$values =
    new BuiltSet<WidgetType>(const <WidgetType>[
  _$GROUP_SUMMARY,
  _$FILTERED_RECEIPTS,
  _$GROUP_ACTIVITY,
]);

class _$WidgetTypeMeta {
  const _$WidgetTypeMeta();
  WidgetType get GROUP_SUMMARY => _$GROUP_SUMMARY;
  WidgetType get FILTERED_RECEIPTS => _$FILTERED_RECEIPTS;
  WidgetType get GROUP_ACTIVITY => _$GROUP_ACTIVITY;
  WidgetType valueOf(String name) => _$valueOf(name);
  BuiltSet<WidgetType> get values => _$values;
}

abstract class _$WidgetTypeMixin {
  // ignore: non_constant_identifier_names
  _$WidgetTypeMeta get WidgetType => const _$WidgetTypeMeta();
}

Serializer<WidgetType> _$widgetTypeSerializer = new _$WidgetTypeSerializer();

class _$WidgetTypeSerializer implements PrimitiveSerializer<WidgetType> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'GROUP_SUMMARY': 'GROUP_SUMMARY',
    'FILTERED_RECEIPTS': 'FILTERED_RECEIPTS',
    'GROUP_ACTIVITY': 'GROUP_ACTIVITY',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'GROUP_SUMMARY': 'GROUP_SUMMARY',
    'FILTERED_RECEIPTS': 'FILTERED_RECEIPTS',
    'GROUP_ACTIVITY': 'GROUP_ACTIVITY',
  };

  @override
  final Iterable<Type> types = const <Type>[WidgetType];
  @override
  final String wireName = 'WidgetType';

  @override
  Object serialize(Serializers serializers, WidgetType object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  WidgetType deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      WidgetType.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
