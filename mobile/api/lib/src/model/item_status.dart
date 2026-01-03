//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'item_status.g.dart';

class ItemStatus extends EnumClass {

  @BuiltValueEnumConst(wireName: r'OPEN')
  static const ItemStatus OPEN = _$OPEN;
  @BuiltValueEnumConst(wireName: r'RESOLVED')
  static const ItemStatus RESOLVED = _$RESOLVED;
  @BuiltValueEnumConst(wireName: r'DRAFT')
  static const ItemStatus DRAFT = _$DRAFT;

  static Serializer<ItemStatus> get serializer => _$itemStatusSerializer;

  const ItemStatus._(String name): super(name);

  static BuiltSet<ItemStatus> get values => _$values;
  static ItemStatus valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class ItemStatusMixin = Object with _$ItemStatusMixin;

