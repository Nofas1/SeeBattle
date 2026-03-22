package bot

import (
	"sea_battle/internal/domain"
)

type SimpleBot struct {
	field *domain.Field
}

func NewSimpleBot(field *domain.Field) *SimpleBot {
	return &SimpleBot{
		field: field,
	}
}

func (sb *SimpleBot) Shoot() domain.Pair {
	for {
		target := domain.Pair{X: domain.GlobalRand.Intn(10), Y: domain.GlobalRand.Intn(10)}
		row, col := target.X, target.Y
		if sb.field.Matrix[row][col] == domain.EMPTY || sb.field.Matrix[row][col] == domain.SHIP {
			return target
		}
	}
}

func (sb *SimpleBot) Getter() *domain.Field {
	return sb.field
}

func (sb *SimpleBot) SetResult(domain.ShotResult) {}
