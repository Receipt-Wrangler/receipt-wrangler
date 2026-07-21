import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;

import '../../client/client.dart';
import '../../shared/widgets/paged_data_list.dart';
import 'report_list_item.dart';

/// Paged list of saved report templates (`POST /report/template/list`), sorted by
/// most-recently-updated to match the desktop list. Each row is a
/// [ReportListItem]; a refresh callback lets row actions (delete) reload the list.
class ReportList extends StatefulWidget {
  const ReportList({super.key});

  @override
  State<ReportList> createState() => _ReportListState();
}

class _ReportListState extends State<ReportList> {
  VoidCallback? _refresh;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        PagedDataList(
          noItemsFoundText: 'No reports found',
          onRefreshCallbackSet: (callback) => _refresh = callback,
          listItemBuilder: (context, item, index) {
            // AnyOf.values is a Map<int, Object?> (type-index -> value); pull the
            // ReportTemplate out by type rather than a brittle fixed index.
            final matches =
                item.anyOf.values.values.whereType<api.ReportTemplate>();
            if (matches.isEmpty) return const SizedBox.shrink();
            return ReportListItem(
              template: matches.first,
              onChanged: () => _refresh?.call(),
            );
          },
          getPagedDataFuture: (pageKey) {
            final command = (api.$PagedRequestCommandBuilder()
                  ..page = pageKey
                  ..pageSize = 20
                  ..orderBy = 'updated_at'
                  ..sortDirection = api.SortDirection.desc)
                .build();
            return OpenApiClient.client
                .getReportApi()
                .getReportTemplates(pagedRequestCommand: command);
          },
        ),
      ],
    );
  }
}
