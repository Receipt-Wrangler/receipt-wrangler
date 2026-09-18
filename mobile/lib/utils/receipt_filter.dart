import 'package:built_value/json_object.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/utils/currency.dart';
import 'package:receipt_wrangler_mobile/utils/date.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

/// One authored filter condition: an operation plus the value it applies to.
///
/// [value] holds the **display** objects rather than the wire payload -- whole
/// `Category` / `Tag` / `Group` / `UserView` / `ReceiptStatus` instances for the
/// list buckets, `DateTime`s for dates -- so the condition card, the editor's
/// chips and [buildReceiptPagedRequestFilter] all read one thing. Ids and wire
/// strings are derived at encode time.
///
/// By field type:
///
///  * text   -> `String`
///  * number -> `double`, or `[double, double]` for `BETWEEN`
///  * date   -> `DateTime`, `[DateTime, DateTime]` for `BETWEEN`, `null` for
///              `WITHIN_CURRENT_MONTH` (the API pins that range itself)
///  * list / users -> `List` of the selected objects
class ReceiptFilterCondition {
  const ReceiptFilterCondition({required this.operation, this.value});

  final api.FilterOperation operation;

  final dynamic value;

  ReceiptFilterCondition copyWith({
    api.FilterOperation? operation,
    dynamic value,
    bool clearValue = false,
  }) {
    return ReceiptFilterCondition(
      operation: operation ?? this.operation,
      value: clearValue ? null : (value ?? this.value),
    );
  }
}

/// Writes [fieldObject] onto the `ReceiptPagedRequestFilter` property named
/// [key], ignoring a key nothing declares.
///
/// The generated filter's ten properties are `JsonObject?` rather than a typed
/// `PagedRequestField`, so there is no way to set one by name -- hence the
/// switch. It lives here so the two producers ([buildReceiptPagedRequestFilter]
/// and `dashboardConfigurationToFilter`) cannot drift over which keys exist.
void setReceiptFilterField(
  api.ReceiptPagedRequestFilterBuilder builder,
  String key,
  JsonObject fieldObject,
) {
  switch (key) {
    case "date":
      builder.date = fieldObject;
      break;
    case "amount":
      builder.amount = fieldObject;
      break;
    case "name":
      builder.name = fieldObject;
      break;
    case "paidBy":
      builder.paidBy = fieldObject;
      break;
    case "categories":
      builder.categories = fieldObject;
      break;
    case "tags":
      builder.tags = fieldObject;
      break;
    case "status":
      builder.status = fieldObject;
      break;
    case "group":
      builder.group = fieldObject;
      break;
    case "resolvedDate":
      builder.resolvedDate = fieldObject;
      break;
    case "createdAt":
      builder.createdAt = fieldObject;
      break;
  }
}

/// The string an operation is sent as.
///
/// Every operation a field offers is named identically in Dart and on the wire,
/// but `FilterOperation.empty` is not -- its wire name is `""` while its Dart
/// name is `"empty"`. It is never offered, and a condition carrying it would be
/// invalid anyway; mapping explicitly means a stray one degrades to the empty
/// operation the API already ignores rather than to a string it does not know.
String filterOperationWireName(api.FilterOperation operation) {
  return operation == api.FilterOperation.empty ? "" : operation.name;
}

/// Encodes authored [conditions] into the filter the receipts endpoint takes.
///
/// A field absent from [conditions] is left unset, which the API reads as "do
/// not narrow on this" -- an unset field is not the same as one carrying an
/// empty value.
api.ReceiptPagedRequestFilterBuilder buildReceiptPagedRequestFilter(
  Map<String, ReceiptFilterCondition> conditions,
) {
  final builder = api.ReceiptPagedRequestFilterBuilder();

  conditions.forEach((key, condition) {
    final field = receiptFilterFieldByKey(key);
    if (field == null) {
      return;
    }

    // A half-authored condition is dropped rather than sent, because on the API
    // an EMPTY field is not the same as an ABSENT one. `initReceiptFilterValues`
    // coerces a null date value to "" and a null amount to 0, and the query
    // builder then runs `date = ''` (matches nothing) or `amount = 0` (matches
    // the wrong rows) -- both silent. An absent field is genuinely inert.
    if (!isReceiptFilterConditionValid(field, condition)) {
      return;
    }

    setReceiptFilterField(
      builder,
      key,
      JsonObject({
        "operation": filterOperationWireName(condition.operation),
        "value": _encodeValue(field, condition),
      }),
    );
  });

  return builder;
}

