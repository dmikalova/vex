package cardtest

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests assert the harness's own failure reporting: the paths a passing
// card test never walks, which is exactly why they need driving from outside. A
// harness that stopped failing would let every card test pass on a broken engine,
// so what is checked here is that each misuse is caught and named.

// recordingT stands in for the *testing.T a card test would pass, capturing what
// the harness reports instead of failing the run. Fatalf and Fatal end the
// calling goroutine as testing.T does, so the harness's control flow after a
// fatal is the same here as in a real test.
type recordingT struct {
	testing.TB
	failures []string
	cleanups []func()
}

func (r *recordingT) Helper() {}

func (r *recordingT) Errorf(format string, args ...any) {
	r.failures = append(r.failures, fmt.Sprintf(format, args...))
}

func (r *recordingT) Fatalf(format string, args ...any) {
	r.Errorf(format, args...)
	runtime.Goexit()
}

func (r *recordingT) Fatal(args ...any) {
	r.failures = append(r.failures, fmt.Sprint(args...))
	runtime.Goexit()
}

func (r *recordingT) Cleanup(f func()) { r.cleanups = append(r.cleanups, f) }

// expectFail runs fn against a recording TB and returns everything it reported,
// failing t if the harness let the misuse through. fn runs on its own goroutine
// because a fatal ends the goroutine that hits it, and the registered cleanups
// run after it the way the testing package would run them.
func expectFail(t *testing.T, fn func(tb testing.TB)) string {
	t.Helper()
	rec := &recordingT{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			for i := len(rec.cleanups) - 1; i >= 0; i-- {
				rec.cleanups[i]()
			}
		}()
		fn(rec)
	}()
	<-done
	if len(rec.failures) == 0 {
		t.Fatal("the harness accepted a misuse it should have reported")
	}
	return strings.Join(rec.failures, "\n")
}

// wants fails unless the report contains every fragment.
func wants(t *testing.T, got string, fragments ...string) {
	t.Helper()
	for _, f := range fragments {
		if !strings.Contains(got, f) {
			t.Errorf("report %q does not mention %q", got, f)
		}
	}
}

func TestZoneString(t *testing.T) {
	for _, tc := range []struct {
		z    Zone
		want string
	}{
		{PlayArea, "play area"},
		{Hand, "hand"},
		{Discard, "discard"},
		{Archives, "archives"},
		{Deck, "deck"},
		{Attached, "attached"},
		{Purge, "purge"},
		{Gone, "gone"},
	} {
		if got := tc.z.String(); got != tc.want {
			t.Errorf("Zone(%d) = %q, want %q", tc.z, got, tc.want)
		}
	}
}

// An unbound handle is the mistake of declaring a ct.Card but never passing it to
// ct.Bind, so every read through it names that omission rather than panicking.
func TestUnboundHandle(t *testing.T) {
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, Setup{})
		Card{h: h}.require()
	})
	wants(t, got, "not bound", "ct.Bind")
}

// Every CardExpect assertion reports with Errorf rather than Fatalf, so one pass
// over a card that matches none of them collects the whole set.
func TestExpectMismatches(t *testing.T) {
	var troll Card
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, Setup{P1: Side{InPlay: []Entry{Bind(&troll, Creature(Power(5), Armor(1)))}}})
		h.Expect(troll).
			Damage(3).
			Power(9).
			Armor(4).
			AmberOn(2).
			Exhausted().
			Stunned(true).
			At(Discard)
		troll.Exhaust()
		h.Expect(troll).Ready()
	})
	wants(t, got,
		"damage = 0, want 3",
		"power = 5, want 9",
		"armor = 1, want 4",
		"Æmber-on-card = 0, want 2",
		"is ready, want exhausted",
		"stunned = false, want true",
		"is in play area, want discard",
		"is exhausted, want ready",
	)
}

func TestExpectPoolMismatches(t *testing.T) {
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, Setup{P1: Side{Amber: 2, Keys: 1}})
		h.P1.ExpectAmber(5)
		h.P1.ExpectKeys(3)
	})
	wants(t, got, "P1 Æmber = 2, want 5", "P1 keys = 1, want 3")
}

// A card whose type the harness cannot dispatch on is a harness gap, not a card
// bug, so it says so with the type it was handed.
func TestPlayUnknownType(t *testing.T) {
	got := expectFail(t, func(tb testing.TB) {
		odd := engine.NewCard("Oddity", DefaultHouse, engine.CardType(99), engine.Common)
		h := Play(tb, Setup{P1: Side{Hand: Cards(odd)}})
		h.P1.Play(odd)
	})
	wants(t, got, "unknown type")
}

