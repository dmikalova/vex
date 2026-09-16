package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Z-Ray Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Connected
//	Bonus:  Æmber
//
//	This creature gains +3 power and +3 splash-attack.
var ZRayBlaster = set.New(
	"Z-Ray Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Connected,
	card.Provenance(card.MM, "355"),
	card.InCluster(zForceCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		PowerBonus:        3,
		SplashAttackBonus: 3,
	}),
)
