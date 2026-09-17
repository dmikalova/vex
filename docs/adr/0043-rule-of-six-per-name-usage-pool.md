# The Rule of Six is one per-name usage pool, counted and enforced at every use

## Context

KeyForge's **Rule of Six** caps how many times a player may use a single card
_name_ in one turn at six. It is a global anti-combo rule: without it, a pair of
cards that use each other (two Replicators reaching for each other's reap effect,
a Bait and Switch that repeats itself, two Legatus Raptors readying and using each
other) would loop without bound. The rule is not per-card — it is per **name**,
summed across every copy of that name in play, so eight Replicators or three
copies of an Automaton all draw from the same six. It caps the **active player**,
who is the only one who uses cards during a turn: if that player controls an
opponent's copy of a name they also own, both copies draw from the one six, not
two separate pools.

"Use" here is broad. The rulebook counts a name's usage every time a copy is
played, discarded, used (reaped, fought, or its `Action:` fired), resolves a
`Destroyed:` ability, or repeats an ability of its own. Each of these draws from
the one pool of six. The rule must hold no matter _how_ the usage is initiated —
whether the active player takes the action directly, or a card's effect forces it
(Replicator triggering another creature's reap, Legatus Raptor using a friendly
creature, a self-repeating effect looping). This last point is the whole
difficulty: the count is easy, but a mechanic that drives usage through a
back-channel (an ability that reaps, fights, triggers, or repeats) can silently
bypass a cap that only the player's own top-level actions check.

Two forces pull on the design. The count wants to live in flat, comparable state
(ADR 0005) so undo is a snapshot and no closures leak in. The enforcement wants to
sit at a small number of chokepoints, not be scattered across every effect that
could conceivably cause a usage — a scattered check is one a new mechanic forgets,
exactly the class of bug ADR 0029 closed for destruction.

## Decision

### One pool, keyed by name alone, held in flat state

The pool is tracked per card in `GameState.UsagesThisTurn` — a fixed
`[maxCards]uint8` indexed by `LocalID`, a flat comparable array with no pointers
(ADR 0005). It is **not** the running total; it is the raw per-copy tally.
`nameUsagesThisTurn(id)` computes the pool on demand by summing
`UsagesThisTurn[i]` across every card `i` whose name matches `id`, regardless of
owner. Keying the sum by name alone is what makes every copy share one pool,
without any per-name registry to keep in state. It is safe to ignore owner because
the ledger resets every turn (see `StartTurn` below) and only the active player
uses cards within a turn — so a usage recorded this turn is always the active
player's, and their control of an opponent's same-named copy correctly draws from
the same six rather than opening a second pool.

`atRuleOfSix(id)` reports `nameUsagesThisTurn(id) >= RuleOfSix` (`RuleOfSix = 6`).
`recordUsage(id)` counts one usage by incrementing `UsagesThisTurn[id]`. The whole
ledger resets to zero in `StartTurn` — every turn is a fresh six-usage window for
every name.

`UsagesThisTurn` lives on `GameState` beside `Cards`, **not** inside `CardCore`,
because the cap is per-name (summed across copies) rather than per-card, and
because `resetCore` would clear a copy's tally when it left play — the pool must
persist across a copy leaving so the name's remaining copies still see the spend.

### Counting points — every kind of usage records one

Each distinct usage records against the pool at the moment it happens:

