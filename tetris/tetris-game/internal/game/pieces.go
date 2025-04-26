package game

import (
    "math/rand"
)

type Piece struct {
    Shape    [][]int
    Color    string
    Position [2]int
}

var Pieces = []Piece{
    {
        Shape: [][]int{
            {1, 1, 1, 1},
        },
        Color: "cyan",
    },
    {
        Shape: [][]int{
            {1, 1, 1},
            {0, 1, 0},
        },
        Color: "yellow",
    },
    {
        Shape: [][]int{
            {1, 1},
            {1, 1},
        },
        Color: "purple",
    },
    {
        Shape: [][]int{
            {0, 1, 1},
            {1, 1, 0},
        },
        Color: "green",
    },
    {
        Shape: [][]int{
            {1, 1, 0},
            {0, 1, 1},
        },
        Color: "red",
    },
    {
        Shape: [][]int{
            {0, 1, 0},
            {1, 1, 1},
        },
        Color: "blue",
    },
    {
        Shape: [][]int{
            {1, 1, 1},
            {1, 0, 0},
        },
        Color: "orange",
    },
}

func NewRandomPiece() *Piece {
    p := Pieces[rand.Intn(len(Pieces))]
    return &Piece{
        Shape:    p.Shape,
        Color:    p.Color,
        Position: [2]int{3, 0}, // 初始位置
    }
}

func (p *Piece) Rotate() {
    newShape := make([][]int, len(p.Shape[0]))
    for i := range newShape {
        newShape[i] = make([]int, len(p.Shape))
    }
    for i := range p.Shape {
        for j := range p.Shape[i] {
            newShape[j][len(p.Shape)-1-i] = p.Shape[i][j]
        }
    }
    p.Shape = newShape
}

func (p *Piece) MoveDown() {
    p.Position[1]++
}