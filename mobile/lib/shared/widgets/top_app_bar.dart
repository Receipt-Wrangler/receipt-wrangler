import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/user_avatar.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

import '../../client/client.dart';

class TopAppBar extends StatefulWidget implements PreferredSizeWidget {
  const TopAppBar(
      {super.key,
      required this.titleText,
      this.leadingArrowRedirect,
      this.leadingArrowExtra,
      this.onLeadingArrowPressed,
      this.leadingArrowPop,
      this.actions,
      this.hideAvatar,
      this.surfaceTintColor});

  final String titleText;

  final String? leadingArrowRedirect;

  final dynamic leadingArrowExtra;

  final bool? leadingArrowPop;

  final Function? onLeadingArrowPressed;

  final List<Widget>? actions;

  final bool? hideAvatar;

  final Color? surfaceTintColor;

  @override
  State<TopAppBar> createState() => _TopAppBar();

  @override
  Size get preferredSize => AppBar().preferredSize;
}

class _TopAppBar extends State<TopAppBar> {
  late final AuthModel authModel =
      Provider.of<AuthModel>(context, listen: false);
  late final LoadingModel loadingModel =
      Provider.of<LoadingModel>(context, listen: true);

  Future<void> _logout() async {
    AuthModel authModel = Provider.of<AuthModel>(context, listen: false);
    try {
      var refreshToken = await authModel.getRefreshToken();
      await OpenApiClient.client.getAuthApi().logout(
            logoutCommand: (api.LogoutCommandBuilder()
                  ..refreshToken = refreshToken ?? "")
                .build(),
          );
      await authModel.purgeTokens();
      context.go("/login");
    } catch (e) {
      showErrorSnackbar(context, e as dynamic);
    }
  }

  Widget? getIconButton() {
    if (widget.leadingArrowRedirect != null) {
      return IconButton(
        icon: const Icon(Icons.arrow_back),
        onPressed: () {
          if (widget.onLeadingArrowPressed != null) {
            widget.onLeadingArrowPressed!();
          }

          if (widget.leadingArrowPop == true) {
            context.pop();
            return;
          }

          context.go(widget.leadingArrowRedirect ?? "/",
              extra: widget.leadingArrowExtra);
        },
      );
    } else {
      return null;
    }
  }

  Widget getUserAvatar() {
    if (widget.hideAvatar == true) {
      return const SizedBox.shrink();
    }

    return PopupMenuButton(
        child: const UserAvatar(),
        itemBuilder: (BuildContext context) {
          return [
            PopupMenuItem(
              child: TextButton(
                onPressed: () => _logout(),
                child: const Text('Logout'),
              ),
            ),
          ];
        });
  }

  @override
  Widget build(BuildContext context) {
    PreferredSizeWidget? getLoadingIndicator(context) {
      if (loadingModel.isLoading) {
        return const PreferredSize(
            preferredSize: Size.square(4), child: LinearProgressIndicator());
      } else {
        return null;
      }
    }

    return AppBar(
      automaticallyImplyLeading: false,
      leading: getIconButton(),
      title: Text(widget.titleText),
      surfaceTintColor: widget.surfaceTintColor,
      bottom: getLoadingIndicator(context),
      actions: [
        Padding(
          padding: const EdgeInsets.only(right: 16.0),
          child: getUserAvatar(),
        ),
        ...widget.actions ?? [],
      ],
      centerTitle: false,
    );
  }
}
