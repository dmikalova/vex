package web

import (
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []queryTerm
	}{
		{"empty", "   ", nil},
		{"single", "damage", []queryTerm{{alts: []string{"damage"}}}},
		{
			"and",
			"deal damage",
			[]queryTerm{{alts: []string{"deal"}}, {alts: []string{"damage"}}},
		},
		{"or", "mars|logos", []queryTerm{{alts: []string{"mars", "logos"}}}},
		{"exclude", "-token", []queryTerm{{negate: true, alts: []string{"token"}}}},
		{
			"phrase",
			`"deal 2 damage"`,
			[]queryTerm{{alts: []string{"deal 2 damage"}}},
		},
		{
			"escaped quote",
			`\"`,
			[]queryTerm{{alts: []string{`"`}}},
		},
		{
			"mixed",
			`gain -chains "each enemy"|"each friendly"`,
			[]queryTerm{
				{alts: []string{"gain"}},
				{negate: true, alts: []string{"chains"}},
				{alts: []string{"each enemy", "each friendly"}},
			},
		},
		{"lone dash dropped", "-", nil},
		{"case folded", "DaMaGe", []queryTerm{{alts: []string{"damage"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseQuery(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseQuery(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestMatchesQuery(t *testing.T) {
	hay := "Play: Deal 2 damage to each enemy creature. Gain 1 Æmber."
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{"empty matches", "", true},
		{"substring", "damage", true},
		{"case insensitive", "DAMAGE", true},
		{"and both present", "damage aember", true},
		{"and one missing", "damage chains", false},
		{"or one present", "chains|damage", true},
		{"or none present", "chains|heal", false},
		{"phrase present", `"each enemy"`, true},
		{"phrase absent", `"each friendly"`, false},
		{"exclude present fails", "-damage", false},
		{"exclude absent passes", "-chains", true},
		{"combo", `damage -heal "each enemy"`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesQuery(parseQuery(tt.query), hay); got != tt.want {
				t.Errorf("matchesQuery(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}
