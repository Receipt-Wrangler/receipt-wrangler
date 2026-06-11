import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;

import 'api.dart';
import 'env.dart';

/// Enables the `aiPoweredReceipts` feature flag (which gates the Quick Scan UI)
/// for the duration of a test, then restores it via `addTearDown`.
///
/// The flag is computed server-side purely as
/// `systemSettings.receiptProcessingSettingsId != null`
/// (`api/internal/services/system_settings.go`), so flipping it needs no real
/// AI provider — we point system settings at a throwaway receipt-processing
/// record. The AI never runs (and isn't what these tests cover); only the flag
/// matters. Teardown restores the pointer to whatever it was before the test
/// (null on a default local backend, so `quick_scan_disabled_test` still sees
/// the flag off) and deletes the throwaway record -- always writing null here
/// would clobber a real AI configuration on non-pristine backends.
Future<void> enableAiPoweredReceiptsForTest() async {
  final jwt = await apiLogin();
  final promptId = await _ensurePromptId(jwt);
  final rpsId = await _createProcessingSettings(jwt, promptId);
  final settings = await _getSystemSettings(jwt);
  final originalRpsId = settings['receiptProcessingSettingsId'] as int?;
  await _putSystemSettingsWithRpsId(jwt, settings, rpsId);

  addTearDown(() async {
    final j = await apiLogin();
    final current = await _getSystemSettings(j);
    await _putSystemSettingsWithRpsId(j, current, originalRpsId);
    await _deleteProcessingSettings(j, rpsId);
  });
}

Map<String, String> _jsonAuth(String jwt) => {
      'Content-Type': 'application/json',
      'Cookie': 'jwt=$jwt',
    };

Future<int> _ensurePromptId(String jwt) async {
  final list = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/prompt/getPagedPrompts'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({
          'page': 1,
          'pageSize': 1,
          'orderBy': 'name',
          'sortDirection': 'asc',
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (list.statusCode == 200) {
    final data = (jsonDecode(list.body)['data'] as List?) ?? const [];
    if (data.isNotEmpty) {
      return (data.first as Map<String, dynamic>)['id'] as int;
    }
  }
  // None exist yet -- create the default prompt.
  final created = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/prompt/createDefaultPrompt'),
        headers: _jsonAuth(jwt),
      )
      .timeout(const Duration(seconds: 10));
  if (created.statusCode != 200) {
    throw StateError('createDefaultPrompt failed: '
        'HTTP ${created.statusCode}: ${created.body}');
  }
  return (jsonDecode(created.body) as Map<String, dynamic>)['id'] as int;
}

Future<int> _createProcessingSettings(String jwt, int promptId) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/receiptProcessingSettings/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({
          'name': 'e2e-ai-flag-${DateTime.now().microsecondsSinceEpoch}',
          'aiType': 'OLLAMA',
          'url': 'http://localhost:11434',
          'model': 'e2e-noop',
          'isVisionModel': true,
          'promptId': promptId,
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('create receiptProcessingSettings failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

Future<Map<String, dynamic>> _getSystemSettings(String jwt) async {
  final res = await http
      .get(Uri.parse('${E2eEnv.baseUrl}/systemSettings/'),
          headers: {'Cookie': 'jwt=$jwt'})
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
        'GET systemSettings failed: HTTP ${res.statusCode}: ${res.body}');
  }
  return jsonDecode(res.body) as Map<String, dynamic>;
}

Future<void> _putSystemSettingsWithRpsId(
  String jwt,
  Map<String, dynamic> settings,
  int? rpsId,
) async {
  final body = Map<String, dynamic>.from(settings)
    ..['receiptProcessingSettingsId'] = rpsId;
  final res = await http
      .put(
        Uri.parse('${E2eEnv.baseUrl}/systemSettings/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode(body),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
        'PUT systemSettings failed: HTTP ${res.statusCode}: ${res.body}');
  }
}

Future<void> _deleteProcessingSettings(String jwt, int rpsId) async {
  try {
    await http
        .delete(Uri.parse('${E2eEnv.baseUrl}/receiptProcessingSettings/$rpsId'),
            headers: {'Cookie': 'jwt=$jwt'})
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // best-effort
  }
}
