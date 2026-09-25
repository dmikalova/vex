package main

import (
	"bytes"
	jsontext "encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/dmikalova/vex/internal/cards/provenance"
)

// importSet pairs a source set with its Master Vault expansion id. There is no
// live expansions endpoint (the decks feed links only houses), so the ids are
// pinned here from decksofkeyforge.com's expansion table. Order is release order,
// the order `import-provenance all` walks the sets.
type importSet struct {
	Set       provenance.SourceSet
	Expansion int
}

// importSets is every set the importer can build a catalog for, with its Master
// Vault expansion id.
var importSets = []importSet{
	{provenance.CallOfTheArchons, 341},
	{provenance.AgeOfAscension, 435},
	{provenance.WorldsCollide, 452},
	{provenance.AnomalyExpansion, 453},
	{provenance.MassMutation, 479},
	{provenance.DarkTidings, 496},
	{provenance.WindsOfExchange, 600},
	{provenance.Unchained2022, 601},
	{provenance.VaultMasters2023, 609},
	{provenance.GrimReminders, 700},
	{provenance.Menagerie, 722},
	{provenance.VaultMasters2024, 737},
	{provenance.AemberSkies, 800},
	{provenance.TokensOfChange, 855},
	{provenance.MoreMutation, 874},
	{provenance.PropheticVisions, 886},
	{provenance.MartianCivilWar, 892},
	{provenance.Discovery, 907},
	{provenance.CrucibleClash, 918},
	{provenance.DraconianMeasures, 928},
	{provenance.VaultMasters2025, 939},
	{provenance.VaultMasters2026, 964},
}

// importSetBySlug returns the importSet for a slug, or false when unknown.
func importSetBySlug(slug string) (importSet, bool) {
	for _, is := range importSets {
		if is.Set.Slug == slug {
			return is, true
		}
	}
	return importSet{}, false
}

// importSlugs lists every importable set slug, for error messages.
func importSlugs() string {
	slugs := make([]string, len(importSets))
	for i, is := range importSets {
		slugs[i] = is.Set.Slug
	}
	return strings.Join(slugs, ", ")
}

// importProvenance builds a source catalog (internal/cards/provenance/<slug>.json)
// by paging the Master Vault decks feed for a set and folding each card it links
// into the catalog shape. With no slug it opens an interactive picker offering
// every set plus "All sets"; the slug "all" imports every set in release order.
// Each set is fetched on its own — there is no combined query.
func importProvenance(args []string) error {
	var slug string
	switch len(args) {
	case 0:
	case 1:
		slug = args[0]
	default:
		return fmt.Errorf("usage: cardlookup import-provenance [setSlug|all]")
	}
	if slug == "" {
		if isInteractive() {
			picked, err := pickImportSet()
			if err != nil {
				return err
			}
			if picked == "" {
				return nil // cancelled
			}
			slug = picked
		} else {
			return fmt.Errorf(
				"import-provenance needs a set slug or `all` when stdin is not a terminal",
			)
		}
	}

	targets := importSets
	if slug != "all" {
		is, ok := importSetBySlug(slug)
		if !ok {
			return fmt.Errorf("unknown set slug %q; known slugs: %s, all", slug, importSlugs())
		}
		targets = []importSet{is}
	}

	// A fresh connection per request (no keep-alive) mirrors curl: reused pooled
	// connections go half-dead under the feed's aggressive throttling and then
	// hang the next request until the client timeout.
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			TLSHandshakeTimeout: 15 * time.Second,
		},
	}
	// One pacer for the whole run so an `all` import carries the throttle thresholds
	// it learned from earlier sets into later ones, instead of every set restarting
	// at the polite base pace and re-provoking the feed's throttling from scratch.
	p := newPacer()
	for _, is := range targets {
		if err := importOneSet(client, p, is); err != nil {
			return fmt.Errorf("importing %s: %w", is.Set.Name, err)
		}
	}
	return nil
}

// pacer holds the throttle thresholds the importer learns from the feed. It lives
// for a whole import run, so during an `all` import each set resumes from the pace
// and backoff floor earlier sets settled on rather than relearning them.
type pacer struct {
	// pace is the delay between successive pages. It starts polite and ratchets up
	// one second per throttled group, settling below the throttle threshold; it
	// never eases back down.
	pace time.Duration
	// longest is the longest throttle wait seen so far, the floor a fresh throttle
	// resumes from instead of restarting at two seconds.
	longest time.Duration
}

