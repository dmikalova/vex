package engine

import "fmt"

// CardDefinition is the immutable blueprint for a card. Definitions are shared,
// read-only data held in the match Catalog; all mutable per-match state lives in
// the flat GameState. This split keeps the runtime state a plain value that can
// be copied cheaply (see GameState.FastCopy).
type CardDefinition struct {
	// Identity printed on the card.
	Name   string
	House  House
	Type   CardType
	Rarity Rarity
	Traits []Trait

	// Creature stats. Zero for non-creatures.
	Power int
	Armor int

	// Keywords are the keywords printed on the card.
	Keywords []Keyword

	// A creature with Assault N deals N damage to the creature it attacks,
	// immediately before combat damage is dealt. Zero means the creature does not
	// have Assault.
	Assault int
	// A creature with Hazardous N deals N damage to any creature that attacks it,
	// before that attacker deals its combat damage. Zero means the creature does
	// not have Hazardous.
	Hazardous int
	// A creature with Splash-attack N deals N damage to each neighbor of the
	// creature it fights, simultaneously with its own fight damage. Zero means the
	// creature does not have Splash-attack.
	SplashAttack int
	// AttackDamage customizes the damage this creature deals when it fights, for the
	// few creatures whose fight damage is not simply their current power (Valdr's
	// flank bonus, Ether Spider dealing none). The zero value deals its power.
	AttackDamage AttackDamage

	// DealsNoDamageWhenAttacked stops this creature from dealing its retaliation
	// damage to an attacker when it is fought (Lollop the Titanic). It leaves the
	// creature's own fight damage untouched.
	DealsNoDamageWhenAttacked bool

	// GrantsEntersReady, when set to Creature or Artifact, makes friendly cards of
	// that type enter play ready instead of exhausted while this card is in play —
	// Duskwitch readies your creatures, The Curator readies your artifacts.
	// TypeUnset grants nothing.
	GrantsEntersReady CardType

	// FightRestriction, when set, limits which enemy creatures this creature may
	// fight to those its Target allows (Bigtwig can only fight stunned creatures).
	// The zero Target imposes no restriction.
	FightRestriction Target

	// CannotBeUsedTo are the ways this card may not be used — Tireless Crocag cannot
	// reap, but still fights. It bars the card itself, unlike the player-wide
	// restrictions in Restricts. Empty imposes no restriction.
	CannotBeUsedTo []UseKind

	// CannotBeUsedWhile, when set, bars the card from being used at all — reaped,
	// fought with, or used for its Action ability — for as long as the condition
	// holds for its controller (Valoocanth cannot be used while the tide is low).
	// It is the conditional, all-ways counterpart to CannotBeUsedTo. The zero value
	// (nil condition) imposes no restriction.
	CannotBeUsedWhile Condition

	// DestroyedWhen, when set, is a condition that puts this creature in a
	// destroyable state for as long as it holds — Tireless Crocag dies while its
	// controller's opponent has no creatures. It is read board-wide every time
	// destruction settles, so it takes effect the moment the board makes it true.
	DestroyedWhen Condition

	// TakesDamageFor, when set, names the creatures whose damage this card takes
	// instead — Shadow Self takes the damage dealt to its non-Specter neighbors.
	// It is read wherever damage lands, so it covers fight damage and effect damage
	// alike. The zero Target shields nobody.
	TakesDamageFor Target

	// AttackIgnores are the defensive keywords this creature ignores while it is
	// attacking — Niffle Ape ignores taunt and elusive, so it may be used to fight a
	// creature its neighbors' taunt shields and no elusive stops the damage.
	AttackIgnores []Keyword

	// AttackKeywords are keywords this creature gains only while attacking — Spyyyder
	// gains poison while attacking a flank creature. The zero value grants none.
	AttackKeywords AttackKeywords

	// AemberBonus is the number of Æmber "pips" printed on the card; the
	// controller gains this much Æmber when the card is played.
	AemberBonus int

	// Static is a continuous modifier an Upgrade applies to its host creature. A
	// creature with PlayableAsUpgrade set also carries what it grants a host here.
	Static StaticModifier

	// PlayableAsUpgrade lets this creature be played as an upgrade instead of a
	// creature, attaching to a host that then gains its Static modifier (Explo-rover
	// grants skirmish, CALV-1N grants a Fight/Reap draw). The card is authored twice
	// over — its own creature keywords and abilities are its identity in the
	// battleline, and Static is what it grants a host as an upgrade. Only a Creature
	// may set it, and its Static must grant something (validated in NewCard).
	PlayableAsUpgrade bool

	// ConstantAbilities are constant abilities this card applies to creatures in play for as
	// long as it stays in play (see Game.constantBonus). Unlike Static (an Upgrade
	// buffing its own host), constant abilities reach whole Targets of creatures.
	ConstantAbilities []ConstantAbility

	// Restricts holds the continuous "cannot" rules the card imposes on its
	// controller while it is in play (e.g. Grommid's "You cannot play creatures").
	Restricts Restrictions

	// CannotPlayWhile is a symmetric, board-wide play bar the card imposes while in
	// play: any player who meets its condition cannot play cards of its type,
	// whoever controls the card (Quixxle Stone). Unlike Restricts.CannotPlay it is
	// evaluated per attempting player, so it bars whichever side the condition
	// names. The zero value (nil condition) imposes no restriction.
	CannotPlayWhile ConditionalPlayBar

	// KeyCostChanges are the continuous changes this card, while in play, makes to
	// key cost — who each affects and by how much (e.g. Grabber Jammer's "Your
	// opponent's keys cost +1 Æmber"). A card may impose several (Grump Buggy raises
	// each player's keys by a different count). Empty changes nothing.
	KeyCostChanges []KeyCostChange

	// HouseLock is a continuous constraint this card, while in play, puts on a
	// player's active-house choice — Pitlord's "you must choose Dis", Restringuntus'
	// bar on the house it named. The zero value constrains nobody.
	HouseLock HouseLock

	// PlayPermission is a continuous grant this card makes while in play, letting
	// its controller play cards of a house on turns where that house is not their
	// active house (Witch of the Wilds). The zero value grants nothing.
	PlayPermission PlayPermission

	// Replaces is a continuous replacement this card applies to a game event's
	// outcome while it is in play, on either end of an Æmber flow: Ether Spider
	// replaces Æmber being added to its opponent's pool (EventAemberAddedToPool, the
	// destination) with capturing it, and Po's Pixies replaces the source of a steal
	// or capture from its own pool (EventAemberTakenFromPool) with the common supply.
	// The zero value carries no replacement. (An Upgrade grants a replacement to its
	// host through StaticModifier.Replaces instead.)
	Replaces Instead

	// DrawModifier is a continuous change this card makes to a player's end-of-turn
	// hand-refill size while in play — Mother refills its controller to one more
	// card, Succubus refills the opponent to one fewer. The zero value changes
	// nothing.
	DrawModifier DrawModifier

	// AemberCannotBeStolen, while the card is in play, makes its controller's Æmber
	// impossible for the opponent to steal (The Vaultkeeper).
	AemberCannotBeStolen bool

	// AemberCannotBeStolenWhileItHasAember protects its controller's Æmber from
	// theft only while the card itself has Æmber on it (Odoac the Patrician).
	AemberCannotBeStolenWhileItHasAember bool

	// SpendableAember lets the Æmber sitting on this card be put toward a key,
	// so it is a private vault its controller can bank into (Safe Place).
	SpendableAember bool

	// GainsForgeAember gives this card's controller all the Æmber their opponent
	// spends forging a key, for as long as it stays in play (The Sting).
	GainsForgeAember bool

	// GuardsOpponentForge lets this card, while in play, interrupt the opponent's
	// key forges: when the opponent would forge a key, they name a house, a random
	// card is revealed from this card's controller's hand, and if that card is not
	// of the named house the card is destroyed and the forge is prevented
	// (Keyforgery).
	GuardsOpponentForge bool

	// PlayRequirement is the Æmber the controller must have — and, when the
	// requirement spends, gives up — to play this card from hand.
	PlayRequirement PlayRequirement

	// Abilities are the triggered abilities on the card.
	Abilities []Ability
}

