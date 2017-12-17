package ai

import (
	"math"
	"math/rand"
	"time"

	"github.com/lukechampine/tenten/game"
)

const monteC = 1.414 // uct tradeoff parameter; low = choose high value nodes, high = choose unexplored nodes

type Move struct {
	Piece game.Piece
	X, Y  int
}

func BestMoves(b *game.Board, bag [3]game.Piece, timeLimit time.Duration) ([3]Move, int) {
	perms := [][3]game.Piece{
		{bag[0], bag[1], bag[2]},
		{bag[0], bag[2], bag[1]},
		{bag[1], bag[0], bag[2]},
		{bag[1], bag[2], bag[0]},
		{bag[2], bag[0], bag[1]},
		{bag[2], bag[1], bag[0]},
	}
	trees := make([]*tree, len(perms))
	for i, perm := range perms {
		trees[i] = newTree(b, perm)
	}
	start := time.Now()
	for time.Since(start) < timeLimit {
		for _, tree := range trees {
			tree.expand()
		}
	}
	bestTree := trees[0]
	movesEvaluated := trees[0].root.N
	for _, tree := range trees[1:] {
		if tree.root.W > bestTree.root.W {
			bestTree = tree
		}
		movesEvaluated += tree.root.N
	}
	return bestTree.bestMoves(), movesEvaluated
}

type tree struct {
	origBoard   *game.Board
	expandBoard *game.Board
	root        *treeNode
	leaves      []*treeNode
}

func (t *tree) expand() {
	// locate node with best UCT, applying each move as we descend
	t.origBoard.Copy(t.expandBoard)
	toExpand := t.root
	for len(toExpand.leaves) > 0 {
		toExpand = toExpand.highestUCTChild()
		t.expandBoard.Place(toExpand.move.Piece, toExpand.move.X, toExpand.move.Y)
	}
	toExpand.expand(t.expandBoard)
}

func (t *tree) bestMoves() [3]Move {
	t1 := t.root.highestValueChild()
	t2 := t1.highestValueChild()
	t3 := t2.highestValueChild()
	return [3]Move{t1.move, t2.move, t3.move}
}

func newTree(b *game.Board, bag [3]game.Piece) *tree {
	t := &tree{
		origBoard:   b,
		expandBoard: new(game.Board),
		root:        &treeNode{leaves: make([]*treeNode, 0, 100)},
	}
	// first move
	for x0 := 0; x0 <= 10-bag[0].Width(); x0++ {
		for y0 := 0; y0 <= 10-bag[0].Height(); y0++ {
			b.Copy(t.expandBoard)
			if t.expandBoard.Place(bag[0], x0, y0) == 0 {
				continue
			}
			move0 := &treeNode{
				move:   Move{bag[0], x0, y0},
				parent: t.root,
				leaves: make([]*treeNode, 0, 100),
			}
			score := rollout(t.expandBoard)
			move0.propagate(score)
			t.root.leaves = append(t.root.leaves, move0)

			// second move
			for x1 := 0; x1 <= 10-bag[1].Width(); x1++ {
				for y1 := 0; y1 <= 10-bag[1].Height(); y1++ {
					b.Copy(t.expandBoard)
					t.expandBoard.Place(bag[0], x0, y0)
					if t.expandBoard.Place(bag[1], x1, y1) == 0 {
						continue
					}
					move1 := &treeNode{
						move:   Move{bag[1], x1, y1},
						parent: move0,
						leaves: make([]*treeNode, 0, 100),
					}
					score := rollout(t.expandBoard)
					move1.propagate(score)
					move0.leaves = append(move0.leaves, move1)

					// third move
					for x2 := 0; x2 <= 10-bag[2].Width(); x2++ {
						for y2 := 0; y2 <= 10-bag[2].Height(); y2++ {
							b.Copy(t.expandBoard)
							t.expandBoard.Place(bag[0], x0, y0)
							t.expandBoard.Place(bag[1], x1, y1)
							if t.expandBoard.Place(bag[2], x2, y2) == 0 {
								continue
							}
							move2 := &treeNode{
								move:   Move{bag[2], x2, y2},
								parent: move1,
								leaves: make([]*treeNode, 0, 100),
							}
							score := rollout(t.expandBoard)
							move2.propagate(score)
							move1.leaves = append(move1.leaves, move2)
						}
					}
				}
			}
		}
	}
	return t
}

type treeNode struct {
	move   Move
	N      int // number of rollouts
	W      int // accumulated value
	parent *treeNode
	leaves []*treeNode
}

func (t *treeNode) highestUCTChild() *treeNode {
	if len(t.leaves) == 0 {
		panic("no children")
	}
	bestUCT := math.Inf(-1)
	var bestChild *treeNode
	for _, l := range t.leaves {
		uct := (float64(l.W) / float64(l.N)) + monteC*math.Sqrt(math.Log(float64(t.N))/float64(l.N))
		if uct > bestUCT {
			bestChild = l
			bestUCT = uct
		}
	}
	return bestChild
}

func (t *treeNode) highestValueChild() *treeNode {
	if len(t.leaves) == 0 {
		panic("no children")
	}
	bestChild := t.leaves[0]
	for _, l := range t.leaves[1:] {
		if l.W > bestChild.W {
			bestChild = l
		}
	}
	return bestChild
}

var rng = rand.New(rand.NewSource(0))

func (t *treeNode) expand(b *game.Board) {
	if len(t.leaves) != 0 {
		panic("node already expanded")
	}
	rolloutBoard := new(game.Board)
	for _, p := range game.Pieces {
		for x := 0; x <= 10-p.Width(); x++ {
			for y := 0; y <= 10-p.Height(); y++ {
				b.Copy(rolloutBoard)
				if rolloutBoard.Place(p, x, y) == 0 {
					continue
				}
				m := &treeNode{
					move:   Move{p, x, y},
					parent: t,
					leaves: make([]*treeNode, 0, 100),
				}
				score := rollout(rolloutBoard)
				m.propagate(score)
				t.leaves = append(t.leaves, m)
			}
		}
	}
}

func rollout(b *game.Board) int {
	score := 0
	for {
		p := game.Pieces[rng.Intn(game.NumPieces)]
		for i := 0; i < 100; i++ {
			x := rng.Intn(11 - p.Width())
			y := rng.Intn(11 - p.Height())
			if b.Place(p, x, y) > 0 {
				break
			} else if i == 99 {
				// give up
				return score
			}
		}
		score++
	}
}

func (t *treeNode) propagate(score int) {
	t.W += score
	t.N++
	if t.parent != nil {
		t.parent.propagate(score)
	}
}
