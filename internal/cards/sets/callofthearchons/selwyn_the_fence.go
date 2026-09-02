package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Selwyn the Fence
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Thief
//
//	Fight/Reap: Move 1 Æmber from a friendly creature or artifact to your pool.
var SelwynTheFence = card.New(
	"Selwyn the Fence",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 309),
	card.WithPower(3),
	card.WithTraits("Elf", "Thief"),
	card.WithFightOrReap(card.MoveAember{
		Amount: 1,
		From:   card.Target.FriendlyCreatureOrArtifact,
		To:     card.Controller,
	}),
)
