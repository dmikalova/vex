package web

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/deckgen"
	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is the /clusters page: an internal data view (like /style, gated
// behind styleEnabled) of the card families deck generation places together
// (ADR 0036). It is not a player surface — the axis is the deck-generation data,
// not the presentation — so it reads the two cluster mechanisms straight off the
// registry and lays them out together: the named clusters (grouped by
// GenerationProfile.Cluster) and the filtered pulls (GenerationProfile.Leads),
// each with its strategy, trigger, lead, spawn parameters, and member faces, so
// every connected family can be reviewed in one place.

// NewClusters returns the root component for the /clusters page.
func NewClusters() app.Composer { return &clusters{} }

// clusterMember is one card in a named cluster, with its role: whether it is the
// lead, and — for a Pull cluster — its own per-partner pull rate (min copies and
// Poisson mean), which the other strategies share at the cluster level.
type clusterMember struct {
	def  *engine.CardDefinition
	lead bool
	min  int
	mean float64
}

// clusterRow is a named cluster resolved to its members: the strategy that fills
// it, the trigger that fires it, the cluster-level spawn rate, and its members.
type clusterRow struct {
	name     string
	strategy deckgen.ClusterStrategy
	trigger  deckgen.ClusterTrigger
	min, max int
	mean     float64
	members  []clusterMember
}

// filteredRow is a deck-wide filtered pull: the lead card whose presence fires
// it, the floor and mean it guarantees, and the catalog cards whose definitions
// its predicate matches (the pull's possible targets).
type filteredRow struct {
	name    string
	leadDef *engine.CardDefinition
	floor   int
	mean    float64
	matches []*engine.CardDefinition
}

// clusters is the /clusters page component. It reads the cluster data off the
// registry once on mount, then renders the named clusters and filtered pulls.
type clusters struct {
	app.Compo

	named    []clusterRow
	filtered []filteredRow
}

// OnMount builds the cluster tables from the registry. The data is static, so it
// is read once here rather than on every render.
func (c *clusters) OnMount(app.Context) {
	c.named, c.filtered = buildClusterData()
}

// buildClusterData reads both cluster mechanisms off the registry: it groups the
// registered cards by named-cluster membership, and collects the filtered pulls
// their leads declare, resolving each pull's matches against the whole catalog.
func buildClusterData() ([]clusterRow, []filteredRow) {
	regs := card.Cards()
	defs := make([]engine.CardDefinition, len(regs))
	for i := range regs {
		defs[i] = regs[i].Def
	}

	byName := map[string]*clusterRow{}
	var order []string
	type lead struct {
		def *engine.CardDefinition
		fc  *deckgen.FilteredCluster
	}
	var leads []lead
	for i := range regs {
		rc := regs[i]
		if m := rc.Profile.Cluster; !m.Empty() {
			row := byName[m.Name]
			if row == nil {
				row = &clusterRow{
					name:     m.Name,
					strategy: m.Strategy,
					trigger:  m.Trigger,
				}
				byName[m.Name] = row
				order = append(order, m.Name)
			}
			def := rc.Def
			row.members = append(row.members, clusterMember{
				def:  &def,
				lead: m.Lead,
				min:  m.Min,
				mean: m.Mean,
			})
			// The lead carries the cluster-level rate; a non-lead only overrides it
			// when the cluster has no lead (RandomCount, SelfPull, OnePerHouse).
			if m.Lead || (row.max == 0 && row.mean == 0 && row.min == 0) {
				row.min, row.max, row.mean = m.Min, m.Max, m.Mean
			}
		}
		if rc.Profile.Leads != nil {
			def := rc.Def
			leads = append(leads, lead{
				def: &def,
				fc:  rc.Profile.Leads,
			})
		}
	}

	named := make([]clusterRow, 0, len(order))
	for _, name := range order {
		row := byName[name]
		sort.SliceStable(row.members, func(i, j int) bool {
			if row.members[i].lead != row.members[j].lead {
				return row.members[i].lead
			}
			return row.members[i].def.Name < row.members[j].def.Name
		})
		named = append(named, *row)
	}
	sort.SliceStable(named, func(i, j int) bool { return named[i].name < named[j].name })

	filtered := make([]filteredRow, 0, len(leads))
	for _, l := range leads {
		row := filteredRow{
			name:    l.fc.Name,
			leadDef: l.def,
			floor:   l.fc.Floor,
			mean:    l.fc.Mean,
		}
		for i := range defs {
			if l.fc.Match == nil || defs[i].Name == l.def.Name || !l.fc.Match(defs[i]) {
				continue
			}
			def := defs[i]
			row.matches = append(row.matches, &def)
		}
		sort.SliceStable(row.matches, func(i, j int) bool {
			return row.matches[i].Name < row.matches[j].Name
		})
		filtered = append(filtered, row)
	}
	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].name < filtered[j].name })

	return named, filtered
}

