# 47. An upgrade is a card in play, but not a source of triggered abilities

## Context

"In play" has two readings in `internal/engine`, and until now they were split by
accident rather than by rule.

The effect layer settled its reading first: `resolverCardsInPlay`
(`effect_cross_zone.go`) lists a player's creatures, their artifacts, **and the
upgrades attached to either**, with each host's upgrades immediately ahead of the
host. A card placed _under_ another card is not in play. That is pinned by
`TestCardsInPlayCountsUpgrades`, `TestInPlayCountsUpgrades` and
`TestEachCardInPlayReachesUpgrades`, and the kinds that name their types were
split off onto the row-only `creaturesAndArtifacts` / `creaturesAndArtifactsOf`.

The `Game` layer kept a single row-only helper, `allInPlay`, whose name claims
more than it delivers: it returns the battleline and the artifact row and nothing
else. About 57 call sites read it — ability scanning, phase processing, play
permissions, combat restrictions, key cost, Æmber capture and redirection, and
the card-conservation invariant.

Widening all 57 is wrong, and so is leaving all 57, because the two groups ask
different questions. The fact that separates them is how upgrade text is
modelled:

- An upgrade's printed text is authored as `Static.Granted` and is collected onto
  its **host**, by `upgradeGrantedTriggers` and `hasTrigger`
  (`game_abilities.go`). The host is the `source` of every ability the upgrade
  grants.
- The upgrade's own `Abilities` slice is consulted in exactly one place,
  `resolveUpgradePlay`, for `TriggerAfterPlay` only, and even that resolves with
  the host as `Source`.
- No implemented upgrade carries a `ConstantAbility`; of 81 upgrade cards, none
  defines one.

So an ability scan that walked upgrades directly would reach the same printed
text twice by two different routes — the divergent-matcher failure already ruled
out for `fireLastingBeforeFight`, whose doc comment records that a Before-Fight
scan must route through the trigger window rather than re-derive its filters.

## Decision

**An upgrade is a card in play. It is not a source of triggered abilities.**

Its triggered text belongs to the creature it is attached to and reaches the game
only through that host. Its standing rules — a constant ability, a play
permission or bar, a key-cost or forge modifier, an Æmber capture or redirect —
apply from where the upgrade sits, exactly as an artifact's would, because a
constant ability carries its own `Target` and needs no host to speak for it.

The `Game` layer mirrors the effect layer's vocabulary with two helpers, and
`allInPlay` is retired rather than kept alongside them: "all in play" is now a
false name, and leaving it is how the next reader picks the wrong helper.

- `g.cardsInPlay(player)` — creatures, artifacts, and the upgrades on either, in
  `resolverCardsInPlay`'s interleaved order (each host's upgrades ahead of the
  host). The `Game`-side twin of `resolverCardsInPlay`, pinned equal to it
  element for element, including order, by `TestCardsInPlayMatchesResolver`.
- `g.creaturesAndArtifacts(player)` — row-only; the rename of `allInPlay`. It
  shares its name with the free function in `target_select.go` on purpose: same
  name, same meaning, different receiver.

Each caller is decided by the question it asks, not by its file:

- **Scans for a standing rule a card in play imposes** read `cardsInPlay`: play
  permissions and play bars, toll, key cost and forge bars, Æmber capture,
  redirection and protection, draw modifiers, house presence, enters-ready
  grants, bonus-icon bars, armor and damage-redirect watchers, and
  `constantAbilitiesOf`.
- **Trigger and reaction scans stay row-only** and read
  `creaturesAndArtifacts`: the ability window's `addBoard` / `addSide` /
  `addTurnScoped`, `emitSide`, `actionReactions`, `fightReactions`,
  `forgeKeyReactions`, `afterBonusReaction`, `EmitAemberStolenFrom`, and the
  discard-from-hand watcher. So does every read of per-card row state — the ready
  and start-of-turn phases, `revertTemporaryCreatures`, `actedThisTurn`, and
  `tallyPlacement`'s chain walk, which reaches upgrades through their hosts.

A caller that stays row-only carries one sentence in its doc comment saying why,
so the next agent does not re-litigate it.

## Consequences

- Most widenings are unobservable to any printed card today, because no upgrade
  carries a constant ability or a play permission. They are still testable and
  still tested: `internal/engine` tests build their own `CardDefinition`
  blueprints (`helpers_test.go`), so an engine test can define an upgrade with a
  `ConstantAbility` and pin the widened behaviour. The 100% coverage gate is met
  by testing the rule, not by avoiding it.
- Order changes for a widened caller only by insertion: the relative order of row
  cards is untouched, so first-match scans (`aemberCaptorFor`,
  `stolenRedirectSource`, `forgeAemberGainer`, `damageRedirect`) keep today's
  answer unless an upgrade genuinely matches.
- Two implementations of one list now exist — `cardsInPlay` reads state directly,
  `resolverCardsInPlay` goes through the port so test doubles see it.
  `TestCardsInPlayMatchesResolver` holds them together; without it they drift.
- The first card that wants an upgrade to fire a triggered ability of its own
  does not get a second scan. It either authors the text as `Static.Granted` like
  every other upgrade, or it changes this decision.
- `inPlay` keeps its own two-zone `contains` loop. It allocates nothing and sits
  in a hot path; measure with `mage profile` before routing it through a helper
  that builds a slice.

## Player-facing wording

The rulebook's Upgrade card-type term (`ruleterms_cardtype.go`) states the
player's half of this rule: an upgrade is in play, and the abilities it grants
are used by the creature it is attached to. No new rulebook term is added — this
is engine structure, and "in play" already reads correctly on every printed card.
