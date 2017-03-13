package ai

import "github.com/lukechampine/tenten/game"

const fullLine = (uint16(1) << 10) - 1

var popcounts [256]byte

// only need 10 bits for these tables, but 16 eliminates the bounds check
var stride3Lookup [1 << 16]int
var stride5Lookup [1 << 16]int
var groupsLookup [1 << 16]int

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
			if u&(1<<i) != 0 {
				groupsLookup[u]++
				for i < 10 && (u&(1<<i) != 0) {
					i++
				}
			}
		}
	}
}

func popcount(x uint16) byte { return popcounts[byte(x)] + popcounts[byte(x>>8)] }

type hboard struct {
	origrows  [10]uint16
	origcols  [10]uint16
	rows      [10]uint16
	cols      [10]uint16
	piecerows []uint16
	piececols []uint16
	weights   [5]float64
}

func (h *hboard) place(x, y int) bool {
	for i, prow := range h.piecerows {
		if (prow<<uint16(x))&h.origrows[y+i] != 0 {
			return false
		}
	}
	h.rows, h.cols = h.origrows, h.origcols
	fullRows, fullCols := make([]byte, 0, 10), make([]byte, 0, 10)
	for i, prow := range h.piecerows {
		dy := byte(y + i)
		h.rows[dy] |= prow << uint16(x)
		if h.rows[dy] == fullLine {
			fullRows = append(fullRows, dy)
		}
	}
	for i, pcol := range h.piececols {
		dx := byte(x + i)
		h.cols[dx] |= pcol << uint16(y)
		if h.cols[dx] == fullLine {
			fullCols = append(fullCols, dx)
		}
	}

	for _, r := range fullRows {
		h.rows[r] = 0
		for i := range h.cols {
			h.cols[i] &^= 1 << uint16(r)
		}
	}
	for _, c := range fullCols {
		h.cols[c] = 0
		for i := range h.rows {
			h.rows[i] &^= 1 << uint16(c)
		}
	}
	return true
}

func (h *hboard) heuristic() float64 {
	var emptyLines, groups, cl15, cl51, csq3 int
	var prev, prev2 uint16
	for i, row := range h.rows {
		// maximize empty lines
		if row == 0 {
			emptyLines++
		}

		// maximize contiguity
		groups += groupsLookup[row]

		// maximize space for "dangerous" pieces
		cl15 += stride5Lookup[row]

		// maximize space for "dangerous" pieces
		if i >= 2 {
			// bitwise | 3 rows, then count strides of 3
			csq3 += stride3Lookup[row|prev|prev2]
		}

		prev2, prev = prev, row
	}

	for _, col := range h.cols {
		// maximize empty lines
		if col == 0 {
			emptyLines++
		}

		// maximize contiguity
		groups += groupsLookup[col]

		// maximize space for "dangerous" pieces
		cl51 += stride5Lookup[col]
	}

	// apply weights
	return 0 +
		h.weights[0]*float64(emptyLines)/20 +
		h.weights[1]*float64(100-groups)/100 +
		h.weights[2]*float64(cl15)/60 +
		h.weights[3]*float64(cl51)/60 +
		h.weights[4]*float64(csq3)/64
}

func makeHBoard(b *game.Board, p game.Piece) *hboard {
	h := new(hboard)
	for i := range b {
		for y := range b {
			if b[y][i] != game.Empty {
				h.origcols[i] |= (1 << uint16(y))
			}
		}
		for x := range b {
			if b[i][x] != game.Empty {
				h.origrows[i] |= (1 << uint16(x))
			}
		}
	}
	h.rows, h.cols = h.origrows, h.origcols

	h.piecerows = make([]uint16, p.Height())
	for _, d := range p.Dots() {
		h.piecerows[d.Y] |= (1 << uint16(d.X))
	}
	h.piececols = make([]uint16, p.Width())
	for _, d := range p.Dots() {
		h.piececols[d.X] |= (1 << uint16(d.Y))
	}

	h.weights = [5]float64{0.0, 0.5, 0.125, 0.125, 0.25}

	return h
}

func copyHBoard(dst, src *hboard) {
	dst.origrows, dst.origcols = src.rows, src.cols
}
