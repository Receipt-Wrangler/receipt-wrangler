import 'package:flutter/material.dart';

import '../../shared/widgets/screen_wrapper.dart';
import '../../shared/widgets/top_app_bar.dart';
import '../widgets/report_list.dart';

/// Standalone "Reports" screen reached from the avatar menu. App-scoped (mirrors
/// the desktop top-level `/reports` route); the route is gated on
/// `app.reports.read` / `app.reports.readAll` by `reportsReadRedirect`, and each
/// row's actions are gated on the server-computed `allowedActions`.
class ReportListScreen extends StatelessWidget {
  const ReportListScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const ScreenWrapper(
      appBarWidget: TopAppBar(
        titleText: 'Reports',
        leadingArrowPop: true,
        hideAvatar: true,
      ),
      child: ReportList(),
    );
  }
}
