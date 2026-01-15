// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_paged_request_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptPagedRequestCommand extends ReceiptPagedRequestCommand {
  @override
  final int page;
  @override
  final int pageSize;
  @override
  final String? orderBy;
  @override
  final SortDirection? sortDirection;
  @override
  final ReceiptPagedRequestFilter? filter;
  @override
  final bool? fullReceipts;

  factory _$ReceiptPagedRequestCommand(
          [void Function(ReceiptPagedRequestCommandBuilder)? updates]) =>
      (new ReceiptPagedRequestCommandBuilder()..update(updates))._build();

  _$ReceiptPagedRequestCommand._(
      {required this.page,
      required this.pageSize,
      this.orderBy,
      this.sortDirection,
      this.filter,
      this.fullReceipts})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        page, r'ReceiptPagedRequestCommand', 'page');
    BuiltValueNullFieldError.checkNotNull(
        pageSize, r'ReceiptPagedRequestCommand', 'pageSize');
  }

  @override
  ReceiptPagedRequestCommand rebuild(
          void Function(ReceiptPagedRequestCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptPagedRequestCommandBuilder toBuilder() =>
      new ReceiptPagedRequestCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptPagedRequestCommand &&
        page == other.page &&
        pageSize == other.pageSize &&
        orderBy == other.orderBy &&
        sortDirection == other.sortDirection &&
        filter == other.filter &&
        fullReceipts == other.fullReceipts;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, page.hashCode);
    _$hash = $jc(_$hash, pageSize.hashCode);
    _$hash = $jc(_$hash, orderBy.hashCode);
    _$hash = $jc(_$hash, sortDirection.hashCode);
    _$hash = $jc(_$hash, filter.hashCode);
    _$hash = $jc(_$hash, fullReceipts.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptPagedRequestCommand')
          ..add('page', page)
          ..add('pageSize', pageSize)
          ..add('orderBy', orderBy)
          ..add('sortDirection', sortDirection)
          ..add('filter', filter)
          ..add('fullReceipts', fullReceipts))
        .toString();
  }
}

class ReceiptPagedRequestCommandBuilder
    implements
        Builder<ReceiptPagedRequestCommand, ReceiptPagedRequestCommandBuilder> {
  _$ReceiptPagedRequestCommand? _$v;

  int? _page;
  int? get page => _$this._page;
  set page(int? page) => _$this._page = page;

  int? _pageSize;
  int? get pageSize => _$this._pageSize;
  set pageSize(int? pageSize) => _$this._pageSize = pageSize;

  String? _orderBy;
  String? get orderBy => _$this._orderBy;
  set orderBy(String? orderBy) => _$this._orderBy = orderBy;

  SortDirection? _sortDirection;
  SortDirection? get sortDirection => _$this._sortDirection;
  set sortDirection(SortDirection? sortDirection) =>
      _$this._sortDirection = sortDirection;

  ReceiptPagedRequestFilterBuilder? _filter;
  ReceiptPagedRequestFilterBuilder get filter =>
      _$this._filter ??= new ReceiptPagedRequestFilterBuilder();
  set filter(ReceiptPagedRequestFilterBuilder? filter) =>
      _$this._filter = filter;

  bool? _fullReceipts;
  bool? get fullReceipts => _$this._fullReceipts;
  set fullReceipts(bool? fullReceipts) => _$this._fullReceipts = fullReceipts;

  ReceiptPagedRequestCommandBuilder() {
    ReceiptPagedRequestCommand._defaults(this);
  }

  ReceiptPagedRequestCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _page = $v.page;
      _pageSize = $v.pageSize;
      _orderBy = $v.orderBy;
      _sortDirection = $v.sortDirection;
      _filter = $v.filter?.toBuilder();
      _fullReceipts = $v.fullReceipts;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptPagedRequestCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ReceiptPagedRequestCommand;
  }

  @override
  void update(void Function(ReceiptPagedRequestCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptPagedRequestCommand build() => _build();

  _$ReceiptPagedRequestCommand _build() {
    _$ReceiptPagedRequestCommand _$result;
    try {
      _$result = _$v ??
          new _$ReceiptPagedRequestCommand._(
              page: BuiltValueNullFieldError.checkNotNull(
                  page, r'ReceiptPagedRequestCommand', 'page'),
              pageSize: BuiltValueNullFieldError.checkNotNull(
                  pageSize, r'ReceiptPagedRequestCommand', 'pageSize'),
              orderBy: orderBy,
              sortDirection: sortDirection,
              filter: _filter?.build(),
              fullReceipts: fullReceipts);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'filter';
        _filter?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'ReceiptPagedRequestCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
