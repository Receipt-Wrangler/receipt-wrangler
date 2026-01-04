// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_task_queue_configuration.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertTaskQueueConfiguration extends UpsertTaskQueueConfiguration {
  @override
  final QueueName? name;
  @override
  final int? priority;

  factory _$UpsertTaskQueueConfiguration(
          [void Function(UpsertTaskQueueConfigurationBuilder)? updates]) =>
      (new UpsertTaskQueueConfigurationBuilder()..update(updates))._build();

  _$UpsertTaskQueueConfiguration._({this.name, this.priority}) : super._();

  @override
  UpsertTaskQueueConfiguration rebuild(
          void Function(UpsertTaskQueueConfigurationBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertTaskQueueConfigurationBuilder toBuilder() =>
      new UpsertTaskQueueConfigurationBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertTaskQueueConfiguration &&
        name == other.name &&
        priority == other.priority;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, priority.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertTaskQueueConfiguration')
          ..add('name', name)
          ..add('priority', priority))
        .toString();
  }
}

class UpsertTaskQueueConfigurationBuilder
    implements
        Builder<UpsertTaskQueueConfiguration,
            UpsertTaskQueueConfigurationBuilder> {
  _$UpsertTaskQueueConfiguration? _$v;

  QueueName? _name;
  QueueName? get name => _$this._name;
  set name(QueueName? name) => _$this._name = name;

  int? _priority;
  int? get priority => _$this._priority;
  set priority(int? priority) => _$this._priority = priority;

  UpsertTaskQueueConfigurationBuilder() {
    UpsertTaskQueueConfiguration._defaults(this);
  }

  UpsertTaskQueueConfigurationBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _priority = $v.priority;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertTaskQueueConfiguration other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertTaskQueueConfiguration;
  }

  @override
  void update(void Function(UpsertTaskQueueConfigurationBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertTaskQueueConfiguration build() => _build();

  _$UpsertTaskQueueConfiguration _build() {
    final _$result = _$v ??
        new _$UpsertTaskQueueConfiguration._(name: name, priority: priority);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
