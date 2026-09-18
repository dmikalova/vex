package web

import (
	"reflect"
	"sort"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is the Card text section of the Style gallery: one specimen per
// feature a card's printed rules text can carry, each a real card the query
// matched rather than a name written down. See
// docs/adr/0046-log-entry-catalog-and-sampled-gallery.md. The queries range over
// the engine's own trigger registry and the card definition's feature fields, so
// a feature added to the engine shows up as a new specimen — as a gap, if no
// loaded card carries it yet.

// defHasTrigger reports whether a card carries a triggered ability of the given
// trigger, the predicate the trigger specimens select on.
func defHasTrigger(d *engine.CardDefinition, t engine.Trigger) bool {
	for _, ab := range d.Abilities {
		if ab.Trigger == t {
			return true
		}
	}
	return false
}

// triggerSpecimens returns one card per printed trigger, ranging over
// engine.Triggers() rather than a list written here so a new printed trigger
// shows up as a new specimen. Unprinted triggers (EntersPlay, AfterChooseHouse)
// print no text prefix, so they are skipped: this section shows text.
func triggerSpecimens() []specimen {
	var out []specimen
	for _, t := range engine.Triggers() {
		if !t.Printed() {
			continue
		}
		t := t
		out = append(out, randomMatch(t.String(), func(d *engine.CardDefinition) bool {
			return defHasTrigger(d, t)
		}))
	}
	return out
}

// cardTextFeatures are the continuous, non-triggered text features a card can
// carry, each named and paired with the field predicate that detects it. Listing
// them here is what lets the gallery show one real card per feature in one place;
// adding a feature is adding a row.
var cardTextFeatures = []struct {
	caption string
	has     func(*engine.CardDefinition) bool
}{
	{"Grants an ability to creatures", func(d *engine.CardDefinition) bool {
		return len(d.ConstantAbilities) > 0
	}},
	{"Changes key cost", func(d *engine.CardDefinition) bool {
		return len(d.KeyCostChanges) > 0
	}},
	{"Æmber bonus", func(d *engine.CardDefinition) bool {
		return d.AemberBonus() > 0
	}},
	{"Æmber cannot be stolen", func(d *engine.CardDefinition) bool {
		return d.AemberCannotBeStolen != nil
	}},
	{"Spendable Æmber", func(d *engine.CardDefinition) bool {
		return d.SpendableAember
	}},
	{"Gains Æmber when opponent forges", func(d *engine.CardDefinition) bool {
		return d.GainsForgeAember
	}},
	{"Playable as an upgrade", func(d *engine.CardDefinition) bool {
		return d.PlayableAsUpgrade
	}},
}

// continuousSpecimens returns one card per continuous text feature.
func continuousSpecimens() []specimen {
	out := make([]specimen, 0, len(cardTextFeatures))
	for _, f := range cardTextFeatures {
		out = append(out, randomMatch(f.caption, f.has))
	}
	return out
}

// targetType is the reflect type walkEffectTargets stops on.
var targetType = reflect.TypeOf(engine.Target{})

// walkEffectTargets calls fn for every Target reachable through an effect's
// exported fields, descending into nested effects (a Sequence's children, a
// Conditional's branches). It reads targets by reflection rather than a
// per-effect switch so a new effect that carries a Target is covered for free;
// an unexported field holding a Target is skipped, not read, so the walk never
// panics on one.
func walkEffectTargets(e engine.Effect, fn func(engine.Target)) {
	walkForTargets(reflect.ValueOf(e), fn)
}

func walkForTargets(v reflect.Value, fn func(engine.Target)) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
		if !v.IsNil() {
			walkForTargets(v.Elem(), fn)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			walkForTargets(v.Index(i), fn)
		}
	case reflect.Struct:
		if v.Type() == targetType {
			if v.CanInterface() {
				fn(v.Interface().(engine.Target))
			}
			return
		}
		for i := 0; i < v.NumField(); i++ {
			walkForTargets(v.Field(i), fn)
		}
	}
}

