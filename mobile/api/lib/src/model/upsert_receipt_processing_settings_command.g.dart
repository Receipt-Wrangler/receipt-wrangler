// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_receipt_processing_settings_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertReceiptProcessingSettingsCommand
    extends UpsertReceiptProcessingSettingsCommand {
  @override
  final String name;
  @override
  final String? description;
  @override
  final AiType aiType;
  @override
  final String? url;
  @override
  final String? key;
  @override
  final String? model;
  @override
  final bool? isVisionModel;
  @override
  final OcrEngine ocrEngine;
  @override
  final int promptId;

  factory _$UpsertReceiptProcessingSettingsCommand(
          [void Function(UpsertReceiptProcessingSettingsCommandBuilder)?
              updates]) =>
      (new UpsertReceiptProcessingSettingsCommandBuilder()..update(updates))
          ._build();

  _$UpsertReceiptProcessingSettingsCommand._(
      {required this.name,
      this.description,
      required this.aiType,
      this.url,
      this.key,
      this.model,
      this.isVisionModel,
      required this.ocrEngine,
      required this.promptId})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        name, r'UpsertReceiptProcessingSettingsCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        aiType, r'UpsertReceiptProcessingSettingsCommand', 'aiType');
    BuiltValueNullFieldError.checkNotNull(
        ocrEngine, r'UpsertReceiptProcessingSettingsCommand', 'ocrEngine');
    BuiltValueNullFieldError.checkNotNull(
        promptId, r'UpsertReceiptProcessingSettingsCommand', 'promptId');
  }

  @override
  UpsertReceiptProcessingSettingsCommand rebuild(
          void Function(UpsertReceiptProcessingSettingsCommandBuilder)
              updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertReceiptProcessingSettingsCommandBuilder toBuilder() =>
      new UpsertReceiptProcessingSettingsCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertReceiptProcessingSettingsCommand &&
        name == other.name &&
        description == other.description &&
        aiType == other.aiType &&
        url == other.url &&
        key == other.key &&
        model == other.model &&
        isVisionModel == other.isVisionModel &&
        ocrEngine == other.ocrEngine &&
        promptId == other.promptId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, aiType.hashCode);
    _$hash = $jc(_$hash, url.hashCode);
    _$hash = $jc(_$hash, key.hashCode);
    _$hash = $jc(_$hash, model.hashCode);
    _$hash = $jc(_$hash, isVisionModel.hashCode);
    _$hash = $jc(_$hash, ocrEngine.hashCode);
    _$hash = $jc(_$hash, promptId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(
            r'UpsertReceiptProcessingSettingsCommand')
          ..add('name', name)
          ..add('description', description)
          ..add('aiType', aiType)
          ..add('url', url)
          ..add('key', key)
          ..add('model', model)
          ..add('isVisionModel', isVisionModel)
          ..add('ocrEngine', ocrEngine)
          ..add('promptId', promptId))
        .toString();
  }
}

class UpsertReceiptProcessingSettingsCommandBuilder
    implements
        Builder<UpsertReceiptProcessingSettingsCommand,
            UpsertReceiptProcessingSettingsCommandBuilder> {
  _$UpsertReceiptProcessingSettingsCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(String? description) => _$this._description = description;

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

  bool? _isVisionModel;
  bool? get isVisionModel => _$this._isVisionModel;
  set isVisionModel(bool? isVisionModel) =>
      _$this._isVisionModel = isVisionModel;

  OcrEngine? _ocrEngine;
  OcrEngine? get ocrEngine => _$this._ocrEngine;
  set ocrEngine(OcrEngine? ocrEngine) => _$this._ocrEngine = ocrEngine;

  int? _promptId;
  int? get promptId => _$this._promptId;
  set promptId(int? promptId) => _$this._promptId = promptId;

  UpsertReceiptProcessingSettingsCommandBuilder() {
    UpsertReceiptProcessingSettingsCommand._defaults(this);
  }

  UpsertReceiptProcessingSettingsCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _description = $v.description;
      _aiType = $v.aiType;
      _url = $v.url;
      _key = $v.key;
      _model = $v.model;
      _isVisionModel = $v.isVisionModel;
      _ocrEngine = $v.ocrEngine;
      _promptId = $v.promptId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertReceiptProcessingSettingsCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertReceiptProcessingSettingsCommand;
  }

  @override
  void update(
      void Function(UpsertReceiptProcessingSettingsCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertReceiptProcessingSettingsCommand build() => _build();

  _$UpsertReceiptProcessingSettingsCommand _build() {
    final _$result = _$v ??
        new _$UpsertReceiptProcessingSettingsCommand._(
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UpsertReceiptProcessingSettingsCommand', 'name'),
            description: description,
            aiType: BuiltValueNullFieldError.checkNotNull(
                aiType, r'UpsertReceiptProcessingSettingsCommand', 'aiType'),
            url: url,
            key: key,
            model: model,
            isVisionModel: isVisionModel,
            ocrEngine: BuiltValueNullFieldError.checkNotNull(ocrEngine,
                r'UpsertReceiptProcessingSettingsCommand', 'ocrEngine'),
            promptId: BuiltValueNullFieldError.checkNotNull(promptId,
                r'UpsertReceiptProcessingSettingsCommand', 'promptId'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
