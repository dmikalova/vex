package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bonkers Killing Machine
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Weapon
//
//	Action: Discard the top card of each player's deck. For each card discarded this way, destroy a creature or artifact of that card's house. If fewer than 2 cards are destroyed this way, destroy Bonkers Killing Machine.
func TestBonkersKillingMachine(t *testing.T) {
	t.Run("destroys one matching card for each discarded house and stays in play", func(t *testing.T) {
		var p1Top, p2Top, marsCreature, marsArtifact, shadowsCreature ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					BonkersKillingMachine,
					ct.Bind(&marsCreature, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
					ct.Bind(&shadowsCreature, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
				),
				Deck: ct.Cards(
					ct.Bind(&p1Top, ct.Tactic(ct.OfHouse(card.House.Mars))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&marsArtifact, ct.Artifact(ct.OfHouse(card.House.Mars))),
				),
				Deck: ct.Cards(
					ct.Bind(&p2Top, ct.Tactic(ct.OfHouse(card.House.Shadows))),
				),
			},
		})

		h.P1.UseAction(BonkersKillingMachine)
		h.P1.ExpectPrompt("Choose a Mars creature or artifact").Source("Bonkers Killing Machine")
		h.P1.ClickCard(marsArtifact)

		h.Expect(p1Top).At(ct.Discard)
		h.Expect(p2Top).At(ct.Discard)
		h.Expect(marsArtifact).At(ct.Discard)
		h.Expect(shadowsCreature).At(ct.Discard)
		h.Expect(marsCreature).At(ct.PlayArea)
		h.Expect(BonkersKillingMachine).At(ct.PlayArea)
	})

	t.Run("destroys itself when fewer than two cards are destroyed", func(t *testing.T) {
		var p1Top, p2Top, bystander ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(BonkersKillingMachine),
				Deck: ct.Cards(
					ct.Bind(&p1Top, ct.Tactic(ct.OfHouse(card.House.Mars))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&bystander, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4)))),
				Deck: ct.Cards(
					ct.Bind(&p2Top, ct.Tactic(ct.OfHouse(card.House.Dis))),
				),
			},
		})

		h.P1.UseAction(BonkersKillingMachine)

		h.Expect(p1Top).At(ct.Discard)
		h.Expect(p2Top).At(ct.Discard)
		h.Expect(bystander).At(ct.PlayArea)
		h.Expect(BonkersKillingMachine).At(ct.Discard)
	})
}
