import 'package:built_value/json_object.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/receipt_filter_test_helpers.dart';
import '../helpers/receipt_form_test_helpers.dart';
import '../helpers/widget_test_helpers.dart';

/// Encoding is where a mistake costs the most: the Go query builder
/// type-asserts every value with no comma-ok (`.(string)`, `.(float64)`,
/// `.([]interface{})`), so a wrong shape is a 500 on the receipts list rather
/// than a condition that quietly does nothing. Hence a case per
/// (field type, operation) pair, asserted against the real serializer.
void main() {
  setUpAll(registerCustomCurrencyForTests);

  ReceiptFilterField field(String key) => receiptFilterFieldByKey(key)!;

  group("text", () {
    test("CONTAINS sends the raw string", () {
      expect(
          serializedField({
            "name": const ReceiptFilterCondition(
                operation: api.FilterOperation.CONTAINS, value: "whole foods")
          }, "name"),
          {"operation": "CONTAINS", "value": "whole foods"});
    });

    test("EQUALS sends the raw string", () {
      expect(
          serializedField({
            "name": const ReceiptFilterCondition(
                operation: api.FilterOperation.EQUALS, value: "Costco")
          }, "name"),
          {"operation": "EQUALS", "value": "Costco"});
    });
  });

  group("number", () {
    test("a single amount is sent as a NUMBER, never a string", () {
      // AmountField's valueTransformer hands back "12.34"; the Go side does
      // `Value.(float64)` and panics on a string.
      final encoded = serializedField({
        "amount": const ReceiptFilterCondition(
            operation: api.FilterOperation.GREATER_THAN, value: 12.34)
      }, "amount");

      expect(encoded!["value"], isA<num>());
      expect(encoded["value"], 12.34);
      expect(encoded["operation"], "GREATER_THAN");
    });

    test("BETWEEN sends a two-element array", () {
      expect(
          serializedField({
            "amount": const ReceiptFilterCondition(
                operation: api.FilterOperation.BETWEEN, value: [10.0, 25.5])
          }, "amount"),
          {
            "operation": "BETWEEN",
            "value": [10.0, 25.5]
          });
    });

    test("a reversed BETWEEN pair is normalised low to high", () {
      // The query is `>= a AND <= b`, so a reversed pair matches nothing rather
      // than the range the user drew.
      final encoded = serializedField({
        "amount": const ReceiptFilterCondition(
            operation: api.FilterOperation.BETWEEN, value: [90.0, 10.0])
      }, "amount");

      expect(encoded!["value"], [10.0, 90.0]);
    });

    test("an int amount is accepted and sent as a number", () {
      final encoded = serializedField({
        "amount": const ReceiptFilterCondition(
            operation: api.FilterOperation.EQUALS, value: 7)
      }, "amount");

      expect(encoded!["value"], 7);
    });

    test("zero is a real filter and is sent", () {
      // Nothing defaults to zero, and the API applies `amount = 0`, so dropping
      // it would leave a condition the user authored with no effect.
      expect(
          serializedField({
            "amount": const ReceiptFilterCondition(
                operation: api.FilterOperation.EQUALS, value: 0.0)
          }, "amount"),
          {"operation": "EQUALS", "value": 0.0});
    });
  });

  group("date", () {
    for (final key in ["date", "resolvedDate", "createdAt"]) {
      test("$key EQUALS sends one zulu string", () {
        expect(
            serializedField({
              key: ReceiptFilterCondition(
                  operation: api.FilterOperation.EQUALS,
                  value: DateTime(2026, 9, 16, 14, 30))
            }, key),
            {"operation": "EQUALS", "value": "2026-09-16T14:30:00Z"});
      });

      test("$key BETWEEN spans start-of-day to end-of-day", () {
        // These are datetime columns, so a bare `<= 2026-09-18T00:00:00Z` upper
        // bound excludes everything actually recorded on the last day picked.
        expect(
            serializedField({
              key: ReceiptFilterCondition(
                  operation: api.FilterOperation.BETWEEN,
                  value: [DateTime(2026, 9, 16, 9), DateTime(2026, 9, 18, 17)])
            }, key),
            {
              "operation": "BETWEEN",
              "value": ["2026-09-16T00:00:00Z", "2026-09-18T23:59:59Z"]
            });
      });

      test("$key WITHIN_CURRENT_MONTH carries no value", () {
        // The API pins the range to month start through today and ignores
        // whatever it is handed.
        expect(
            serializedField({
              key: const ReceiptFilterCondition(
                  operation: api.FilterOperation.WITHIN_CURRENT_MONTH)
            }, key),
            {"operation": "WITHIN_CURRENT_MONTH", "value": null});
      });
    }

    test("a reversed BETWEEN range is normalised earliest to latest", () {
      final encoded = serializedField({
        "date": ReceiptFilterCondition(
            operation: api.FilterOperation.BETWEEN,
            value: [DateTime(2026, 9, 20), DateTime(2026, 9, 1)])
      }, "date");

      expect(encoded!["value"],
          ["2026-09-01T00:00:00Z", "2026-09-20T23:59:59Z"]);
    });
  });

  group("list and users", () {
    test("categories send ids, not objects, in selection order", () {
      expect(
          serializedField({
            "categories": ReceiptFilterCondition(
                operation: api.FilterOperation.CONTAINS,
                value: [buildCategory(7, "Dining"), buildCategory(3, "Fuel")])
          }, "categories"),
          {
            "operation": "CONTAINS",
            "value": [7, 3]
          });
    });

    test("tags send ids", () {
      expect(
          serializedField({
            "tags": ReceiptFilterCondition(
                operation: api.FilterOperation.CONTAINS,
                value: [buildTag(2, "Reimbursable")])
          }, "tags")!["value"],
          [2]);
    });

    test("paidBy sends user ids", () {
      expect(
          serializedField({
            "paidBy": ReceiptFilterCondition(
                operation: api.FilterOperation.CONTAINS,
                value: [buildUserView(id: 4, displayName: "Dana Kim")])
          }, "paidBy")!["value"],
          [4]);
    });

    test("status sends WIRE STRINGS, not the labels the UI shows", () {
      expect(
          serializedField({
            "status": const ReceiptFilterCondition(
                operation: api.FilterOperation.CONTAINS,
                value: [api.ReceiptStatus.OPEN, api.ReceiptStatus.NEEDS_ATTENTION])
          }, "status"),
          {
            "operation": "CONTAINS",
            "value": ["OPEN", "NEEDS_ATTENTION"]
          });
    });

    test("a list value is always an array, even for a single selection", () {
      // The Go side does `Value.([]interface{})` for all five list fields.
      final encoded = serializedField({
        "categories": ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS,
            value: [buildCategory(1, "Groceries")])
      }, "categories");

      expect(encoded!["value"], isA<List<dynamic>>());
    });
  });

  group("what is NOT sent", () {
    test("no conditions leaves every field unset", () {
      // An unfiltered list must produce the request it always did.
      expect(serializeFilter({}), isEmpty);
    });

    test("a field with no condition is left unset", () {
      final encoded = serializeFilter({
        "name": const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS, value: "Costco")
      });

      expect(encoded.keys, ["name"]);
    });

    test("an unknown key is ignored", () {
      expect(
          serializeFilter({
            "nonsense": const ReceiptFilterCondition(
                operation: api.FilterOperation.EQUALS, value: "x")
          }),
          isEmpty);
    });

    group("a half-authored condition is dropped rather than sent", () {
      // An EMPTY field is not an ABSENT one: initReceiptFilterValues coerces a
      // null date value to "" and a null amount to 0, and the query builder then
      // runs `date = ''` (matches nothing) or `amount = 0` (matches the wrong
      // rows). Both are silent.
      test("a date with no date", () {
        expect(
            serializeFilter({
              "date": const ReceiptFilterCondition(
                  operation: api.FilterOperation.EQUALS)
            }),
            isEmpty);
      });

      test("a half-filled date range", () {
        expect(
            serializeFilter({
              "date": ReceiptFilterCondition(
                  operation: api.FilterOperation.BETWEEN,
                  value: [DateTime(2026, 9, 1), null])
            }),
            isEmpty);
      });

      test("an amount with no amount", () {
        expect(
            serializeFilter({
              "amount": const ReceiptFilterCondition(
                  operation: api.FilterOperation.LESS_THAN)
            }),
            isEmpty);
      });

      test("a half-filled amount range", () {
        expect(
            serializeFilter({
              "amount": const ReceiptFilterCondition(
                  operation: api.FilterOperation.BETWEEN, value: [10.0, null])
            }),
            isEmpty);
      });

      test("blank text", () {
        expect(
            serializeFilter({
              "name": const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS, value: "   ")
            }),
            isEmpty);
      });

      test("an empty selection", () {
        expect(
            serializeFilter({
              "tags": const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS, value: [])
            }),
            isEmpty);
      });
    });
  });

  group("setReceiptFilterField", () {
    test("sets each of the ten slots", () {
      for (final declared in receiptFilterFields) {
        final builder = api.ReceiptPagedRequestFilterBuilder();
        setReceiptFilterField(builder, declared.key, JsonObject({"a": 1}));

        final serialized = Map<String, dynamic>.from(
            api.standardSerializers.serializeWith(
                api.ReceiptPagedRequestFilter.serializer, builder.build()) as Map);

        expect(serialized.keys, [declared.key]);
      }
    });

    test("ignores a key nothing declares", () {
      final builder = api.ReceiptPagedRequestFilterBuilder();
      setReceiptFilterField(builder, "nope", JsonObject({"a": 1}));

      expect(
          api.standardSerializers.serializeWith(
              api.ReceiptPagedRequestFilter.serializer, builder.build()),
          isEmpty);
    });
  });

  group("filterOperationWireName", () {
    test("every real operation is its own name", () {
      for (final operation in filterOperationLabels.keys) {
        expect(filterOperationWireName(operation), operation.name);
      }
    });

    test("the empty operation maps to the empty string, not 'empty'", () {
      expect(filterOperationWireName(api.FilterOperation.empty), "");
    });
  });

  group("isReceiptFilterConditionValid", () {
    test("text needs non-whitespace", () {
      expect(
          isReceiptFilterConditionValid(
              field("name"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS, value: "  ")),
          isFalse);
      expect(
          isReceiptFilterConditionValid(
              field("name"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS, value: " x ")),
          isTrue);
    });

    test("a number needs a number, and zero counts", () {
      expect(
          isReceiptFilterConditionValid(
              field("amount"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.EQUALS)),
          isFalse);
      expect(
          isReceiptFilterConditionValid(
              field("amount"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.EQUALS, value: 0.0)),
          isTrue);
      expect(
          isReceiptFilterConditionValid(
              field("amount"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.EQUALS, value: -5.0)),
          isTrue);
    });

    test("a BETWEEN needs both ends", () {
      expect(
          isReceiptFilterConditionValid(
              field("amount"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.BETWEEN, value: [1.0, null])),
          isFalse);
      expect(
          isReceiptFilterConditionValid(
              field("amount"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.BETWEEN, value: [1.0, 2.0])),
          isTrue);
    });

    test("WITHIN_CURRENT_MONTH is valid with no value at all", () {
      expect(
          isReceiptFilterConditionValid(
              field("date"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.WITHIN_CURRENT_MONTH)),
          isTrue);
    });

    test("a selection needs at least one entry", () {
      expect(
          isReceiptFilterConditionValid(
              field("tags"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS, value: [])),
          isFalse);
      expect(
          isReceiptFilterConditionValid(
              field("tags"),
              ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS,
                  value: [buildTag(1, "Tax")])),
          isTrue);
    });
  });

  group("receiptFilterValueLabel", () {
    test("text reads as itself", () {
      expect(
          receiptFilterValueLabel(
              field("name"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS, value: "Costco")),
          "Costco");
    });

    test("a date reads in the app's display format", () {
      expect(
          receiptFilterValueLabel(
              field("date"),
              ReceiptFilterCondition(
                  operation: api.FilterOperation.EQUALS,
                  value: DateTime(2026, 9, 16))),
          "09/16/2026");
    });

    test("a date range reads as both ends", () {
      expect(
          receiptFilterValueLabel(
              field("date"),
              ReceiptFilterCondition(
                  operation: api.FilterOperation.BETWEEN,
                  value: [DateTime(2026, 9, 1), DateTime(2026, 9, 30)])),
          "09/01/2026 - 09/30/2026");
    });

    test("WITHIN_CURRENT_MONTH says what the server will do", () {
      expect(
          receiptFilterValueLabel(
              field("date"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.WITHIN_CURRENT_MONTH)),
          "This month, through today");
    });

    test("a selection is comma joined", () {
      expect(
          receiptFilterValueLabel(
              field("categories"),
              ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS,
                  value: [buildCategory(1, "Dining"), buildCategory(2, "Fuel")])),
          "Dining, Fuel");
    });

    test("a status reads through the single owner of status labels", () {
      // Not a second hard-coded list: receiptStatusLabel is what the receipts
      // list and the form dropdown already use.
      expect(
          receiptFilterValueLabel(
              field("status"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS,
                  value: [api.ReceiptStatus.NEEDS_ATTENTION])),
          receiptStatusLabel(api.ReceiptStatus.NEEDS_ATTENTION));
    });

    test("an empty selection says so rather than reading blank", () {
      expect(
          receiptFilterValueLabel(
              field("tags"),
              const ReceiptFilterCondition(
                  operation: api.FilterOperation.CONTAINS, value: [])),
          "Nothing selected");
    });
  });

  group("defaultReceiptFilterValue", () {
    test("a single amount does not survive the switch to a range", () {
      // The shapes differ, so carrying it over would encode a lone value as a
      // two-element array.
      expect(
          defaultReceiptFilterValue(
              field("amount"), api.FilterOperation.BETWEEN, 42.0),
          isNull);
    });

    test("a range does not survive the switch back to a single value", () {
      expect(
          defaultReceiptFilterValue(
              field("amount"), api.FilterOperation.EQUALS, [1.0, 2.0]),
          isNull);
    });

    test("a same-shape value is carried over", () {
      expect(
          defaultReceiptFilterValue(
              field("amount"), api.FilterOperation.LESS_THAN, 42.0),
          42.0);
      expect(
          defaultReceiptFilterValue(
              field("amount"), api.FilterOperation.BETWEEN, [1.0, 2.0]),
          [1.0, 2.0]);
    });

    test("WITHIN_CURRENT_MONTH drops whatever date was picked", () {
      expect(
          defaultReceiptFilterValue(field("date"),
              api.FilterOperation.WITHIN_CURRENT_MONTH, DateTime(2026, 1, 1)),
          isNull);
    });

    test("a selection starts empty and survives an operation it shares", () {
      expect(
          defaultReceiptFilterValue(
              field("tags"), api.FilterOperation.CONTAINS, null),
          isEmpty);
    });
  });
}
