import 'package:intl/intl.dart';

const defaultDateFormat = "MM/dd/yyyy";

const defaultLongDateFormat = "MMM dd, yyyy, hh:mm:ss a";

const formFormat = "yyyy-dd-MM";

const zuluDateFormat = "yyyy-MM-dd'T'HH:mm:ss'Z'";

String formatDate(String dateFormat, DateTime date) {
  var formatter = DateFormat(dateFormat);
  return formatter.format(date.toLocal());
}

/// Midnight at the start of [date]'s day, in its own zone.
DateTime startOfDay(DateTime date) {
  return DateTime(date.year, date.month, date.day);
}

/// The last representable instant of [date]'s day, in its own zone.
///
/// The receipt date columns are datetimes, so a range whose upper bound is that
/// day's midnight excludes everything actually recorded on it.
DateTime endOfDay(DateTime date) {
  return DateTime(date.year, date.month, date.day, 23, 59, 59, 999);
}