const (
	pacerBasePace = time.Second
	pacerMaxPace  = 20 * time.Second
	pacerFloor    = 2 * time.Second
)

// newPacer returns a pacer at the polite starting thresholds.
func newPacer() *pacer {
	return &pacer{
		pace:    pacerBasePace,
		longest: pacerFloor,
	}
}

// throttled ratchets the inter-page pace up one second, capped at pacerMaxPace.
func (p *pacer) throttled() {
	p.pace += time.Second
	if p.pace > pacerMaxPace {
		p.pace = pacerMaxPace
	}
}

// importOneSet fetches, transforms, and writes one set's catalog.
func importOneSet(client *http.Client, p *pacer, is importSet) error {
	fmt.Printf("%s (expansion %d): fetching…\n", is.Set.Name, is.Expansion)
	raw, err := fetchSetCards(client, p, is.Set.Slug, is.Expansion)
	if err != nil {
		return err
	}
	cards := make([]catalogCard, 0, len(raw))
	for _, r := range raw {
		c, err := transformCard(r)
		if err != nil {
			return err
		}
		cards = append(cards, c)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].Number < cards[j].Number })

	path := filepath.Join("internal", "cards", "provenance", is.Set.Slug+".json")
	if err := writeCatalog(path, cards); err != nil {
		return err
	}
	fmt.Printf("%s: wrote %d cards to %s\n", is.Set.Name, len(cards), path)
	return nil
}

// --- Master Vault decks feed -------------------------------------------------

const mvDecksURL = "https://www.keyforgegame.com/api/decks/"

// mvResponse is the slice of the decks feed the importer reads: the total deck
// count and the card objects the page's decks link to (links=cards folds them
// into _linked.cards).
type mvResponse struct {
	Count  int `json:"count"`
	Linked struct {
		Cards []mvCard `json:"cards"`
	} `json:"_linked"`
}

// mvCard is the subset of a Master Vault card object the catalog needs. Numeric
// stats arrive as strings on some cards and as JSON numbers on others, so they
// are decoded through flexString.
type mvCard struct {
	CardTitle  string     `json:"card_title"`
	House      string     `json:"house"`
	CardType   string     `json:"card_type"`
	CardText   string     `json:"card_text"`
	Traits     string     `json:"traits"`
	Amber      flexString `json:"amber"`
	Power      flexString `json:"power"`
	Armor      flexString `json:"armor"`
	Rarity     string     `json:"rarity"`
	CardNumber string     `json:"card_number"`
	Expansion  int        `json:"expansion"`
	IsMaverick bool       `json:"is_maverick"`
	IsAnomaly  bool       `json:"is_anomaly"`
}

// flexString decodes a JSON value that may be either a string or a number into a
// string, since the decks feed spells the same stat both ways across cards.
type flexString string

func (f *flexString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = flexString(s)
		return nil
	}
	*f = flexString(strings.Trim(string(data), `"`))
	return nil
}

// fetchSetCards pages the decks feed for one expansion, collecting every distinct
// card it links, until the set is covered. It stops when the decks run out, or a
// run of pages adds nothing new: a run of 10 empty pages with no interior gap in
// the collected collector numbers (the set looks complete), or a hard cap of 100
// empty pages (a lingering gap that never fills). Maverick cards are skipped; a
// card belonging to another expansion (a cross-set maverick) is skipped too. The
// pacer carries the throttle thresholds across sets, so an `all` import does not
// relearn the feed's throttling for every set.
func fetchSetCards(client *http.Client, p *pacer, slug string, expansion int) ([]mvCard, error) {
	const pageSize = 25
	seen := map[string]int{}
	var collected []mvCard

	emptyStreak := 0
	for page := 1; ; page++ {
		resp, throttled, err := fetchDecksPage(client, &p.longest, expansion, page, pageSize)
		if err != nil {
			return nil, err
		}
		if throttled {
			p.throttled()
		}
		added := 0
		for _, c := range resp.Linked.Cards {
			if c.IsMaverick || (c.Expansion != 0 && c.Expansion != expansion) {
				continue
			}
			var isNew bool
			collected, isNew = mergeCard(collected, seen, c)
			if isNew {
				added++
			}
		}
		if added > 0 {
			emptyStreak = 0
		} else {
			emptyStreak++
		}

		decksSeen := min(page*pageSize, resp.Count)
		fmt.Printf("%s: %d\n", slug, len(collected))
		fmt.Printf(
			"  page %d: %d/%d decks scanned, %d cards so far\n",
			page, decksSeen, resp.Count, len(collected),
		)

		// Out of decks: the last page came back short.
		if page*pageSize >= resp.Count {
			break
		}
		if emptyStreak >= 100 {
			break
		}
		if emptyStreak >= 10 && !hasInteriorGap(collected) {
			break
		}
		time.Sleep(p.pace)
	}
	return collected, nil
}