// DrawModifier is a continuous change a card in play makes to how many cards a
// player draws back up to during their "draw cards" phase. Player is relative to
// the card's controller (Controller, Opponent, or EachPlayer), and Amount is added
// to the normal hand size (+1 for Mother and The Howling Pit, -1 for Succubus).
type DrawModifier struct {
	Player Player
	Amount int
	// OnlyWhileOffFlank restricts the modifier to while the source card is not on a
	// flank of its battleline (Streke).
	OnlyWhileOffFlank bool
	// OnlyWhileInCenter restricts the modifier to while the source card is in the
	// center of its battleline (Zenzizenzizenzic).
	OnlyWhileInCenter bool
}

// affects reports whether a draw modifier owned by owner applies to target's draw.
func (m DrawModifier) affects(owner, target int) bool {
	switch m.Player {
	case Controller:
		return target == owner
	case Opponent:
		return target != owner
	default: // EachPlayer
		return true
	}
}

// Restrictions are the continuous "cannot" rules a card imposes while it stays in
// play. They are consulted by the matching gate (Game.cannotFight,
// Game.cannotPlayCreatures, Game.cannotPlayCard) alongside any timed restriction.
type Restrictions struct {
	// Fighting bars the controller from using creatures to fight.
	Fighting bool
	// Reaping bars a player's creatures from reaping while this card stays in play,
	// relative to the card's controller (Barrister Joya's Opponent bars the enemy's
	// creatures, "Enemy creatures cannot reap."). Its zero value (playerUnset)
	// imposes no reaping restriction.
	Reaping Player
	// CannotPlay bars the controller from playing cards of this type (e.g. Creature
	// for Grommid's "You cannot play creatures"). The zero value (an unset CardType)
	// imposes no play restriction.
	CannotPlay CardType
	// PlayCardLimit caps cards a relative player may play each turn (Ember Imp's
	// "your opponent cannot play more than 2 cards each turn"). Its zero value
	// imposes no limit.
	PlayCardLimit PlayCardLimit
	// Toll is Æmber the controller's opponent must pay the controller to play or
	// use an artifact (Customs Office, Tentacus). Its zero value imposes no toll.
	Toll Toll
	// UseCondition is a Condition that must be met for the controller to use this card (Giant Sloth).
	UseCondition Condition
	// SkipForge bars the controller from forging a key during their "forge a
	// key" step (The Sting).
	SkipForge bool
	// NoForgeKeyNumber bars every player from forging the key of this ordinal
	// (1 = first, 2 = second, 3 = third) while this card stays in play — the Key
	// Imps' "Players cannot forge their first key." Its zero value bars nothing.
	NoForgeKeyNumber int
	// NoForgeWhileAheadOnKeys bars every player from forging while they have more
	// forged keys than their opponent, whoever controls this card (Heart of the
	// Forest keeps the leader from pulling further ahead).
	NoForgeWhileAheadOnKeys bool
	// MustFightIfAble makes every creature on the board that could fight an enemy
	// have to fight when used — it cannot reap or use an Action ability while a legal
	// fight target exists (Little Rapscal). Affects both players' creatures.
	MustFightIfAble bool
}

