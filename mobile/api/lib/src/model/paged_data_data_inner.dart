//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/system_email.dart';
import 'package:openapi/src/model/prompt.dart';
import 'package:openapi/src/model/custom_field_value.dart';
import 'package:openapi/src/model/tag.dart';
import 'package:openapi/src/model/group_member.dart';
import 'package:openapi/src/model/group.dart';
import 'package:openapi/src/model/receipt.dart';
import 'package:openapi/src/model/group_receipt_settings.dart';
import 'package:openapi/src/model/activity.dart';
import 'package:openapi/src/model/group_settings.dart';
import 'package:openapi/src/model/associated_entity_type.dart';
import 'package:openapi/src/model/comment.dart';
import 'package:openapi/src/model/custom_field_type.dart';
import 'package:openapi/src/model/file_data.dart';
import 'package:openapi/src/model/item.dart';
import 'package:openapi/src/model/system_task_status.dart';
import 'package:openapi/src/model/tag_view.dart';
import 'package:openapi/src/model/custom_field.dart';
import 'package:openapi/src/model/ocr_engine.dart';
import 'package:openapi/src/model/system_task.dart';
import 'package:openapi/src/model/ai_type.dart';
import 'package:openapi/src/model/custom_field_option.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/category.dart';
import 'package:openapi/src/model/receipt_processing_settings.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';
import 'package:one_of/any_of.dart';

part 'paged_data_data_inner.g.dart';

/// PagedDataDataInner
///
/// Properties:
/// * [amount] - Receipt total amount
/// * [categories] - Categories associated to receipt
/// * [comments] - Comments associated to receipt
/// * [customFields] - Custom fields associated to receipt
/// * [createdAt] 
/// * [createdBy] 
/// * [date] - Receipt date
/// * [groupId] 
/// * [id] 
/// * [imageFiles] - Files associated to receipt
/// * [name] - Custom Field name
/// * [paidByUserId] - User paid foreign key
/// * [receiptItems] - Items associated to receipt
/// * [resolvedDate] - Date resolved
/// * [status] 
/// * [tags] - Tags associated to receipt
/// * [updatedAt] 
/// * [createdByString] - Created by entity's name
/// * [description] - Custom Field description
/// * [prompt] 
/// * [groupSettings] 
/// * [groupReceiptSettings] 
/// * [groupMembers] - Members of the group
/// * [isDefault] - Is default group (not used yet)
/// * [isAllGroup] - Is all group for user
/// * [numberOfReceipts] - Number of receipts associated with this tag
/// * [type] 
/// * [startedAt] 
/// * [endedAt] 
/// * [associatedEntityId] 
/// * [associatedEntityType] 
/// * [ranByUserId] 
/// * [receiptId] 
/// * [resultDescription] 
/// * [apiKeyId] 
/// * [childSystemTasks] 
/// * [aiType] 
/// * [url] - URL for custom endpoints
/// * [key] - Key for endpoints that require authentication
/// * [model] - LLM model
/// * [isVisionModel] - Is vision model
/// * [ocrEngine] 
/// * [promptId] - Prompt foreign key
/// * [host] - IMAP host
/// * [port] - IMAP port
/// * [username] - IMAP username
/// * [password] - IMAP password
/// * [useStartTLS] - Whether to use STARTTLS
/// * [canBeRestarted] 
/// * [options] 
@BuiltValue()
abstract class PagedDataDataInner implements Built<PagedDataDataInner, PagedDataDataInnerBuilder> {
  /// Any Of [Activity], [Category], [CustomField], [Group], [Prompt], [Receipt], [ReceiptProcessingSettings], [SystemEmail], [SystemTask], [Tag], [TagView]
  AnyOf get anyOf;

  PagedDataDataInner._();

  factory PagedDataDataInner([void updates(PagedDataDataInnerBuilder b)]) = _$PagedDataDataInner;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(PagedDataDataInnerBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<PagedDataDataInner> get serializer => _$PagedDataDataInnerSerializer();
}

class _$PagedDataDataInnerSerializer implements PrimitiveSerializer<PagedDataDataInner> {
  @override
  final Iterable<Type> types = const [PagedDataDataInner, _$PagedDataDataInner];

  @override
  final String wireName = r'PagedDataDataInner';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    PagedDataDataInner object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
  }

  @override
  Object serialize(
    Serializers serializers,
    PagedDataDataInner object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final anyOf = object.anyOf;
    return serializers.serialize(anyOf, specifiedType: FullType(AnyOf, anyOf.valueTypes.map((type) => FullType(type)).toList()))!;
  }

  @override
  PagedDataDataInner deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = PagedDataDataInnerBuilder();
    Object? anyOfDataSrc;
    final targetType = const FullType(AnyOf, [FullType(Receipt), FullType(Category), FullType(Tag), FullType(Prompt), FullType(Group), FullType(TagView), FullType(SystemTask), FullType(ReceiptProcessingSettings), FullType(SystemEmail), FullType(Activity), FullType(CustomField), ]);
    anyOfDataSrc = serialized;
    result.anyOf = serializers.deserialize(anyOfDataSrc, specifiedType: targetType) as AnyOf;
    return result.build();
  }
}