// mergeCard folds one card into the collected set, keyed by collector number and
// title. A gigantic creature's two halves share both a number and a title, so the
// base half — which carries the power, armor, and card text — replaces an art half
// already collected under the same key, and an art half never displaces a base. It
// returns the (possibly grown) slice and whether a new card was added; a
// replacement is not a new card.
func mergeCard(collected []mvCard, seen map[string]int, c mvCard) ([]mvCard, bool) {
	key := c.CardNumber + "|" + c.CardTitle
	if idx, ok := seen[key]; ok {
		if isGiganticBase(c.CardType) && !isGiganticBase(collected[idx].CardType) {
			collected[idx] = c
		}
		return collected, false
	}
	seen[key] = len(collected)
	return append(collected, c), true
}

// isGiganticBase reports whether a Master Vault card type is the base half of a
// gigantic creature — the half that carries its power, armor, and card text (its
// art half prints "Gigantic Creature Art" with none of those).
func isGiganticBase(cardType string) bool {
	return strings.EqualFold(cardType, "gigantic creature base")
}

// hasCreatureStats reports whether a card type carries power and armor: a plain
// creature or a gigantic creature's base half.
func hasCreatureStats(cardType string) bool {
	return strings.EqualFold(cardType, "creature") || isGiganticBase(cardType)
}

// so a 429 is never fatal: it waits and keeps retrying until the page comes back.
// Each wait starts from *longest — the longest wait seen so far across all pages —
// and doubles up to a one-minute cap, so throttling that returns picks up where the
// backoff left off rather than restarting at two seconds. *longest is updated to
// the longest wait used. It reports whether the page was throttled so the caller
// can slow its pace.
func fetchDecksPage(
	client *http.Client,
	longest *time.Duration,
	expansion, page, pageSize int,
) (*mvResponse, bool, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	q.Set("links", "cards")
	q.Set("expansion", strconv.Itoa(expansion))
	q.Set("ordering", "date")
	reqURL := mvDecksURL + "?" + q.Encode()

	throttled := false
	const maxBackoff = 60 * time.Second
	wait := *longest
	sleepBackoff := func() {
		if wait > *longest {
			*longest = wait
		}
		time.Sleep(wait)
		wait *= 2
		if wait > maxBackoff {
			wait = maxBackoff
		}
	}
	for {
		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, throttled, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "vex-cardlookup/1.0")
		resp, err := client.Do(req)
		if err != nil {
			// A timeout or reset under load is transient — wait it out too.
			throttled = true
			fmt.Printf(
				"  page %d request failed (%v), waiting %s then retrying…\n",
				page, err, wait.Round(time.Second),
			)
			sleepBackoff()
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			_ = resp.Body.Close()
			throttled = true
			if ra := retryAfter(resp); ra > wait {
				wait = ra
			}
			fmt.Printf(
				"  page %d throttled (429), waiting %s then retrying…\n",
				page, wait.Round(time.Second),
			)
			sleepBackoff()
			continue
		}
		// A 5xx is a transient gateway blip (502/503/504 from the feed's proxy),
		// not a problem with this set — wait it out and retry like a 429 rather
		// than aborting the whole run.
		if resp.StatusCode >= 500 {
			_ = resp.Body.Close()
			throttled = true
			fmt.Printf(
				"  page %d gateway error (%d), waiting %s then retrying…\n",
				page, resp.StatusCode, wait.Round(time.Second),
			)
			sleepBackoff()
			continue
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
			_ = resp.Body.Close()
			return nil, throttled, fmt.Errorf(
				"decks feed status %d: %s",
				resp.StatusCode,
				string(bytes.TrimSpace(body)),
			)
		}
		var out mvResponse
		err = json.UnmarshalRead(resp.Body, &out)
		_ = resp.Body.Close()
		if err != nil {
			return nil, throttled, fmt.Errorf("decoding decks feed: %w", err)
		}
		return &out, throttled, nil
	}
}

