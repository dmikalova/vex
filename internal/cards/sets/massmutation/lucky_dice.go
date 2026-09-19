package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Lucky Dice
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Lucky Dice. During your opponent's next turn, each friendly creature cannot be dealt damage.
var LuckyDice = set.New(
	"Lucky Dice",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "267"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.CannotBeDealtDamage{
				Target:   card.Target.EachFriendlyCreature,
				Duration: card.Duration.OpponentNextTurn,
			},
		}}),
)
