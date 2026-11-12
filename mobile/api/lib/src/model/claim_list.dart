//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/claim.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/claim_list_summary.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim_list.g.dart';

/// A group's claims and a summary of them
///
/// Properties:
/// * [claims] - Newest first. Always present; empty when none.
/// * [summary] 
@BuiltValue()
abstract class ClaimList implements Built<ClaimList, ClaimListBuilder> {
  /// Newest first. Always present; empty when none.
  @BuiltValueField(wireName: r'claims')
  BuiltList<Claim> get claims;

  @BuiltValueField(wireName: r'summary')
  ClaimListSummary get summary;

  ClaimList._();

  factory ClaimList([void updates(ClaimListBuilder b)]) = _$ClaimList;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimListBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ClaimList> get serializer => _$ClaimListSerializer();
}

class _$ClaimListSerializer implements PrimitiveSerializer<ClaimList> {
  @override
  final Iterable<Type> types = const [ClaimList, _$ClaimList];

  @override
  final String wireName = r'ClaimList';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ClaimList object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'claims';
    yield serializers.serialize(
      object.claims,
      specifiedType: const FullType(BuiltList, [FullType(Claim)]),
    );
    yield r'summary';
    yield serializers.serialize(
      object.summary,
      specifiedType: const FullType(ClaimListSummary),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ClaimList object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimListBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'claims':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Claim)]),
          ) as BuiltList<Claim>;
          result.claims.replace(valueDes);
          break;
        case r'summary':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ClaimListSummary),
          ) as ClaimListSummary;
          result.summary.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ClaimList deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimListBuilder();
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

