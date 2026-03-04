package game

import (
	"sea_battle/internal/domain"
)

func UserShot(field *domain.Field, row, col int) domain.ShotResult {
	return field.UserShoot(row, col)
}

func BotShot(field *domain.Field) domain.ShotResult {
	return field.BotShoot()
}
