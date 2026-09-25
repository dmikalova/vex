package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Hidden Stash
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Reveal your opponent's hand. Archive a card from your opponent's hand.
var HiddenStash = set.New(
	"Hidden Stash",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "271"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.RevealHand{Player: card.Opponent},
			card.ArchiveCard{
				From:      card.Opponent,
				Zone:      card.Hand,
				Selection: card.Chosen{},
			},
		}},
	),
)
