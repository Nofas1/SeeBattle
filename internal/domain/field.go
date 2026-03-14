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
					Dir:   GlobalRand.Intn(4),
					Point: Pair{GlobalRand.Intn(10), GlobalRand.Intn(10)},
					Feedback: feedback,
				}
				if check := <-feedback; check {
					break
				}
			}
			close(feedback)
		}
	}()
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

func (f *Field) FillSunkArea(row, col int) {
    shipCells := []Pair{{X: row, Y: col}}
    for _, d := range Directions {
        newRow, newCol := row + d[0], col + d[1]
        for newRow >= 0 && newRow < Size && newCol >= 0 && newCol < Size {
            if f.Matrix[newRow][newCol] != SHOOTED {
                break
            }
            shipCells = append(shipCells, Pair{X: newRow, Y: newCol})
            newRow += d[0]
            newCol += d[1]
        }
    }

    for _, cell := range shipCells {
        for i := cell.X - 1; i <= cell.X + 1; i++ {
            for j := cell.Y - 1; j <= cell.Y + 1; j++ {
                if i >= 0 && i < Size && j >= 0 && j < Size {
                    if f.Matrix[i][j] == EMPTY {
                        f.Matrix[i][j] = FILL
                    }
                }
            }
        }
    }
}

func (f *Field) IsSunk(row, col int) bool {
	for _, d := range Directions {
		newRow, newCol := row + d[0], col + d[1]
		for newCol < Size && newCol >= 0 && newRow < Size && newRow >= 0 {
			if f.Matrix[newRow][newCol] == SHIP {
				return false
			}
			if f.Matrix[newRow][newCol] == EMPTY || f.Matrix[newRow][newCol] == MISSED {
				break
			}
			newRow += d[0]
			newCol += d[1]
		}
	}

	return true
}

func (f *Field) Shoot(row, col int) ShotResult {
	if f.Matrix[row][col] == SHOOTED || f.Matrix[row][col] == MISSED || f.Matrix[row][col] == FILL {
		return Already
	}
	if f.Matrix[row][col] == SHIP {
		f.Matrix[row][col] = SHOOTED
		if f.IsSunk(row, col) {
			f.FillSunkArea(row, col)
			return Sink
		}
		return Hit
	}
	f.Matrix[row][col] = MISSED
	return Miss
}

func (f *Field) UserShoot(row, col int) ShotResult {
	return f.Shoot(row, col)
}

func (f *Field) BotShoot() ShotResult {
	for {
		target := Pair{GlobalRand.Intn(10), GlobalRand.Intn(10)}
		row, col := target.Y, target.X
		if f.Matrix[row][col] == EMPTY || f.Matrix[row][col] == SHIP {
			return f.Shoot(row, col)
		}
	}
}
