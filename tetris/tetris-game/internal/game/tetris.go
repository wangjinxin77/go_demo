package game

import (
    "fmt"
    "time"
)

type Tetris struct {
    board    *Board
    current  *Piece
    next     *Piece
    gameOver bool
}

func NewTetris() *Tetris {
    return &Tetris{
        board:    NewBoard(10, 20), // 设置宽度和高度
        current:  NewRandomPiece(),
        next:     NewRandomPiece(),
        gameOver: false,
    }
}

func (t *Tetris) Start() {
    for !t.gameOver {
        t.Update()
        t.Render()
        time.Sleep(500 * time.Millisecond) // 控制游戏速度
    }
    fmt.Println("Game Over!")
}

func (t *Tetris) Update() {
    if !t.board.CanPlace(t.current, t.current.Position) {
        t.board.PlacePiece(t.current)
        t.current = t.next
        t.next = NewRandomPiece()
        if !t.board.CanPlace(t.current, t.current.Position) {
            t.gameOver = true
        }
    } else {
        t.current.MoveDown()
    }
    lines := t.board.CheckFilledLines()
    t.board.ClearLines(lines)
}

func (t *Tetris) Render() {
    t.board.Render()
    fmt.Printf("Next Piece: %v\n", t.next.Shape)
}

func (t *Tetris) IsGameOver() bool {
    return t.gameOver
}