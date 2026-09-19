# Future-set backlog

Work that is real and decided, but that has **no consuming card in an implemented
set yet** — so it waits for the set that first introduces a card that needs it.
This is the sibling of [todo-agent.md](todo-agent.md): that file is work to do
now; this file is work parked against a set that has not been stood up.

Rules:

- **When you begin implementing (or stubbing) a set, scan this file first.** If an
  item names that set — or a card that set introduces — fold it into the run and
  build the primitive alongside its first real consumer, so a live card pins its
  shape and a card test covers it. The `implement-cards` and `stub-cards` skills
  both point here for exactly this reason.
- Each item carries the **decision, not just the task** — the chosen behavior, the
  cards affected, and the design already settled — so a future agent who was not in
  the conversation can build it without re-deriving it.
- **When the work lands, delete the item** (do not mark it done), same as
  [todo-agent.md](todo-agent.md). This file only ever shows parked work.

## Result-set reference — "choose one of the cards this effect just moved"

**Trigger set:** whichever set first stands up a consuming card — Junk Restoration
(Æmber Skies / Grim Reminders / Menagerie) or Haunting Measures (Draconian
Measures). Build it with that card, not before.

**What it is.** A reference to the exact cards a preceding sub-effect just moved —
KeyForge's "a card discarded **this way**" — where the player **chooses one** of
those cards. The producer half already exists: `DiscardTop` /
`DiscardTopOfEachDeck` record what they discard into `ctx.Produced.Discarded` (a
`[]LocalID` result set) in `internal/engine/effect_deck.go`. Two consumers of that
set already exist — `ForEachDiscarded` (act on **all** of them) and
`DiscardDeckUntil` + `PutDiscardedIntoHand` (find **the first match** by digging).
The missing consumer is **"the player picks one from the produced set."**

Why the result set is its own thing and not "the discard pile": it can be a strict
subset of the pile (discard 3 onto an already-full pile → only those 3 are
eligible), and it can be empty even when the pile is not.

**Cards (both unimplemented, in unbuilt sets):**

- **Junk Restoration** — "Discard the top 3 cards of your deck. You may put a card
  discarded this way into your hand." Discard 3, then optionally choose one of
  those 3 to keep.
- **Haunting Measures** — "Discard the top 6 cards of your deck. You may put a
  non-Geistoid card discarded this way into your hand." Same, with a house
  exclusion (non-Geistoid) and 6 cards.

**Design decided:**

- Build a **reusable "choose one from the produced result set" selection** — it
  picks one card from `ctx.Produced.Discarded` into `ctx.It`, then reuses the
  existing tail (`PutDiscardedIntoHand`, or a future purge-/archive-from-discard).
  Prefer this over a single fused `PutOneDiscardedIntoHand` node so hand / purge /
  archive all share one reference. This matches the "reference primitive" framing
  and covers the purge/archive-top cards below for free.
- The selection is **optional** ("you may") and carries a **filter**: a house
  include/exclude (Haunting Measures excludes Geistoid; Junk Restoration has no
  filter) plus a card-type refinement (for the purge/archive-top cards). One small
  filter on the chooser covers all three shapes.
- The **"purge/archive the top card"** scrap-adjacent cards are the same family
  with the tail swapped (purge or archive instead of put-into-hand) and no choice
  when only one card was produced — the `Amount: 1` case. Build them on the same
  reference when their set arrives.

This is deliberately distinct from the top-of-deck **zone slice**
(`LookAtTopOfDeck`, already built): a zone slice reads deck **positions**; a result
set reads **what an effect produced**.

## Live "considered an artifact/creature" type grant (ADR 0033 live route)

**Trigger set:** whichever set first stands up a card that converts a type on a
lifetime **other than the converted card's own** — Deanimator is the model. Build
it with that card, not before.

