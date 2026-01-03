import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class DateBlock extends StatefulWidget {
  const DateBlock({
    super.key,
    required this.date,
  });

  final String date;

  @override
  State<DateBlock> createState() => _DateBlock();
}

class _DateBlock extends State<DateBlock> {
  @override
  Widget build(BuildContext context) {
    var date = DateTime.parse(widget.date);
    DateFormat format = DateFormat("MMM d");
    var formattedDate = format.format(date);
    var formattedDateParts = formattedDate.split(" ");

    var dateTextWidgets = formattedDateParts.map((e) {
      return Text(
        e,
        style: const TextStyle(
          fontSize: 14,
          fontWeight: FontWeight.w600,
        ),
      );
    });

    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [...dateTextWidgets],
    );
  }
}
