import 'dart:convert';

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
  var errorObject = jsonDecode(error.response.toString() ?? "{}");
  ScaffoldMessenger.of(context).showSnackBar(SnackBar(
    content: Text(errorObject["errorMsg"] ?? 'An error occurred'),
    backgroundColor: Colors.red,
  ));
}