// ConditionalPlayBar is a symmetric, board-wide play restriction a card imposes
// while it stays in play: any player for whom When is met cannot play cards of
// Type, whoever controls the card (Quixxle Stone bars whichever player controls
// more creatures). A nil When imposes no restriction.
type ConditionalPlayBar struct {
	Type CardType
	When Condition
}

// PlayCardLimit caps how many cards Player may play in a turn while its source
// card remains in play. Player is relative to the source card's controller, so
// Controller, Opponent, and EachPlayer compose naturally. Amount zero means no
// limit.
type PlayCardLimit struct {
	Player Player
	Amount int
}

// affects reports whether the limit on a card owned by controller applies to
// target.
func (l PlayCardLimit) affects(controller, target int) bool {
	switch l.Player {
	case Controller:
		return target == controller
	case Opponent:
		return target != controller
	case EachPlayer:
		return true
	default:
		return false
	}
}

// KeyCostChange is a continuous change to the cost of forging a key that a card in
// play imposes while it stays in play. Build it with NewKeyCostChange: the
// affected player is a required argument, so authors state whose keys change
// (Controller, Opponent, or EachPlayer) rather than lean on a zero-value default —
// an unset player would be indistinguishable from Controller. The zero value
// (which NewKeyCostChange never produces) changes nothing. (A Duration will later
// bound how long the change lasts; today every key-cost change is continuous.)
type KeyCostChange struct {
	// amount is the Æmber added to the affected keys' cost; player is whose keys
	// change (Controller, Opponent, or EachPlayer). Build with NewKeyCostChange.
	amount int
	player Player
	// per scales the amount by a running count read from the source card's point of
	// view — Iron Obelisk charges +1 per friendly damaged Brobnar creature.
	per Count
	// whileOnFlank suspends the change unless the source card holds a flank of its
	// controller's battleline (Titan Mechanic).
	whileOnFlank bool
	// whileCondition suspends the change unless the condition holds on the live
	// board (Proclamation 346E charges +2 only while the opponent controls creatures
	// from fewer than three houses).
	whileCondition Condition
}

// Per scales the change by a running count, so a card can charge per creature it
// sees rather than a flat amount.
func (kc KeyCostChange) Per(c Count) KeyCostChange {
	kc.per = c
	return kc
}

// WhileOnFlank applies the change only while the source card is on a flank.
func (kc KeyCostChange) WhileOnFlank() KeyCostChange {
	kc.whileOnFlank = true
	return kc
}

