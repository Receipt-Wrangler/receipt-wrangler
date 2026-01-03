import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/auth/set-homeserver-url/screens/set_homeserver_url.dart';

class Home extends StatefulWidget {
  const Home({super.key});

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  @override
  Widget build(BuildContext context) {
    return const SetHomeserverUrl();
  }
}
