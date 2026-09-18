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
	// of an effect that filters by trait (card.Filter.Trait).
	Trait = engine.Trait
	// KeywordValue is the value type for a card.Keyword.X constant, e.g. as an
	// element of card.AttackKeywords.Keywords.
	KeywordValue = engine.Keyword
	// BonusIcon is the value type for a card.Bonus.X constant, an element of
	// card.WithBonus / card.WithEnhance.
	BonusIcon = engine.BonusIcon
	// BonusInstead is the continuous bonus-icon substitution a card offers through
	// card.WithBonusInstead (Amphora Captura, Scrivener Favian).
	BonusInstead = engine.BonusInstead
	// Player is the relative player an effect targets: card.Controller or card.Opponent.
	Player = engine.Player
)

// Bonus groups the bonus-icon kinds, e.g. card.Bonus.Aember.
var Bonus = bonusIcons{
	Aember:  engine.BonusAember,
	Capture: engine.BonusCapture,
	Damage:  engine.BonusDamage,
	Draw:    engine.BonusDraw,
}

type bonusIcons struct {
	// Aember gains 1 Æmber when the card is played.
	Aember engine.BonusIcon
	// Capture has a friendly creature capture 1 Æmber from the opponent.
	Capture engine.BonusIcon
	// Damage deals 1 damage to a creature in play.
	Damage engine.BonusIcon
	// Draw draws 1 card.
	Draw engine.BonusIcon
}

// House groups the faction values, e.g. card.House.Brobnar.
var House = houses{
	None:         engine.HouseNone,
	Brobnar:      engine.Brobnar,
	Dis:          engine.Dis,
	Logos:        engine.Logos,
	Mars:         engine.Mars,
	Sanctum:      engine.Sanctum,
	Saurian:      engine.Saurian,
	Shadows:      engine.Shadows,
	StarAlliance: engine.StarAlliance,
	Untamed:      engine.Untamed,
	Self:         engine.SelfHouse,
}

type houses struct {
	// None is no house at all, for a houseless card (a Special stamped with its
	// pod's house at deck generation). Pair it with card.Houseless.
	None engine.House
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
	// Saurian is the house of dinosaurs and Roman legionnaires.
	Saurian engine.House
	// StarAlliance is the house of the united starship crews.
	StarAlliance engine.House
	// Untamed is the house of nature and beasts.
	Untamed engine.House
	// Self is the card's own house, filled in when the card is built. Use it
	// whenever an ability names the house the card itself belongs to — Battle
	// Fleet reveals card.House.Self cards — so the two can never drift apart.
	// A card that names a *different* house names it outright.
	Self engine.House
}

// KeyColor groups the forged-key colours, e.g. card.KeyColor.Red.
var KeyColor = keyColors{
	Red:    engine.KeyColorRed,
	Blue:   engine.KeyColorBlue,
	Yellow: engine.KeyColorYellow,
}

