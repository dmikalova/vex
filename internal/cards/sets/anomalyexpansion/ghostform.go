package anomalyexpansion

import "github.com/dmikalova/vactrol/internal/card"

// Ghostform
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Special
//	Bonus:  Æmber
//
//	This creature gains invulnerable.
//	This creature gains, "Fight/Reap: Archive Ghostform."
var Ghostform = set.New(
	"Ghostform",
	card.House.Brobnar,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A01"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Keywords: card.Keywords(card.Keyword.Invulnerable),
		Granted:  card.FightReap(card.ArchiveGrantingUpgrade{}),
	}),
)