// retryAfter reads the feed's Retry-After header (seconds), doubling it as the
// keyteki converter does, and falls back to two seconds when it is absent.
func retryAfter(resp *http.Response) time.Duration {
	if v := resp.Header.Get("Retry-After"); v != "" {
		if secs, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return time.Duration(secs*2) * time.Second
		}
	}
	return 2 * time.Second
}

// hasInteriorGap reports whether the collected numeric collector numbers skip a
// value below their maximum — a hole that means a card is still missing, so the
// importer keeps paging. Lettered reference numbers (S01, A21, P07) are ignored;
// they fill as they appear.
func hasInteriorGap(cards []mvCard) bool {
	maxNum := 0
	present := map[int]bool{}
	for _, c := range cards {
		n, err := strconv.Atoi(strings.TrimLeft(c.CardNumber, "0"))
		if err != nil {
			continue
		}
		present[n] = true
		if n > maxNum {
			maxNum = n
		}
	}
	for n := 1; n < maxNum; n++ {
		if !present[n] {
			return true
		}
	}
	return false
}

// --- transform ---------------------------------------------------------------

// catalogCard is one card as written to a source catalog. Field order is the
// order the JSON is emitted; empty optional fields are omitted. Printed carries
// the original card title only when it differs from the ASCII-folded name.
type catalogCard struct {
	Number  string   `json:"number"`
	Name    string   `json:"name"`
	Printed string   `json:"printed,omitempty"`
	House   string   `json:"house"`
	Type    string   `json:"type"`
	Rarity  string   `json:"rarity"`
	Traits  []string `json:"traits,omitempty"`
	Power   int      `json:"power,omitempty"`
	Armor   int      `json:"armor,omitempty"`
	Amber   int      `json:"amber,omitempty"`
	Text    string   `json:"text,omitempty"`
}

// transformCard folds one Master Vault card object into the catalog shape: it
// ASCII-folds the title and text, expands the amber/damage markup, normalizes the
// house and rarity, and classifies an anomaly. It fails loudly on an unknown
// markup token so a new icon is caught rather than silently dropped.
func transformCard(m mvCard) (catalogCard, error) {
	name := asciiFold(m.CardTitle)
	text, err := expandAndFoldText(m.CardText)
	if err != nil {
		return catalogCard{}, fmt.Errorf("card %q: %w", m.CardTitle, err)
	}

	house := strings.ToLower(strings.ReplaceAll(m.House, " ", ""))
	rarity := m.Rarity
	if strings.EqualFold(rarity, "FIXED") {
		rarity = "Special"
	}
	if m.IsAnomaly {
		house = "brobnar"
		rarity = "Special"
	}

	c := catalogCard{
		Number: normalizeCatalogNumber(m.CardNumber),
		Name:   name,
		House:  house,
		Type:   strings.ToLower(m.CardType),
		Rarity: rarity,
		Traits: foldTraits(m.Traits),
		Amber:  atoiSafe(string(m.Amber)),
		Text:   text,
	}
	if name != m.CardTitle {
		c.Printed = m.CardTitle
	}
	if hasCreatureStats(m.CardType) {
		c.Power = atoiSafe(string(m.Power))
		c.Armor = atoiSafe(string(m.Armor))
	}
	return c, nil
}

// normalizeCatalogNumber zero-pads an all-digit collector number to three digits
// (the printed form, "4" -> "004") and leaves a lettered reference number as is
// ("S01", "A21", "P07").
func normalizeCatalogNumber(s string) string {
	s = strings.TrimSpace(s)
	if n, err := strconv.Atoi(s); err == nil {
		return fmt.Sprintf("%03d", n)
	}
	return s
}

