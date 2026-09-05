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
	// SelfHouse is the sentinel a card uses to name its own house instead of
	// spelling it out: Battle Fleet, a Mars card, reveals SelfHouse cards. NewCard
	// substitutes the card's house for it when the definition is built (see
	// self_house.go), so the sentinel never reaches text, resolution, or state —
	// and the printed house can never drift from the house on the card.
	SelfHouse
	// NumHouses is the number of house slots a card can actually occupy: HouseNone
	// through Untamed, but not SelfHouse, which is resolved away when the card is
	// built and so never indexes state.
	NumHouses = int(Untamed) + 1
)

// houseNames maps a House to its printed name, indexed by the enum value.
var houseNames = [...]string{
	"None", "Brobnar", "Dis", "Logos", "Mars", "Sanctum", "Shadows", "Untamed",
}

// String returns the printed house name.
func (h House) String() string {
	if h == SelfHouse {
		return "this card's house"
	}
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
	// Connected is the rarity of a card that never rolls on its own; it enters a
	// deck only when another card's connection pulls it into the same pod.
	Connected Rarity = "Connected"
)

// CardType is one of the four Vactrol/KeyForge card types. It is a small enum
// rather than a string so a type-scoped bar stays a couple of bytes of flat state
// instead of a string header (see GameState.CannotPlayTypeThis).
type CardType uint8

const (
	// TypeUnset is the zero value: a type-scoped effect that names no type matches
	// every card rather than none.
	TypeUnset CardType = iota
	// A creature is a unit you play into your battleline. Once it is ready, it can
	// reap for Æmber, fight an enemy creature, or use an "Action:" ability.
	Creature
	// A tactic (KeyForge's "action" card type, renamed to free the word "Action"
	// for the ability) is a one-shot card: its effect resolves as you play it, and
	// it then goes straight to your discard pile.
	Tactic
	// An artifact is a permanent card you play alongside your creatures. It stays
	// in play until something removes it and is typically used for its "Action:"
	// ability.
	Artifact
	// An upgrade attaches to a creature as you play it, changing that creature's
	// stats or granting it keywords and abilities for as long as it stays attached.
	Upgrade
	// AnyType is the wildcard a type-scoped effect uses to mean "every card type"
	// rather than one of them — Treasure Map bars playing cards outright, not cards
	// of a particular type. No card is ever of this type; it exists so a bar can say
	// "all" while staying one comparable CardType.
	AnyType
)

// cardTypeNames is the printed word for each type; the unset zero renders empty,
// which is what the "names no type" callers expect.
var cardTypeNames = map[CardType]string{
	Creature: "Creature",
	Tactic:   "Tactic",
	Artifact: "Artifact",
	Upgrade:  "Upgrade",
	AnyType:  "Card",
}

// String renders the type as its printed word.
func (t CardType) String() string { return cardTypeNames[t] }

// CardTypes lists every real card type in rulebook order. It excludes the
// TypeUnset and AnyType sentinels (which are not types a card can be), so
// anything that must cover all card types — the rulebook completeness check —
// ranges over this and cannot silently miss one.
func CardTypes() []CardType {
	return []CardType{Creature, Tactic, Artifact, Upgrade}
}

// reacts reports whether a filter of this type matches a card of type other.
// TypeUnset matches everything, AnyType matches the two types that stay in play
// under their own name (a creature or an artifact), and a named type matches only
// itself.
func (t CardType) reacts(other CardType) bool {
	switch t {
	case TypeUnset:
		return true
	case AnyType:
		return other == Creature || other == Artifact
	default:
		return other == t
	}
}

// Trait is a flavor/type label printed on a card (e.g. "Giant", "Weapon").
// Traits carry no inherent rules meaning on their own; other cards reference
// them (a Destroy that targets "each Scientist trait creature"). KeyForge keeps
// coining new traits with every set, but each one still comes from a fixed,
// printed vocabulary, so — like Keyword — this is a closed enum rather than a
// freeform string: a misspelled trait cannot compile. A set that prints a trait
// not yet listed below needs it added here first.
type Trait int