// While applies the change only while the condition holds on the live board.
func (kc KeyCostChange) While(c Condition) KeyCostChange {
	kc.whileCondition = c
	return kc
}

// selfHouseResolved fills the card's own house in for any SelfHouse sentinel the
// scaling count names (Iron Obelisk counts its own house's damaged creatures). A
// key-cost change keeps that count unexported, so it resolves itself rather than
// being rewritten by reflection (see self_house.go).
func (kc KeyCostChange) selfHouseResolved(house House) any {
	if kc.per != nil {
		kc.per = resolvedIn(kc.per, house)
	}
	return kc
}

// NewKeyCostChange builds a key-cost change of amount Æmber on the keys of player —
// one of Controller, Opponent, or EachPlayer. The player is mandatory: there is no
// default, so a key-cost change cannot be constructed without stating whose keys
// it changes (omitting it is a compile error at the call site).
func NewKeyCostChange(player Player, amount int) KeyCostChange {
	return KeyCostChange{amount: amount, player: player}
}

// affects reports whether a change on a card owned by owner applies to the key
// cost of target.
func (kc KeyCostChange) affects(owner, target int) bool {
	switch kc.player {
	case Opponent:
		return target == 1-owner
	case EachPlayer:
		return true
	default: // Controller
		return target == owner
	}
}

// StaticModifier is a continuous change applied by an Upgrade to the creature it
// is attached to.
type StaticModifier struct {
	// Flat stat bonuses the Upgrade adds to its host creature.
	PowerBonus        int
	ArmorBonus        int
	AssaultBonus      int
	HazardousBonus    int
	SplashAttackBonus int

	// Per scales the flat stat bonuses (PowerBonus, ArmorBonus) by a count read
	// off the host creature — Light of the Archons gives its host +1 power and
	// +1 armor for each upgrade attached to it (Per: UpgradesOnIt). The zero value
	// leaves the bonuses flat.
	Per PerTarget

	// Granted are triggered abilities the Upgrade grants its host creature. The
	// host fires them as if they were printed on it (see Game.triggerAbilities).
	Granted []Ability

	// Keywords are keywords the Upgrade grants its host creature; the host has
	// them in addition to its own (see Game.hasKeyword).
	Keywords []Keyword

	// KeywordsToNeighbors are keywords the Upgrade grants its host creature AND
	// each of the host's battleline neighbors — Cloaking Dongle gives the host and
	// both its neighbors Elusive. The host gains them in addition to Keywords (see
	// Game.hasKeyword).
	KeywordsToNeighbors []Keyword

	// KeyCostChange is a key-cost change an Upgrade grants its host; while attached
	// the host imposes it (e.g. "Your opponent's keys cost +2 Æmber").
	KeyCostChange KeyCostChange

	// AemberCannotBeStolen, while the Upgrade is attached, keeps the host's
	// controller's Æmber from being stolen (Discombobulator grants the host "Your
	// Æmber cannot be stolen.").
	AemberCannotBeStolen bool

	// Replaces is a continuous replacement the Upgrade applies to a game event's
	// outcome for its host while attached — Armageddon Cloak replaces the host's
	// destruction (EventCreatureDestroyed) with an effect that fully heals it and
	// destroys the Upgrade. The zero value carries no replacement.
	Replaces Replace

	// WhileOnFlank suspends the whole modifier unless the host creature holds a
	// flank of its controller's battleline — Shoulder Armor only armors a creature
	// standing at the edge of the line.
	WhileOnFlank bool

	// ProtectsFromNonFlank bars creatures that are not on a flank from being used
	// to fight the host creature — Camouflage. A creature on a flank may still
	// fight it.
	ProtectsFromNonFlank bool

	// HouseOverride, while the Upgrade is attached and its controller controls the
	// host, makes the host belong to this house instead of its printed one — Academy
	// Training makes its creature a Logos creature. HouseNone carries no override.
	HouseOverride House

	// SpendAemberOnCard lets the Æmber sitting on the host creature be spent to pay
	// Æmber costs as if it were in its controller's pool — The Callipygian Ideal
	// grants the creature it upgrades this permission. The pay path
	// (spendAsPoolCreatures) consults it; the zero value grants nothing.
	SpendAemberOnCard bool
}

