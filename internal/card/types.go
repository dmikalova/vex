package card

import "github.com/dmikalova/vactrol/internal/engine"

// The enum-like card categories, re-exported as grouped namespaces so related
// values stay together and read unambiguously — card.House.Brobnar is plainly a
// house, card.Type.Tactic plainly a card type. This mirrors the engine's
// types.go. Treat these package-level vars as read-only.

// Trait and Player are the loose value types authors name directly (Player as an
// effect field, whose values are card.Controller / card.Opponent below).
type (
	// Trait is the value type for a card.Traits.X constant, e.g. as the field type
	// of an effect that filters by trait (PutFromDiscard.Trait).
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
	Self:    engine.SelfHouse,
}

type houses struct {
	// Brobnar is the house of giants and brawlers.
	Brobnar engine.House
	// Dis is the house of demons.
	Dis engine.House
	// Logos is the house of scientists and invention.
	Logos engine.House
	// Mars is the house of the martians.
	Mars engine.House
	// Sanctum is the house of knights and spirits.
	Sanctum engine.House
	// Shadows is the house of thieves and elves.
	Shadows engine.House
	// Untamed is the house of nature and beasts.
	Untamed engine.House
	// Self is the card's own house, filled in when the card is built. Use it
	// whenever an ability names the house the card itself belongs to — Battle
	// Fleet reveals card.House.Self cards — so the two can never drift apart.
	// A card that names a *different* house names it outright.
	Self engine.House
}

// Type groups the card-type values, e.g. card.Type.Creature.
var Type = cardTypes{
	Creature: engine.Creature,
	Tactic:   engine.Tactic,
	Artifact: engine.Artifact,
	Upgrade:  engine.Upgrade,
	Any:      engine.AnyType,
}

type cardTypes struct {
	// Creature is a creature card.
	Creature engine.CardType
	// Tactic is an action card (KeyForge's "action" type; card-wording rule 19).
	Tactic engine.CardType
	// Artifact is an artifact card.
	Artifact engine.CardType
	// Upgrade is an upgrade card that attaches to a creature.
	Upgrade engine.CardType
	// Any is the wildcard type, printed as "card". Where an effect waits for a card
	// to enter play it means "creature or artifact" — the two types that stay there.
	Any engine.CardType
}

// Subject groups the cards a condition can name instead of saying "it", e.g.
// card.Subject.DiscardedCard.
var Subject = subjects{
	DiscardedCard: engine.DiscardedCard,
}

type subjects struct {
	// DiscardedCard names the card an effect just discarded.
	DiscardedCard engine.Subject
}

// Rarity groups the rarity values, e.g. card.Rarity.Common.
var Rarity = rarities{
	Common:    engine.Common,
	Uncommon:  engine.Uncommon,
	Rare:      engine.Rare,
	Special:   engine.Special,
	Connected: engine.Connected,
}

type rarities struct {
	// Common is the common rarity.
	Common engine.Rarity
	// Uncommon is the uncommon rarity.
	Uncommon engine.Rarity
	// Rare is the rare rarity.
	Rare engine.Rarity
	// Special is the special rarity.
	Special engine.Rarity
	// Connected is the rarity of a card that only enters a deck through another
	// card's connection (see card.Connects).
	Connected engine.Rarity
}

// Traits groups the trait values, e.g. card.Traits.Beast. Named plural, unlike
// House/Type/Keyword/Rarity, because the singular Trait already names the value
// type above.
var Traits = traits{
	Agent:     engine.Agent,
	Ally:      engine.Ally,
	Angel:     engine.Angel,
	Beast:     engine.Beast,
	Cleric:    engine.Cleric,
	Cyborg:    engine.Cyborg,
	Demon:     engine.Demon,
	Dragon:    engine.Dragon,
	Elf:       engine.Elf,
	Equation:  engine.Equation,
	Faerie:    engine.Faerie,
	Fungus:    engine.Fungus,
	Giant:     engine.Giant,
	Goblin:    engine.Goblin,
	Horseman:  engine.Horseman,
	Human:     engine.Human,
	Imp:       engine.Imp,
	Insect:    engine.Insect,
	Item:      engine.Item,
	Knight:    engine.Knight,
	Law:       engine.Law,
	Location:  engine.Location,
	Martian:   engine.Martian,
	Merchant:  engine.Merchant,
	Monk:      engine.Monk,
	Mutant:    engine.Mutant,
	Niffle:    engine.Niffle,
	Power:     engine.Power,
	Priest:    engine.Priest,
	Quest:     engine.Quest,
	Ranger:    engine.Ranger,
	Rat:       engine.Rat,
	Robot:     engine.Robot,
	Scientist: engine.Scientist,
	Shard:     engine.Shard,
	Soldier:   engine.Soldier,
	Specter:   engine.Specter,
	Spirit:    engine.Spirit,
	Thief:     engine.Thief,
	Vehicle:   engine.Vehicle,
	Weapon:    engine.Weapon,
	Witch:     engine.Witch,
}