// foldTraits splits a Master Vault traits string on its bullet separator, folds
// each to lowercase ASCII, and drops the empties. A card with no traits yields a
// nil slice, which is omitted from the JSON.
func foldTraits(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, " • ")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToLower(strings.TrimSpace(asciiFold(p)))
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// atoiSafe parses a Master Vault numeric string, treating blank or non-numeric as
// zero (the stat is absent).
func atoiSafe(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// markupWords maps a Master Vault card-text icon token to the word the catalog
// spells it as. These are the four enhancement bonus icons (Aember, Capture,
// Damage, Draw), so an "Enhance" line reads "Enhance Aember Capture Damage Draw"
// (Mutagenesis Researcher). Only the tokens the feed is known to use are listed;
// expandMarkup fails on any other so a new icon surfaces instead of vanishing.
var markupWords = map[string]string{
	"A": "Aember",
	"C": "Capture",
	"D": "Damage",
	"R": "Draw",
}

// glyphWords maps each KeyForge icon-font Private Use Area glyph the feed embeds
// directly (rather than as an <A>-style token) to the word the catalog spells it
// as: the resource and enhancement bonus icons, the house enhancement icons
// introduced in Aember Skies, and Draconian Measures' power counter. Capture also
// appears as an F36F+F560 pair and the tide icon (U+F566) carries no printed word;
// both are handled in glyphWord, not here.
var glyphWords = map[rune]string{
	'\uF360': "Aember",
	'\uF361': "Damage",
	'\uF36E': "Draw",
	'\uF565': "Capture",
	'\uF372': "Discard",
	'\uF379': "Brobnar",
	'\uF37A': "Dis",
	'\uF37B': "Ekwidon",
	'\uF37C': "Geistoid",
	'\uF37D': "Logos",
	'\uF37E': "Mars",
	'\uF37F': "Skyborn",
	'\uF386': "Redemption",
	'\uF387': "Sanctum",
	'\uF388': "Saurian",
	'\uF389': "Shadows",
	'\uF38A': "Star Alliance",
	'\uF38B': "Untamed",
	'\uF390': "Unfathomable",
	'\uF391': "Ouboros",
	'\uF392': "+1 power counter",
}

// expandAndFoldText expands a card's markup, spells any bare-letter Enhance icon
// codes, ASCII-folds the result, strips the templated keyword reminder text, and
// tidies stray whitespace.
func expandAndFoldText(s string) (string, error) {
	expanded, err := expandMarkup(s)
	if err != nil {
		return "", err
	}
	return tidyText(stripKeywordReminders(asciiFold(expandEnhanceLetters(expanded)))), nil
}

// enhanceLetterRe matches an Enhance keyword followed by the run of bare letter
// icon codes some feed rows use in place of glyphs (A, D, R, and the P+T pair for
// Capture), up to the closing period.
var enhanceLetterRe = regexp.MustCompile(`Enhance ([ACDRPT]+)\.`)

// expandEnhanceLetters spells each bare-letter Enhance icon code: A Aember, D
// Damage, R Draw, and the P+T pair Capture. It only rewrites a run that directly
// follows Enhance and ends at a period, so ordinary prose is left untouched. An
// unexpected letter leaves the run as-is rather than mangling it.
func expandEnhanceLetters(s string) string {
	return enhanceLetterRe.ReplaceAllStringFunc(s, func(m string) string {
		codes := m[len("Enhance ") : len(m)-1]
		var words []string
		for i := 0; i < len(codes); i++ {
			switch codes[i] {
			case 'A':
				words = append(words, "Aember")
			case 'C':
				words = append(words, "Capture")
			case 'D':
				words = append(words, "Damage")
			case 'R':
				words = append(words, "Draw")
			case 'P':
				if i+1 < len(codes) && codes[i+1] == 'T' {
					words = append(words, "Capture")
					i++
					continue
				}
				return m
			default:
				return m
			}
		}
		return "Enhance " + strings.Join(words, " ") + "."
	})
}

// keywordReminderTexts is the exact parenthetical reminder text KeyForge prints
// after a keyword. Only these templated blurbs are stripped; any other
// parenthetical (a card-specific clarification such as "(rounding down the
// loss)") is a rules clarification, not a keyword reminder, and is left untouched.
var keywordReminderTexts = []string{
	"The first time this creature is attacked each turn, no damage is dealt.",
	"When you use this creature to fight, it is dealt no damage in return.",
	"When you use a creature with skirmish to fight, it is dealt no damage in return.",
	"This creature's neighbors cannot be attacked unless they have taunt.",
	"Any damage dealt by this creature's power during a fight destroys the damaged creature.",
	"You can only play this card before doing anything else this step.",
	"After you play this card, end this step.",
	"This creature can enter play anywhere in your battleline.",
	"These icons have already been added to cards in your deck.",
	"This card may be used as if it belonged to the active house.",
	"This card enters play under your opponent's control.",
}

// keywordReminderPatterns matches the reminder texts that carry a number
// (Assault N, Hazardous N, Splash-attack N), anchored so only the whole blurb
// qualifies. The trailing period is optional because the feed sometimes prints it
// outside the parenthesis.
var keywordReminderPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^Before this creature attacks, deal \d+ Damage to the attacked enemy\.$`),
	regexp.MustCompile(
		`^Before this creature is attacked, deal \d+ Damage to the attacking enemy\.$`,
	),
	regexp.MustCompile(
		`^When this creature attacks, also deal \d+ Damage to each of the attacked creature's neighbors\.?$`,
	),
}

