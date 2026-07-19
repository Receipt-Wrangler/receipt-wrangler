# openapi.api.ReportApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createReportTemplate**](ReportApi.md#createreporttemplate) | **POST** /report/template | Save a report template
[**deleteReportTemplate**](ReportApi.md#deletereporttemplate) | **DELETE** /report/template/{id} | Delete a report template
[**duplicateReportTemplate**](ReportApi.md#duplicatereporttemplate) | **POST** /report/template/{id}/duplicate | Duplicate a report template
[**generateReport**](ReportApi.md#generatereport) | **POST** /report/generate | Generate a report
[**generateReportFromTemplate**](ReportApi.md#generatereportfromtemplate) | **POST** /report/template/{id}/generate | Generate a report from a saved template
[**getReportTemplate**](ReportApi.md#getreporttemplate) | **GET** /report/template/{id} | Get a report template
[**getReportTemplateOptions**](ReportApi.md#getreporttemplateoptions) | **GET** /report/template/options | Get report template options
[**getReportTemplates**](ReportApi.md#getreporttemplates) | **POST** /report/template/list | Get paged report templates
[**previewReport**](ReportApi.md#previewreport) | **POST** /report/preview | Preview a report
[**renderReportTemplate**](ReportApi.md#renderreporttemplate) | **POST** /report/template/{id}/render | Render a saved template as HTML for the dashboard report widget
[**updateReportTemplate**](ReportApi.md#updatereporttemplate) | **PUT** /report/template/{id} | Update a report template


# **createReportTemplate**
> ReportTemplate createReportTemplate(reportRequestCommand)

Save a report template

Saves the current report configuration as a reusable template. App-scoped (requires app.reports.create) — it persists a configuration and touches no group's receipts.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final ReportRequestCommand reportRequestCommand = ; // ReportRequestCommand | The report builder configuration to save

try {
    final response = api.createReportTemplate(reportRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->createReportTemplate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **reportRequestCommand** | [**ReportRequestCommand**](ReportRequestCommand.md)| The report builder configuration to save | 

### Return type

[**ReportTemplate**](ReportTemplate.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteReportTemplate**
> deleteReportTemplate(id)

Delete a report template

Deletes a saved report template by id. Requires app.reports.delete.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final int id = 56; // int | Id of the report template to delete

try {
    api.deleteReportTemplate(id);
} catch on DioException (e) {
    print('Exception when calling ReportApi->deleteReportTemplate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the report template to delete | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **duplicateReportTemplate**
> ReportTemplate duplicateReportTemplate(id)

Duplicate a report template

Duplicates a saved report template by id, returning the new copy. App-scoped (requires app.reports.duplicate).

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final int id = 56; // int | Id of the report template to duplicate

try {
    final response = api.duplicateReportTemplate(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->duplicateReportTemplate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the report template to duplicate | 

### Return type

[**ReportTemplate**](ReportTemplate.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **generateReport**
> Uint8List generateReport(reportRequestCommand)

Generate a report

Generates a report over the receipts of one or more groups and streams it back as a file, or a zip when several formats are requested. Runs under the caller's group permissions and requires group.reports.read in every covered group.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final ReportRequestCommand reportRequestCommand = ; // ReportRequestCommand | The report builder configuration

try {
    final response = api.generateReport(reportRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->generateReport: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **reportRequestCommand** | [**ReportRequestCommand**](ReportRequestCommand.md)| The report builder configuration | 

### Return type

[**Uint8List**](Uint8List.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **generateReportFromTemplate**
> Uint8List generateReportFromTemplate(id)

Generate a report from a saved template

Generates a saved report template by id and streams the file. Enforces the per-template generate grant (generateAll, or the per-group ceiling plus the matrix), loading the stored configuration server-side.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final int id = 56; // int | Id of the report template to generate

try {
    final response = api.generateReportFromTemplate(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->generateReportFromTemplate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the report template to generate | 

### Return type

[**Uint8List**](Uint8List.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReportTemplate**
> ReportTemplate getReportTemplate(id)

Get a report template

Returns a saved report template by id. Access is resolved per template (readAll, or the per-group ceiling plus the per-template matrix).

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final int id = 56; // int | Id of the report template to get

try {
    final response = api.getReportTemplate(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->getReportTemplate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the report template to get | 

### Return type

[**ReportTemplate**](ReportTemplate.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReportTemplateOptions**
> BuiltList<ReportTemplateOption> getReportTemplateOptions()

Get report template options

Returns every report template as a lightweight {id, name, groupIds} option for the role-form access matrix. Gated on app.roles.read (the admin role editor may not personally hold report permissions), not a report permission.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();

try {
    final response = api.getReportTemplateOptions();
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->getReportTemplateOptions: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BuiltList&lt;ReportTemplateOption&gt;**](ReportTemplateOption.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReportTemplates**
> PagedData getReportTemplates(pagedRequestCommand)

Get paged report templates

Returns a paged, sorted list of saved report templates. App-scoped (requires app.reports.read).

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final PagedRequestCommand pagedRequestCommand = ; // PagedRequestCommand | Paging and sorting data

try {
    final response = api.getReportTemplates(pagedRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->getReportTemplates: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pagedRequestCommand** | [**PagedRequestCommand**](PagedRequestCommand.md)| Paging and sorting data | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **previewReport**
> ReportPreviewResponse previewReport(reportRequestCommand)

Preview a report

Renders the current report configuration as HTML for the builder's live preview, along with the number of receipts it covers. Requires the app-level app.reports.read (or app.reports.readAll) and group.reports.read in every covered group. The preview is a row-capped sample of the engine's own output.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final ReportRequestCommand reportRequestCommand = ; // ReportRequestCommand | The report builder configuration

try {
    final response = api.previewReport(reportRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->previewReport: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **reportRequestCommand** | [**ReportRequestCommand**](ReportRequestCommand.md)| The report builder configuration | 

### Return type

[**ReportPreviewResponse**](ReportPreviewResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **renderReportTemplate**
> ReportPreviewResponse renderReportTemplate(id)

Render a saved template as HTML for the dashboard report widget

Renders a saved report template by id as a self-contained HTML document over the full dataset (unlike /report/preview, which caps the sample) for the dashboard report widget, loading the stored configuration server-side. Re-resolves the caller's access on every call; when the caller may not view the template (or it was deleted) a restricted-notice HTML document is returned at 200 with empty allowedActions, so the widget always has HTML to render. allowedActions carries the actions the caller may perform (drives the widget's download button).

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final int id = 56; // int | Id of the report template to render

try {
    final response = api.renderReportTemplate(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->renderReportTemplate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the report template to render | 

### Return type

[**ReportPreviewResponse**](ReportPreviewResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateReportTemplate**
> ReportTemplate updateReportTemplate(id, reportRequestCommand)

Update a report template

Updates a saved report template in place, replacing its name and stored configuration. App-scoped (requires app.reports.update).

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReportApi();
final int id = 56; // int | Id of the report template to update
final ReportRequestCommand reportRequestCommand = ; // ReportRequestCommand | The report builder configuration to save over the existing template

try {
    final response = api.updateReportTemplate(id, reportRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReportApi->updateReportTemplate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the report template to update | 
 **reportRequestCommand** | [**ReportRequestCommand**](ReportRequestCommand.md)| The report builder configuration to save over the existing template | 

### Return type

[**ReportTemplate**](ReportTemplate.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

