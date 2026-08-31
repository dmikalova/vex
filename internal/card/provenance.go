package card

import "github.com/dmikalova/vactrol/internal/cards/provenance"

// Provenance tags a card as derived from an original source card (set + collector
// number), for coverage tracking. Optional and repeatable — a card may draw from
// more than one original — e.g. card.Provenance(card.CotA, 1).
//
// It is purely a bookkeeping tag: it records which original KeyForge card an
// implementation is based on, so the author can confirm every original is
// eventually covered (see `mage missing`/`mage coverage`). Nothing in the engine
// or in deck generation ever reads it, and a card's behavior never depends on it.
func Provenance(set provenance.SourceSet, number int) Option {
	return func(b *builder) { b.prov = append(b.prov, provenance.Ref{Set: set, Number: number}) }
}

// Source sets to tag a card's Provenance with, e.g. card.Provenance(card.CotA, 1).
var (
	// CotA is Call of the Archons.
	CotA = provenance.CallOfTheArchons
	// AoA is Age of Ascension.
	AoA = provenance.AgeOfAscension
	// WC is Worlds Collide.
	WC = provenance.WorldsCollide
	// MM is Mass Mutation.
	MM = provenance.MassMutation
	// DT is Dark Tidings.
	DT = provenance.DarkTidings
	// WoE is Winds of Exchange.
	WoE = provenance.WindsOfExchange
	// GR is Grim Reminders.
	GR = provenance.GrimReminders
	// AS is Æmber Skies.
	AS = provenance.AemberSkies
	// ToC is Tokens of Change.
	ToC = provenance.TokensOfChange
	// MoM is More Mutation.
	MoM = provenance.MoreMutation
	// Men is Menagerie.
	Men = provenance.Menagerie
	// VM is Vault Masters 2025.
	VM = provenance.VaultMasters2025
	// PV is Prophetic Visions.
	PV = provenance.PropheticVisions
	// CC is Crucible Clash.
	CC = provenance.CrucibleClash
	// DM is Draconian Measures.
	DM = provenance.DraconianMeasures
)
