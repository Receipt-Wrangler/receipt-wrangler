import 'dart:convert';

import 'package:http/http.dart' as http;

import 'env.dart';

/// Logs into the Go API as the admin and returns the JWT cookie value.
///
/// The API issues auth via `Set-Cookie: jwt=…` (the body's `jwt` field
/// is empty — confirmed against a live demo response). We just parse
/// the cookie out of the Set-Cookie header.
Future<String> apiLogin() async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/login/'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': E2eEnv.adminUsername,
          'password': E2eEnv.adminPassword,
        }),
      )
      .timeout(const Duration(seconds: 10));

  if (res.statusCode != 200) {
    throw StateError(
      'apiLogin failed: HTTP ${res.statusCode}: ${res.body}',
    );
  }
  final setCookie = res.headers['set-cookie'] ?? '';
  final match = RegExp(r'jwt=([^;]+)').firstMatch(setCookie);
  if (match == null) {
    throw StateError(
      'apiLogin succeeded but no jwt cookie in Set-Cookie: $setCookie',
    );
  }
  return match.group(1)!;
}

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
