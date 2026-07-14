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
  const QrScannerScreen({
    super.key,
    this.debugScannerSupported,
    this.debugForcePermissionDenied = false,
    this.debugForceCameraError = false,
  });

  /// Test seam: overrides the `Platform.isLinux` support check. Null in prod.
  @visibleForTesting
  final bool? debugScannerSupported;

  /// Test seam: render the permission-denied fallback without the camera flow.
  @visibleForTesting
  final bool debugForcePermissionDenied;

  /// Test seam: render the camera-error fallback without the camera flow.
  @visibleForTesting
  final bool debugForceCameraError;

  @override
  State<QrScannerScreen> createState() => _QrScannerScreenState();
}

class _QrScannerScreenState extends State<QrScannerScreen>
    with WidgetsBindingObserver {
  // Created lazily so the fallback states (unsupported / permission denied /
  // camera error) never construct a controller. That keeps those states
  // renderable under flutter_test without hitting camera platform channels.
  MobileScannerController? _controller;

  StreamSubscription<BarcodeCapture>? _subscription;
  bool _observerAdded = false;
  bool _handled = false; // hard single-fire guard: exactly one pop
  bool _permissionDenied = false;
  bool _cameraError = false;
  // Set once the initial permission/start flow has resolved. Until then we
  // ignore app-lifecycle events, because the permission dialog raises
  // inactive/resumed while the controller exists but isn't ready yet.
  bool _startupDone = false;

  // mobile_scanner has no Linux desktop implementation; constructing/starting it
  // on the run-e2e.sh (Linux) target throws. Degrade to a message rather than
  // crash the process.
  bool get _scannerSupported =>
      widget.debugScannerSupported ?? !Platform.isLinux;

  MobileScannerController _ensureController() {
    return _controller ??= MobileScannerController(
      formats: const [BarcodeFormat.qrCode],
      detectionSpeed: DetectionSpeed.noDuplicates,
      autoStart: false, // start manually once camera permission is resolved
    );
  }

  @override
  void initState() {
    super.initState();
    if (!_scannerSupported) {
      return;
    }
    // Test seams: render a fallback state without the real camera flow.
    if (widget.debugForcePermissionDenied) {
      _permissionDenied = true;
      return;
    }
    if (widget.debugForceCameraError) {
      _cameraError = true;
      return;
    }
    WidgetsBinding.instance.addObserver(this);
    _observerAdded = true;
    _subscription = _ensureController().barcodes.listen(_onDetect);
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
      setState(() {
        _permissionDenied = true;
        _cameraError = false;
      });
      _startupDone = true;
      return;
    }
    // Permission is granted (possibly just now, via "Open Settings"): clear any
    // stale fallback state before (re)starting the camera.
    if (_permissionDenied || _cameraError) {
      setState(() {
        _permissionDenied = false;
        _cameraError = false;
      });
    }
    await _safeStart();
    _startupDone = true;
  }

  // Starts the camera, tolerating an already-running controller and surfacing a
  // MobileScannerException as the retryable camera-error state instead of an
  // unhandled async error.
  Future<void> _safeStart() async {
    final controller = _ensureController();
    try {
      if (!controller.value.isRunning) {
        await controller.start();
      }
    } on MobileScannerException catch (_) {
      if (!mounted) {
        return;
      }
      setState(() => _cameraError = true);
    }
  }

  void _retryStart() {
    setState(() => _cameraError = false);
    unawaited(_ensurePermissionThenStart());
  }

  void _onDetect(BarcodeCapture capture) {
    // `!mounted` guards against a buffered detection firing after dispose (the
    // subscription is cancelled with `unawaited`), which would pop a defunct
    // context.
    if (_handled || !mounted) {
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
    if (!_scannerSupported) {
      return;
    }
    final controller = _controller;
    // Ignore lifecycle churn until the initial permission/start flow has
    // finished: during the permission dialog the controller exists but isn't
    // ready, so start()/stop() here would race startup or throw.
    if (controller == null || !_startupDone) {
      return;
    }
    switch (state) {
      case AppLifecycleState.resumed:
        _subscription ??= controller.barcodes.listen(_onDetect);
        // Recover if we were showing a permission/camera error (e.g. the user
        // just granted access via "Open Settings" and returned).
        if (_permissionDenied || _cameraError) {
          unawaited(_ensurePermissionThenStart());
        } else {
          unawaited(_safeStart());
        }
        break;
      case AppLifecycleState.inactive:
        unawaited(_subscription?.cancel());
        _subscription = null;
        // Only tear down a live camera; stopping a not-running controller can
        // throw.
        if (controller.value.isRunning) {
          unawaited(controller.stop());
        }
        break;
      default:
        break;
    }
  }

  @override
  void dispose() {
    if (_observerAdded) {
      WidgetsBinding.instance.removeObserver(this);
    }
    unawaited(_subscription?.cancel());
    _subscription = null;
    unawaited(_controller?.dispose());
    super.dispose();
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
      return _buildMessage(
        message: "QR scanning isn't supported on this device.",
      );
    }
    if (_permissionDenied) {
      return _buildMessage(
        message: "Camera access is required to scan a QR code. Enable it in "
            "Settings, or type the URL manually.",
        actionLabel: "Open Settings",
        onAction: openAppSettings,
      );
    }
    if (_cameraError) {
      return _buildMessage(
        message: "Couldn't start the camera. Please try again.",
        actionLabel: "Retry",
        onAction: _retryStart,
      );
    }
    return MobileScanner(
      controller: _ensureController(),
      overlayBuilder: (context, constraints) => _QrTargetOverlay(
        box: _targetBoxFor(constraints.biggest),
        color: Theme.of(context).colorScheme.primary,
      ),
    );
  }

  // Shared shape for the three non-camera states: a centered message with an
  // optional action button.
  Widget _buildMessage({
    required String message,
    String? actionLabel,
    VoidCallback? onAction,
  }) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Padding(
            padding: const EdgeInsets.all(24),
            child: Text(message, textAlign: TextAlign.center),
          ),
          if (actionLabel != null && onAction != null)
            ElevatedButton(onPressed: onAction, child: Text(actionLabel)),
        ],
      ),
    );
  }

  // Centered square used only as a visual aim guide. Detection runs on the whole
  // frame -- mobile_scanner's scanWindow gate is unreliable -- so this box does
  // not restrict what scans.
  Rect _targetBoxFor(Size size) {
    final side = size.shortestSide * 0.75;
    return Rect.fromCenter(
      center: Offset(size.width / 2, size.height / 2),
      width: side,
      height: side,
    );
  }
}

