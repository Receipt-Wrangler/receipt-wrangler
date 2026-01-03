import 'package:openapi/openapi.dart';

String getActivityTypeDisplay(SystemTaskType type) {
  switch (type) {
    case SystemTaskType.QUICK_SCAN:
      return "Quick Scan";
    case SystemTaskType.RECEIPT_UPDATED:
      return "Updated Receipt";
    case SystemTaskType.RECEIPT_UPLOADED:
      return "Uploaded Receipt";
    case SystemTaskType.EMAIL_UPLOAD:
      return "Email Upload";
    default:
      return "Unknown";
  }
}

String getActivityStatusDisplay(SystemTaskStatus status) {
  switch (status) {
    case SystemTaskStatus.SUCCEEDED:
      return "Succeeded";
    case SystemTaskStatus.FAILED:
      return "Failed";
    default:
      return "Unknown";
  }
}
