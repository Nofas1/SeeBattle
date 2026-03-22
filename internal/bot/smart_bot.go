package bot

import (
	// "fmt"
	"sea_battle/internal/domain"
)

type SmartBot struct {
	field       *domain.Field
	targets     []domain.Pair
	last_shot   domain.Pair
	dir         []int
	first_shot  domain.Pair
	try_to_sink bool
}

func NewSmartBot(field *domain.Field) *SmartBot {
	return &SmartBot{
		field:     field,
		targets:   make([]domain.Pair, 0),
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
			row, col := target.X, target.Y
			if sb.field.Matrix[row][col] == domain.EMPTY || sb.field.Matrix[row][col] == domain.SHIP {
				sb.last_shot = target
				return target
			}
		}
	}
}

func (sb *SmartBot) Getter() *domain.Field {
	return sb.field
}

func (sb *SmartBot) addIfValid(p domain.Pair) {
    row, col := p.X, p.Y
    if row < 0 || row >= domain.Size || col < 0 || col >= domain.Size {
        return
    }
    if sb.field.Matrix[row][col] == domain.EMPTY || sb.field.Matrix[row][col] == domain.SHIP {
        sb.targets = append(sb.targets, p)
    }
}

func (sb *SmartBot) SetResult(shotRes domain.ShotResult) {
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
                    X: sb.last_shot.X + d[0],
                    Y: sb.last_shot.Y + d[1],
                })
            }
        } else {
            dx := sb.last_shot.X - sb.first_shot.X
            if dx != 0 {
                sb.dir = []int{1, 0}
            } else {
                sb.dir = []int{0, 1}
            }
            sb.addIfValid(domain.Pair{
                X: sb.last_shot.X + sb.dir[0],
                Y: sb.last_shot.Y + sb.dir[1],
            })
            filtered := sb.targets[:0]
            for _, t := range sb.targets {
                tdx := t.X - sb.first_shot.X
                tdy := t.Y - sb.first_shot.Y
                if (sb.dir[0] != 0 && tdx != 0) || (sb.dir[1] != 0 && tdy != 0) {
                    filtered = append(filtered, t)
                }
            }
            sb.targets = filtered
        }
        return
    }
    if shotRes == domain.Miss {
        if sb.try_to_sink && sb.dir != nil && len(sb.targets) == 0 {
            sb.addIfValid(domain.Pair{
                X: sb.first_shot.X - sb.dir[0],
                Y: sb.first_shot.Y - sb.dir[1],
            })
        }
    }
}
