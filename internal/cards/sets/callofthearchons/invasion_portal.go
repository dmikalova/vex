//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// InvasionPortal
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Discard cards from the top of your deck until you discard a Mars creature or run out of cards. If you discard a Mars creature this way, put it into your hand.
var InvasionPortal = card.New(
	"Invasion Portal",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 185),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
