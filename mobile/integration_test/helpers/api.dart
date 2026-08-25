import 'dart:convert';

import 'package:http/http.dart' as http;

import 'env.dart';

/// Logs into the Go API as [username]/[password] and returns the JWT cookie
/// value.
///
/// The API issues auth via `Set-Cookie: jwt=…` (the body's `jwt` field
/// is empty — confirmed against a live demo response). We just parse
/// the cookie out of the Set-Cookie header.
Future<String> apiLoginAs(String username, String password) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/login/'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'password': password,
        }),
      )
      .timeout(const Duration(seconds: 10));

  if (res.statusCode != 200) {
    throw StateError(
      'apiLoginAs($username) failed: HTTP ${res.statusCode}: ${res.body}',
    );
  }
  final setCookie = res.headers['set-cookie'] ?? '';
  final match = RegExp(r'jwt=([^;]+)').firstMatch(setCookie);
  if (match == null) {
    throw StateError(
      'apiLoginAs($username) succeeded but no jwt cookie in '
      'Set-Cookie: $setCookie',
    );
  }
  return match.group(1)!;
}

/// Logs into the Go API as the admin and returns the JWT cookie value.
/// Thin wrapper over [apiLoginAs] kept for the existing cleanup helpers
/// (`scheduleReceiptCleanup`, `getReceipt`, …) that always act as admin.
Future<String> apiLogin() =>
    apiLoginAs(E2eEnv.adminUsername, E2eEnv.adminPassword);

/// Best-effort DELETE of a receipt. Swallows errors so cleanup failures
/// don't mask test failures. Auth is via the Cookie header, matching how
/// the production API consumes it.
Future<void> deleteReceipt(int receiptId, {required String jwt}) async {
  try {
    await http
        .delete(
          Uri.parse('${E2eEnv.baseUrl}/receipt/$receiptId'),
          headers: {'Cookie': 'jwt=$jwt'},
        )
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // Swallowed on purpose -- best-effort cleanup.
  }
}

/// GETs a receipt by id and returns the parsed JSON body. Used by tests
/// that want to assert server-side state (e.g. "the receipt has 1 image"
/// or "1 item") rather than just the URL we landed on.
Future<Map<String, dynamic>> getReceipt(
  int receiptId, {
  required String jwt,
}) async {
  final res = await http
      .get(
        Uri.parse('${E2eEnv.baseUrl}/receipt/$receiptId'),
        headers: {'Cookie': 'jwt=$jwt'},
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
      'getReceipt($receiptId) failed: HTTP ${res.statusCode}: ${res.body}',
    );
  }
  return jsonDecode(res.body) as Map<String, dynamic>;
}

