# 35. House choice is a delayed constraint table; must and cannot stack, cannot overrides must, and wagers are reactions to the choice

## Context

Three cards reach forward to a player's next house choice, and each grew its own
paired state slot, resolver method, and "next turn" promotion:

- **Must** — Control the Weak (a chosen house) and Snag (the house of the creature
  it fought) make the opponent choose a particular house next turn. Modeled by
  `ForcedHouse`/`ForcedHouseNext` (`Bar[House]`), `ForceActiveHouseNextTurn`, and
  the `OpponentMustChooseHouse` effect.
- **Cannot** — Tezmal (a chosen house) and Snag's Mirror (the house a player just
  chose) bar a player from a house next turn. Modeled by
  `ForbiddenHouse`/`ForbiddenHouseNext`, `ForbidActiveHouseNextTurn`, and
  `OpponentCannotChooseHouse`.
- **Wager** — Snaglet bets the opponent will pick a chosen house next turn and
  steals if they do. Modeled by `HouseWager`/`HouseWagerNext`,
  `WagerOnHouseNextTurn`, `payOffHouseWager`, and
  `WagerOpponentChoosesChosenHouse`.

`StartTurn` promotes each `…Next` slot to its current-turn slot, and `ChooseHouse`
enforces them: it rejects a choice that isn't the forced house (unless the player
doesn't have that house, in which case cannot-overrides-must lets any house
through), rejects the forbidden house, and pays off the wager once a house locks
in. The "after you choose a house" trigger window (`TriggerAfterChooseHouse`)
fires for every card in play — the chooser's own abilities and every other card's
each-player-scoped ones — which is how Snag's Mirror (an `EachPlayer` ability)
reads the house a player just chose.

Four problems pushed this to a decision:

1. **Each slot holds only one constraint.** `ForcedHouse` is a single `Bar[House]`,
   so a second "must" overwrites the first, and there is no way to express _must A,
   must B, cannot A → must B_. The single-slot model cannot stack.
2. **Snag is modeled statically but plays dynamically.** `OpponentMustChooseHouse`
   with a `Fought` source reads the fought creature's house **at fight time** and
   freezes it. A creature's houses can change before the opponent's next turn
   (house-granting and house-changing effects), so the constraint must read the
   creature's _current_ houses at **choice time**, not fight time.
3. **The wager is not a constraint at all.** It neither forces nor forbids a house;
   it observes the choice and pays out. It is a reaction to the "house chosen"
   event — exactly the shape Sci. Officer Qincan already uses through the
   `TriggerAfterChooseHouse` window — wearing a bespoke promotion-slot costume.
4. **The choosable set is wrong, and the empty case is unmodelled.**
   `playerHasHouse` checks only the three houses on the Archon identity card, but a
   player may also choose the house of any card they currently control (rulebook
   lines 454, 2355): a controlled off-house creature or artifact adds its house to
   the choice. That also means a _must_ for a house the player neither has on their
   identity card nor controls is void — a maverick Pitlord forcing Dis in a
   deck without Dis. And nothing modelled the case where every choosable house is
   forbidden: the player then has _no active house_ (rulebook line 788), a valid
   outcome the single slots could not express.

## Decision

**A player's forward house constraints are one delayed constraint table**, not
three paired slots. Each entry records who imposed it, which player it binds, and
what it says — a _must_ (this house, or a live reference resolved at choice time)
or a _cannot_ (this house). Entries accumulate; nothing overwrites. `StartTurn`
promotes the next-turn entries for the player whose turn is beginning and clears
them, the same lifecycle the single slots had.