// grants reports whether the modifier gives its host anything at all — a stat
// bonus, a keyword, a granted ability, a key-cost change, a replacement, or
// flank protection. It gates a creature-as-upgrade at init: playing such a card
// as an upgrade must actually do something (see NewCard).
func (m StaticModifier) grants() bool {
	return m.PowerBonus != 0 ||
		m.ArmorBonus != 0 ||
		m.AssaultBonus != 0 ||
		m.HazardousBonus != 0 ||
		m.SplashAttackBonus != 0 ||
		len(m.Granted) > 0 ||
		len(m.Keywords) > 0 ||
		len(m.KeywordsToNeighbors) > 0 ||
		m.KeyCostChange.amount != 0 ||
		m.Replaces.valid() ||
		m.ProtectsFromNonFlank ||
		m.HouseOverride != HouseNone ||
		m.AemberCannotBeStolen ||
		m.SpendAemberOnCard
}

// ConstantAbility is a continuous stat modifier a card in play applies to
// creatures — "Each friendly creature gains +1 power" — lasting only while the
// source card remains in play. It reuses Target to say which cards it reaches,
// evaluated from the source card's point of view. An unset Target (the zero
// value) reaches every card in play — creatures and artifacts, the source
// included.
type ConstantAbility struct {
	// Flat stat bonuses added to each creature the Target reaches.
	PowerBonus int
	ArmorBonus int
	// HazardousBonus is Hazardous the ability grants each creature the Target
	// reaches — Armsmaster Molina gives each of its neighbors hazardous 3.
	HazardousBonus int
	// Target says which cards the ability reaches, read from the source's point of
	// view; the zero value reaches every card in play.
	Target Target
	// Per scales the bonuses by a running count read from the source's point of
	// view — Mushroom Man gets +3 power for each unforged key its controller has,
	// Primus Unguis +2 power for each Æmber on itself. The same count reaches every
	// creature the Target names; use PerTarget when the count is read per creature.
	Per Count
	// PerTarget scales the bonuses by a count read separately for each creature the
	// Target reaches — Tribune Pompitus gives each friendly creature +2 power for
	// each Æmber on that creature, so a creature holding no Æmber gains nothing.
	PerTarget PerTarget
	// Keywords are keywords the card grants to every creature its Target reaches,
	// for as long as it stays in play — Round Table grants friendly Knights taunt.
	Keywords []Keyword
	// Granted are triggered abilities the card grants to every creature its Target
	// reaches, for as long as the card stays in play — Annihilation Ritual grants
	// each creature a "Destroyed: purge this creature." The reached creatures fire
	// them as if printed on them (see Game.triggerAbilities).
	Granted []Ability
	// CannotBeUsedTo bars the creatures its Target reaches from these ways of being
	// used, for as long as the card stays in play — Narp stops its neighbors from
	// reaping. It is the grantable form of CardDefinition.CannotBeUsedTo.
	CannotBeUsedTo []UseKind
	// Morphs are trigger morphs the card grants to every creature its Target reaches,
	// for as long as it stays in play — Kompsos Haruspex makes each friendly
	// creature's play effect also fire on reap. Each pair fires an ability under one
	// trigger when another occurs (see Game.morphedTriggers).
	Morphs []TriggerMorph
	// WhileOffFlank suspends the whole ability unless the source card is off a
	// flank (in the interior of its controller's battleline) — Gub's "While Gub is
	// not on a flank, it gets +5 power and gains taunt."
	WhileOffFlank bool
	// WhileInCenter suspends the whole ability unless the source card sits in the
	// center of its controller's battleline — Kaloch Stonefather grants friendly
	// creatures skirmish only while it is centered.
	WhileInCenter bool
	// WhileCondition suspends the whole ability unless the condition holds, read
	// from the source's point of view — The Red Baron grants itself a reap only
	// while your red key is forged. It is nil when the ability is always active.
	WhileCondition Condition
	// SpendAemberOnCard lets the Æmber sitting on each creature the Target reaches
	// be spent to pay Æmber costs as if it were in its controller's pool — Senator
	// Bracchus grants this to every friendly creature. The pay path
	// (spendAsPoolCreatures) consults it; the zero value grants nothing.
	SpendAemberOnCard bool
}

// target returns the constant ability's effective Target: an unset Target reaches
// every card in play (creatures and artifacts, including the source).
func (c ConstantAbility) target() Target {
	if c.Target == (Target{}) {
		return Target{Kind: TargetEachCardInPlay}
	}
	return c.Target
}

// Ability pairs a trigger with the effect that resolves when it fires.
type Ability struct {
	Trigger Trigger
	Effect  Effect
}

// hasKeyword reports whether the definition has the given keyword.
func (d *CardDefinition) hasKeyword(k Keyword) bool {
	for _, kw := range d.Keywords {
		if kw == k {
			return true
		}
	}
	return false
}

// hasTrait reports whether the definition has the given trait.
func (d *CardDefinition) hasTrait(t Trait) bool {
	for _, tr := range d.Traits {
		if tr == t {
			return true
		}
	}
	return false
}

