// Package engine is the core rules engine for Vactrol, a KeyForge-style
// card game. This file defines the fundamental enumerations shared across the
// engine: houses, rarities, card types, traits, keywords, and ability triggers.
package engine

import "strings"

// House is the faction a card belongs to. It is a small integer enum so it can
// live inside the flat, value-copyable GameState without introducing pointers.
type House uint8

// The houses a card can belong to, in canonical order.
const (
	// HouseNone means no house is chosen/assigned.
	HouseNone House = iota
	Brobnar
	Dis
	Logos
	Mars
	Sanctum
	Shadows
	Untamed
	// NumHouses is the number of house enum slots, including HouseNone, so house-
	// indexed state arrays can use a flat fixed size.
	NumHouses = int(Untamed) + 1
)

// houseNames maps a House to its printed name, indexed by the enum value.
var houseNames = [...]string{
	"None", "Brobnar", "Dis", "Logos", "Mars", "Sanctum", "Shadows", "Untamed",
}

// String returns the printed house name.
func (h House) String() string {
	if int(h) < len(houseNames) {
		return houseNames[h]
	}
	return "Unknown"
}

// ParseHouse resolves a house name (case-insensitive) to its House value. The
// boolean result is false when the name matches no house.
func ParseHouse(name string) (House, bool) {
	for h := HouseNone + 1; int(h) < len(houseNames); h++ {
		if strings.EqualFold(houseNames[h], name) {
			return h, true
		}
	}
	return HouseNone, false
}

// Rarity describes how frequently a card appears when generating a deck. It only
// ever lives in the read-only card catalog, so a string is fine here.
type Rarity string

// The card rarities.
const (
	Common   Rarity = "Common"
	Uncommon Rarity = "Uncommon"
	Rare     Rarity = "Rare"
	Special  Rarity = "Special"
	// Fixed is the rarity of a card that always appears in a fixed slot of its
	// house's decks (rendered "FIXED", as KeyForge prints it).
	Fixed Rarity = "FIXED"
)

// CardType is one of the four Vactrol/KeyForge card types.
type CardType string

const (
	// A creature is a unit you play into your battleline. Once it is ready, it can
	// reap for Aember, fight an enemy creature, or use an "Action:" ability.
	//
	//rulebook:cardtype Creature
	Creature CardType = "Creature"
	// A tactic (KeyForge's "action" card type, renamed to free the word "Action"
	// for the ability) is a one-shot card: its effect resolves as you play it, and
	// it then goes straight to your discard pile.
	//
	//rulebook:cardtype Tactic
	Tactic CardType = "Tactic"
	// An artifact is a permanent card you play alongside your creatures. It stays
	// in play until something removes it and is typically used for its "Action:"
	// ability.
	//
	//rulebook:cardtype Artifact
	Artifact CardType = "Artifact"
	// An upgrade attaches to a creature as you play it, changing that creature's
	// stats or granting it keywords and abilities for as long as it stays attached.
	//
	//rulebook:cardtype Upgrade
	Upgrade CardType = "Upgrade"
)

// Trait is a flavor/type label printed on a card (e.g. "Giant", "Weapon").
// Traits carry no inherent rules meaning on their own; other cards reference them.
type Trait string

// Keyword is a rules shorthand a card can have (e.g. Skirmish, Poison).
type Keyword string

const (
	// A creature with Skirmish takes no damage when it is used to fight: it deals
	// its power to the enemy creature but takes none back.
	//
	//rulebook:keyword Skirmish
	Skirmish Keyword = "Skirmish"
	// Any amount of damage dealt to a creature with Poison destroys it, however
	// much power it has left.
	//
	//rulebook:keyword Poison
	Poison Keyword = "Poison"
	// Elusive: defined for completeness; behavior not yet implemented.
	Elusive Keyword = "Elusive"
	// Taunt: defined for completeness; behavior not yet implemented.
	Taunt Keyword = "Taunt"
	// A card with Versatile may, once in play, be used (reap/fight/action) as if
	// it belonged to the active house. It does not relax playing from hand — a
	// Versatile card is still played only when its own house is the one chosen
	// this turn.
	//
	//rulebook:keyword Versatile
	Versatile Keyword = "Versatile"
)

// Trigger identifies when an ability's effect resolves.
type Trigger int