// traits backs the Traits namespace. A trait carries no rules meaning of its
// own (unlike a keyword), so a per-field comment here would only restate the
// field's name.
type traits struct {
	Agent,
	Ally,
	Angel,
	Beast,
	Cleric,
	Cyborg,
	Demon,
	Dragon,
	Elf,
	Equation,
	Faerie,
	Fungus,
	Giant,
	Goblin,
	Horseman,
	Human,
	Imp,
	Insect,
	Item,
	Knight,
	Law,
	Location,
	Martian,
	Merchant,
	Monk,
	Mutant,
	Niffle,
	Power,
	Priest,
	Quest,
	Ranger,
	Rat,
	Robot,
	Scientist,
	Shard,
	Soldier,
	Specter,
	Spirit,
	Thief,
	Vehicle,
	Weapon,
	Witch engine.Trait
}

// Keyword groups the keyword values, e.g. card.Keyword.Skirmish.
var Keyword = keywords{
	Skirmish:  engine.Skirmish,
	Poison:    engine.Poison,
	Elusive:   engine.Elusive,
	Taunt:     engine.Taunt,
	Versatile: engine.Versatile,
	Alpha:     engine.Alpha,
	Omega:     engine.Omega,
	Deploy:    engine.Deploy,
}

type keywords struct {
	// Skirmish: this creature deals no retaliation damage when it fights.
	Skirmish engine.Keyword
	// Poison: any damage this creature deals to a creature destroys it.
	Poison engine.Keyword
	// Elusive: the first time this creature is attacked each turn, no damage is dealt.
	Elusive engine.Keyword
	// Taunt: neighboring non-Taunt creatures cannot be attacked or fought.
	Taunt engine.Keyword
	// Versatile: this card may be played from any house (its Action: is an Omni).
	Versatile engine.Keyword
	// Alpha: this card can only be played as the first card its player plays,
	// uses, or discards on their turn.
	Alpha engine.Keyword
	// Omega: after this card is played, the current step of the turn ends — no
	// more cards this step except through pending abilities still resolving.
	Omega engine.Keyword
	// Deploy: this creature may enter play at any position in its controller's
	// battleline, not only on a flank.
	Deploy engine.Keyword
}

// Keywords builds the keyword slice for an upgrade's granted keywords, e.g.
// card.StaticModifier{Keywords: card.Keywords(card.Keyword.Skirmish)}. It exists
// because card.Keyword is the value namespace, so a []card.Keyword literal can't
// be written directly.
func Keywords(k ...engine.Keyword) []engine.Keyword { return k }

// UseKind groups the ways a card in play can be used, e.g. card.UseKind.Reap.
var UseKind = useKinds{
	Reap:   engine.ReapUse,
	Fight:  engine.FightUse,
	Action: engine.ActionUse,
}

type useKinds struct {
	// Reap is using a creature to reap.
	Reap engine.UseKind
	// Fight is using a creature to fight.
	Fight engine.UseKind
	// Action is using a card's "Action:" ability.
	Action engine.UseKind
}

