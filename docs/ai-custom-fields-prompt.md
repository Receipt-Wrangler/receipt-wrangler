# AI population of custom fields — prompt guidance

This change adds the `@customFields` prompt template variable, which serialises
the configured custom fields (id, name, type, description) into the receipt
extraction prompt. To make the AI actually populate custom fields, the prompt
must instruct it to return values. Add a block like the following to the prompt
(it pairs with the `@customFields` variable):

```
Custom fields: from the list below, extract a value for each custom field where
the receipt clearly provides one. Return them in a "customFields" array. Each
entry MUST contain the field's "customFieldId" and exactly ONE typed value
property matching the field type:
- CURRENCY: "currencyValue" as a number (e.g. 4.66). Use for VAT/tax amounts.
- TEXT: "stringValue" as a string.
- DATE: "dateValue" in ISO 8601 UTC.
- BOOLEAN: "booleanValue".
- SELECT: "selectValue" as the option id.
Example: {"customFieldId": 1, "currencyValue": 4.66}. Omit a field entirely if
the receipt does not provide its value (do not guess). Return an empty array if
there are none.

Custom fields to populate: @customFields
```

## Notes
- Custom-field values are NOT run through `UpsertReceiptCommand.Validate`, so the
  prompt must return well-typed values keyed by a valid `customFieldId`. The
  persistence path saves `Receipt.CustomFields []CustomFieldValue` via GORM
  association, same as categories/tags.
- The prompt validator only accepts known template variables; `@customFields`
  becomes valid once this change is deployed.
