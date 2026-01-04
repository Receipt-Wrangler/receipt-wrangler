import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for CommentApi
void main() {
  final instance = Openapi().getCommentApi();

  group(CommentApi, () {
    // Add comment
    //
    // This will add a comment to a receipt, [SYSTEM USER]
    //
    //Future<Comment> addComment(UpsertCommentCommand upsertCommentCommand) async
    test('test addComment', () async {
      // TODO
    });

    // Delete comment
    //
    // This will delete a comment by id [SYSTEM User]
    //
    //Future deleteComment(int commentId) async
    test('test deleteComment', () async {
      // TODO
    });

  });
}
