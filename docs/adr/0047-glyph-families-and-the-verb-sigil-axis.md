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

### Structure is punctuation drawn between glyphs, never a third badge

`decor` and `sigil` are both drawn **on** a noun: the subject glows, the verb wears
a corner mark. Neither can say where one clause ends and the next begins, so an
effect that puts more than one clause on a line renders them fused. An
`engine.Conditional` concatenated its `Then` and its `Else` with nothing between
them, so "if x, destroy a creature, otherwise draw a card" drew as destroy,
creature, draw, card and read as one sentence; `engine.ChooseOne` did the same,
prefixing the choose glyph and then running every option together.

**Structure is a third axis, `punct`, drawn between glyphs rather than inside
them.** Three marks, defined once and reused everywhere:

- a **gate** mark opening a condition clause,
- a **branch** mark between a `Conditional`'s `Then` and its `Else`,
- a **divider** between the options of a `ChooseOne`.

A glyph carries the mark that precedes it, so a family punctuates an
already-composed run by stamping its head (`punctuate`), exactly as `markVerb`
stamps a verb. What each mark draws is `punctMark`'s business alone
(`view_icons.go`), the same seam `sigilMark` owns for a verb, so promoting a mark
to art touches no family file.

The marks are **not** a badge on the glyph, because a badge has to mean something
about the noun it sits on, and a tint meaning "enemy" on one noun and
"hypothetical" on another is exactly the ambiguity the two badge axes exist to
remove. Structure is about the space *between* two mechanics, so it is drawn
there.

The existing arrow is **not** reused for the gate. The arrow is a flow mark
meaning "and then this happens to that" — æmber → creature, discard → hand — and
overloading it with "if" is the same conflation the rules-text side is currently
undoing for a self-repeat. It joins `punct` as a fourth member so there is one
punctuation channel rather than two, but it keeps its own character and meaning,
and `TestPunctMarkCoversEveryPunct` fails if two marks ever draw the same one.

The cost is horizontal. Measured over all 1231 materialized defs the median card
draws 3 glyph slots and the widest 16; the estimate after the whole grammar lands
is a median of 4 and a widest of 22. `iconStrip` scales `.card-icons-fit`
horizontally through `cardFitScript`, so a busy card squeezes rather than wraps.
That squeeze is accepted as-is and the fit behaviour is deliberately untouched; it
can be reworked if it becomes untenable.

### A comparison is a mark on the quantity badge, not a seventh picture

`engine.Comparison` — `AtLeast`, `AtMost`, `Exactly`, `MoreThanYou`,
`MoreThanOpponent`, `Even`, `Odd` — is the shared vocabulary behind the power
filters and most conditions. Drawing seven pictures for it would give the strip
seven nouns that are not nouns; "power 3 or more" is the power glyph with a 3 and
an at-least mark, not a new drawing.

**A comparison qualifies the quantity badge the glyph already prints.** It is a
field on `glyph` typed as the engine's own `Comparison`, so the vocabulary cannot
drift from the enum it draws, and `comparisonMark` (`view_icons.go`) owns what
each member draws — text first, exactly as `sigilMark` allows for a verb, so
promoting one to art later touches no family file.

It rides the existing badge rather than taking a third badge axis beside `decor`
and `sigil`, for the reason above: a strip where every glyph wears one more badge
is noise, and a comparison is a qualifier on a number the badge is already
printing.

Two members print no number at all. `Even` and `Odd` ignore `Amount` — they
compare a quantity's shape, not its size — so their mark is the whole badge, which
is why the badge draws for a bare comparison as well as for a count. The two
relative comparisons share one mark: both mean "more than the other player's
pool", and which pool is the subject is already carried by the glyph's
enemy/friendly decor, so a second mark would encode the subject twice.
`TestComparisonMarkCoversEveryComparison` walks `engine.Comparisons()`, so a
comparison added to the engine cannot reach the strip undrawn.

### A timing window is drawn by composition, not by a fourth axis

The strip had no duration channel at all. `internal/web` mentioned `Duration` in
one place — `icon_compose.go` recursing through `engine.ForDuration` into whatever
it wrapped — so the wrapper rendered as its child and the window was dropped. Two
measured pairs came out byte-identical: `BatchDestroy` and
`DestroyEachCreatureAtEndOfTurn` both drew `[destroy, → each creature]`, and
`GainAssault` and `GainAssaultUntilNextTurn` both drew `[assault, → this
creature]`.

**An effect bound to a window appends the `phase-turn` glyph to its run** —
`inWindow`, and `windowed` for a wrapper that still holds a covered flag. The
end-of-turn wipe reads as destroy, each creature, turn marker. `phase-turn` is
already the trigger glyph for the start and end of a turn, so the marker is a word
the player has met rather than a new one, and the cost is one horizontal slot on
the handful of effects that have a window.

*Which* window is deliberately not drawn. The marker says "not now"; saying more
would need a picture per `Duration`, and six pictures for a distinction the printed
text already carries is the seventh-picture mistake the comparison marks avoid.

Declaring the assault pair synonyms was the alternative — a lasting buff and a
permanent one are one act with a different shelf life. They are drawn apart
instead: the marker costs one slot on one rare effect, and how long a buff lasts is
a thing a player plays around. `TestTimedAssaultDrawsItsWindow` and
`TestDelayedWipeDrawsItsWindow` pin both pairs.

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
- A glyph's `arrow bool` becomes `punct punct`, so there is one punctuation
  channel and `arrowTo` is one of its callers rather than its own mechanism. The
  CSS class `.card-glyph-arrow` becomes `.card-glyph-punct` with a modifier per
  mark.
- The four defensive `covered = false` paths are carried across unchanged. None is
  reachable for a shipped card — `TestIconTotality` would be red if one were — but
  they are what forces the three-valued family return.
