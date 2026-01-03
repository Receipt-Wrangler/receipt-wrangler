import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';

class BottomSubmitButton extends StatefulWidget {
  const BottomSubmitButton(
      {super.key, required this.onPressed, this.buttonText, this.disabled});

  final String? buttonText;

  final disabled;

  final void Function() onPressed;

  @override
  State<BottomSubmitButton> createState() => _BottomSubmitButtonState();
}

class _BottomSubmitButtonState extends State<BottomSubmitButton> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      height: 50,
      child: Consumer<LoadingModel>(
        builder: (context, loadingModel, child) {
          return MaterialButton(
            onPressed: (loadingModel.isLoading || (widget.disabled ?? false))
                ? null
                : widget.onPressed,
            color: Theme.of(context).primaryColor,
            child: loadingModel.isLoading
                ? const CircularProgressIndicator(
                    color: Colors.white,
                    strokeWidth: 2,
                  )
                : Text(
                    widget.buttonText ?? "Submit",
                    style: const TextStyle(color: Colors.white),
                  ),
          );
        },
      ),
    );
  }
}