const (
	// triggerUnset is the invalid zero value: an ability must name its trigger
	// rather than leave it unset.
	triggerUnset Trigger = iota
	// A Play ability resolves right after you play the card from your hand. On a
	// creature or artifact it fires as the card enters play; on an action it is the
	// card's one-shot effect.
	//
	//rulebook:ability Play
	TriggerAfterPlay
	// A Reap ability resolves after you use a ready creature to reap. Reaping gains
	// you 1 Aember and exhausts the creature; the ability resolves in addition.
	//
	//rulebook:ability Reap
	TriggerAfterReap
	// A Fight ability resolves after a creature you used to fight has dealt and
	// taken its combat damage and any resulting destruction has been carried out.
	//
	//rulebook:ability Fight
	TriggerAfterFight
	// An Action ability is one you resolve by using the card directly, without
	// reaping or fighting; using it this way exhausts the card.
	//
	//rulebook:ability Action
	TriggerAction
	// This ability resolves after its controller forges a key.
	//
	//rulebook:ability After You Forge a Key
	TriggerAfterForgeKey
	// This ability resolves after any creature enters play, including creatures
	// your opponent plays.
	//
	//rulebook:ability After a Creature Enters Play
	TriggerAfterCreatureEnters
	// A Destroyed ability resolves as the card is destroyed, before it reaches the
	// discard pile, so it can still act on the board it is leaving.
	//
	//rulebook:ability Destroyed
	TriggerDestroyed
	// A Before Fight ability resolves when a creature is used to fight, before any
	// combat damage is dealt.
	//
	//rulebook:ability Before Fight
	TriggerBeforeFight
	// This ability resolves on a creature that survives a fight in which the other
	// combatant was destroyed; the destroyed creature is the one referred to as
	// "it".
	//
	//rulebook:ability After a Creature Is Destroyed Fighting
	TriggerAfterDestroyedFighting
	// This ability resolves after its controller plays a card — a creature,
	// artifact, or action — from hand. Putting a card into play by another effect is
	// not "playing" it and does not fire this (that is TriggerAfterCreatureEnters).
	// AfterCardPlayed narrows it to a house and/or type.
	//
	//rulebook:ability After You Play a Card
	TriggerAfterCardPlayed
	// An Enters Play ability resolves on a creature as it enters play, whatever
	// brought it in — the creature's own reaction to arriving, such as Chuff Ape
	// entering stunned. It is fired on the entering creature by the enter-play event
	// (Game.emitCreatureEnters), so a new "as it enters play" behavior is just
	// another ability rather than a special case in the play path.
	TriggerEntersPlay
	// An End of Turn ability resolves during the end of its controller's turn,
	// before cards ready and the controller draws (Shaffles drains the opponent at
	// each turn's end).
	TriggerEndOfTurn
	// This ability resolves after its controller chooses their active house at the
	// start of the turn — the only "choose a house" step it watches. Changing houses
	// mid-turn by another effect is not this start-of-turn choice and does not fire
	// it. The ability names the house it cares about (Jehu the Bureaucrat gains Æmber
	// only when Sanctum is chosen), which need not be the card's own house.
	TriggerAfterChooseHouse
	// This ability resolves after an enemy creature is destroyed during its
	// controller's turn — Pile of Skulls has a friendly creature capture Æmber
	// whenever an enemy creature is destroyed on your turn. It is a persistent
	// reaction on an in-play card, fired only for the active player's cards.
	TriggerAfterEnemyCreatureDestroyed
)

// prefix returns the printed text prefix for a trigger and whether the effect
// clause that follows should be capitalized (colon-style triggers start a new
// sentence; comma-style triggers continue one).
// valid reports whether t names a real trigger (not the unset zero value).
func (t Trigger) valid() bool { return t != triggerUnset }

func (t Trigger) prefix() (text string, capitalizeEffect bool) {
	switch t {
	case TriggerAfterPlay:
		return "Play: ", true
	case TriggerAfterReap:
		return "Reap: ", true
	case TriggerAfterFight:
		return "Fight: ", true
	case TriggerBeforeFight:
		return "Before Fight: ", true
	case TriggerAction:
		return "Action: ", true
	case TriggerDestroyed:
		return "Destroyed: ", true
	case TriggerAfterForgeKey:
		return "After you forge a key, ", false
	case TriggerAfterCreatureEnters:
		return "After a creature enters play, ", false
	case TriggerAfterDestroyedFighting:
		return "After a creature is destroyed fighting " + SelfName + ", ", false
	case TriggerAfterEnemyCreatureDestroyed:
		return "After an enemy creature is destroyed during your turn, ", false
	case TriggerAfterCardPlayed:
		return "After you play a card, ", false
	case TriggerEndOfTurn:
		return "At the end of your turn, ", false
	default:
		return "", true
	}
}
