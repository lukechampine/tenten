package ai

import "github.com/lukechampine/tenten/game"

const fullLine = (uint16(1) << 10) - 1

var popcounts [256]byte
var stride3Lookup [1 << 16]int // only need 10 bits, but 16 eliminates the bounds check
var stride5Lookup [1 << 16]int // only need 10 bits, but 16 eliminates the bounds check

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
}

func (h *hboard) place(p game.Piece, x, y int) bool {
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

func (h *hboard) heuristic() int {
	var emptyLines, cl15, cl51, csq3, contiguous, disparate int
	for i, row := range h.rows {
		// maximize empty lines
		if row == 0 {
			emptyLines++
		}

		// maximize space for "dangerous" pieces
		cl15 += stride5Lookup[row]

		// maximize space for "dangerous" pieces
		if i >= 2 {
			// bitwise | 3 rows, then count strides of 3
			csq3 += stride3Lookup[h.rows[i]|h.rows[i-1]|h.rows[i-2]]
		}

		// maximize contiguity
		if i >= 1 {
			// bitwise ^ 2 rows -- 0 means contiguous, 1 means disparate
			ones := int(popcount(h.rows[i] ^ h.rows[i-1]))
			contiguous += 10 - ones
			disparate += ones
		}
	}

	for _, col := range h.cols {
		// maximize empty lines
		if col == 0 {
			emptyLines++
		}
		// maximize space for "dangerous" pieces
		cl51 += stride5Lookup[col]
	}

	// apply weights
	return emptyLines*20 + contiguous*2 + disparate*-10 + cl15*20 + cl51*20 + csq3*50
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

	return h
}
