# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Work that has **no consuming card in an implemented set yet** does not belong
  here — park it in [todo-future-set.md](todo-future-set.md), keyed to the set that
  first needs it.
- Cite the ADR or doc that decided an item where one exists.
- **When you need a decision from the human, ask in the reply itself, grill-me
  style** — a numbered list of `❓ **Q1** - **title**: <question>` with a `➡️`
  recommended answer under each — at the end of the turn, then stop. Do **not**
  reach for an interactive question tool: under autopilot it is auto-answered with
  "work autonomously" and the human never sees it. Plain-text questions at the end
  of the turn are the channel the human actually reads.

## Maverick houses complete every OnePerHouse cycle to all 9 Houses (ADR 0036)

Deck generation will grow a **maverick House pod** (deck-generation.md §1): very
rarely a whole pod is a House not in the deck's Set at all, drawn wholesale from
another Set. That breaks the current per-set `OnePerHouse` gate: a maverick Saurian
pod dropped into an Age of Ascension deck fires the Shard cycle for a House AoA has
no Shard for. So **when maverick houses land**:

- Expand the `OnePerHouse` completeness gate (`Set.validateClusters`) from the
  set's own Houses to **every House a maverick pod can introduce** — effectively all
  nine. AoA (and every Shard-bearing set) then needs a Shard for **all** Houses,
  including Saurian and Star Alliance, which AoA never printed.
- Add those missing Shards as `//go:build todo` stubs so the gate stays honest
  until they are implemented (Saurian and Star Alliance were introduced after AoA,
  so their Shards are Vactrol-generative, like Master of 4/5).
- The Ambassador and Plant templates already cover this for free — they materialize
  per partner House from `DeckHouses`, so a maverick House gets its Ambassador/Plant
  with no new card. Only the Shard cycle, which is concrete per-House cards, needs
  the stubs.

Until then, a legacy-drawn Shard does not fire the cycle either: `expandClusters`
reads only the drawing set's own `clusters`. Revisit both together when maverick/
legacy whole-pod overlays are built.

## Engine — effects & mechanics: decompose fused one-off effect nodes

A survey found many effect nodes named after a single card that weld several
actions into one `Resolve` — one-offs instead of a hardened shared core, which has
been breeding bugs. Retire them by building ~6 shared primitives and re-expressing
each card as a composition. Build each core once, then knock out its whole group;
land each green with its `effect_*_test.go` and ratchet the shape into
`internal/engine/AGENTS.md` / `docs/style-guide.md`. Groups, cheapest first:

- **Core A — "choose any number, accumulating a value into `ctx.Produced`"
  multi-pick producer.** Retires the forge-fusion nodes and the multi-zone shuffle
  selection. Decision: the producer picks any number through `pickCards`, tallies
  the chosen (count and/or power) into `ctx.Produced`, and the follow-up reads that
  tally instead of a local `len`/sum.
  - `SacrificeToForge` (Obsidian Forge) → `Sequence{<pick-any-number>, Destroy,
ForgeKey{ReducedBy: CreaturesDestroyed}, PurgeSource-on-forge}`. `ForgeKey`
    already has `ReducedBy`; `Destroy` already tallies `ctx.Produced.Destroyed`.
    Today it uses local `len(chosen)` — switch it to the tally.
  - `DestroyFriendlyCreaturesToForge` (Might Makes Right) → same producer tallying
    **power**, gated `Then` above the threshold.
  - `ShuffleChosenCreaturesFromZones` → the producer over a multi-zone candidate
    union feeding the existing shuffle-movers.
- **Core B — cross-zone `Named` selector with an EXPLICIT destination** (human:
  do not infer the destination from the card's zone — pass it explicitly; the
  primitive only dispatches the right per-source-zone move). Retires:
  - `ReturnNamedToHand` (battleline+discard → hand).
  - `SearchForName` (deck+discard → hand, choose-one/take-all + reveal + gate).
  - `ShuffleNamedFromDiscardIntoDeck` (discard → deck).
- **Core C — parameterized resource-redistribution node** (read → clear →
  redistribute one unit at a time), with lethal-destruction as a SEPARATE follow-up
  pass, not fused in. Retires `RedistributeCapturedAember` (Equalize) and
  `RedistributeDamage` (Entropic Manipulator; its lethal pass becomes a conditional
  destroy after the shared core).
- **Core D — deck "dig until" family (DEFERRED — owned by another agent right
  now).** See the item below; do not touch `DiscardDeckUntil` /
  `RevealDeckUntilHouse` while that work is in flight.

Pending human decisions before building (asked in the reply): the house-choice **wager** unification
(`WagerOpponentChoosesChosenHouse` → a house-choice reaction; revises ADR 0035).

Explicitly REJECTED (atomic, not fused — do not re-propose): `EndTurn`,
`PutDiscardedIntoHand`, `GainTextBox` (one semantic op), `PlayFromOpponent` (HARD —
divergent ownership/play-legality, would become a branchy blob).
(`MakeItsHouseActive` → `ChangeActiveHouse{To: HouseChoice}` is DONE.)

## Card wording / authoring

_No outstanding items._

## Card catalog / provenance

_No outstanding items._

## Web — mobile, previews, layout

_No outstanding items._

## Tooling / tests / docs

_No outstanding items._
