package engine

// PromptKind is the closed catalog of decision points a Chooser can be asked to
// answer — one kind per capability method on the Chooser interface and its
// optional extensions. It is the forward form of the decision-point context ADR
// 0040 folds into a suspendable Request: today a prompt is a synchronous call on
// a capability method, and each method here names the kind of decision that call
// makes. The catalog is closed by construction — a new prompt route is a new
// capability method and a new kind here together — so a human-facing renderer can
// be proven to cover every kind (ADR 0045).
type PromptKind uint8

const (
	// promptKindUnset is the invalid zero value (ADR 0010): a PromptKind must name
	// a real decision point, never the default of an unset field.
	promptKindUnset PromptKind = iota

	// PromptCreature is a mandatory pick of one candidate: Chooser.ChooseCreature.
	PromptCreature
	// PromptCardOrDecline is an optional card pick the player may pass on:
	// DeclinableChooser.ChooseCardOrDecline.
	PromptCardOrDecline
	// PromptOption is a labeled multiple choice: OptionChooser.ChooseOption.
	PromptOption
	// PromptPosition is a battleline placement for a Deploy creature:
	// PositionChooser.ChoosePosition.
	PromptPosition
	// PromptReaction is picking which reaction in a trigger window resolves next:
	// ReactionChooser.ChooseReaction.
	PromptReaction
	// PromptOrder is arranging ids into a resolution order: Orderer.OrderCreatures.
	PromptOrder
	// PromptBadge is the status preview a client draws on a candidate before it is
	// chosen: BadgeChooser.PreviewBadge. It carries no decision of its own — the
	// pick it decorates is a PromptCreature or PromptCardOrDecline — but it is a
	// distinct capability a renderer must handle, so it is its own kind.
	PromptBadge
)

// PromptKinds returns every real prompt kind, in declaration order, excluding the
// invalid zero. A renderer's totality test iterates this to prove it covers them
// all.
func PromptKinds() []PromptKind {
	return []PromptKind{
		PromptCreature,
		PromptCardOrDecline,
		PromptOption,
		PromptPosition,
		PromptReaction,
		PromptOrder,
		PromptBadge,
	}
}

// String names the kind for logs and test failures.
func (k PromptKind) String() string {
	switch k {
	case PromptCreature:
		return "creature"
	case PromptCardOrDecline:
		return "card-or-decline"
	case PromptOption:
		return "option"
	case PromptPosition:
		return "position"
	case PromptReaction:
		return "reaction"
	case PromptOrder:
		return "order"
	case PromptBadge:
		return "badge"
	case promptKindUnset:
		return "unset"
	default:
		return "invalid"
	}
}
