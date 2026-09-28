# Reused effect shapes

Effects recur in **shapes**: the same small field cluster and the same logic to
turn it into a value, a phrase, or a validation error, across unrelated
mechanics. When a shape shows up in a third effect, factor its logic into one
helper the effects share. Each effect keeps its own authoring fields, so the
card API stays ergonomic, but not its own copy of the behavior. These are the
shapes factored so far; use the helper rather than writing a second copy.

- **Scale by a `Per` count.** "N for each ..." is a base times a `Per Count`.
  `scaled(base, per, ctx)` (`effect_count.go`) is the value companion to
  `forEach`'s text; `Draw`, `GainAember`, `DealDamage`, `AddPowerCounter`,
  `CaptureAember`, and the economy amounts all scale through it. `Per` means
  the same on every effect: multiply one target's amount. Re-running with a
  fresh target is `Times` (`Repeat.Times`, `CaptureAember.Times`), never an
  overloaded `Per`.
- **Amount, or an alternative magnitude.** A fixed `Amount` or one other way to
  say the size (a `By` pool share, an `All`/`Fully` flag) validates through
  `errAmountOr(name, alt, amount, altSet)` (`effect.go`), the shared "set one,
  not both" check beside the `errUnset*` helpers.
- **How much Æmber against a pool.** `StealAember`, `LoseAember`,
  `CaptureAember` (and, minus the share, `GainAember`) take a fixed base scaled
  by `Per`, or a `By Loss` share (`Half`, `AllBut(n)`, `AllAember`). The
  helpers in `effect_aember.go` carry it: `poolAmount(base, by, per, ctx, pool)`
  for the value (uncapped; the caller `min`s against the pool and records the
  capped figure) and `aemberObject(amount, by, possessive)` for the printed
  object. A `bool` like Capture's `All` is `By: AllAember` internally; keep an
  authoring `bool` only where the printed wording differs ("captures **all**
  your opponent's Æmber" has no "from" clause).
- **Count what the previous effect produced.** "For each ... this way" is a
  `Count` reading a tally the prior effect left on `ctx.Produced`
  (`CreaturesHealed`, `DamageHealed`, `CardsDestroyed`). A new tally is a field
  on `Produced` plus a small `Count`, never a fused effect. Measure a tally
  about a card that leaves play **just before** it is removed, so a modifier it
  carried still counts and a ward that keeps it in play does not:
  `Destroy.destroy` snapshots each creature's board `Power` before
  `DestroyEachFrom` and adds it to `Produced.DestroyedPower` only for creatures
  that left, which `PowerDestroyedThisWay` reads (Might Makes Right). A `Count`
  used only by a `CountIs` threshold still implements `CountText`, though only
  its `CountClause` renders.
- **Move a card between zones.** Archive, Discard, Purge, Shuffle, and Put are
  one unexported mechanism (ADR 0031): a source zone, a selection strategy
  (chosen / any-number / all / random / top-N, with house/type/name/trait
  filters), a destination, and the optional `ctx.It`/`ctx.Produced` side
  effects, with each KeyForge verb a thin authoring struct that fixes the
  destination and delegates. A new zone-movement effect adds a selection
  filter or a destination, never a new bespoke type; fold the old families into
  it as their cards are touched. Authors write the verb the card prints, not a
  generic `Move`.
- **Move a card that could be in several zones to one destination.** Faygin
  (`ReturnNamedToHand`: play/discard to hand), the search tutors
  (`SearchForName`: deck/discard to hand), and Song of Spring
  (`ShuffleIntoDeck`: play/hand/discard to deck) share
  `crossZoneMover{Player, Dest, Sources}` (`effect_cross_zone.go`). It gathers
  with a `keep` predicate (`gather`) and sends each pick through the
  per-`(origin, Dest)` move from the zone the card actually sits in
  (`move`/`originOf`). The **destination is passed explicitly**; a mover never
  infers it. Each verb keeps its own printed text, filter, prompt, reveal, and
  `All`/gate specifics; only the zone-to-zone move matrix is shared.
- **Pick any number of cards one at a time.** "Destroy any number of ...",
  "purge any number of ...", "keep N ..." use
  `pickCards(ctx, prompt, limit, optional, avail)` (`target_select.go`):
  `limit <= 0` is unbounded, `optional` makes each prompt declinable, and
  `avail` is re-read each round so a shrinking pool stays current (Obsidian
  Forge, Destructive Analysis, Unnatural Selection, Tertiate). Do not re-roll
  the pick-until-decline loop.
- **Redistribute a resource among a side's creatures.** Equalize (Æmber) and
  Entropic Manipulator (damage) drain a resource into a pool and hand it back
  one unit at a time onto a chooser-picked creature (falling back to the
  first): `placeAmong(ctx, creatures, prompt, pool, place)`
  (`effect_redistribute.go`). Keep the two **separate nodes**
  (`RedistributeCapturedAember`, `RedistributeDamage`): the drain, the side
  selection, the accept gate, and the lethal follow-up all diverge, and one
  `Resource`-switched node would be a branchy blob. Share only the loop.

**Prefer a shared helper over a shared embeddable value type.**
`{Amount, By, Per}` is uniform across the Æmber effects, but a
`Quantity{Amount, By, Per}` does not pay off: Go composite literals do not
promote embedded fields, and the `card` facade aliases these structs, so every
flat authoring site (`{Amount: 2, Per: X}`) would become a nested
`{Quantity: Quantity{…}}` across ~120 cards and their tests. Factor the logic,
keep the fields.
