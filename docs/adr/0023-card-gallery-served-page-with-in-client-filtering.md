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

The page keeps every matching card mounted in the DOM — so browser find,
anchors, and accessibility all work — and bounds rendering cost with CSS
`content-visibility: auto` plus `contain-intrinsic-size`, which skips layout and
paint for off-screen cards without JavaScript virtualization. We deliberately
reject windowing/virtualization first (it breaks Ctrl-F and needs scroll math);
if the CSS ceiling proves insufficient at extreme sizes we layer infinite scroll
on top later.

## Consequences

- The gallery is generated, not curated: it reflects exactly what the registry
  holds, so a newly implemented card appears without touching the page.
- Filtering logic and the query tokenizer live in the web/gallery package and are
  unit-testable without a running game.
- Bounding cost with `content-visibility` keeps the full result set findable and
  scrollable; the fallback to infinite scroll is a known, deferred escape hatch,
  not a commitment.
