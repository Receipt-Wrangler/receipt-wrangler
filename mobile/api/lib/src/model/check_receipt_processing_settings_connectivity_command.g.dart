// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'check_receipt_processing_settings_connectivity_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CheckReceiptProcessingSettingsConnectivityCommand
    extends CheckReceiptProcessingSettingsConnectivityCommand {
  @override
  final int? id;
  @override
  final String? name;
  @override
  final AiType? aiType;
  @override
  final String? url;
  @override
  final String? key;
  @override
  final String? model;
  @override
  final int? numWorkers;
  @override
  final OcrEngine? ocrEngine;
  @override
  final int? promptId;

  factory _$CheckReceiptProcessingSettingsConnectivityCommand(
          [void Function(
                  CheckReceiptProcessingSettingsConnectivityCommandBuilder)?
              updates]) =>
      (new CheckReceiptProcessingSettingsConnectivityCommandBuilder()
            ..update(updates))
          ._build();

  _$CheckReceiptProcessingSettingsConnectivityCommand._(
      {this.id,
      this.name,
      this.aiType,
      this.url,
      this.key,
      this.model,
      this.numWorkers,
      this.ocrEngine,
      this.promptId})
      : super._();

  @override
  CheckReceiptProcessingSettingsConnectivityCommand rebuild(
          void Function(
                  CheckReceiptProcessingSettingsConnectivityCommandBuilder)
              updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CheckReceiptProcessingSettingsConnectivityCommandBuilder toBuilder() =>
      new CheckReceiptProcessingSettingsConnectivityCommandBuilder()
        ..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CheckReceiptProcessingSettingsConnectivityCommand &&
        id == other.id &&
        name == other.name &&
        aiType == other.aiType &&
        url == other.url &&
        key == other.key &&
        model == other.model &&
        numWorkers == other.numWorkers &&
        ocrEngine == other.ocrEngine &&
        promptId == other.promptId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, aiType.hashCode);
    _$hash = $jc(_$hash, url.hashCode);
    _$hash = $jc(_$hash, key.hashCode);
    _$hash = $jc(_$hash, model.hashCode);
    _$hash = $jc(_$hash, numWorkers.hashCode);
    _$hash = $jc(_$hash, ocrEngine.hashCode);
    _$hash = $jc(_$hash, promptId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(
            r'CheckReceiptProcessingSettingsConnectivityCommand')
          ..add('id', id)
          ..add('name', name)
          ..add('aiType', aiType)
          ..add('url', url)
          ..add('key', key)
          ..add('model', model)
          ..add('numWorkers', numWorkers)
          ..add('ocrEngine', ocrEngine)
          ..add('promptId', promptId))
        .toString();
  }
}

class CheckReceiptProcessingSettingsConnectivityCommandBuilder
    implements
        Builder<CheckReceiptProcessingSettingsConnectivityCommand,
            CheckReceiptProcessingSettingsConnectivityCommandBuilder> {
  _$CheckReceiptProcessingSettingsConnectivityCommand? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  AiType? _aiType;
  AiType? get aiType => _$this._aiType;
  set aiType(AiType? aiType) => _$this._aiType = aiType;

  String? _url;
  String? get url => _$this._url;
  set url(String? url) => _$this._url = url;

  String? _key;
  String? get key => _$this._key;
  set key(String? key) => _$this._key = key;

  String? _model;
  String? get model => _$this._model;
  set model(String? model) => _$this._model = model;

  int? _numWorkers;
  int? get numWorkers => _$this._numWorkers;
  set numWorkers(int? numWorkers) => _$this._numWorkers = numWorkers;

  OcrEngine? _ocrEngine;
  OcrEngine? get ocrEngine => _$this._ocrEngine;
  set ocrEngine(OcrEngine? ocrEngine) => _$this._ocrEngine = ocrEngine;

  int? _promptId;
  int? get promptId => _$this._promptId;
  set promptId(int? promptId) => _$this._promptId = promptId;

  CheckReceiptProcessingSettingsConnectivityCommandBuilder() {
    CheckReceiptProcessingSettingsConnectivityCommand._defaults(this);
  }

  CheckReceiptProcessingSettingsConnectivityCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _name = $v.name;
      _aiType = $v.aiType;
      _url = $v.url;
      _key = $v.key;
      _model = $v.model;
      _numWorkers = $v.numWorkers;
      _ocrEngine = $v.ocrEngine;
      _promptId = $v.promptId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(CheckReceiptProcessingSettingsConnectivityCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$CheckReceiptProcessingSettingsConnectivityCommand;
  }

  @override
  void update(
      void Function(CheckReceiptProcessingSettingsConnectivityCommandBuilder)?
          updates) {
    if (updates != null) updates(this);
  }

  @override
  CheckReceiptProcessingSettingsConnectivityCommand build() => _build();

  _$CheckReceiptProcessingSettingsConnectivityCommand _build() {
    final _$result = _$v ??
        new _$CheckReceiptProcessingSettingsConnectivityCommand._(
            id: id,
            name: name,
            aiType: aiType,
            url: url,
            key: key,
            model: model,
            numWorkers: numWorkers,
            ocrEngine: ocrEngine,
            promptId: promptId);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
