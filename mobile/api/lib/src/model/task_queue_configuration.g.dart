// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'task_queue_configuration.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$TaskQueueConfiguration extends TaskQueueConfiguration {
  @override
  final QueueName? name;
  @override
  final int? priority;
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

  factory _$TaskQueueConfiguration(
          [void Function(TaskQueueConfigurationBuilder)? updates]) =>
      (new TaskQueueConfigurationBuilder()..update(updates))._build();

  _$TaskQueueConfiguration._(
      {this.name,
      this.priority,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'TaskQueueConfiguration', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'TaskQueueConfiguration', 'createdAt');
  }

  @override
  TaskQueueConfiguration rebuild(
          void Function(TaskQueueConfigurationBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  TaskQueueConfigurationBuilder toBuilder() =>
      new TaskQueueConfigurationBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is TaskQueueConfiguration &&
        name == other.name &&
        priority == other.priority &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, priority.hashCode);
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
    return (newBuiltValueToStringHelper(r'TaskQueueConfiguration')
          ..add('name', name)
          ..add('priority', priority)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class TaskQueueConfigurationBuilder
    implements
        Builder<TaskQueueConfiguration, TaskQueueConfigurationBuilder>,
        BaseModelBuilder {
  _$TaskQueueConfiguration? _$v;

  QueueName? _name;
  QueueName? get name => _$this._name;
  set name(covariant QueueName? name) => _$this._name = name;

  int? _priority;
  int? get priority => _$this._priority;
  set priority(covariant int? priority) => _$this._priority = priority;

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

  TaskQueueConfigurationBuilder() {
    TaskQueueConfiguration._defaults(this);
  }

  TaskQueueConfigurationBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _priority = $v.priority;
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
  void replace(covariant TaskQueueConfiguration other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$TaskQueueConfiguration;
  }

  @override
  void update(void Function(TaskQueueConfigurationBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  TaskQueueConfiguration build() => _build();

  _$TaskQueueConfiguration _build() {
    final _$result = _$v ??
        new _$TaskQueueConfiguration._(
            name: name,
            priority: priority,
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'TaskQueueConfiguration', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'TaskQueueConfiguration', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
