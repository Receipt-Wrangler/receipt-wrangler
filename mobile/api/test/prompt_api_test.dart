import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for PromptApi
void main() {
  final instance = Openapi().getPromptApi();

  group(PromptApi, () {
    // Create default prompt
    //
    // This will create a default prompt
    //
    //Future<Prompt> createDefaultPrompt() async
    test('test createDefaultPrompt', () async {
      // TODO
    });

    // Create prompt
    //
    // This will create a prompt
    //
    //Future<Prompt> createPrompt(UpsertPromptCommand upsertPromptCommand) async
    test('test createPrompt', () async {
      // TODO
    });

    // Delete prompt by id
    //
    // This will delete a prompt by id
    //
    //Future deletePromptById(int id) async
    test('test deletePromptById', () async {
      // TODO
    });

    // Gets paged prompts
    //
    // This will return paged prompts
    //
    //Future<PagedData> getPagedPrompts(PagedRequestCommand pagedRequestCommand) async
    test('test getPagedPrompts', () async {
      // TODO
    });

    // Get prompt by id
    //
    // This will get a prompt by id
    //
    //Future<Prompt> getPromptById(int id) async
    test('test getPromptById', () async {
      // TODO
    });

    // Update prompt by id
    //
    // This will update a prompt by id
    //
    //Future<Prompt> updatePromptById(int id, UpsertPromptCommand upsertPromptCommand) async
    test('test updatePromptById', () async {
      // TODO
    });

  });
}
