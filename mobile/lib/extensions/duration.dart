extension StringDurationExtension on String? {
  String toDurationString() {
    if (this == null || this!.isEmpty) {
      return '';
    }

    try {
      final dateTime = DateTime.parse(this!);
      return dateTime.toDurationString();
    } catch (e) {
      return '';
    }
  }
}

extension DateTimeDurationExtension on DateTime? {
  String toDurationString() {
    if (this == null) {
      return '';
    }

    final now = DateTime.now().toUtc();
    final date = this!.toUtc();
    final difference = now.difference(date);

    // Convert to days, hours, minutes, seconds
    final days = difference.inDays;
    final hours = difference.inHours % 24;
    final minutes = difference.inMinutes % 60;
    final seconds = difference.inSeconds % 60;

    // Check if date is invalid or in the future
    if (days < 0 || hours < 0 || minutes < 0) {
      return '';
    }

    // Check if date is today
    final isToday =
        now.year == date.year && now.month == date.month && now.day == date.day;

    if (isToday) {
      if (hours > 0) {
        return '${hours} ${_pluralize(hours, 'hour')} ago';
      }
      if (minutes > 0) {
        return '${minutes} ${_pluralize(minutes, 'minute')} ago';
      }
      if (seconds > 0) {
        return 'just now';
      }
    } else {
      if (days > 0) {
        return '${days} ${_pluralize(days, 'day')} ago';
      }
      if (hours > 0) {
        return '${hours} ${_pluralize(hours, 'hour')} ago';
      }
    }

    return '';
  }
}

String _pluralize(int count, String word) {
  return count == 1 ? word : '${word}s';
}

// Usage examples:
// String dateStr = '2024-01-29T10:00:00Z';
// String duration1 = dateStr.toDurationString();
//
// DateTime dateTime = DateTime.parse('2024-01-29T10:00:00Z');
// String duration2 = dateTime.toDurationString();
