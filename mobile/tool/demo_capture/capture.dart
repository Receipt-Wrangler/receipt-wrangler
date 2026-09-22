import 'dart:io';
import 'dart:ui' as ui;

import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:image/image.dart' as img;

import 'fake_keyboard.dart';

/// Frame capture for the keyboard demos. See `mobile/CLAUDE.md` ->
/// "Recording a demo GIF".
///
/// The repo's other demo (`tool/record_tap_target_demo.sh`) records a real
/// Linux-desktop binary under Xvfb with ffmpeg. That cannot be reused here:
/// **there is no software keyboard on Linux desktop**, so `viewInsets.bottom`
/// is permanently 0 and the bug is literally unreproducible on that target.
/// `CLAUDE.md` already prescribes the alternative for exactly this case -- drive
/// it in a widget test and capture `RenderRepaintBoundary.toImage` per frame --
/// which is what this implements. It also needs no ffmpeg, ImageMagick,
/// xdotool or Android SDK, none of which the sandbox has.
///
/// **Fidelity.** The ramp is driven through `tester.view.viewInsets`, the same
/// engine-level field a real OS keyboard sets, so `Scaffold` and `EditableText`
/// (which reads `View.of`, not `MediaQuery`) both react exactly as on device.
/// The before/after panels are two recordings of the **real** screen, differing
/// only by `debugDisableKeyboardLift`, stitched side by side at encode time --
/// so neither side is a hand-copied imitation that could drift from the
/// shipping tree. Nothing about the layout is simulated; only the keyboard's
/// *pixels* are ours, because the engine reserves that space and the platform
/// is what normally paints into it.

const demoPhoneWidth = 390.0;
const demoPhoneHeight = 760.0;
const demoPhoneSize = Size(demoPhoneWidth, demoPhoneHeight);
const demoKeyboardHeight = 300.0;
const _labelBarHeight = 44.0;
const _panelGap = 16;

/// Height of one captured panel: the phone plus its caption strip.
const demoPanelHeight = demoPhoneHeight + _labelBarHeight;

final _boundaryKey = GlobalKey();

/// One step of the ramp: how far the keyboard is up, and how long to hold it.
///
/// Durations rather than duplicate hold frames, because `GifEncoder` writes
/// every frame in full (no inter-frame diffing like the ffmpeg pipeline), so a
/// held frame costs the same as a moving one. This is the main lever on size.
typedef DemoStep = ({double inset, int centis});

/// The keyboard's rise: a beat of context, an eased slide, then a long hold on
/// the outcome, which is the frame a reviewer actually reads.
List<DemoStep> demoInsetRamp() {
  const rise = 10;
  return <DemoStep>[
    (inset: 0.0, centis: 120),
    for (var i = 1; i <= rise; i++)
      (
        inset: demoKeyboardHeight * Curves.easeOutCubic.transform(i / rise),
        centis: 5,
      ),
    (inset: demoKeyboardHeight, centis: 260),
  ];
}

/// Makes text render as text rather than as boxes.
///
/// `flutter test` always passes `--use-test-fonts` to `flutter_tester`, whose
/// stub font draws every glyph as a filled rectangle -- so an unloaded
/// screenshot of a form is unreadable. Raleway is the app's real family; it is
/// *also* registered as `Roboto` because that is what Material's default
/// `Typography` asks for on the test binding's target platform, and most text
/// in these screens takes the theme default. MaterialIcons is needed for the
/// sheet's app-bar actions.
Future<void> loadDemoFonts() async {
  final raleway = await File('fonts/Raleway-Regular.ttf').readAsBytes();
  for (final family in const ['Raleway', 'Roboto']) {
    await (FontLoader(family)
          ..addFont(Future.value(ByteData.view(raleway.buffer))))
        .load();
  }

  // Repo-owned fonts come off disk; MaterialIcons ships with the SDK, so try
  // the asset bundle and fall back to the SDK cache.
  Future<ByteData> iconBytes() async {
    try {
      return await rootBundle.load('fonts/MaterialIcons-Regular.otf');
    } catch (_) {
      final root = Platform.environment['FLUTTER_ROOT'] ?? '/tmp/flutter';
      final file = File(
          '$root/bin/cache/artifacts/material_fonts/MaterialIcons-Regular.otf');
      return ByteData.view((await file.readAsBytes()).buffer);
    }
  }

  await (FontLoader('MaterialIcons')..addFont(iconBytes())).load();
}

