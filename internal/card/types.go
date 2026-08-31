package card

import "github.com/dmikalova/vactrol/internal/engine"

// The enum-like card categories, re-exported as grouped namespaces so related
// values stay together and read unambiguously — card.House.Brobnar is plainly a
// house, card.Type.Tactic plainly a card type. This mirrors the engine's
// types.go. Treat these package-level vars as read-only.

// Trait and Player are the loose value types authors name directly (Player as an
// effect field, whose values are card.Controller / card.Opponent below).
type (
	// Trait is a freeform label such as "Giant" or "Weapon".
	Trait = engine.Trait
	// Player is the relative player an effect targets: card.Controller or card.Opponent.
	Player = engine.Player
)

// House groups the faction values, e.g. card.House.Brobnar.
var House = houses{
	Brobnar: engine.Brobnar,
	Dis:     engine.Dis,
	Logos:   engine.Logos,
	Mars:    engine.Mars,
	Sanctum: engine.Sanctum,
	Shadows: engine.Shadows,
	Untamed: engine.Untamed,
}

type houses struct {
	Brobnar, Dis, Logos, Mars, Sanctum, Shadows, Untamed engine.House
}

// Type groups the card-type values, e.g. card.Type.Creature.
var Type = cardTypes{
	Creature: engine.Creature,
	Tactic:   engine.Tactic,
	Artifact: engine.Artifact,
	Upgrade:  engine.Upgrade,
}

type cardTypes struct {
	Creature, Tactic, Artifact, Upgrade engine.CardType
}

// Rarity groups the rarity values, e.g. card.Rarity.Common.
var Rarity = rarities{
	Common:   engine.Common,
	Uncommon: engine.Uncommon,
	Rare:     engine.Rare,
	Special:  engine.Special,
	Fixed:    engine.Fixed,
}

type rarities struct {
	Common, Uncommon, Rare, Special, Fixed engine.Rarity
}

// Keyword groups the keyword values, e.g. card.Keyword.Skirmish.
var Keyword = keywords{
	Skirmish:  engine.Skirmish,
	Poison:    engine.Poison,
	Elusive:   engine.Elusive,
	Taunt:     engine.Taunt,
	Versatile: engine.Versatile,
}

type keywords struct {
	Skirmish, Poison, Elusive, Taunt, Versatile engine.Keyword
}

// Keywords builds the keyword slice for an upgrade's granted keywords, e.g.
// card.StaticModifier{Keywords: card.Keywords(card.Keyword.Skirmish)}. It exists
// because card.Keyword is the value namespace, so a []card.Keyword literal can't
// be written directly.
func Keywords(k ...engine.Keyword) []engine.Keyword { return k }

// Trigger groups the ability triggers, e.g. card.Trigger.Play or
// card.Trigger.AfterForgeKey.
var Trigger = triggers{
	Play:                        engine.TriggerAfterPlay,
	Reap:                        engine.TriggerAfterReap,
	Fight:                       engine.TriggerAfterFight,
	BeforeFight:                 engine.TriggerBeforeFight,
	Action:                      engine.TriggerAction,
	AfterForgeKey:               engine.TriggerAfterForgeKey,
	AfterCreatureEnters:         engine.TriggerAfterCreatureEnters,
	Destroyed:                   engine.TriggerDestroyed,
	AfterDestroyedFighting:      engine.TriggerAfterDestroyedFighting,
	AfterCardPlayed:             engine.TriggerAfterCardPlayed,
	EndOfTurn:                   engine.TriggerEndOfTurn,
	AfterChooseHouse:            engine.TriggerAfterChooseHouse,
	AfterEnemyCreatureDestroyed: engine.TriggerAfterEnemyCreatureDestroyed,
}

type triggers struct {
	Play,
	Reap,
	Fight,
	BeforeFight,
	Action,
	AfterForgeKey,
	AfterCreatureEnters,
	Destroyed,
	AfterDestroyedFighting,
	AfterCardPlayed,
	EndOfTurn,
	AfterChooseHouse,
	AfterEnemyCreatureDestroyed engine.Trigger
}

// Controller and Opponent are the two players an effect can target, relative to
// the card's controller: card.Controller (the player who controls the card) or
// card.Opponent (their opponent). card.EachPlayer means both, for effects that
// reach everyone at once (e.g. a key-cost change on each player's keys).
var (
	Controller = engine.Controller
	Opponent   = engine.Opponent
	EachPlayer = engine.EachPlayer
	ItsOwner   = engine.ItsOwner
)