// Types builds the card-type slice for a Types filter, e.g.
// card.DiscardFromHand{Types: card.Types(card.Type.Creature)}. Like Keywords, it
// exists because card.Type is the value namespace, so a []card.CardType literal
// can't be written directly.
func Types(t ...engine.CardType) []engine.CardType { return t }

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
	StartOfTurn:                 engine.TriggerStartOfTurn,
	AfterChooseHouse:            engine.TriggerAfterChooseHouse,
	AfterEnemyCreatureDestroyed: engine.TriggerAfterEnemyCreatureDestroyed,
	AfterEnemyCardPlayed:        engine.TriggerAfterEnemyCardPlayed,
	AfterUse:                    engine.TriggerAfterUse,
	AfterDiscardFromHand:        engine.TriggerAfterDiscardFromHand,
	UsedSelf:                    engine.TriggerAfterUsedSelf,
	AfterCreatureReaps:          engine.TriggerAfterCreatureReaps,
	AfterEnemyCreatureReaps:     engine.TriggerAfterEnemyCreatureReaps,
	LeavesPlay:                  engine.TriggerLeavesPlay,
}

type triggers struct {
	// Play fires when the card is played ("Play:").
	Play engine.Trigger
	// Reap fires after this creature reaps ("Reap:").
	Reap engine.Trigger
	// Fight fires after this creature fights ("Fight:").
	Fight engine.Trigger
	// BeforeFight fires before this creature's fight resolves ("Before Fight:").
	BeforeFight engine.Trigger
	// Action is an ability the controller activates on their turn ("Action:").
	Action engine.Trigger
	// AfterForgeKey fires after the controller forges a key.
	AfterForgeKey engine.Trigger
	// AfterCreatureEnters fires after another creature enters play.
	AfterCreatureEnters engine.Trigger
	// Destroyed fires when this creature is destroyed ("Destroyed:").
	Destroyed engine.Trigger
	// AfterDestroyedFighting fires when a creature is destroyed fighting this one.
	AfterDestroyedFighting engine.Trigger
	// AfterCardPlayed fires after the controller plays a card.
	AfterCardPlayed engine.Trigger
	// EndOfTurn fires at the end of the controller's turn.
	EndOfTurn engine.Trigger
	// StartOfTurn fires at the start of the controller's turn, before they forge.
	StartOfTurn engine.Trigger
	// AfterChooseHouse fires after the controller chooses their active house.
	AfterChooseHouse engine.Trigger
	// AfterEnemyCreatureDestroyed fires after an enemy creature is destroyed during your turn.
	AfterEnemyCreatureDestroyed engine.Trigger
	// AfterEnemyCardPlayed fires after the opponent plays a card.
	AfterEnemyCardPlayed engine.Trigger
	// AfterUse fires after the controller uses a card (reap, fight, or Action:).
	AfterUse engine.Trigger
	// AfterDiscardFromHand fires after the controller discards a card from hand.
	AfterDiscardFromHand engine.Trigger
	// UsedSelf fires after this creature is itself used (reap, fight, or Action:),
	// so an upgrade can punish its own host ("After this creature is used, ...").
	UsedSelf engine.Trigger
	// AfterCreatureReaps fires after any creature reaps (friendly or enemy), with
	// the reaper as "it" (Orb of Invidius stuns whatever just reaped).
	AfterCreatureReaps engine.Trigger
	// AfterEnemyCreatureReaps fires after an enemy creature reaps, with the reaper
	// as "it" (Pip Pip stuns the enemy that just reaped).
	AfterEnemyCreatureReaps engine.Trigger
	// LeavesPlay fires as this card leaves play by any route ("Leaves Play:").
	LeavesPlay engine.Trigger
}

// Controller and Opponent are the two players an effect can target, relative to
// the card's controller: card.Controller (the player who controls the card) or
// card.Opponent (their opponent). card.EachPlayer means both, for effects that
// reach everyone at once (e.g. a key-cost change on each player's keys).
var (
	// Controller is the player who controls the card whose ability is resolving.
	Controller = engine.Controller
	// Opponent is the controller's opponent.
	Opponent = engine.Opponent
	// EachPlayer means both players, for effects that reach everyone at once.
	EachPlayer = engine.EachPlayer
	// ItsOwner is the owner of the card in context (ctx.It).
	ItsOwner = engine.ItsOwner
	// ItsOpponent is the opponent of the card in context — for a capture, of the
	// capturing creature, so each side's creatures draw from a different pool.
	ItsOpponent = engine.ItsOpponent
)
