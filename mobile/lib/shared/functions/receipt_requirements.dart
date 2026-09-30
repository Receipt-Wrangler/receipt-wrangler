import 'package:openapi/openapi.dart';

import '../../enums/form_state.dart';
import '../../models/receipt_model.dart';

/// Client-side mirror of the server's role-required receipt fields
/// (`AppData.groupReceiptRequirements`, see `PermissionsModel.receiptRequirements`).
///
/// The server is the enforcement point — it 400s a create/update missing a
/// required comment or image, and refuses deleting the last one. These helpers
/// only let the UI say so before the round trip, so they apply the same tests
/// the server does.

/// Whether [comments] holds at least one comment the server counts: a
/// whitespace-only comment does not satisfy the requirement.
bool hasNonBlankComment(Iterable<Comment> comments) =>
    countNonBlankComments(comments) > 0;

/// How many of [comments] the server counts toward the requirement.
int countNonBlankComments(Iterable<Comment> comments) =>
    comments.where((c) => c.comment.trim().isNotEmpty).length;

/// The message to show when a receipt misses what its group's role requires,
/// or null when nothing is missing. [hasComment] / [hasImage] describe the
/// receipt as it will be saved.
String? missingReceiptRequirementsMessage(
  ReceiptRequirements requirements, {
  required bool hasComment,
  required bool hasImage,
}) {
  final missingComment = requirements.commentRequired && !hasComment;
  final missingImage = requirements.imageRequired && !hasImage;

  if (missingComment && missingImage) {
    return "Your role in this group requires at least one image and one "
        "comment on each receipt.";
  }
  if (missingImage) {
    return "Your role in this group requires at least one image on each "
        "receipt.";
  }
  if (missingComment) {
    return "Your role in this group requires at least one comment on each "
        "receipt.";
  }
  return null;
}

/// Whether the receipt on [receiptModel] has an image, as the server will
/// judge it when [formState] is submitted.
///
/// **Add:** the staged images, which ride the create call itself.
///
/// **Edit:** the saved images — edit uploads and deletes each image through its
/// own endpoint immediately, and the server checks the stored receipt. The
/// loaded image list ([ReceiptModel.imageBehaviorSubject]) is only populated
/// once the images screen has been opened, so until then the receipt's own
/// `imageFiles` stand in. The one case that reads wrong — the images screen
/// loaded and every image deleted — cannot happen in a group that requires an
/// image, because deleting the last one is disabled there; the server still
/// refuses anything this misjudges.
bool receiptHasImage(ReceiptModel receiptModel, WranglerFormState formState) {
  if (formState == WranglerFormState.add) {
    return receiptModel.imagesToUploadBehaviorSubject.value.isNotEmpty;
  }
  if (receiptModel.imageBehaviorSubject.value.any((image) => image != null)) {
    return true;
  }
  return receiptModel.receipt.imageFiles?.isNotEmpty ?? false;
}

/// The message blocking a submit of [receiptModel] in [formState] under
/// [requirements], or null when the submit may go ahead. Comments are the
/// model's in both states: staged ones in add, and in edit the list the
/// comment screen keeps in step with each add/delete call.
String? receiptSubmitRequirementsMessage(
  ReceiptRequirements requirements, {
  required ReceiptModel receiptModel,
  required WranglerFormState formState,
}) => missingReceiptRequirementsMessage(
  requirements,
  hasComment: hasNonBlankComment(receiptModel.comments),
  hasImage: receiptHasImage(receiptModel, formState),
);

/// Whether removing one of [remaining] required items would leave none — the
/// server refuses that delete, so the UI disables it instead of letting it fail.
bool isLastRequiredItem({required bool required, required int remaining}) =>
    required && remaining <= 1;