// hasTrigger reports whether the definition has an ability with the trigger.
func (d *CardDefinition) hasTrigger(t Trigger) bool {
	for _, ab := range d.Abilities {
		if ab.Trigger == t {
			return true
		}
	}
	return false
}

// CardOption configures a CardDefinition. Definitions use the functional options
// pattern so optional fields read clearly and defaults are centralized.
type CardOption func(*CardDefinition)

// NewCard builds a CardDefinition. Required fields (including rarity) are
// positional; everything optional is supplied via options.
func NewCard(
	name string,
	house House,
	ct CardType,
	rarity Rarity,
	opts ...CardOption,
) CardDefinition {
	c := CardDefinition{
		Name:   name,
		House:  house,
		Type:   ct,
		Rarity: rarity,
	}
	for _, opt := range opts {
		opt(&c)
	}
	c = resolveSelfHouse(c)
	for _, ab := range c.Abilities {
		if !ab.Trigger.valid() {
			panic(fmt.Sprintf("card %q: an ability has no trigger set", name))
		}
		if err := validateEffect(ab.Effect); err != nil {
			panic(fmt.Sprintf("card %q: %v", name, err))
		}
	}
	if err := c.Static.Replaces.validate(); err != nil {
		panic(fmt.Sprintf("card %q: %v", name, err))
	}
	if c.Replaces.valid() {
		if err := c.Replaces.validate(); err != nil {
			panic(fmt.Sprintf("card %q: %v", name, err))
		}
	}
	if err := c.PlayPermission.validate(); err != nil {
		panic(fmt.Sprintf("card %q: %v", name, err))
	}
	if uc := c.Restricts.UseCondition; uc != nil {
		if err := validateCondition(uc); err != nil {
			panic(fmt.Sprintf("card %q: %v", name, err))
		}
	}
	if dw := c.DestroyedWhen; dw != nil {
		if err := validateCondition(dw); err != nil {
			panic(fmt.Sprintf("card %q: %v", name, err))
		}
	}
	for _, k := range c.CannotBeUsedTo {
		if !k.valid() {
			panic(fmt.Sprintf("card %q: CannotBeUsedTo has an unset use kind", name))
		}
	}
	for _, ca := range c.ConstantAbilities {
		for _, k := range ca.CannotBeUsedTo {
			if !k.valid() {
				panic(fmt.Sprintf(
					"card %q: a constant ability's CannotBeUsedTo has an unset use kind",
					name,
				))
			}
		}
		for _, m := range ca.Morphs {
			if !m.valid() {
				panic(fmt.Sprintf(
					"card %q: a constant ability's trigger morph names a non-action trigger",
					name,
				))
			}
		}
	}
	if c.PlayableAsUpgrade {
		if c.Type != Creature {
			panic(fmt.Sprintf("card %q: only a creature may be played as an upgrade", name))
		}
		if !c.Static.grants() {
			panic(fmt.Sprintf(
				"card %q: a creature played as an upgrade must grant its host something",
				name,
			))
		}
	}
	return c
}

// WithCannotBeUsedTo bars a card from the named ways of being used.
func WithCannotBeUsedTo(kinds ...UseKind) CardOption {
	return func(c *CardDefinition) { c.CannotBeUsedTo = append(c.CannotBeUsedTo, kinds...) }
}

// WithCannotBeUsedWhile bars a card from being used in any way for as long as cond
// holds for its controller (Valoocanth cannot be used while the tide is low).
func WithCannotBeUsedWhile(cond Condition) CardOption {
	return func(c *CardDefinition) { c.CannotBeUsedWhile = cond }
}

// WithDestroyedWhen makes a creature destroyable for as long as cond holds.
func WithDestroyedWhen(cond Condition) CardOption {
	return func(c *CardDefinition) { c.DestroyedWhen = cond }
}

// WithTakesDamageFor makes this card take the damage dealt to the creatures its
// Target names, instead of them (Shadow Self shields its non-Specter neighbors).
func WithTakesDamageFor(t Target) CardOption {
	return func(c *CardDefinition) { c.TakesDamageFor = t }
}

// WithPower sets a creature's power.
func WithPower(p int) CardOption { return func(c *CardDefinition) { c.Power = p } }

// WithArmor sets a creature's armor.
func WithArmor(a int) CardOption { return func(c *CardDefinition) { c.Armor = a } }

// WithTraits appends traits to the card.
func WithTraits(traits ...Trait) CardOption {
	return func(c *CardDefinition) { c.Traits = append(c.Traits, traits...) }
}

