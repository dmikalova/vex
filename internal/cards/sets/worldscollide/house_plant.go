package worldscollide

import (
	"math/rand"

	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/engine"
)

// House Plant
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//
//	Template: its concrete card is materialized per deck at generation.
var HousePlant = set.New(
	"House Plant",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "286"),
	card.Provenance(card.WC, "287"),
	card.Provenance(card.WC, "288"),
	card.Provenance(card.WC, "289"),
	card.Provenance(card.WC, "290"),
	card.Provenance(card.WC, "291"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.Template(plantFor),
)

// plantFor materializes the Plant for a partner House. With no deck Houses in
// context — the uniqueness sampling test — it falls back to a random partner so
// the whole cycle is still exercised.
func plantFor(ctx card.SlotContext, r *rand.Rand) card.Definition {
	partner := plantPartner(ctx, r)
	return card.Build(
		partner.String()+" Plant",
		card.House.Shadows,
		card.Type.Creature,
		card.Rarity.Rare,
		card.WithPower(1),
		card.WithTraits(card.Traits.Elf, card.Traits.Thief),
		card.WithKeywords(card.Keyword.Elusive),
		card.WithEachPlayerAbility(
			card.Trigger.AfterChooseHouse, card.Conditional{
				Cond: card.ChoseHouse{House: partner},
				Then: card.GainAember{
					Player: card.Controller,
					Amount: 1,
				},
			}),
	)
}

// plantPartner picks the House a Plant watches for: one of the deck's Houses
// other than its own Shadows, or a random non-Shadows House when the context
// carries no deck Houses.
func plantPartner(ctx card.SlotContext, r *rand.Rand) engine.House {
	var partners []engine.House
	for _, h := range ctx.DeckHouses {
		if h != engine.HouseNone && h != engine.Shadows {
			partners = append(partners, h)
		}
	}
	if len(partners) == 0 {
		partners = nonShadowsHouses()
	}
	return partners[r.Intn(len(partners))]
}

// nonShadowsHouses is every real House but Shadows, in enum order.
func nonShadowsHouses() []engine.House {
	houses := make([]engine.House, 0, engine.NumHouses)
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		if h != engine.Shadows {
			houses = append(houses, h)
		}
	}
	return houses
}
