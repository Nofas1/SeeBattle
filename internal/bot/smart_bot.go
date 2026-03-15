package bot

import (
	"fmt"
	"sea_battle/internal/domain"
)

type SmartBot struct {
	field *domain.Field
	targets []domain.Pair
	last_shot domain.Pair
	dir []int
	first_shot domain.Pair
	try_to_sink bool
}

func NewSmartBot(field *domain.Field) *SmartBot {
	return &SmartBot{
		field: field,
		targets: make([]domain.Pair, 0),
		last_shot: domain.Pair{},
	}
}

func (sb *SmartBot) reset() {
	sb.targets = make([]domain.Pair, 0)
	sb.dir = nil
	sb.first_shot = domain.Pair{}
	sb.try_to_sink = false
}

func (sb *SmartBot) Shoot() domain.Pair {
	if len(sb.targets) == 0 && sb.try_to_sink && sb.dir != nil {
        dx, dy := sb.dir[0], sb.dir[1]
        sb.addIfValid(domain.Pair{
            X: sb.first_shot.X - dx,
            Y: sb.first_shot.Y - dy,
        })
        sb.dir = nil
    }
	if len(sb.targets) > 0 {
		sb.last_shot = sb.targets[0]
		sb.targets = sb.targets[1:]
		return sb.last_shot
	} else {
		for {
			target := domain.Pair{X: domain.GlobalRand.Intn(10), Y: domain.GlobalRand.Intn(10)}
			row, col := target.Y, target.X
			if sb.field.Matrix[row][col] == domain.EMPTY || sb.field.Matrix[row][col] == domain.SHIP {
				return target
			}
		}
	}
}

func (sb *SmartBot) Getter() *domain.Field {
	return sb.field
}

func (sb *SmartBot) addIfValid(p domain.Pair) {
	fmt.Printf("addIfValid: p=(%d,%d)\n", p.X, p.Y)
	if p.X < 0 || p.X >= domain.Size || p.Y < 0 || p.Y >= domain.Size {
		return
	}
	cell := sb.field.Matrix[p.Y][p.X]
	if cell != domain.EMPTY && cell != domain.SHIP {
		return
	}
	for _, t := range sb.targets {
		if t == p {
			return
		}
	}
	sb.targets = append(sb.targets, p)
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func (sb *SmartBot) SetResult(shotRes domain.ShotResult) {
	fmt.Printf("SetResult: shot=(%d,%d) result=%d try_to_sink=%v targets=%v\n", sb.last_shot.X, sb.last_shot.Y, shotRes, sb.try_to_sink, sb.targets)
	if shotRes == domain.Already {
		return
	}
	if shotRes == domain.Sink {
		sb.reset()
		return
	}
	if shotRes == domain.Hit {
		if !sb.try_to_sink {
			sb.try_to_sink = true
			sb.first_shot = sb.last_shot
			for _, d := range domain.Directions {
				sb.addIfValid(domain.Pair{
					X: sb.last_shot.X + d[1],
					Y: sb.last_shot.Y + d[0],
				})
			}
		} else {
			dx := sign(sb.last_shot.X - sb.first_shot.X)
			dy := sign(sb.last_shot.Y - sb.first_shot.Y)
			sb.dir = []int{dx, dy}
			sb.targets = nil
			
			sb.addIfValid(domain.Pair{
				X: sb.last_shot.X + dx,
				Y: sb.last_shot.Y + dy,
			})
			
			sb.addIfValid(domain.Pair{
				X: sb.first_shot.X - dx,
				Y: sb.first_shot.Y - dy,
			})
		}
		return
	}
	if shotRes == domain.Miss {
		if sb.try_to_sink && sb.dir != nil && len(sb.targets) == 0 {
			dx, dy := sb.dir[0], sb.dir[1]
			sb.addIfValid(domain.Pair{
				X: sb.first_shot.X - dx,
				Y: sb.first_shot.Y - dy,
			})
		}
		return
	}
}