/// The wire value for [condition], shaped the way the Go query builder
/// type-asserts it: a string for text and single dates, a number for a single
/// amount, a two-element array for either `BETWEEN`, and an array of ids (or
/// status wire strings) for the list buckets.
dynamic _encodeValue(
    ReceiptFilterField field, ReceiptFilterCondition condition) {
  switch (field.type) {
    case ReceiptFilterFieldType.text:
      return condition.value as String? ?? "";

    case ReceiptFilterFieldType.number:
      if (condition.operation == api.FilterOperation.BETWEEN) {
        final pair = _numberPair(condition.value);
        if (pair == null) {
          return null;
        }
        // Low first: the query is `>= a AND <= b`, so a reversed pair would
        // match nothing rather than the range the user drew.
        return [
          pair[0] <= pair[1] ? pair[0] : pair[1],
          pair[0] <= pair[1] ? pair[1] : pair[0],
        ];
      }
      return (condition.value as num?)?.toDouble();

    case ReceiptFilterFieldType.date:
      if (condition.operation == api.FilterOperation.WITHIN_CURRENT_MONTH) {
        // The API pins this to month-start through today and ignores any value.
        return null;
      }
      if (condition.operation == api.FilterOperation.BETWEEN) {
        final pair = _datePair(condition.value);
        if (pair == null) {
          return null;
        }
        final low = pair[0].isAfter(pair[1]) ? pair[1] : pair[0];
        final high = pair[0].isAfter(pair[1]) ? pair[0] : pair[1];
        // These columns are datetimes, so a bare `<= 2026-09-16T00:00:00Z`
        // would exclude everything recorded on the last day the user picked.
        return [
          formatDate(zuluDateFormat, startOfDay(low)),
          formatDate(zuluDateFormat, endOfDay(high)),
        ];
      }
      final date = condition.value as DateTime?;
      return date == null ? null : formatDate(zuluDateFormat, date);

    case ReceiptFilterFieldType.list:
    case ReceiptFilterFieldType.users:
      return _selectionValues(condition.value);
  }
}

List<double>? _numberPair(dynamic value) {
  if (value is! List || value.length != 2) {
    return null;
  }
  final low = value[0];
  final high = value[1];
  if (low is! num || high is! num) {
    return null;
  }
  return [low.toDouble(), high.toDouble()];
}

List<DateTime>? _datePair(dynamic value) {
  if (value is! List || value.length != 2) {
    return null;
  }
  final low = value[0];
  final high = value[1];
  if (low is! DateTime || high is! DateTime) {
    return null;
  }
  return [low, high];
}

/// Selected objects -> the ids (or status wire strings) the API filters on.
List<dynamic> _selectionValues(dynamic value) {
  if (value is! List) {
    return const [];
  }

  return value
      .map(receiptFilterOptionValue)
      .where((entry) => entry != null)
      .toList();
}

/// The wire value identifying one selected option.
dynamic receiptFilterOptionValue(dynamic option) {
  if (option is api.Category) return option.id;
  if (option is api.Tag) return option.id;
  if (option is api.Group) return option.id;
  if (option is api.UserView) return option.id;
  if (option is api.ReceiptStatus) return option.name;
  return null;
}

/// How one selected option reads. Status goes through [receiptStatusLabel], the
/// single owner of status presentation, rather than a second label list.
String receiptFilterOptionLabel(dynamic option) {
  if (option is api.Category) return option.name ?? "";
  if (option is api.Tag) return option.name;
  if (option is api.Group) return option.name;
  if (option is api.UserView) return option.displayName;
  if (option is api.ReceiptStatus) return receiptStatusLabel(option);
  return "";
}

