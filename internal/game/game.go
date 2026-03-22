package game

import (
	"sea_battle/internal/domain"
)

type Bot interface{
	Shoot() domain.Pair
	Getter() *domain.Field
	SetResult(domain.ShotResult)
}

func Shoot(field *domain.Field, row, col int) domain.ShotResult {
	if field.Matrix[row][col] == domain.SHOOTED || field.Matrix[row][col] == domain.MISSED || field.Matrix[row][col] == domain.FILL {
		return domain.Already
	}
	if field.Matrix[row][col] == domain.SHIP {
		field.Matrix[row][col] = domain.SHOOTED
		if field.IsSunk(row, col) {
			field.FillSunkArea(row, col)
			return domain.Sink
		}
		return domain.Hit
	}
	field.Matrix[row][col] = domain.MISSED
	return domain.Miss
}

func UserShot(field *domain.Field, row, col int) domain.ShotResult {
	return Shoot(field, row, col)}

func BotShot(bot Bot) domain.ShotResult {
	shot := bot.Shoot()
	if shot.X == 11 && shot.Y == 11 {
		shot.X = 5
		shot.Y = 5
	}
	shotRes := Shoot(bot.Getter(), shot.X, shot.Y)
	bot.SetResult(shotRes)
	return shotRes
}
