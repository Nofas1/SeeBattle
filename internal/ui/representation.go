package ui

import (
	"sea_battle/internal/domain"
	"sea_battle/internal/game"

	//   "sea_battle/internal/game"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	CELL = 50
	ROWS = HEIGHT / CELL
	COLS = WIDTH / CELL

	PADDING = 40
	GRID    = CELL * 10
	WIDTH   = GRID + PADDING*2
	HEIGHT  = GRID*2 + PADDING*3
)

const (
	userOffsetX = int32(PADDING)
	botOffsetX  = int32(PADDING*2 + GRID)
	offsetY     = int32(PADDING)
)

func DrawGrid(offsetX, offsetY int32, matrix [][]int, hideShips bool) {
	for i := int32(0); i < 10; i++ {
		for j := int32(0); j < 10; j++ {
			x := offsetX + j*CELL
			y := offsetY + i*CELL

			rect := rl.Rectangle{
				X:      float32(x),
				Y:      float32(y),
				Width:  CELL,
				Height: CELL,
			}

			color := rl.RayWhite
			switch matrix[i][j] {
			case domain.SHIP:
				if hideShips {
					color = rl.RayWhite
				} else {
					color = rl.Gray
				}
			case domain.SHOOTED:
				color = rl.Red
			case domain.MISSED:
				color = rl.Blue
			}

			rl.DrawRectangleRec(rect, color)
			rl.DrawRectangleLinesEx(rect, 1, rl.Black)
		}
	}

	if !hideShips {
		mp := rl.GetMousePosition()
		col := (int32(mp.X) - offsetX) / CELL
		row := (int32(mp.Y) - offsetY) / CELL
		if col > 0 && col < domain.Size && row > 0 && row < domain.Size {
			rl.DrawRectangle(offsetX+col*CELL, offsetY+row*CELL, CELL, CELL, rl.Fade(rl.Red, 0.5))
		}
	}
}

func Placer(userField *domain.Field, cancel <-chan struct{}, music rl.Music) {
    ship_index := 0
    dir := domain.Up
    input := make(chan domain.PlaceRequest)
    placed := make(chan bool)

	go userField.BuildField(domain.UserPlacer(input), cancel)

	for !rl.WindowShouldClose() && ship_index < len(domain.ShipSizes) {
		rl.UpdateMusicStream(music)
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		DrawGrid(userOffsetX, offsetY, userField.Matrix, false)

		select {
		case ok := <-placed:
			if ok {
				ship_index++
			}
		default:
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			dir = (dir + 1) % 4
		}

		mp := rl.GetMousePosition()
		col := (int32(mp.X) - userOffsetX) / CELL
		row := (int32(mp.Y) - offsetY) / CELL

		if ship_index < len(domain.ShipSizes) {
			pc, pr := col, row
			for i := 0; i < domain.ShipSizes[ship_index]; i++ {
				rl.DrawRectangle(userOffsetX+pc*CELL, offsetY+pr*CELL, CELL, CELL, rl.Fade(rl.Blue, 0.4))
				pr += int32(domain.Directions[dir][0])
				pc += int32(domain.Directions[dir][1])
			}
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			feedback := make(chan bool)
			req := domain.PlaceRequest{
				ShipSize: domain.ShipSizes[ship_index],
				Dir:      dir,
				Point:    domain.Pair{X: int(row), Y: int(col)},
				Feedback: feedback,
			}
			go func() {
				input <- req
				placed <- <-feedback
			}()
		}

		rl.EndDrawing()
	}
}

func Battle(userField, botField *domain.Field, music rl.Music) {
    // проверка убийства -> покрас всех на границе, 
    user_sunk := 0
    bot_sunk := 0
    turn := true
    hit_sound := rl.LoadSound("sounds/hit.wav")
	defer rl.UnloadSound(hit_sound)
	rl.SetSoundVolume(hit_sound, 0.1)
    for !rl.WindowShouldClose() {
        rl.UpdateMusicStream(music)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		DrawGrid(userOffsetX, offsetY, userField.Matrix, false)
		DrawGrid(botOffsetX, offsetY, botField.Matrix, true)

        if user_sunk == 10 {
            break
        }
        if bot_sunk == 10 {
            break
        }

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && turn == true {
			mp := rl.GetMousePosition()
			col := (int32(mp.X) - botOffsetX) / CELL
			row := (int32(mp.Y) - offsetY) / CELL
			if col >= 0 && col < domain.Size && row >= 0 && row < domain.Size {
				shotRes := game.UserShot(botField, int(row), int(col))
                rl.PlaySound(hit_sound)
                if shotRes == domain.Sink {
                    turn = true
                    bot_sunk++
                } else if shotRes == domain.Hit {
                    turn = true
                } else {
                    turn = false
                }
			}
		} else if turn == false {
            shotRes := game.BotShot(userField)
            if shotRes == domain.Sink {
                turn = false
                user_sunk++
            } else if shotRes == domain.Hit {
                turn = false
            } else {
                turn = true
            }
        }

		rl.EndDrawing()
	}
}

func Run(userField, botField *domain.Field) {
	rl.InitWindow(HEIGHT, WIDTH, "Sea Battle")
	defer rl.CloseWindow()

	cancel := make(chan struct{})
	defer close(cancel)
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	music := rl.LoadMusicStream("sounds/theme.mp3")
	defer rl.UnloadMusicStream(music)
	rl.SetMusicVolume(music, 0.1)
	rl.PlayMusicStream(music)
	Placer(userField, cancel, music)
	Battle(userField, botField, music)
}
