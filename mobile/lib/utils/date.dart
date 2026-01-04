import 'package:intl/intl.dart';

const defaultDateFormat = "MM/dd/yyyy";

const defaultLongDateFormat = "MMM dd, yyyy, hh:mm:ss a";

const formFormat = "yyyy-dd-MM";

const zuluDateFormat = "yyyy-MM-dd'T'HH:mm:ss'Z'";

String formatDate(String dateFormat, DateTime date) {
  var formatter = DateFormat(dateFormat);
  return formatter.format(date.toLocal());
}
