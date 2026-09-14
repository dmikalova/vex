# Runtime type conversion via LastingType

## Context

A few cards change a card's type while it is in play. Auto-Legionary is a Saurian
artifact whose `Action:` gives itself five +1 power counters and then moves it to
a flank of the battleline **as a creature** — from that point it is a 5-power
creature that can reap, fight, and be fought, until it leaves play. Effigy of
Melerukh and The Mysticeti do the same to themselves. The card's printed type is
Artifact, but its runtime type is Creature for the rest of its time in play.

This is a different axis from ADR 0026's creatures-as-upgrades, where the runtime
type is read from the card's **attachment** (a card in an upgrade chain reads as
an Upgrade). Here there is no attachment: an artifact in the artifact row becomes
a creature on the battleline. The type has to be stored on the card.

State is flat, pointerless, and comparable (ADR 0005), so the stored type is a
plain enum field, not a behavior closure.

## Decision

`CardCore` carries a `LastingType CardType` field. `TypeUnset` (the zero value)
means the card keeps its printed type; any other value overrides it for as long
as the card stays in play. `resetCore` zeroes the whole core when a card leaves
play, so the override never outlives the card.

`TypeOf(id)` is the single authoritative reader of a card's current type, in
priority order:

1. A card in an upgrade chain (`HostPlus != 0`) reads as `Upgrade`, whatever its
   printed or lasting type (ADR 0026).
2. An in-play card with a non-unset `LastingType` reads that type.
3. Otherwise it reads its printed type.

Every in-play gate that used to branch on `g.cat.def(id).Type` now branches on
`g.TypeOf(id)` — combat, use/reap/fight legality, armor refresh, aember capture,
control placement, text-blanking. Play-time hand checks and
`creaturesPlayedThisTurn` deliberately keep the **printed** type: turning into a
creature is not "playing a creature" (no creature entered play, nothing was
played), so it must not count toward "for each creature you played this turn."

The conversion itself is one resolver method, `PutIntoBattlelineAsCreature`,
behind a single effect node `TurnIntoCreature` (in `effect_battleline.go`,
alongside `MoveToFlank`). The effect prompts the controller for a flank and moves
the card from the artifact row onto that flank, topping up its `ArmorRemaining`
to its full armor so it can absorb hits as a creature this turn. The card keeps
its exhaustion, Æmber, and power counters — its counters now count toward its
power because `Power` reads them regardless of type.

The view shows power and armor only when `TypeOf(id) == Creature`, so an artifact
(printed or a creature-as-upgrade, ADR 0026) shows no power/armor line — "null
power" is a display concern, not a change to `Power(id)`'s signature, which stays
`int` for its 200-plus callers.

## Two routes for a type change: stored vs. live

`LastingType` is the right tool for **one** shape of type change, and the wrong
tool for another. Pick the route by how the change ends, not by what type it
produces.

**Stored (`LastingType`) — a permanent self-conversion.** The card converts
**itself**, the change never reverts, and it ends only when that card leaves
play. Auto-Legionary, Effigy of Melerukh, and The Mysticeti are the whole set:
each turns itself into a creature for good, so the type is a plain enum written on
the card and zeroed by `resetCore` on exit. There is no external source, no
duration, and no counter — nothing to unwind but the card leaving play. A single
field holds every such card at once, because two of them can never fight over the
same card (a card only ever converts itself, once, one way).

**Live (a `ConstantAbility`, never stored) — a "considered an artifact" grant.**
The card is a given type only while some other condition holds — a counter on
it, a source card in play, a while-condition. Deanimator is the model: "each card
that has a mineralize counter on it is considered an artifact." Per ADR 0024 this
is a `ConstantAbility` composed on top of a `CounterInPlay{Kind, Target: This}`
read, computed live at read time (the way `Power`/`Armor` fold `constantBonus`),
**not** a `LastingType` write. It reverts for free: when Deanimator leaves play or
the counter is removed, the next `TypeOf` read simply no longer sees the grant, so
there is nothing to pop and no source lifetime to track. Two Deanimators compose
the same way — each grant is recomputed every read. The live route is **not built
yet** — no implemented set has a card that needs it — so it is parked in
[todo-future-set.md](../todo-future-set.md), to build alongside Deanimator (or the
first live-conversion card) when its set is stood up.

For the live route, `TypeOf` grows one step after the `HostPlus`/`LastingType`
checks: ask whether any active constant ability converts this card's type, exactly
as the stat reads scan for bonuses. Do not store the result.

### What to look out for — signals you need the live route (or, later, a stack)

When implementing a type-conversion card, read the card text for the ending
condition and route on it:

- **"becomes a creature" / "is an artifact until it leaves play", self-targeted →
  stored `LastingType`.** Permanent, self-sourced, one-way.
- **"is considered an artifact" (or any type) gated on anything — a counter, a card
  in play, a while-condition → live `ConstantAbility`.** The tell is that the change must end
  on a lifetime **other than the converted card's own** (a _source_ leaving play, a
  counter being removed, a condition ceasing). Never store this; storing it would
  need an unwind path the live layer gives for free.
- **Two _stored_ overrides could apply to one card at once, on different lifetimes
  → this ADR's single field is no longer enough.** No such card exists today —
  every known type change is either a permanent self-conversion (stored) or a live
  "considered" grant (not stored), so the two routes never overlap. If one ever
  appears, do not add a second field: fold `LastingType` into a global, flat,
  comparable `TypeEntry{Card, Type, Source}` LIFO side-table exactly like the
  control stack (ADR 0028) — self-conversions become self-sourced entries (ended by
  the card leaving play), gated conversions name their source, entries supersede by
  `(card, source)` to stay board-bounded (ADR 0024), and `TypeOf` reads the top of
  the stack. The one genuinely new decision at that point is the precedence between
  a stored override and a live constant grant on the same card; KeyForge resolves
  overlapping continuous effects by last-applied-wins, so record that rule when it
  first matters.

## Consequences

- One place decides a card's current type. Adding a new type-conversion mechanic
  is a new `LastingType` write, not a new branch in every gate.
- The stored type never leaks: it lives in `CardCore`, which `resetCore` zeroes
  on every real leave-play exit. Upgrade mode is deliberately **not** stored here
  (it is read from `HostPlus`), because the shed-upgrade discard path does not
  reset the core (ADR 0026).
- Converting an existing card's `def.Type` gate to `TypeOf` is behavior-neutral
  for every card that never converts or attaches (`TypeOf == def.Type` then), so
  existing tests and their coverage are unchanged.
- `TurnIntoCreature` renders "move it to a flank of your battleline as a
  creature", matching KeyForge's wording for Effigy of Melerukh and The
  Mysticeti; the pronoun reads against the effect that named the card first.
- The mechanic needs no rulebook term (ADR 0018): the printed clause carries its
  whole meaning, and the type it produces is an existing card type.
