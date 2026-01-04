import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';

class CustomFieldModel extends ChangeNotifier {
  List<CustomField> _customFields = [];

  List<CustomField> get customFields => _customFields;

  bool _isLoading = false;

  bool get isLoading => _isLoading;

  void setCustomFields(List<CustomField> customFields) {
    _customFields = customFields;
    notifyListeners();
  }

  Future<void> loadCustomFields() async {
    _isLoading = true;

    try {
      var client = await OpenApiClient.client;
      var pagedRequest = ($PagedRequestCommandBuilder()
            ..page = 1
            ..pageSize = -1
            ..orderBy = "name"
            ..sortDirection = SortDirection.desc)
          .build();

      var response = await client.getCustomFieldApi().getPagedCustomFields(
            pagedRequestCommand: pagedRequest,
          );

      if (response.data?.data != null) {
        var customFields = response.data!.data
            .map((field) => field.anyOf.values[10] as CustomField);
        setCustomFields(customFields.toList());
      }
    } catch (e) {
      print('Error loading custom fields: $e');
      setCustomFields([]);
    } finally {
      _isLoading = false;
    }
  }

  void resetModel() {
    _customFields = [];
    _isLoading = false;
  }
}
