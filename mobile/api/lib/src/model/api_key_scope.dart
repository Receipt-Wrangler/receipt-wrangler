//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'api_key_scope.g.dart';

class ApiKeyScope extends EnumClass {

  /// Scope/permissions for API keys
  @BuiltValueEnumConst(wireName: r'r')
  static const ApiKeyScope r = _$r;
  /// Scope/permissions for API keys
  @BuiltValueEnumConst(wireName: r'w')
  static const ApiKeyScope w = _$w;
  /// Scope/permissions for API keys
  @BuiltValueEnumConst(wireName: r'rw')
  static const ApiKeyScope rw = _$rw;

  static Serializer<ApiKeyScope> get serializer => _$apiKeyScopeSerializer;

  const ApiKeyScope._(String name): super(name);

  static BuiltSet<ApiKeyScope> get values => _$values;
  static ApiKeyScope valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class ApiKeyScopeMixin = Object with _$ApiKeyScopeMixin;

