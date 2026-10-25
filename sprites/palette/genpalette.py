#!/usr/bin/env python3
"""Generate the Vex pixel-art palette from the Material Design colors.

The palette is the Material swatches (shades 900..50 and accents
A700..A100) from MUI's color source, plus a red grey and a green grey. These
two are Material blue grey with a different OKLCH hue. Each shade keeps the
lightness and chroma of the blue grey shade.

Outputs (in --out, named after --name): <name>.hex (Lospec), <name>.gpl
(GIMP/Aseprite), <name>.json, <name>.png (1 px per color), and
<name>-preview.png (scaled).

Run from anywhere:  python3 sprites/palette/genpalette.py
"""

import argparse
import json
import math
import struct
import zlib
from pathlib import Path

# name: (shades 900..50, accents A700 A400 A200 A100). The grey, brown and blue
# grey "accents" in MUI repeat their shades, so they are left out.
MATERIAL_SWATCHES = {
    "red": ("B71C1C C62828 D32F2F E53935 F44336 EF5350 E57373 EF9A9A FFCDD2 FFEBEE",
            "D50000 FF1744 FF5252 FF8A80"),
    "pink": ("880E4F AD1457 C2185B D81B60 E91E63 EC407A F06292 F48FB1 F8BBD0 FCE4EC",
             "C51162 F50057 FF4081 FF80AB"),
    "purple": ("4A148C 6A1B9A 7B1FA2 8E24AA 9C27B0 AB47BC BA68C8 CE93D8 E1BEE7 F3E5F5",
               "AA00FF D500F9 E040FB EA80FC"),
    "deep-purple": ("311B92 4527A0 512DA8 5E35B1 673AB7 7E57C2 9575CD B39DDB D1C4E9 EDE7F6",
                    "6200EA 651FFF 7C4DFF B388FF"),
    "indigo": ("1A237E 283593 303F9F 3949AB 3F51B5 5C6BC0 7986CB 9FA8DA C5CAE9 E8EAF6",
               "304FFE 3D5AFE 536DFE 8C9EFF"),
    "blue": ("0D47A1 1565C0 1976D2 1E88E5 2196F3 42A5F5 64B5F6 90CAF9 BBDEFB E3F2FD",
             "2962FF 2979FF 448AFF 82B1FF"),
    "light-blue": ("01579B 0277BD 0288D1 039BE5 03A9F4 29B6F6 4FC3F7 81D4FA B3E5FC E1F5FE",
                   "0091EA 00B0FF 40C4FF 80D8FF"),
    "cyan": ("006064 00838F 0097A7 00ACC1 00BCD4 26C6DA 4DD0E1 80DEEA B2EBF2 E0F7FA",
             "00B8D4 00E5FF 18FFFF 84FFFF"),
    "teal": ("004D40 00695C 00796B 00897B 009688 26A69A 4DB6AC 80CBC4 B2DFDB E0F2F1",
             "00BFA5 1DE9B6 64FFDA A7FFEB"),
    "green": ("1B5E20 2E7D32 388E3C 43A047 4CAF50 66BB6A 81C784 A5D6A7 C8E6C9 E8F5E9",
              "00C853 00E676 69F0AE B9F6CA"),
    "light-green": ("33691E 558B2F 689F38 7CB342 8BC34A 9CCC65 AED581 C5E1A5 DCEDC8 F1F8E9",
                    "64DD17 76FF03 B2FF59 CCFF90"),
    "lime": ("827717 9E9D24 AFB42B C0CA33 CDDC39 D4E157 DCE775 E6EE9C F0F4C3 F9FBE7",
             "AEEA00 C6FF00 EEFF41 F4FF81"),
    "yellow": ("F57F17 F9A825 FBC02D FDD835 FFEB3B FFEE58 FFF176 FFF59D FFF9C4 FFFDE7",
               "FFD600 FFEA00 FFFF00 FFFF8D"),
    "amber": ("FF6F00 FF8F00 FFA000 FFB300 FFC107 FFCA28 FFD54F FFE082 FFECB3 FFF8E1",
              "FFAB00 FFC400 FFD740 FFE57F"),
    "orange": ("E65100 EF6C00 F57C00 FB8C00 FF9800 FFA726 FFB74D FFCC80 FFE0B2 FFF3E0",
               "FF6D00 FF9100 FFAB40 FFD180"),
    "deep-orange": ("BF360C D84315 E64A19 F4511E FF5722 FF7043 FF8A65 FFAB91 FFCCBC FBE9E7",
                    "DD2C00 FF3D00 FF6E40 FF9E80"),
    "brown": ("3E2723 4E342E 5D4037 6D4C41 795548 8D6E63 A1887F BCAAA4 D7CCC8 EFEBE9",
              ""),
    "grey": ("212121 424242 616161 757575 9E9E9E BDBDBD E0E0E0 EEEEEE F5F5F5 FAFAFA",
             ""),
    "blue-grey": ("263238 37474F 455A64 546E7A 607D8B 78909C 90A4AE B0BEC5 CFD8DC ECEFF1",
                  ""),
}

# Blue grey with these OKLCH hues.
TINTED_GREYS = {
    "red-grey": 20.0,
    "green-grey": 150.0,
}

# The rows of the palette, top to bottom.
RAMP_ORDER = [
    "red", "pink", "purple", "deep-purple", "indigo", "blue", "light-blue", "cyan",
    "teal", "green", "light-green", "lime", "yellow", "amber",
    "orange", "deep-orange", "brown", "grey", "red-grey", "green-grey", "blue-grey",
]


# ---------- color math (OKLab, Björn Ottosson) ----------

