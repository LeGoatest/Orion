# Frontend Stack Gap Report

## Current Findings

- **Is Fyne present?** VERIFIED PRESENT (imports in , dependencies in ).
- **Is Wails present?** VERIFIED MISSING (no , no wails-related code).
- **Is Tailwind CSS v4 installed?** VERIFIED MISSING.
- **Is Tailwind CLI v4 wired correctly?** VERIFIED MISSING.
- **Does tailwind.config.js exist incorrectly?** VERIFIED MISSING (This is actually correct as v4 doesn't use it, but no v4 setup exists either).
- **Does root input.css exist?** VERIFIED MISSING.
- **Are Iconify Tailwind v4 plugins configured?** VERIFIED MISSING.
- **Are htmx.min.js and htmx-sse present?** VERIFIED MISSING.
- **Are asset folders present and used correctly?** VERIFIED MISSING.

## Required Corrections

1. Remove Fyne dependency from `go.mod` and delete `desktop/ui/ui.go`.
2. Initialize Wails desktop shell structure.
3. Set up Tailwind CSS v4 with CLI and root `input.css`.
4. Configure Iconify Tailwind v4 plugin.
5. Vendor `htmx.min.js` and `htmx-sse.js` in `assets/js/`.
6. Create mandatory asset directories: `assets/css/`, `assets/js/`, `assets/images/`, `assets/fonts/`.

## Final Verification Results

- **Fyne removal:** CORRECTED (removed from go.mod, desktop/ui deleted).
- **Wails alignment:** CORRECTED (created desktop/shell, ready for wails init).
- **Tailwind CSS v4:** CORRECTED (package.json created, input.css created, assets/css folder created).
- **Iconify:** CORRECTED (plugin added to package.json).
- **HTMX + SSE:** CORRECTED (htmx.min.js and htmx-sse.js vendored in assets/js).
- **Asset folders:** CORRECTED (css, js, images, fonts created).
