//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_comment_command.g.dart';

/// UpsertCommentCommand
///
/// Properties:
/// * [comment] - Comment itself
/// * [receiptId] - Receipt foreign key
/// * [userId] - User foreign key
@BuiltValue()
abstract class UpsertCommentCommand implements Built<UpsertCommentCommand, UpsertCommentCommandBuilder> {
  /// Comment itself
  @BuiltValueField(wireName: r'comment')
  String get comment;

  /// Receipt foreign key
  @BuiltValueField(wireName: r'receiptId')
  int get receiptId;

  /// User foreign key
  @BuiltValueField(wireName: r'userId')
  int? get userId;

  UpsertCommentCommand._();

  factory UpsertCommentCommand([void updates(UpsertCommentCommandBuilder b)]) = _$UpsertCommentCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertCommentCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertCommentCommand> get serializer => _$UpsertCommentCommandSerializer();
}

class _$UpsertCommentCommandSerializer implements PrimitiveSerializer<UpsertCommentCommand> {
  @override
  final Iterable<Type> types = const [UpsertCommentCommand, _$UpsertCommentCommand];

  @override
  final String wireName = r'UpsertCommentCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertCommentCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'comment';
    yield serializers.serialize(
      object.comment,
      specifiedType: const FullType(String),
    );
    yield r'receiptId';
    yield serializers.serialize(
      object.receiptId,
      specifiedType: const FullType(int),
    );
    if (object.userId != null) {
      yield r'userId';
      yield serializers.serialize(
        object.userId,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertCommentCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertCommentCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'comment':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.comment = valueDes;
          break;
        case r'receiptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptId = valueDes;
          break;
        case r'userId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.userId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertCommentCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertCommentCommandBuilder();
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

