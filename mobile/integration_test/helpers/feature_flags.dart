import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;

import 'api.dart';
import 'env.dart';

/// Name prefix stamped on the throwaway receipt-processing records this suite
/// creates. [_createProcessingSettings] writes it and
/// [_isFixtureProcessingSettings] reads it back to tell a leaked fixture record
/// apart from a real AI configuration, so the two must never drift -- hence a
/// constant rather than two literals.
const _fixtureRpsNamePrefix = 'e2e-ai-flag-';

/// Enables the `aiPoweredReceipts` feature flag (which gates the Quick Scan UI)
/// for the duration of a test, then restores it via `addTearDown`. Call it
/// BEFORE logging in -- `featureConfig` hydrates from AppData at login.
///
/// The flag is computed server-side purely as
/// `systemSettings.receiptProcessingSettingsId != null`
/// (`api/internal/services/system_settings.go`), so flipping it needs no real
/// AI provider — we point system settings at a throwaway receipt-processing
/// record. The AI never runs (and isn't what these tests cover); only the flag
/// matters.
///
/// Teardown restores the pointer to whatever it was before the test and deletes
/// the throwaway record -- always writing null would clobber a real AI
/// configuration on non-pristine backends. The one exception is a pointer at a
/// record THIS SUITE created, which is a leak rather than a configuration; see
/// [_restorePointerPatch].
///
/// Unlike [disableAiPoweredReceiptsForTest], this heals the leaked *pointer*
/// but does NOT delete the record it pointed at: a concurrently-running job's
/// live fixture record is indistinguishable from a leaked orphan by name, and
/// deleting that would leave a dangling pointer nothing can heal. The orphan is
/// inert (the flag depends only on the pointer) and gets swept the next time a
/// disable finds it live.
///
/// The counterpart is [disableAiPoweredReceiptsForTest] -- specs that need the
/// flag OFF must call it rather than assume, because the flag is install-wide.
Future<void> enableAiPoweredReceiptsForTest() async {
  final jwt = await apiLogin();
  final promptId = await _ensurePromptId(jwt);
  final rpsId = await _createProcessingSettings(jwt, promptId);
  final settings = await getSystemSettings(jwt);
  final restoreTo = await _restorePointerPatch(jwt, settings);

  // Registered BEFORE the override so LIFO runs it AFTER the restore: system
  // settings still point at this record until then. Same ordering trap as the
  // role / custom-field teardowns in permission_fixtures.dart.
  //
  // Reuses the setup jwt rather than logging in again: login is rate-limited
  // (429s on tight reruns), and splitting this out of the restore teardown must
  // not turn one login per fixture call into two. The token outlives any single
  // test by a wide margin (20 min), and the delete is best-effort anyway.
  addTearDown(() async => _deleteProcessingSettings(jwt, rpsId));

  await overrideSystemSettingsForTest(
    jwt,
    settings,
    overrides: {'receiptProcessingSettingsId': rpsId},
    restoreTo: restoreTo,
  );
}

/// Pins the `aiPoweredReceipts` feature flag OFF for the duration of a test,
/// restoring the previous pointers via `addTearDown`. The counterpart to
/// [enableAiPoweredReceiptsForTest]; call it BEFORE logging in, since
/// `featureConfig` hydrates from AppData at login and is read `listen: false`.
///
/// The flag is INSTALL-WIDE, so "off" is AMBIENT state on a shared backend, not
/// a default. Any run that aborts between a fixture's PUT and its teardown
/// leaves it stuck on for every later run -- one such leak kept it on for eight
/// weeks and read as a product bug (`Timed out after 10s waiting for text
/// "Name"`, because the scan slot really does open a scanner instead of falling
/// through to the manual form). A spec that depends on the flag being off has to
/// say so.
///
/// SELF-HEALING: when the captured pointer names a record this suite created it
/// is a leak, so the restore writes null instead of replaying it and the leaked
/// record is deleted. Both happen in SETUP rather than teardown -- the state
/// being healed was itself produced by a teardown that never ran, so a heal that
/// only runs on teardown would inherit the same failure mode. A pointer at any
/// other record is a real AI configuration and is restored exactly.
///
/// No-ops (no write, no teardown) when the flag is already off, which is the
/// normal case -- on a healthy backend this fixture costs nothing and touches no
/// shared state.
Future<void> disableAiPoweredReceiptsForTest() async {
  final jwt = await apiLogin();
  final settings = await getSystemSettings(jwt);
  final originalRpsId = settings['receiptProcessingSettingsId'] as int?;
  if (originalRpsId == null) return;

  final restoreTo = await _restorePointerPatch(jwt, settings);
  // A null restore for a non-null original is _restorePointerPatch saying "this
  // pointer is one of ours".
  final leaked = restoreTo['receiptProcessingSettingsId'] == null;

  await overrideSystemSettingsForTest(
    jwt,
    settings,
    overrides: _pointerPatch(null, null),
    restoreTo: restoreTo,
  );

  // Only after the PUT above has pointed away from it. The fallback guard is
  // theoretical (fixtures never set a fallback) but deleting a row the fallback
  // still references would leave a dangling pointer.
  if (leaked &&
      settings['fallbackReceiptProcessingSettingsId'] != originalRpsId) {
    await _deleteProcessingSettings(jwt, originalRpsId);
  }
}

