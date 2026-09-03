# The Style gallery renders real components, gated at runtime

## Context

The client's visual vocabulary has no single place it can be seen. It is spread
across one 2000-line stylesheet, 38 SVG assets, eight house palettes of five
tokens each, and a dozen paired keyframe animations — and the only way to look at
any of it is to reach a game state that happens to show it. Some of it is
effectively unreachable: `cardFlight` needs a card to leave play, the Æmber
capture segment needs a capture, and a Sanctum Upgrade needs that card to be
drawn. Choosing a font, or checking that a new icon sits right next to the old
ones, means playing until the board shows you.

Every project eventually writes a page like this, and every one of them rots. The
two ways it happens are worth naming, because they drive the decision: the page
duplicates the app's markup and slowly stops matching it, or the page hardcodes a
list of examples and silently omits everything added afterwards.

## Decision

Add a **Style gallery** at `/style` — one scrolling page showing the tokens,
icons, fonts, card faces, a Player bar, and the animations together.

**It renders the real components.** Specimens go through `cardView`, `scorePill`
and friends, over a real `engine.Game` filled through the engine's own setup
calls. Those are methods on the unexported `game` component, so the gallery is a
file inside `internal/web` that builds a synthetic `game` — and where the two
disagree, the gallery changes. The app is never reshaped to make the gallery
easier.

**Its specimens are queries, not lists.** The card faces are derived from
`cards.All()` by walking House × Card type and by feature predicate, so the page
cannot omit something that exists. Combinations no card covers are drawn as gaps,
because absence is information: it is how you see that a set has no Dis Upgrade.

**It is gated by an environment switch, not a build tag.** `mage web` sets
`VACTROL_STYLE=1`; nothing else does, so `/style` exists on a development server
and 404s everywhere else. The switch has to be read on both sides of the build
and agree — go-app 404s a path the _server_ has no route for, and renders nothing
for a path the _client_ has no route for — which `app.Getenv` does, reading the
process environment on the server and the `Env` map it passes down on the client.

**Its tests compare lists, never markup.** Each asserts that two enumerations
agree — every asset appears, every house has a palette and a font token, every
card type, rarity and keyword has a specimen or an explicit gap marker. No test
pins layout, class names, or text, so restyling never breaks them.

## Consequences

- **Gallery code ships in the production wasm bundle.** This looks like a mistake
  and is not: go-app routes on the client, so a build tag is the only thing that
  would truly remove it, and tagged code is not compiled by default — it would rot
  exactly as the `//go:build todo` card stubs do. The marginal size is near zero,
  since every component and the whole card database are already bundled for the
  game and the Card picker. Do not "fix" this by adding a tag.
- **A route must be registered on the server as well as the client.** The obvious
  gate — checking `location.hostname` in the client — cannot work on its own,
  because go-app's handler 404s any path it has no route for and the page never
  reaches the client at all. Anything else that wants to be conditionally routed
  has the same constraint.
- The gallery is coupled to the shape of the game component, and a refactor there
  will break it. That is the intended direction of the dependency: it breaking is
  the signal that the page still reflects reality.
- A canonical `Keyword` enumeration has to exist for the keyword coverage test,
  which is what motivates making `Keyword` an enum rather than a string.
- The synthetic position is built by the engine's own setup calls but its display
  values (chains, keys, Æmber) are then set outright, because no single reachable
  game state exercises every visual branch of a Player bar at once.
