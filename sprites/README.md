# Sprites

This directory holds the source files for the Vex pixel art. The source format
is `.aseprite`. Each sprite has one source file. The game does not load these
sprites yet; the UI uses the SVG icons in [web/assets](../web/assets).

## Status

- [ ] Set up the tablet ([Tablet setup](#tablet-setup)).
- [ ] Set up Aseprite ([Aseprite setup](#aseprite-setup)).
- [x] Choose the palette and save it as `palette/material.gpl`
      ([Palette](#palette)).
- [ ] Do the learning exercises ([Learning path](#learning-path)).
- [ ] Draw tier 1 ([Checklist](#checklist)).
- [ ] Decide the art direction and record it in an ADR
      ([Open decisions](#open-decisions)).
- [ ] Add a card art slot and PNG loading to the UI
      ([Integration](#integration)).

## Tablet setup

1. Install the Wacom driver: `brew install --cask wacom-tablet`. The dotfiles
   Brewfile has this cask, so `chezmoi apply` also installs it.
2. Restart the Mac.
3. Open **System Settings > Privacy & Security**. Allow the Wacom driver
   entries in **Accessibility** and **Input Monitoring**. If the pen moves the
   cursor but does not click, this step is the usual cause.
4. Connect the tablet. The tablet is an Intuos S Bluetooth (CTL-4100WL). It
   has a micro-USB port and comes with a USB-A to micro-USB cable.
   - **USB (recommended first):** connect the cable to a USB-A port on the dock
     or hub, or use a USB-C adapter. USB also charges the tablet.
   - **Bluetooth:** disconnect the cable. Push the power button and hold it
     until the blue LED flashes. Open **System Settings > Bluetooth** and
     select **Connect** for the tablet. Do the pairing from Wacom Center if
     macOS does not show the tablet.
5. Open **Wacom Center** and select the tablet. Install the firmware update if
   Wacom Center shows one.
6. In **Mapping**, map the tablet to one display only. With more than one
   display, the tablet area is otherwise stretched across all displays.
7. In **Mapping**, turn on **Force Proportions**. This keeps a drawn circle
   round.
8. Set the pen buttons:

   | Control           | Action             | Why                                  |
   | ----------------- | ------------------ | ------------------------------------ |
   | Lower side button | Pan (middle click) | Aseprite pans with the middle button |
   | Upper side button | Right click        | Erase (see Aseprite setup step 3)    |
   | Tip double click  | Off (distance 0)   | Stops accidental double clicks       |

9. Set the 4 ExpressKeys on the top edge of the tablet to these Aseprite keys:

   | Key | Aseprite action                     |
   | --- | ----------------------------------- |
   | 1   | `Cmd+Z` undo                        |
   | 2   | `Cmd+Shift+Z` redo                  |
   | 3   | `Option` (hold) to pick a color     |
   | 4   | `X` swap the two colors             |

The pressure of the pen is not important for pixel art. A 1 px pencil
makes the same pixel at all pressures.

## Aseprite setup

1. Open **Edit > Preferences > General**. Set **Screen Scaling** and **UI
   Scaling** so that the text is easy to read on the Retina display.
2. Open **Edit > Preferences > Editor**. Turn on **Zoom with scroll wheel** if
   you want the mouse wheel to zoom.
3. In **Edit > Preferences > Editor**, set **Right-click** to **Erase**. Then
   the upper pen button erases.
4. Select the pencil (`B`). In the context bar, turn on **Pixel-perfect**.
   This removes the extra corner pixels from freehand lines.
5. In the context bar, set the pencil size to 1 and turn off **Dynamics**
   (pen pressure).
6. Turn on the pixel grid: **View > Show > Pixel Grid**. Aseprite shows it
   only at high zoom.
7. Use **View > Duplicate View** to see the sprite at two zoom levels at the
   same time. Keep one view at 100% or 200%. Judge the sprite at that size,
   because the game shows it at that size.

### New sprite settings

Use these settings in **File > New** for every sprite:

| Setting            | Value                    |
| ------------------ | ------------------------ |
| Color mode         | Indexed                  |
| Background         | Transparent              |
| Pixel aspect ratio | Square pixels (1:1)      |
| Size               | See the size table below |

Indexed color mode locks the sprite to the palette. Change the palette, and
every sprite that uses it changes too.

### Sizes and ratio

Pixel art must show at a whole-number scale: 1x, 2x, 3x. A scale of 2.5x makes
some pixels wider than others. The sizes below divide into the sizes that the
UI uses.

The card face is `9rem` by `12rem`, which is 144 x 192 CSS px. That is a 3:4
ratio.

| Asset             | Canvas  | Ratio | Shown at                                   |
| ----------------- | ------- | ----- | ------------------------------------------ |
| Status icon       | 16 x 16 | 1:1   | 16 px (1x) or 32 px (2x)                   |
| House logo        | 32 x 32 | 1:1   | 32 px (1x) or 64 px (2x)                   |
| House logo, small | 16 x 16 | 1:1   | 16 px, inline next to text                 |
| Character         | 48 x 48 | 1:1   | 96 px (2x); 144 px (3x) is the card width  |
| Card back         | 72 x 96 | 3:4   | 144 x 192 px (2x) is the full card face    |

Do not scale a sprite down to make a smaller version. Draw the 16 x 16 house
logo again by hand.

## Palette

Use one palette for all sprites. A shared palette makes the sprites look like
one set.

The palette is the Material Design colors. The script
`palette/genpalette.py` generates it. It has no dependencies. The palette has
the 19 Material ramps. Each ramp has the shades 900 to 50, from dark to light.
Each chromatic ramp also has the accents A700 to A100. The palette also adds a
red grey and a green grey. These are Material blue grey with a different hue,
and they keep the lightness and chroma of each blue grey shade.

The Material ramps do not have the same lightness at each shade. For example,
yellow 500 is much lighter than blue 500. If you swap one ramp for another,
check the value structure of the sprite again.

1. Run `python3 sprites/palette/genpalette.py`. The script writes these files
   to `palette/`:
   - `material.gpl`: the palette for Aseprite and GIMP.
   - `material.hex`: the palette for Lospec.
   - `material.json`: the OKLCH lightness, chroma, and hue of each color.
   - `material.png`: one pixel for each color.
   - `material-preview.png`: the swatches at a larger scale. The shades are on
     the left and the accents are on the right.
2. To change the order of the ramps, change `RAMP_ORDER` in the script. Then
   run the script again.
3. In Aseprite, open the palette options menu and select **Load Palette**.
   Select `palette/material.gpl`.
4. Give each house a ramp of 3 to 5 colors from the palette. Start near the
   color that the SVG house icon uses now:

   | House         | SVG color now        |
   | ------------- | -------------------- |
   | Brobnar       | `#e08a3c`            |
   | Dis           | `#cc4657`            |
   | Logos         | `#3fc1e6`            |
   | Mars          | `#6d3f96`            |
   | Sanctum       | `#1c6adb`            |
   | Saurian       | `#c9a24c`, `#1f9e94` |
   | Shadows       | `#aeb4bb`            |
   | Star Alliance | `#1c2f5c`, `#d9a41f` |
   | Untamed       | `#f2ab2b`            |

5. Keep the generated files next to the script.

The UI background is dark (`#1c1c1b`). Test each sprite on that color, not
only on the checkerboard.

## File layout

```text
sprites/
  palette/genpalette.py           # generates the material.* palette files
  palette/material.gpl
  icons/<name>.aseprite           # name = the SVG name in web/assets
  houses/<house>.aseprite         # 32 x 32
  houses/<house>-16.aseprite      # 16 x 16
  characters/<house>/<name>.aseprite
  card-back.aseprite
```

Use lowercase names with hyphens. Use the house names from the SVG files
(`staralliance`, not `star-alliance`).

## Learning path

Watch
[Everything You Need to Know About ASEPRITE](https://www.youtube.com/watch?v=Ep-_nJNHPc4)
first. Then do the exercises in this order. Each exercise makes a sprite from
tier 1, so the practice is also real work.

| #   | Exercise                    | Skill                                   | Sprite         |
| --- | --------------------------- | --------------------------------------- | -------------- |
| 1   | Clean lines                 | Pixel-perfect lines, even slopes        | `damage` (16)  |
| 2   | Three-color ramp            | Ramps, light from one direction         | `aember` (16)  |
| 3   | Hue shift                   | Shadows cooler, highlights warmer       | `key` (16)     |
| 4   | Silhouette first            | Readable shape at 1x                    | `power` (16)   |
| 5   | Big shape, then small shape | Detail at 32 px, then simplify to 16 px | One house logo |
| 6   | Full character              | Silhouette, flat colors, shade, outline | One character  |

### The rules to learn

- **Lines.** Use even steps: 1:1, 1:2, 2:1, 1:3. A line with steps of 1, 2, 1,
  3 looks jagged. A step that touches another step only at a corner is a
  "jaggy"; remove it.
- **Silhouette.** Fill the shape with one dark color first. If you cannot tell
  what it is at 1x, the details will not help.
- **Ramps.** A ramp is 3 to 5 colors from dark to light. Use one ramp for each
  material.
- **Hue shift.** Move the hue as the value changes: shadows toward blue or
  purple, highlights toward yellow. A ramp with one hue looks flat.
- **Value.** Check contrast in grayscale. Make a grayscale copy of the palette,
  or look at the sprite with the saturation at zero.
- **One light source.** Put the light at the top left for every sprite. Do not
  shade from the outline inward ("pillow shading").
- **Outlines.** Use a dark outline at 16 px. At 48 px, try a selective outline
  (sel-out): color the outline with a darker tone of the color next to it.
- **Anti-aliasing and dithering.** Use them only a little, and not at 16 px.
- **Clusters.** Avoid single pixels of a color that stand alone. Small groups
  of pixels read better and look cleaner.

### Courses and references

- [Pixel Art Class](https://www.youtube.com/playlist?list=PLLdxW--S_0h4dlWUpl-TzBp-ulqK3NiM_)
  by AdamCYounis. A long video course that uses Aseprite. Start with the
  basics videos.
- [Getting Started with Pixel Art!](https://www.youtube.com/playlist?list=PLxfQIomHccxvoTON6hXhfZyAUdFXd-z1P)
  by Brandon James Greer. Short videos about lines, color, and sizes.
- [Pixelblog catalogue](https://www.slynyrd.com/pixelblog-catalogue) by
  Slynyrd. Written lessons. Read the lessons about palettes and icons.
- [Pixel art tutorials](https://saint11.art/blog/pixel-art-tutorials/) by
  Pedro Medeiros (saint11). Short image tutorials, one topic each.
- [Lospec tutorials](https://lospec.com/pixel-art-tutorials). An index of
  tutorials by topic.
- [Aseprite docs](https://www.aseprite.org/docs/).

## Checklist

Mark a sprite done when the `.aseprite` file is in this directory and the
sprite looks correct at 1x and 2x on the dark background.

### Tier 1: core icons (16 x 16)

- [ ] `aember` — Æmber
- [ ] `damage` — Damage
- [ ] `power` — Power
- [ ] `shield` — Armor
- [ ] `key` — Key
- [ ] `key-red` — Red key
- [ ] `key-blue` — Blue key
- [ ] `key-yellow` — Yellow key
- [ ] `forge` — Forge
- [ ] `capture` — Capture
- [ ] `chains` — Chains
- [ ] `stun` — Stun
- [ ] `enrage` — Enrage
- [ ] `ward` — Ward
- [ ] `exhausted` — Exhausted
- [ ] `draw` — Draw

### Tier 1: house logos

| House         | 32 x 32 | 16 x 16 |
| ------------- | ------- | ------- |
| Brobnar       | [ ]     | [ ]     |
| Dis           | [ ]     | [ ]     |
| Logos         | [ ]     | [ ]     |
| Mars          | [ ]     | [ ]     |
| Sanctum       | [ ]     | [ ]     |
| Saurian       | [ ]     | [ ]     |
| Shadows       | [ ]     | [ ]     |
| Star Alliance | [ ]     | [ ]     |
| Untamed       | [ ]     | [ ]     |

### Tier 2: characters (48 x 48, 3 for each house)

The ideas are starting points only. Make original designs.

- [ ] Brobnar — ideas: a giant, a goblin, a brawler
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Dis — ideas: a demon, an imp, a tormentor
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Logos — ideas: a scientist, a robot, a lab creature
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Mars — ideas: a Martian soldier, a saucer pilot, an alien beast
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Sanctum — ideas: a knight, an angel, a priest
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Saurian — ideas: a dinosaur senator, a legionnaire, a gladiator
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Shadows — ideas: a thief, an elf rogue, a faerie
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Star Alliance — ideas: a starship captain, an engineer, an alien crew
      member
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3
- [ ] Untamed — ideas: a beast, a druid, a forest spirit
  - [ ] Character 1
  - [ ] Character 2
  - [ ] Character 3

### Tier 2: card faces

- [ ] `card-back` — Facedown card (72 x 96)
- [ ] `type-creature` — Creature (16)
- [ ] `type-tactic` — Action (16)
- [ ] `type-artifact` — Artifact (16)
- [ ] `type-upgrade` — Upgrade (16)

### Tier 3: optional

Draw these only if the pixel icons and the SVG icons look wrong together. Each
group has one SVG file for each item in [web/assets](../web/assets).

- [ ] Keywords (`kw-*`, 12 icons)
- [ ] Generic counters (`generic-counter-*`, 7 icons)
- [ ] Sigils (`sigil-*`, 3 icons)
- [ ] Rarity (`rarity-*`, 5 icons)
- [ ] Zones (`zone-*`, 7 icons)
- [ ] Ability glyphs (`glyph-*`, 23 icons)

UI controls (undo, redo, restart, deck list, manual mode) stay SVG.

## Open decisions

- **Art direction.** Step 5 of the [roadmap](../docs/roadmap.md) plans
  procedural art from layered SVG parts. Pixel art replaces or changes that
  plan. Record the decision in an ADR before the UI loads the sprites.
- **House theme.** The roadmap plans original houses (Alchemia, Mechanica,
  Geometria). The engine uses the 9 KeyForge houses now. Characters for the
  KeyForge houses can become unused if the theme changes.
- **Intellectual property.** The KeyForge house names, logos, and characters
  belong to their owner. Draw original designs. Do not trace official art.

## Integration

The UI does not show PNG sprites yet. An integration task must:

- Add an art slot to the card face.
- Export the sprites as PNG at 1x with the Aseprite CLI (see below).
- Show each sprite with the CSS `image-rendering: pixelated`.
- Give each sprite a fixed size in `px` that is a whole multiple of the canvas.
  The SVG icons use `em` sizes (for example `1.1em`). An `em` size gives a
  scale that is not a whole number, which blurs pixel art.

The export command:

```sh
/Applications/Aseprite.app/Contents/MacOS/aseprite -b in.aseprite --save-as out.png
```
