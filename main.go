package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
  rl.InitWindow(800, 450, "test")
  defer rl.CloseWindow()

  for !rl.WindowShouldClose() {

    rl.BeginDrawing()
    rl.ClearBackground(rl.RayWhite)
    rl.DrawText("WORKS!", 350, 200, 20, rl.Black)
    rl.EndDrawing()
  }
}
