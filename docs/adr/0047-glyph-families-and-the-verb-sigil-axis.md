# Glyphs split into mechanic-domain families, and the verb carries a sigil

## Context

ADR 0022 established the Iconography pass: a Visitor in `internal/web` that
type-switches over the engine's `Effect` AST and emits the composed glyphs of a
card's icon strip. Three years of mechanics later that Visitor is one function,
`effectGlyphs` in `internal/web/icon.go`, with 155 case labels across 633 lines
inside a 1518-line file. Its cognitive complexity is 90 by the gate's measure and
99 by a current `gocognit`, and it is the **last remaining exclusion** in the
`tools.golangci` exclusions list in `mklv.config.json`.

Two separate problems hide behind that one number.

**Structure.** `gocognit` does not charge per case; a bare dispatch switch is
complexity 1. The 99 comes almost entirely from the ~40 `if` branches inside case
bodies, so hoisting those bodies alone would delete the exclusion without
splitting anything. The gate is therefore not the reason to act — readability is.
A new effect node's glyph has no obvious home today, and nothing tells the agent
who added `effect_aember.go` where to draw its icon.

## Vocabulary.** The pass has a rich **subject** vocabulary and no **verb

vocabulary. `decor` is a five-bit flag set — enemy, friendly, each, chosen, this —
and every bit describes *who* the glyph is about, never *what happens to them*.
The verb is left to the base asset, so verbs collapse onto each other. Measured
across every ability on every card, twelve pairs of genuinely different leaf
mechanics render byte-identical glyph strips:

| pair | both draw |
| --- | --- |
| `CaptureAember` / `GainAember` | æmber, qty, enemy tint |
| `CaptureAember` / `DistributeCapture` | æmber, enemy tint |
| `Exalt` / `MoveAemberFromPool` | æmber qty → target |
| `CannotBeDealtDamage` / `Ward` | shield → target |
| `PlaceCounter` / `RemoveCounters` | counter → target |
| `LookAtTopOfDeck` / `RevealTopOfDeck` | look → deck |
| `PutCard` / `ShuffleIntoDeck` | card → deck |
| `PutChosen` / `PutFromPlay` | creature → hand |
| `PutUnderIntoPlay` / `TriggerGraftedPlayEffect` | card → play |
| `CreaturesCannot` / `Restrict` | verb, ban |
| `MayPlayOrUse` / `PlayOrUse` | play/use pair |

A player cannot tell "place a counter" from "remove counters", and `Ward` draws
`shield.svg` — armour's icon — while a dedicated `ward.svg` sits unused in
`web/assets`. `LoseAember` has no modifier channel at all: it fakes one by
writing `"−"` into the quantity slot, the only effect in the pass that does.

## Decision

### Families are mechanic domains, mirroring `effect_<mechanic>.go`

`effectGlyphs` splits along **effect category**, a coarsened mirror of the axis
`internal/engine` already splits on. Not glyph shape — `DealDamage` draws a damage
glyph *and* a creature noun, so most cases belong to two shape groups at once. Not
render target — there is one strip, and `icons.go` / `view_icons.go` are the asset
table and the renderer, neither of which varies by effect. A 1:1 mirror of the
engine's ~90 `effect_*.go` files would give ~90 three-case functions, so the
families are coarsened to nine of 13–25 cases each: æmber and keys, damage and
removal, creature state, stats and abilities, zones, deck and reveal, play and
board, turn and house, and composition (the wrappers that recurse).

The payoff is that a new effect node's glyph has one obvious home.

### Dispatch is a family chain, policed by a disjointness test

`effectGlyphs` becomes a loop over a package-level slice of family functions; the
first family that claims the effect wins. Each family returns
`(glyphs []glyph, covered, ok bool)`. The third value is required because
`covered = false` already means something else: `Then`, `PutFromPlay`, `PutChosen`
and `PutCard` return real glyphs with `covered = false` today, so a family cannot
signal "not mine" by returning `false`.