func TestFindInHandFailures(t *testing.T) {
	dup := Creature()
	var onBoard Card
	for _, tc := range []struct {
		name string
		play func(h *Harness)
		want []string
	}{
		{
			name: "handle for a card that is not in hand",
			play: func(h *Harness) { h.P1.Play(onBoard) },
			want: []string{"is not in P1's hand"},
		},
		{
			name: "two copies named by definition",
			play: func(h *Harness) { h.P1.Play(dup) },
			want: []string{"2 copies", "ct.Bind"},
		},
		{
			name: "definition that is not in hand",
			play: func(h *Harness) { h.P1.Play(Creature()) },
			want: []string{"not in P1's hand"},
		},
		{
			name: "something that is not a card at all",
			play: func(h *Harness) { h.P1.Play("Troll") },
			want: []string{"cannot play a string"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := expectFail(t, func(tb testing.TB) {
				tc.play(Play(tb, Setup{P1: Side{
					InPlay: []Entry{Bind(&onBoard, Creature())},
					Hand:   Cards(dup, dup),
				}}))
			})
			wants(t, got, tc.want...)
		})
	}
}

func TestResolveFailures(t *testing.T) {
	dup := Creature()
	for _, tc := range []struct {
		name string
		use  func(h *Harness)
		want []string
	}{
		{
			name: "two matching cards in play",
			use:  func(h *Harness) { h.P1.Reap(dup) },
			want: []string{"Reap", "2 matching cards", "ct.Bind"},
		},
		{
			name: "no matching card in play",
			use:  func(h *Harness) { h.P1.Reap(Creature()) },
			want: []string{"no matching card found"},
		},
		{
			name: "something that is not a card at all",
			use:  func(h *Harness) { h.P1.Reap(42) },
			want: []string{"cannot resolve a int"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := expectFail(t, func(tb testing.TB) {
				tc.use(Play(tb, Setup{P1: Side{InPlay: Cards(dup, dup)}}))
			})
			wants(t, got, tc.want...)
		})
	}
}

// The "cannot" assertions exist for cards that gate themselves; each fails when
// the card is in fact wide open.
func TestExpectCannotWhenItCan(t *testing.T) {
	var troll, tactic Card
	setup := Setup{P1: Side{
		InPlay: []Entry{Bind(&troll, Creature())},
		Hand:   []Entry{Bind(&tactic, Tactic())},
	}}
	for _, tc := range []struct {
		name string
		call func(h *Harness)
		want string
	}{
		{
			name: "play",
			call: func(h *Harness) { h.P1.ExpectCannotPlay(tactic) },
			want: "may be played, want it blocked",
		},
		{
			name: "use",
			call: func(h *Harness) { h.P1.ExpectCannotUse(troll) },
			want: "may be used, want it blocked",
		},
		{
			name: "use to reap",
			call: func(h *Harness) { h.P1.ExpectCannotUseTo(troll, engine.ReapUse) },
			want: "may be used to",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := expectFail(t, func(tb testing.TB) { tc.call(Play(tb, setup)) })
			wants(t, got, tc.want)
		})
	}
}

func TestToEntryAndToUpgradePanic(t *testing.T) {
	for _, tc := range []struct {
		name string
		call func()
		want string
	}{
		{"entry", func() { toEntry("Troll") }, "as a card entry"},
		{"upgrade", func() { toUpgrade(7) }, "as an upgrade"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				got, ok := recover().(string)
				if !ok {
					t.Fatal("no panic")
				}
				wants(t, got, tc.want)
			}()
			tc.call()
		})
	}
}

// The vanilla builders' remaining knobs: a card no card test has needed yet is
// still part of the harness's vocabulary, so it is checked here.
func TestVanillaOptions(t *testing.T) {
	pip := Tactic(AemberBonus(2))
	if pip.AemberBonus != 2 {
		t.Errorf("AemberBonus = %d, want 2", pip.AemberBonus)
	}
	up := Upgrade(PowerBonus(3), ArmorBonus(1))
	if up.Type != engine.Upgrade {
		t.Errorf("Upgrade type = %v, want %v", up.Type, engine.Upgrade)
	}
	if up.Static.PowerBonus != 3 || up.Static.ArmorBonus != 1 {
		t.Errorf("Upgrade static = %+v, want power 3 armor 1", up.Static)
	}
}
