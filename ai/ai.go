package ai

import "github.com/lukechampine/tenten/game"

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

func popcount(x uint16) (n byte) {
	return popcounts[byte(x)] + popcounts[byte(x>>8)]
}

const fullLine = (uint16(1) << 10) - 1

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

func (h *hboard) heuristic() int {
	// maximize empty lines
	var emptyLines int
	for i := range h.rows {
		if h.rows[i] == 0 {
			emptyLines++
		}
		if h.cols[i] == 0 {
			emptyLines++
		}
	}

	// maximize space for "dangerous" pieces
	var cl15 int
	for _, row := range h.rows {
		cl15 += stride5Lookup[row]
	}
	var cl51 int
	for _, col := range h.cols {
		cl51 += stride5Lookup[col]
	}
	var csq3 int
	for i := range h.rows[2:] {
		// bitwise | 3 rows, then count strides of 3
		csq3 += stride3Lookup[h.rows[i]|h.rows[i+1]|h.rows[i+2]]
	}

	// maximize contiguity
	var contiguous, disparate int
	for i := range h.rows[1:] {
		// bitwise ^ 2 rows -- 0 means contiguous, 1 means disparate
		ones := int(popcount(h.rows[i] ^ h.rows[i+1]))
		contiguous += 10 - ones
		disparate += ones
	}

	// apply weights
	return emptyLines*20 + contiguous*2 + disparate*-10 + cl15*20 + cl51*20 + csq3*50
}

func Heuristic(b *game.Board) int {
	return makeHBoard(b, 0).heuristic()
}

type Move struct {
	Piece game.Piece
	X, Y  int
}

func BestMoves(b *game.Board, bag [3]game.Piece) [3]Move {
	perms := [][3]game.Piece{
		{bag[0], bag[1], bag[2]},
		{bag[0], bag[2], bag[1]},
		{bag[1], bag[0], bag[2]},
		{bag[1], bag[2], bag[0]},
		{bag[2], bag[0], bag[1]},
		{bag[2], bag[1], bag[0]},
	}
	var scratch1, scratch2 game.Board
	bestPerm := perms[0]
	var bestX, bestY [3]int
	maxH := -1000000
	for _, perm := range perms {
		for x1 := 0; x1 <= 10-perm[0].Width(); x1++ {
			for y1 := 0; y1 <= 10-perm[0].Height(); y1++ {
				if !b.IsEmpty(x1, y1) {
					continue
				}
				b.Copy(&scratch1)
				if scratch1.Place(perm[0], x1, y1) <= 0 {
					continue
				}
				for x2 := 0; x2 <= 10-perm[1].Width(); x2++ {
					for y2 := 0; y2 <= 10-perm[1].Height(); y2++ {
						if !scratch1.IsEmpty(x2, y2) {
							continue
						}
						scratch1.Copy(&scratch2)
						if scratch2.Place(perm[1], x2, y2) <= 0 {
							continue
						}

						hboard := makeHBoard(&scratch2, perm[2])
						for x3 := 0; x3 <= 10-perm[2].Width(); x3++ {
							for y3 := 0; y3 < 10-perm[2].Height(); y3++ {
								if !hboard.place(perm[2], x3, y3) {
									continue
								} else if h := hboard.heuristic(); h > maxH {
									maxH = h
									bestPerm = perm
									bestX[0], bestX[1], bestX[2], bestY[0], bestY[1], bestY[2] = x1, x2, x3, y1, y2, y3
								}
							}
						}
					}
				}
			}
		}
	}
	return [3]Move{
		{bestPerm[0], bestX[0], bestY[0]},
		{bestPerm[1], bestX[1], bestY[1]},
		{bestPerm[2], bestX[2], bestY[2]},
	}
}
