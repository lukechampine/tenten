package game

import (
	"math/rand"
	"time"
)

type Color int

const (
	Empty Color = iota
	Red
	Pink
	Orange
	Yellow
	Green
	Teal
	Cyan
	Blue
	Purple
)

type Dot struct {
	X, Y int
}

type Piece int

const (
	invalid Piece = iota

	Line1x2 // vertical line, length 2
	Line1x3 // vertical line, length 3
	Line1x4 // vertical line, length 4
	Line1x5 // vertical line, length 5
	Line2x1 // horizontal line, length 2
	Line3x1 // horizontal line, length 3
	Line4x1 // horizontal line, length 4
	Line5x1 // horizontal line, length 5
	Ltr2x2  // small l, missing top right
	Ltl2x2  // small l, missing top left
	Lbr2x2  // small l, missing bottom right
	Lbl2x2  // small l, missing bottom left
	Ltr3x3  // big l, missing top right
	Ltl3x3  // big l, missing top left
	Lbr3x3  // big l, missing bottom right
	Lbl3x3  // big l, missing bottom left
	Sq1x1   // 1x1 square
	Sq2x2   // 2x2 square
	Sq3x3   // 3x3 square

	NumPieces = Sq3x3
)

func (p Piece) Color() Color { return pieceColors[p] }
func (p Piece) Dots() []Dot  { return pieceDots[p] }

var pieceColors = [NumPieces + 1]Color{
	invalid: Empty,
	Line1x2: Yellow,
	Line1x3: Orange,
	Line1x4: Pink,
	Line1x5: Red,
	Line2x1: Yellow,
	Line3x1: Orange,
	Line4x1: Pink,
	Line5x1: Red,
	Ltr2x2:  Teal,
	Ltl2x2:  Teal,
	Lbr2x2:  Teal,
	Lbl2x2:  Teal,
	Ltr3x3:  Blue,
	Ltl3x3:  Blue,
	Lbr3x3:  Blue,
	Lbl3x3:  Blue,
	Sq1x1:   Purple,
	Sq2x2:   Green,
	Sq3x3:   Cyan,
}

var pieceDots = [NumPieces + 1][]Dot{
	invalid: {},
	Line1x2: {
		{0, 0},
		{0, 1}},
	Line1x3: {
		{0, 0},
		{0, 1},
		{0, 2}},
	Line1x4: {
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3}},
	Line1x5: {
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
		{0, 4}},
	Line2x1: {
		{0, 0}, {1, 0}},
	Line3x1: {
		{0, 0}, {1, 0}, {2, 0}},
	Line4x1: {
		{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	Line5x1: {
		{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}},
	Ltr2x2: {
		{0, 0},
		{0, 1}, {1, 1}},
	Ltl2x2: {
		/*   */ {1, 0},
		{0, 1}, {1, 1}},
	Lbr2x2: {
		{0, 0}, {1, 0},
		{0, 1}},
	Lbl2x2: {
		{0, 0}, {1, 0},
		/*   */ {1, 1}},
	Ltr3x3: {
		{0, 0},
		{0, 1},
		{0, 2}, {1, 2}, {2, 2}},
	Ltl3x3: {
		/*           */ {2, 0},
		/*           */ {2, 1},
		{0, 2}, {1, 2}, {2, 2}},
	Lbr3x3: {
		{0, 0}, {1, 0}, {2, 0},
		{0, 1},
		{0, 2}},
	Lbl3x3: {
		{0, 0}, {1, 0}, {2, 0},
		/*           */ {2, 1},
		/*           */ {2, 2}},
	Sq1x1: {
		{0, 0}},
	Sq2x2: {
		{0, 0}, {1, 0},
		{0, 1}, {1, 1}},
	Sq3x3: {
		{0, 0}, {1, 0}, {2, 0},
		{0, 1}, {1, 1}, {2, 1},
		{0, 2}, {1, 2}, {2, 2}},
}

var Bags [NumPieces * NumPieces * NumPieces][3]Piece

func init() {
	pieces := []Piece{Line1x2, Line1x3, Line1x4, Line1x5, Line2x1, Line3x1, Line4x1, Line5x1, Ltr2x2, Ltl2x2, Lbr2x2, Lbl2x2, Ltr3x3, Ltl3x3, Lbr3x3, Lbl3x3, Sq1x1, Sq2x2, Sq3x3}
	i := 0
	for _, p1 := range pieces {
		for _, p2 := range pieces {
			for _, p3 := range pieces {
				Bags[i] = [3]Piece{p1, p2, p3}
				i++
			}
		}
	}
}

// A Board represents a tenten game board.
type Board [10][10]Color

func (b *Board) IsEmpty(x, y int) bool { return b[y][x] == Empty }
func (b *Board) set(x, y int, c Color) { b[y][x] = c }

// Copy copies b into dst.
func (b *Board) Copy(dst *Board) { copy(dst[:], b[:]) }

func (b *Board) clearLines(newdots []Dot) int {
	var checkRow, checkCol [10]bool
	for _, d := range newdots {
		checkCol[d.X] = true
		checkRow[d.Y] = true
	}

	var rows, cols []int
	for i := range b {
		if !checkRow[i] {
			// check row i
			clear := true
			for x := range b {
				if b[x][i] == Empty {
					clear = false
					break
				}
			}
			if clear {
				rows = append(rows, i)
			}
		}
		if !checkCol[i] {
			// check column i
			clear := true
			for y := range b {
				if b[i][y] == Empty {
					clear = false
					break
				}
			}
			if clear {
				cols = append(cols, i)
			}
		}
	}
	for _, r := range rows {
		for x := range b {
			b[x][r] = Empty
		}
	}
	for _, c := range cols {
		for y := range b {
			b[c][y] = Empty
		}
	}
	lines := len(rows) + len(cols)
	return 5 * lines * (lines + 1) // 1 -> 10, 2 -> 30, 3 -> 60, etc.
}

// Place places p at location (x,y), updating b as necessary. It returns the
// point value of the move. If the move is invalid, Place returns 0 and b is
// not affected.
func (b *Board) Place(p Piece, x, y int) int {
	if x < 0 || y < 0 || x >= 10 || y >= 10 {
		return 0
	}

	dots := p.Dots()
	for _, d := range dots {
		if d.X+x >= 10 || d.Y+y >= 10 || !b.IsEmpty(d.X+x, d.Y+y) {
			return 0
		}
	}
	for _, d := range dots {
		b.set(d.X+x, d.Y+y, p.Color())
	}
	// point value is added dots + cleared dots
	return len(dots) + b.clearLines(dots)
}

type Game struct {
	b     *Board
	score int
	rnd   *rand.Rand
}

func (g Game) Score() int    { return g.score }
func (g Game) Board() *Board { return g.b }

// Place places p at location (x,y) on the Board, updating it as necessary. If
// the move is invalid, Place returns false and the Board is not affected.
func (g *Game) Place(p Piece, x, y int) bool {
	points := g.b.Place(p, x, y)
	g.score += points
	return points > 0
}

func (g Game) NextBag() [3]Piece {
	return Bags[g.rnd.Intn(len(Bags))]
}

func New() Game {
	return Game{
		b:   new(Board),
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}
