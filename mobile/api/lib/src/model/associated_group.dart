//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'associated_group.g.dart';

class AssociatedGroup extends EnumClass {

  @BuiltValueEnumConst(wireName: r'MINE')
  static const AssociatedGroup MINE = _$MINE;
  @BuiltValueEnumConst(wireName: r'ALL')
  static const AssociatedGroup ALL = _$ALL;

  static Serializer<AssociatedGroup> get serializer => _$associatedGroupSerializer;

  const AssociatedGroup._(String name): super(name);

  static BuiltSet<AssociatedGroup> get values => _$values;
  static AssociatedGroup valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class AssociatedGroupMixin = Object with _$AssociatedGroupMixin;

