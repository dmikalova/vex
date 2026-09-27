package engine

import "testing"

// TestInterfaceTotality binds the interface classification to the interfaces the
// package's own source declares, so declaring a new interface without saying
// which catalog covers it — or why it needs none — fails the build (ADR 0018).
// This is the hole one level up from TestFamilyTotality: that test holds a
// family's members, this one holds the set of families, which is where Quantity,
// Gather and TopAct each entered.
func TestInterfaceTotality(t *testing.T) {
	declared := declaredInterfaces(t)
	classified := map[string]bool{}
	for _, row := range Interfaces() {
		if classified[row.Name] {
			t.Errorf("interface %s is classified twice", row.Name)
		}
		classified[row.Name] = true
	}
	for name := range declared {
		if !classified[name] {
			t.Errorf("interface %s is unclassified; name the catalog covering its "+
				"implementations, or the reason it is not a family (ADR 0018)", name)
		}
	}
	for name := range classified {
		if _, ok := declared[name]; !ok {
			t.Errorf("classification row %s is not an interface the package declares", name)
		}
	}
}

// TestInterfaceRowsClassified checks each row itself: it sets exactly one of a
// catalog or a reason for having none, and a catalog it names is a catalog that
// exists. The second half is what stops a row pointing at a family that was
// renamed or never written, which would read as covered while nothing enumerated
// its implementations.
func TestInterfaceRowsClassified(t *testing.T) {
	catalogs := map[string]bool{logEntryCatalog: true}
	for _, family := range Families() {
		catalogs[family.Name] = true
	}
	for _, row := range Interfaces() {
		switch {
		case row.Role.Catalog != "" && row.Role.NotFamily != "":
			t.Errorf("%s names both a catalog (%q) and a reason for none (%q)",
				row.Name, row.Role.Catalog, row.Role.NotFamily)
		case row.Role.Catalog == "" && row.Role.NotFamily == "":
			t.Errorf("%s is unclassified: name the catalog covering its "+
				"implementations, or the reason it is not a family", row.Name)
		case row.Role.Catalog != "" && !catalogs[row.Role.Catalog]:
			t.Errorf("%s names catalog %q, which no Families() entry carries",
				row.Name, row.Role.Catalog)
		}
	}
}

// TestDeclaredInterfacesReportsAScanFailure pins the error arm of the scan: a
// directory that cannot be read is reported rather than read as a package
// declaring no interfaces, which would pass the totality check while checking
// nothing.
func TestDeclaredInterfacesReportsAScanFailure(t *testing.T) {
	const missing = "no-such-directory"
	if _, err := DeclaredInterfaces(missing); err == nil {
		t.Errorf("DeclaredInterfaces(%q) returned no error", missing)
	}
}

// declaredInterfaces returns every interface the package's non-test source
// declares, failing the test when the scan finds none — a broken scan would
// otherwise report every row as an orphan.
func declaredInterfaces(t *testing.T) map[string]string {
	t.Helper()
	declared, err := DeclaredInterfaces(".")
	if err != nil {
		t.Fatalf("scanning for interfaces: %v", err)
	}
	if len(declared) == 0 {
		t.Fatalf("found no interfaces; the source scan is broken")
	}
	return declared
}
