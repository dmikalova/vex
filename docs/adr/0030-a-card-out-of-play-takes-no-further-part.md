# A card out of play takes no further part

## Context

Two recurring bug classes both come from acting on a card that is no longer in
play:

- **Stale write** — a counter, Æmber, damage, stun, ward, or enrage marked on a
  card that already reached its discard pile, where it later leaks back into the
  game. [REDACTED] plus Strange Gizmo was the reported case: forging a key mid-
  window put Æmber back onto the already-destroyed [REDACTED].
- **Out-of-play source ability** — an ability resolving from a source that an
  earlier ability in the same window removed from play.

Both were patched at individual call sites. The guarded accessor `stateOf` returns
`nil` for a card not in play and so no-ops most writes; the choose-house window
skips an out-of-play source at its own loop (`game_turn.go`); the destruction
window skips an out-of-play source at its own loop (`game_leaves_play.go`). But the
write path for Æmber (`addAmberOn`) bypassed the accessor, and the forge-key window
had no source guard at all. Scattered guards mean "know every call site" — the
exact thing we want to stop needing to know.

Master Rulebook §190 settles the rule: an ability checks that its source is in
play at the instant it begins resolving, and once it has started it is locked in.
Vactrol adopts this uniformly. It is **not** a divergence — the ordinary trigger
windows were simply missing the guard the destruction window already had.

The model, stated precisely: **a card that has left play can only ever be the
subject of a resolving ability, never its source.**

- `Destroyed:` is a triggered ability whose trigger is the creature being _tagged
  for destruction_. It resolves while the creature is still in play (tagged, its
  buffs and power still computed), before it is considered destroyed — so its
  source is present. It runs in the destruction window's own loop
  (`destroyTogether`), not `triggerAbilitiesAs`.
- The "after … is destroyed" and "after … leaves play" reactions (Neffru via
  `emitAfterCreatureDestroyed`, Pile of Skulls via `emitAfterEnemyDestroyed`) are
  triggered abilities on a _different, in-play_ card that watches the after-window.
  The reacting card is the source and is present; the card that left play is the
  subject (`it`).
- `TriggerLeavesPlay` is an engine-internal trigger name, **not** a KeyForge term,
  and it fires (`emitLeavesPlay`) while the card is still listed on the board.
- A **tactic** is the one exception: it never enters play, so it resolves its own
  `Play:` ability while not on the board. "Source in play" does not apply to a
  card that has no in-play state to begin with — the tactic is the actor of its
  own play, not a card an earlier ability removed.

So for a card that _has_ a place in play, no legitimate resolution has an
out-of-play source: an out-of-play source arises only from an earlier ability in
the same window taking the source out. A tactic is not such a case — it is out of
play throughout, by construction, and its play is legitimate.

## Decision

Two structural guards, each the single seam for its concern.

**1. Guarded write path.** Every effect- or `Resolver`-driven write to a card's
`CardCore` — Æmber, damage, stun, ward, enrage, power counters, all of them — goes
through the guarded accessor (`stateOf`, `nil` for a card not in play) and silently
no-ops on a card that has left play. `addAmberOn` is routed through it rather than
writing `CardCore.Amber` directly. The effect-level `resolverInPlay`-before-write
checks scattered through the effects come out, because the write itself is now safe.

Leave-play **teardown** is exempt and keeps direct field access: `resetCore`,
`clearCounters`, and releasing a card's on-card Æmber as it leaves are not an
ability acting on the card — they are the act of removing it. The guard means "an
ability cannot touch a card that already left," not "no code may write a card that
is leaving."

**2. Source-in-play resolution guard.** A triggered ability resolves only if its
source is in play at the instant it begins resolving. This is one guard in
`resolveWindow` (the shared trigger-window resolution loop) — `if !g.inPlay(src) &&
g.cat.def(src).Type != Tactic { continue }` — applied per ability so it covers both
a source removed between two
sibling cards' abilities and a source that removes itself between its own two
abilities. The single exception is a **tactic**: it resolves its `Play:` while not
in play (a tactic never enters play), so the guard lets a tactic through and drops
only a non-tactic source that an earlier ability removed. Every other trigger the
engine routes through `triggerAbilitiesAs` has a present source by construction
(reactors are in play; a tagged creature is still in play; the leaves-play emit
fires while the card is still listed), so the guard only ever drops a source an
earlier ability removed. The choose-house call-site guard is deleted, and the
forge-key window — which never had a guard — is covered for the first time. The
destruction window runs its own re-gathering resolution loop
(`resolveDestroyedWindow`) rather than the up-front-ordered `resolveWindow`, but
both resolve each entry through the shared `resolveTriggered` step, so the source
guard lives in one place and applies identically everywhere.

## Consequences

- Stale writes become impossible at the write, not at each effect; the per-effect
  `resolverInPlay`-before-write guards are removed.
- An out-of-play source drops its remaining abilities everywhere, by one line; the
  [REDACTED]+Strange Gizmo case and the forge-key window are covered by the same
  guard, and the choose-house special case stops being special.
- A card that has left play, appearing as a **subject**, stays handled where
  subjects are handled — the write no-op above and target/selection filters — not
  by the source guard.
- Semantic control flow that reads `inPlay` to decide a follow-up ("deal damage,
  then if it was destroyed draw a card") is untouched: that is a condition, not a
  guard, and it stays.
- The destruction window's guard and `triggerAbilitiesAs`'s guard are siblings of
  the same one-line shape in two different loops, not duplicates to merge.

## Refinement: a departed subject reads from its last-known scalars

The model above leaves one half open. A card that left play "can be the subject of
a resolving ability" — but a subject the ability still needs to _read_. When an
effect removes a card and a later effect in the **same resolution** reads that
card's mutable state (Power of Fire destroys a friendly creature, then each player
loses Æmber equal to half **its power**), the live read returns the zeroed core
`resetCore` left behind: printed power only, with power counters, buffs, upgrades,
and constants all dropped. The ability must instead read the value the card had the
instant before it left.

