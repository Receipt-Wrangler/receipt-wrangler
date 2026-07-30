# App Links / Universal Links hosting for `receiptwrangler.io`

These files power the **login-QR deep link**: the QR on the desktop login page encodes
`https://receiptwrangler.io/app/setup#url=<server url>`. Scanning it opens the mobile app (App Links
on Android / Universal Links on iOS) and pre-fills the server URL; if the app isn't installed, the
link opens the `/app/setup` landing page instead. See:

- `api/CLAUDE.md` → "Login QR & mobile deep link"
- `mobile/CLAUDE.md` → "App Links / Universal Links — server-URL pre-fill (login)"

**These are hosted on the project's own `receiptwrangler.io` domain — not on a self-hoster's
instance.** The published Play/App Store builds verify against `receiptwrangler.io`; the self-hoster's
actual server URL only travels inside the link's fragment.

## What to deploy

| File | Serve at | Notes |
|---|---|---|
| `.well-known/assetlinks.json` | `https://receiptwrangler.io/.well-known/assetlinks.json` | Android App Links (Digital Asset Links) |
| `.well-known/apple-app-site-association` | `https://receiptwrangler.io/.well-known/apple-app-site-association` | iOS Universal Links. **No file extension.** |
| `app/setup/index.html` | `https://receiptwrangler.io/app/setup` | "Get the app" fallback when the app isn't installed |

## Serving requirements (App/Universal Links will silently fail otherwise)

- **HTTPS**, `Content-Type: application/json` for both `.well-known` files.
- **No redirects** to reach the `.well-known` files (Apple/Google fetch them directly).
- The `.well-known` files must be reachable **without** any auth wall.
- Publish the `.well-known` files **before** releasing the mobile app builds — the OS verifies the
  domain↔app association around install/update time.

## Identifiers baked into these files

- Android `package_name` = `io.receiptwrangler` (the installed `applicationId`), with both the
  Play **app-signing** and **upload-key** SHA-256 fingerprints so Play-installed and local builds
  both verify.
- iOS `appID` = `3VD3YNZ3KA.io.receiptwrangler` (Apple Team ID + bundle id).

## Before you ship

- Replace the App Store placeholder link in `app/setup/index.html`
  (`https://apps.apple.com/app/id0000000000`) with the real numeric App Store ID.
- Validate with Google's **Statement List Tester** / Digital Asset Links API and Apple's
  **AASA validator** before the store release.
