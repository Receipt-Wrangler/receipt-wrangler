//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'associated_api_keys.g.dart';

class AssociatedApiKeys extends EnumClass {

  @BuiltValueEnumConst(wireName: r'MINE')
  static const AssociatedApiKeys MINE = _$MINE;
  @BuiltValueEnumConst(wireName: r'ALL')
  static const AssociatedApiKeys ALL = _$ALL;

  static Serializer<AssociatedApiKeys> get serializer => _$associatedApiKeysSerializer;

  const AssociatedApiKeys._(String name): super(name);

  static BuiltSet<AssociatedApiKeys> get values => _$values;
  static AssociatedApiKeys valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class AssociatedApiKeysMixin = Object with _$AssociatedApiKeysMixin;

