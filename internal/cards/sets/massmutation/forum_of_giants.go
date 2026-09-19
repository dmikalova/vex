package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Forum of Giants
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	At the start of your turn, choose the most powerful creature. Its controller gains 1 Æmber.
var ForumOfGiants = set.New(
	"Forum of Giants",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "219"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(card.Trigger.StartOfTurn, card.ChooseCreatureThen{
		Target: card.Target.EachCreature.Refine(card.MostPowerful),
		Then:   card.GainAember{Player: card.ItsController, Amount: 1},
	}),
)
