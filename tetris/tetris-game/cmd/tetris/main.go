package main

import (
    "tetris-game/internal/game"
)

func main() {
    tetris := game.NewTetris()
    tetris.Start()
}