# Card implementation reference

One place to answer "can the engine already do this?" while implementing a card.
It is a **capability catalog**: what exists, what it is called, and how to use it.

It is not the authoring guide. For file layout, generated comments, the
one-field-per-line style, and tests, read
[../internal/cards/AGENTS.md](../internal/cards/AGENTS.md). For printed wording,
read [card-wording-rules.md](card-wording-rules.md). For why the engine is shaped
this way, read [../internal/engine/AGENTS.md](../internal/engine/AGENTS.md).

## How to look a capability up

The failure mode this page exists to prevent is **building something the engine
already has**. Before deciding a card is gated on new work, in this order:

1. **Scan this page.** The tables below list every effect node, target, filter,
   condition, count, trigger, and card-level option by name.
2. **Read the facade.** `internal/card/` is the whole authoring surface, and every
   entry carries a doc comment with an example card:
   [effects.go](../internal/card/effects.go) (effects, conditions, counts),
   [target.go](../internal/card/target.go) (targets and refinements),
   [options.go](../internal/card/options.go) (card-level options),
   [types.go](../internal/card/types.go) (houses, types, traits, keywords,
   triggers, players), plus `duration.go`, `destination.go`, `zone.go`.
3. **Grep for the mechanic word, not the card.** Grepping `internal/card/` for a
   word like "neighbor" finds more than reading one card does.
4. **Find a card that already does it.** `mage tool:lookup "<name>"` gives the
   printed text; the implementation is `internal/cards/sets/<slug>/<snake>.go`.

Only then is the card gated. When it is, extend an existing node (a field, a
`Strategy`, a target filter, a count, a condition) before adding a new one — a
brand-new effect node is grill-gated, see the `implement-cards` skill.

## The shape of a card

```go
var AmmoniaClouds = set.New(
  "Ammonia Clouds",
  card.House.Mars,
  card.Type.Action,
  card.Rarity.Common,
  card.Provenance(card.CotA, "199"),
  card.WithAbility(
    card.Trigger.Play, card.DealDamage{
      Amount: 3,
      Target: card.Target.EachCreature,
    }),
)
```

Four positional arguments (name, house, type, rarity), then options. An ability is
a **trigger** plus an **effect**; an effect is a tree of nodes, each of which
renders its own printed text and resolves itself.

The vocabulary an effect composes from:

