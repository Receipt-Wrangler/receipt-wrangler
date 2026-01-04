import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:openapi/openapi.dart' as api;

class PagedDataList extends StatefulWidget {
  const PagedDataList({
    super.key,
    required this.getPagedDataFuture,
    required this.listItemBuilder,
    required this.noItemsFoundText,
    this.onRefreshCallbackSet,
  });

  final Future<Response<api.PagedData>> Function(int) getPagedDataFuture;

  final Widget Function(BuildContext, api.PagedDataDataInner, int)
      listItemBuilder;

  final String noItemsFoundText;

  final void Function(VoidCallback)? onRefreshCallbackSet;

  @override
  _PagedDataListState createState() {
    return _PagedDataListState();
  }
}

class _PagedDataListState extends State<PagedDataList> {
  late final PagingController<int, api.PagedDataDataInner> _pagingController;
  int? _totalCount;

  @override
  void initState() {
    super.initState();

    _pagingController = PagingController<int, api.PagedDataDataInner>(
      getNextPageKey: (state) {
        if (_totalCount == null) return state.nextIntPageKey;
        final itemsLength = state.items?.length ?? 0;
        if (itemsLength >= _totalCount!) return null;
        return state.nextIntPageKey;
      },
      fetchPage: (pageKey) async {
        final response = await widget.getPagedDataFuture(pageKey);
        if (response.data != null) {
          _totalCount = response.data!.totalCount;
          return response.data!.data.toList();
        }
        return [];
      },
    );

    // Provide refresh callback to parent widget
    widget.onRefreshCallbackSet?.call(() {
      _pagingController.refresh();
    });
  }

  @override
  void dispose() {
    _pagingController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: PagingListener<int, api.PagedDataDataInner>(
        controller: _pagingController,
        builder: (context, state, fetchNextPage) {
          return PagedListView<int, api.PagedDataDataInner>(
            state: state,
            fetchNextPage: fetchNextPage,
            builderDelegate: PagedChildBuilderDelegate<api.PagedDataDataInner>(
              itemBuilder: (context, item, index) {
                return widget.listItemBuilder(context, item, index);
              },
              noItemsFoundIndicatorBuilder: (context) => Center(
                child: Text(widget.noItemsFoundText),
              ),
              invisibleItemsThreshold: 5,
            ),
          );
        },
      ),
    );
  }
}