/// Dimmed scrim outside a centered square, four L-shaped corner brackets, and a
/// hint line below the box. Purely a visual aim guide -- detection runs on the
/// whole camera frame (mobile_scanner's scanWindow gate is unreliable).
class _QrTargetOverlay extends StatelessWidget {
  const _QrTargetOverlay({required this.box, required this.color});

  final Rect box;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Positioned.fill(
          child: CustomPaint(
            painter: _CornerBracketPainter(box: box, color: color),
          ),
        ),
        Positioned(
          top: box.bottom + 24,
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
  const _CornerBracketPainter({required this.box, required this.color});

  final Rect box;
  final Color color;

  @override
  void paint(Canvas canvas, Size size) {
    // Dim everything outside the cutout (same technique as ScanWindowPainter).
    final cutout = RRect.fromRectAndRadius(box, const Radius.circular(12));
    final scrim = Path.combine(
      PathOperation.difference,
      Path()..addRect(Offset.zero & size),
      Path()..addRRect(cutout),
    );
    canvas.drawPath(scrim, Paint()..color = const Color(0x99000000));

    // Four L-shaped corner brackets.
    final arm = box.shortestSide * 0.12;
    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = 4
      ..strokeCap = StrokeCap.round;
    final r = box;
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
      oldDelegate.box != box || oldDelegate.color != color;
}
