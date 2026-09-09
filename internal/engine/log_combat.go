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

// SkirmishAvoidedReturn narrates skirmish sparing an attacker the return damage a
// fight would otherwise deal it.
type SkirmishAvoidedReturn struct{ Attacker LocalID }

// Text renders skirmish sparing an attacker its return damage.
func (e SkirmishAvoidedReturn) Text(n Namer) string {
	return fmt.Sprintf("%s is skirmish — it takes no damage in return", n.Name(e.Attacker))
}

// PoisonKills narrates a poison creature's fight damage proving lethal: the
// creature it damaged is destroyed, however much power it had left.
type PoisonKills struct {
	Source LocalID
	Victim LocalID
}

// Text renders poison destroying a creature the poison combatant damaged.
func (e PoisonKills) Text(n Namer) string {
	return fmt.Sprintf("%s's poison is lethal to %s", n.Name(e.Source), n.Name(e.Victim))
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

// AssaultDealt narrates the pre-fight Assault an attacker deals the creature it
// attacks: the source, the Assault value, the damage that landed, and the target.
type AssaultDealt struct {
	Source LocalID
	Value  int
	Amount int
	Target LocalID
}

// Text renders the Assault damage, naming its source and keyword value.
func (e AssaultDealt) Text(n Namer) string {
	return fmt.Sprintf("%s's %d Assault deals %d damage to %s",
		n.Name(e.Source), e.Value, e.Amount, n.Name(e.Target))
}

// HazardousDealt narrates the pre-fight Hazardous a defender deals its attacker:
// the source, the Hazardous value, the damage that landed, and the target.
type HazardousDealt struct {
	Source LocalID
	Value  int
	Amount int
	Target LocalID
}

// Text renders the Hazardous damage, naming its source and keyword value.
func (e HazardousDealt) Text(n Namer) string {
	return fmt.Sprintf("%s's %d Hazardous deals %d damage to %s",
		n.Name(e.Source), e.Value, e.Amount, n.Name(e.Target))
}

// damageEntry chooses how a landed hit narrates: pre-fight Assault or Hazardous
// names its striking creature and keyword value, and any other damage is the bare
// DamageTaken line with the creature's new total. dealt is the damage left after
// armor; total is the creature's damage after taking it.
func (t DamageTarget) damageEntry(target LocalID, dealt, total int) LogEntry {
	switch t.SourceKeyword {
	case assaultDamage:
		return AssaultDealt{Source: t.Source, Value: t.Amount, Amount: dealt, Target: target}
	case hazardousDamage:
		return HazardousDealt{Source: t.Source, Value: t.Amount, Amount: dealt, Target: target}
	default:
		return DamageTaken{Creature: target, Amount: dealt, Total: total}
	}
}
