// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_processing_settings.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptProcessingSettings extends ReceiptProcessingSettings {
  @override
  final OcrEngine? ocrEngine;
  @override
  final bool? isVisionModel;
  @override
  final AiType? aiType;
  @override
  final int? promptId;
  @override
  final String? name;
  @override
  final String? description;
  @override
  final String? model;
  @override
  final Prompt? prompt;
  @override
  final String? url;
  @override
  final String? key;
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

  factory _$ReceiptProcessingSettings(
          [void Function(ReceiptProcessingSettingsBuilder)? updates]) =>
      (new ReceiptProcessingSettingsBuilder()..update(updates))._build();

  _$ReceiptProcessingSettings._(
      {this.ocrEngine,
      this.isVisionModel,
      this.aiType,
      this.promptId,
      this.name,
      this.description,
      this.model,
      this.prompt,
      this.url,
      this.key,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        id, r'ReceiptProcessingSettings', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'ReceiptProcessingSettings', 'createdAt');
  }

  @override
  ReceiptProcessingSettings rebuild(
          void Function(ReceiptProcessingSettingsBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptProcessingSettingsBuilder toBuilder() =>
      new ReceiptProcessingSettingsBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptProcessingSettings &&
        ocrEngine == other.ocrEngine &&
        isVisionModel == other.isVisionModel &&
        aiType == other.aiType &&
        promptId == other.promptId &&
        name == other.name &&
        description == other.description &&
        model == other.model &&
        prompt == other.prompt &&
        url == other.url &&
        key == other.key &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, ocrEngine.hashCode);
    _$hash = $jc(_$hash, isVisionModel.hashCode);
    _$hash = $jc(_$hash, aiType.hashCode);
    _$hash = $jc(_$hash, promptId.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, description.hashCode);
    _$hash = $jc(_$hash, model.hashCode);
    _$hash = $jc(_$hash, prompt.hashCode);
    _$hash = $jc(_$hash, url.hashCode);
    _$hash = $jc(_$hash, key.hashCode);
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
    return (newBuiltValueToStringHelper(r'ReceiptProcessingSettings')
          ..add('ocrEngine', ocrEngine)
          ..add('isVisionModel', isVisionModel)
          ..add('aiType', aiType)
          ..add('promptId', promptId)
          ..add('name', name)
          ..add('description', description)
          ..add('model', model)
          ..add('prompt', prompt)
          ..add('url', url)
          ..add('key', key)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class ReceiptProcessingSettingsBuilder
    implements
        Builder<ReceiptProcessingSettings, ReceiptProcessingSettingsBuilder>,
        BaseModelBuilder {
  _$ReceiptProcessingSettings? _$v;

  OcrEngine? _ocrEngine;
  OcrEngine? get ocrEngine => _$this._ocrEngine;
  set ocrEngine(covariant OcrEngine? ocrEngine) =>
      _$this._ocrEngine = ocrEngine;

  bool? _isVisionModel;
  bool? get isVisionModel => _$this._isVisionModel;
  set isVisionModel(covariant bool? isVisionModel) =>
      _$this._isVisionModel = isVisionModel;

  AiType? _aiType;
  AiType? get aiType => _$this._aiType;
  set aiType(covariant AiType? aiType) => _$this._aiType = aiType;

  int? _promptId;
  int? get promptId => _$this._promptId;
  set promptId(covariant int? promptId) => _$this._promptId = promptId;

  String? _name;
  String? get name => _$this._name;
  set name(covariant String? name) => _$this._name = name;

  String? _description;
  String? get description => _$this._description;
  set description(covariant String? description) =>
      _$this._description = description;

  String? _model;
  String? get model => _$this._model;
  set model(covariant String? model) => _$this._model = model;

  PromptBuilder? _prompt;
  PromptBuilder get prompt => _$this._prompt ??= new PromptBuilder();
  set prompt(covariant PromptBuilder? prompt) => _$this._prompt = prompt;

  String? _url;
  String? get url => _$this._url;
  set url(covariant String? url) => _$this._url = url;

  String? _key;
  String? get key => _$this._key;
  set key(covariant String? key) => _$this._key = key;

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

  ReceiptProcessingSettingsBuilder() {
    ReceiptProcessingSettings._defaults(this);
  }

  ReceiptProcessingSettingsBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _ocrEngine = $v.ocrEngine;
      _isVisionModel = $v.isVisionModel;
      _aiType = $v.aiType;
      _promptId = $v.promptId;
      _name = $v.name;
      _description = $v.description;
      _model = $v.model;
      _prompt = $v.prompt?.toBuilder();
      _url = $v.url;
      _key = $v.key;
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
  void replace(covariant ReceiptProcessingSettings other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ReceiptProcessingSettings;
  }

  @override
  void update(void Function(ReceiptProcessingSettingsBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptProcessingSettings build() => _build();

  _$ReceiptProcessingSettings _build() {
    _$ReceiptProcessingSettings _$result;
    try {
      _$result = _$v ??
          new _$ReceiptProcessingSettings._(
              ocrEngine: ocrEngine,
              isVisionModel: isVisionModel,
              aiType: aiType,
              promptId: promptId,
              name: name,
              description: description,
              model: model,
              prompt: _prompt?.build(),
              url: url,
              key: key,
              id: BuiltValueNullFieldError.checkNotNull(
                  id, r'ReceiptProcessingSettings', 'id'),
              createdAt: BuiltValueNullFieldError.checkNotNull(
                  createdAt, r'ReceiptProcessingSettings', 'createdAt'),
              createdBy: createdBy,
              createdByString: createdByString,
              updatedAt: updatedAt);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'prompt';
        _prompt?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'ReceiptProcessingSettings', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
