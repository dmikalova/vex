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

// targetType and durationType are the reflect types the effect walk collects.
// conditionType and countType are the two interface types it collects instead by
// implementation: unlike Target and Duration they are not one concrete type, so
// the walk matches any value declared as the interface (a Conditional's Cond, a
// scaled effect's Per) and reads its own CondText/CountText.
var (
	targetType    = reflect.TypeFor[engine.Target]()
	durationType  = reflect.TypeFor[engine.Duration]()
	conditionType = reflect.TypeFor[engine.Condition]()
	countType     = reflect.TypeFor[engine.Count]()
)

// walkEffectValues calls visit for every value reachable through an effect's
// exported fields, descending into nested effects (a Sequence's children, a
// Conditional's branches), slices, and pointers. It reads the tree by reflection
// rather than a per-effect switch so a new effect that carries a Target or a
// Duration is covered for free. It stops at a Target so its internal filter
// fields are not walked, and it never reads an unexported field, so the walk
// never panics on one.
func walkEffectValues(v reflect.Value, visit func(reflect.Value)) {
	if !v.IsValid() {
		return
	}
	visit(v)
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer:
		if !v.IsNil() {
			walkEffectValues(v.Elem(), visit)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			walkEffectValues(v.Index(i), visit)
		}
	case reflect.Struct:
		if v.Type() == targetType {
			return
		}
		for _, field := range v.Fields() {
			walkEffectValues(field, visit)
		}
	}
}

// walkEffectTargets calls fn for every Target reachable through an effect.
func walkEffectTargets(e engine.Effect, fn func(engine.Target)) {
	walkEffectValues(reflect.ValueOf(e), func(v reflect.Value) {
		if v.Type() == targetType && v.CanInterface() {
			if t, ok := v.Interface().(engine.Target); ok {
				fn(t)
			}
		}
	})
}

// walkEffectDurations calls fn for every Duration reachable through an effect.
func walkEffectDurations(e engine.Effect, fn func(engine.Duration)) {
	walkEffectValues(reflect.ValueOf(e), func(v reflect.Value) {
		if v.Type() == durationType && v.CanInterface() {
			if d, ok := v.Interface().(engine.Duration); ok {
				fn(d)
			}
		}
	})
}

// walkEffectConditions calls fn for every Condition an effect gates on — a
// Conditional's Cond and any condition nested inside it — matched by the field's
// interface type so a new condition is covered without listing it here.
func walkEffectConditions(e engine.Effect, fn func(engine.Condition)) {
	walkEffectValues(reflect.ValueOf(e), func(v reflect.Value) {
		if v.Kind() == reflect.Interface && v.Type() == conditionType &&
			!v.IsNil() && v.CanInterface() {
			if c, ok := v.Interface().(engine.Condition); ok {
				fn(c)
			}
		}
	})
}

// walkEffectCounts calls fn for every Count an effect scales by — a Per or Times
// clause — matched by the field's interface type. A Fixed count prints nothing
// (its magnitude shows as a plain number), so callers drop the empty CountText.
func walkEffectCounts(e engine.Effect, fn func(engine.Count)) {
	walkEffectValues(reflect.ValueOf(e), func(v reflect.Value) {
		if v.Kind() == reflect.Interface && v.Type() == countType &&
			!v.IsNil() && v.CanInterface() {
			if n, ok := v.Interface().(engine.Count); ok {
				fn(n)
			}
		}
	})
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
	return phraseSpecimens(func(e engine.Effect, add func(string)) {
		walkEffectTargets(e, func(t engine.Target) { add(t.Text()) })
	}, defTargetsPhrase)
}

// defUsesDuration reports whether a card has an ability effect carrying the given
// timing window, the predicate the duration specimens select on.
func defUsesDuration(d *engine.CardDefinition, dur engine.Duration) bool {
	found := false
	for _, ab := range d.Abilities {
		walkEffectDurations(ab.Effect, func(got engine.Duration) {
			if got == dur {
				found = true
			}
		})
	}
	return found
}

// durationSpecimens returns one real card per timing window a timed effect can
// take, ranging over engine.Durations() rather than a list written here so a new
// duration shows up as a new specimen. A duration no loaded card carries yet is
// drawn as a gap. The caption names the window (Duration.String()); the card
// itself shows how that window is actually phrased in context, which varies by
// effect.
func durationSpecimens() []specimen {
	out := make([]specimen, 0, len(engine.Durations()))
	for _, dur := range engine.Durations() {
		out = append(out, randomMatch(dur.String(), func(d *engine.CardDefinition) bool {
			return defUsesDuration(d, dur)
		}))
	}
	return out
}

