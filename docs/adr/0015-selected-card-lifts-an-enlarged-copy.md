# Clicking a card lifts an enlarged copy with its actions on it

## Context

A card on the board is small. Its slot is a fixed `9rem`/`12rem` box, and
`.card-body`/`.card-rules` are `flex: 1 1 0%; min-height: 0` so a long rules block
is clipped to fit. That is right for the board — a battleline of full-height cards
would not fit — but it means the one card the player is deciding about is the one
whose text they cannot fully read.

The verbs for that decision (reap, fight, use, play) needed a home too. Putting
them in the sidebar makes the player look away from the card to act on it, and
puts the answer to "what can this do" a glance away from the thing it is about.
It is worse on a small screen: the sidebar is where the board is not, so on mobile
the buttons and the card they act on may not even be on screen together, and the
tap travel between them is the whole width of the layout. Anchoring the verbs on
the enlarged card puts the action within a thumb's reach of the thing being acted
upon — better on mobile specifically, and a shorter eye-and-hand path to the card
for everyone.

The obvious way to enlarge a card is `transform: scale()` on the board card. It is
also wrong twice over: a transform does not reflow text, so the clipped rules stay
clipped, only blurrier; and scaling the real card in place shoves its neighbours
around, so the board reflows every time a selection changes.

## Decision

**Clicking a card selects it, and the selection lifts an enlarged copy** —
`cardFocus` in [view_focus.go](../../internal/web/view_focus.go) — onto a layer
above the board, with a button for each thing the card can do under it and the
note saying why it can do nothing.

**The real card never moves.** A copy is laid out over the slot it was lifted from
and grown from there. The board never reflows, the card stays put while its
buttons are up, and a click that lands on the copy is swallowed rather than
falling through to the board it covers, so blowing a card up cannot cost a
misclick on a card the player can no longer see.

**The copy is grown by resizing its box, not scaling a picture.** Only its width
is set (`--focus-w`); the height follows its content — a stack of upgrades comes
out taller — floored at `focusMinGrow` and ceilinged at what the window has left.
Resizing means setting three things together, the way `.card-preview` and
`.card-focus` both do: the box, the nominal `--card-full-h` (derived from the
_width_, not the rendered height, so the ogee mask is drawn for the right box),
and the inner font sizes. A content-sized card also has to undo the
fill-and-clip `flex`/`min-height` its slot form relies on.

**Where the card sits is a measurement, and the only thing the markup carries.**
`measureFocus` reads the selected card's `getBoundingClientRect` and the window
size and hands them to `app.css` as custom properties (`--focus-x`, `--focus-w`,
…). The stylesheet owns everything done with them — values, never styling, the
same contract house colours arrive under. A rule needing a new declaration in Go
is a rule that belongs in the stylesheet.

**It is driven by the selection, gated by phase.** `focusCardID` lifts only in the
main and flank phases and only when nothing modal is up (busy, a prompt, the
picker, forging). Picking a flank keeps the lift, because the question is about the
card being placed and is asked on it; picking a fight target drops it, because the
answer is another card on the board that a copy blown up over it would cover.

## Consequences

- **The measure/render pair settles in two frames.** The copy is on its own layer
  and reflows nothing, so the render after a measurement measures the same numbers
  and stops. `measureFocus` is also called from the selection itself, so the common
  case places the copy on its first render rather than a frame later.
- **The overlay is placed from the edge it is nearest, and tracks what moves it.**
  Its height is unknown until laid out, and a frozen tab will not give the second
  render pass a self-centring overlay would need — so `.card-focus` anchors `top`
  in the opponent's half and `bottom` in the player's, growing away from the near
  edge. It re-measures from go-app's `OnResize` and a document-level `scroll`
  listener in the **capture** phase, since scroll does not bubble and every card
  strip scrolls on its own.
- **A card that overhangs its row leaves the board's coordinate space**
  (`position: fixed`): `.card-strip` and `.board-area` both clip their overflow,
  so a copy grown past its slot would be cut off inside either.
- **The copy has to take the pointer, carefully.** `pointer-events: none` would let
  a drag fall through to whatever card the enlarged face lies over, so the face is
  `pointer-events: auto` and its own drag source, and the wheel is handed to the
  strip underneath in Go (`wheelOverFocus`). `filter` is kept off the drag source's
  ancestor — a filtered ancestor makes the browser cut the drag image out of it and
  drag in a sliver of the neighbour — so the drop-shadow goes on the face and the
  verbs separately.
- **The grow animation replays with the `-a`/`-b` parity trick.** go-app patches
  the same element and a CSS animation only restarts when its name changes, so
  `focusParity` flips whenever the lift moves to another card.
- This lives entirely in the DOM-bound layer of `internal/web`, whose coverage is
  deliberately ungated: the placement is a runtime measurement no unit test can
  stand in for, and the tests here assert the geometry (`focusBox`, edge anchor)
  rather than pinning markup.