// The traits a card can carry, in alphabetical order.
const (
	// traitUnset is the invalid zero value: a Trait field left unset (e.g.
	// PutFromDiscard.Trait, Target.trait) means "no trait filter", not a real
	// trait.
	traitUnset Trait = iota
	Agent
	Ally
	Angel
	Beast
	Cleric
	Cyborg
	Demon
	Dragon
	Elf
	Faerie
	Fungus
	Giant
	Goblin
	Equation
	Horseman
	Human
	Imp
	Insect
	Item
	Knight
	Law
	Location
	Martian
	Merchant
	Monk
	Mutant
	Niffle
	Power
	Priest
	Quest
	Ranger
	Rat
	Robot
	Scientist
	Shard
	Soldier
	Specter
	Spirit
	Thief
	Vehicle
	Weapon
	Witch
	// traitCount bounds the enum; it is not a trait.
	traitCount
)

// traitNames maps a Trait to its printed word, indexed by the enum value.
var traitNames = [traitCount]string{
	Agent:     "Agent",
	Ally:      "Ally",
	Angel:     "Angel",
	Beast:     "Beast",
	Cleric:    "Cleric",
	Cyborg:    "Cyborg",
	Demon:     "Demon",
	Dragon:    "Dragon",
	Elf:       "Elf",
	Equation:  "Equation",
	Faerie:    "Faerie",
	Fungus:    "Fungus",
	Giant:     "Giant",
	Goblin:    "Goblin",
	Horseman:  "Horseman",
	Human:     "Human",
	Imp:       "Imp",
	Insect:    "Insect",
	Item:      "Item",
	Knight:    "Knight",
	Law:       "Law",
	Location:  "Location",
	Martian:   "Martian",
	Merchant:  "Merchant",
	Mutant:    "Mutant",
	Niffle:    "Niffle",
	Monk:      "Monk",
	Power:     "Power",
	Priest:    "Priest",
	Quest:     "Quest",
	Ranger:    "Ranger",
	Rat:       "Rat",
	Robot:     "Robot",
	Scientist: "Scientist",
	Shard:     "Shard",
	Soldier:   "Soldier",
	Specter:   "Specter",
	Spirit:    "Spirit",
	Thief:     "Thief",
	Vehicle:   "Vehicle",
	Weapon:    "Weapon",
	Witch:     "Witch",
}

// String returns the trait's printed word, or "" for the unset zero value.
func (t Trait) String() string {
	if t <= traitUnset || t >= traitCount {
		return ""
	}
	return traitNames[t]
}

// Keyword is a rules shorthand a card can have (e.g. Skirmish, Poison). Unlike
// Trait, which any card may coin, the keywords are a closed set the rulebook
// defines, so this is an enum: a keyword that does not exist cannot be written.
type Keyword int

const (
	// keywordUnset is the invalid zero value: an effect that names a keyword must
	// say which one rather than leave it unset.
	keywordUnset Keyword = iota
	// A creature with Skirmish takes no damage when it is used to fight: it deals
	// its power to the enemy creature but takes none back.
	Skirmish
	// Any amount of damage dealt to a creature with Poison destroys it, however
	// much power it has left.
	Poison
	// Elusive: the first time this creature is chosen to be fought each turn, no
	// pending fight damage is dealt by or to it. Later fights that same turn deal
	// damage normally.
	Elusive
	// Taunt: this creature's neighbors cannot be chosen to be fought unless they
	// have taunt themselves, so a Taunt creature shields the creatures beside it.
	Taunt
	// A card with Versatile may, once in play, be used (reap/fight/action) as if
	// it belonged to the active house. It does not relax playing from hand — a
	// Versatile card is still played only when its own house is the one chosen
	// this turn.
	Versatile
	// A card with Alpha can only be played as the first thing its player does on
	// their turn: it cannot be played once that player has played, used, or
	// discarded any other card this turn (First Blood, Unlocked Gateway's
	// counterpart). It is a restriction on playing, so it never appears as a
	// granted or lost creature keyword.
	Alpha
	// A card with Omega ends the current step of the turn the moment it resolves:
	// no more cards may be played, used, or discarded for the rest of that step,
	// except through pending abilities and effects still resolving (Unlocked
	// Gateway). Play then continues to the next step, so more cards can still be
	// played later that turn. Like Alpha it constrains playing, not combat.
	Omega
	// A creature with Deploy may enter play at any position in its controller's
	// battleline, not only on a flank — its controller chooses the spot as it is
	// played ("Lion" Bautrem, Challe the Safeguard). It matters only while the
	// creature is being played, so it too is never granted or lost.
	Deploy
	// keywordCount bounds the enum; it is not a keyword.
	keywordCount
)

// keywordNames is the printed word for each keyword; the unset zero renders
// empty.
var keywordNames = [keywordCount]string{
	Skirmish:  "Skirmish",
	Poison:    "Poison",
	Elusive:   "Elusive",
	Taunt:     "Taunt",
	Versatile: "Versatile",
	Alpha:     "Alpha",
	Omega:     "Omega",
	Deploy:    "Deploy",
}

