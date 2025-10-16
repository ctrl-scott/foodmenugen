# Menu Generator — Web + Go/Rust Hooks

A single‑file **web‑based restaurant menu generator** that runs entirely in the browser and can connect to optional **Go or Rust backends** for storage, rendering, or PDF generation.


ChatGPT Link: https://chatgpt.com/share/68f14002-f474-800c-acf3-16528285f344

## ✨ Features

- 🧩 **Single HTML file**, fully offline‑friendly (no dependencies)
- 🖥️ **Builder interface** with live preview
- 💾 Import/export menus as JSON or self‑contained HTML
- 🖨️ Print‑ready with selectable layouts and page presets
- ⚙️ **Go and Rust backend hooks** for server‑side rendering or storage
- 🎨 Layout & Design:
  - Classic (single column)
  - Two‑column (grid per section)
  - Bordered (decorative SVG dividers)
  - Bistro (leaders between name and price)
- 🖼️ Logo upload and positioning (None / Left / Center)
- 🧭 SVG ornament themes: Fork & Knife, Laurel, Wave, or None
- 🪙 **Flexible currency system** with ISO or symbol display, placement options, and custom symbol support
- 📏 Adjustable logo size, menu width, column count, and leader density

## 🧰 Technologies

- Frontend: HTML5, CSS3, and Vanilla JavaScript
- Backend: Go (net/http + template) or Rust (Axum + Askama)
- No external libraries required

## 🚀 Usage

1. **Open `menu.html` in a browser** — it works offline.
2. Use the left‑hand builder panel to:
   - Edit restaurant info, layout, ornaments, and pricing.
   - Add or remove sections and items.
3. **Export** using the buttons in the top bar:
   - `Export JSON` → Saves your editable data.
   - `Export HTML` → Creates a printable standalone menu.
   - `Print` → Opens print dialog (A4 or Letter).

## 🖧 Optional Backend Setup

You can connect to a local Go or Rust backend to persist data or generate server‑side PDFs.

### Go Example
```bash
go run .  # runs at http://localhost:8080
```
Set **Backend base URL** in the app to `http://localhost:8080`.

Endpoints:
- `POST /api/menu` → Store menu JSON
- `POST /api/render` → Render HTML and return `{ html }`

### Rust Example
```bash
cargo run  # runs at http://localhost:8080
```
Endpoints identical to Go version.

## 🧪 Self‑Tests
Use the **Run Self‑Tests** button to verify:
- DOM builder safety
- Layout and SVG injection
- Bistro leader rendering
- Logo injection
- Currency formatting logic
- Print column toggling

## 🧩 Customization
- Edit CSS in the `<style>` block for colors, borders, and typography.
- Extend `ornamentSVG()` for new decorative SVG themes.
- Modify the Go/Rust templates to include restaurant branding or analytics.

## 📜 License
Creative Commons Attribution‑ShareAlike 4.0 International (CC BY‑SA 4.0)

## 🧾 Credits
Developed collaboratively with ChatGPT‑5 (OpenAI) and Scott Owen — 2025.

---
**Tip:** For deploying as part of a local restaurant intranet or kiosk, serve the HTML file and backend binaries via Nginx or Caddy for offline resilience.
