// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_request_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

abstract class PagedRequestCommandBuilder {
  void replace(PagedRequestCommand other);
  void update(void Function(PagedRequestCommandBuilder) updates);
  int? get page;
  set page(int? page);

  int? get pageSize;
  set pageSize(int? pageSize);

  String? get orderBy;
  set orderBy(String? orderBy);

  SortDirection? get sortDirection;
  set sortDirection(SortDirection? sortDirection);
}

class _$$PagedRequestCommand extends $PagedRequestCommand {
  @override
  final int page;
  @override
  final int pageSize;
  @override
  final String? orderBy;
  @override
  final SortDirection? sortDirection;

  factory _$$PagedRequestCommand(
          [void Function($PagedRequestCommandBuilder)? updates]) =>
      (new $PagedRequestCommandBuilder()..update(updates))._build();

  _$$PagedRequestCommand._(
      {required this.page,
      required this.pageSize,
      this.orderBy,
      this.sortDirection})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        page, r'$PagedRequestCommand', 'page');
    BuiltValueNullFieldError.checkNotNull(
        pageSize, r'$PagedRequestCommand', 'pageSize');
  }

  @override
  $PagedRequestCommand rebuild(
          void Function($PagedRequestCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  $PagedRequestCommandBuilder toBuilder() =>
      new $PagedRequestCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is $PagedRequestCommand &&
        page == other.page &&
        pageSize == other.pageSize &&
        orderBy == other.orderBy &&
        sortDirection == other.sortDirection;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, page.hashCode);
    _$hash = $jc(_$hash, pageSize.hashCode);
    _$hash = $jc(_$hash, orderBy.hashCode);
    _$hash = $jc(_$hash, sortDirection.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'$PagedRequestCommand')
          ..add('page', page)
          ..add('pageSize', pageSize)
          ..add('orderBy', orderBy)
          ..add('sortDirection', sortDirection))
        .toString();
  }
}

class $PagedRequestCommandBuilder
    implements
        Builder<$PagedRequestCommand, $PagedRequestCommandBuilder>,
        PagedRequestCommandBuilder {
  _$$PagedRequestCommand? _$v;

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

  $PagedRequestCommandBuilder() {
    $PagedRequestCommand._defaults(this);
  }

  $PagedRequestCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _page = $v.page;
      _pageSize = $v.pageSize;
      _orderBy = $v.orderBy;
      _sortDirection = $v.sortDirection;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant $PagedRequestCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$$PagedRequestCommand;
  }

  @override
  void update(void Function($PagedRequestCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  $PagedRequestCommand build() => _build();

  _$$PagedRequestCommand _build() {
    final _$result = _$v ??
        new _$$PagedRequestCommand._(
            page: BuiltValueNullFieldError.checkNotNull(
                page, r'$PagedRequestCommand', 'page'),
            pageSize: BuiltValueNullFieldError.checkNotNull(
                pageSize, r'$PagedRequestCommand', 'pageSize'),
            orderBy: orderBy,
            sortDirection: sortDirection);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