// String returns the keyword's printed word, capitalized as a card prints it.
func (k Keyword) String() string {
	if !k.valid() {
		return ""
	}
	return keywordNames[k]
}

// valid reports whether k names a real keyword (not the unset zero value).
func (k Keyword) valid() bool { return k > keywordUnset && k < keywordCount }

// bit is the bit k occupies in GameState.KeywordsLost, so a "for the remainder of
// the turn, each creature loses <keyword>" effect can be held as one flat
// comparable value. Being derived from the enum, it cannot fall out of step with
// it the way a hand-maintained table could.
func (k Keyword) bit() uint8 {
	if !k.valid() {
		return 0
	}
	return 1 << (k - 1)
}

// Keywords lists every keyword in rulebook order. It is the canonical
// enumeration: anything that must cover all keywords ranges over this, so a
// keyword added above cannot be silently missed.
func Keywords() []Keyword {
	all := make([]Keyword, 0, keywordCount-1)
	for k := keywordUnset + 1; k < keywordCount; k++ {
		all = append(all, k)
	}
	return all
}

// Trigger identifies when an ability's effect resolves.
type Trigger int

const (
	// triggerUnset is the invalid zero value: an ability must name its trigger
	// rather than leave it unset.
	triggerUnset Trigger = iota
	// A Play ability resolves right after you play the card from your hand. On a
	// creature or artifact it fires as the card enters play; on an action it is the
	// card's one-shot effect.
	TriggerAfterPlay
	// A Reap ability resolves after you use a ready creature to reap. Reaping gains
	// you 1 Æmber and exhausts the creature; the ability resolves in addition.
	TriggerAfterReap
	// A Fight ability resolves after a creature you used to fight has dealt and
	// taken its combat damage and any resulting destruction has been carried out.
	TriggerAfterFight
	// An Action ability is one you resolve by using the card directly, without
	// reaping or fighting; using it this way exhausts the card.
	TriggerAction
	// This ability resolves after its controller forges a key.
	TriggerAfterForgeKey
	// This ability resolves after any creature enters play, including creatures
	// your opponent plays.
	TriggerAfterCreatureEnters
	// A Destroyed ability resolves as the card is destroyed, before it reaches the
	// discard pile, so it can still act on the board it is leaving.
	TriggerDestroyed
	// A Before Fight ability resolves when a creature is used to fight, before any
	// combat damage is dealt.
	TriggerBeforeFight
	// This ability resolves on a creature that survives a fight in which the other
	// combatant was destroyed; the destroyed creature is the one referred to as
	// "it".
	TriggerAfterDestroyedFighting
	// This ability resolves after its controller plays a card — a creature,
	// artifact, or action — from hand. Putting a card into play by another effect is
	// not "playing" it and does not fire this (that is TriggerAfterCreatureEnters).
	// AfterCardPlayed narrows it to a house and/or type.
	TriggerAfterCardPlayed
	// An Enters Play ability resolves on a creature as it enters play, whatever
	// brought it in — the creature's own reaction to arriving, such as Chuff Ape
	// entering stunned. It is fired on the entering creature by the enter-play event
	// (Game.emitCreatureEnters), so a new "as it enters play" behavior is just
	// another ability rather than a special case in the play path.
	TriggerEntersPlay
	// An End of Turn ability resolves during the end of its controller's turn,
	// after cards ready and the controller draws (Shaffles drains the opponent at
	// each turn's end).
	TriggerEndOfTurn
	// A Start of Turn ability resolves at the start of its controller's turn,
	// before they forge, so an ability that changes what a key costs still has time
	// to.
	TriggerStartOfTurn
	// This ability resolves after its controller chooses their active house at the
	// start of the turn — the only "choose a house" step it watches. Changing houses
	// mid-turn by another effect is not this start-of-turn choice and does not fire
	// it. The ability names the house it cares about (Jehu the Bureaucrat gains Æmber
	// only when Sanctum is chosen), which need not be the card's own house.
	TriggerAfterChooseHouse
	// This ability resolves after an enemy creature is destroyed during its
	// controller's turn (Pile of Skulls captures Æmber onto a friendly creature
	// whenever an enemy creature is destroyed on your turn).
	TriggerAfterEnemyCreatureDestroyed
	// This ability resolves after your opponent plays a card (Teliga gains its
	// controller Æmber whenever the opponent plays a card).
	TriggerAfterEnemyCardPlayed
	// This ability resolves after its controller uses a card — reaps or fights with a
	// creature, or fires an "Action:" — with the used card as "it" (Veylan Analyst
	// gains Æmber whenever you use an artifact).
	TriggerAfterUse
	// This ability resolves after its controller discards a card from their hand,
	// with the discarded card as "it" (Rock-Hurling Giant). Discarding from anywhere
	// else — the top of a deck, the archives — is not this.
	TriggerAfterDiscardFromHand
	// A Leaves Play ability resolves as the card leaves play by any route — destroyed,
	// purged, returned to hand, archived, or shuffled away. It fires before the card's
	// teardown, so the card is still on the board when it resolves. TriggerDestroyed is
	// the narrower "only when destroyed" version.
	TriggerLeavesPlay
	// This ability resolves after the card that carries it is itself used — reaped
	// or fought with, or fired as an "Action:". Unlike TriggerAfterUse (which fires
	// on the user's other cards with the used card as "it"), this fires on the used
	// card itself, so an upgrade can punish its own host: Containment Field destroys
	// its host after it is used, Bonerot Venom deals it 2 damage.
	TriggerAfterUsedSelf
	// triggerCount bounds the enum so Triggers can range it; it is not a trigger.
	triggerCount
)

