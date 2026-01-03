//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'system_task_status.g.dart';

class SystemTaskStatus extends EnumClass {

  @BuiltValueEnumConst(wireName: r'SUCCEEDED')
  static const SystemTaskStatus SUCCEEDED = _$SUCCEEDED;
  @BuiltValueEnumConst(wireName: r'FAILED')
  static const SystemTaskStatus FAILED = _$FAILED;

  static Serializer<SystemTaskStatus> get serializer => _$systemTaskStatusSerializer;

  const SystemTaskStatus._(String name): super(name);

  static BuiltSet<SystemTaskStatus> get values => _$values;
  static SystemTaskStatus valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class SystemTaskStatusMixin = Object with _$SystemTaskStatusMixin;

