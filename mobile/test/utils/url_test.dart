import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/utils/url.dart';

void main() {
  group('normalizeServerUrl', () {
    test('accepts an https URL with a path', () {
      expect(
        normalizeServerUrl('https://demo.receiptwrangler.io/api'),
        'https://demo.receiptwrangler.io/api',
      );
    });

    test('accepts a plain-http LAN URL with a port', () {
      expect(
        normalizeServerUrl('http://192.168.1.5:8081/api'),
        'http://192.168.1.5:8081/api',
      );
    });

    test('trims surrounding whitespace and newlines', () {
      expect(
        normalizeServerUrl('  https://demo.receiptwrangler.io/api\n'),
        'https://demo.receiptwrangler.io/api',
      );
    });

    test('returns null for empty or whitespace-only input', () {
      expect(normalizeServerUrl(''), isNull);
      expect(normalizeServerUrl('   '), isNull);
    });

    test('returns null for a scheme-less host', () {
      expect(normalizeServerUrl('demo.receiptwrangler.io'), isNull);
    });

    test('returns null for a non-http(s) scheme', () {
      expect(normalizeServerUrl('ftp://demo.receiptwrangler.io'), isNull);
      expect(normalizeServerUrl('javascript:alert(1)'), isNull);
    });

    test('returns null when the host is empty', () {
      expect(normalizeServerUrl('http://'), isNull);
    });

    test('returns null for arbitrary non-URL text', () {
      expect(normalizeServerUrl('hello world'), isNull);
    });
  });
}
