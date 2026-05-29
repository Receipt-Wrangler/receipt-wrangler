import 'package:dio/dio.dart';
import 'package:flutter/material.dart';

void showSuccessSnackbar(BuildContext context, String message,
    {SnackBarAction? action}) {
  ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      content: Text(message), backgroundColor: Colors.green, action: action));
}

void showErrorSnackbar(BuildContext context, String message,
    {SnackBarAction? action}) {
  ScaffoldMessenger.of(context).showSnackBar(SnackBar(
    content: Text(message),
    backgroundColor: Colors.red,
    action: action,
  ));
}

void showApiErrorSnackbar(BuildContext context, DioException error) {
  String? message;
  final data = error.response?.data;
  if (data is Map) {
    final raw = data['errorMsg'];
    if (raw is String && raw.isNotEmpty) {
      message = raw;
    }
  }
  ScaffoldMessenger.of(context).showSnackBar(SnackBar(
    content: Text(message ?? 'An error occurred'),
    backgroundColor: Colors.red,
  ));
}