/// Renders a stored (USD) amount in the install's configured currency.
///
/// Filter amounts are held in USD, the same unit `AmountField`'s
/// `valueTransformer` produces and the API stores, so the display conversion is
/// the same one the receipt list uses.
String _formatAmount(double usdAmount) {
  return exchangeUSDToCustom(usdAmount.toString()).toString();
}

/// The one-line summary of [condition] shown on its card.
String receiptFilterValueLabel(
    ReceiptFilterField field, ReceiptFilterCondition condition) {
  switch (field.type) {
    case ReceiptFilterFieldType.text:
      final text = (condition.value as String? ?? "").trim();
      return text.isEmpty ? "Any" : text;

    case ReceiptFilterFieldType.number:
      if (condition.operation == api.FilterOperation.BETWEEN) {
        final pair = _numberPair(condition.value);
        if (pair == null) return "Any";
        return "${_formatAmount(pair[0])} - ${_formatAmount(pair[1])}";
      }
      final amount = (condition.value as num?)?.toDouble();
      return amount == null ? "Any" : _formatAmount(amount);

    case ReceiptFilterFieldType.date:
      if (condition.operation == api.FilterOperation.WITHIN_CURRENT_MONTH) {
        return "This month, through today";
      }
      if (condition.operation == api.FilterOperation.BETWEEN) {
        final pair = _datePair(condition.value);
        if (pair == null) return "Any";
        return "${formatDate(defaultDateFormat, pair[0])} - "
            "${formatDate(defaultDateFormat, pair[1])}";
      }
      final date = condition.value as DateTime?;
      return date == null ? "Any" : formatDate(defaultDateFormat, date);

    case ReceiptFilterFieldType.list:
    case ReceiptFilterFieldType.users:
      final selection = condition.value as List? ?? const [];
      if (selection.isEmpty) return "Nothing selected";
      return selection.map(receiptFilterOptionLabel).join(", ");
  }
}

/// Whether [condition] narrows anything, i.e. whether the editor may save it.
///
/// A zero (or negative) amount is deliberately **valid**: nothing defaults to
/// zero, and the API applies an `amount EQUALS 0`, so rejecting it would leave
/// the user unable to author a filter the backend honours.
bool isReceiptFilterConditionValid(
    ReceiptFilterField field, ReceiptFilterCondition condition) {
  switch (field.type) {
    case ReceiptFilterFieldType.text:
      return (condition.value as String? ?? "").trim().isNotEmpty;

    case ReceiptFilterFieldType.number:
      if (condition.operation == api.FilterOperation.BETWEEN) {
        return _numberPair(condition.value) != null;
      }
      return condition.value is num;

    case ReceiptFilterFieldType.date:
      if (condition.operation == api.FilterOperation.WITHIN_CURRENT_MONTH) {
        return true;
      }
      if (condition.operation == api.FilterOperation.BETWEEN) {
        return _datePair(condition.value) != null;
      }
      return condition.value is DateTime;

    case ReceiptFilterFieldType.list:
    case ReceiptFilterFieldType.users:
      final selection = condition.value;
      return selection is List && selection.isNotEmpty;
  }
}

/// The value a freshly-picked [operation] starts from.
///
/// Switching between a single value and a `BETWEEN` pair changes the *shape* the
/// encoder expects, so a value that no longer fits is dropped rather than
/// carried over -- otherwise a stale single would be encoded as a range.
dynamic defaultReceiptFilterValue(
  ReceiptFilterField field,
  api.FilterOperation operation,
  dynamic previousValue,
) {
  switch (field.type) {
    case ReceiptFilterFieldType.text:
      return previousValue is String ? previousValue : "";

    case ReceiptFilterFieldType.number:
      if (operation == api.FilterOperation.BETWEEN) {
        return _numberPair(previousValue);
      }
      return previousValue is num ? previousValue.toDouble() : null;

    case ReceiptFilterFieldType.date:
      if (operation == api.FilterOperation.WITHIN_CURRENT_MONTH) {
        return null;
      }
      if (operation == api.FilterOperation.BETWEEN) {
        return _datePair(previousValue);
      }
      return previousValue is DateTime ? previousValue : null;

    case ReceiptFilterFieldType.list:
    case ReceiptFilterFieldType.users:
      return previousValue is List ? previousValue : <dynamic>[];
  }
}
