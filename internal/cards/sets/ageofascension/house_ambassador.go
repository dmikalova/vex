package ageofascension

import (
	"math/rand"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/engine"
)

// House Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Human
//
//	Elusive.
//
//	Template: its concrete card is materialized per deck at generation.
var HouseAmbassador = set.New(
	"House Ambassador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "229"),
	card.Provenance(card.AoA, "230"),
	card.Provenance(card.AoA, "237"),
	card.Provenance(card.AoA, "238"),
	card.Provenance(card.AoA, "243"),
	card.Provenance(card.AoA, "247"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	card.Template(ambassadorFor),
)

// ambassadorFor materializes the Ambassador for a partner House. With no deck
// Houses in context — the uniqueness sampling test — it falls back to a random
// partner so the whole cycle is still exercised.
func ambassadorFor(ctx card.SlotContext, r *rand.Rand) card.Definition {
	partner := ambassadorPartner(ctx, r)
	return card.Build(
		partner.String()+" Ambassador",
		card.House.Sanctum,
		card.Type.Creature,
		card.Rarity.Rare,
		card.WithPower(1),
		card.WithTraits(card.Traits.Human),
		card.WithKeywords(card.Keyword.Elusive),
		card.WithAbility(
			card.Trigger.FightReap,
			card.MayPlayOrUse{
				Houses: card.GrantHouses.Named(partner),
				Grant:  card.GrantPlay | card.GrantUse,
			},
		),
	)
}

// ambassadorPartner picks the House an Ambassador vouches for: one of the deck's
// Houses other than its own Sanctum, or a random non-Sanctum House when the
// context carries no deck Houses.
func ambassadorPartner(ctx card.SlotContext, r *rand.Rand) engine.House {
	var partners []engine.House
	for _, h := range ctx.DeckHouses {
		if h != engine.HouseNone && h != engine.Sanctum {
			partners = append(partners, h)
		}
	}
	if len(partners) == 0 {
		partners = nonSanctumHouses()
	}
	return partners[r.Intn(len(partners))]
}

// nonSanctumHouses is every real House but Sanctum, in enum order.
func nonSanctumHouses() []engine.House {
	houses := make([]engine.House, 0, engine.NumHouses)
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		if h != engine.Sanctum {
			houses = append(houses, h)
		}
	}
	return houses
}
