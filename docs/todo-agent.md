# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Cite the ADR or doc that decided an item where one exists.

## Engine refactor sweep (survey of `internal/engine`)

File-naming decision for the splits: keep each family's existing top-level prefix
so `ls` groupings stay intact — `target.go` is a bare concept file, so its splits
are `target_*.go`; `effect_condition.go` / `effect_count.go` are part of the
`effect_*` vocabulary, so their splits keep the prefix (`effect_condition_*.go`,
`effect_count_*.go`). Ratchet the split convention into `internal/engine/AGENTS.md`.

Structural (decompose / atomize / recompose):

- **A1** — Fuse the per-zone move-log clones (`CardArchivedFromHand/FromDiscard`,
  `TopOfDeckArchived`, `CardDiscardedFromDeck/FromArchives`,
  `CardPurgedFrom{Discard,Hand,Archives,Deck}`) into one `CardMoved{Player, Card,
Verb, From}` with the "from <zone>" phrasing in a table — the log analog of the
  zone-movement consolidation (ADR 0031). Pure renderers, no `Resolve` divergence.
- **A2** — Extract `randomCardsFrom(ctx, zone, owner, n)`, the random analog of the
  `handCardsWhere`/`discardCardsWhere` gathers, and route `DiscardRandomFromHand/
FromArchives`, `PurgeRandomFromHand`, `ArchiveRandomFromHand` through it (ADR 0031).
- **A3** — Selection mode (chosen / random / each) is a `Selector` Strategy, not a
  node per (verb × zone × mode): fold `PurgeFromHand` / `PurgeRandomFromHand` /
  `PurgeEachFromHand` / `PurgeCreatureFromHand` (and the archive/discard mirrors)
  onto one node per verb carrying a selector. **Do all zones**, not just hand.
- **A4** — `MayPlayOrUseFriendlyHouse` welds play/use: atomize into a play/use axis
  (enum/Strategy) plus a `May` wrapper.
- **A6** — `PurgeArchivedCardThen` fuses a May-gate with a follow-up: compose it
  from a `May` wrapper, a purge gate, and `Then`, like the other `…Then` gates.
- **A7** — Extract the loop constructs (`RepeatWhile`, `RepeatOnCondition`,
  `MayRepeat`, `Overwhelmed`, `CountIs`) out of `effect_condition.go` into
  `effect_repeat.go`; they are not `Condition` predicates.
- **A8** — Unify the flank predicates asked of two subjects (`SourceOnFlank` /
  `ItIsOnFlank` / `ItIsOnNamedFlank`) into one predicate parameterized by subject
  (source vs `ctx.It`), the way `PoolAember{Player}` unified the pool mirrors.
- **A10** — With `Or`/`OrAmount` landed, sweep the ~45 conditions for any remaining
  baked-in disjunction a one-off would reach for; confirm the `ctx.It` atoms are
  what `Or` composes.

Cleanup (splits / comments / naming — not method refactors):

- **B1** — Split `effect_condition.go` (1289 lines, 45 types) into
  `effect_condition_*.go` by category (pool/board, source-position, `ctx.It`,
  turn/play-history), plus A7's repeat file.
- **B2** — Split `target.go` (1479 lines) into `target_*.go` (type + kinds vs
  selection/filter helpers).
- **B3** — Split `effect_count.go` further. `effect_count_produced.go` is extracted
  (the `ctx.Produced` "this way" tallies); the file is now 898 lines. Still to
  pull: the board block (`InPlay` / `CardsPlayed` / `CreaturesUsed` /
  `ExcessCreatures` — the big cluster; `countOfHouse` moves with it) into
  `effect_count_board.go`, then the Æmber and zone-count clusters.
- **B4** — Split `resolver.go` (1246) and `text.go` (1106): separate the role-
  interface declarations from the `*Game` method bodies; group the text helpers.

- **B6** — When A1 lands, move `log_zone.go`'s header comment about `nameMoved`
  onto the single `CardMoved` type.
- **B7** — Comment the struct fields whose unit/why the name doesn't carry
  (`bar.go` prediction bar; `card.go` Assault/Hazardous/SplashAttack bonuses).
- **B8** — Rename `MayPlayOrUseFriendlyHouse` toward KeyForge vernacular once A4
  settles the axis.
- **B9** — `effect_blank.go` is a live mechanic (Shadow of Dis; more variants
  coming) — leave it; just confirm the file earns its own name as variants land.

## Card wording / authoring

- **Festering Touch** reword, e.g. "Play: Choose 2 creatures. Deal 1 damage to the
  chosen creatures with no damage. Deal 3 damage to the chosen creatures with
  damage."
- **Orator Hissaro** could read: "Play: Exalt and ready each neighboring creature.
  For the remainder of the turn, those creatures belong to house Saurian."
- **Borr-Nit** and similar could be atomized and recomposed further (decompose
  fused effects into shared nodes).
- **Memory Chip**: after choosing self house, archive a card from hand (check the
  original printed card text and match it).

## Card catalog / provenance

_No outstanding items._

## Web — mobile, previews, layout

- **Log hover preview on mobile = single tap.** Currently needs a double tap on the
  card name to raise the preview; make it one tap.
- **Facedown cards face down for BOTH players.** A facedown card shows its back to
  everyone; only the controller can hover (desktop) / tap (mobile) to peek at its
  face. Must be visually distinguishable from a faceup card.
- **Card title dynamic shrink (QUESTION).** Yshi's title shrinks by a fixed step
  rather than to the exact needed size. Explain whether we still use predetermined
  shrink amounts, whether dynamic fit is possible, and why smooth fit is/ isn't
  feasible.
- **Glyph bar margin** — remove the margin around the glyph bar (or the surrounding
  boxes) to save space.
- **Glyph strip — residual effect-level unknowns.** Card-level features, all
  triggers (incl. phase triggers), Static/Restrictions/Toll/Replaces/KeyCost, and
  Splash-attack now transcribe. Still falling back to the abstract `glyph-unknown`
  (`*`): the inner `.Then` effect of `ChooseHouseThen`/`ChooseCreatureThen`
  (Restringuntus, Deep Probe, Niffle Grounds) and the granted-ability effects on
  the Blasters, Evasion Sigil, and Rocket Boots. Map each inner effect to a glyph
  (ADR 0022). The remaining empty strips (Dust Pixie, Toad, Mega, etc.) are
  Æmber-bonus-pip-only creatures — the pips render on the card face, not the strip,
  so those are intentional, not gaps.

## Tooling / tests

- **Unused-asset test.** Add a test that fails when an asset in web/assets is not
  referenced by the code (no dead assets).