// defHasConditionPhrase reports whether a card gates an ability on the given
// rendered condition clause, the predicate the condition specimens select on.
func defHasConditionPhrase(d *engine.CardDefinition, phrase string) bool {
	found := false
	for _, ab := range d.Abilities {
		walkEffectConditions(ab.Effect, func(c engine.Condition) {
			if c.CondText() == phrase {
				found = true
			}
		})
	}
	return found
}

// conditionSpecimens returns one real card per distinct "if ..." clause the pool
// gates on, gathered from effect conditions rather than a list written here, so
// the gallery shows exactly the condition phrasings the client draws. Every
// phrase comes from a real card, so none is a gap.
func conditionSpecimens() []specimen {
	return phraseSpecimens(func(e engine.Effect, add func(string)) {
		walkEffectConditions(e, func(c engine.Condition) { add(c.CondText()) })
	}, defHasConditionPhrase)
}

// defHasCountPhrase reports whether a card scales an effect by the given "for
// each ..." clause, the predicate the count specimens select on.
func defHasCountPhrase(d *engine.CardDefinition, phrase string) bool {
	found := false
	for _, ab := range d.Abilities {
		walkEffectCounts(ab.Effect, func(c engine.Count) {
			if c.CountText() == phrase {
				found = true
			}
		})
	}
	return found
}

// countSpecimens returns one real card per distinct "for each ..." clause the
// pool scales by, gathered from effect counts. A Fixed count renders no clause
// and is dropped, so every specimen shown is a scaling count carried by a real
// card.
func countSpecimens() []specimen {
	return phraseSpecimens(func(e engine.Effect, add func(string)) {
		walkEffectCounts(e, func(c engine.Count) { add(c.CountText()) })
	}, defHasCountPhrase)
}

// phraseSpecimens is the shared shape of the target/condition/count galleries:
// gather every non-empty phrase the pool renders through gather, then draw one
// random real card per phrase via has. It keeps the three data-driven sections
// from repeating the same collect-sort-match loop.
func phraseSpecimens(
	gather func(engine.Effect, func(string)),
	has func(*engine.CardDefinition, string) bool,
) []specimen {
	all := cards.All()
	phraseSet := map[string]bool{}
	for i := range all {
		for _, ab := range all[i].Abilities {
			gather(ab.Effect, func(p string) {
				if p != "" {
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
		out = append(out, randomMatch(p, func(d *engine.CardDefinition) bool {
			return has(d, p)
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

// cardTextSpecimens holds the Card text section's four specimen rows. They are
// queried once per page load and cached on the style component, because the
// target-shape and duration rows walk every card's effect tree by reflection and
// the page re-renders in full on every hover.
type cardTextSpecimens struct {
	triggers, continuous, targetShapes, durations, conditions, counts []specimen
}

// cardTextSpecs returns the cached specimens, computing them on first use. The
// queries depend only on styleSeed, which is fixed for the life of a page load,
// so one computation serves every re-render.
func (s *style) cardTextSpecs() *cardTextSpecimens {
	if s.cardText == nil {
		s.cardText = &cardTextSpecimens{
			triggers:     triggerSpecimens(),
			continuous:   continuousSpecimens(),
			targetShapes: targetShapeSpecimens(),
			durations:    durationSpecimens(),
			conditions:   conditionSpecimens(),
			counts:       countSpecimens(),
		}
	}
	return s.cardText
}

// cardTextSection shows one card per text feature: the triggers a card can fire
// on, then the continuous abilities it can carry. Each specimen is a real card
// the query matched, so the section shows how the client actually renders that
// feature's rules text; an empty query is drawn as a gap rather than skipped.
func (s *style) cardTextSection() app.UI {
	specs := s.cardTextSpecs()
	return app.Div().Body(
		app.P().Class("style-note").Body(
			app.Text("One card per text feature, each a random real card the query "+
				"matched, reshuffled on reload. Terms are defined in the "),
			app.A().Class("style-link").Href("/rulebook").Text("rulebook"),
			app.Text("."),
		),
		app.H3().Class("style-h3").Text("Triggers"),
		s.specimenRow(specs.triggers),
		app.H3().Class("style-h3").Text("Continuous abilities"),
		s.specimenRow(specs.continuous),
		app.H3().Class("style-h3").Text("Target shapes"),
		s.specimenRow(specs.targetShapes),
		app.H3().Class("style-h3").Text("Durations"),
		s.specimenRow(specs.durations),
		app.H3().Class("style-h3").Text("Conditions"),
		s.specimenRow(specs.conditions),
		app.H3().Class("style-h3").Text("Counts"),
		s.specimenRow(specs.counts),
		app.H3().Class("style-h3").Text("Combiners"),
		combinerList(),
	)
}
