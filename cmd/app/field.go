package app

import (
	"time"
	"math/rand"
)

const (
	Size = 10
)

type Cell int

const (
	EMPTY = iota
	SHIP
	SHOOTED
	MISSED
	FILL
)

type Pair struct {
	x int
	y int
}

type Field struct {
	matrix [][]int
}

var shipSizes = []int{4, 3, 3, 2, 2, 2, 1, 1, 1, 1}
var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type PlaceRequest struct {
    Dir   int
    Point Pair
}

func Constructor() *Field {
	return &Field{matrix: make([][]int, Size)}
}

type PlacerFunc func(ship int) <-chan PlaceRequest

func RandomPlacer(ship int) <-chan PlaceRequest {
    ch := make(chan PlaceRequest, 1)
    ch <- PlaceRequest{
        Dir:   globalRand.Intn(4),
        Point: Pair{globalRand.Intn(10), globalRand.Intn(10)},
    }
    return ch
}

func UserPlacer(input <-chan PlaceRequest) PlacerFunc {
    return func(ship int) <-chan PlaceRequest {
        return input
    }
}

func (f *Field) Validation(point Pair) bool {
	for i := point.x - 1; i < point.x + 2; i++ {
		for j := point.y - 1; j < point.y + 2; j++ {
			if i < 0 || i > 9 || j < 0 || j > 9 {
				continue
			}
			if f.matrix[i][j] == SHIP {
				return false
			}
		}
	}
	return true
}

func (f *Field) placeUp(ship int, point Pair) []Pair {
    cells := make([]Pair, 0, ship)
    for i := 0; i < ship; i++ {
        p := Pair{point.x - i, point.y}
        if !f.Validation(p) {
            return nil
        }
        cells = append(cells, p)
    }
    return cells
}

func (f *Field) placeDown(ship int, point Pair) []Pair {
	cells := make([]Pair, 0, ship)
    for i := 0; i < ship; i++ {
        p := Pair{point.x + i, point.y}
        if !f.Validation(p) {
            return nil
        }
        cells = append(cells, p)
    }
    return cells
}
func (f *Field) placeLeft(ship int, point Pair) []Pair {
	cells := make([]Pair, 0, ship)
    for i := 0; i < ship; i++ {
        p := Pair{point.x, point.y - i}
        if !f.Validation(p) {
            return nil
        }
        cells = append(cells, p)
    }
    return cells
}
func (f *Field) placeRight(ship int, point Pair) []Pair {
	cells := make([]Pair, 0, ship)
    for i := 0; i < ship; i++ {
        p := Pair{point.x, point.y + i}
        if !f.Validation(p) {
            return nil
        }
        cells = append(cells, p)
    }
    return cells
}

func (f *Field) PlaceShip(ship, dir int, point Pair) bool {
	var cells []Pair
    switch dir {
	case 0:
        cells = f.placeUp(ship, point)
    case 1:
        cells = f.placeRight(ship, point)
    case 2:
        cells = f.placeDown(ship, point)
    default:
        cells = f.placeLeft(ship, point)
    }

	if cells == nil {
        return false
    }
    for _, cell := range cells {
        f.matrix[cell.x][cell.y] = SHIP
    }
	return true
}

func (f *Field) BuildField(placer PlacerFunc, cancel <-chan struct{}) *Field {
	for i := range f.matrix {
        f.matrix[i] = make([]int, Size)
    }

	for cnt := 0; cnt < len(shipSizes); {
		ship := shipSizes[cnt]
		requests := placer(ship)
		select {
        case req := <-requests:
            if f.PlaceShip(ship, req.Dir, req.Point) {
                cnt++
            }
        case <-cancel:
            return nil
        }
	}

	return f
}