type keyColors struct {
	// Red is the red key.
	Red engine.KeyColor
	// Blue is the blue key.
	Blue engine.KeyColor
	// Yellow is the yellow key.
	Yellow engine.KeyColor
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

// ItNoun groups the nouns a condition can print instead of saying "it", e.g.
// card.ItNoun.DiscardedCard. It changes only the wording, never which card the
// condition reads — that is card.Subject.
var ItNoun = itNouns{
	DiscardedCard:  engine.DiscardedCard,
	ThatCard:       engine.ThatCard,
	FoughtCreature: engine.FoughtCreature,
}

type itNouns struct {
	// DiscardedCard names the card an effect just discarded.
	DiscardedCard engine.ItNoun
	// ThatCard names the card an effect just acted on when "it" would be ambiguous.
	ThatCard engine.ItNoun
	// FoughtCreature names the creature the source is fighting.
	FoughtCreature engine.ItNoun
}

// Subject groups the cards a condition can read, e.g. card.Subject.This. Unlike
// card.ItNoun it changes the question's referent, not its wording.
var Subject = subjects{
	It:   engine.It,
	This: engine.This,
}

type subjects struct {
	// It is the default: the card a trigger or a preceding effect put in context.
	It engine.Subject
	// This is the card the ability is printed on.
	This engine.Subject
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
	// Connected is the rarity of a card that only enters a deck through a cluster
	// (see card.InCluster).
	Connected engine.Rarity
}

// Traits groups the trait values, e.g. card.Traits.Beast. Named plural, unlike
// House/Type/Keyword/Rarity, because the singular Trait already names the value
// type above.
var Traits = traits{
	Agent:        engine.Agent,
	Ai:           engine.Ai,
	Alien:        engine.Alien,
	Ally:         engine.Ally,
	Angel:        engine.Angel,
	Aquan:        engine.Aquan,
	Assassin:     engine.Assassin,
	Beast:        engine.Beast,
	Cat:          engine.Cat,
	Changeling:   engine.Changeling,
	Cleric:       engine.Cleric,
	Cyborg:       engine.Cyborg,
	Demon:        engine.Demon,
	Dinosaur:     engine.Dinosaur,
	Dragon:       engine.Dragon,
	Egg:          engine.Egg,
	Elf:          engine.Elf,
	Equation:     engine.Equation,
	Experiment:   engine.Experiment,
	Faerie:       engine.Faerie,
	Fungus:       engine.Fungus,
	Giant:        engine.Giant,
	Goblin:       engine.Goblin,
	Handuhan:     engine.Handuhan,
	Horseman:     engine.Horseman,
	Human:        engine.Human,
	Hunter:       engine.Hunter,
	Imp:          engine.Imp,
	Insect:       engine.Insect,
	Item:         engine.Item,
	Jelly:        engine.Jelly,
	Knight:       engine.Knight,
	Krxix:        engine.Krxix,
	Law:          engine.Law,
	Leader:       engine.Leader,
	Location:     engine.Location,
	Martian:      engine.Martian,
	Merchant:     engine.Merchant,
	Monk:         engine.Monk,
	Mutant:       engine.Mutant,
	Niffle:       engine.Niffle,
	Philosopher:  engine.Philosopher,
	Pilot:        engine.Pilot,
	Pirate:       engine.Pirate,
	Politician:   engine.Politician,
	Power:        engine.Power,
	Priest:       engine.Priest,
	Proximan:     engine.Proximan,
	Psion:        engine.Psion,
	Quest:        engine.Quest,
	Ranger:       engine.Ranger,
	Rat:          engine.Rat,
	Redacted:     engine.Redacted,
	Robot:        engine.Robot,
	Scientist:    engine.Scientist,
	Shapeshifter: engine.Shapeshifter,
	Shard:        engine.Shard,
	Sin:          engine.Sin,
	Soldier:      engine.Soldier,
	Specter:      engine.Specter,
	Spirit:       engine.Spirit,
	Thief:        engine.Thief,
	Tree:         engine.Tree,
	Vehicle:      engine.Vehicle,
	Weapon:       engine.Weapon,
	Witch:        engine.Witch,
	Wolf:         engine.Wolf,
}

// traits backs the Traits namespace. A trait carries no rules meaning of its
// own (unlike a keyword), so a per-field comment here would only restate the
// field's name.
type traits struct {
	Agent,
	Ally,
	Angel,
	Beast,
	Cat,
	Changeling,
	Cleric,
	Cyborg,
	Demon,
	Dragon,
	Egg,
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
	Redacted,
	Robot,
	Scientist,
	Shard,
	Sin,
	Soldier,
	Specter,
	Spirit,
	Thief,
	Tree,
	Vehicle,
	Weapon,
	Witch,
	Ai,
	Alien,
	Aquan,
	Dinosaur,
	Experiment,
	Handuhan,
	Hunter,
	Jelly,
	Krxix,
	Leader,
	Philosopher,
	Pilot,
	Pirate,
	Politician,
	Proximan,
	Psion,
	Shapeshifter,
	Wolf,
	Assassin engine.Trait
}

// Keyword groups the keyword values, e.g. card.Keyword.Skirmish.
var Keyword = keywords{
	Skirmish:     engine.Skirmish,
	Poison:       engine.Poison,
	Elusive:      engine.Elusive,
	Taunt:        engine.Taunt,
	Versatile:    engine.Versatile,
	Alpha:        engine.Alpha,
	Omega:        engine.Omega,
	Deploy:       engine.Deploy,
	Treachery:    engine.Treachery,
	Invulnerable: engine.Invulnerable,
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
	// Treachery: this card enters play under your opponent's control.
	Treachery engine.Keyword
	// Invulnerable: this creature cannot be destroyed or dealt damage.
	Invulnerable engine.Keyword
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

// UseKinds builds the use-kind slice for a ConstantAbility.CannotBeUsedTo field,
// because card.UseKind is the value namespace, so a []card.UseKind literal can't
// be written directly. E.g. card.UseKinds(card.UseKind.Reap).
func UseKinds(k ...engine.UseKind) []engine.UseKind { return k }

// Triggers builds the trigger slice for a ConstantAbility.DisableTriggers field,
// because card.Trigger is the value namespace, so a []card.Trigger literal can't
// be written directly. E.g. card.Triggers(card.Trigger.Destroyed).
func Triggers(t ...engine.Trigger) []engine.Trigger { return t }

// SpendScope names which player a spend-as-pool permission benefits, e.g.
// card.SpendScope.Controller (The Callipygian Ideal) or card.SpendScope.Opponent
// (Mole). It has no default: a card that grants the permission names one.
var SpendScope = spendScopes{
	Controller: engine.SpendByController,
	Opponent:   engine.SpendByOpponent,
}

type spendScopes struct {
	Controller engine.SpendScope
	Opponent   engine.SpendScope
}

// Trigger groups the ability triggers, e.g. card.Trigger.Play or
// card.Trigger.AfterForgeKey.
var Trigger = triggers{
	Action:                         engine.TriggerAction,
	AfterAemberStolenFromYou:       engine.TriggerAfterAemberStolenFromYou,
	AfterAnyPlayerChoosesHouse:     engine.TriggerAfterAnyPlayerChoosesHouse,
	AfterAnyPlayerStartOfTurn:      engine.TriggerAfterAnyPlayerStartOfTurn,
	AfterAnyPlayerEndOfTurn:        engine.TriggerAfterAnyPlayerEndOfTurn,
	AfterArmorPrevents:             engine.TriggerAfterArmorPrevents,
	AfterAssaultDestroys:           engine.TriggerAfterAssaultDestroys,
	AfterBonusDamage:               engine.TriggerAfterBonusDamage,
	AfterBonusDraw:                 engine.TriggerAfterBonusDraw,
	AfterCardPlayed:                engine.TriggerAfterCardPlayed,
	AfterChooseHouse:               engine.TriggerAfterChooseHouse,
	AfterCreatureDestroyed:         engine.TriggerAfterCreatureDestroyed,
	AfterCreatureEnters:            engine.TriggerAfterCreatureEnters,
	AfterCreatureFights:            engine.TriggerAfterCreatureFights,
	AfterCreaturePlayed:            engine.TriggerAfterCreaturePlayed,
	AfterCreaturePlayedAdjacent:    engine.TriggerAfterCreaturePlayedAdjacent,
	AfterCreatureReaps:             engine.TriggerAfterCreatureReaps,
	AfterDestroyedFighting:         engine.TriggerAfterDestroyedFighting,
	AfterDiscardFromHand:           engine.TriggerAfterDiscardFromHand,
	AfterEnemyCardPlayed:           engine.TriggerAfterEnemyCardPlayed,
	AfterEnemyDestroyedFighting:    engine.TriggerAfterEnemyDestroyedFighting,
	AfterForgeKey:                  engine.TriggerAfterForgeKey,
	AfterNeighborFights:            engine.TriggerAfterNeighborFights,
	AfterOpponentForgesKey:         engine.TriggerAfterOpponentForgesKey,
	AfterPlayerForgesKey:           engine.TriggerAfterPlayerForgesKey,
	AfterTacticPlayedBeforeResolve: engine.TriggerAfterTacticPlayedBeforeResolve,
	AfterUse:                       engine.TriggerAfterUse,
	AfterUpgradeEnters:             engine.TriggerAfterUpgradeEnters,
	BeforeFight:                    engine.TriggerBeforeFight,
	BeforeOpponentForgesKey:        engine.TriggerBeforeOpponentForgesKey,
	Destroyed:                      engine.TriggerDestroyed,
	EndOfReadyStep:                 engine.TriggerEndOfReadyStep,
	EndOfTurn:                      engine.TriggerEndOfTurn,
	Fight:                          engine.TriggerAfterFight,
	FightReap:                      triggerFightReap,
	LeavesPlay:                     engine.TriggerLeavesPlay,
	Play:                           engine.TriggerAfterPlay,
	PlayFight:                      triggerPlayFight,
	PlayFightReap:                  triggerPlayFightReap,
	PlayReap:                       triggerPlayReap,
	Reap:                           engine.TriggerAfterReap,
	StartOfTurn:                    engine.TriggerStartOfTurn,
	UsedSelf:                       engine.TriggerAfterUsedSelf,
}

// Composite triggers are facade-only fan-out markers: card.WithAbility expands
// each into the atomic Play/Fight/Reap abilities, so the engine runtime only ever
// sees atomic triggers (ADR 0006), and text rendering merges the atomic abilities
// back into one "Play/Fight/Reap:" line. Their negative values can never collide
// with the engine's non-negative Trigger constants.
const (
	triggerPlayFightReap engine.Trigger = -1 - iota
	triggerFightReap
	triggerPlayReap
	triggerPlayFight
)

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
	// AfterCreaturePlayedAdjacent fires after a creature is played adjacent to this card.
	AfterCreaturePlayedAdjacent engine.Trigger
	// AfterNeighborFights fires after a battleline neighbor of this card is used to fight.
	AfterNeighborFights engine.Trigger
	// AfterBonusDamage fires after you resolve a Damage bonus icon, with the creature it hit as "it".
	AfterBonusDamage engine.Trigger
	// AfterBonusDraw fires after you resolve a Draw bonus icon (only when a card was drawn).
	AfterBonusDraw engine.Trigger
	// Destroyed fires when this creature is destroyed ("Destroyed:").
	Destroyed engine.Trigger
	// AfterDestroyedFighting fires when a creature is destroyed in a fight with this one.
	AfterDestroyedFighting engine.Trigger
	// AfterEnemyDestroyedFighting fires on a bystander when an enemy creature is destroyed while fighting.
	AfterEnemyDestroyedFighting engine.Trigger
	// AfterAssaultDestroys fires when this creature's Assault damage destroys the creature it attacks.
	AfterAssaultDestroys engine.Trigger
	// AfterArmorPrevents fires after this card prevents damage with its own armor.
	AfterArmorPrevents engine.Trigger
	// AfterCardPlayed fires after the controller plays a card.
	AfterCardPlayed engine.Trigger
	// EndOfTurn fires at the end of the controller's turn.
	EndOfTurn engine.Trigger
	// StartOfTurn fires at the start of the controller's turn, before they forge.
	StartOfTurn engine.Trigger
	// EndOfReadyStep fires at the end of the controller's "ready cards" step, after
	// every card has readied (Greater Oxtet).
	EndOfReadyStep engine.Trigger
	// AfterChooseHouse fires after the controller chooses their active house.
	AfterChooseHouse engine.Trigger
	// AfterAnyPlayerChoosesHouse fires after either player chooses their active
	// house, whoever's turn it is (Snag's Mirror).
	AfterAnyPlayerChoosesHouse engine.Trigger
	// AfterCreatureDestroyed fires after any creature is destroyed, with the
	// destroyed creature as "it"; it fires only for cards that survive the batch.
	AfterCreatureDestroyed engine.Trigger
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
	// the reaper as "it" (Orb of Invidius stuns whatever just reaped). A reaction
	// narrowed to an enemy reaper (Pip Pip) is this trigger gated on ItIsEnemy.
	AfterCreatureReaps engine.Trigger
	// AfterCreatureFights fires after any creature is used to fight (friendly or
	// enemy), with the fighting creature as "it" (Shattered Throne makes it capture).
	AfterCreatureFights engine.Trigger
	// AfterPlayerForgesKey fires after any player forges a key, acting on the
	// player who forged (Forgemaster Og).
	AfterPlayerForgesKey engine.Trigger
	// AfterOpponentForgesKey fires after the opponent forges a key, on the
	// non-forging player's cards (Forge Compiler).
	AfterOpponentForgesKey engine.Trigger
	// BeforeOpponentForgesKey fires when the opponent would forge a key, before the
	// forge, so the ability can cancel it (Keyforgery). The forging opponent is
	// "they"; a cancelled forge leaves their Æmber unspent.
	BeforeOpponentForgesKey engine.Trigger
	// AfterCreaturePlayed fires after any creature is played from hand (friendly or
	// enemy), with the played creature as "it" (The Big One).
	AfterCreaturePlayed engine.Trigger
	// AfterUpgradeEnters fires after any upgrade enters play (friendly or enemy),
	// with the entering upgrade as "it" (Armory Officer Nel).
	AfterUpgradeEnters engine.Trigger
	// AfterAemberStolenFromYou fires after Æmber is stolen from this card's
	// controller, with the number stolen in that theft available as a count
	// (Molephin).
	AfterAemberStolenFromYou engine.Trigger
	// AfterAnyPlayerStartOfTurn fires at the start of every player's turn — its own
	// controller's and the opponent's — resolving as the player whose turn is
	// starting, so "they"/"that player" is that active player (Gambling Den, General
	// Order 24).
	AfterAnyPlayerStartOfTurn engine.Trigger
	// AfterAnyPlayerEndOfTurn fires at the end of every player's turn — its own
	// controller's and the opponent's — resolving as the player whose turn is
	// ending, so "they"/"that player" is that active player (Pincerator).
	AfterAnyPlayerEndOfTurn engine.Trigger
	// LeavesPlay fires as this card leaves play by any route ("Leaves Play:").
	LeavesPlay engine.Trigger
	// AfterTacticPlayedBeforeResolve fires after a Tactic is played, by either
	// player, before that Tactic's own effect resolves.
	AfterTacticPlayedBeforeResolve engine.Trigger
	// PlayFightReap fires the effect as a Play, a Fight, and a Reap ability, printed
	// as one "Play/Fight/Reap:" line. It is a composite: card.WithAbility fans it
	// into the three atomic abilities.
	PlayFightReap engine.Trigger
	// FightReap fires the effect as both a Fight and a Reap ability, printed as
	// one "Fight/Reap:" line. It is a composite fanned out by card.WithAbility.
	FightReap engine.Trigger
	// PlayReap fires the effect as both a Play and a Reap ability, printed as one
	// "Play/Reap:" line. It is a composite fanned out by card.WithAbility.
	PlayReap engine.Trigger
	// PlayFight fires the effect as both a Play and a Fight ability, printed as one
	// "Play/Fight:" line. It is a composite fanned out by card.WithAbility.
	PlayFight engine.Trigger
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
	// ItsController is the player who controlled the card in context when a
	// preceding effect touched it — the "its controller" referent for a destroyed
	// creature, captured before it leaves play so a stolen creature pays its
	// controller, not its owner (Saury About That).
	ItsController = engine.ItsController
	// ThatPlayer is the player named by the ability's trigger — for a cross-player
	// reaction, whoever caused it (Forgemaster Og drains the player who forged).
	ThatPlayer = engine.ThatPlayer
	// ChosenPlayer is a player the controller chooses at resolution — used where a
	// zone-movement effect acts on "a discard pile" the controller picks (Creeping
	// Oblivion purges up to 2 cards from a discard pile).
	ChosenPlayer = engine.ChosenPlayer
)
