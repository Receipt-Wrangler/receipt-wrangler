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

  group('extractDeepLinkServerUrl', () {
    test('extracts the inner url from a valid app/setup deep link', () {
      expect(
        extractDeepLinkServerUrl(
          'https://receiptwrangler.io/app/setup'
          '#url=https%3A%2F%2Fdemo.receiptwrangler.io%2Fapi',
        ),
        'https://demo.receiptwrangler.io/api',
      );
    });

    test('decodes a percent-encoded inner url (http LAN + port)', () {
      expect(
        extractDeepLinkServerUrl(
          'https://receiptwrangler.io/app/setup'
          '#url=http%3A%2F%2F192.168.1.5%3A8081%2Fapi',
        ),
        'http://192.168.1.5:8081/api',
      );
    });

    test('returns null when the url fragment param is missing', () {
      expect(
        extractDeepLinkServerUrl('https://receiptwrangler.io/app/setup'),
        isNull,
      );
      expect(
        extractDeepLinkServerUrl(
            'https://receiptwrangler.io/app/setup#other=1'),
        isNull,
      );
    });

    test('returns null for the wrong path', () {
      expect(
        extractDeepLinkServerUrl(
          'https://receiptwrangler.io/app/other'
          '#url=https%3A%2F%2Fdemo.receiptwrangler.io%2Fapi',
        ),
        isNull,
      );
    });

    test('returns null for the wrong host', () {
      expect(
        extractDeepLinkServerUrl(
          'https://evil.example.com/app/setup'
          '#url=https%3A%2F%2Fdemo.receiptwrangler.io%2Fapi',
        ),
        isNull,
      );
    });

    test('rejects a non-http(s) inner url via normalizeServerUrl', () {
      expect(
        extractDeepLinkServerUrl(
          'https://receiptwrangler.io/app/setup#url=ftp%3A%2F%2Fx.io',
        ),
        isNull,
      );
      expect(
        extractDeepLinkServerUrl(
          'https://receiptwrangler.io/app/setup#url=not-a-url',
        ),
        isNull,
      );
    });

    test('returns null for arbitrary non-deep-link text', () {
      expect(extractDeepLinkServerUrl('hello world'), isNull);
      expect(
        extractDeepLinkServerUrl('https://demo.receiptwrangler.io/api'),
        isNull,
      );
    });
  });
}
