package ai

import "github.com/lukechampine/tenten/game"

func Heuristic(b *game.Board) int {
	// maximize empty lines
	var emptyLines int
	for i := range b {
		// check row i
		clear := true
		for x := range b {
			if b[x][i] != game.Empty {
				clear = false
				break
			}
		}
		if clear {
			emptyLines++
		}
		// check column i
		clear = true
		for y := range b {
			if b[i][y] != game.Empty {
				clear = false
				break
			}
		}
		if clear {
			emptyLines++
		}
	}

	// maximize contiguity
	var contiguous, disparate int
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			if b[i][j] != game.Empty {
				if b[i-1][j] != game.Empty && b[i][j-1] != game.Empty {
					contiguous += 4
				} else if b[i-1][j] != game.Empty || b[i][j-1] != game.Empty {
					contiguous++
				} else {
					disparate++
				}
			}

		}
	}

	// maximize space for "dangerous" pieces
	cl15 := capacityLine1x5(b)
	cl51 := capacityLine5x1(b)
	csq3 := capacitySq3x3(b)

	// apply weights
	h := emptyLines*20 + contiguous*2 + disparate*-10 + cl15*20 + cl51*20 + csq3*50
	return h
}

func capacityLine1x5(b *game.Board) (n int) {
	for y := range b {
		for x := 0; x < 10; x++ {
			l := 0
			for ; x < 10 && b.IsEmpty(x, y); x++ {
				l++
				if l >= 5 {
					n++
				}
			}
		}
	}
	return
}

func capacityLine5x1(b *game.Board) (n int) {
	for x := range b {
		for y := 0; y < 10; y++ {
			l := 0
			for ; y < 10 && b.IsEmpty(x, y); y++ {
				l++
				if l >= 5 {
					n++
				}
			}
		}
	}
	return
}

func capacitySq3x3(b *game.Board) (n int) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if b.IsEmpty(x, y) && b.IsEmpty(x+1, y) && b.IsEmpty(x+2, y) &&
				b.IsEmpty(x, y+1) && b.IsEmpty(x+1, y+1) && b.IsEmpty(x+2, y+1) &&
				b.IsEmpty(x, y+2) && b.IsEmpty(x+1, y+2) && b.IsEmpty(x+2, y+2) {
				n++
			}
		}
	}
	return
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
	var scratch1, scratch2, scratch3 game.Board
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
						for x3 := 0; x3 <= 10-perm[2].Width(); x3++ {
							for y3 := 0; y3 < 10-perm[2].Height(); y3++ {
								if !scratch2.IsEmpty(x3, y3) {
									continue
								}
								scratch2.Copy(&scratch3)
								if scratch3.Place(perm[2], x3, y3) <= 0 {
									continue
								}
								if h := Heuristic(&scratch3); h > maxH {
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
