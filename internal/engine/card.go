package engine

import "fmt"

// CardDefinition is the immutable blueprint for a card. Definitions are shared,
// read-only data held in the match Catalog; all mutable per-match state lives in
// the flat GameState. This split keeps the runtime state a plain value that can
// be copied cheaply (see GameState.FastCopy).
type CardDefinition struct {
	Name   string
	House  House
	Type   CardType
	Rarity Rarity
	Traits []Trait

	// Creature stats. Zero for non-creatures.
	Power int
	Armor int

	Keywords []Keyword

	// A creature with Assault N deals N damage to the creature it attacks,
	// immediately before combat damage is dealt. Zero means the creature does not
	// have Assault.
	//
	//rulebook:keyword Assault
	Assault int
	// A creature with Hazardous N deals N damage to any creature that attacks it,
	// before that attacker deals its combat damage. Zero means the creature does
	// not have Hazardous.
	//
	//rulebook:keyword Hazardous
	Hazardous int
	// AttackDamage customizes the damage this creature deals when it fights, for the
	// few creatures whose fight damage is not simply their current power (Valdr's
	// flank bonus, Ether Spider dealing none). The zero value deals its power.
	AttackDamage AttackDamage

	// FightRestriction, when set, limits which enemy creatures this creature may
	// fight to those its Target allows (Bigtwig can only fight stunned creatures).
	// The zero Target imposes no restriction.
	FightRestriction Target

	// AemberBonus is the number of Æmber "pips" printed on the card; the
	// controller gains this much Æmber when the card is played.
	AemberBonus int

	// Static is a continuous modifier an Upgrade applies to its host creature.
	Static StaticModifier

	// Constant is a constant ability this card applies to creatures in play for as
	// long as it stays in play (see Game.constantBonus). Unlike Static (an Upgrade
	// buffing its own host), a constant ability reaches a whole Target of creatures.
	Constant ConstantAbility

	// Restricts holds the continuous "cannot" rules the card imposes on its
	// controller while it is in play (e.g. Grommid's "You cannot play creatures").
	Restricts Restrictions

	// KeyCostChange is a continuous change this card, while in play, makes to key
	// cost — who it affects and by how much (e.g. Grabber Jammer's "Your opponent's
	// keys cost +1 Æmber"). The zero value changes nothing.
	KeyCostChange KeyCostChange

	// PlayPermission is a continuous grant this card makes while in play, letting
	// its controller play cards of a house on turns where that house is not their
	// active house (Witch of the Wilds). The zero value grants nothing.
	PlayPermission PlayPermission

	// Replaces is a continuous replacement this card applies to a game event's
	// outcome while it is in play — Ether Spider replaces Æmber being added to its
	// opponent's pool (EventAemberAddedToPool, scoped by Player) with capturing it.
	// The zero value carries no replacement. (An Upgrade grants a replacement to its
	// host through StaticModifier.Replaces instead.)
	Replaces Instead

	// Abilities are the triggered abilities on the card.
	Abilities []Ability
}

// Restrictions are the continuous "cannot" rules a card imposes while it stays in
// play. They are consulted by the matching gate (Game.cannotFight,
// Game.cannotPlayCreatures, Game.cannotPlayCard) alongside any timed restriction.
type Restrictions struct {
	// Fighting bars the controller from using creatures to fight.
	Fighting bool
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
	amount int
	player Player
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
	PowerBonus     int
	ArmorBonus     int
	AssaultBonus   int
	HazardousBonus int

	// Granted are triggered abilities the Upgrade grants its host creature. The
	// host fires them as if they were printed on it (see Game.triggerAbilities).
	Granted []Ability

	// Keywords are keywords the Upgrade grants its host creature; the host has
	// them in addition to its own (see Game.hasKeyword).
	Keywords []Keyword

	// KeyCostChange is a key-cost change an Upgrade grants its host; while attached
	// the host imposes it (e.g. "Your opponent's keys cost +2 Æmber").
	KeyCostChange KeyCostChange

	// Replaces is a continuous replacement the Upgrade applies to a game event's
	// outcome for its host while attached — Armageddon Cloak replaces the host's
	// destruction (EventCreatureDestroyed) with an effect that fully heals it and
	// destroys the Upgrade. The zero value carries no replacement.
	Replaces Replace
}

// ConstantAbility is a continuous stat modifier a card in play applies to
// creatures — "Each friendly creature gains +1 power" — lasting only while the
// source card remains in play. It reuses Target to say which cards it reaches,
// evaluated from the source card's point of view. An unset Target (the zero
// value) reaches every card in play — creatures and artifacts, the source
// included.
type ConstantAbility struct {
	PowerBonus int
	ArmorBonus int
	Target     Target
	// Granted are triggered abilities the card grants to every creature its Target
	// reaches, for as long as the card stays in play — Annihilation Ritual grants
	// each creature a "Destroyed: purge this creature." The reached creatures fire
	// them as if printed on them (see Game.triggerAbilities).
	Granted []Ability
}

// target returns the constant ability's effective Target: an unset Target reaches
// every card in play (creatures and artifacts, including the source).
func (c ConstantAbility) target() Target {
	if c.Target == (Target{}) {
		return Target{Kind: TargetEachInPlay}
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
func NewCard(name string, house House, ct CardType, rarity Rarity, opts ...CardOption) CardDefinition {
	c := CardDefinition{
		Name:   name,
		House:  house,
		Type:   ct,
		Rarity: rarity,
	}
	for _, opt := range opts {
		opt(&c)
	}
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
	return c
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

// WithFightRestriction limits which creatures a creature may fight to those the
// Target allows (e.g. card-level "can only fight stunned creatures").
func WithFightRestriction(t Target) CardOption {
	return func(c *CardDefinition) { c.FightRestriction = t }
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

// WithConstantAbility sets the constant ability this card applies to creatures
// while it is in play.
func WithConstantAbility(c ConstantAbility) CardOption {
	return func(d *CardDefinition) { d.Constant = c }
}

// WithRestrictions sets the continuous "cannot" rules a card imposes on its
// controller while it is in play.
func WithRestrictions(r Restrictions) CardOption {
	return func(c *CardDefinition) { c.Restricts = r }
}

// WithKeyCost makes the card, while in play, impose the given key-cost change (who
// it affects and by how much).
func WithKeyCost(kc KeyCostChange) CardOption {
	return func(c *CardDefinition) { c.KeyCostChange = kc }
}

// WithReplaces sets a continuous replacement the card applies to a game event's
// outcome while in play (Ether Spider capturing Æmber added to its opponent's pool).
func WithReplaces(r Instead) CardOption {
	return func(c *CardDefinition) { c.Replaces = r }
}

// WithAbility appends a triggered ability to the card.
func WithAbility(trigger Trigger, effect Effect) CardOption {
	return func(c *CardDefinition) {
		c.Abilities = append(c.Abilities, Ability{Trigger: trigger, Effect: effect})
	}
}