// defTargetsPhrase reports whether a card has an ability effect that targets the
// given noun phrase, the predicate the target-shape specimens select on.
func defTargetsPhrase(d *engine.CardDefinition, phrase string) bool {
	found := false
	for _, ab := range d.Abilities {
		walkEffectTargets(ab.Effect, func(t engine.Target) {
			if t.Text() == phrase {
				found = true
			}
		})
	}
	return found
}

// targetShapeSpecimens returns one real card per distinct target noun phrase the
// card pool actually renders — "an enemy creature", "each friendly creature", and
// the rest — gathered from effect targets rather than a list written here, so the
// gallery shows exactly the target phrasings the client has to draw. Every phrase
// comes from a real card, so none is a gap.
func targetShapeSpecimens() []specimen {
	all := cards.All()
	phraseSet := map[string]bool{}
	for i := range all {
		for _, ab := range all[i].Abilities {
			walkEffectTargets(ab.Effect, func(t engine.Target) {
				if p := t.Text(); p != "" {
					phraseSet[p] = true
				}
			})
		}
	}
	phrases := make([]string, 0, len(phraseSet))
	for p := range phraseSet {
		phrases = append(phrases, p)
	}
	sort.Strings(phrases)
	out := make([]specimen, 0, len(phrases))
	for _, p := range phrases {
		p := p
		out = append(out, randomMatch(p, func(d *engine.CardDefinition) bool {
			return defTargetsPhrase(d, p)
		}))
	}
	return out
}

// textCombiner is a constructed noun phrase that isolates one text-helper axis,
// so the article and quantifier inflections read as minimal pairs side by side
// rather than being hunted for across cards.
type textCombiner struct {
	Caption string
	Phrase  string
}

// combinerRows renders a few bare Targets to expose the text helpers on purpose:
// the a/an article turning on the following vowel, and the single-vs-collective
// quantifier. The phrases are built here, not matched from cards, because the
// point is the helper, not any one card.
func combinerRows() []textCombiner {
	return []textCombiner{
		{"Article — consonant", engine.Target{Kind: engine.TargetChosenFriendlyCreature}.Text()},
		{"Article — vowel", engine.Target{Kind: engine.TargetChosenEnemyCreature}.Text()},
		{"Chosen — single", engine.Target{Kind: engine.TargetChosenCreature}.Text()},
		{"Each — collective", engine.Target{Kind: engine.TargetEachCreature}.Text()},
	}
}

// combinerList draws the combiner minimal pairs as caption/phrase rows.
func combinerList() app.UI {
	rows := combinerRows()
	items := make([]app.UI, 0, len(rows))
	for _, r := range rows {
		items = append(items, app.Div().Class("style-combiner").Body(
			app.Span().Class("style-combiner-caption").Text(r.Caption),
			app.Span().Class("style-combiner-phrase").Text(r.Phrase),
		))
	}
	return app.Div().Class("style-combiners").Body(items...)
}

// cardTextSection shows one card per text feature: the triggers a card can fire
// on, then the continuous abilities it can carry. Each specimen is a real card
// the query matched, so the section shows how the client actually renders that
// feature's rules text; an empty query is drawn as a gap rather than skipped.
func (s *style) cardTextSection() app.UI {
	return app.Div().Body(
		app.P().Class("style-note").Body(
			app.Text("One card per text feature, each a random real card the query "+
				"matched, reshuffled on reload. Terms are defined in the "),
			app.A().Class("style-link").Href("/rulebook").Text("rulebook"),
			app.Text("."),
		),
		app.H3().Class("style-h3").Text("Triggers"),
		s.specimenRow(triggerSpecimens()),
		app.H3().Class("style-h3").Text("Continuous abilities"),
		s.specimenRow(continuousSpecimens()),
		app.H3().Class("style-h3").Text("Target shapes"),
		s.specimenRow(targetShapeSpecimens()),
		app.H3().Class("style-h3").Text("Combiners"),
		combinerList(),
	)
}