def srgb_to_linear(c):
    return c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4


def linear_to_srgb(c):
    return 12.92 * c if c <= 0.0031308 else 1.055 * c ** (1 / 2.4) - 0.055


def hex_to_rgb(h):
    h = h.lstrip("#")
    return tuple(int(h[i:i + 2], 16) / 255 for i in (0, 2, 4))


def rgb_to_hex(rgb):
    return "#" + "".join(f"{round(min(1.0, max(0.0, v)) * 255):02X}" for v in rgb)


def hex_to_oklch(h):
    r, g, b = (srgb_to_linear(c) for c in hex_to_rgb(h))
    l_ = (0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b) ** (1 / 3)
    m_ = (0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b) ** (1 / 3)
    s_ = (0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b) ** (1 / 3)
    L = 0.2104542553 * l_ + 0.7936177850 * m_ - 0.0040720468 * s_
    a = 1.9779984951 * l_ - 2.4285922050 * m_ + 0.4505937099 * s_
    bb = 0.0259040371 * l_ + 0.7827717662 * m_ - 0.8086757660 * s_
    return L, math.hypot(a, bb), math.degrees(math.atan2(bb, a)) % 360


def oklch_to_hex(L, C, H):
    a, b = C * math.cos(math.radians(H)), C * math.sin(math.radians(H))
    l_ = (L + 0.3963377774 * a + 0.2158037573 * b) ** 3
    m_ = (L - 0.1055613458 * a - 0.0638541728 * b) ** 3
    s_ = (L - 0.0894841775 * a - 1.2914855480 * b) ** 3
    rgb = (4.0767416621 * l_ - 3.3077115913 * m_ + 0.2309699292 * s_,
           -1.2684380046 * l_ + 2.6097574011 * m_ - 0.3413193965 * s_,
           -0.0041960863 * l_ - 0.7034186147 * m_ + 1.7076147010 * s_)
    return rgb_to_hex([linear_to_srgb(max(0.0, v)) for v in rgb])


# ---------- palette ----------

def color(h):
    L, C, H = hex_to_oklch(h)
    return {"hex": h, "L": round(L, 4), "C": round(C, 4), "H": round(H, 2)}


def build():
    ramps = {}
    for name, (shades, accents) in MATERIAL_SWATCHES.items():
        ramps[name] = {"name": name, "steps": [color("#" + h) for h in shades.split()],
                       "accents": [color("#" + h) for h in accents.split()]}
    for name, hue in TINTED_GREYS.items():
        steps = []
        for h in MATERIAL_SWATCHES["blue-grey"][0].split():
            L, C, _ = hex_to_oklch("#" + h)
            steps.append(color(oklch_to_hex(L, C, hue)))
        ramps[name] = {"name": name, "steps": steps, "accents": []}
    missing = set(ramps) - set(RAMP_ORDER)
    if missing:
        raise SystemExit(f"RAMP_ORDER is missing: {', '.join(sorted(missing))}")
    return [ramps[n] for n in RAMP_ORDER]


# ---------- output ----------

def write_png(path, pixels, w, h, scale=1):
    raw = bytearray()
    for y in range(h * scale):
        raw.append(0)
        for x in range(w * scale):
            raw.extend(pixels[y // scale][x // scale])

    def chunk(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xFFFFFFFF)

    data = b"\x89PNG\r\n\x1a\n"
    data += chunk(b"IHDR", struct.pack(">IIBBBBB", w * scale, h * scale, 8, 2, 0, 0, 0))
    data += chunk(b"IDAT", zlib.compress(bytes(raw), 9)) + chunk(b"IEND", b"")
    Path(path).write_bytes(data)


def rgb255(h):
    return tuple(round(v * 255) for v in hex_to_rgb(h))


def emit(ramps, out, name):
    out.mkdir(parents=True, exist_ok=True)

    hexes = [c["hex"] for r in ramps for c in r["steps"] + r["accents"]]
    (out / f"{name}.hex").write_text("\n".join(h[1:].lower() for h in hexes) + "\n")

    n_steps = max(len(r["steps"]) for r in ramps)
    n_acc = max(len(r["accents"]) for r in ramps)
    gpl = ["GIMP Palette", f"Name: {name} (Material Design swatches)", f"Columns: {n_steps + n_acc}", "#"]
    for r in ramps:
        for label, group in (("", r["steps"]), ("A", r["accents"])):
            for i, c in enumerate(group):
                rr, gg, bb = rgb255(c["hex"])
                gpl.append(f"{rr:3d} {gg:3d} {bb:3d}\t{r['name']}-{label}{i + 1}")
    (out / f"{name}.gpl").write_text("\n".join(gpl) + "\n")

    (out / f"{name}.json").write_text(json.dumps(ramps, indent=2) + "\n")

    bg = (24, 24, 24)
    w = n_steps + 1 + n_acc
    grid = []
    for r in ramps:
        row = [rgb255(c["hex"]) for c in r["steps"]] + [bg] + [rgb255(c["hex"]) for c in r["accents"]]
        grid.append(row + [bg] * (w - len(row)))
    write_png(out / f"{name}.png", grid, w, len(grid))
    write_png(out / f"{name}-preview.png", grid, w, len(grid), scale=24)


def main():
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("--out", type=Path, default=Path(__file__).resolve().parent)
    p.add_argument("--name", default="material")
    args = p.parse_args()

    ramps = build()
    emit(ramps, args.out, args.name)
    for r in ramps:
        print(f"{r['name']:<12} {' '.join(c['hex'] for c in r['steps'])}  | {' '.join(c['hex'] for c in r['accents'])}")


if __name__ == "__main__":
    main()
