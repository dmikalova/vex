package engine

import "fmt"

// This file holds the log entries that narrate combat and damage (ADR 0011):
// the fight itself, what stopped it, and what each point of damage actually did
// once armor and immunity had their say.

// FightCancelled narrates a fight that an effect stopped before it happened.
type FightCancelled struct{ Attacker LocalID }

// Text renders a fight an effect stopped before it happened.
func (e FightCancelled) Text(n Namer) string {
	return fmt.Sprintf("%s's fight does not occur", n.Name(e.Attacker))
}

// Fought narrates two creatures fighting, with the power each brought to it.
type Fought struct {
	Attacker      LocalID
	AttackerPower int
	Defender      LocalID
	DefenderPower int
}

// Text renders the fight, and the power each creature brought to it.
func (e Fought) Text(n Namer) string {
	return fmt.Sprintf("%s (%d power) fights %s (%d power)",
		n.Name(e.Attacker), e.AttackerPower, n.Name(e.Defender), e.DefenderPower)
}

// ElusiveAvoidedFight narrates elusive turning a fight's damage aside.
type ElusiveAvoidedFight struct{ Defender LocalID }

// Text renders elusive turning a fight's damage aside.
func (e ElusiveAvoidedFight) Text(n Namer) string {
	return fmt.Sprintf("%s is elusive — no fight damage is dealt", n.Name(e.Defender))
}

// DamageRefused narrates damage that a creature could not be dealt at all.
type DamageRefused struct{ Creature LocalID }

// Text renders damage a creature could not be dealt at all.
func (e DamageRefused) Text(n Namer) string {
	return fmt.Sprintf("%s cannot be dealt damage", n.Name(e.Creature))
}

// ArmorAbsorbed narrates the part of an incoming hit that armor soaked up.
type ArmorAbsorbed struct {
	Creature LocalID
	Amount   int
}

// Text renders the part of an incoming hit that armor soaked up.
func (e ArmorAbsorbed) Text(n Namer) string {
	return fmt.Sprintf("%s's armor absorbs %d damage", n.Name(e.Creature), e.Amount)
}

// DamageTaken narrates damage that actually landed, and the total now on the
// creature.
type DamageTaken struct {
	Creature LocalID
	Amount   int
	Total    int
}

// Text renders the damage that landed, and the creature's new total.
func (e DamageTaken) Text(n Namer) string {
	return fmt.Sprintf("%s takes %d damage (%d total)",
		n.Name(e.Creature), e.Amount, e.Total)
}
