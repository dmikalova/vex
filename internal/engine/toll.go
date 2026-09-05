package engine

// A toll is Æmber a card in play makes its controller's opponent give in order to
// take an action with an artifact — playing an artifact, or using an artifact's
// ability. The opponent cannot take the action unless they can pay the toll, and
// the Æmber they give goes to the toll card's controller. (The mechanic keeps the
// name Toll, but its printed text always reads "give", never "pay".)
// A Toll is Æmber a card, while in play, makes its controller's opponent pay in
// order to take an action with an artifact — Customs Office's "pay 1 Æmber to
// play an artifact", Tentacus's "pay 1 Æmber to use an artifact". The opponent
// cannot take the action unless they can pay, and the Æmber goes to the toll
// card's controller.
type Toll struct {
	// Action is the action the opponent is charged for.
	Action TollAction
	// Amount is the Æmber owed; zero imposes no toll.
	Amount int
}

// TollAction names the action a Toll charges the opponent for.
type TollAction uint8

const (
	// tollActionUnset is the invalid zero value; a toll must name its action.
	tollActionUnset TollAction = iota
	// TollPlayArtifact charges the opponent to play an artifact.
	TollPlayArtifact
	// TollUseArtifact charges the opponent to use an artifact's action ability.
	TollUseArtifact
)

// phrase renders the toll's action for card text and the log, e.g. "play an
// artifact".
func (a TollAction) phrase() string {
	if a == TollUseArtifact {
		return "use an artifact"
	}
	return "play an artifact"
}
