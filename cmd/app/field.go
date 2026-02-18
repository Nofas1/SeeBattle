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

var directions = [][]int{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

type PlaceRequest struct {
	ShipSize int
    Dir   int
    Point Pair
	Feedback chan bool
}

func Constructor() *Field {
	m := make([][]int, Size)
	for i := range m {
        m[i] = make([]int, Size)
    }
	return &Field{matrix: m}
}

type PlacerFunc func() <-chan PlaceRequest

func RandomPlacer() <-chan PlaceRequest {
    ch := make(chan PlaceRequest)
	go func() {
		defer close(ch)
		for i := 0; i < len(shipSizes); i++ {
			feedback := make(chan bool)
			for {
				ch <- PlaceRequest{
					ShipSize: shipSizes[i],
					Dir:   globalRand.Intn(4),
					Point: Pair{globalRand.Intn(10), globalRand.Intn(10)},
					Feedback: feedback,
				}
				if check := <-feedback; check {
					break
				}
			}
		}
	}()

	// cansel channel if does not stop randomizing(no option for ship placing)
    
    return ch
}

func UserPlacer(input <-chan PlaceRequest) PlacerFunc {
    return func() <-chan PlaceRequest {
        return input
    }
}

func (f *Field) Validation(point Pair) bool {
	for i := point.x - 1; i < point.x + 2; i++ {
		for j := point.y - 1; j < point.y + 2; j++ {
			if i < 0 || i >= Size || j < 0 || j >= Size {
				continue
			}
			if f.matrix[i][j] == SHIP {
				return false
			}
		}
	}
	return true
}

func (f *Field) PlaceShip(ship, dir int, point Pair) bool {
	var cells []Pair
    myDir := directions[dir]

	for i := 0; i < ship; i++ {
        p := Pair{point.x + myDir[0] * i, point.y + myDir[1] * i}
        if !f.Validation(p) {
            return false
        }
        cells = append(cells, p)
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
	requests := placer()
	for cnt := 0; cnt < len(shipSizes); {
		select {
        case req, ok := <-requests:
			if !ok {
				return nil
			}
            if f.PlaceShip(req.ShipSize, req.Dir, req.Point) {
				req.Feedback <- true
                cnt++
            } else {
				req.Feedback <- false
			}
        case <-cancel:
            return nil
        }
	}

	return f
}
