package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Dr. Milli
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Scientist
//
//	Play: For each creature your opponent controls in excess of you, not counting Dr. Milli, archive a card from your hand.
var DrMilli = card.New(
	"Dr. Milli",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "150"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Scientist),
	card.WithAbility(card.Trigger.Play, card.ArchiveCard{
		Zone:      card.Hand,
		Selection: card.Chosen{},
		Per:       card.ExcessCreatures{Player: card.Opponent, NotCountingSelf: true},
	}),
)