/// The two receipt-processing pointers, always written together.
///
/// The backend rejects a fallback without a primary ("Fallback receipt
/// processing settings ID cannot be set without receipt processing settings ID",
/// `api/internal/commands/upsert_system_settings_command.go`), so nulling the
/// primary while leaving a configured fallback in place would 400 the PUT.
Map<String, dynamic> _pointerPatch(int? rpsId, int? fallbackRpsId) => {
      'receiptProcessingSettingsId': rpsId,
      'fallbackReceiptProcessingSettingsId':
          rpsId == null ? null : fallbackRpsId,
    };

/// What a fixture should put back for the AI pointers, given the [settings] it
/// captured before touching anything.
///
/// Normally the captured values. But if the captured primary points at a record
/// THIS SUITE created, the capture is a LEAK from a run that aborted before its
/// teardown -- replaying it would faithfully re-wedge the flag ON for the next
/// spec, which is exactly the failure this fixture pair exists to stop. In that
/// case the restore is null. The name check is what separates the two cases:
/// always writing null would clobber a real AI configuration.
Future<Map<String, dynamic>> _restorePointerPatch(
  String jwt,
  Map<String, dynamic> settings,
) async {
  final rpsId = settings['receiptProcessingSettingsId'] as int?;
  if (rpsId != null && await _isFixtureProcessingSettings(jwt, rpsId)) {
    return _pointerPatch(null, null);
  }
  return _pointerPatch(
    rpsId,
    settings['fallbackReceiptProcessingSettingsId'] as int?,
  );
}

/// Whether [rpsId] names a record this suite created (`e2e-ai-flag-...`).
///
/// Only a positively-read matching name counts. A read failure answers false:
/// `GET /receiptProcessingSettings/{id}` returns 500 for BOTH "no such record"
/// and a transient error (the handler maps every repository error to 500), so a
/// missing record can't be told from a blip. The asymmetry is deliberate --
/// being wrong this way costs one ambient-state test failure, while being wrong
/// the other way would silently unhook and delete a real install's AI
/// configuration.
Future<bool> _isFixtureProcessingSettings(String jwt, int rpsId) async {
  try {
    final res = await http
        .get(
          Uri.parse('${E2eEnv.baseUrl}/receiptProcessingSettings/$rpsId'),
          headers: {'Cookie': 'jwt=$jwt'},
        )
        .timeout(const Duration(seconds: 10));
    if (res.statusCode != 200) return false;
    final name = (jsonDecode(res.body) as Map<String, dynamic>)['name'];
    return name is String && name.startsWith(_fixtureRpsNamePrefix);
  } catch (_) {
    return false;
  }
}

Future<int> _ensurePromptId(String jwt) async {
  final list = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/prompt/getPagedPrompts'),
        headers: jsonAuthHeaders(jwt),
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
        headers: jsonAuthHeaders(jwt),
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
        headers: jsonAuthHeaders(jwt),
        body: jsonEncode({
          'name': '$_fixtureRpsNamePrefix'
              '${DateTime.now().microsecondsSinceEpoch}',
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