// WithKeywords appends keywords to the card.
func WithKeywords(keywords ...Keyword) CardOption {
	return func(c *CardDefinition) { c.Keywords = append(c.Keywords, keywords...) }
}

// WithAssault gives a creature Assault N: it deals N damage to the creature it
// attacks, before fight damage.
func WithAssault(n int) CardOption { return func(c *CardDefinition) { c.Assault = n } }

// WithHazardous gives a creature Hazardous N: a creature that attacks it is dealt
// N damage before fight damage.
func WithHazardous(n int) CardOption { return func(c *CardDefinition) { c.Hazardous = n } }

// WithSplashAttack gives a creature Splash-attack N: when it fights, it deals N
// damage to each neighbor of the creature it fights, at the same time as its own
// fight damage.
func WithSplashAttack(n int) CardOption {
	return func(c *CardDefinition) { c.SplashAttack = n }
}

// AttackDamage customizes the damage a creature deals when it fights. The zero
// value leaves fight damage equal to the creature's power; the fields override or
// adjust it for the handful of creatures that need it.
type AttackDamage struct {
	// Amount is the number. When Fixed it is the whole fight damage, replacing the
	// creature's power (Ether Spider deals 0); otherwise it is a bonus added to the
	// creature's power (Valdr's +2).
	Amount int
	// Fixed deals Amount as the entire fight damage instead of adding it to power.
	Fixed bool
	// FlankOnly limits a bonus (a non-Fixed Amount) to attacks on a defender that is
	// on a flank (Valdr). It does not restrict a Fixed amount.
	FlankOnly bool
}

// WithAttackDamage customizes the damage a creature deals when it fights.
func WithAttackDamage(ad AttackDamage) CardOption {
	return func(c *CardDefinition) { c.AttackDamage = ad }
}

// WithNoDamageWhenAttacked makes a creature deal no retaliation damage to an
// attacker that fights it (Lollop the Titanic).
func WithNoDamageWhenAttacked() CardOption {
	return func(c *CardDefinition) { c.DealsNoDamageWhenAttacked = true }
}

// WithFriendlyEntersPlayReady makes friendly cards of the given type enter play
// ready while this card is in play (Duskwitch for creatures, The Curator for artifacts).
func WithFriendlyEntersPlayReady(t CardType) CardOption {
	return func(c *CardDefinition) { c.GrantsEntersReady = t }
}

// WithFightRestriction limits which creatures a creature may fight to those the
// Target allows (e.g. card-level "can only fight stunned creatures").
func WithFightRestriction(t Target) CardOption {
	return func(c *CardDefinition) { c.FightRestriction = t }
}

// WithAttackIgnores makes a creature ignore defensive keywords while it attacks
// (Niffle Ape ignores taunt and elusive).
func WithAttackIgnores(kws ...Keyword) CardOption {
	return func(c *CardDefinition) { c.AttackIgnores = kws }
}

// AttackKeywords are keywords a creature gains only while it is attacking. The
// zero value grants none.
type AttackKeywords struct {
	// Keywords are the keywords the attacker gains for the fight.
	Keywords []Keyword
	// FlankOnly limits the grant to attacks on a defender that is on a flank
	// (Spyyyder gains poison only against a flank creature).
	FlankOnly bool
}

// WithAttackKeywords makes a creature gain keywords while it is attacking — with
// FlankOnly the grant applies only against a defender on a flank (Spyyyder).
func WithAttackKeywords(ak AttackKeywords) CardOption {
	return func(c *CardDefinition) { c.AttackKeywords = ak }
}

// WithEntersPlay makes a creature apply an effect to itself as it enters play
// (Chuff Ape stunning itself with Stun) by giving it that effect as an Enters Play
// ability — an ability the enter-play event fires, so the play path needs no
// special case for it.
func WithEntersPlay(e Effect) CardOption {
	return WithAbility(TriggerEntersPlay, e)
}

// WithAemberBonus sets the number of Æmber pips on the card.
func WithAemberBonus(n int) CardOption { return func(c *CardDefinition) { c.AemberBonus = n } }

// WithStatic sets the continuous modifier an Upgrade applies to its host.
func WithStatic(m StaticModifier) CardOption { return func(c *CardDefinition) { c.Static = m } }

// WithPlayableAsUpgrade lets a creature be played as an upgrade instead of a
// creature, granting its host the card's Static modifier. The card must be a
// creature and its Static must grant something (both validated in NewCard).
func WithPlayableAsUpgrade() CardOption {
	return func(c *CardDefinition) { c.PlayableAsUpgrade = true }
}