**What it is.** The **live route** of ADR 0033's two type-conversion routes. The
stored route (`LastingType` on `CardCore`) is built and covers permanent
self-conversions (Auto-Legionary; Effigy of Melerukh and The Mysticeti are its
still-unimplemented siblings, needing no new primitive). The **live** route is not
built: a card is a given type only while some other condition holds — a counter on
it, a source in play, a while-condition — and the change must **revert for free**
when that condition ends.

**Card (unimplemented, in an unbuilt set):**

- **Deanimator** — "Each card that has a mineralize counter on it is considered an
  artifact." The grant ends when Deanimator leaves play or the counter is removed,
  not when the converted card leaves — so it cannot be stored.

**Design decided (per ADR 0033):**

- Author it as a `ConstantAbility` composed over a `CounterInPlay{Kind, Target:
This}` read (per ADR 0024), computed **live at read time** the way `Power`/`Armor`
  fold `constantBonus` — **never** a `LastingType` write. It reverts for free: the
  next `TypeOf` read simply no longer sees the grant.
- Grow `TypeOf` one step **after** the `HostPlus`/`LastingType` checks: ask whether
  any active constant ability converts this card's type, exactly as the stat reads
  scan for bonuses. Do not store the result. Two Deanimators compose the same way,
  each grant recomputed every read.
- Only if a **stored** and a **live** override could ever land on one card at once
  does the single `LastingType` field become insufficient — then fold it into a
  flat comparable `TypeEntry{Card, Type, Source}` LIFO side-table (like the control
  stack, ADR 0028) and record the last-applied-wins precedence. No such card exists
  today; do not build the stack pre-emptively.

## Curse of Forgery — "when you would forge a key" self-veto

