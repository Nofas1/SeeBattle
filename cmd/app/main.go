package main

import (
  "sea_battle/internal/domain"
  "sea_battle/internal/ui"
)

func main() {
  cancel := make(chan struct{})
  defer close(cancel)

  botField := domain.Constructor()
  userField := domain.Constructor()

  go botField.BuildField(domain.RandomPlacer, cancel)

  ui.Run(userField, botField)
}
