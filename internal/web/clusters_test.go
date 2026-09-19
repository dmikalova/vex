package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// TestBuildClusterDataGroupsBothMechanisms checks buildClusterData reads both
// cluster mechanisms off the registry: named clusters grouped by membership (with
// the cross-set Shard cluster among them, every cluster holding members) and the
// filtered pulls their leads declare.
func TestBuildClusterDataGroupsBothMechanisms(t *testing.T) {
	named, filtered := buildClusterData()
	if len(named) == 0 {
		t.Fatal("no named clusters found")
	}
	foundShard := false
	for _, r := range named {
		if len(r.members) == 0 {
			t.Errorf("named cluster %q has no members", r.name)
		}
		if r.name == "Shard" {
			foundShard = true
		}
	}
	if !foundShard {
		t.Error("expected the Shard cluster among the named clusters")
	}
	if len(filtered) == 0 {
		t.Fatal("no filtered pulls found")
	}
	for _, r := range filtered {
		if r.leadDef == nil {
			t.Errorf("filtered pull %q has no lead", r.name)
		}
	}
}

// TestClustersRenderWithoutPanic draws the page through its production Render with
// the real registry data behind it, the same faces the gallery draws, so a change
// to the printed face reaches it too.
func TestClustersRenderWithoutPanic(t *testing.T) {
	named, filtered := buildClusterData()
	c := &clusters{
		named:    named,
		filtered: filtered,
	}
	html := app.HTMLString(c.Render())
	if !strings.Contains(html, "Named clusters") || !strings.Contains(html, "Filtered pulls") {
		t.Error("clusters page did not render its section headings")
	}
}
