// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_api_key_request_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PagedApiKeyRequestCommand extends PagedApiKeyRequestCommand {
  @override
  final ApiKeyFilter? filter;
  @override
  final int page;
  @override
  final int pageSize;
  @override
  final String? orderBy;
  @override
  final SortDirection? sortDirection;

  factory _$PagedApiKeyRequestCommand(
          [void Function(PagedApiKeyRequestCommandBuilder)? updates]) =>
      (new PagedApiKeyRequestCommandBuilder()..update(updates))._build();

  _$PagedApiKeyRequestCommand._(
      {this.filter,
      required this.page,
      required this.pageSize,
      this.orderBy,
      this.sortDirection})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        page, r'PagedApiKeyRequestCommand', 'page');
    BuiltValueNullFieldError.checkNotNull(
        pageSize, r'PagedApiKeyRequestCommand', 'pageSize');
  }

  @override
  PagedApiKeyRequestCommand rebuild(
          void Function(PagedApiKeyRequestCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PagedApiKeyRequestCommandBuilder toBuilder() =>
      new PagedApiKeyRequestCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PagedApiKeyRequestCommand &&
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
    return (newBuiltValueToStringHelper(r'PagedApiKeyRequestCommand')
          ..add('filter', filter)
          ..add('page', page)
          ..add('pageSize', pageSize)
          ..add('orderBy', orderBy)
          ..add('sortDirection', sortDirection))
        .toString();
  }
}

class PagedApiKeyRequestCommandBuilder
    implements
        Builder<PagedApiKeyRequestCommand, PagedApiKeyRequestCommandBuilder>,
        PagedRequestCommandBuilder {
  _$PagedApiKeyRequestCommand? _$v;

  ApiKeyFilterBuilder? _filter;
  ApiKeyFilterBuilder get filter =>
      _$this._filter ??= new ApiKeyFilterBuilder();
  set filter(covariant ApiKeyFilterBuilder? filter) => _$this._filter = filter;

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

  PagedApiKeyRequestCommandBuilder() {
    PagedApiKeyRequestCommand._defaults(this);
  }

  PagedApiKeyRequestCommandBuilder get _$this {
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
  void replace(covariant PagedApiKeyRequestCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PagedApiKeyRequestCommand;
  }

  @override
  void update(void Function(PagedApiKeyRequestCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PagedApiKeyRequestCommand build() => _build();

  _$PagedApiKeyRequestCommand _build() {
    _$PagedApiKeyRequestCommand _$result;
    try {
      _$result = _$v ??
          new _$PagedApiKeyRequestCommand._(
              filter: _filter?.build(),
              page: BuiltValueNullFieldError.checkNotNull(
                  page, r'PagedApiKeyRequestCommand', 'page'),
              pageSize: BuiltValueNullFieldError.checkNotNull(
                  pageSize, r'PagedApiKeyRequestCommand', 'pageSize'),
              orderBy: orderBy,
              sortDirection: sortDirection);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'filter';
        _filter?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'PagedApiKeyRequestCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