- **Playing** a card — `recordUsage` in the play path (`game_play.go`).
- **Discarding** a card — `recordUsage` in the discard path.
- **Using** a creature or artifact (reap, fight, `Action:`) — `recordUse`, which
  calls `recordUsage` and additionally bumps the creature's own
  `TimesUsedThisTurn` / turn tallies. Reap, fight, and action all funnel through
  it (`reapWith`, `game_combat.go`'s `fight`, `useActionOf`).
- **A `Destroyed:` resolution** — `recordUsage` in `game_leaves_play.go`, and the
  matching gathered-trigger site.
- **A self-repeating ability** past its free first loop — `RecordUsage(ctx.Source)`
  inside the `While`/`MayWhile` loop bodies (`effect_repeat.go`).
- **A chained Replicator-style trigger** past the free first — `RecordUsage(root)`
  in `TriggerAbility.Resolve` (`effect_trigger.go`).

### The free first resolution, and cascade-root attribution

A genuine use already counts one (reaping records a usage). That single use
**buys the first resolution of its own repeat/trigger ability for free** — the
head rides on the use and is not double-counted. Every resolution _after_ the
first costs one more:

- A **repeat** loop resolves its first iteration outside the counting loop (free),
  then records one per further iteration, so Bait and Switch's "steal 1 Æmber →
  repeat" resolves the initial steal plus at most five repeats.
- A **chained trigger** charges not the card whose ability happens to resolve
  mid-chain, but the **cascade root** — the card that was genuinely used to start
  the chain. `TriggerAbility.Resolve` threads the root through `EffectContext`
  (`Root`/`HasRoot`; absence marks the free head) and charges `root`'s name pool
  for each trigger past the head. So two Replicators reaching for each other spend
  **one** pool of six between them (six total resolutions), not six each; a
  Replicator and a differently-named Doppelganger charge only the Replicator's
  pool, never twelve.

### Enforcement points — checked at both the player and ability boundaries

The cap is enforced at two kinds of boundary, never scattered into the effects
themselves:

1. **The active player's own top-level actions.** `usable` (the gate behind
   reap/fight/action/unstun) and the play path both return `ErrRuleOfSix` when
   `atRuleOfSix` holds, so a seventh play or use of a name is refused before it
   resolves. A `Destroyed:` resolution and a repeat/trigger loop each consult
   `atRuleOfSix` at their own boundary and stop.

2. **Ability-driven use.** The ability-driven use API — `ReapWith`, `FightWith`,
   `UseActionOf` on the resolver — is the back-channel a card uses to make another
   card reap, fight, or fire its action (Replicator, Legatus Raptor, Sergeant
   Zakiel). These share one gate, `usableByAbility`, which no-ops the use when the
   target is exhausted **or** at the Rule of Six. This is why two Legatus Raptors
   cannot ready-and-use each other past the six between them: the gate stops the
   seventh use even though no top-level action was taken. The gate mirrors the
   free-first-plus-cap shape of the repeat and trigger paths — the player's own
   uses already checked the cap in `usable`, so gating the ability API bounds the
   forced uses without touching the top-level path.

The reason the gate lives at these boundaries rather than inside each effect is
the same reason destruction settles at boundaries (ADR 0029): a new mechanic that
drives a usage reaches for one of these existing entry points and is bounded for
free, whereas a per-effect check is one a new mechanic forgets.

## Consequences

- A card-name combo cannot loop past six by construction, whichever way the usage
  is driven — player action, forced reap/fight/action, `Destroyed:` cascade,
  repeat loop, or Replicator trigger chain — because every path either records
  against the shared pool and checks it at its boundary, or routes through the one
  ability-driven use gate that does.
- The pool stays flat and comparable: it is a per-copy `uint8` array summed on
  demand, so undo remains a snapshot and no name registry has to be kept in sync.
- New code inherits the cap by using the existing seams. A mechanic that forces a
  use calls `ReapWith`/`FightWith`/`UseActionOf` (bounded by `usableByAbility`); a
  mechanic that repeats or triggers records past its free first and consults
  `atRuleOfSix`. Adding a bespoke `if atRuleOfSix` inside a new effect, or driving
  a usage around these seams, is the smell this decision exists to prevent.
- The free-first-resolution rule keeps a genuine use from being double-charged: a
  single Replicator reap costs one, not two, and the cascade-root attribution
  keeps a two-name chain honest — one shared pool, not one per name.
  </content>
  </invoke>