**`ChooseHouse` resolves the whole table at choice time.** It collects the player's
active entries, resolves any live _must_ references (Snag reads the fought
creature's houses _now_), and computes the allowed set:

- Start from the player's **choosable houses**: every house on their Archon
  identity card, plus every house of a card they currently control. A controlled
  off-house card — a creature or artifact taken from the opponent — contributes
  its house, so a Mars deck controlling an enemy Dis creature may choose Dis
  (rulebook line 454). Only a card in play in its own right counts — a creature in
  the battleline or an artifact in the artifact row; an attached upgrade or a card
  under a host does not contribute a house. This set is read at choice time: a
  controlled card that has since left play no longer contributes its house.
- Remove every _cannot_ house. **Cannot overrides must**: a house that is both
  required and barred is barred.
- The surviving _musts_ are the required houses that are still choosable and not
  forbidden. If any survive, the allowed set is exactly those houses and the player
  picks any one of them — _must A, must B_ leaves _{A, B}_ (either is legal);
  _must A, must B, cannot A_ leaves _{B}_.
- A _must_ naming a house that is forbidden (cannot overrides must) **or no longer
  choosable** does not survive. A must for a house the player neither has on their
  identity card nor controls is void — a maverick Pitlord forcing Dis in a deck
  without Dis, or a Mark of Dis "must choose D" after the controlled house-D
  creature has left play, so the player falls back to their identity houses
  (rulebook line 2355). When no must survives, every not-forbidden choosable house
  is allowed.
- If the allowed set is empty — every choosable house is forbidden — the player
  has **no active house** this turn (rulebook line 788). This is a valid resolved
  outcome, not an error. The player makes the choice explicitly: the client presents
  a **No House** option they must click, so the state is unmistakable rather than
  silently auto-resolved. Choosing No House _is_ their house choice for the phase —
  it consumes an armed wager (Snaglet, which then pays nothing) and opens the
  after-choose-house window, where house-keyed reactions (Snag's Mirror, Qincan)
  find no chosen house and no-op. The player then enters their main phase, may play
  or use the cards a card ability specifically permits (an out-of-house play
  allowance such as Witch of the Wilds), and ends the turn normally.

A house is rejected only when it is not choosable, is forbidden, or a surviving
must names a different house. This subsumes today's two guards and the
`playerHasHouse` escape hatch — which today checks only identity houses — into one
set computation over the choosable set.

**A must that reads a creature is dynamic.** Snag stores a reference to the fought
creature, not a frozen house. At choice time the constraint reads that creature's
current houses; a multi-house creature contributes each of its houses to the must
set, and a creature no longer in play (or whose houses changed) resolves against
its state _now_. Control the Weak and Tezmal, which name a house the card already
chose, store that fixed house.

**The wager is a reaction, not a constraint.** Snaglet is authored as a reaction to
the house-chosen event (the `TriggerAfterChooseHouse` window), reading the house
the opponent just chose and stealing if it matches the predicted house — the same
mechanism as Qincan, not a promotion slot. Its one extra requirement is memory: the
predicted house must survive until the opponent's **next** turn, and a player may
take several turns before that (extra-turn effects), so the wager's armed state
rides the same next-turn promotion the constraints use and is consumed the first
time that player chooses a house — including a No House choice, which resolves the
bet as lost.

### Why one table instead of keeping the three slots

The three slots are the same shape — a house constraint that arms this turn, waits,
and fires at the next choice — spelled three times with three promotion pairs and
three resolver methods. Stacking is impossible without a list, dynamic musts are
impossible without a reference, and the wager was never a house constraint. One
table with typed entries expresses all three, makes stacking and
cannot-overrides-must a single set computation the choice already half-implements,
and deletes two of the three resolver method families.

### Why not a public constraint type for authors

Card authors keep the printed KeyForge verbs — `OpponentMustChooseHouse`,
`OpponentCannotChooseHouse`, and a reaction for the wager. Following ADR 0006 and
ADR 0031, the table is the private mechanism; the authoring vocabulary stays the
words the card prints. The effects build table entries instead of calling three
different arm-next-turn resolver methods.

## Consequences

- _Must A, must B, cannot A → must B_ and every other stack becomes expressible;
  the single-slot overwrite bug is gone.
- Snag respects house changes between the fight and the opponent's choice, because
  the must resolves the creature's houses at choice time.
- Snaglet stops being a bespoke promotion slot and becomes a reaction to the
  house-chosen event, so a future "when a house is chosen" card reuses that window
  instead of a fourth slot.
- `ChooseHouse` gains one allowed-set computation and loses its three ad-hoc guards
  and the `payOffHouseWager` call; the constraint logic lives in one place that is
  exercised by the stacking cases directly.
- The choosable set widens from the identity card's three houses to identity houses
  plus the houses of cards the player controls, read at choice time: controlling an
  enemy off-house creature makes its house legal to choose, and losing that creature
  makes it illegal again — the same live read that voids a must the creature was
  granting. This widened read is a **new** `StateReader.AllowedHouses(player)` port
  method (the choosable set minus cannots, then surviving musts), which `ChooseHouse`
  and the client's house picker both consult. `PlayerHasHouse` deliberately stays
  **identity-only** and does not widen: it answers "is this house on the player's
  identity card", which is what `ItIsOffIdentity` (Sneklifter's off-identity check)
  means by house — a controlled off-house card must not make a creature count as
  in-identity. Widening the choice and answering the identity question are two
  different reads, so they are two methods.
- "No active house" becomes a first-class resolved outcome: when cannots bar every
  choosable house, the turn proceeds with no active house rather than blocking the
  choice, and the player relies on out-of-house permissions for the turn.
- Realized as a single pass, not incrementally: the state slots, resolver methods,
  and the three effects change together, and the card sites (Control the Weak,
  Tezmal, Snag and its Mirror, Snaglet) move to the table and the reaction in the
  same change. Coverage stays at 100% because every allowed-set branch —
  surviving must, fully-cancelled must, plain cannot, dynamic multi-house must — is
  reachable from the existing cards plus the stacking case a test constructs.
  </content>
  </invoke>