// Triggers lists every real trigger in declaration order (excluding the unset
// zero). Printed reports which of them a card actually prints a prefix for, and
// String gives that trigger's rulebook heading; the two implicit triggers
// (EntersPlay, AfterChooseHouse) fire with no printed text and no heading.
func Triggers() []Trigger {
	all := make([]Trigger, 0, triggerCount-1)
	for t := triggerUnset + 1; t < triggerCount; t++ {
		all = append(all, t)
	}
	return all
}

// Printed reports whether a card prints a text prefix for this trigger. A printed
// trigger is player-facing and must carry a rulebook entry (its String heading);
// an unprinted one (EntersPlay, AfterChooseHouse) fires implicitly.
func (t Trigger) Printed() bool {
	text, _ := t.prefix()
	return text != ""
}

// String returns the trigger's name. For a printed trigger this is also its
// rulebook heading (the Title of its SectionAbility RuleTerm); the two implicit
// triggers name themselves for logging though they head no rulebook entry, and
// the unset zero renders empty.
func (t Trigger) String() string {
	switch t {
	case TriggerAfterPlay:
		return "Play"
	case TriggerAfterReap:
		return "Reap"
	case TriggerAfterFight:
		return "Fight"
	case TriggerAction:
		return "Action"
	case TriggerAfterForgeKey:
		return "After You Forge a Key"
	case TriggerAfterCreatureEnters:
		return "After a Creature Enters Play"
	case TriggerDestroyed:
		return "Destroyed"
	case TriggerBeforeFight:
		return "Before Fight"
	case TriggerAfterDestroyedFighting:
		return "After a Creature Is Destroyed Fighting"
	case TriggerAfterCardPlayed:
		return "After You Play a Card"
	case TriggerAfterEnemyCreatureDestroyed:
		return "After an Enemy Creature Is Destroyed"
	case TriggerAfterEnemyCardPlayed:
		return "After Your Opponent Plays a Card"
	case TriggerAfterUse:
		return "After You Use a Card"
	case TriggerAfterDiscardFromHand:
		return "After You Discard a Card From Your Hand"
	case TriggerAfterUsedSelf:
		return "After This Creature Is Used"
	case TriggerEndOfTurn:
		return "End of Turn"
	case TriggerStartOfTurn:
		return "Start of Turn"
	case TriggerLeavesPlay:
		return "Leaves Play"
	case TriggerEntersPlay:
		return "Enters Play"
	case TriggerAfterChooseHouse:
		return "After Choosing a House"
	default:
		return ""
	}
}

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
	case TriggerAfterEnemyCardPlayed:
		return "After your opponent plays a card, ", false
	case TriggerAfterUse:
		return "After you use a card, ", false
	case TriggerAfterDiscardFromHand:
		return "After you discard a card from your hand, ", false
	case TriggerAfterUsedSelf:
		return "After " + SelfName + " is used, ", false
	case TriggerLeavesPlay:
		return "Leaves Play: ", true
	case TriggerEndOfTurn:
		return "At the end of your turn, ", false
	case TriggerStartOfTurn:
		return "At the start of your turn, ", false
	default:
		return "", true
	}
}
