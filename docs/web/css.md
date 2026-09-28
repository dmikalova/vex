# Client CSS conventions (`web/app.css`)

- BEM-ish: a block (`.card`) and `--modifier` classes (`.card--dimmed`,
  `.log-group--p0`). No inline styles from Go.
- House colours are custom properties (`--nm`, `--tp`, `--edge`) supplied by
  the `.card-<house>` class from `palette.go`; markup carries only class names.
- The one exception is a **measurement**: where a card sits on screen is a
  runtime fact no class can name, so `view_focus.go` hands the lifted card's
  rect over as custom properties (`--focus-x`, `--focus-w`, …) and `app.css`
  owns everything done with them. Values, never styling: a rule needing a new
  declaration in Go belongs in the stylesheet.
- Keep every animation's `-a`/`-b` pair in sync; their keyframes must be
  identical.
- **Enlarge a card by resizing its box, not by `transform: scale()`**, which
  does not reflow the text the board was clipping. Resize the box and the inner
  font sizes together (they are hardcoded px), as `.card-preview` and
  `.card-focus` do. A card whose height follows its content also undoes the
  `flex: 1 1 0%; min-height: 0` on `.card-body`/`.card-rules`, which exist to
  fill and clip a fixed slot.
- **The card frame is a flat two-tone gradient, not a clip or a mask.** `.card`
  paints a hard-stop `linear-gradient`: the name-banner colour (`--nm`) upper
  left, the type band (`--tp`) lower right, split by a straight diagonal. It
  scales with the card and respects `border-radius`, so the card needs no
  `overflow: hidden` or SVG mask, which lets the left-edge house and bonus
  icons (`.card-bonuses`) hang off the card. The keybar and `.card-kind` round
  their own corners. The `125deg` angle and `48%` stop are the two knobs.
- **A card that overhangs its row has to leave the board's coordinate space**
  (`position: fixed`): `.card-strip` is `overflow-x: auto`, which forces
  `overflow-y: hidden`, and `.board-area` is `overflow: hidden`. The same trap
  catches tooltips in `.score-pill` (`overflow-x: auto`), so every `data-tip`
  shares **one** floating label (`#tip-float`, `installTips`) placed with
  `position: fixed` from a measurement of the element under the pointer, not a
  bubble per icon.
- **Place a content-sized overlay from the edge it is nearest**, not its
  centre: its height is unknown until laid out, and a second measuring pass may
  never come (a frozen tab, a dropped frame). `.card-focus` anchors `top` for a
  card in the opponent's half and `bottom` for one in the player's.
- **An overlay placed from a measurement tracks everything that moves it.**
  `.card-focus` re-measures from go-app's `Resizer` (`OnResize`) and from a
  document-level `scroll` listener in the **capture** phase; scroll does not
  bubble, and every card strip scrolls on its own.
- **A copy of a card that covers its neighbours takes the pointer.**
  `pointer-events: none` on `.card-focus` let a drag land on whatever card lay
  underneath. The face is `pointer-events: auto` and its own drag source, and
  the wheel is handed to the strip underneath in Go (`wheelOverFocus`).
- **Keep `filter` off the ancestor of a drag source.** A filtered ancestor
  stops the dragged element being its own layer, and the drag image then pulls
  in a sliver of its neighbour. `.card-focus` gives its drop-shadow to the face
  and the verbs separately.
