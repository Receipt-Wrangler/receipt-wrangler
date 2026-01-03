import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:rxdart/rxdart.dart';

import '../../client/client.dart';
import '../../models/auth_model.dart';
import '../../models/receipt_model.dart';
import '../../enums/form_state.dart';
import '../../shared/widgets/screen_wrapper.dart';
import '../../utils/snackbar.dart';
import '../widgets/receipt_comments.dart';
import '../widgets/receipt_comment_app_bar.dart';

class ReceiptCommentScreen extends StatefulWidget {
  const ReceiptCommentScreen({super.key, required this.formState});
  
  final WranglerFormState formState;

  @override
  State<ReceiptCommentScreen> createState() => _ReceiptCommentScreenState();
}

class _ReceiptCommentScreenState extends State<ReceiptCommentScreen> {
  late final formState = widget.formState;
  late final receipt =
      Provider.of<ReceiptModel>(context, listen: false).receipt;
  late final textBehaviorSubject = BehaviorSubject<String>();

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    textBehaviorSubject.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ScreenWrapper(
      appBarWidget: ReceiptCommentAppBar(
        formState: formState,
        title: 'Receipt Comments',
        onBack: () => Navigator.of(context).pop(),
        onEditMode: formState == WranglerFormState.view ? () {
          Navigator.of(context).pushReplacement(
            MaterialPageRoute(
              builder: (context) => ReceiptCommentScreen(formState: WranglerFormState.edit),
            ),
          );
        } : null,
      ),
      bottomSheetWidget: _buildCommentBottomSheet(),
      child: Consumer<ReceiptModel>(
        builder: (context, receiptModel, child) {
          return ReceiptComments(
            comments: receiptModel.comments,
            formState: formState,
          );
        },
      ),
    );
  }

  Widget _buildCommentBottomSheet() {
    if (formState == WranglerFormState.view) {
      return SizedBox.shrink();
    }

    var formKey = GlobalKey<FormBuilderState>();

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _buildCommentTextField(formKey),
        StreamBuilder(
            stream: textBehaviorSubject.stream,
            builder: (context, snapshot) {
              return _buildSubmitButton(formKey, snapshot.data);
            }),
      ],
    );
  }

  Widget _buildCommentTextField(GlobalKey<FormBuilderState> formKey) {
    return Expanded(
      child: FormBuilder(
          key: formKey,
          child: FormBuilderTextField(
            name: "comment",
            decoration: const InputDecoration(labelText: "Comment"),
            validator: FormBuilderValidators.required(),
            onChanged: (value) {
              textBehaviorSubject.add(value ?? "");
            },
          )),
    );
  }

  Widget _buildSubmitButton(GlobalKey<FormBuilderState> formKey, String? comment) {
    var isValid = comment?.isNotEmpty;
    return IconButton(
        icon: Icon(Icons.send),
        onPressed: (isValid ?? false) ? () => _submitComment(formKey) : null);
  }

  void _submitComment(GlobalKey<FormBuilderState> formKey) async {
    if (formKey.currentState?.saveAndValidate() ?? false) {
      if (formState == WranglerFormState.edit) {
        await _submitCommentToApi(formKey);
      } else if (formState == WranglerFormState.add) {
        _addCommentToModel(formKey);
      }
    }
  }

  void _addCommentToModel(GlobalKey<FormBuilderState> formKey) {
    var commentText = formKey.currentState?.value['comment'];
    var userId = Provider.of<AuthModel>(context, listen: false).claims?.userId ?? 0;
    var receiptModel = Provider.of<ReceiptModel>(context, listen: false);

    var comment = (api.CommentBuilder()
          ..id = 0
          ..comment = commentText
          ..receiptId = 0
          ..userId = userId
          ..createdAt = DateTime.now().toString())
        .build();

    var currentComments = receiptModel.comments.toList();
    currentComments.add(comment);
    receiptModel.setComments(currentComments);

    formKey.currentState?.reset();
    textBehaviorSubject.add("");
  }

  Future<void> _submitCommentToApi(GlobalKey<FormBuilderState> formKey) async {
    try {
      var commentText = formKey.currentState?.value['comment'];
      var receiptModel = Provider.of<ReceiptModel>(context, listen: false);
      
      var commentToCreate = (api.UpsertCommentCommandBuilder()
            ..comment = commentText
            ..receiptId = receiptModel.receipt.id)
          .build();

      var response = await OpenApiClient.client
          .getCommentApi()
          .addComment(upsertCommentCommand: commentToCreate);

      var currentComments = receiptModel.comments.toList();
      if (response.data != null) {
        currentComments.add(response.data!);
        receiptModel.setComments(currentComments);
      }

      formKey.currentState?.reset();
      textBehaviorSubject.add("");
      showSuccessSnackbar(context, "Comment added successfully");
    } catch (e) {
      showErrorSnackbar(context, "Failed to add comment");
    }
  }
}