| Axis      | What it answers                | Section                                            |
| --------- | ------------------------------ | -------------------------------------------------- |
| Trigger   | when it fires                  | [Triggers](#triggers)                              |
| Target    | which cards in play it touches | [Targets](#targets)                                |
| Effect    | what happens                   | [Effects](#effects)                                |
| Condition | whether it happens             | [Conditions](#conditions)                          |
| Count     | how much / how many            | [Counts](#counts)                                  |
| Duration  | how long it lasts              | [Durations](#durations-zones-destinations-players) |
| Selection | which cards in a zone it moves | [Zone movement](#zone-movement)                    |

## Triggers

`card.Trigger.X`, passed as the first argument to `card.WithAbility`. Com.
Officer Kirby fires the same ability on three timings:

```go
card.WithAbility(card.Trigger.PlayFightReap, card.PlayFrom{
  From:  card.Hand,
  House: card.Houses.Except(card.House.Self),
  Types: card.Types.Of(card.Type.Artifact, card.Type.Upgrade, card.Type.Tactic),
}),
```

**Own-card timing:**

| Trigger                  | Fires                                       |
| ------------------------ | ------------------------------------------- |
| `Play`                   | after this card is played                   |
| `Reap`                   | after this creature reaps                   |
| `Fight`                  | after this creature fights                  |
| `BeforeFight`            | before this creature fights                 |
| `Action`                 | when this card's action is used             |
| `Destroyed`              | when this card is destroyed                 |
| `LeavesPlay`             | when this card leaves play any way          |
| `UsedSelf`               | after this creature is used at all          |
| `AfterDestroyedFighting` | when the creature it fought is destroyed    |
| `AfterAssaultDestroys`   | when its assault damage destroys the target |
| `AfterArmorPrevents`     | after its armor prevents damage             |

**Turn structure:**

| Trigger                      | Fires                                |
| ---------------------------- | ------------------------------------ |
| `StartOfTurn`                | at its controller's turn start       |
| `EndOfTurn`                  | at its controller's turn end         |
| `EndOfReadyStep`             | after the ready step                 |
| `AfterChooseHouse`           | after its controller chooses a house |
| `AfterAnyPlayerChoosesHouse` | after either player chooses a house  |
| `AfterAnyPlayerStartOfTurn`  | at either player's turn start        |
| `AfterAnyPlayerEndOfTurn`    | at either player's turn end          |

**Reactions to the board:**

| Trigger                          | Fires                                     |
| -------------------------------- | ----------------------------------------- |
| `AfterCardPlayed`                | after its controller plays a card         |
| `AfterEnemyCardPlayed`           | after the opponent plays a card           |
| `AfterCreaturePlayed`            | after a creature is played                |
| `AfterCreatureEnters`            | after a creature enters play any way      |
| `AfterCreaturePlayedAdjacent`    | after a creature enters next to this one  |
| `AfterUpgradeEnters`             | after an upgrade enters play              |
| `AfterTacticPlayedBeforeResolve` | before a played tactic resolves           |
| `AfterUse`                       | after its controller uses a card          |
| `AfterCreatureReaps`             | after a creature reaps                    |
| `AfterCreatureFights`            | after a creature fights                   |
| `AfterNeighborFights`            | after a neighbor fights                   |
| `AfterCreatureDestroyed`         | after a creature is destroyed             |
| `AfterEnemyDestroyedFighting`    | after an enemy dies fighting              |
| `AfterDiscardFromHand`           | after its controller discards from hand   |
| `AfterBonusDamage`               | after a damage bonus icon resolves        |
| `AfterBonusDraw`                 | after a draw bonus icon resolves          |
| `AfterAemberStolenFromYou`       | after Æmber is stolen from its controller |

**Keys:**

| Trigger                   | Fires                                      |
| ------------------------- | ------------------------------------------ |
| `AfterForgeKey`           | after its controller forges                |
| `AfterPlayerForgesKey`    | after either player forges                 |
| `AfterOpponentForgesKey`  | after the opponent forges                  |
| `BeforeOpponentForgesKey` | before the opponent forges — can cancel it |

**Composites** fan out into the atomic triggers and print as one line.

| Trigger         | Fires on             |
| --------------- | -------------------- |
| `PlayFightReap` | play, fight, or reap |
| `FightReap`     | fight or reap        |
| `PlayReap`      | play or reap         |
| `PlayFight`     | play or fight        |

For a _granted_ pair use `card.FightReap(effect)`, which returns the
`[]card.Ability` a `Granted` list wants.

There is **no Omni trigger**. `Omni:` is authored as
`card.WithKeywords(card.Keyword.Versatile)` plus a `Trigger.Action` ability
(ADR 0009).

A card in play fires only its printed triggers; `card.WithTriggersFromDiscard()`
keeps them live from the discard pile (Relentless Creeper).

## Targets

`card.Target.X` names a set of cards in play, from the source's point of view.
Tribute captures onto the most powerful friendly creature:

```go
card.CaptureAember{
  Amount: 2,
  Target: card.Target.EachFriendlyCreature.Refine(card.MostPowerful),
  Source: card.Opponent,
}
```

**Chosen** — one card, the controller picks.

| Target                       | Selects                             |
| ---------------------------- | ----------------------------------- |
| `Creature`                   | any creature                        |
| `FriendlyCreature`           | a friendly creature                 |
| `EnemyCreature`              | an enemy creature                   |
| `OtherFriendlyCreature`      | a friendly creature other than this |
| `OtherCreature`              | a creature other than this          |
| `Artifact`                   | any artifact                        |
| `FriendlyArtifact`           | a friendly artifact                 |
| `EnemyArtifact`              | an enemy artifact                   |
| `Upgrade`                    | any upgrade                         |
| `CreatureOrArtifact`         | either, the controller's choice     |
| `FriendlyCreatureOrArtifact` | either, friendly                    |
| `EnemyCreatureOrArtifact`    | either, enemy                       |

**Whole sets:**

| Target                      | Selects                               |
| --------------------------- | ------------------------------------- |
| `EachCreature`              | every creature in play                |
| `EachFriendlyCreature`      | every friendly creature               |
| `EachEnemyCreature`         | every enemy creature                  |
| `EachOtherFriendlyCreature` | every friendly creature except this   |
| `EachArtifact`              | every artifact in play                |
| `EachFriendlyArtifact`      | every friendly artifact               |
| `EachEnemyArtifact`         | every enemy artifact                  |
| `EachFriendlyCardInPlay`    | every friendly creature and artifact  |
| `EachNeighbor`              | this card's neighbors                 |
| `EachUpgradeOnThis`         | the upgrades on this card             |
| `FormerNeighbors`           | the neighbors a departed creature had |

**Contextual** — a card an earlier step or the trigger bound.

| Target              | Selects                                  |
| ------------------- | ---------------------------------------- |
| `This`              | the source card itself                   |
| `Triggering`        | the card that fired the trigger          |
| `TheSameCreature`   | the creature the previous step used      |
| `TheOtherCreature`  | the other creature in a two-card step    |
| `TheChosenCreature` | the creature an earlier choice picked    |
| `CreatureFought`    | the creature being fought right now      |
| `TheFoughtCreature` | the creature this card fought            |
| `AttachedHost`      | the creature this upgrade is attached to |
| `GrantingCard`      | the card that granted this ability       |

### Target filters

Chain off any target; they conjoin.

| Filter                                                      | Keeps                                  |
| ----------------------------------------------------------- | -------------------------------------- |
| `.House(m)`                                                 | cards a `card.Houses` matcher admits   |
| `.WithTrait(t)` / `.ExceptTrait(t)`                         | by trait                               |
| `.Named(n)`                                                 | by printed name                        |
| `.SharingTrait()`                                           | creatures sharing a trait with another |
| `.PowerAtMost(n)` / `.PowerAtLeast(n)` / `.PowerExactly(n)` | by power                               |
| `.OddPower()` / `.EvenPower()`                              | by power parity                        |
| `.Damaged()` / `.Undamaged()`                               | by damage                              |
| `.WithAember()` / `.WithoutAember()`                        | by Æmber on the card                   |
| `.WithArmor()` / `.WithUpgrade()` / `.WithCounter(k)`       | by what it carries                     |
| `.WithoutBonusIcons()`                                      | cards with no bonus icons              |
| `.Keyword(k)`                                               | creatures with a keyword               |
| `.Stunned()` / `.Ready()`                                   | by state                               |
| `.OnFlank()` / `.NotOnFlank()` / `.InCenter()`              | by battleline position                 |
| `.ToLeftOfSource()` / `.ToRightOfSource()`                  | by side of the source                  |
| `.Neighboring()` / `.AndNeighbors()` / `.NeighborsOf()`     | by adjacency                           |
| `.SharesHouseWithNeighbors(n)`                              | by neighboring houses                  |
| `.OfHouseWithMostCreatures()`                               | the most-represented house             |
| `.Other()`                                                  | excludes the source                    |

### Refinements

`.Refine(r)` narrows relative to the whole selected set — so "each enemy creature
except the most powerful" is
`card.Target.EachEnemyCreature.Refine(card.Except(card.MostPowerful))`.

**Power:**

| Refinement              | Keeps                                   |
| ----------------------- | --------------------------------------- |
| `MostPowerful`          | exactly one, the controller breaks ties |
| `LeastPowerful`         | exactly one, the controller breaks ties |
| `HighestPower`          | every creature tied at the top          |
| `LowestPower`           | every creature tied at the bottom       |
| `MostPowerfulN(n)`      | the top n creatures                     |
| `PowerLessThan(count)`  | creatures below a running count         |
| `PowerLessThanSource()` | creatures weaker than the source        |

**Shape:**

| Refinement                 | Keeps                                    |
| -------------------------- | ---------------------------------------- |
| `KeepPerSide(n)`           | all but n creatures per side             |
| `PortionPerSide(fraction)` | a fraction of each side                  |
| `HouseWithAtLeast(n)`      | only houses with at least n creatures    |
| `WithoutSharedTrait()`     | creatures sharing no trait with another  |
| `SamePowerAsChosen`        | creatures matching the chosen power      |
| `SamePowerAsEitherChosen`  | creatures matching either chosen's power |

**Combinators:**

| Refinement    | Keeps                        |
| ------------- | ---------------------------- |
| `Except(r)`   | everything `r` would discard |
| `AnyOf(r...)` | the union of the refinements |

> For any "most/least" comparison, ties **pass**. Where an effect must pick one
> tied candidate, the active player chooses.

### House matchers

A `Target`'s `.House(m)` filter and any effect field named `House` take a
`card.Houses.X` matcher.

| Matcher      | Admits                                     |
| ------------ | ------------------------------------------ |
| `Named(h)`   | one named house                            |
| `Except(h)`  | every house but one                        |
| `Chosen`     | the house an earlier step chose            |
| `Active`     | the active house this turn                 |
| `Each`       | the house the enclosing `ForEachHouse` set |
| `Contextual` | the house of the card in context           |
| `Any`        | every house                                |

### A card that names its own house writes `card.House.Self`

Most cards that name a house name their own: Battle Fleet (Mars) reveals Mars
cards, Pitlord (Dis) locks you into Dis, Witch of the Wilds (Untamed) lets you
play an Untamed card off-house. Spelling that house out a second time lets the
two drift and does not work for Mavericks, so the ability names `card.House.Self`
and `card.New` fills the card's own house in when the definition is built:

```go
card.WithPlayPermission(card.PlayPermission{House: card.House.Self, Amount: 1}),
```

The sentinel never survives past `card.New`, so the printed text, resolution, and
state all see the concrete house — the generated comment still reads "one Untamed
card". The same holds for a `Target` filtered by house: Ixxyxli Fixfinger (Mars)
buffing each other Martian creature writes
`card.Target.EachOtherFriendlyCreature.House(card.Houses.Named(card.House.Self))`.

The test is what the card is _about_, not which house it happens to print. Take
That, Smarty Pants names Logos because it is about Logos creatures, whichever
house the card itself belongs to — that house is written out. So:

- **Named house == the card's own house?** Use `card.House.Self`.
- **Named house is a different house** (or the card would still say "Logos" if it
  were reprinted in another house)? Write the house out.

`TestNoCardHardcodesItsOwnHouse` enforces this by parsing every `card.New` call:
because `card.New` resolves the sentinel to the concrete house, the two are
indistinguishable afterward, so the test reads source.

## Effects

Every entry below is `card.X`, documented with an example card in
[effects.go](../internal/card/effects.go).

### Æmber

Special Agent Fingers steals one Æmber:

```go
card.WithAbility(
  card.Trigger.Action, card.StealAember{Amount: 1}),
```

| Effect                       | What it does                                |
| ---------------------------- | ------------------------------------------- |
| `GainAember`                 | supply → a player's pool                    |
| `LoseAember`                 | a player's pool → supply                    |
| `StealAember`                | opponent's pool → yours                     |
| `GiveAember`                 | opponent hands you a toll, or all of it     |
| `CaptureAember`              | a pool → onto a capturing creature          |
| `CaptureFromAnyPlayer`       | captures from both pools in any split       |
| `DistributeCapture`          | captures a pool one Æmber at a time         |
| `Exalt`                      | supply → onto a chosen card                 |
| `PlaceAemberOnThis`          | supply → onto this card                     |
| `MoveAember`                 | off a card into a pool or onto another card |
| `MoveAemberFromPool`         | banks your pool Æmber onto a card           |
| `MoveAemberToSupply`         | removes Æmber sitting on a card             |
| `RedistributeCapturedAember` | reshuffles captured Æmber between creatures |

How much: set exactly one of these on the effect.

| Field     | Meaning                                    |
| --------- | ------------------------------------------ |
| `Amount`  | a fixed number                             |
| `Per`     | `Amount` once per unit of a `Count`        |
| `EqualTo` | "equal to X", straight from a `Count`      |
| `By`      | a `LoseAember` portion instead of a number |

`By` takes `card.HalfRoundedDown`, `card.AllBut(5)`, or `card.AllAember`.

### Damage and combat

A Vinda deals one damage and, if it kills, follows up:

```go
card.WithAbility(
  card.Trigger.Reap, card.DealDamage{
    Amount: 1,
    After:  card.IfDestroyed,
    Target: card.Target.Creature,
    Then:   card.DiscardCard{Player: card.Opponent, Zones: []card.Zone{card.Hand}, Selection: card.Random{Count: 1}},
  }),
```

| Effect                     | What it does                                     |
| -------------------------- | ------------------------------------------------ |
| `DealDamage`               | deals damage to each creature the target selects |
| `DealDamage`               | deals damage, then resolves a follow-up          |
| `Heal`                     | removes damage                                   |
| `LoseArmor`                | strips remaining armor                           |
| `GainStats`                | grants power or armor for a duration             |
| `GainAssault`              | grants assault this turn                         |
| `GainAssaultUntilNextTurn` | grants assault through the next turn             |
| `RedistributeDamage`       | moves existing damage between creatures          |
| `RedirectFightDamage`      | sends fight damage somewhere else                |
| `CancelFight`              | stops the fight before damage                    |
| `TakesExtraDamage`         | adds to damage a creature will take              |
| `CannotBeDealtDamage`      | makes a creature untouchable by damage           |
| `OverrideStats`            | replaces a creature's printed stats              |

`DealDamage` fields:

| Field        | Meaning                             |
| ------------ | ----------------------------------- |
| `Amount`     | base damage                         |
| `Per`        | `Amount` once per unit of a `Count` |
| `AmountFrom` | damage read straight off a `Count`  |
| `PerTarget`  | damage scaled per target hit        |
| `Spread`     | how the damage fans out             |

`DealDamage{After: …}` gates the follow-up on `Always`, `IfDestroyed`, or
`IfSurvives`.

Spreads:

| Spread                 | What it does                                                 |
| ---------------------- | ------------------------------------------------------------ |
| `CreatureAndNeighbors` | a creature plus its neighbors (`Scope: OneNeighbor` for one) |
| `DifferentCreatures`   | splits across distinct creatures                             |
| `UpToCreatures`        | up to N creatures, the controller picks                      |
| `DivideDamage`         | divides a pool of damage as chosen                           |
| `FlankWalk`            | walks inward from a flank, decreasing                        |

Per-target scaling (`PerTarget`):

| Value              | Scales by               |
| ------------------ | ----------------------- |
| `AemberOnIt`       | Æmber on that creature  |
| `DamageOnIt`       | damage on that creature |
| `ArmorLostThisWay` | armor it just lost      |
| `UpgradesOnIt`     | upgrades on it          |

### Destruction and purging

Igon the Green purges itself, then fetches its counterpart:

```go
card.WithAbility(
  card.Trigger.Destroyed, card.Sequence{Effects: []card.Effect{
    card.PurgeCreature{Target: card.Target.This},
    card.PutFromDiscard{Selection: card.Chosen{Name: IgonTheTerrible.Name}, Destination: card.To.Hand},
  }}),
```

| Effect                           | What it does                              |
| -------------------------------- | ----------------------------------------- |
| `Destroy`                        | destroys everything the target selects    |
| `DestroyChosen`                  | destroys a creature chosen mid-resolution |
| `BatchDestroy`                   | destroys a `Gather`ed set at once         |
| `DestroyEachCreatureAtEndOfTurn` | schedules a board wipe for end of turn    |
| `PurgeCreature`                  | purges a creature from play               |
| `PurgeCard`                      | purges a card from a pile                 |
| `PurgeFromHand`                  | purges a card out of hand                 |
| `PurgeSource`                    | purges the source card                    |
| `PurgeArchives`                  | purges a player's archives                |
| `PurgeArchivedCardThen`          | purges from archives, then resolves more  |

`BatchDestroy` takes a `Gather` such as `EachPlayerUnless`, which spares the
creatures a nested count matches.

`Sacrifice <self>` is authored as `Destroy` — there is one destruction verb.

#### An effect can read a creature it just destroyed

When one effect destroys (or deals lethal damage to) a creature and a following
effect in the **same** ability reads that creature — "Destroy a friendly creature.
Each player loses Æmber equal to half **its power**" (Power of Fire) — the read
sees the creature's power, Æmber-on-card, and damage as they were the instant
before it left play, counters and buffs included, not the zeroed card it becomes.
This is automatic for `PowerOfChosen`, `AemberOnThis` / `AemberOnIt`, and
`DamageOnThis` / `DamageOnIt`; compose them freely after a `Destroy` or a
`DealDamage{After: card.IfDestroyed}`.

Only those three mutable dimensions are captured — printed power, house, traits,
keywords, and bonus icons survive on their own; a departed creature's _granted_
keywords do not (no card reads them). See
[ADR 0030](adr/0030-a-card-out-of-play-takes-no-further-part.md).

### Creature state

Stunner grants its host a stun on fight or reap:

```go
card.WithStatic(card.StaticModifier{
  Granted: card.FightReap(card.May{Do: card.Stun{Target: card.Target.Creature}}),
}),
```

| Effect             | What it does                             |
| ------------------ | ---------------------------------------- |
| `Stun`             | stuns the selected creatures             |
| `Unstun`           | removes stun                             |
| `Enrage`           | enrages, forcing it to fight             |
| `Ward`             | wards, absorbing the next harm           |
| `RemoveWard`       | strips a ward                            |
| `Exhaust`          | exhausts the selected creatures          |
| `ExhaustCreatures` | exhausts creatures chosen mid-resolution |
| `Ready`            | readies the selected creatures           |
| `ReadyCreatures`   | readies creatures chosen mid-resolution  |
| `ReadyIfFirstUse`  | readies only on the first use this turn  |
| `AddPowerCounter`  | adds permanent power counters            |
| `PlaceCounter`     | places a counter of a named kind         |
| `RemoveCounters`   | removes counters of a named kind         |

Counter kinds: `card.Counter.Doom`, `.Fuse`, `.Growth`, `.Glory`, `.Disruption`,
`.Scheme`.

### Keywords, traits, text boxes

Mutation of Fury folds two next-turn grants into one clause:

```go
card.GainUntilNextTurn{Effects: []card.Effect{
  card.GainAssaultUntilNextTurn{Target: card.Target.TheChosenCreature, Amount: card.Fixed(3)},
  card.GainTrait{Target: card.Target.TheChosenCreature, Trait: card.Traits.Mutant},
}}
```

| Effect                | What it does                                   |
| --------------------- | ---------------------------------------------- |
| `GainKeywords`        | grants keywords for a duration                 |
| `LoseKeywords`        | removes several keywords                       |
| `LoseKeyword`         | removes one keyword                            |
| `GainTrait`           | grants a trait                                 |
| `GainUntilNextTurn`   | folds several next-turn grants into one clause |
| `GainAbility`         | grants a triggered ability                     |
| `GainTextBox`         | copies another card's text box                 |
| `LendTextBoxFromHand` | lends a text box off a card in hand            |
| `BlankEnemyText`      | blanks an enemy card's text                    |
| `ConsiderFlank`       | treats a creature as if on a flank             |

### Zone movement

Archiving, discarding, purging, shuffling into the deck, and putting a card into
play / hand / on top of a deck are **one mechanism** (ADR 0031): a **source
zone**, a **Selection**, and a **destination**. Author each with the KeyForge verb
the card prints.

When a card needs a movement an existing verb does not cover, it is almost always
a new source zone, selection, or destination on that family — **not** a new
bespoke type. These families are being consolidated onto the one mechanism as
their cards are touched, so shape a new movement effect to fit it rather than
adding another one-off `…FromHand` / `…TopOfDeck` variant to unwind later.

**Selections** — which cards in the source zone move.

| Selection | Takes                              |
| --------- | ---------------------------------- |
| `Chosen`  | one the controller picks           |
| `Random`  | one at random                      |
| `Each`    | every card that matches            |
| `Named`   | a card named by name               |
| `Self`    | the source card                    |
| `Top`     | the top card(s) of an ordered pile |
| `Bottom`  | the bottom card(s)                 |

`Chosen` and `Each` narrow with `House`, `Type`, `Trait`, `Name`, and `Or`;
`Chosen` also takes `Optional`.

**Zones** (the source): `card.Hand`, `card.Deck`, `card.Discard`,
`card.Archives`.

**Destinations:**

| Destination            | Sends it             |
| ---------------------- | -------------------- |
| `card.To.Hand`         | into hand            |
| `card.To.TopOfDeck`    | on top of the deck   |
| `card.To.BottomOfDeck` | under the deck       |
| `card.To.DeckShuffled` | into a shuffled deck |
| `card.To.Archives`     | into archives        |

Deck-top routing steps use the parallel `card.Into.Hand`, `.Archives`,
`.Discard`, `.Purge`, `.BottomOfDeck`.

**Verbs:**

| Effect                            | What it does                             |
| --------------------------------- | ---------------------------------------- |
| `ArchiveCard`                     | archives from a pile or hand             |
| `ArchiveFromPlay`                 | archives a card in play                  |
| `ArchiveSource`                   | archives the source card                 |
| `ArchiveGrantingUpgrade`          | archives the upgrade that granted this   |
| `ArchivePurgedCard`               | archives a card out of the purge pile    |
| `ArchiveDiscardedThisWay`         | archives what an earlier step discarded  |
| `ArchiveCardUnder`                | archives the card sitting under this one |
| `DiscardCard`                     | discards from hand, deck, or archives    |
| `DiscardHand`                     | discards a whole hand                    |
| `DiscardArchives`                 | discards a whole archives                |
| `DiscardTop`                      | discards off the top of a deck           |
| `DiscardUntil`                    | discards until a match turns up          |
| `PutFromPlay`                     | moves a card out of play                 |
| `PutChosen`                       | moves a card chosen mid-resolution       |
| `PutFromDiscard`                  | moves a card out of the discard pile     |
| `PutFromHand`                     | moves a card out of hand                 |
| `PutIntoPlay`                     | puts a card into play                    |
| `PutDiscardedIntoHand`            | returns what an earlier step discarded   |
| `PutDiscardedIntoPlay`            | puts what was discarded into play        |
| `ReturnNamedToHand`               | returns a card named by name             |
| `ReturnItToHand`                  | returns the card in context              |
| `Shuffle`                         | shuffles a deck                          |
| `ShuffleFromDiscard`              | shuffles discard cards back in           |
| `ShuffleFriendlyCardsIntoDeck`    | shuffles friendly cards in play back in  |
| `ShuffleChosenCreaturesFromZones` | shuffles chosen creatures back in        |
| `ShuffleNamedFromDiscardIntoDeck` | shuffles a named discarded card back in  |
| `SwapDeckAndDiscard`              | swaps the two piles                      |
| `Search`                          | searches named zones for a card          |
| `Draw`                            | draws cards                              |
| `RefillHand`                      | refills to the hand size                 |

`Search` names its zones explicitly (`Sources`) and never shuffles on its own:

```go
card.Search{
  Sources: []card.Zone{card.Deck, card.Discard},
  Filter:  card.Filter{Trait: card.Traits.Beast},
  Reveal:  true,
}
```

### Reading and playing off the top

`LookAtTopOfDeck` and `RevealTopOfDeck` read N cards and route them through
ordered `Then` steps. With no steps they are a pure peek/reveal — Navigator Ali
looks at three and puts them back in any order:

```go
card.WithAbility(card.Trigger.PlayFightReap, card.LookAtTopOfDeck{
  Amount: 3,
  Then:   []card.TopAct{card.ReorderRest{}},
}),
```

Routing steps (`[]card.TopAct`):

| Step                     | What it does                                |
| ------------------------ | ------------------------------------------- |
| `ChooseAndMove`          | moves `Count` chosen cards to `Dest`        |
| `PartitionByChosenHouse` | splits the read cards by a chosen house     |
| `ReorderRest`            | puts the rest back in any order — goes last |
| `MayDiscardLookedAt`     | optionally discards what was looked at      |

Playing:

| Effect                             | What it does                             |
| ---------------------------------- | ---------------------------------------- |
| `PlayTopOfDeck`                    | plays the top card of a deck             |
| `PlayRevealedCard`                 | plays a card an earlier step revealed    |
| `PutRevealedCard`                  | moves a revealed card instead of playing |
| `PlayFrom`                         | plays a card out of a named zone         |
| `PlayOrUse`                        | plays or uses one matching card          |
| `PlayFromOpponent`                 | plays a card from the opponent's zones   |
| `PlayItFromOpponentDiscard`        | plays the card in context from discard   |
| `DiscardOpponentArchivesOrDeckTop` | discards one of the opponent's cards     |

Revealing hands: `RevealHand`, `RevealChosenFromHand`, `RevealRandomFromHand`.

Bonus icons: `ResolveBonusIcons` resolves another card's icons;
`ExtraBonusIconResolution` makes future icons resolve twice.

### Under and graft

Spangler Box grafts a creature under itself, then hands itself to the opponent:

```go
card.WithAbility(
  card.Trigger.Action, card.Sequence{Effects: []card.Effect{
    card.Graft{Target: card.Target.Creature},
    card.TakeControl{Target: card.Target.This, ToOpponent: true, Duration: card.Duration.Forever},
  }}),
```

| Effect                     | What it does                         |
| -------------------------- | ------------------------------------ |
| `PutUnderFromHand`         | puts a card from hand under this one |
| `PlayCardUnder`            | plays a card that is under this one  |
| `PutUnderIntoPlay`         | puts the cards underneath into play  |
| `ArchiveCardUnder`         | archives the card underneath         |
| `Graft`                    | moves a card in play under this one  |
| `TriggerGraftedPlayEffect` | fires a grafted card's Play ability  |

See ADR 0016 for what "under" is.

### Using and moving creatures

Mega Ganger Chieftain readies a neighbor and makes it fight:

```go
card.WithAbility(
  card.Trigger.Play, card.OnChooseCreature{
    Target: card.Target.Creature.Neighboring(),
    Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
  }),
```

| Effect                 | What it does                                |
| ---------------------- | ------------------------------------------- |
| `OnChooseCreature`     | picks a creature and applies `Verbs` to it  |
| `ChooseCreatureThen`   | picks a creature, then resolves an effect   |
| `OneAtATime`           | repeats over a different creature each time |
| `RepeatedFight`        | fights again with the same creature         |
| `Use`                  | uses a card as if its controller did        |
| `TriggerAbility`       | fires another card's ability                |
| `TakeControl`          | moves a card to the other player's side     |
| `AttachSelfTo`         | attaches this upgrade to a creature         |
| `Swap`                 | swaps this creature with another            |
| `SwapChosen`           | swaps two chosen creatures                  |
| `RearrangeBattleline`  | reorders a whole battleline                 |
| `MoveToFlank`          | moves a creature to a flank                 |
| `MoveWithinBattleline` | moves a creature one step or more           |
| `TurnIntoCreature`     | makes a non-creature card a creature        |

Verbs for `OnChooseCreature`: `ReadyVerb`, `ReapVerb`, `FightVerb`, `UseVerb`,
`StunVerb`, `ExhaustVerb`, `GainKeywordVerb`.

### Houses, keys, chains, restrictions

United Action grants out-of-house play, then bans using cards this turn:

```go
card.Sentences{Effects: []card.Effect{
  card.MayPlayOrUse{Houses: card.GrantHouses.Controlled, Grant: card.GrantPlay},
  card.Restrict{Player: card.Controller, Action: card.Restricted.Use, Duration: card.Duration.RemainderOfPlayerTurn},
}}
```

| Effect                            | What it does                           |
| --------------------------------- | -------------------------------------- |
| `Restrict`                        | bars `Fighting`, `Reaping`, or `Use`   |
| `CannotPlay`                      | bars a player from playing a card type |
| `PlayersCannotPlay`               | bars both players                      |
| `CreaturesCannot`                 | bars creatures from an action          |
| `MayPlayOrUse`                    | grants out-of-house play or use        |
| `BelongToHouse`                   | changes which house a card belongs to  |
| `ChangeActiveHouse`               | switches the active house              |
| `NameHouse`                       | names a house for a later step         |
| `NameCard`                        | names a card for a later step          |
| `MustChooseHouse`                 | forces the next house choice           |
| `CannotChooseHouse`               | bars a house from being chosen         |
| `OpponentNamesHouse`              | the opponent names the house           |
| `WagerOpponentChoosesChosenHouse` | a wager on the opponent's guess        |
| `ForgeKey`                        | forges a key now                       |
| `UnforgeKey`                      | takes a forged key back                |
| `CancelForge`                     | stops a forge that is about to happen  |
| `RaiseKeyCost`                    | raises the next key's cost             |
| `RaiseKeyCostPerHouseCreature`    | raises it per creature of a house      |
| `LowerKeyCost`                    | lowers the next key's cost             |
| `SkipForgePhase`                  | skips the forge step entirely          |
| `GainChains`                      | gives a player chains                  |
| `EndTurn`                         | ends the turn immediately              |

`MayPlayOrUse` is the single out-of-house permission node (ADR 0037):

| Field    | Takes                                                 |
| -------- | ----------------------------------------------------- |
| `Houses` | `card.GrantHouses.Named/Chosen/Any/Except/Controlled` |
| `Grant`  | `card.GrantPlay`, `card.GrantUse`, `card.GrantFight`  |
| `Types`  | `card.Types.Of(...)`, to narrow the card types        |
| `Trait`  | a trait the card must have                            |
| `Count`  | how many cards the permission covers                  |

### Control flow

Psionic Officer Lang gates its reaction on who reaped:

```go
card.WithAbility(
  card.Trigger.AfterCreatureReaps, card.Conditional{
    Cond: card.ItIsEnemy{},
    Then: card.ArchiveCard{Zone: card.Deck, Selection: card.Top{}},
  }),
```

| Node                                 | Use                                                               |
| ------------------------------------ | ----------------------------------------------------------------- |
| `Sequence`                           | several effects in order, joined with ", and"                     |
| `Sentences`                          | several effects rendered as separate sentences on one line        |
| `May`                                | the whole effect is optional                                      |
| `Then{First, Result}`                | the `A -> B` result gate — `Result` only if `First` did something |
| `Conditional{Cond, Then, Otherwise}` | gated on a condition                                              |
| `ChooseOne`                          | the controller picks one of several effects                       |
| `ChooseHouseThen`                    | pick a house, then resolve                                        |
| `ForEach`                            | once per running count, choosing afresh each time                 |
| `ForEachHouse`                       | once per house, binding that house                                |
| `ForEachDiscarded`                   | once per card a preceding discard removed                         |
| `Repeat{Do, Gate}`                   | `While`, `MayWhile`, or `ByExalting`                              |
| `ForDuration`                        | several timed effects sharing one duration clause                 |
| `ByActivePlayer`                     | resolve as the active player, not the controller                  |
| `OrAmount`                           | "N, or M if <cond>" without an `Otherwise` branch                 |

Keep a `May` only when declining has a real downside; pure-upside "may"s are
authored as mandatory.

## Conditions

Gate a `Conditional`, a `Repeat`, `CountIs`, or a card-level option — as in the
Psionic Officer Lang example above.

**Pools and keys:**

| Condition             | Asks                                |
| --------------------- | ----------------------------------- |
| `PoolAember`          | how much Æmber a player has         |
| `ForgedKey`           | whether a key was forged            |
| `HasMoreForgedKeys`   | whether a player leads on keys      |
| `KeyColorForged`      | whether a given key color is forged |
| `AemberStolenFromYou` | whether Æmber was stolen from you   |
| `HasAember`           | whether Æmber sits on its `Subject` |

`PoolAember{Player, Is, Amount}` compares with `AtLeast`, `AtMost`, `Exactly`,
`MoreThanYou`, `MoreThanOpponent`, `Even`, or `Odd`.

**Board:**

| Condition                         | Asks                                    |
| --------------------------------- | --------------------------------------- |
| `ControlsMoreCreatures`           | who has more creatures                  |
| `Overwhelmed`                     | whether the opponent has more           |
| `ControlsNamed`                   | whether a named card is in play         |
| `ControlsCreaturesOfHouses`       | whether creatures of N houses are out   |
| `PlayerControlsFewerHousesThan`   | who spans fewer houses                  |
| `SourceInCenterOfBattleline`      | whether the source is in the center     |
| `SourceHasNoNeighbor`             | whether the source stands alone         |
| `SourceReady`                     | whether the source is ready             |
| `SourceIsFighting`                | whether the source is in a fight now    |
| `OnFlank`                         | whether a card is on a flank            |
| `ItIsAmong`                       | whether the card in context is in a set |
| `CounterInPlay`                   | whether a counter kind is out           |
| `CountersOnThisAtLeast`           | how many counters this card carries     |
| `ActiveHouseMatchesNoCardsInPlay` | whether the active house is absent      |

**Piles:**

| Condition               | Asks                                  |
| ----------------------- | ------------------------------------- |
| `Haunted`               | whether a discard pile has 10+ cards  |
| `CardsInDiscardAtLeast` | how many matching cards are discarded |
| `NamedCardInDiscard`    | whether a named card is discarded     |
| `NamedCardPurged`       | whether a named card is purged        |

**The card in context ("it"):**

| Condition                    | Asks                                  |
| ---------------------------- | ------------------------------------- |
| `ItIs`                       | its house and/or type                 |
| `ItIsNamed`                  | its printed name                      |
| `ItIsOfTrait`                | whether it has a trait                |
| `ItIsFriendly`               | whether you control it                |
| `ItIsEnemy`                  | whether the opponent controls it      |
| `ItIsStunned`                | whether it is stunned                 |
| `ItIsOffIdentity`            | whether its house is outside yours    |
| `ItHasBonusIcon`             | whether it prints a bonus icon        |
| `ItIsNotOfNamedHouse`        | whether it dodges a named house       |
| `ItAttachedToThisOrNeighbor` | whether it sits on this or a neighbor |

**This turn:**

| Condition                     | Asks                                   |
| ----------------------------- | -------------------------------------- |
| `ItIsYourTurn`                | whether you are the active player      |
| `ChoseHouse`                  | which house was chosen                 |
| `FirstCreaturePlayedThisTurn` | whether this is the first creature     |
| `NoCreaturesPlayedThisTurn`   | whether no creature was played         |
| `UsedCreatureToReap`          | whether a creature reaped              |
| `UsedCreatureToFight`         | whether a creature fought              |
| `UsedNoCreatures`             | whether no creature was used           |
| `FirstReapOfTurn`             | whether this is the first reap         |
| `SourceFirstUseThisTurn`      | whether this is the source's first use |
| `CardsDiscarded`              | whether cards were discarded           |
| `EnemyCreatureDestroyed`      | whether an enemy creature died         |
| `FriendlyCreatureDestroyed`   | whether a friendly creature died       |

**Within this resolution:**

| Condition                     | Asks                                   |
| ----------------------------- | -------------------------------------- |
| `CardsDestroyedFewerThan`     | how many this resolution has destroyed |
| `ArchivedCreaturesShareHouse` | whether the archived creatures match   |
| `MovedAnyAember`              | whether an earlier step moved Æmber    |

**Tide:** `TideIsLow`, `TideIsHigh`.

**Composition:**

| Condition   | Asks                                  |
| ----------- | ------------------------------------- |
| `And`       | every listed condition holds          |
| `Or`        | at least one holds                    |
| `Not`       | the inner condition fails             |
| `AlwaysMet` | nothing — it always passes            |
| `CountIs`   | how a `Count` compares to an `Amount` |

`CountIs{Count, Is, Amount}` turns any Count below into a condition, so a
missing condition is usually a Count plus `CountIs`.

## Counts

Feed an effect's `Per`, `EqualTo`, `AmountFrom`, or `Times`, and a `CountIs`
condition. Phalanx Strike deals one damage per friendly creature:

```go
card.DealDamage{
  Amount: 1,
  Per:    card.InPlay{Player: card.Controller, Type: card.Type.Creature},
  Target: card.Target.Creature,
}
```

**Board:**

| Count                             | Counts                               |
| --------------------------------- | ------------------------------------ |
| `Fixed`                           | a constant                           |
| `InPlay`                          | matching cards in play               |
| `ArtifactsInPlay`                 | artifacts in play                    |
| `ExcessCreatures`                 | creatures beyond the opponent's      |
| `HousesInPlay`                    | distinct houses in play              |
| `HousesAmong`                     | distinct houses among selected cards |
| `HousesRepresented`               | houses represented in a zone         |
| `NeighborsOfThis`                 | this card's neighbors                |
| `NeighborsMatching`               | neighbors that match a filter        |
| `CombinedPowerOfNeighborsWithout` | neighbor power, excluding a trait    |
| `UpgradesOn`                      | upgrades on a target                 |

**Cards and zones:**

| Count             | Counts                      |
| ----------------- | --------------------------- |
| `CardsInHand`     | matching cards in a hand    |
| `CardsInZone`     | cards in a named zone       |
| `CopiesInDiscard` | copies of a card in discard |
| `PurgedCards`     | cards in the purge pile     |

**Æmber and keys:**

| Count                       | Counts                        |
| --------------------------- | ----------------------------- |
| `AemberInPool`              | Æmber in a player's pool      |
| `AemberOnThis`              | Æmber on this card            |
| `AemberOnFriendlyCreatures` | Æmber on your creatures       |
| `AemberStolenThisEvent`     | Æmber the current theft moved |
| `ForgedKeys`                | keys already forged           |
| `UnforgedKeys`              | keys still to forge           |

**This turn:** `CardsPlayed`, `CreaturesUsed`, and `TurnCount` (a named
`card.TurnStat` tally).

**Counters:** `CountersOnThis`, `PowerCountersOnThis`.

**The card in context:**

| Count                | Counts                                         |
| -------------------- | ---------------------------------------------- |
| `PowerOfChosen`      | the chosen creature's power (`Of:` a fraction) |
| `DamageOnThis`       | damage on this card                            |
| `TraitsOfChosen`     | the chosen card's traits                       |
| `BonusIconsOfChosen` | the chosen card's bonus icons                  |

**"This way" tallies** — what an earlier step in the same resolution produced.

| Count                   | Counts                            |
| ----------------------- | --------------------------------- |
| `CardsDestroyed`        | cards destroyed this way          |
| `CreaturesDestroyed`    | creatures destroyed this way      |
| `PowerDestroyedThisWay` | total power destroyed this way    |
| `CardsPurged`           | cards purged this way             |
| `PurgedAemberBonus`     | Æmber bonus icons purged this way |
| `CardsRevealed`         | cards revealed this way           |
| `CardsShuffledIntoDeck` | cards shuffled back this way      |
| `CreaturesHealed`       | creatures healed this way         |
| `DamageHealed`          | damage healed this way            |
| `DamagePrevented`       | damage prevented this way         |
| `ProducedThisWay`       | any `card.Tally` this way         |

`ProducedThisWay{Tally, Player}` covers `card.Tally.CreaturesDestroyed`,
`.CreaturesShuffledIntoDeck`, `.AemberLost`, `.CardsReturned`, `.CardsPurged`.

**Fractions** (the `Of:` field of a count): `HalfRoundedDown`, `HalfRoundedUp`,
`ThirdRoundedDown`, `ThirdRoundedUp`.

## Durations, zones, destinations, players

Creed of Nature scopes its grants to the rest of the turn:

```go
card.ForDuration{
  Duration: card.Duration.RemainderOfPlayerTurn,
  Effects: []card.Effect{
    card.GainKeywords{
      Target:   card.Target.Triggering,
      Keywords: []card.KeywordValue{card.Keyword.Skirmish},
      Duration: card.Duration.RemainderOfPlayerTurn,
    },
  },
}
```

| Duration                | Lasts until                    |
| ----------------------- | ------------------------------ |
| `RemainderOfPlayerTurn` | the end of this turn           |
| `OpponentNextTurn`      | the end of the opponent's turn |
| `StartOfPlayerNextTurn` | that player's next turn begins |
| `EndOfPlayerNextTurn`   | that player's next turn ends   |
| `UntilThisLeavesPlay`   | this card leaves play          |
| `Forever`               | the game ends                  |

Players an effect names:

| Player          | Names                            |
| --------------- | -------------------------------- |
| `Controller`    | the card's controller            |
| `Opponent`      | the controller's opponent        |
| `EachPlayer`    | both, in turn order              |
| `ChosenPlayer`  | a player an earlier step chose   |
| `ThatPlayer`    | the player the trigger bound     |
| `ItsOwner`      | the owner of the card in context |
| `ItsController` | its current controller           |
| `ItsOpponent`   | that controller's opponent       |

Zones (`card.Hand`, `.Deck`, `.Discard`, `.Archives`) and destinations
(`card.To.*`, `card.Into.*`) are tabulated under
[Zone movement](#zone-movement).

## Lasting effects and replacements

Some effects last "for the remainder of the turn" and attach to a later game
event — Full Moon gains Æmber whenever you play a creature, Charge! deals damage
whenever you play a creature, Dimension Door makes reaping steal instead of gain.
These are **lasting effects**, routed through a flat registry. Never add a
bespoke branch to the play or reap path for one (ADR 0007).

```go
card.ForRemainderOfTurn{
  On: card.Event.CreaturePlayed,
  Do: card.GainAember{Amount: 1},
}
```

Events (`card.Event.X`):

| Event                    | Fires when                     |
| ------------------------ | ------------------------------ |
| `CreaturePlayed`         | a creature is played           |
| `CardPlayed`             | any card is played             |
| `Reap`                   | a creature reaps               |
| `ReapAember`             | reaping is about to pay out    |
| `Fight`                  | a creature fights              |
| `Destroyed`              | a creature is destroyed        |
| `EnemyCreatureDestroyed` | an enemy creature is destroyed |
| `AemberAddedToPool`      | Æmber enters a pool            |
| `AemberTakenFromPool`    | Æmber leaves a pool            |
| `AemberStolen`           | Æmber is stolen                |
| `Forge`                  | a key is forged                |

Two flavors:

- A **reaction** runs _after_ the event, ordered together with the card abilities
  that fire on the same event (ADR 0013) — so "gain Æmber after you play a
  creature" and "deal damage after you play a creature" order for free. The nodes
  are `ForRemainderOfTurn`, `ForOpponentNextTurn`, `NextPlayed`,
  `FuseTriggersForTurn`, and `AlsoTriggersOn`. `Do` is a small composed effect.
- A **replacement** changes the event's _own outcome_ before it happens —
  `Instead{Of, With}` with `card.Steal`, `card.Capture`, or
  `card.FromCommonSupply`. A permanent that replaces an event for as long as it
  stays in play uses `card.WithReplaces(...)` instead.

The engine side — the flat `LastingEffect` record, `lastingActionOf`, and what
adding a new event costs — is in
[../internal/engine/AGENTS.md](../internal/engine/AGENTS.md).

## Always-on card shapes

These are card-level options, not effects — reach for them when the card's text
has no trigger.

### `card.WithConstant(card.ConstantAbility{...})`

A continuous effect a card in play applies to the creatures its `Target` reaches.

| Field                    | Sets                                       |
| ------------------------ | ------------------------------------------ |
| `Target`                 | which cards it reaches                     |
| `PowerBonus`             | a standing power bonus                     |
| `ArmorBonus`             | a standing armor bonus                     |
| `AssaultBonus`           | a standing assault bonus                   |
| `HazardousBonus`         | a standing hazardous bonus                 |
| `Per`                    | scales the bonuses by a `Count`            |
| `PerTarget`              | scales them per affected card              |
| `Keywords`               | keywords every target gains                |
| `Granted`                | abilities every target gains               |
| `CannotBeUsedTo`         | actions the targets may not take           |
| `AlsoTriggers`           | an extra timing the targets fire on        |
| `SpendAemberOnCard`      | lets Æmber on the card be spent            |
| `DisableTriggers`        | silences the targets' triggers             |
| `BlankText`              | blanks the targets' text boxes             |
| `RemovesTraits`          | strips the targets' traits                 |
| `SelectiveArchivePickup` | controller picks up any number of archives |
| `WhileOffFlank`          | applies only while this is off a flank     |
| `WhileInCenter`          | applies only while this is in the center   |
| `WhileCondition`         | applies only while a condition holds       |

"Each creature gains an ability" is a `Granted` list here — Annihilation Ritual
grants each creature the ability that purges that creature when it is destroyed:

```go
card.WithConstant(card.ConstantAbility{
  Target:  card.Target.EachCreature,
  Granted: []card.Ability{{
    Trigger: card.Trigger.Destroyed,
    Effect:  card.PurgeCreature{Target: card.Target.This},
  }},
}),
```

The engine gathers every Destroyed ability for all the creatures being destroyed,
then lets the active player order them. A `PurgeCreature{Target:
card.Target.This}` ability takes its creature out of play immediately, so that
creature's remaining Destroyed abilities do not resolve and final destruction
cleanup does not move it to its discard pile. **Never** implement this kind of
card as a global override in `discardDestroyed` or another leave-play path.

### `card.WithStatic(card.StaticModifier{...})`

What an upgrade grants its host. Plasma Nozzle is two flat bonuses:

```go
card.WithStatic(card.StaticModifier{
  AssaultBonus:      2,
  SplashAttackBonus: 2,
}),
```

| Field                  | Sets                                            |
| ---------------------- | ----------------------------------------------- |
| `PowerBonus`           | host power bonus                                |
| `ArmorBonus`           | host armor bonus                                |
| `AssaultBonus`         | host assault bonus                              |
| `HazardousBonus`       | host hazardous bonus                            |
| `SplashAttackBonus`    | host splash-attack bonus                        |
| `Per`                  | scales the bonuses by a `Count`                 |
| `Granted`              | abilities the host gains                        |
| `Keywords`             | keywords the host gains                         |
| `KeywordGrants`        | keywords with explicit `Host`/`Neighbors` reach |
| `KeyCostChange`        | how much the key cost moves                     |
| `AemberCannotBeStolen` | protects Æmber on the host                      |
| `Replaces`             | an event outcome the host replaces              |
| `WhileOnFlank`         | applies only on a flank                         |
| `ProtectsFromNonFlank` | shields the host from non-flank attackers       |
| `HouseOverride`        | changes the host's house                        |
| `SpendAemberOnCard`    | lets Æmber on the card be spent                 |
| `CannotBeUsedTo`       | actions the host may not take                   |

### `card.WithRestrictions(card.Restrictions{...})`

Continuous "cannot" rules. Gold Key Imp bars one key number:

```go
card.WithRestrictions(card.Restrictions{NoForgeKeyNumber: 3}),
```

| Field                     | Bars                               |
| ------------------------- | ---------------------------------- |
| `Fighting`                | a player's creatures from fighting |
| `Reaping`                 | a player's creatures from reaping  |
| `BonusIcons`              | bonus icons from resolving         |
| `CannotPlay`              | a card type from being played      |
| `PlayCardLimit`           | plays beyond a per-turn limit      |
| `Toll`                    | playing without paying a toll      |
| `UseCondition`            | use unless a condition holds       |
| `SkipForge`               | the forge step                     |
| `NoForgeKeyNumber`        | forging that numbered key          |
| `NoForgeWhileAheadOnKeys` | forging while ahead on keys        |
| `MustFightIfAble`         | anything but fighting, when able   |

### Other card-level options

**Stats and identity:**

| Option         | Sets                      |
| -------------- | ------------------------- |
| `WithPower`    | printed power             |
| `WithArmor`    | printed armor             |
| `WithPowerX`   | power read from a `Count` |
| `WithTraits`   | printed traits            |
| `WithKeywords` | printed keywords          |
| `WithBonus`    | printed bonus icons       |

**Combat shape:**

| Option                                | Sets                                     |
| ------------------------------------- | ---------------------------------------- |
| `WithAssault`                         | assault value                            |
| `WithHazardous`                       | hazardous value                          |
| `WithSplashAttack`                    | splash-attack value                      |
| `WithAttackDamage`                    | the damage it deals fighting             |
| `WithAttackKeywords`                  | keywords it attacks with                 |
| `WithAttackIgnores`                   | defensive keywords it ignores            |
| `WithNoDamageWhenAttacked`            | that it deals no damage back             |
| `WithFightRestriction`                | which creatures it may fight             |
| `WithTakesDamageFor`                  | damage it soaks for others               |
| `WithAlsoTakesNeighborFightDamage`    | that it shares a neighbor's fight damage |
| `WithTauntReachingNeighborsNeighbors` | extended taunt reach                     |
| `WithCannotBeDealtDamageBy`           | sources that cannot damage it            |

**Availability:**

| Option                        | Sets                                     |
| ----------------------------- | ---------------------------------------- |
| `WithCannotBeUsedTo`          | actions it may never take                |
| `WithCannotBeUsedWhile`       | a condition that locks it                |
| `WithDestroyedWhen`           | a condition that destroys it             |
| `WithCannotPlayWhile`         | a condition that bars playing it         |
| `WithHouseLock`               | an active-house requirement              |
| `WithPlayPermission`          | off-house cards it lets you play         |
| `WithPlayableAsUpgrade`       | that a creature may be played as upgrade |
| `WithTriggersFromDiscard`     | that its triggers work from discard      |
| `WithFriendlyEntersPlayReady` | that friendly arrivals enter ready       |
| `WithEntersPlay`              | an effect on entering play               |

**Economy:**

| Option                     | Sets                                 |
| -------------------------- | ------------------------------------ |
| `WithAemberThreshold`      | Æmber needed before an effect fires  |
| `WithAemberCost`           | Æmber paid to play it                |
| `WithAemberCannotBeStolen` | that its controller's Æmber is safe  |
| `WithSpendableAember`      | that Æmber banked on it can be spent |
| `WithGainsForgeAember`     | that it takes the forge's Æmber      |
| `WithKeyCost`              | a change to the key cost             |
| `WithBonusInstead`         | bonus icons substituted for others   |
| `WithDrawModifier`         | a change to the draw count           |
| `WithDrawModifierPer`      | a draw change scaled by a `Count`    |
| `WithDrawModifierOffFlank` | a draw change while off a flank      |
| `WithDrawModifierInCenter` | a draw change while in the center    |

## Keywords, bonus icons, card types

Special Agent Fingers prints one keyword and one bonus icon:

```go
card.WithKeywords(card.Keyword.Elusive),
card.WithBonus(card.Bonus.Aember),
```

| Keyword        | What it does                                  |
| -------------- | --------------------------------------------- |
| `Skirmish`     | takes no damage back when it fights           |
| `Poison`       | destroys any creature it damages              |
| `Elusive`      | ignores the first attack against it each turn |
| `Taunt`        | shields its neighbors from being fought       |
| `Versatile`    | may be played or used out of house            |
| `Alpha`        | must be the first card played this turn       |
| `Omega`        | ends the turn's play step after it resolves   |
| `Deploy`       | may enter anywhere in the battleline          |
| `Treachery`    | enters under the opponent's control           |
| `Invulnerable` | cannot be damaged or destroyed                |

Numeric keywords (assault, hazardous, splash-attack) are card options, not enum
values.

| Bonus icon      | Resolves to       |
| --------------- | ----------------- |
| `Bonus.Aember`  | gain one Æmber    |
| `Bonus.Capture` | capture one Æmber |
| `Bonus.Damage`  | deal one damage   |
| `Bonus.Draw`    | draw one card     |

Author bonus icons in printed top-to-bottom order.

Card types: `card.Type.Creature`, `.Artifact`, `.Upgrade`, `.Tactic` (KeyForge's
"action"), `.Any`.

## Deck-generation metadata

Not gameplay, but authored on the card — Alaka's Brew leads its own cluster:

```go
card.Provenance(card.WC, "2"),
card.LeadsCluster(alakasBrewCluster),
```

| Option             | What it does                                |
| ------------------ | ------------------------------------------- |
| `Provenance`       | records the source set and collector number |
| `OneCopyPerDeck()` | caps the card at one copy                   |
| `Houseless()`      | lets the card take any house                |
| `RarityWeight(w)`  | tunes how often it is drafted               |
| `Template(f)`      | materializes a family of generated cards    |
| `InCluster`        | joins a card family                         |
| `LeadsCluster`     | leads a card family                         |
| `Pulled`           | is pulled in by its cluster's lead          |
| `PullsMatching`    | pulls in whatever matches a filter          |
| `Gigantic(...)`    | authors a two-card gigantic creature        |

A card that names another card by name must share a cluster with it (ADR 0036);
ask the human for the pull rate rather than guessing.

## When the capability really is missing

Extend along the cheapest surface first:

1. A new **field** on an existing effect.
2. A new **Strategy** — a `Refinement`, `Count`, `Condition`, `Selection`,
   `Spread`, or `Chooser` — each of which carries its own text fragment.
3. A new **target filter**.
4. A new **effect node** in `effect_<mechanic>.go`. Grill-gated: present the case
   and stop (see the `implement-cards` skill).
5. A new **Resolver** capability, then new state. Both are last resorts; state
   must stay flat, pointerless, and comparable (ADR 0005).

Name what you add for the **mechanic**, not the card: prefer a `Target` or
`Refinement` filter, a `Count`, a `Duration` field, or a portion (`By: Half`)
over a name that spells out the whole card sentence (see
[style-guide.md](style-guide.md), "Composition and design").

Whatever you add: an engine test in the matching `effect_*_test.go`
(`internal/engine` is gated at 100%), a `RuleTerm` in the matching
`ruleterms_<section>.go` if it is player-facing (ADR 0018), and a client prompt if
it asks the player a new question (`internal/web/AGENTS.md`).
