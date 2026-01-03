// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'activity.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Activity extends Activity {
  @override
  final int id;
  @override
  final SystemTaskType type;
  @override
  final SystemTaskStatus status;
  @override
  final String startedAt;
  @override
  final String endedAt;
  @override
  final int? ranByUserId;
  @override
  final int? receiptId;
  @override
  final int? groupId;
  @override
  final bool? canBeRestarted;

  factory _$Activity([void Function(ActivityBuilder)? updates]) =>
      (new ActivityBuilder()..update(updates))._build();

  _$Activity._(
      {required this.id,
      required this.type,
      required this.status,
      required this.startedAt,
      required this.endedAt,
      this.ranByUserId,
      this.receiptId,
      this.groupId,
      this.canBeRestarted})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'Activity', 'id');
    BuiltValueNullFieldError.checkNotNull(type, r'Activity', 'type');
    BuiltValueNullFieldError.checkNotNull(status, r'Activity', 'status');
    BuiltValueNullFieldError.checkNotNull(startedAt, r'Activity', 'startedAt');
    BuiltValueNullFieldError.checkNotNull(endedAt, r'Activity', 'endedAt');
  }

  @override
  Activity rebuild(void Function(ActivityBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ActivityBuilder toBuilder() => new ActivityBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Activity &&
        id == other.id &&
        type == other.type &&
        status == other.status &&
        startedAt == other.startedAt &&
        endedAt == other.endedAt &&
        ranByUserId == other.ranByUserId &&
        receiptId == other.receiptId &&
        groupId == other.groupId &&
        canBeRestarted == other.canBeRestarted;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, startedAt.hashCode);
    _$hash = $jc(_$hash, endedAt.hashCode);
    _$hash = $jc(_$hash, ranByUserId.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, canBeRestarted.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Activity')
          ..add('id', id)
          ..add('type', type)
          ..add('status', status)
          ..add('startedAt', startedAt)
          ..add('endedAt', endedAt)
          ..add('ranByUserId', ranByUserId)
          ..add('receiptId', receiptId)
          ..add('groupId', groupId)
          ..add('canBeRestarted', canBeRestarted))
        .toString();
  }
}

class ActivityBuilder implements Builder<Activity, ActivityBuilder> {
  _$Activity? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  SystemTaskType? _type;
  SystemTaskType? get type => _$this._type;
  set type(SystemTaskType? type) => _$this._type = type;

  SystemTaskStatus? _status;
  SystemTaskStatus? get status => _$this._status;
  set status(SystemTaskStatus? status) => _$this._status = status;

  String? _startedAt;
  String? get startedAt => _$this._startedAt;
  set startedAt(String? startedAt) => _$this._startedAt = startedAt;

  String? _endedAt;
  String? get endedAt => _$this._endedAt;
  set endedAt(String? endedAt) => _$this._endedAt = endedAt;

  int? _ranByUserId;
  int? get ranByUserId => _$this._ranByUserId;
  set ranByUserId(int? ranByUserId) => _$this._ranByUserId = ranByUserId;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  bool? _canBeRestarted;
  bool? get canBeRestarted => _$this._canBeRestarted;
  set canBeRestarted(bool? canBeRestarted) =>
      _$this._canBeRestarted = canBeRestarted;

  ActivityBuilder() {
    Activity._defaults(this);
  }

  ActivityBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _type = $v.type;
      _status = $v.status;
      _startedAt = $v.startedAt;
      _endedAt = $v.endedAt;
      _ranByUserId = $v.ranByUserId;
      _receiptId = $v.receiptId;
      _groupId = $v.groupId;
      _canBeRestarted = $v.canBeRestarted;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Activity other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Activity;
  }

  @override
  void update(void Function(ActivityBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Activity build() => _build();

  _$Activity _build() {
    final _$result = _$v ??
        new _$Activity._(
            id: BuiltValueNullFieldError.checkNotNull(id, r'Activity', 'id'),
            type: BuiltValueNullFieldError.checkNotNull(
                type, r'Activity', 'type'),
            status: BuiltValueNullFieldError.checkNotNull(
                status, r'Activity', 'status'),
            startedAt: BuiltValueNullFieldError.checkNotNull(
                startedAt, r'Activity', 'startedAt'),
            endedAt: BuiltValueNullFieldError.checkNotNull(
                endedAt, r'Activity', 'endedAt'),
            ranByUserId: ranByUserId,
            receiptId: receiptId,
            groupId: groupId,
            canBeRestarted: canBeRestarted);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
