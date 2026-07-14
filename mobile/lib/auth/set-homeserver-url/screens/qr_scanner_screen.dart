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
    return MobileScanner(controller: _controller);
  }
}
