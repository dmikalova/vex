# 46. Log entries form an enumerable catalog the style gallery covers by sampling games

This decision records how the `/style` gallery gains a section showing every kind
of rendered card text and game-log line. It is the display counterpart to ADR
0014 (the style gallery renders real components), ADR 0018 (a closed catalog is
complete by construction), and ADR 0022 (a Visitor guarded by a totality test).

## Context

The log narrates resolved outcomes (ADR 0011): ~130 concrete `LogEntry` variants
across `log_<family>.go`, each rendering its own past-tense sentence, plus the
card-link and icon spans `RenderRecord` layers on. Card text (`RenderCardRules`)
renders a parallel present-tense imperative from the effect AST. Between them
they are the game's entire player-facing prose surface, and there is no one place
to _see_ that surface — to look at an example of every kind of line the game can
print and confirm it reads well.

Two naive shapes both fail:

- **A curated catalog** — one hand-written example per variant — duplicates the
  wording the engine already owns, drifts the moment an entry changes, and is
  exactly the snapshot-in-markup that ADR 0014 forbids on the gallery.
- **A random-game log dump** — play a game, print its log — is cheap and real but
  repetitive and blind: it over-shows the common lines, under-shows the rare ones,
  and can never tell "this kind never came up" from "this kind does not exist."

The blind spot is the crux. A pass that discovers kinds by walking `%T` over a
game's log can enumerate what it _saw_, but not what it _missed_, so it has no
stopping condition, no coverage number, and no way to flag a gap. The rendering
surface, though, is not combinatorially infinite: the variation that changes how
a line reads is a small bounded set of shapes (the entry variants) crossed with a
handful of **combiner** axes (article `a`/`an`, singular/plural, small-number
words, `Target` phrasing, `Condition` And/Or/negation, `Count` phrasing). Card
names are substitutions into those shapes, not new shapes.

The engine already holds the shape list implicitly: `TestLogEntryText` is a
near-exhaustive, hand-maintained table that constructs one instance of every
entry variant and pins its wording. It is the enumeration we need — trapped in a
`_test.go` file where no non-test code can reach it.

## Decision

Make the log's entry kinds an **enumerable catalog in the engine**, then have the
gallery **cover the catalog by sampling real games**, falling back to the
catalog's own constructed instance only for kinds a game run never produced.

- **Catalog (engine, the keystone).** Promote `TestLogEntryText`'s constructed
  instances into an exported `engine.LogEntrySamples()` — one representative
  `LogEntry` value per variant — co-located with the entries the same way
  `Keywords()` sits beside `Keyword`. A **totality test** discovers the variant
  set the way the card tests discover cards — by parsing the package's own source
  for every type with a `Text(Namer)` method (the `LogEntry` interface, distinct
  from `Effect.Text()`, which takes no argument) — and fails if any variant has no
  sample. `TestLogEntryText` becomes a consumer of the catalog: it renders every
  sample and asserts its wording against an index-aligned table, so the pins are
  preserved but the instances live in one place. This gives the gallery a _known
  target set_: a stopping condition, a coverage number, and the difference between
  "not observed" and "does not exist."

- **Sample real games against the catalog (web).** A collection loop plays seeded
  games and buckets every produced `Record` by its entry's shape (`%T`), keeping
  real bubbles, until every catalogued kind has at least one real example or a
  1000-game cap is hit. Coverage is measured against the catalog, not guessed.

- **Cover, then drill in.** The hero view is a greedy **set-cover** of log
  _bubbles_: pick the fewest real bubbles whose union covers every catalogued
  kind, rendered with the production `logBlockView`. Clicking a bubble opens the
  full log of the game that produced it — the real `logBlocks` render, bubble
  highlighted — so every specimen is anchored in the context that made it. A
  secondary index isolates single shapes and the combiner axes a bubble cannot
  isolate on its own (article, 1-vs-many), and the seeded whole-game dump stays
  as a replayable end-to-end example.

- **Synthetic fallback, three-state gaps.** A catalogued kind a game run never
  produced renders from its _constructed_ `LogEntrySamples()` instance, badged
  "synthetic — not observed in N games", with no drill-in context. A gap is thus
  three-state: **observed** (real bubble + context), **synthetic** (renders,
  flagged, no context), and **truly missing** — which a complete catalog prevents.

- **Card text, in parallel.** The card-text section shows one specimen per text
  feature (trigger kinds, `Target` shapes, `Condition`/`Count` composition,
  durations, keywords, static/constant/granted, Enhance), each a _random real
  card_ matched by a predicate over `cards.All()` with a re-roll, plus a
  **Combiners** subsection that isolates the text-helper axes on purpose. It
  reuses `RenderCardRules` / `printedFace`; empty partitions show as visible gaps.

The gallery links to `/rulebook` for term definitions rather than duplicating
them (ADR 0018 owns that surface).

## Consequences

- **The catalog is the load-bearing addition.** Promoting the test table to an
  exported, totality-gated registry is worth doing on its own: it closes the
  enumeration that made every other option blind, and it is the only part inside
  the 100%-coverage engine gate. The gallery is a pure consumer in the ungated
  `internal/web`.

- **Tests assert coverage, never wording (ADR 0014).** Web tests check that every
  catalogued kind is sampled-or-flagged and that specimens render without panic;
  they never pin markup or CSS. Wording stays pinned by `TestLogEntryText` in the
  engine, now driven off the same catalog.

- **A new entry variant fails the build until catalogued.** Adding a `LogEntry`
  without a sample lights up the totality test, which reads the source for every
  `Text(Namer)` type — the same pressure ADR 0018 and ADR 0022 apply to their
  catalogs. The gallery then shows it for free.

- **Bounded, not exhaustive.** The gallery gives practical coverage — every shape
  once, each combiner axis isolated — not the cartesian product of shapes and card
  names, which is infinite and uninformative. This is the equivalence-partition
  reading of "see a bunch of different examples."

- **Staged.** Stage 1 is the engine catalog + totality test. Stage 2 is the
  sampling and set-cover collection. Stage 3 is the two web sections and their
  coverage tests. Each stage is independently green; see `docs/todo-agent.md`.
