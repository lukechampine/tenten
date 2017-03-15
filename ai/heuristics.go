package ai

var (
	popcounts [256]byte

	// only need 10 bits for these tables, but 16 eliminates the bounds check
	stride3Lookup [1 << 16]int
	stride5Lookup [1 << 16]int
	gapsLookup    [1 << 16]int
	dotsLookup    [1 << 16]int
	streaksLookup [1 << 16]int
)

func init() {
	for i := range popcounts {
		var n byte
		for x := i; x != 0; x >>= 1 {
			if x&1 != 0 {
				n++
			}
		}
		popcounts[i] = n
	}

	for u := range stride3Lookup[:1<<10] {
		for i := uint16(0); i < 10; i++ {
			l := 0
			for ; i < 10 && (u&(1<<i) == 0); i++ {
				l++
				if l >= 3 {
					stride3Lookup[u]++
				}
				if l >= 5 {
					stride5Lookup[u]++
				}
			}
		}
		for i := uint16(0); i < 10; i++ {
			if u&(1<<i) == 0 {
				gapsLookup[u]++
				for i < 10 && (u&(1<<i) == 0) {
					i++
				}
			}
		}
		for i := uint16(0); i < 10; i++ {
			d1 := i == 0 || u&(1<<(i-1)) != 0
			d2 := u&(1<<i) != 0
			d3 := i == 9 || u&(1<<(i+1)) != 0
			if d1 && !d2 && d3 {
				dotsLookup[u]++
			}
		}

		var i uint16
		var n int
		if u&1 == 0 {
			for i < 10 && (u&(1<<i) == 0) {
				i++
			}
			for i < 10 && (u&(1<<i) != 0) {
				i++
				n++
			}
		} else {
			for i < 10 && (u&(1<<i) != 0) {
				i++
				n++
			}
			for i < 10 && (u&(1<<i) == 0) {
				i++
			}
		}
		if i == 10 {
			streaksLookup[u] = n
		}
	}
	streaksLookup[0] = 10 // incentivize empty lines too
	gapsLookup[0] = 0
}

func popcount(x uint16) byte { return popcounts[byte(x)] + popcounts[byte(x>>8)] }

// number of rows and cols that do not contain any dots -- larger is better
func (h *hboard) emptyLines() float64 {
	var emptyLines int
	for _, row := range h.rows {
		if row == 0 {
			emptyLines++
		}
	}
	for _, col := range h.cols {
		if col == 0 {
			emptyLines++
		}
	}
	return float64(emptyLines)
}

// smallest rectangle that contains every dot -- smaller is better
func (h *hboard) smallestRect() float64 {
	var rx1, rx2, ry1, ry2 int
	for i, row := range h.rows {
		if row > 0 {
			rx1 = i
			break
		}
	}
	for i := 9; i >= 0; i-- {
		if h.rows[i] > 0 {
			rx2 = i + 1
			break
		}
	}
	for i, col := range h.cols {
		if col > 0 {
			ry1 = i
			break
		}
	}
	for i := 9; i >= 0; i-- {
		if h.cols[i] > 0 {
			ry2 = i + 1
			break
		}
	}
	rect := (rx2 - rx1) * (ry2 - ry1)
	return float64(rect)
}

// number of gaps -- smaller is better
func (h *hboard) gaps() float64 {
	var gaps int
	for i := range h.rows {
		gaps += gapsLookup[h.rows[i]]
		gaps += gapsLookup[h.cols[i]]
	}
	return float64(gaps)
}

// number of dots -- smaller is better
func (h *hboard) dots() float64 {
	var dots int
	for i := range h.rows {
		dots += dotsLookup[h.rows[i]]
		dots += dotsLookup[h.cols[i]]
	}
	return float64(dots)
}

// spaces for long pieces -- larger is better
func (h *hboard) lineSpaces() float64 {
	var longs int
	for i := range h.rows {
		longs += stride5Lookup[h.rows[i]]
		longs += stride5Lookup[h.cols[i]]
	}
	return float64(longs)
}

// spaces for large squares or Ls -- larger is better
func (h *hboard) squareSpaces() float64 {
	var squares int
	for i := range h.rows[2:] {
		squares += stride3Lookup[h.rows[i]|h.rows[i+1]|h.rows[i+2]]
	}
	return float64(squares)
}

// lines are empty or contain a single streak bordering an edge -- larger is better
func (h *hboard) streaks() float64 {
	var streaks int
	for i := range h.rows {
		streaks += streaksLookup[h.rows[i]]
		streaks += streaksLookup[h.cols[i]]
	}
	return float64(streaks)
}

// smallest number of rows and/or cols that cover every dot -- smaller is better
// NOTE: NP-complete
func (h *hboard) mincover() float64 {
	return 0
}
