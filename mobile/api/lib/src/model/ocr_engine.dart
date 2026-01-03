//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'ocr_engine.g.dart';

class OcrEngine extends EnumClass {

  @BuiltValueEnumConst(wireName: r'TESSERACT')
  static const OcrEngine TESSERACT = _$TESSERACT;
  @BuiltValueEnumConst(wireName: r'EASY_OCR')
  static const OcrEngine EASY_OCR = _$EASY_OCR;

  static Serializer<OcrEngine> get serializer => _$ocrEngineSerializer;

  const OcrEngine._(String name): super(name);

  static BuiltSet<OcrEngine> get values => _$values;
  static OcrEngine valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class OcrEngineMixin = Object with _$OcrEngineMixin;

