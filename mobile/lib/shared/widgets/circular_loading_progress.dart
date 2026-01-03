import 'package:flutter/material.dart';

class CircularLoadingProgress extends StatefulWidget {
  const CircularLoadingProgress({
    super.key,
  });

  @override
  State<CircularLoadingProgress> createState() => _CircularLoadingProgress();
}

class _CircularLoadingProgress extends State<CircularLoadingProgress> {
  @override
  Widget build(BuildContext context) {
    return const Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: <Widget>[
        CircularProgressIndicator(),
      ],
    );
  }
}