// parenReminderRe captures a parenthetical along with any spaces or tabs before it
// and an optional period after it, so removing a keyword reminder leaves no
// dangling separator when the feed prints the sentence's period outside the
// parenthesis (Splash-attack's "…neighbors).").
var parenReminderRe = regexp.MustCompile(`[ \t]*\([^)]*\)\.?`)

// stripKeywordReminders removes each parenthetical whose text is a templated
// keyword reminder, leaving card-specific parentheticals in place.
func stripKeywordReminders(s string) string {
	return parenReminderRe.ReplaceAllStringFunc(s, func(m string) string {
		open := strings.IndexByte(m, '(')
		closeIdx := strings.LastIndexByte(m, ')')
		inner := strings.TrimSpace(m[open+1 : closeIdx])
		if isKeywordReminder(inner) {
			return ""
		}
		return m
	})
}

// isKeywordReminder reports whether inner is one of the templated keyword blurbs.
func isKeywordReminder(inner string) bool {
	if slices.Contains(keywordReminderTexts, inner) {
		return true
	}
	for _, re := range keywordReminderPatterns {
		if re.MatchString(inner) {
			return true
		}
	}
	return false
}

// multiSpaceRe matches a run of two or more spaces or tabs, which tidyText
// collapses to one so a dropped icon or a stray double space reads cleanly.
var multiSpaceRe = regexp.MustCompile(`[ \t]{2,}`)

// tidyText collapses each line's runs of spaces to one, trims the leading and
// trailing whitespace of every line, and trims the whole text.
func tidyText(s string) string {
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		lines[i] = strings.TrimSpace(multiSpaceRe.ReplaceAllString(ln, " "))
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// expandMarkup turns Master Vault card-text markup into plain words: the vertical
// tab line break becomes a newline, a carriage return (alone or as \r\n) becomes a
// newline too, each icon token (<A>, <D>) becomes its word, and each icon-font
// Private Use Area glyph the feed embeds directly (Aember, house enhancements, and
// the like) becomes its word too, with a space inserted where an icon abutted a
// letter or digit so "2<A>" reads "2 Aember". An unrecognized or unterminated
// token, or an unrecognized glyph, is an error.
func expandMarkup(s string) (string, error) {
	s = strings.ReplaceAll(s, "\u000b", "\n")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	var b strings.Builder
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		switch {
		case rs[i] == '<':
			j := i + 1
			for j < len(rs) && rs[j] != '>' {
				j++
			}
			if j >= len(rs) {
				return "", fmt.Errorf("unterminated markup token in %q", s)
			}
			tok := string(rs[i+1 : j])
			word, ok := markupWords[tok]
			if !ok {
				return "", fmt.Errorf("unknown markup token <%s> in %q", tok, s)
			}
			writeExpandedIcon(&b, word, j+1 < len(rs) && isAlnumASCII(rs[j+1]))
			i = j
		case isPrivateUse(rs[i]):
			word, extra, err := glyphWord(rs, i)
			if err != nil {
				return "", err
			}
			i += extra
			if word == "" { // tide icon: dropped
				continue
			}
			writeExpandedIcon(&b, word, i+1 < len(rs) && isAlnumASCII(rs[i+1]))
		default:
			b.WriteRune(rs[i])
		}
	}
	return b.String(), nil
}

