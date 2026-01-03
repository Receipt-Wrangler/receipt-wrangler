import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

class UserAvatar extends StatefulWidget {
  const UserAvatar({super.key, this.userId});

  final String? userId;

  @override
  State<UserAvatar> createState() => _UserAvatar();
}

class _UserAvatar extends State<UserAvatar> {
  String _getAvatarText() {
    var authModel = Provider.of<AuthModel>(context, listen: true);
    var userModel = Provider.of<UserModel>(context, listen: true);

    if (widget.userId?.isNotEmpty == true) {
      var user = userModel.getUserById(widget.userId!);
      return user!.displayName[0].toUpperCase();
    } else {
      return authModel.claims!.displayName[0]!.toUpperCase();
    }
  }

  @override
  Widget build(BuildContext context) {
    return CircleAvatar(
      child: Text(_getAvatarText()),
    );
  }
}
