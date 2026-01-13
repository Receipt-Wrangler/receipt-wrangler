// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'system_task.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$SystemTask extends SystemTask {
  @override
  final int? associatedEntityId;
  @override
  final String? endedAt;
  @override
  final AssociatedEntityType? associatedEntityType;
  @override
  final int? groupId;
  @override
  final BuiltList<SystemTask>? childSystemTasks;
  @override
  final String? startedAt;
  @override
  final String? apiKeyId;
  @override
  final String? resultDescription;
  @override
  final SystemTaskType? type;
  @override
  final int? receiptId;
  @override
  final int? ranByUserId;
  @override
  final SystemTaskStatus? status;
  @override
  final int id;
  @override
  final String createdAt;
  @override
  final int? createdBy;
  @override
  final String? createdByString;
  @override
  final String? updatedAt;

  factory _$SystemTask([void Function(SystemTaskBuilder)? updates]) =>
      (new SystemTaskBuilder()..update(updates))._build();

  _$SystemTask._(
      {this.associatedEntityId,
      this.endedAt,
      this.associatedEntityType,
      this.groupId,
      this.childSystemTasks,
      this.startedAt,
      this.apiKeyId,
      this.resultDescription,
      this.type,
      this.receiptId,
      this.ranByUserId,
      this.status,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'SystemTask', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'SystemTask', 'createdAt');
  }

  @override
  SystemTask rebuild(void Function(SystemTaskBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  SystemTaskBuilder toBuilder() => new SystemTaskBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is SystemTask &&
        associatedEntityId == other.associatedEntityId &&
        endedAt == other.endedAt &&
        associatedEntityType == other.associatedEntityType &&
        groupId == other.groupId &&
        childSystemTasks == other.childSystemTasks &&
        startedAt == other.startedAt &&
        apiKeyId == other.apiKeyId &&
        resultDescription == other.resultDescription &&
        type == other.type &&
        receiptId == other.receiptId &&
        ranByUserId == other.ranByUserId &&
        status == other.status &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, associatedEntityId.hashCode);
    _$hash = $jc(_$hash, endedAt.hashCode);
    _$hash = $jc(_$hash, associatedEntityType.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, childSystemTasks.hashCode);
    _$hash = $jc(_$hash, startedAt.hashCode);
    _$hash = $jc(_$hash, apiKeyId.hashCode);
    _$hash = $jc(_$hash, resultDescription.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, ranByUserId.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, createdByString.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'SystemTask')
          ..add('associatedEntityId', associatedEntityId)
          ..add('endedAt', endedAt)
          ..add('associatedEntityType', associatedEntityType)
          ..add('groupId', groupId)
          ..add('childSystemTasks', childSystemTasks)
          ..add('startedAt', startedAt)
          ..add('apiKeyId', apiKeyId)
          ..add('resultDescription', resultDescription)
          ..add('type', type)
          ..add('receiptId', receiptId)
          ..add('ranByUserId', ranByUserId)
          ..add('status', status)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class SystemTaskBuilder
    implements Builder<SystemTask, SystemTaskBuilder>, BaseModelBuilder {
  _$SystemTask? _$v;

  int? _associatedEntityId;
  int? get associatedEntityId => _$this._associatedEntityId;
  set associatedEntityId(covariant int? associatedEntityId) =>
      _$this._associatedEntityId = associatedEntityId;

  String? _endedAt;
  String? get endedAt => _$this._endedAt;
  set endedAt(covariant String? endedAt) => _$this._endedAt = endedAt;

  AssociatedEntityType? _associatedEntityType;
  AssociatedEntityType? get associatedEntityType =>
      _$this._associatedEntityType;
  set associatedEntityType(
          covariant AssociatedEntityType? associatedEntityType) =>
      _$this._associatedEntityType = associatedEntityType;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(covariant int? groupId) => _$this._groupId = groupId;

  ListBuilder<SystemTask>? _childSystemTasks;
  ListBuilder<SystemTask> get childSystemTasks =>
      _$this._childSystemTasks ??= new ListBuilder<SystemTask>();
  set childSystemTasks(covariant ListBuilder<SystemTask>? childSystemTasks) =>
      _$this._childSystemTasks = childSystemTasks;

  String? _startedAt;
  String? get startedAt => _$this._startedAt;
  set startedAt(covariant String? startedAt) => _$this._startedAt = startedAt;

  String? _apiKeyId;
  String? get apiKeyId => _$this._apiKeyId;
  set apiKeyId(covariant String? apiKeyId) => _$this._apiKeyId = apiKeyId;

  String? _resultDescription;
  String? get resultDescription => _$this._resultDescription;
  set resultDescription(covariant String? resultDescription) =>
      _$this._resultDescription = resultDescription;

  SystemTaskType? _type;
  SystemTaskType? get type => _$this._type;
  set type(covariant SystemTaskType? type) => _$this._type = type;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(covariant int? receiptId) => _$this._receiptId = receiptId;

  int? _ranByUserId;
  int? get ranByUserId => _$this._ranByUserId;
  set ranByUserId(covariant int? ranByUserId) =>
      _$this._ranByUserId = ranByUserId;

  SystemTaskStatus? _status;
  SystemTaskStatus? get status => _$this._status;
  set status(covariant SystemTaskStatus? status) => _$this._status = status;

  int? _id;
  int? get id => _$this._id;
  set id(covariant int? id) => _$this._id = id;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(covariant String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(covariant int? createdBy) => _$this._createdBy = createdBy;

  String? _createdByString;
  String? get createdByString => _$this._createdByString;
  set createdByString(covariant String? createdByString) =>
      _$this._createdByString = createdByString;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(covariant String? updatedAt) => _$this._updatedAt = updatedAt;

  SystemTaskBuilder() {
    SystemTask._defaults(this);
  }

  SystemTaskBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _associatedEntityId = $v.associatedEntityId;
      _endedAt = $v.endedAt;
      _associatedEntityType = $v.associatedEntityType;
      _groupId = $v.groupId;
      _childSystemTasks = $v.childSystemTasks?.toBuilder();
      _startedAt = $v.startedAt;
      _apiKeyId = $v.apiKeyId;
      _resultDescription = $v.resultDescription;
      _type = $v.type;
      _receiptId = $v.receiptId;
      _ranByUserId = $v.ranByUserId;
      _status = $v.status;
      _id = $v.id;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _createdByString = $v.createdByString;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant SystemTask other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$SystemTask;
  }

  @override
  void update(void Function(SystemTaskBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  SystemTask build() => _build();

  _$SystemTask _build() {
    _$SystemTask _$result;
    try {
      _$result = _$v ??
          new _$SystemTask._(
              associatedEntityId: associatedEntityId,
              endedAt: endedAt,
              associatedEntityType: associatedEntityType,
              groupId: groupId,
              childSystemTasks: _childSystemTasks?.build(),
              startedAt: startedAt,
              apiKeyId: apiKeyId,
              resultDescription: resultDescription,
              type: type,
              receiptId: receiptId,
              ranByUserId: ranByUserId,
              status: status,
              id: BuiltValueNullFieldError.checkNotNull(
                  id, r'SystemTask', 'id'),
              createdAt: BuiltValueNullFieldError.checkNotNull(
                  createdAt, r'SystemTask', 'createdAt'),
              createdBy: createdBy,
              createdByString: createdByString,
              updatedAt: updatedAt);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'childSystemTasks';
        _childSystemTasks?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'SystemTask', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
