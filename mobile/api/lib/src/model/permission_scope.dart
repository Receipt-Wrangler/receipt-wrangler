//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'permission_scope.g.dart';

class PermissionScope extends EnumClass {

  /// Whether a permission applies app-wide or within a single group.
  @BuiltValueEnumConst(wireName: r'APP')
  static const PermissionScope APP = _$APP;
  /// Whether a permission applies app-wide or within a single group.
  @BuiltValueEnumConst(wireName: r'GROUP')
  static const PermissionScope GROUP = _$GROUP;

  static Serializer<PermissionScope> get serializer => _$permissionScopeSerializer;

  const PermissionScope._(String name): super(name);

  static BuiltSet<PermissionScope> get values => _$values;
  static PermissionScope valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class PermissionScopeMixin = Object with _$PermissionScopeMixin;

