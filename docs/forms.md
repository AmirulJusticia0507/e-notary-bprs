# Form Specifications & Validation

## Layout Architecture

- Form disajikan dalam bentuk **Multi-step Wizard** (`Steppers.vue`).
- Floating labels atau top-aligned labels dengan `text-xs font-semibold text-slate-600`.

## Input Rules & Stack Integration

- **State Management:** Reactive Form State via Vue 3 `reactive()`.
- **Validation:** VeeValidate / Zod schema validation.
- **Auto-format:** Currency masking (Rupiah) & Real-time File Upload preview (`.pdf` max 10MB).

## Tailwind Form Components

### Input Text Standard

```html
<div class="space-y-1">
  <label class="block text-xs font-semibold text-slate-700">Nomor SK Notaris</label>
  <input type="text" class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500" placeholder="SK-123/KEMENKUMHAM..." />
</div>
```