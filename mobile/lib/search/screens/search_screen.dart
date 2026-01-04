import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_list_item.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';

import '../../models/search_model.dart';

class SearchScreen extends StatefulWidget {
  const SearchScreen({super.key});

  @override
  State<SearchScreen> createState() => _SearchScreen();
}

class _SearchScreen extends State<SearchScreen> {
  late final searchModel = Provider.of<SearchModel>(context);
  late final groupModel = Provider.of<GroupModel>(context);

  Widget buildSearchTextText(String searchText) {
    return Text(searchText);
  }

  Widget buildReceiptList(List<api.SearchResult> searchResults) {
    return ListView.builder(
      itemCount: searchResults.length,
      itemBuilder: (context, index) {
        var searchResult = searchResults[index];
        var receipt = (api.ReceiptBuilder()
              ..amount = searchResult.amount
              ..date = searchResult.date
              ..groupId = searchResult.groupId
              ..id = searchResult.id
              ..name = searchResult.name
              ..paidByUserId = searchResult.paidByUserId
              ..status = searchResult.receiptStatus
              ..createdAt = searchResult.createdAt)
            .build();

        return ReceiptListItem(displayGroup: true, receipt: receipt);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Consumer(builder: (context, SearchModel searchModel, child) {
      return buildReceiptList(searchModel.searchResults);
    });
  }
}