The engine already did this ad hoc for three dimensions — `Produced.DestroyedPower`
(summed power), `Produced.Neighbors` (position), and `ctx.ItController` (controller)
— each captured before the exit and consumed by one reader. This refinement
generalizes the same idea rather than adding a persistent snapshot:

- **Ephemeral, not stateful.** The last-known values live on the resolving
  `EffectContext` (`ctx.Departed`, a `map[LocalID]departedSubject`), never in
  `GameState`. `EffectContext` is short-lived execution context, so this pays none
  of the flat/comparable cost ADR 0005 imposes on `GameState` and grows no undo
  snapshot. A fresh context per ability means the captures never leak between
  abilities.
- **Captured at the exit boundary.** `captureDepartingSubject(ctx, id)` records a
  card's power, Æmber-on-card, and damage the instant before an effect removes it.
  It is called by the effects that remove a card they will keep referencing —
  `Destroy` captures every creature it destroys, `DealDamage{IfDestroyed}` captures
  the creature it deals lethal damage to (alongside the neighbor snapshot it
  already took).
- **Consulted only once the card is gone.** The scalar readers route through
  `ctx.powerOf` / `ctx.amberOn` / `ctx.damageOn` (`PowerOfChosen`, `AemberOnThis`,
  `DamageOnThis`, and the per-target `AemberOnIt` / `DamageOnIt`). While the card is
  still in play its live value is authoritative; only after it has left play do they
  fall back to the captured value. This is invisible to card authors — the reader
  is the same node either way.

Only these three mutable dimensions are captured. House, traits, printed keywords,
printed power, and bonus icons live on the immutable `CardDefinition` and survive a
card leaving play unaided; a departed subject's _granted_ keywords/house are not
read by any card, so no snapshot of them is kept. The write no-op (`stateOf` returns
`nil` off-board) is unchanged: a departed subject is still read-only — this only
lets an in-flight ability _read_ what it was.

## Refinement: a card active in its owner's discard pile is a second source exception

A tactic is not the only card that resolves an ability while off the board.
KeyForge has cards that stay live **in their owner's discard pile** and trigger
from there — Relentless Creeper (Mass Mutation #029) returns _itself_ from the
discard pile to hand after its controller chooses Dis. Its `AfterChooseHouse`
ability resolves with the card sitting in the discard pile, so the source-in-play
guard would otherwise drop it.

This is a bounded, opt-in exception, not a hole in the model:

- **Opt-in per card.** A card carries the exception only by setting the
  `TriggersFromDiscard` field (`WithTriggersFromDiscard()`); every other card is
  unaffected. The guard reads `if !g.inPlay(src) && def.Type != Tactic &&
!g.activeInDiscard(src) { return false }`, where`activeInDiscard` is true only
  for a `TriggersFromDiscard` card sitting in its owner's discard pile.
- **Scanned at one window only.** The choose-house window is the sole site that
  gathers these abilities: after adding the in-play sources, `ChooseHouse` scans
  the acting player's discard pile and adds each `TriggersFromDiscard` card's
  `TriggerAfterChooseHouse` ability. No other trigger window looks in the discard
  pile, so a discard-active card fires only where its printed text says it does.
- **Still read-only as a subject elsewhere.** The card is a legitimate _source_
  only for its own discard-pile ability; everywhere else the write no-op and
  target/selection filters keep it read-only, exactly as for any other off-board
  card.

So the precise model becomes: a card that has left play can only ever be the
subject of a resolving ability, never its source — **except** a tactic resolving
its own `Play:`, and a `TriggersFromDiscard` card resolving its own choose-house
ability from its owner's discard pile.

## Refinement: a deferred "Leaves Play:" window is a third source exception

A `Leaves Play:` ability is the one trigger whose whole purpose is to fire
*because* its card left play. It escapes the guard today only by timing: the card
is still listed on the board when `emitLeavesPlay` resolves it, so `inPlay(src)`
is still true.

That timing is what makes a multi-card move non-simultaneous. When an effect moves
several cards, the first card's `Leaves Play:` ability resolved before the rest had
moved, so it could destroy a card the same effect had already selected but not yet
reached — and the outcome depended on the order the selection happened to be
visited. `simultaneously` now holds those windows until the whole batch has moved,
which means they resolve with their source already gone.

The exception is as bounded as the other two:

- **Gathered on the board, resolved off it.** `leavesPlayWindow` collects the
  abilities while the card is still listed, so its controller, its grantors, and
  its ability list all read exactly as before. Only resolution moves.
- **Marked, not inferred.** The entry carries `fromLeave`, set at gather time and
  nowhere else. The guard reads `if !t.fromLeave && !g.inPlay(src) && …`, so no
  other trigger gains the exemption by accident.
- **Drained by the batch that opened it.** A nested batch does not flush; the
  outermost one owns the moment, so the queue never outlives the effect that
  created it. It is runtime scheduling, not game state, and stays off `GameState`
  and out of the undo snapshot (ADR 0005).

So the final model: a card that has left play can only ever be the subject of a
resolving ability, never its source — **except** a tactic resolving its own
`Play:`, a `TriggersFromDiscard` card resolving its own choose-house ability from
its owner's discard pile, and a `Leaves Play:` window gathered while its card was
still on the board. Pinned by `TestLeavesPlayWaitsForTheWholeBatch` and
`TestPutFromPlayMovesTheWholeSelectionTogether`.
