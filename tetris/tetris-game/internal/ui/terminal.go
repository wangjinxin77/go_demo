package ui

import (
    "fmt"
    "tetris-game/internal/game"
    "os"
    "os/exec"
    "runtime"
    "time"
)

type TerminalUI struct {
    board *game.Board
}

func NewTerminalUI(board *game.Board) *TerminalUI {
    return &TerminalUI{board: board}
}

func (ui *TerminalUI) ClearScreen() {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "linux":
        cmd = exec.Command("clear")
    case "windows":
        cmd = exec.Command("cmd", "/c", "cls")
    default:
        cmd = exec.Command("clear")
    }
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func (ui *TerminalUI) Render() {
    ui.ClearScreen()
    for _, row := range ui.board.Grid {
        for _, cell := range row {
            if cell {
                fmt.Print("[]")
            } else {
                fmt.Print("  ")
            }
        }
        fmt.Println()
    }
    fmt.Println("Controls: [A] Left, [D] Right, [S] Down, [W] Rotate, [Q] Quit")
}

func (ui *TerminalUI) DisplayMessage(message string) {
    fmt.Println(message)
    time.Sleep(2 * time.Second)
}