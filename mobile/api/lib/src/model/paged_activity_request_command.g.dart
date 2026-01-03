// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_activity_request_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PagedActivityRequestCommand extends PagedActivityRequestCommand {
  @override
  final BuiltList<int>? groupIds;
  @override
  final int page;
  @override
  final int pageSize;
  @override
  final String? orderBy;
  @override
  final SortDirection? sortDirection;

  factory _$PagedActivityRequestCommand(
          [void Function(PagedActivityRequestCommandBuilder)? updates]) =>
      (new PagedActivityRequestCommandBuilder()..update(updates))._build();

  _$PagedActivityRequestCommand._(
      {this.groupIds,
      required this.page,
      required this.pageSize,
      this.orderBy,
      this.sortDirection})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        page, r'PagedActivityRequestCommand', 'page');
    BuiltValueNullFieldError.checkNotNull(
        pageSize, r'PagedActivityRequestCommand', 'pageSize');
  }

  @override
  PagedActivityRequestCommand rebuild(
          void Function(PagedActivityRequestCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PagedActivityRequestCommandBuilder toBuilder() =>
      new PagedActivityRequestCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PagedActivityRequestCommand &&
        groupIds == other.groupIds &&
        page == other.page &&
        pageSize == other.pageSize &&
        orderBy == other.orderBy &&
        sortDirection == other.sortDirection;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, groupIds.hashCode);
    _$hash = $jc(_$hash, page.hashCode);
    _$hash = $jc(_$hash, pageSize.hashCode);
    _$hash = $jc(_$hash, orderBy.hashCode);
    _$hash = $jc(_$hash, sortDirection.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PagedActivityRequestCommand')
          ..add('groupIds', groupIds)
          ..add('page', page)
          ..add('pageSize', pageSize)
          ..add('orderBy', orderBy)
          ..add('sortDirection', sortDirection))
        .toString();
  }
}

class PagedActivityRequestCommandBuilder
    implements
        Builder<PagedActivityRequestCommand,
            PagedActivityRequestCommandBuilder>,
        PagedRequestCommandBuilder {
  _$PagedActivityRequestCommand? _$v;

  ListBuilder<int>? _groupIds;
  ListBuilder<int> get groupIds => _$this._groupIds ??= new ListBuilder<int>();
  set groupIds(covariant ListBuilder<int>? groupIds) =>
      _$this._groupIds = groupIds;

  int? _page;
  int? get page => _$this._page;
  set page(covariant int? page) => _$this._page = page;

  int? _pageSize;
  int? get pageSize => _$this._pageSize;
  set pageSize(covariant int? pageSize) => _$this._pageSize = pageSize;

  String? _orderBy;
  String? get orderBy => _$this._orderBy;
  set orderBy(covariant String? orderBy) => _$this._orderBy = orderBy;

  SortDirection? _sortDirection;
  SortDirection? get sortDirection => _$this._sortDirection;
  set sortDirection(covariant SortDirection? sortDirection) =>
      _$this._sortDirection = sortDirection;

  PagedActivityRequestCommandBuilder() {
    PagedActivityRequestCommand._defaults(this);
  }

  PagedActivityRequestCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _groupIds = $v.groupIds?.toBuilder();
      _page = $v.page;
      _pageSize = $v.pageSize;
      _orderBy = $v.orderBy;
      _sortDirection = $v.sortDirection;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant PagedActivityRequestCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PagedActivityRequestCommand;
  }

  @override
  void update(void Function(PagedActivityRequestCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PagedActivityRequestCommand build() => _build();

  _$PagedActivityRequestCommand _build() {
    _$PagedActivityRequestCommand _$result;
    try {
      _$result = _$v ??
          new _$PagedActivityRequestCommand._(
              groupIds: _groupIds?.build(),
              page: BuiltValueNullFieldError.checkNotNull(
                  page, r'PagedActivityRequestCommand', 'page'),
              pageSize: BuiltValueNullFieldError.checkNotNull(
                  pageSize, r'PagedActivityRequestCommand', 'pageSize'),
              orderBy: orderBy,
              sortDirection: sortDirection);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'groupIds';
        _groupIds?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'PagedActivityRequestCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
