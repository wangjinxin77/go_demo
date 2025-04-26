package game

import (
    "fmt"
)

type Board struct {
    width  int
    height int
    Grid   [][]int
}

// NewBoard initializes a new game board with the specified width and height.
func NewBoard(width, height int) *Board {
    grid := make([][]int, height)
    for i := range grid {
        grid[i] = make([]int, width)
    }
    return &Board{
        width:  width,
        height: height,
        Grid:   grid,
    }
}

// CanPlace checks if a piece can be placed at the given position.
func (b *Board) CanPlace(piece *Piece, position [2]int) bool {
    for i, row := range piece.Shape {
        for j, cell := range row {
            if cell == 1 {
                x, y := position[0]+j, position[1]+i
                if x < 0 || x >= b.width || y < 0 || y >= b.height || b.Grid[y][x] == 1 {
                    return false
                }
            }
        }
    }
    return true
}

// PlacePiece places a piece on the board.
func (b *Board) PlacePiece(piece *Piece) {
    for i, row := range piece.Shape {
        for j, cell := range row {
            if cell == 1 {
                x, y := piece.Position[0]+j, piece.Position[1]+i
                b.Grid[y][x] = 1
            }
        }
    }
}

// CheckFilledLines checks for filled lines and returns their indices.
func (b *Board) CheckFilledLines() []int {
    filledLines := []int{}
    for i, row := range b.Grid {
        isFilled := true
        for _, cell := range row {
            if cell == 0 {
                isFilled = false
                break
            }
        }
        if isFilled {
            filledLines = append(filledLines, i)
        }
    }
    return filledLines
}

// ClearLines clears the specified lines and shifts the rows above down.
func (b *Board) ClearLines(lines []int) {
    for _, line := range lines {
        for i := line; i > 0; i-- {
            b.Grid[i] = b.Grid[i-1]
        }
        b.Grid[0] = make([]int, b.width) // Clear the top row
    }
}

// Render displays the current state of the board in the terminal.
func (b *Board) Render() {
    for _, row := range b.Grid {
        for _, cell := range row {
            if cell == 0 {
                fmt.Print(". ")
            } else {
                fmt.Print("# ")
            }
        }
        fmt.Println()
    }
}