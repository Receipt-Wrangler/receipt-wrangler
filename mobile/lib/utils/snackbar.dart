import 'package:dio/dio.dart';
import 'package:flutter/material.dart';

void showSuccessSnackbar(BuildContext context, String message,
    {SnackBarAction? action}) {
  ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      content: Text(message), backgroundColor: Colors.green, action: action));
}

/// A neutral, non-error notice — used where the app quietly changes course and
/// wants to say why (e.g. falling back to the gallery when the camera is off).
/// Left on the theme's default background so it doesn't read as a failure.
void showInfoSnackbar(BuildContext context, String message,
    {SnackBarAction? action}) {
  ScaffoldMessenger.of(context)
      .showSnackBar(SnackBar(content: Text(message), action: action));
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
