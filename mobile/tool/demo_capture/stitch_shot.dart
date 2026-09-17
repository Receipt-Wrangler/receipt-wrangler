import 'dart:io';

import 'package:image/image.dart' as img;

/// Lays two captured panels side by side and writes the result.
///
/// Plain `dart run`, not a widget test: this touches no Flutter binding, and a
/// test harness would only add the fake-async I/O trap for nothing.
///
///     dart run tool/demo_capture/stitch_shot.dart <before> <after> <out>
void main(List<String> args) {
  if (args.length != 3) {
    stderr.writeln('usage: stitch_shot.dart <before.png> <after.png> <out.png>');
    exit(64);
  }

  final left = img.decodePng(File(args[0]).readAsBytesSync())!;
  final right = img.decodePng(File(args[1]).readAsBytesSync())!;

  const gap = 16;
  final canvas = img.Image(
    width: left.width + gap + right.width,
    height: left.height,
    numChannels: 3,
  );
  img.fill(canvas, color: img.ColorRgb8(255, 255, 255));
  img.compositeImage(canvas, left, dstX: 0, dstY: 0);
  img.compositeImage(canvas, right, dstX: left.width + gap, dstY: 0);

  final bytes = img.encodePng(canvas);
  final out = File(args[2]);
  out.parent.createSync(recursive: true);
  out.writeAsBytesSync(bytes);

  stdout.writeln('wrote ${args[2]} — ${(bytes.length / 1024).toStringAsFixed(0)} KB, '
      '${canvas.width}x${canvas.height}');
}