// writeExpandedIcon appends an expanded icon word to b, inserting a space before
// it when it abuts a preceding letter or digit and after it when nextAlnum is set,
// so "2<A>" reads "2 Aember" and two adjacent icons read as two words.
func writeExpandedIcon(b *strings.Builder, word string, nextAlnum bool) {
	if out := b.String(); out != "" && isAlnumASCII(rune(out[len(out)-1])) {
		b.WriteByte(' ')
	}
	b.WriteString(word)
	if nextAlnum {
		b.WriteByte(' ')
	}
}

// isPrivateUse reports whether r is in the Basic Multilingual Plane Private Use
// Area, where KeyForge's icon font places its glyphs.
func isPrivateUse(r rune) bool {
	return r >= '\uE000' && r <= '\uF8FF'
}

// glyphWord returns the word for the icon glyph at rs[i] and how many extra runes
// it consumed (1 for the F36F+F560 Capture pair, 0 otherwise). The tide icon
// (U+F566) returns an empty word so it is dropped. An unrecognized glyph is an
// error so a new icon surfaces rather than vanishing.
func glyphWord(rs []rune, i int) (word string, extra int, err error) {
	switch rs[i] {
	case '\uF566': // tide: an ability marker with no printed word
		return "", 0, nil
	case '\uF36F': // Capture, encoded as an F36F+F560 pair
		if i+1 < len(rs) && rs[i+1] == '\uF560' {
			return "Capture", 1, nil
		}
		return "", 0, fmt.Errorf("unpaired capture glyph U+F36F in %q", string(rs))
	case '\uF560':
		return "", 0, fmt.Errorf("unpaired capture glyph U+F560 in %q", string(rs))
	}
	if w, ok := glyphWords[rs[i]]; ok {
		return w, 0, nil
	}
	return "", 0, fmt.Errorf("unknown icon glyph U+%04X in %q", rs[i], string(rs))
}

// isAlnumASCII reports whether r is a plain ASCII letter or digit.
func isAlnumASCII(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
}

// asciiFold rewrites a string into plain ASCII: curly quotes straighten, Æ becomes
// Ae, the multiplication sign becomes x, en/em dashes, the minus sign, and the
// non-breaking/figure hyphens become a hyphen, the non-breaking space becomes a
// space, zero-width joiners and the byte-order mark are dropped, and accented
// letters lose their diacritics (é -> e). English card text and names fold
// cleanly; anything left non-ASCII is passed through unchanged.
func asciiFold(s string) string {
	s = strings.NewReplacer(
		"\u2019", "'", // ’ right single quote
		"\u2018", "'", // ‘ left single quote
		"\u02bc", "'", // ʼ modifier letter apostrophe
		"\u201c", `"`, // “ left double quote
		"\u201d", `"`, // ” right double quote
		"\u00c6", "Ae", // Æ
		"\u00e6", "ae", // æ
		"\u00d7", "x", // × multiplication sign
		"\u2010", "-", // ‐ hyphen
		"\u2011", "-", // ‑ non-breaking hyphen (feed's "non‑Mars", "Nine‑Toes")
		"\u2013", "-", // – en dash
		"\u2014", "-", // — em dash
		"\u2212", "-", // − minus sign
		"\u00a0", " ", //   non-breaking space
		"\ufeff", "", // zero-width no-break space (BOM) the feed leaves at line ends
		"\u200b", "", // zero-width space
		"\u2026", "...", // … ellipsis
	).Replace(s)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	folded, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return folded
}

// --- output ------------------------------------------------------------------

// writeCatalog writes the cards as a source catalog: a two-space-indented JSON
// array with a trailing newline, matching the existing catalog files. HTML escaping
// is off so card text keeps its literal &, <, and > rather than \u escapes.
func writeCatalog(path string, cards []catalogCard) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := json.MarshalWrite(
		f,
		cards,
		jsontext.WithIndent("  "),
		jsontext.EscapeForHTML(false),
	); err != nil {
		_ = f.Close()
		return err
	}
	if _, err := f.WriteString("\n"); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