The cost is that the compiler cannot see across functions, so two families listing
the same effect type would silently make the second dead code — a wrong glyph, not
an unknown one, which no existing test catches. A **disjointness test** walks every
effect on every card and fails unless exactly one family claims it. That covers the
same population `TestIconTotality` already polices.

The rejected alternative was a flat delegating switch: keep one 155-case switch
whose every case body is a one-line call into a family file. It gets
compiler-checked disjointness for free, but it is not a split — the dispatch stays
~400 lines and spawns ~40 named one-call helpers.

### The verb takes a sigil; the subject keeps the glow

A verb modifier is a small mark composited into the corner of its base glyph: `+`
on the æmber for gain, `−` for lose, `capture.svg` for capture, a reaching hand for
steal, and so on. It is a closed enum on `glyph` alongside `decor`, with its own
renderer, exactly as `decor` has `decorClass`.

The verb does **not** take the glow, even though a green/red/orange glow is the
most direct reading of "gain, lose, exalt". The glow channel is already spent:
`card-glyph--enemy` is a red drop-shadow and `card-glyph--friendly` a blue one, so
a red glow for "lose" would stack red on red for "the opponent loses æmber", and
"you gain" would be green while "the opponent gains" stayed red. The two axes would
fight for the same pixels, and moving enemy/friendly off the glow would rewrite a
subject vocabulary that works and is used by every target glyph in the pass.

A sigil is also **shape, not hue**, which follows the precedent `rarityMark`
already set: it carries tier by distinct polygon silhouettes, and its comment says
that is "not colour, … so it stays legible to colour-blind players". A sigil may
still be tinted, so green and red survive as a second cue rather than the only one.

A verb glyph of its own in the strip was rejected: it spends a horizontal slot on a
strip already squeezed to fit by `cardFitScript`, and it is not composition.

### No collision is exempt, but synonyms are declared

A collision test asserts that two distinct leaf mechanics never render the same
strip. Four measured pairs are benign — the same mechanic under two node names
(`Destroy`/`DestroyChosen`/`BatchDestroy`, `Ready`/`ReadyCreatures`,
`PlayFrom`/`PlayTopOfDeck`/`PutIntoPlay`, `ArchiveCard`/`ArchiveFromPlay`). These
are declared in the test as **synonym groups with a reason**, which is a positive
assertion that they *should* look alike, not an allowlist exempting undrawn art.
The distinction matters: `iconFallbackAllowed` was deleted precisely because an
unmapped mechanic is a bug rather than an exemption, and nothing here reintroduces
a way to opt out of drawing a glyph.

### The gate stays at 30

With `effectGlyphs` split, every remaining function over complexity 30 is in a
`_test.go` file and already excluded by the `path: _test\.go` rule. `min-complexity`
stays at 30 and the "next ratchet step" note is deleted rather than carried
forward; ratcheting to 25 would cost four splits in `magefiles` and `deckgen` that
have nothing to do with glyphs.

## Consequences

- `internal/web/icon.go` shrinks from 1518 lines to roughly 200: the ADR 0022
  header, the `glyph`/`decor`/`glyphLine` types, `cardGlyphs`, the family chain,
  and `composeGlyphs`/`fallbackGlyphs`. The continuous-rules renderers, the trigger
  table and the shared noun vocabulary move to `icon_static.go`,
  `icon_trigger.go` and `icon_target.go`; the nine families to
  `icon_<family>.go`. Tests follow their source per `AGENTS.md`.
- The last `gocognit` exclusion leaves `mklv.config.json`, and `.golangci.yaml` is
  regenerated by `mage ci:fix`.
- Roughly a third of cards change what they draw once the sigil axis lands. That
  is the point, but it means the split must land first and separately: a pure move
  diff can be reviewed by confirming every case label survived exactly once, and a
  move-plus-edit diff cannot.
- New sigil SVGs join `web/assets` and must be registered in `galleryIcons`, which
  `TestGalleryShowsEveryIcon` holds equal to the assets directory.
- The four defensive `covered = false` paths are carried across unchanged. None is
  reachable for a shipped card — `TestIconTotality` would be red if one were — but
  they are what forces the three-valued family return.
