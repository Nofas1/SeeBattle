package domain

func Constructor() *Field {
	m := make([][]int, Size)
	for i := range m {
        m[i] = make([]int, Size)
    }
	return &Field{Matrix: m}
}

type PlacerFunc func() <-chan PlaceRequest

func RandomPlacer() <-chan PlaceRequest {
    ch := make(chan PlaceRequest)
	go func() {
		defer close(ch)
		for i := 0; i < len(ShipSizes); i++ {
			feedback := make(chan bool)
			for {
				ch <- PlaceRequest{
					ShipSize: ShipSizes[i],
					Dir:   globalRand.Intn(4),
					Point: Pair{globalRand.Intn(10), globalRand.Intn(10)},
					Feedback: feedback,
				}
				if check := <-feedback; check {
					break
				}
			}
			close(feedback)
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
	if point.X < 0 || point.X >= Size || point.Y < 0 || point.Y >= Size {
		return false
	}
	for i := point.X - 1; i < point.X + 2; i++ {
		for j := point.Y - 1; j < point.Y + 2; j++ {
			if i < 0 || i >= Size || j < 0 || j >= Size {
				continue
			}
			if f.Matrix[i][j] == SHIP {
				return false
			}
		}
	}
	return true
}

func (f *Field) PlaceShip(ship, dir int, point Pair) bool {
	var cells []Pair
    myDir := Directions[dir]

	for i := 0; i < ship; i++ {
        p := Pair{point.X + myDir[0] * i, point.Y + myDir[1] * i}
        if !f.Validation(p) {
            return false
        }
        cells = append(cells, p)
    }

	if cells == nil {
        return false
    }
    for _, cell := range cells {
        f.Matrix[cell.X][cell.Y] = SHIP
    }
	return true
}

func (f *Field) BuildField(placer PlacerFunc, cancel <-chan struct{}) error {
	requests := placer()
	for cnt := 0; cnt < len(ShipSizes); {
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
	return nil
}

func (f *Field) Shoot(row, col int) ShotResult {
	if f.Matrix[row][col] == SHOOTED || f.Matrix[row][col] == MISSED {
		return Already
	}
	if f.Matrix[row][col] == SHIP {
		f.Matrix[row][col] = SHOOTED
		return Hit
	}
	f.Matrix[row][col] = MISSED
	return Miss
}
