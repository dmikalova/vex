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