// clusterStrategyLabel names a fill strategy for display.
func clusterStrategyLabel(s deckgen.ClusterStrategy) string {
	switch s {
	case deckgen.WholePool:
		return "Whole pool"
	case deckgen.RandomCount:
		return "Random count"
	case deckgen.SelfPull:
		return "Self pull"
	case deckgen.PullExact:
		return "Pull exact"
	case deckgen.Pull:
		return "Pull"
	case deckgen.OnePerHouse:
		return "One per House"
	case deckgen.PerGigantic:
		return "Per gigantic"
	}
	return "—"
}

// clusterTriggerLabel names a trigger for display.
func clusterTriggerLabel(t deckgen.ClusterTrigger) string {
	switch t {
	case deckgen.ByLead:
		return "By lead"
	case deckgen.ByAnyMember:
		return "By any member"
	}
	return "—"
}

// clusterSpawnSummary describes, in one line, how many cards a cluster places and
// at what rate — the strategy-specific reading of its Min/Max/Mean.
func clusterSpawnSummary(r clusterRow) string {
	switch r.strategy {
	case deckgen.WholePool:
		return fmt.Sprintf("places all %d members", len(r.members))
	case deckgen.RandomCount:
		return fmt.Sprintf("places %d–%d members", r.min, r.max)
	case deckgen.SelfPull:
		return fmt.Sprintf("≥%d copies, ~%s average (cap 12)", r.min, formatMean(r.mean))
	case deckgen.PullExact:
		return "one of each non-lead per lead in the pod"
	case deckgen.Pull:
		return "per-partner pull (see member rates)"
	case deckgen.OnePerHouse:
		return "one member in each of the deck's Houses"
	case deckgen.PerGigantic:
		return "one member per gigantic base"
	}
	return ""
}

// formatMean renders a Poisson mean without trailing zeros.
func formatMean(m float64) string { return strconv.FormatFloat(m, 'g', -1, 64) }

// Render draws the two cluster tables. printedFace needs the icon-outline filter
// the game frame normally supplies, so it is injected here as the gallery does.
func (c *clusters) Render() app.UI {
	return app.Div().Class("clusters-page").Body(
		app.Raw(iconOutlineFilter),
		app.Header().Class("clusters-header").Body(
			app.Span().Class("clusters-brand").Text("Vactrol"),
			app.Span().Class("clusters-title").Text("Clusters"),
			app.Nav().Class("clusters-nav").Body(
				app.A().Class("clusters-nav-link").Href("/cards").Text("Cards"),
				app.A().Class("clusters-nav-link").Href("/style").Text("Style"),
			),
		),
		app.Main().Class("clusters-main").Body(
			app.P().Class("clusters-overview").Text(fmt.Sprintf(
				"%d named clusters · %d filtered pulls",
				len(c.named), len(c.filtered))),
			app.Section().Class("clusters-section").Body(
				app.H2().Class("clusters-h2").Text("Named clusters"),
				c.namedList(),
			),
			app.Section().Class("clusters-section").Body(
				app.H2().Class("clusters-h2").Text("Filtered pulls"),
				c.filteredList(),
			),
		),
	)
}