/// Wraps [app] in the demo surface: a caption strip, a phone-sized viewport,
/// and a drawn keyboard filling whatever the current inset reserved.
///
/// The keyboard height is read from `MediaQuery` rather than passed in, so one
/// pumped tree animates as the ramp advances. Re-pumping per frame would
/// rebuild the app and close the very sheet being recorded.
///
/// The opaque [ColoredBox] is load-bearing: `rawRgba` is premultiplied, and
/// `GifEncoder` switches on GIF transparency the moment it finds a zero-alpha
/// palette entry, which makes the whole recording flicker.
Widget buildDemoSurface({
  required Widget app,
  required String label,
  required Color labelColor,
}) {
  return MediaQuery.fromView(
    view: WidgetsBinding.instance.platformDispatcher.views.first,
    child: Directionality(
      textDirection: TextDirection.ltr,
      child: RepaintBoundary(
        key: _boundaryKey,
        child: ColoredBox(
          color: Colors.white,
          child: Builder(
            builder: (BuildContext context) {
              final inset = MediaQuery.viewInsetsOf(context).bottom;
              return Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    height: _labelBarHeight,
                    width: demoPhoneSize.width,
                    color: labelColor,
                    alignment: Alignment.center,
                    child: Text(
                      label,
                      style: const TextStyle(
                        // Explicit: this caption is outside any MaterialApp, so
                        // there is no theme to inherit a family from and the
                        // test font would draw every glyph as a box.
                        fontFamily: 'Raleway',
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ),
                  SizedBox(
                    width: demoPhoneSize.width,
                    height: demoPhoneSize.height,
                    child: ClipRect(
                      child: Stack(
                        children: [
                          Positioned.fill(child: app),
                          Positioned(
                            left: 0,
                            right: 0,
                            bottom: 0,
                            child:
                                IgnorePointer(child: FakeKeyboard(height: inset)),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              );
            },
          ),
        ),
      ),
    ),
  );
}

/// Loads the demo fonts, failing the test if they do not load.
///
/// Call this rather than `tester.runAsync(loadDemoFonts)`. `loadDemoFonts`
/// returns `Future<void>`, and `runAsync` returns null **both** when the
/// callback throws and when it simply has nothing to return -- so a bare call
/// cannot tell the two apart and silently carries on. What follows is a whole
/// recording in `--use-test-fonts`' stub font, where every glyph is a filled
/// box: the demo passes and writes a GIF nobody can read.
///
/// Returning a sentinel from inside [WidgetTester.runAsync] is what makes the
/// failure detectable, and the shape matches [grabFrame] below, which has the
/// same problem for the same reason.
Future<void> loadDemoFontsOrFail(WidgetTester tester) async {
  final loaded = await tester.runAsync(() async {
    await loadDemoFonts();
    return true;
  });

  if (loaded != true) {
    fail('demo fonts failed to load: ${tester.takeException()}');
  }
}

/// Grabs the current frame off the repaint boundary.
///
/// `toImage` is asynchronous and `testWidgets` runs in fake-async, where the
/// engine never services its future -- so without [WidgetTester.runAsync] the
/// capture hangs until the test times out rather than failing. `runAsync` also
/// **swallows** any error and returns null, surfacing it through
/// `takeException`, hence the explicit `fail`.
Future<img.Image> grabFrame(WidgetTester tester) async {
  final frame = await tester.runAsync(() async {
    final boundary = _boundaryKey.currentContext!.findRenderObject()!
        as RenderRepaintBoundary;
    final image = await boundary.toImage();
    try {
      final data = (await image.toByteData(format: ui.ImageByteFormat.rawRgba))!;
      return img.Image.fromBytes(
        width: image.width,
        height: image.height,
        bytes: data.buffer,
        bytesOffset: data.offsetInBytes,
        numChannels: 4,
        order: img.ChannelOrder.rgba,
      );
    } finally {
      image.dispose();
    }
  });

  if (frame == null) {
    fail('frame capture failed: ${tester.takeException()}');
  }
  return frame;
}

/// Runs the inset ramp against an already-mounted tree, capturing every frame.
Future<List<img.Image>> recordInsetRamp(WidgetTester tester) async {
  final frames = <img.Image>[];
  for (final step in demoInsetRamp()) {
    tester.view.viewInsets = FakeViewPadding(bottom: step.inset);
    await tester.pump(const Duration(milliseconds: 40));
    frames.add(await grabFrame(tester));
  }
  return frames;
}

/// Lays two equally-sized panels side by side on an opaque canvas.
///
/// Opaque on purpose: `rawRgba` is premultiplied and `GifEncoder` turns on GIF
/// transparency the moment it finds a zero-alpha palette entry, which makes the
/// whole clip flicker.
img.Image stitchPanels(img.Image left, img.Image right) {
  final canvas = img.Image(
    width: left.width + _panelGap + right.width,
    height: left.height,
    numChannels: 3,
  );
  img.fill(canvas, color: img.ColorRgb8(255, 255, 255));
  img.compositeImage(canvas, left, dstX: 0, dstY: 0);
  img.compositeImage(canvas, right, dstX: left.width + _panelGap, dstY: 0);
  return canvas;
}

/// Writes a single before/after still.
///
/// Synchronous for the same reason [writeSideBySideGif] is: `testWidgets` runs
/// in fake-async, where `File.writeAsBytes` never completes and the test simply
/// hangs until it times out.
void writeSideBySidePng({
  required img.Image before,
  required img.Image after,
  required String path,
}) {
  final bytes = img.encodePng(stitchPanels(before, after));
  final file = File(path);
  file.parent.createSync(recursive: true);
  file.writeAsBytesSync(bytes);

  // ignore: avoid_print
  print('wrote $path — ${(bytes.length / 1024).toStringAsFixed(0)} KB');
}

/// Writes a single-panel looping GIF from frames that each carry their own
/// hold time, in **1/100 s**.
///
/// [writeSideBySideGif] below reads its durations off `demoInsetRamp()`, which
/// only makes sense for a demo driven by the keyboard inset. A demo driven by
/// *taps* has no ramp -- each step holds for as long as that step needs to be
/// read -- so the frames carry their own timing instead. The encoder settings
/// are the shared part and the reason this lives here rather than in a demo
/// file: they are not the defaults, and getting them wrong is not obvious until
/// the GIF is already committed (see `mobile/CLAUDE.md` -> "Recording a demo
/// GIF").
///
/// Synchronous for the same reason the others are: `testWidgets` runs in
/// fake-async, where `File.writeAsBytes` never completes and the test hangs
/// until it times out.
void writeGif({
  required List<({img.Image frame, int centis})> frames,
  required String path,
}) {
  final encoder = img.GifEncoder(
    repeat: 0,
    // 256, the GIF maximum, where the side-by-side writer below uses 128.
    // That one records flat UI art, which quantizes cleanly; a receipts list
    // does not -- `ListItemTrailingStatus` paints a LinearGradient in every
    // row, and with dithering off (which it must be, or LZW loses the long
    // pixel runs it compresses) a short palette bands that gradient into
    // visible stripes. It also cleans up antialiased white text on a saturated
    // caption strip, which ghosts at 128.
    numColors: 256,
    quantizerType: img.QuantizerType.octree,
    dither: img.DitherKernel.none,
  );

  for (final step in frames) {
    // addFrame encodes the PREVIOUS image and finish() flushes the last, so N
    // calls plus finish yield N frames.
    encoder.addFrame(step.frame, duration: step.centis);
  }

  final bytes = encoder.finish()!;
  final file = File(path);
  file.parent.createSync(recursive: true);
  file.writeAsBytesSync(bytes);

  // ignore: avoid_print
  print('wrote $path — ${(bytes.length / 1024).toStringAsFixed(0)} KB, '
      '${frames.length} frames, '
      '${frames.first.frame.width}x${frames.first.frame.height}');
}

/// Stitches the two recordings into one looping GIF, before on the left.
///
/// Side by side rather than sequential for the reason the tap-target demo is:
/// a strip asks the viewer to hold frame 3 in mind while looking at frame 20,
/// whereas two panels make the same keyboard rise against both layouts at once.
///
/// Synchronous, and that is deliberate: `testWidgets` runs in fake-async, where
/// `File.writeAsBytes` never completes — the test simply hangs until it times
/// out. Encoding is CPU-bound anyway, so sync I/O sidesteps the whole problem
/// without a `runAsync` wrapper.
/// [durationsCentis] gives one hold time per frame, in 1/100 s. It defaults to the
/// keyboard ramp's own timings, which is what the two keyboard demos want; any demo
/// with a different frame count MUST pass its own, or this indexes past the end of a
/// 13-step ramp (or silently applies keyboard timings to a scroll).
void writeSideBySideGif({
  required List<img.Image> before,
  required List<img.Image> after,
  required String path,
  List<int>? durationsCentis,
  int numColors = 128,
}) {
  assert(before.length == after.length, 'panels must have the same frame count');

  final encoder = img.GifEncoder(
    repeat: 0,
    // The defaults are wrong for this: the neural quantizer is a per-frame
    // neural net (slow), and dithering destroys the long runs of identical
    // pixels that LZW needs. Flat UI art quantizes cleanly without either.
    //
    // [numColors] is raised by a demo whose screen carries a GRADIENT -- at 128
    // the octree banding across ListItemTrailingStatus's status chips reads as
    // hard black stripes, i.e. as a rendering bug rather than as compression.
    // Raise the palette rather than switching the dither on, which would cost
    // far more bytes by breaking LZW's runs everywhere else.
    numColors: numColors,
    quantizerType: img.QuantizerType.octree,
    dither: img.DitherKernel.none,
  );

  final steps =
      durationsCentis ?? demoInsetRamp().map((step) => step.centis).toList();
  assert(steps.length == before.length, 'one duration per frame');
  for (var i = 0; i < before.length; i++) {
    final canvas = stitchPanels(before[i], after[i]);
    // Durations are in 1/100 s, not ms. Note addFrame encodes the PREVIOUS
    // image and finish() flushes the last, so N calls plus finish yield N
    // frames.
    encoder.addFrame(canvas, duration: steps[i]);
  }

  final bytes = encoder.finish()!;
  final file = File(path);
  file.parent.createSync(recursive: true);
  file.writeAsBytesSync(bytes);

  // ignore: avoid_print
  print('wrote $path — ${(bytes.length / 1024).toStringAsFixed(0)} KB, '
      '${before.length} frames, ${canvasLabel(before.first)}');
}

String canvasLabel(img.Image frame) =>
    '${frame.width * 2 + _panelGap}x${frame.height}';
