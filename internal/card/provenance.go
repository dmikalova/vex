package card

import "github.com/dmikalova/vactrol/internal/cards/provenance"

// Provenance tags a card as derived from an original source card (set + collector
// number), for coverage tracking. Optional and repeatable — a card may draw from
// more than one original — e.g. card.Provenance(card.CotA, "001"). The collector
// number is a string so lettered reference cards (anomalies, prophecies) fit.
//
// It is purely a bookkeeping tag: it records which original KeyForge card an
// implementation is based on, so the author can confirm every original is
// eventually covered (see `mage tool:missing`/`mage tool:coverage`). Nothing in the engine
// or in deck generation ever reads it, and a card's behavior never depends on it.
func Provenance(set provenance.SourceSet, number string) Option {
	return func(b *builder) {
		b.prov = append(b.prov, provenance.Ref{
			Set:    set,
			Number: number,
		})
	}
}

// InSet declares the source set a card belongs to for deck generation, decoupled
// from Provenance. Deck generation groups a card into its set's pool by this set
// alone; it never reads Provenance (ADR 0003). Every card declares its set — a set
// package's registrar (set.New) stamps InSet for it — so the set is always
// explicit and never inferred. This option is the low-level primitive that
// registrar applies; author cards through set.New rather than calling it directly.
func InSet(set provenance.SourceSet) Option {
	return func(b *builder) { b.set = set }
}

// Source sets to tag a card's Provenance with, e.g. card.Provenance(card.CotA, "001").
var (
	// CotA is Call of the Archons.
	CotA = provenance.CallOfTheArchons
	// AoA is Age of Ascension.
	AoA = provenance.AgeOfAscension
	// WC is Worlds Collide.
	WC = provenance.WorldsCollide
	// AE is Anomaly Expansion, a reservoir set (never itself draftable): it holds
	// the Vactrol-invented Shards that enter a deck only through a cross-set
	// cluster, plus the housed anomaly cards, which stay legacy-drawable by other
	// sets.
	AE = provenance.AnomalyExpansion
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
	// MoMu is More Mutation.
	MoMu = provenance.MoreMutation
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
