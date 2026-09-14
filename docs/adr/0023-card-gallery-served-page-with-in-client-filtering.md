# The Card gallery is an always-served /cards page filtered in the client

## Context

We want a browsable **Card gallery** at `/cards` that shows every card in the
database and filters by House, Set, Type, and by name and rules text, with both
union and disjoint searches. The registry already enumerates every card
(`cards.All()`), and the client already renders a printed-only face
(`printedFace`) reused across piles, previews, and the Style gallery. Two
existing surfaces are nearby but wrong for this: the Style gallery at `/style` is
a development-only surface served solely by `mage web`, and the Card picker is a
manual-mode control, not a browsing page.

The gallery is large — hundreds to thousands of cards, each a full card face plus
the generated Icon strip — so how the page bounds its rendering cost is a real
decision, and how filtering is expressed is a small language of its own.

## Decision

`/cards` is a first-class route registered beside `/rulebook` and `/glossary`,
always served (not gated behind `mage web` like `/style`). It reuses
`printedFace` and the Icon strip, so the gallery doubles as a living showcase of
the iconography.

Filtering runs entirely in the client over `cards.All()` — no new data store and
no server round-trip. Facet semantics are **OR within a category, AND across
categories**: choosing House=Mars, House=Sanctum, Type=Creature matches
`(Mars OR Sanctum) AND Creature`. Name and rules-text search takes a small query
syntax parsed by a hand-rolled, quote-aware tokenizer (`term term` = all-of,
`a|b` = either, `-x` = exclude, `"phrase"` = exact run, `\"` = a literal quote).
The syntax is tiny and the WASM binary size matters, so we do **not** pull in a
full-text search library.

The page bounds rendering cost by **windowing**: only the cards on or near the
viewport are mounted as full printed faces; every other matching card renders as a
fixed-size text placeholder that keeps its name and rules text (`window` and
`galleryPlaceholder` in `web/gallery.go`). A document scroll listener and a
frame-measure callback slide the window as the viewport moves (`installWindowing`,
`recomputeWindow`), and the placeholders preserve the full scroll height and the
find-in-page text of the whole result set — so browser find, anchors, and
accessibility keep working across every match, not just the drawn window. Keeping
the placeholder text mounted is what lets windowing coexist with Ctrl-F, which is
why the earlier CSS-only plan (`content-visibility`) was replaced.

## Consequences

- The gallery is generated, not curated: it reflects exactly what the registry
  holds, so a newly implemented card appears without touching the page.
- Filtering logic and the query tokenizer live in the web/gallery package and are
  unit-testable without a running game.
- Windowing bounds paint and layout to a screenful of faces regardless of catalog
  size, while the text placeholders keep every match findable and scrollable.

## Update — windowing replaced the CSS-only ceiling

The original decision bounded cost with CSS `content-visibility: auto` plus
`contain-intrinsic-size` and deliberately rejected JavaScript
windowing/virtualization (it breaks Ctrl-F and needs scroll math), keeping every
match mounted as a full face and deferring windowing as an escape hatch. The
implementation instead adopted windowing directly: full faces are mounted only near
the viewport and the rest become text placeholders. The two objections were
resolved — the placeholders keep the searchable text (so Ctrl-F still spans the
whole result set) and a scroll listener plus frame measure carry the scroll math —
so windowing became the primary mechanism rather than the deferred fallback.
