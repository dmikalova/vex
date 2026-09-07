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
`triggerAbilitiesAs` — `if !g.inPlay(src) && g.cat.def(src).Type != Tactic {
continue }` — applied per ability so it covers both a source removed between two
sibling cards' abilities and a source that removes itself between its own two
abilities. The single exception is a **tactic**: it resolves its `Play:` while not
in play (a tactic never enters play), so the guard lets a tactic through and drops
only a non-tactic source that an earlier ability removed. Every other trigger the
engine routes through `triggerAbilitiesAs` has a present source by construction
(reactors are in play; a tagged creature is still in play; the leaves-play emit
fires while the card is still listed), so the guard only ever drops a source an
earlier ability removed. The choose-house call-site guard is deleted, and the
forge-key window — which never had a guard — is covered for the first time. The
destruction window keeps its own copy of the guard because it runs its own
resolution loop (`destroyTogether`), which `triggerAbilitiesAs` does not reach.

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
