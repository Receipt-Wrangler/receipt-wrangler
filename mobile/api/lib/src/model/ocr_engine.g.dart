// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'ocr_engine.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const OcrEngine _$TESSERACT = const OcrEngine._('TESSERACT');
const OcrEngine _$EASY_OCR = const OcrEngine._('EASY_OCR');

OcrEngine _$valueOf(String name) {
  switch (name) {
    case 'TESSERACT':
      return _$TESSERACT;
    case 'EASY_OCR':
      return _$EASY_OCR;
    default:
      throw new ArgumentError(name);
  }
}

final BuiltSet<OcrEngine> _$values = new BuiltSet<OcrEngine>(const <OcrEngine>[
  _$TESSERACT,
  _$EASY_OCR,
]);

class _$OcrEngineMeta {
  const _$OcrEngineMeta();
  OcrEngine get TESSERACT => _$TESSERACT;
  OcrEngine get EASY_OCR => _$EASY_OCR;
  OcrEngine valueOf(String name) => _$valueOf(name);
  BuiltSet<OcrEngine> get values => _$values;
}

abstract class _$OcrEngineMixin {
  // ignore: non_constant_identifier_names
  _$OcrEngineMeta get OcrEngine => const _$OcrEngineMeta();
}

Serializer<OcrEngine> _$ocrEngineSerializer = new _$OcrEngineSerializer();

class _$OcrEngineSerializer implements PrimitiveSerializer<OcrEngine> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'TESSERACT': 'TESSERACT',
    'EASY_OCR': 'EASY_OCR',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'TESSERACT': 'TESSERACT',
    'EASY_OCR': 'EASY_OCR',
  };

  @override
  final Iterable<Type> types = const <Type>[OcrEngine];
  @override
  final String wireName = 'OcrEngine';

  @override
  Object serialize(Serializers serializers, OcrEngine object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  OcrEngine deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      OcrEngine.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
