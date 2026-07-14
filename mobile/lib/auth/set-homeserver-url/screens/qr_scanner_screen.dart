import 'dart:async';
import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:receipt_wrangler_mobile/utils/permissions.dart';

/// Full-screen QR scanner used to fill the server URL on the "Connect to Server"
/// screen. Pops the raw scanned string back to the caller (validation happens
/// there); pops null when the user cancels or camera access is unavailable.
class QrScannerScreen extends StatefulWidget {
  const QrScannerScreen({super.key});

  @override
  State<QrScannerScreen> createState() => _QrScannerScreenState();
}

class _QrScannerScreenState extends State<QrScannerScreen>
    with WidgetsBindingObserver {
  final MobileScannerController _controller = MobileScannerController(
    formats: const [BarcodeFormat.qrCode],
    detectionSpeed: DetectionSpeed.noDuplicates,
    autoStart: false, // start manually once camera permission is resolved
  );

  StreamSubscription<BarcodeCapture>? _subscription;
  bool _handled = false; // hard single-fire guard: exactly one pop
  bool _permissionDenied = false;

  // mobile_scanner has no Linux desktop implementation; constructing/starting it
  // on the run-e2e.sh (Linux) target throws. Degrade to a message rather than
  // crash the process.
  bool get _scannerSupported => !Platform.isLinux;

  @override
  void initState() {
    super.initState();
    if (!_scannerSupported) {
      return;
    }
    WidgetsBinding.instance.addObserver(this);
    _subscription = _controller.barcodes.listen(_onDetect);
    _ensurePermissionThenStart();
  }

  Future<void> _ensurePermissionThenStart() async {
    // Reuse the app-startup camera request (a shared in-flight future) so we
    // never fire a second *concurrent* request, which would throw
    // ERROR_ALREADY_REQUESTING_PERMISSIONS on Android.
    await requestPermissions();
    final status = await Permission.camera.status;
    if (!mounted) {
      return;
    }
    if (!status.isGranted) {
      setState(() => _permissionDenied = true);
      return;
    }
    await _controller.start();
  }

  void _onDetect(BarcodeCapture capture) {
    if (_handled) {
      return;
    }
    final raw =
        capture.barcodes.isNotEmpty ? capture.barcodes.first.rawValue : null;
    if (raw == null || raw.isEmpty) {
      return;
    }
    _handled = true;
    Navigator.of(context).pop(raw);
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (!_scannerSupported || _permissionDenied) {
      return;
    }
    switch (state) {
      case AppLifecycleState.resumed:
        _subscription = _controller.barcodes.listen(_onDetect);
        unawaited(_controller.start());
        break;
      case AppLifecycleState.inactive:
        unawaited(_subscription?.cancel());
        _subscription = null;
        unawaited(_controller.stop());
        break;
      default:
        break;
    }
  }

  @override
  Future<void> dispose() async {
    if (_scannerSupported) {
      WidgetsBinding.instance.removeObserver(this);
    }
    unawaited(_subscription?.cancel());
    _subscription = null;
    super.dispose();
    await _controller.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Scan server QR code"),
        leading: IconButton(
          key: const ValueKey("qr-scanner-close"),
          icon: const Icon(Icons.close),
          onPressed: () => Navigator.of(context).pop(), // null => cancel
        ),
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (!_scannerSupported) {
      return const Center(
        child: Text("QR scanning isn't supported on this device."),
      );
    }
    if (_permissionDenied) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Padding(
              padding: EdgeInsets.all(24),
              child: Text(
                "Camera access is required to scan a QR code. Enable it in "
                "Settings, or type the URL manually.",
                textAlign: TextAlign.center,
              ),
            ),
            ElevatedButton(
              onPressed: openAppSettings,
              child: const Text("Open Settings"),
            ),
          ],
        ),
      );
    }
    return LayoutBuilder(
      builder: (context, constraints) {
        final scanWindow = _scanWindowFor(constraints.biggest);
        return MobileScanner(
          controller: _controller,
          scanWindow: scanWindow, // gate detection to the box
          overlayBuilder: (context, _) => _QrTargetOverlay(
            scanWindow: scanWindow,
            color: Theme.of(context).colorScheme.primary,
          ),
        );
      },
    );
  }

  // Centered square, 70% of the shorter side. The same Rect drives both the
  // detection window and the drawn box, so they stay aligned.
  Rect _scanWindowFor(Size size) {
    final side = size.shortestSide * 0.7;
    return Rect.fromCenter(
      center: Offset(size.width / 2, size.height / 2),
      width: side,
      height: side,
    );
  }
}

/// Dimmed scrim outside a centered square, four L-shaped corner brackets, and a
/// hint line below the box. Drawn from the same [scanWindow] Rect passed to
/// [MobileScanner.scanWindow] so the box matches the detection area.
class _QrTargetOverlay extends StatelessWidget {
  const _QrTargetOverlay({required this.scanWindow, required this.color});

  final Rect scanWindow;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Positioned.fill(
          child: CustomPaint(
            painter: _CornerBracketPainter(scanWindow: scanWindow, color: color),
          ),
        ),
        Positioned(
          top: scanWindow.bottom + 24,
          left: 24,
          right: 24,
          child: const Text(
            "Point your camera at the server's QR code",
            textAlign: TextAlign.center,
            style: TextStyle(color: Colors.white, fontSize: 16),
          ),
        ),
      ],
    );
  }
}

class _CornerBracketPainter extends CustomPainter {
  const _CornerBracketPainter({required this.scanWindow, required this.color});

  final Rect scanWindow;
  final Color color;

  @override
  void paint(Canvas canvas, Size size) {
    // Dim everything outside the cutout (same technique as ScanWindowPainter).
    final cutout = RRect.fromRectAndRadius(scanWindow, const Radius.circular(12));
    final scrim = Path.combine(
      PathOperation.difference,
      Path()..addRect(Offset.zero & size),
      Path()..addRRect(cutout),
    );
    canvas.drawPath(scrim, Paint()..color = const Color(0x99000000));

    // Four L-shaped corner brackets.
    final arm = scanWindow.shortestSide * 0.12;
    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = 4
      ..strokeCap = StrokeCap.round;
    final r = scanWindow;
    canvas
      ..drawPath(
          Path()
            ..moveTo(r.left, r.top + arm)
            ..lineTo(r.left, r.top)
            ..lineTo(r.left + arm, r.top),
          paint)
      ..drawPath(
          Path()
            ..moveTo(r.right - arm, r.top)
            ..lineTo(r.right, r.top)
            ..lineTo(r.right, r.top + arm),
          paint)
      ..drawPath(
          Path()
            ..moveTo(r.left, r.bottom - arm)
            ..lineTo(r.left, r.bottom)
            ..lineTo(r.left + arm, r.bottom),
          paint)
      ..drawPath(
          Path()
            ..moveTo(r.right - arm, r.bottom)
            ..lineTo(r.right, r.bottom)
            ..lineTo(r.right, r.bottom - arm),
          paint);
  }

  @override
  bool shouldRepaint(_CornerBracketPainter oldDelegate) =>
      oldDelegate.scanWindow != scanWindow || oldDelegate.color != color;
}