// WithConstantAbility appends a constant ability this card applies to creatures
// while it is in play.
func WithConstantAbility(c ConstantAbility) CardOption {
	return func(d *CardDefinition) { d.ConstantAbilities = append(d.ConstantAbilities, c) }
}

// WithRestrictions sets the continuous "cannot" rules a card imposes on its
// controller while it is in play.
func WithRestrictions(r Restrictions) CardOption {
	return func(c *CardDefinition) { c.Restricts = r }
}

// WithCannotPlayWhile sets the symmetric, board-wide play bar a card imposes while
// in play, barring any player who meets the condition from playing that type.
func WithCannotPlayWhile(b ConditionalPlayBar) CardOption {
	return func(c *CardDefinition) { c.CannotPlayWhile = b }
}

// WithHouseLock sets the continuous constraint a card puts on a player's
// active-house choice while it is in play.
func WithHouseLock(l HouseLock) CardOption {
	return func(c *CardDefinition) { c.HouseLock = l }
}

// WithKeyCost makes the card, while in play, impose the given key-cost change (who
// it affects and by how much). Called more than once, each change stacks (Grump
// Buggy raises each player's keys by a separate per-creature count).
func WithKeyCost(kc KeyCostChange) CardOption {
	return func(c *CardDefinition) { c.KeyCostChanges = append(c.KeyCostChanges, kc) }
}

// WithReplaces sets a continuous replacement the card applies to a game event's
// outcome while in play (Ether Spider capturing Æmber added to its opponent's pool).
func WithReplaces(r Instead) CardOption {
	return func(c *CardDefinition) { c.Replaces = r }
}

// WithDrawModifier makes the card, while in play, change a player's end-of-turn
// hand-refill size by amount (Mother +1 for its controller, Succubus -1 for the
// opponent, The Howling Pit +1 for each player).
func WithDrawModifier(player Player, amount int) CardOption {
	return func(c *CardDefinition) { c.DrawModifier = DrawModifier{Player: player, Amount: amount} }
}

// WithDrawModifierOffFlank is WithDrawModifier gated on the source not being on a
// flank of its battleline (Streke slows the opponent's refill only while buried in
// the middle of the line).
func WithDrawModifierOffFlank(player Player, amount int) CardOption {
	return func(c *CardDefinition) {
		c.DrawModifier = DrawModifier{Player: player, Amount: amount, OnlyWhileOffFlank: true}
	}
}

// WithDrawModifierInCenter is WithDrawModifier gated on the source sitting in the
// center of its battleline (Zenzizenzizenzic refills extra only from the middle).
func WithDrawModifierInCenter(player Player, amount int) CardOption {
	return func(c *CardDefinition) {
		c.DrawModifier = DrawModifier{Player: player, Amount: amount, OnlyWhileInCenter: true}
	}
}

// WithAemberCannotBeStolen makes the card, while in play, keep its controller's
// Æmber from being stolen (The Vaultkeeper).
func WithAemberCannotBeStolen() CardOption {
	return func(c *CardDefinition) { c.AemberCannotBeStolen = true }
}

// WithAemberCannotBeStolenWhileItHasAember keeps its controller's Æmber from
// being stolen only while the card itself has Æmber on it (Odoac the Patrician).
func WithAemberCannotBeStolenWhileItHasAember() CardOption {
	return func(c *CardDefinition) { c.AemberCannotBeStolenWhileItHasAember = true }
}

// WithSpendableAember lets the Æmber banked on the card be spent when its
// controller forges a key (Safe Place, Pocket Universe).
func WithSpendableAember() CardOption {
	return func(c *CardDefinition) { c.SpendableAember = true }
}

// WithGainsForgeAember gives the card's controller all the Æmber their opponent
// spends forging a key, for as long as it stays in play (The Sting).
func WithGainsForgeAember() CardOption {
	return func(c *CardDefinition) { c.GainsForgeAember = true }
}

// WithGuardsOpponentForge makes the card interrupt the opponent's key forges
// while it is in play (Keyforgery).
func WithGuardsOpponentForge() CardOption {
	return func(c *CardDefinition) { c.GuardsOpponentForge = true }
}

// WithPlayRequirement puts an Æmber requirement on playing the card, either a
// threshold it only checks (Kelifi Dragon) or a cost it charges (Truebaru).
func WithPlayRequirement(r PlayRequirement) CardOption {
	return func(c *CardDefinition) { c.PlayRequirement = r }
}

// WithAbility appends a triggered ability to the card.
func WithAbility(trigger Trigger, effect Effect) CardOption {
	return func(c *CardDefinition) {
		c.Abilities = append(c.Abilities, Ability{Trigger: trigger, Effect: effect})
	}
}