// namedList draws every named cluster as its own card.
func (c *clusters) namedList() app.UI {
	rows := make([]app.UI, 0, len(c.named))
	for _, r := range c.named {
		rows = append(rows, c.namedCard(r))
	}
	return app.Div().Class("cluster-list").Body(rows...)
}

// namedCard draws one named cluster: its header of strategy/trigger/spawn, then
// its member faces annotated with their roles.
func (c *clusters) namedCard(r clusterRow) app.UI {
	members := make([]app.UI, 0, len(r.members))
	for _, m := range r.members {
		members = append(members, c.memberCell(r, m))
	}
	return app.Div().Class("cluster").Body(
		app.Div().Class("cluster-head").Body(
			app.H3().Class("cluster-name").Text(r.name),
			app.Div().Class("cluster-meta").Body(
				metaTag("Strategy", clusterStrategyLabel(r.strategy)),
				metaTag("Trigger", clusterTriggerLabel(r.trigger)),
				metaTag("Spawn", clusterSpawnSummary(r)),
				metaTag("Members", strconv.Itoa(len(r.members))),
			),
		),
		app.Div().Class("cluster-members").Body(members...),
	)
}

// memberCell draws one member's face and the annotations its role adds: a Lead
// badge, and — for a Pull cluster — its own per-partner rate.
func (c *clusters) memberCell(r clusterRow, m clusterMember) app.UI {
	badges := make([]app.UI, 0, 2)
	if m.lead {
		badges = append(badges, app.Span().
			Class("cluster-badge", "cluster-badge--lead").Text("Lead"))
	}
	if r.strategy == deckgen.Pull && !m.lead {
		badges = append(badges, app.Span().Class("cluster-badge").
			Text(fmt.Sprintf("≥%d · ~%s", m.min, formatMean(m.mean))))
	}
	return app.Div().Class("cluster-member").Body(
		printedFace(m.def),
		app.Div().Class("cluster-member-meta").Body(badges...),
	)
}

// filteredList draws every filtered pull as its own card.
func (c *clusters) filteredList() app.UI {
	rows := make([]app.UI, 0, len(c.filtered))
	for _, r := range c.filtered {
		rows = append(rows, c.filteredCard(r))
	}
	return app.Div().Class("cluster-list").Body(rows...)
}

// filteredCard draws one filtered pull: its header, the lead face, then a face
// for every catalog card its predicate can pull.
func (c *clusters) filteredCard(r filteredRow) app.UI {
	members := make([]app.UI, 0, len(r.matches)+1)
	members = append(members, app.Div().Class("cluster-member").Body(
		printedFace(r.leadDef),
		app.Div().Class("cluster-member-meta").Body(
			app.Span().Class("cluster-badge", "cluster-badge--lead").Text("Lead"),
		),
	))
	for _, def := range r.matches {
		members = append(members, app.Div().Class("cluster-member").Body(printedFace(def)))
	}
	return app.Div().Class("cluster").Body(
		app.Div().Class("cluster-head").Body(
			app.H3().Class("cluster-name").Text(r.name),
			app.Div().Class("cluster-meta").Body(
				metaTag("Lead", r.leadDef.Name),
				metaTag("Floor", strconv.Itoa(r.floor)),
				metaTag("Mean", formatMean(r.mean)),
				metaTag("Matches", strconv.Itoa(len(r.matches))),
			),
		),
		app.Div().Class("cluster-members").Body(members...),
	)
}

// metaTag draws one labelled value in a cluster's header.
func metaTag(label, value string) app.UI {
	return app.Span().Class("cluster-tag").Body(
		app.Span().Class("cluster-tag-label").Text(label+": "),
		app.Span().Class("cluster-tag-value").Text(value),
	)
}
