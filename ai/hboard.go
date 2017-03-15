package ai

import "github.com/lukechampine/tenten/game"

const fullLine = (uint16(1) << 10) - 1

type hboard struct {
	origrows  [10]uint16
	origcols  [10]uint16
	rows      [10]uint16
	cols      [10]uint16
	piecerows []uint16
	piececols []uint16
	weights   [7]float64
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

	var gaps int
	var dots int
	var longs int
	var streaks int
	var emptyLines int
	for _, row := range h.rows {
		if row == 0 {
			emptyLines++
		}
		streaks += streaksLookup[row]
		longs += stride5Lookup[row]
		dots += dotsLookup[row]
		gaps += gapsLookup[row]
	}
	for _, col := range h.cols {
		if col == 0 {
			emptyLines++
		}
		streaks += streaksLookup[col]
		longs += stride5Lookup[col]
		dots += dotsLookup[col]
		gaps += gapsLookup[col]
	}

	var squares int
	for i := range h.rows[2:] {
		squares += stride3Lookup[h.rows[i]|h.rows[i+1]|h.rows[i+2]]
	}

	return 0 +
		h.weights[0]*float64(emptyLines) +
		h.weights[1]*float64(rect) +
		h.weights[2]*float64(gaps) +
		h.weights[3]*float64(dots) +
		h.weights[4]*float64(longs) +
		h.weights[5]*float64(squares) +
		h.weights[6]*float64(streaks)
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

	h.weights = [7]float64{0.4308687822582667, -0.05941428063493283, -0.5438121475790749, -0.1938080851557288, 0.17377008297367327, 0.7251280120639546, 0.033662995850984694}

	return h
}

func copyHBoard(dst, src *hboard) {
	dst.origrows, dst.origcols = src.rows, src.cols
}
