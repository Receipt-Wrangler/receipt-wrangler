// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_group_request_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PagedGroupRequestCommand extends PagedGroupRequestCommand {
  @override
  final GroupFilter? filter;
  @override
  final int page;
  @override
  final int pageSize;
  @override
  final String? orderBy;
  @override
  final SortDirection? sortDirection;

  factory _$PagedGroupRequestCommand(
          [void Function(PagedGroupRequestCommandBuilder)? updates]) =>
      (new PagedGroupRequestCommandBuilder()..update(updates))._build();

  _$PagedGroupRequestCommand._(
      {this.filter,
      required this.page,
      required this.pageSize,
      this.orderBy,
      this.sortDirection})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        page, r'PagedGroupRequestCommand', 'page');
    BuiltValueNullFieldError.checkNotNull(
        pageSize, r'PagedGroupRequestCommand', 'pageSize');
  }

  @override
  PagedGroupRequestCommand rebuild(
          void Function(PagedGroupRequestCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PagedGroupRequestCommandBuilder toBuilder() =>
      new PagedGroupRequestCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PagedGroupRequestCommand &&
        filter == other.filter &&
        page == other.page &&
        pageSize == other.pageSize &&
        orderBy == other.orderBy &&
        sortDirection == other.sortDirection;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, filter.hashCode);
    _$hash = $jc(_$hash, page.hashCode);
    _$hash = $jc(_$hash, pageSize.hashCode);
    _$hash = $jc(_$hash, orderBy.hashCode);
    _$hash = $jc(_$hash, sortDirection.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PagedGroupRequestCommand')
          ..add('filter', filter)
          ..add('page', page)
          ..add('pageSize', pageSize)
          ..add('orderBy', orderBy)
          ..add('sortDirection', sortDirection))
        .toString();
  }
}

class PagedGroupRequestCommandBuilder
    implements
        Builder<PagedGroupRequestCommand, PagedGroupRequestCommandBuilder>,
        PagedRequestCommandBuilder {
  _$PagedGroupRequestCommand? _$v;

  GroupFilterBuilder? _filter;
  GroupFilterBuilder get filter => _$this._filter ??= new GroupFilterBuilder();
  set filter(covariant GroupFilterBuilder? filter) => _$this._filter = filter;

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

  PagedGroupRequestCommandBuilder() {
    PagedGroupRequestCommand._defaults(this);
  }

  PagedGroupRequestCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _filter = $v.filter?.toBuilder();
      _page = $v.page;
      _pageSize = $v.pageSize;
      _orderBy = $v.orderBy;
      _sortDirection = $v.sortDirection;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant PagedGroupRequestCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PagedGroupRequestCommand;
  }

  @override
  void update(void Function(PagedGroupRequestCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PagedGroupRequestCommand build() => _build();

  _$PagedGroupRequestCommand _build() {
    _$PagedGroupRequestCommand _$result;
    try {
      _$result = _$v ??
          new _$PagedGroupRequestCommand._(
              filter: _filter?.build(),
              page: BuiltValueNullFieldError.checkNotNull(
                  page, r'PagedGroupRequestCommand', 'page'),
              pageSize: BuiltValueNullFieldError.checkNotNull(
                  pageSize, r'PagedGroupRequestCommand', 'pageSize'),
              orderBy: orderBy,
              sortDirection: sortDirection);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'filter';
        _filter?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'PagedGroupRequestCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
