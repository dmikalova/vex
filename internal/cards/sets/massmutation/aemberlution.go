package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Aemberlution
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Omega.
//	Play: Destroy each creature, and each player reveals their hand and puts each creature from their hand into play ready.
var Aemberlution = set.New(
	"Aemberlution",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "394"),
	card.WithKeywords(card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Destroy{Target: card.Target.EachCreature},
				card.EachPlayerPutsHandCreaturesIntoPlay{Ready: true},
			},
		}),
)
