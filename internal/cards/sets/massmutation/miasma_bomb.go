package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Miasma Bomb
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Destroy Miasma Bomb -> your opponent skips the "forge a key" phase during their next turn.
//	Enhance Damage.
var MiasmaBomb = set.New(
	"Miasma Bomb",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "269"),
	card.WithEnhance(card.Bonus.Damage),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First:  card.Destroy{Target: card.Target.This},
			Result: card.SkipForgePhase{Player: card.Opponent},
		}),
)
