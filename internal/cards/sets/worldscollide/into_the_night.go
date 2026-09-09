package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Into the Night
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Until the start of your next turn, non-Shadows creatures cannot be used to fight.
var IntoTheNight = card.New(
	"Into the Night",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "256"),
	card.WithAbility(
		card.Trigger.Play, card.CreaturesCannot{
			Action:      card.UseKind.Fight,
			ExceptHouse: card.House.Self,
			Duration:    card.Duration.NextTurn,
		}),
)