/// Lists all custom fields the admin has access to. Used together with
/// [ensureCustomField] for tests that need a known-name field present.
///
/// The endpoint is `POST /api/customField/getPagedCustomFields`. The
/// API rejects `orderBy: "id"` (server-side bug -- HTTP 500 "Error
/// getting custom fields"); `name` and `type` work, so we order by
/// `name`.
Future<List<Map<String, dynamic>>> listCustomFields({
  required String jwt,
  int limit = 100,
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/customField/getPagedCustomFields'),
        headers: {
          'Content-Type': 'application/json',
          'Cookie': 'jwt=$jwt',
        },
        body: jsonEncode({
          'page': 1,
          'pageSize': limit,
          'orderBy': 'name',
          'sortDirection': 'asc',
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
      'listCustomFields failed: HTTP ${res.statusCode}: ${res.body}',
    );
  }
  final body = jsonDecode(res.body) as Map<String, dynamic>;
  return ((body['data'] as List?) ?? const [])
      .cast<Map<String, dynamic>>();
}

/// Creates a custom field via `POST /api/customField/`. Returns the
/// created field as parsed JSON. [type] is one of TEXT, DATE, SELECT,
/// CURRENCY, BOOLEAN.
Future<Map<String, dynamic>> createCustomField({
  required String jwt,
  required String name,
  required String type,
  String? description,
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/customField/'),
        headers: {
          'Content-Type': 'application/json',
          'Cookie': 'jwt=$jwt',
        },
        body: jsonEncode({
          'name': name,
          'type': type,
          if (description != null) 'description': description,
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
      'createCustomField($name) failed: HTTP ${res.statusCode}: ${res.body}',
    );
  }
  return jsonDecode(res.body) as Map<String, dynamic>;
}

/// Best-effort `DELETE /api/customField/{id}`. Swallows errors so a cleanup
/// failure doesn't mask the test result, matching [deleteReceipt].
///
/// Deleting a custom field destroys every value stored against it, so a spec
/// that seeded receipts holding those values must delete the receipts FIRST.
Future<void> deleteCustomField(int customFieldId, {required String jwt}) async {
  try {
    await http
        .delete(
          Uri.parse('${E2eEnv.baseUrl}/customField/$customFieldId'),
          headers: {'Cookie': 'jwt=$jwt'},
        )
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // Swallowed on purpose -- best-effort cleanup.
  }
}

/// Idempotent: returns the existing custom field with [name] if one
/// exists, otherwise creates one with [type]. Lets tests provision their
/// own fixtures instead of relying on hand-seeded data on the demo
/// backend. Once created the field persists across runs and subsequent
/// calls are pure list-and-filter.
Future<Map<String, dynamic>> ensureCustomField({
  required String jwt,
  required String name,
  required String type,
}) async {
  final existing = await listCustomFields(jwt: jwt);
  for (final f in existing) {
    if (f['name'] == name) return f;
  }
  return createCustomField(jwt: jwt, name: name, type: type);
}

/// Lists the latest [limit] receipts in [groupId] (newest first by id).
/// Used for "exactly one receipt with this name" assertions in flows
/// that need to detect duplicates server-side. Filters client-side --
/// the server's filter shape is fiddly and we only care about a small
/// recent window.
Future<List<Map<String, dynamic>>> listReceiptsForGroup(
  int groupId, {
  required String jwt,
  int limit = 50,
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/receipt/group/$groupId'),
        headers: {
          'Content-Type': 'application/json',
          'Cookie': 'jwt=$jwt',
        },
        body: jsonEncode({
          'page': 1,
          'pageSize': limit,
          // No `orderBy` -- the API rejects "id" with HTTP 500
          // ("Error getting receipts"), and the default ordering
          // returns newest first which is what we want anyway.
          'sortDirection': 'desc',
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
      'listReceiptsForGroup($groupId) failed: '
      'HTTP ${res.statusCode}: ${res.body}',
    );
  }
  final body = jsonDecode(res.body) as Map<String, dynamic>;
  return ((body['data'] as List?) ?? const [])
      .cast<Map<String, dynamic>>();
}

/// JSON + admin-cookie headers for the write endpoints below.
Map<String, String> jsonAuthHeaders(String jwt) => {
      'Content-Type': 'application/json',
      'Cookie': 'jwt=$jwt',
    };

/// GETs the global system settings. Shared by the fixtures that flip a
/// server-wide flag for the duration of a test (`feature_flags.dart`,
/// `login_qr_fixtures.dart`) -- both need the current object to capture the
/// original value AND to build the restore payload.
Future<Map<String, dynamic>> getSystemSettings(String jwt) async {
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

/// PUTs the global system settings. The endpoint is an UPSERT with required
/// fields (currency, taskConcurrency, ...), so callers must pass a full
/// settings object -- read one with [getSystemSettings] and patch the keys
/// they care about rather than sending a partial body.
Future<void> putSystemSettings(
  String jwt,
  Map<String, dynamic> settings,
) async {
  final res = await http
      .put(
        Uri.parse('${E2eEnv.baseUrl}/systemSettings/'),
        headers: jsonAuthHeaders(jwt),
        body: jsonEncode(settings),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
        'PUT systemSettings failed: HTTP ${res.statusCode}: ${res.body}');
  }
}
