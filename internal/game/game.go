package game

import (
	"sea_battle/internal/domain"
)

func Shot(field *domain.Field, row, col int) domain.ShotResult {
	return field.Shoot(row, col)
}