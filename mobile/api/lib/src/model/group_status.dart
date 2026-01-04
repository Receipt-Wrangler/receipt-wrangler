//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'group_status.g.dart';

class GroupStatus extends EnumClass {

  @BuiltValueEnumConst(wireName: r'ACTIVE')
  static const GroupStatus ACTIVE = _$ACTIVE;
  @BuiltValueEnumConst(wireName: r'ARCHIVED')
  static const GroupStatus ARCHIVED = _$ARCHIVED;

  static Serializer<GroupStatus> get serializer => _$groupStatusSerializer;

  const GroupStatus._(String name): super(name);

  static BuiltSet<GroupStatus> get values => _$values;
  static GroupStatus valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class GroupStatusMixin = Object with _$GroupStatusMixin;

