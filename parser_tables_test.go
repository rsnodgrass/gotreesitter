package gotreesitter

import "testing"

func TestLookupActionIndexSmallUsesCompactDenseTokenRows(t *testing.T) {
	lang := &Language{
		TokenCount:         64,
		LargeStateCount:    1,
		SmallParseTableMap: []uint32{0},
		// groupCount=2
		// action 11 for token symbols 1..13
		// action 17 for nonterminal symbol 70
		SmallParseTable: []uint16{
			2,
			11, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
			17, 1, 70,
		},
	}

	smallTokenLookup := buildSmallTokenLookup(lang)
	p := &Parser{
		language:         lang,
		smallBase:        int(lang.LargeStateCount),
		smallLookup:      buildSmallLookup(lang, smallTokenLookup),
		smallTokenLookup: smallTokenLookup,
	}

	if got, want := p.lookupActionIndexSmall(1, 1), uint16(11); got != want {
		t.Fatalf("lookupActionIndexSmall token 1 = %d, want %d", got, want)
	}
	if got, want := p.lookupActionIndexSmall(1, 13), uint16(11); got != want {
		t.Fatalf("lookupActionIndexSmall token 13 = %d, want %d", got, want)
	}
	if got := p.lookupActionIndexSmall(1, 14); got != 0 {
		t.Fatalf("lookupActionIndexSmall missing token = %d, want 0", got)
	}
	if got, want := p.lookupActionIndexSmall(1, 70), uint16(17); got != want {
		t.Fatalf("lookupActionIndexSmall nonterminal = %d, want %d", got, want)
	}
	if len(p.smallTokenLookup) != 1 || len(p.smallTokenLookup[0]) != 14 {
		t.Fatalf("smallTokenLookup row missing or wrong size: %+v", p.smallTokenLookup)
	}
	if len(p.smallLookup) != 1 || len(p.smallLookup[0]) != 1 {
		t.Fatalf("smallLookup should retain only nonterminals for dense token rows: %+v", p.smallLookup)
	}

	seen := map[Symbol]uint16{}
	p.forEachActionIndexInState(1, func(sym Symbol, idx uint16) bool {
		seen[sym] = idx
		return true
	})
	if got, want := seen[1], uint16(11); got != want {
		t.Fatalf("forEachActionIndexInState token 1 = %d, want %d; seen=%v", got, want, seen)
	}
	if got, want := seen[13], uint16(11); got != want {
		t.Fatalf("forEachActionIndexInState token 13 = %d, want %d; seen=%v", got, want, seen)
	}
	if got, want := seen[70], uint16(17); got != want {
		t.Fatalf("forEachActionIndexInState nonterminal = %d, want %d; seen=%v", got, want, seen)
	}
	if _, ok := seen[14]; ok {
		t.Fatalf("forEachActionIndexInState included missing token 14: seen=%v", seen)
	}
}

func TestLookupActionIndexSmallUsesFullRowsForWideDenseTokenRanges(t *testing.T) {
	lang := &Language{
		TokenCount:         64,
		LargeStateCount:    1,
		SmallParseTableMap: []uint32{0},
		// groupCount=2
		// action 11 for token symbols 1..8 plus 63
		// action 17 for nonterminal symbol 70
		SmallParseTable: []uint16{
			2,
			11, 9, 1, 2, 3, 4, 5, 6, 7, 8, 63,
			17, 1, 70,
		},
	}

	smallTokenLookup := buildSmallTokenLookup(lang)
	p := &Parser{
		language:         lang,
		smallBase:        int(lang.LargeStateCount),
		smallLookup:      buildSmallLookup(lang, smallTokenLookup),
		smallTokenLookup: smallTokenLookup,
	}

	if got, want := p.lookupActionIndexSmall(1, 8), uint16(11); got != want {
		t.Fatalf("lookupActionIndexSmall token 8 = %d, want %d", got, want)
	}
	if got, want := p.lookupActionIndexSmall(1, 63), uint16(11); got != want {
		t.Fatalf("lookupActionIndexSmall token 63 = %d, want %d", got, want)
	}
	if got := p.lookupActionIndexSmall(1, 62); got != 0 {
		t.Fatalf("lookupActionIndexSmall missing high token = %d, want 0", got)
	}
	if got, want := p.lookupActionIndexSmall(1, 70), uint16(17); got != want {
		t.Fatalf("lookupActionIndexSmall nonterminal = %d, want %d", got, want)
	}
	if len(p.smallTokenLookup) != 1 || len(p.smallTokenLookup[0]) != int(lang.TokenCount) {
		t.Fatalf("smallTokenLookup should keep full rows for wide token ranges: %+v", p.smallTokenLookup)
	}
	if len(p.smallLookup) != 1 || len(p.smallLookup[0]) != 1 {
		t.Fatalf("smallLookup should retain only nonterminals for dense token rows: %+v", p.smallLookup)
	}
}
