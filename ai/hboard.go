package ai

import "github.com/lukechampine/tenten/game"

const fullLine = (uint16(1) << 10) - 1

var popcounts [256]byte
var stride3Lookup [1 << 10]int
var stride5Lookup [1 << 10]int

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

	for u := range stride3Lookup {
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
	origrows [10]uint16
	origcols [10]uint16
	rows     [10]uint16
	cols     [10]uint16
}

func (h *hboard) place(p game.Piece, x, y int) bool {
	h.rows, h.cols = h.origrows, h.origcols
	fullRows, fullCols := make([]byte, 0, 10), make([]byte, 0, 10)
	for _, d := range p.Dots() {
		// TODO: convert dots to row bits and use bits&row > 0 to check for collisions
		dx, dy := uint16(d.X+x), uint16(d.Y+y)
		if h.rows[dy]&(1<<dx) != 0 {
			return false
		}
		h.rows[dy] |= 1 << dx
		if h.rows[dy] == fullLine {
			fullRows = append(fullRows, byte(dy))
		}
		if h.cols[dx]&(1<<dy) != 0 {
			return false
		}
		h.cols[dx] |= 1 << dy
		if h.cols[dx] == fullLine {
			fullCols = append(fullCols, byte(dx))
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
		for x := range b {
			if b[x][i] != game.Empty {
				h.origcols[i] |= (1 << uint16(x))
			}
		}
		for y := range b {
			if b[i][y] != game.Empty {
				h.origrows[i] |= (1 << uint16(y))
			}
		}
	}
	h.rows, h.cols = h.origrows, h.origcols
	return h
}
