import 'package:flutter/material.dart';

class TotalDisplayWidget extends StatelessWidget {
  final String value;
  final bool isValid;
  final String? errorMessage;

  const TotalDisplayWidget({
    super.key,
    required this.value,
    this.isValid = true,
    this.errorMessage,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: isValid ? Colors.green.shade50 : Colors.red.shade50,
        border: Border.all(
          color: isValid ? Colors.green.shade300 : Colors.red.shade300,
          width: 1,
        ),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Icon(
            isValid ? Icons.check_circle : Icons.error,
            color: isValid ? Colors.green.shade600 : Colors.red.shade600,
            size: 20,
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              "Total: $value${errorMessage != null ? ' ($errorMessage)' : ''}",
              style: TextStyle(
                color: isValid ? Colors.green.shade800 : Colors.red.shade800,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }
}