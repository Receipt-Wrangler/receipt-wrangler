import 'package:flutter/material.dart';

/// Draws a keyboard.
///
/// Setting `viewInsets` moves the *layout* but paints nothing — the engine
/// reserves the space, the platform draws the keys. A recording of the fix with
/// nothing in that space is unreadable: the bar just appears to jump for no
/// reason. So the demo paints its own keys into exactly the region the inset
/// reserved.
///
/// Purely decorative, and deliberately inert (wrapped in [IgnorePointer] by the
/// harness) — nothing here participates in the layout being demonstrated.
class FakeKeyboard extends StatelessWidget {
  const FakeKeyboard({super.key, required this.height});

  final double height;

  static const _rows = <String>[
    'q w e r t y u i o p',
    'a s d f g h j k l',
    'z x c v b n m',
  ];

  @override
  Widget build(BuildContext context) {
    if (height <= 0) {
      return const SizedBox.shrink();
    }

    return Container(
      height: height,
      width: double.infinity,
      color: const Color(0xFFD1D5DB),
      padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 8),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          for (final row in _rows)
            Expanded(
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  for (final key in row.split(' '))
                    Expanded(
                      child: Container(
                        margin: const EdgeInsets.all(2),
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(4),
                        ),
                        alignment: Alignment.center,
                        child: Text(
                          key,
                          style: const TextStyle(
                            // Explicit: the keyboard is painted outside any
                            // MaterialApp, so there is no theme to inherit a
                            // family from and `--use-test-fonts` would draw
                            // every key as a filled box.
                            fontFamily: 'Raleway',
                            fontSize: 13,
                            color: Color(0xFF111827),
                          ),
                        ),
                      ),
                    ),
                ],
              ),
            ),
          // The space bar row, which is most of what makes the block read as a
          // keyboard rather than a grid.
          Expanded(
            child: Container(
              margin: const EdgeInsets.symmetric(horizontal: 40, vertical: 2),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(4),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