**Trigger set:** Crucible Clash (CC #214) / Vault Masters 2026 (VM26 #171) —
whichever stands up first introduces this card. Build it then.

**What it is.** A Treachery upgrade that vetoes its own controller's next forge:

> Treachery.
> When you would forge a key, purge Curse of Forgery and do not forge that
> key (no Æmber is spent).

**It reuses the Keyforgery machinery, mirrored to the forger's own side.**
Keyforgery (Worlds Collide #271) already built the veto: a before-forge trigger
window (`beforeForgePrevented` in `internal/engine/forge_guard.go`), the transient
`State.ForgePrevented` flag, the `CancelForge{}` effect + `CancelCurrentForge()`
resolver, and the `CancelForge` icon case. Keyforgery's trigger
(`TriggerBeforeOpponentForgesKey`) gathers the window over the **opponent's**
in-play cards. Curse of Forgery needs the **own-forge** variant:

- Add a new `Trigger.BeforeYouForgeKey` (engine `TriggerBeforeForgeKey`) with its
  own rulebook term, gathered by `beforeForgePrevented` over the **forger's own**
  in-play cards (add a second `allInPlay(forger)` pass alongside the existing
  opponent pass — the flag and window resolution are shared).
- The ability is `PurgeCreature{Target: Target.This}` (purge on an artifact
  self-target) then `CancelForge{}`, composed the same decomposed way Keyforgery is.
- The "(no Æmber is spent)" clarifier is dropped from card text per the Keyforgery
  divergence already recorded in
  [keyforge-divergences.md](keyforge-divergences.md) — the veto simply leaves the
  Æmber unspent.

## House wager — generalize over polarity and payoff

**Trigger set:** whichever set first stands up **Allusions of Grandeur**
(Discovery #255 / Vault Masters 2024 #296 / 2025 #264 / 2026 #227 — all
Unfathomable, all unimplemented). Build the generalization with that card, and
refactor Snaglet onto it in the same change, so a second live consumer pins the
shape.

**What it is.** A **wager on the opponent's next house choice**: arm a bet now,
settle it when the opponent locks in their active house on their next turn (or
does not — the no-house outcome, rulebook 788, is a valid settlement). The
scheduling half already exists and ships with Snaglet: the `constraintWager` row
in [house_constraint.go](../internal/engine/house_constraint.go), armed by
`WagerOnHouseNextTurn`, promoted onto the opponent's turn by `StartTurn`, and
settled at choice time by `resolveHouseWagers` in
[game_turn.go](../internal/engine/game_turn.go). It rides ADR 0035's delayed
house-constraint table as a **reaction**, not a constraint — it does not restrict
the choice, it pays out on it.

**The two cards are the same machinery, mirrored on two axes:**

|              | Snaglet (WC #118, built)     | Allusions of Grandeur (unbuilt)       |
| ------------ | ---------------------------- | ------------------------------------- |
| Arm          | choose any house             | choose a house on opp's identity card |
| Schedule     | opponent's next house choice | opponent's next house choice          |
| House read   | Chosen                       | Chosen                                |
| **Polarity** | pays **on match**            | pays **on mismatch**                  |
| **Payoff**   | **steal** 2                  | **gain** 3                            |

Allusions' printed text: "Play: Choose a house on your opponent's identity card.
If your opponent does not choose that house as their active house on their next
turn, gain 3 Æmber." Snaglet's: "Action: Choose a house — if your opponent chooses
that house as their active house on their next turn, steal 2 Æmber."

**Design decided:**

- The shared axis is **polarity + payoff, not a `HouseChoiceReference`.** Both
  cards read the **Chosen** house — neither wagers on a Fought or JustChosen house
  — so the wager still must **not** grow a `HouseChoiceReference` field (that stays
  the rejected, speculative wrapper). What the second card earns is:
  - a **polarity** flag on the wager row — pays-on-match | pays-on-mismatch — read
    by `resolveHouseWagers` (Snaglet matches its predicted house; Allusions settles
    when the choice is anything **other** than its house, including no-house);
  - an **enum-tagged payoff** replacing the hardcoded `StealAember{Amount}` in
    `resolveHouseWagers` — steal (Snaglet) vs gain (Allusions), with amount. Model
    it the way the lasting-reaction registry keeps effects out of flat state
    (ADR 0005, ADR 0007): a small enum-tagged payoff, **not** an `Effect` value
    stored on the row.
- This is the state-resident form of the naive `Then: Conditional{Cond:
houseMatches, Then: <payoff>}` body. The body genuinely parameterizes over polarity
  and payoff — Allusions proves it — but the **schedule and the prediction memory**
  (predicted house, amount, predictor) must live in the flat comparable constraint
  table, so it cannot be an inline synchronous `Then`.
- The **arming choice** differs but is a separate concern: Allusions restricts the
  choice to the **opponent's identity houses** (a refinement on `ChooseHouseThen`),
  where Snaglet chooses any house. That refinement is the choosing half, not the
  wager row — build it on the choose side.
- When this lands, rename `WagerOpponentChoosesChosenHouse` if the match-only name
  no longer fits the polarity-carrying node, and update its rulebook/icon coverage.

## DiscardUntil — grow the from-hand-at-random source, Player, and hand-size stop

**Trigger set:** whichever set first stands up a consuming card — either
**High Street Churn** (Ekwidon, Æmber Skies #104) or **Catch and Release**
(Unfathomable, Winds of Exchange #369 and reprints). Build each axis alongside
its first real card.

**What it is.** `DiscardUntil` (`internal/engine/effect_deck.go`) is the shared
"discard until a condition, or the source runs out" node. Today it discards only
from the **top of the controller's own deck** and stops on a **type/house match**
(Sound the Horns, Invasion Portal, Old Boomy). The two future cards discard from
**hand, at random**, one for the **opponent** and one for **each player**, and
Catch and Release stops on a **hand-size** threshold rather than a card match. Grow
the node's axes as each card lands — do **not** build any of these before its card
exists (they would be uncovered).

**Cards (both unimplemented, in unbuilt sets):**

- **High Street Churn** — "Play: Choose a house on your opponent's identity card.
  Your opponent discards a random card from their hand until they discard a card of
  the chosen house or run out of cards. They refill their hand as if it were their
  'draw cards' step."
- **Catch and Release** — "Play: Return each creature to its owner's hand. Each
  player discards random cards from their hand until they have 6 or fewer cards in
  hand. Gain 2 chains."

**Design decided:**

- **Keep it one cohesive `DiscardUntil` node — do not build a generic
  `Until{Condition}` combinator.** There is no other kind of "until" in the game, so
  a general combinator is speculative; the axes below stay local to this node.
- **`From` axis** — add a source strategy: the current top-of-deck source vs
  **hand-at-random**. Its zero value stays top-of-deck so the three built cards are
  untouched. High Street Churn and Catch and Release both use hand-at-random.
- **`Player` axis** — add an explicit player (default the controller). High Street
  Churn targets **the opponent**; Catch and Release targets **each player**. Render
  "your opponent discards…" / "each player discards…" from the field, the way the
  other player-bearing effects do.
- **Terminator axis** — the current type/house match stays; add a **hand-size**
  stop (`until they have N or fewer cards in hand`) for Catch and Release. Keep the
  house filter **local to this node** (it already is) — do **not** add a `House`
  axis to `CardFilter`; that overlap was the old blocker and is rejected.
- **The house choice, the hand refill, the creature return, and the chains are
  separate composed effects**, not part of `DiscardUntil`: High Street Churn is
  `MustChooseHouse` (restricted to the opponent's identity houses) → `DiscardUntil`
  → `RefillHand`; Catch and Release is `ReturnEachCreatureToHand` → per-player
  `DiscardUntil{hand-size}` → `GainChains`. Compose them with `Sentences`.

## InExcessOf — a floored "in excess of" Count combinator

**Trigger set:** whichever set first stands up **Change Agent** (Æmber Skies #101 /
Dark Tidings #202, both reprints, all unimplemented). Build it with that card, not
before.

**What it is.** A generalization of `ExcessCreatures`
([internal/engine/effect_count_board.go](../internal/engine/effect_count_board.go))
into a reusable `Count` that subtracts one measure from another and **floors the
result at 0** — KeyForge's "in excess of" idiom. `ExcessCreatures` is the only
arithmetic `Count` today (one side's creature count minus the other's, clamped at
0); every other `Count` reads state directly and returns a non-negative int. The
grilling that parked this rejected a fully-general signed `Difference{Count,
Minus}` for three reasons, all still binding:

- **The text does not compose.** KeyForge never prints "X minus Y" — it prints "in
  excess of". Mechanically joining two child `CountText()`s reads nothing like a
  real card, so any combinator still needs a hand-written `CountText()`; the
  generic form buys nothing on the text side.
- **Signedness is a live hazard.** `scaled(base, per)` in
  [effect_count.go](../internal/engine/effect_count.go) is an unfloored multiply, so
  a `Count` that can go negative silently inverts gain/damage/draw amounts.
  `ExcessCreatures` is safe only because it self-clamps at 0. Any real combinator
  **must** floor at 0 — at which point it is the "in excess of" shape, not arbitrary
  arithmetic.
- **No card needs two independently-varying dynamic Counts subtracted.** The one
  waiting consumer is count-vs-**fixed-number**, not count-vs-count.

**Card (unimplemented, in an unbuilt set):**

- **Change Agent** — "For each card in your opponent's hand in excess of 5, they
  lose 1 Æmber." That is `CardsInHand{Player: Opponent}` in excess of the **fixed
  constant 5**, not two dynamic counts subtracted.

**Design decided:**

- Introduce `InExcessOf{Count Count, Of Count}` — `Value` = `max(0, Count.Value −
Of.Value)`, floored at 0 so it can never feed a negative into `scaled`. The `Of`
  side may be a `Fixed` (Change Agent's "in excess of 5") or another `Count`
  (creature-vs-creature). It carries a bespoke `CountText()` ("in excess of" idiom),
  not a mechanical join of the two children.
- **Re-express `ExcessCreatures` as a thin constructor over `InExcessOf`** in the
  same change, so a second live consumer pins the shape: its two operands are the
  two sides' creature counts (an `InPlay`-style creature count per player, with the
  trait filter applied identically to both sides). Thread `NotCountingSelf` through
  the operand that sits on the source's side (Dr. Milli's self-exclusion) —
  `InPlay.Other` already expresses "not counting the source".
- Keep the `CardFilter` fold `ExcessCreatures` gained in the InPlay/ExcessCreatures
  refactor: the trait axis routes through `filter().admits`, applied to both sides.

## Winds of Exchange — zone visibility as a computed fact

**Trigger set:** **Winds of Exchange**, where house Ekwidon makes looking into a
hidden zone a house mechanic rather than a one-off. Three cards force it:

- **Flea Market** (WoE 064, Ekwidon) — "Look at **a random card** in your
  opponent's hand. You may give your opponent 1 Æmber. If you do, play that card
  as if it were yours." The look is granted for **one randomly chosen card**, not
  the hand.
- **Talent Scout** (WoE 069, Ekwidon) — "Look at your opponent's hand and play a
  creature from it as if it were yours." A whole-hand grant, then a `Chosen`
  filtered to creatures.
- **Abyssal Sight** (WoE 384, Unfathomable) — "look at your opponent's hand and
  choose a card from it. That player discards that card." Grant, `Chosen`, then a
  pile verb — the shape this whole item exists for.

**Not Æmber Skies' Clipped Wings** (AS 071), which reads "Purge **a random** card
from your opponent's hand" — that is plain `Random{}` and needs no visibility at
all. The AS card that would need it is Talent Scout again (AS 097, the reprint).

**The grant is partial in this set, which today's seam cannot express.** The only
grant the engine has is `RevealHand`, which opens the whole hand. Flea Market
grants a look at exactly one card, and WoE's Plunder (356) reveals "a random
**unrevealed** card", which needs per-card revealed state within a zone. Design
the predicate to take the granted scope, not a hand-wide boolean.

**What it is.** A predicate answering whether a given viewer can see the
individual cards in a given player's zone. It is a function of all three of zone,
owner, and viewer — a discard pile, the board, and the purge pile are open to
both players; a hand and archives are open only to their owner; a deck is closed
to both, including its owner. `Zone.public()` in `internal/engine/zone.go` is the
two-players-only corner of it and should be derived from it once it lands, not
left beside it.

**What the chooser does with it.** A `Chosen` aimed at a zone the chooser cannot
see picks uniformly at random instead, so the invalid combination is
unrepresentable rather than merely rejected.

**The grant is an input, not an exception.** A card can hand a player a look they
would not otherwise have — today a preceding `RevealHand` (Hidden Stash, Imperial
Traitor), in WoE a single card or a peek at the top N of a deck. The grant must be
an argument to the predicate, so a granted look and a natural one answer through
the same seam.

**Two traps, both already paid for once — read before designing:**

- **`Text()` renders with no game state, and the grant is a runtime fact.** A
  reveal is a sibling effect earlier in the same `Sequence`, so static visibility
  would print "at random" for Imperial Traitor while the resolve let the
  controller choose — exactly the ADR 0006 desync the feature is supposed to
  prevent. The design has to say what `Text()` renders for a grant it cannot see,
  and "derive the wording from visibility" is not an answer on its own.
- **Do NOT fold `ownerActs` into it.** `ownerActsSelection` / `Random.ownerActs`
  (`internal/engine/effect_selection.go`) looks like the same idea and is not: it
  asks whether the **pick** is the controller's to make, not whether the zone is
  visible. Imperial Traitor chooses from a hidden hand and still reads in the
  controller's imperative voice. The two coincide today only because every blind
  pick happens to target a hidden zone. Pinned by
  `TestBlindPickVoiceIsUniform`.
